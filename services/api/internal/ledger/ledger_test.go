package ledger_test

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/paybridge/paybridge-ai/services/api/internal/ledger"
)

func TestLedger_BalancedCaptureJournal(t *testing.T) {
	eng := ledger.NewEngine()
	ctx := context.Background()

	gross := decimal.NewFromFloat(1250.00)
	fee := decimal.NewFromFloat(37.50)

	tx, err := eng.RecordCaptureJournal(ctx, uuid.New(), uuid.New(), gross, fee, "USD", "INV-2026-0881")
	if err != nil {
		t.Fatalf("balanced journal failed: %v", err)
	}

	if len(tx.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(tx.Entries))
	}
}

func TestLedger_RejectsUnbalancedTransaction(t *testing.T) {
	eng := ledger.NewEngine()
	ctx := context.Background()

	entries := []ledger.Entry{
		{
			ID:        uuid.New(),
			AccountID: "acct_a",
			Direction: ledger.Debit,
			Amount:    decimal.NewFromFloat(100.00),
			Currency:  "USD",
		},
		{
			ID:        uuid.New(),
			AccountID: "acct_b",
			Direction: ledger.Credit,
			Amount:    decimal.NewFromFloat(90.00), // Off by $10
			Currency:  "USD",
		},
	}

	_, err := eng.PostTransaction(ctx, uuid.New(), uuid.New(), "REF-1", "Unbalanced Test", entries)
	if err == nil {
		t.Fatal("CRITICAL INVARIANT VIOLATION: Unbalanced ledger transaction was permitted")
	}
}

func TestLedger_IdempotentDuplicatePosting(t *testing.T) {
	eng := ledger.NewEngine()
	ctx := context.Background()

	tenantID := uuid.New()
	caseID := uuid.New()
	gross := decimal.NewFromFloat(500.00)
	fee := decimal.NewFromFloat(15.00)
	ref := "TX-PAY-IDEM-001"

	tx1, err := eng.RecordCaptureJournal(ctx, tenantID, caseID, gross, fee, "USD", ref)
	if err != nil {
		t.Fatalf("initial posting failed: %v", err)
	}

	// Post again with identical reference
	tx2, err := eng.RecordCaptureJournal(ctx, tenantID, caseID, gross, fee, "USD", ref)
	if err != nil {
		t.Fatalf("idempotent posting failed: %v", err)
	}

	if tx1.ID != tx2.ID {
		t.Fatalf("expected identical transaction ID on idempotent posting, got %s vs %s", tx1.ID, tx2.ID)
	}
}

func TestLedger_SettlementAndRefundJournals(t *testing.T) {
	eng := ledger.NewEngine()
	ctx := context.Background()

	tenantID := uuid.New()
	caseID := uuid.New()

	// 1. Settlement journal
	payout := decimal.NewFromFloat(485.00)
	settleTx, err := eng.RecordSettlementJournal(ctx, tenantID, caseID, payout, "USD", "SETTLE-001")
	if err != nil {
		t.Fatalf("settlement journal failed: %v", err)
	}
	if len(settleTx.Entries) != 2 {
		t.Errorf("expected 2 settlement legs, got %d", len(settleTx.Entries))
	}

	// 2. Refund journal
	refundGross := decimal.NewFromFloat(100.00)
	refundFee := decimal.NewFromFloat(3.00)
	refundTx, err := eng.RecordRefundJournal(ctx, tenantID, caseID, refundGross, refundFee, "USD", "REFUND-001")
	if err != nil {
		t.Fatalf("refund journal failed: %v", err)
	}
	if len(refundTx.Entries) != 3 {
		t.Errorf("expected 3 refund legs, got %d", len(refundTx.Entries))
	}
}

func TestLedger_TenantIsolation(t *testing.T) {
	eng := ledger.NewEngine()
	ctx := context.Background()

	tenantA := uuid.New()
	tenantB := uuid.New()

	tx, err := eng.RecordCaptureJournal(ctx, tenantA, uuid.New(), decimal.NewFromInt(100), decimal.NewFromInt(3), "USD", "REF-A")
	if err != nil {
		t.Fatalf("posting failed: %v", err)
	}

	// Tenant A can retrieve
	retrieved, err := eng.GetTransaction(tenantA, tx.ID)
	if err != nil || retrieved == nil {
		t.Fatalf("tenant A failed to retrieve own transaction: %v", err)
	}

	// Tenant B is rejected
	_, err = eng.GetTransaction(tenantB, tx.ID)
	if err == nil {
		t.Fatal("CRITICAL INVARIANT VIOLATION: Cross-tenant ledger query succeeded")
	}
}

func TestLedger_ConcurrentPosting(t *testing.T) {
	eng := ledger.NewEngine()
	ctx := context.Background()

	const numWorkers = 16
	const postingsPerWorker = 100

	var wg sync.WaitGroup
	tenantID := uuid.New()

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for i := 0; i < postingsPerWorker; i++ {
				ref := fmt.Sprintf("CONCURRENT-%d-%d", workerID, i)
				gross := decimal.NewFromFloat(float64(100 + i))
				fee := decimal.NewFromFloat(2.50)

				_, err := eng.RecordCaptureJournal(ctx, tenantID, uuid.New(), gross, fee, "USD", ref)
				if err != nil {
					t.Errorf("worker %d posting %d failed: %v", workerID, i, err)
					return
				}
			}
		}(w)
	}

	wg.Wait()
}
