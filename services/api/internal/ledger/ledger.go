package ledger

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var (
	ErrUnbalancedTransaction = errors.New("ledger invariant violation: transaction sum of debits and credits is not zero")
	ErrEmptyTransaction      = errors.New("ledger: transaction must contain at least two entries")
	ErrDuplicatePosting      = errors.New("ledger: transaction with this reference has already been posted")
	ErrImmutableJournal      = errors.New("ledger: posted journal records cannot be modified or deleted")
	ErrTenantIsolation       = errors.New("ledger: tenant cross-access violation")
)

type Direction string

const (
	Debit  Direction = "DEBIT"
	Credit Direction = "CREDIT"
)

type Entry struct {
	ID        uuid.UUID       `json:"id"`
	AccountID string          `json:"account_id"`
	Direction Direction       `json:"direction"`
	Amount    decimal.Decimal `json:"amount"`
	Currency  string          `json:"currency"`
}

type Transaction struct {
	ID          uuid.UUID `json:"id"`
	TenantID    uuid.UUID `json:"tenant_id"`
	CaseID      uuid.UUID `json:"case_id"`
	Reference   string    `json:"reference"`
	Description string    `json:"description"`
	Entries     []Entry   `json:"entries"`
	OccurredAt  time.Time `json:"occurred_at"`
}

type LedgerEngine struct {
	mu           sync.RWMutex
	transactions map[uuid.UUID]*Transaction
	byReference  map[string]uuid.UUID // (tenant_id:reference) -> transaction_id (idempotency)
}

func NewEngine() *LedgerEngine {
	return &LedgerEngine{
		transactions: make(map[uuid.UUID]*Transaction),
		byReference:  make(map[string]uuid.UUID),
	}
}

// PostTransaction validates the zero-sum invariant: SUM(DEBITS) - SUM(CREDITS) == 0
// and commits an immutable, idempotent balanced journal entry
func (le *LedgerEngine) PostTransaction(
	ctx context.Context,
	tenantID, caseID uuid.UUID,
	reference, description string,
	entries []Entry,
) (*Transaction, error) {
	if len(entries) < 2 {
		return nil, ErrEmptyTransaction
	}

	// 1. Invariant Validation: Zero-Sum Balance
	netBalance := decimal.Zero
	for _, e := range entries {
		if e.Amount.LessThanOrEqual(decimal.Zero) {
			return nil, fmt.Errorf("ledger entry amount must be positive, got %s", e.Amount.String())
		}
		if e.Direction == Debit {
			netBalance = netBalance.Add(e.Amount)
		} else if e.Direction == Credit {
			netBalance = netBalance.Sub(e.Amount)
		} else {
			return nil, fmt.Errorf("invalid entry direction: %s", e.Direction)
		}
	}

	if !netBalance.IsZero() {
		return nil, fmt.Errorf("%w: net balance is %s (must be exactly zero)", ErrUnbalancedTransaction, netBalance.String())
	}

	le.mu.Lock()
	defer le.mu.Unlock()

	// 2. Idempotency Check: reject duplicate journal posting with same reference for tenant
	dedupeKey := fmt.Sprintf("%s:%s", tenantID.String(), reference)
	if existingID, exists := le.byReference[dedupeKey]; exists {
		return le.transactions[existingID], nil // idempotent return
	}

	txID := uuid.New()
	tx := &Transaction{
		ID:          txID,
		TenantID:    tenantID,
		CaseID:      caseID,
		Reference:   reference,
		Description: description,
		Entries:     entries,
		OccurredAt:  time.Now().UTC(),
	}

	le.transactions[txID] = tx
	le.byReference[dedupeKey] = txID

	return tx, nil
}

