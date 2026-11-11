package multirail

import (
	"context"
	"crypto/sha256"
	"fmt"
	"sync"
	"time"

	"github.com/shopspring/decimal"
)

type WiseConfig struct {
	APIToken  string
	ProfileID string
	BaseURL   string
}

type WiseAdapter struct {
	mu              sync.RWMutex
	cfg             WiseConfig
	transactions    map[string]*RailTransaction
	idempotencyKeys map[string]*RailTransaction
}

var (
	_ RailReader  = (*WiseAdapter)(nil)
	_ RailMutator = (*WiseAdapter)(nil)
)

func NewWiseAdapter(cfg WiseConfig) *WiseAdapter {
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.wise.com/v1"
	}
	return &WiseAdapter{
		cfg:             cfg,
		transactions:    make(map[string]*RailTransaction),
		idempotencyKeys: make(map[string]*RailTransaction),
	}
}

func (w *WiseAdapter) InitiatePayment(ctx context.Context, auth HumanAuthorization, amount decimal.Decimal, currency, reference string) (*RailTransaction, error) {
	if auth.OperatorID == "" {
		return nil, fmt.Errorf("%w: wise payout initiation requires human authorization", ErrUnauthorizedMutation)
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	// Idempotency check on reference
	if existing, found := w.idempotencyKeys[reference]; found {
		if !existing.Amount.Equal(amount) || existing.Currency != currency {
			return nil, ErrIdempotencyViolation
		}
		return existing, nil
	}

	transferID := fmt.Sprintf("wise_%d", time.Now().UnixNano()%1000000)
	now := time.Now().UTC()

	// Note on UETR: UETR is ONLY assigned when returned by SWIFT gpi or provider webhook.
	// In pending transfer stage, BankClearRef remains empty to prevent fabricated evidence.
	h := sha256.Sum256([]byte(transferID + amount.String() + currency))

	tx := &RailTransaction{
		Rail:           RailWise,
		TransactionID:  transferID,
		Amount:         amount,
		Currency:       currency,
		Status:         "processing",
		BankClearRef:   "", // Empty until genuine bank receipt is received
		OccurredAt:     now,
		RawPayloadHash: fmt.Sprintf("%x", h),
		IdempotencyKey: reference,
	}

	w.transactions[transferID] = tx
	w.idempotencyKeys[reference] = tx
	return tx, nil
}

func (w *WiseAdapter) CapturePayment(ctx context.Context, auth HumanAuthorization, txID string, amount decimal.Decimal) (*RailTransaction, error) {
	// Wise wire transfers are push payments; capture is not applicable
	return nil, fmt.Errorf("capture not supported on Wise push payment rail")
}

func (w *WiseAdapter) IssueRefund(ctx context.Context, auth HumanAuthorization, txID string, amount decimal.Decimal, reason string) (*RailTransaction, error) {
	if auth.OperatorID == "" {
		return nil, fmt.Errorf("%w: wise reversal requires operator authorization", ErrUnauthorizedMutation)
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	tx, exists := w.transactions[txID]
	if !exists {
		return nil, fmt.Errorf("wise transfer %s not found", txID)
	}

	tx.Status = "cancelled_refunded"
	return tx, nil
}

func (w *WiseAdapter) GetTransaction(ctx context.Context, txID string) (*RailTransaction, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	tx, exists := w.transactions[txID]
	if !exists {
		return nil, fmt.Errorf("wise transfer %s not found", txID)
	}
	return tx, nil
}

func (w *WiseAdapter) GetDispute(ctx context.Context, disputeID string) (*DisputeRecord, error) {
	return nil, fmt.Errorf("disputes not supported directly on wire transfer rail")
}

// IngestBankClearance updates the transaction with a verified SWIFT gpi UETR reference
func (w *WiseAdapter) IngestBankClearance(txID string, verifiedUETR string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	tx, exists := w.transactions[txID]
	if !exists {
		return fmt.Errorf("wise transfer %s not found", txID)
	}

	tx.BankClearRef = verifiedUETR
	tx.Status = "outgoing_payment_sent"
	return nil
}
