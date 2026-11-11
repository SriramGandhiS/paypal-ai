package multirail_test

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/paybridge/paybridge-ai/services/api/internal/multirail"
)

func TestStripeAdapter_PaymentAndDispute(t *testing.T) {
	adapter := multirail.NewStripeAdapter(multirail.StripeConfig{
		SecretKey:     "sk_test_placeholder",
		WebhookSecret: "whsec_test_secret_12345",
	})
	ctx := context.Background()

	auth := multirail.HumanAuthorization{
		OperatorID: "ops-user-001",
		Reason:     "Order #ORD-101 fulfillment checkout",
		Timestamp:  time.Now().UTC(),
	}

	amount := decimal.NewFromFloat(150.00)
	tx, err := adapter.InitiatePayment(ctx, auth, amount, "USD", "ORD-101")
	if err != nil {
		t.Fatalf("stripe payment initiation failed: %v", err)
	}

	if tx.Status != "requires_capture" {
		t.Errorf("expected status requires_capture, got %s", tx.Status)
	}

	// Capture payment
	tx, err = adapter.CapturePayment(ctx, auth, tx.TransactionID, amount)
	if err != nil {
		t.Fatalf("stripe capture failed: %v", err)
	}
	if tx.Status != "succeeded" {
		t.Errorf("expected succeeded status, got %s", tx.Status)
	}

	// Ingest dispute
	d := adapter.IngestDisputeWebhook("dp_001", tx.TransactionID, amount, "USD", "fraudulent")
	if d.Status != multirail.DisputeStatusNeedsEvidence {
		t.Errorf("expected dispute status NEEDS_EVIDENCE, got %s", d.Status)
	}

	fetched, err := adapter.GetDispute(ctx, "dp_001")
	if err != nil {
		t.Fatalf("failed to fetch dispute: %v", err)
	}
	if fetched.TransactionID != tx.TransactionID {
		t.Errorf("transaction ID mismatch: got %s, want %s", fetched.TransactionID, tx.TransactionID)
	}
}

func TestStripeAdapter_WebhookSignatureVerification(t *testing.T) {
	secret := "whsec_test_secret_abc123"
	adapter := multirail.NewStripeAdapter(multirail.StripeConfig{
		WebhookSecret: secret,
		ToleranceSec:  300,
	})

	rawBody := []byte(`{"id":"evt_test_123","type":"payment_intent.succeeded"}`)
	now := time.Now().Unix()

	// Compute valid HMAC-SHA256 signature
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(fmt.Sprintf("%d", now)))
	mac.Write([]byte("."))
	mac.Write(rawBody)
	validSig := hex.EncodeToString(mac.Sum(nil))

	validHeader := fmt.Sprintf("t=%d,v1=%s", now, validSig)

	// 1. Valid signature passes
	evtID, err := adapter.VerifyWebhook(rawBody, validHeader)
	if err != nil {
		t.Fatalf("expected valid webhook to pass, got: %v", err)
	}
	if evtID == "" {
		t.Fatal("expected non-empty event ID")
	}

	// 2. Replay of identical webhook body is rejected
	_, err = adapter.VerifyWebhook(rawBody, validHeader)
	if err == nil {
		t.Fatal("expected replay attack rejection, got nil error")
	}

	// 3. Forged signature is rejected
	forgedHeader := fmt.Sprintf("t=%d,v1=deadbeefcafebabe0000111122223333444455556666777788889999aaaabbbb", now)
	differentBody := []byte(`{"id":"evt_test_456","type":"charge.refunded"}`)
	_, err = adapter.VerifyWebhook(differentBody, forgedHeader)
	if err == nil {
		t.Fatal("expected forged signature rejection, got nil error")
	}

	// 4. Expired timestamp (> 300s old) is rejected
	oldTimestamp := now - 600
	macOld := hmac.New(sha256.New, []byte(secret))
	macOld.Write([]byte(fmt.Sprintf("%d", oldTimestamp)))
	macOld.Write([]byte("."))
	macOld.Write(differentBody)
	oldSig := hex.EncodeToString(macOld.Sum(nil))

	oldHeader := fmt.Sprintf("t=%d,v1=%s", oldTimestamp, oldSig)
	_, err = adapter.VerifyWebhook(differentBody, oldHeader)
	if err == nil {
		t.Fatal("expected expired timestamp rejection, got nil error")
	}
}

func TestWiseAdapter_PayoutWithUETR(t *testing.T) {
	adapter := multirail.NewWiseAdapter(multirail.WiseConfig{
		APIToken:  "wise_api_token_test",
		ProfileID: "profile_101",
	})
	ctx := context.Background()

	auth := multirail.HumanAuthorization{
		OperatorID: "compliance-officer-44",
		Reason:     "Supplier Invoice #SUPP-992 wire payment",
		Timestamp:  time.Now().UTC(),
	}

	amount := decimal.NewFromFloat(8500.00)
	tx, err := adapter.InitiatePayment(ctx, auth, amount, "EUR", "SUPP-992")
	if err != nil {
		t.Fatalf("wise payment initiation failed: %v", err)
	}

	// BankClearRef must be initially empty to prevent fabricated evidence
	if tx.BankClearRef != "" {
		t.Errorf("BankClearRef must be empty until genuine bank clearance is reported, got %s", tx.BankClearRef)
	}

	// Verify idempotency on identical reference
	dupTx, err := adapter.InitiatePayment(ctx, auth, amount, "EUR", "SUPP-992")
	if err != nil {
		t.Fatalf("idempotent call failed: %v", err)
	}
	if dupTx.TransactionID != tx.TransactionID {
		t.Errorf("expected same transaction ID for idempotent request, got %s vs %s", dupTx.TransactionID, tx.TransactionID)
	}

	// Ingest genuine SWIFT gpi UETR bank receipt
	realUETR := "c47c0b02-601e-450b-9dfd-b5826f4f2081"
	err = adapter.IngestBankClearance(tx.TransactionID, realUETR)
	if err != nil {
		t.Fatalf("failed to ingest bank clearance: %v", err)
	}

	updated, err := adapter.GetTransaction(ctx, tx.TransactionID)
	if err != nil {
		t.Fatalf("failed to get transaction: %v", err)
	}
	if updated.BankClearRef != realUETR {
		t.Errorf("expected bank clearance ref %s, got %s", realUETR, updated.BankClearRef)
	}
	if updated.Status != "outgoing_payment_sent" {
		t.Errorf("expected outgoing_payment_sent, got %s", updated.Status)
	}
}

func TestMultirail_RejectsUnauthorizedMutations(t *testing.T) {
	stripe := multirail.NewStripeAdapter(multirail.StripeConfig{})
	wise := multirail.NewWiseAdapter(multirail.WiseConfig{})
	ctx := context.Background()

	emptyAuth := multirail.HumanAuthorization{} // No operator authorization

	_, err := stripe.InitiatePayment(ctx, emptyAuth, decimal.NewFromInt(100), "USD", "REF-1")
	if err == nil {
		t.Fatal("CRITICAL INVARIANT VIOLATION: Stripe mutation allowed without human authorization")
	}

	_, err = wise.InitiatePayment(ctx, emptyAuth, decimal.NewFromInt(100), "EUR", "REF-2")
	if err == nil {
		t.Fatal("CRITICAL INVARIANT VIOLATION: Wise mutation allowed without human authorization")
	}
}