// RecordCaptureJournal: Balanced 3-leg entry for successful payment capture
// Leg 1: DEBIT  Gateway Clearing Asset (Gross Amount)
// Leg 2: CREDIT Merchant Payable Liability (Net Amount)
// Leg 3: CREDIT Platform Processing Revenue (Fee Amount)
// Net: Gross - (Net + Fee) == 0
func (le *LedgerEngine) RecordCaptureJournal(
	ctx context.Context,
	tenantID, caseID uuid.UUID,
	gross, fee decimal.Decimal,
	currency, reference string,
) (*Transaction, error) {
	netPayable := gross.Sub(fee)
	if netPayable.IsNegative() {
		return nil, fmt.Errorf("invalid fee: fee %s exceeds gross %s", fee, gross)
	}

	entries := []Entry{
		{
			ID:        uuid.New(),
			AccountID: "asset_gateway_clearing",
			Direction: Debit,
			Amount:    gross,
			Currency:  currency,
		},
		{
			ID:        uuid.New(),
			AccountID: "liability_merchant_payable",
			Direction: Credit,
			Amount:    netPayable,
			Currency:  currency,
		},
		{
			ID:        uuid.New(),
			AccountID: "revenue_platform_fees",
			Direction: Credit,
			Amount:    fee,
			Currency:  currency,
		},
	}

	return le.PostTransaction(ctx, tenantID, caseID, reference, "Payment Capture Settlement", entries)
}

// RecordCaptureJournalHelper provides backward compatibility with initial prototype
func (le *LedgerEngine) RecordCaptureJournalHelper(
	ctx context.Context,
	tenantID, caseID uuid.UUID,
	gross, fee decimal.Decimal,
	currency, reference string,
) (*Transaction, error) {
	return le.RecordCaptureJournal(ctx, tenantID, caseID, gross, fee, currency, reference)
}

// RecordSettlementJournal: Balanced 2-leg entry for payout / bank deposit
// Leg 1: DEBIT  Merchant Payable Liability (Net Payout)
// Leg 2: CREDIT Bank Operating Cash Asset (Net Payout)
func (le *LedgerEngine) RecordSettlementJournal(
	ctx context.Context,
	tenantID, caseID uuid.UUID,
	payoutAmount decimal.Decimal,
	currency, reference string,
) (*Transaction, error) {
	entries := []Entry{
		{
			ID:        uuid.New(),
			AccountID: "liability_merchant_payable",
			Direction: Debit,
			Amount:    payoutAmount,
			Currency:  currency,
		},
		{
			ID:        uuid.New(),
			AccountID: "asset_bank_operating_cash",
			Direction: Credit,
			Amount:    payoutAmount,
			Currency:  currency,
		},
	}

	return le.PostTransaction(ctx, tenantID, caseID, reference, "Merchant Bank Payout Settlement", entries)
}

// RecordRefundJournal: Balanced entry for full or partial refund
// Leg 1: DEBIT  Merchant Payable Liability (Net Refund)
// Leg 2: DEBIT  Platform Fee Refund Expense (Fee Reversal)
// Leg 3: CREDIT Gateway Clearing Asset (Total Refunded)
func (le *LedgerEngine) RecordRefundJournal(
	ctx context.Context,
	tenantID, caseID uuid.UUID,
	refundGross, refundFee decimal.Decimal,
	currency, reference string,
) (*Transaction, error) {
	merchantDebit := refundGross.Sub(refundFee)

	entries := []Entry{
		{
			ID:        uuid.New(),
			AccountID: "liability_merchant_payable",
			Direction: Debit,
			Amount:    merchantDebit,
			Currency:  currency,
		},
		{
			ID:        uuid.New(),
			AccountID: "expense_platform_fee_reversal",
			Direction: Debit,
			Amount:    refundFee,
			Currency:  currency,
		},
		{
			ID:        uuid.New(),
			AccountID: "asset_gateway_clearing",
			Direction: Credit,
			Amount:    refundGross,
			Currency:  currency,
		},
	}

	return le.PostTransaction(ctx, tenantID, caseID, reference, "Customer Refund Reversal", entries)
}

// GetTransaction retrieves a posted transaction, enforcing tenant isolation
func (le *LedgerEngine) GetTransaction(tenantID, txID uuid.UUID) (*Transaction, error) {
	le.mu.RLock()
	defer le.mu.RUnlock()

	tx, exists := le.transactions[txID]
	if !exists {
		return nil, fmt.Errorf("ledger transaction %s not found", txID)
	}

	if tx.TenantID != tenantID {
		return nil, ErrTenantIsolation
	}

	return tx, nil
}
