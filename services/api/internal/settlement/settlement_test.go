package settlement_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/paybridge/paybridge-ai/services/api/internal/domain"
	"github.com/paybridge/paybridge-ai/services/api/internal/reducer"
	"github.com/paybridge/paybridge-ai/services/api/internal/settlement"
)

func TestDepositMatcher_SuccessfulBankVerification(t *testing.T) {
	red := reducer.NewReducer()
	matcher := settlement.NewDepositMatcher(red)
	ctx := context.Background()

	amt := decimal.NewFromFloat(1250.00)
	now := time.Now().UTC()

	initialCase := &domain.PaymentCase{
		ID:                uuid.New(),
		TenantID:          uuid.New(),
		ExternalReference: "INV-2026-0881",
		Amount:            amt,
		Currency:          "USD",
		Environment:       domain.EnvPayPalSimulator,
		ProviderState:     domain.ProviderStateCaptureCompleted,
		ReceiverState:     domain.ReceiverStateNotObservable,
		StateVersion:      3,
	}

	deposit := settlement.BankTransaction{
		AccountID:   "acct_hdfc_9912",
		BankName:    "HDFC Bank",
		Amount:      amt,
		Currency:    "USD",
		Narration:   "NEFT-UTR981240-INV-2026-0881-PAYPAL-SETTLE",
		OccurredAt:  now.Add(10 * time.Minute),
		SignedUTR:   "UTR98124001923",
		ProviderRef: "CAP-9912",
		TrustLevel:  settlement.TrustAuthenticatedProviderEvidence,
	}

	updated, evid, err := matcher.MatchAndVerify(ctx, initialCase, now, deposit)
	if err != nil {
		t.Fatalf("deposit match failed: %v", err)
	}

	if updated.ReceiverState != domain.ReceiverStateVerifiedBankCredited {
		t.Errorf("expected ReceiverState=VERIFIED_BANK_CREDITED, got %s", updated.ReceiverState)
	}

	if !evid.Verified {
		t.Error("expected generated bank evidence to be verified")
	}

	if evid.Tier != domain.Tier1ProviderDirect {
		t.Errorf("expected Tier1ProviderDirect, got %d", evid.Tier)
	}
}

func TestDepositMatcher_DistinguishesMockFromAuthenticated(t *testing.T) {
	red := reducer.NewReducer()
	matcher := settlement.NewDepositMatcher(red)
	ctx := context.Background()

	amt := decimal.NewFromFloat(500.00)
	now := time.Now().UTC()

	initialCase := &domain.PaymentCase{
		ID:                uuid.New(),
		TenantID:          uuid.New(),
		ExternalReference: "INV-MOCK-001",
		Amount:            amt,
		Currency:          "USD",
		Environment:       domain.EnvPayPalSimulator,
		ProviderState:     domain.ProviderStateCaptureCompleted,
		ReceiverState:     domain.ReceiverStateNotObservable,
		StateVersion:      1,
	}

	// Mock bank transaction must never be labelled PROVIDER_VERIFIED or transition receiver_state
	mockDeposit := settlement.BankTransaction{
		AccountID:   "acct_sim_001",
		BankName:    "Simulated Bank",
		Amount:      amt,
		Currency:    "USD",
		Narration:   "TRANSFER-INV-MOCK-001",
		OccurredAt:  now.Add(5 * time.Minute),
		SignedUTR:   "SIM-UTR-001",
		TrustLevel:  settlement.TrustMockBankTransaction,
	}

	updated, evid, err := matcher.MatchAndVerify(ctx, initialCase, now, mockDeposit)
	if err != nil {
		t.Fatalf("unexpected error on mock deposit: %v", err)
	}

	if evid.Tier != domain.Tier6Simulated {
		t.Errorf("expected Tier6Simulated for mock deposit, got %d", evid.Tier)
	}

	if evid.SourceTrust != "SIMULATED" {
		t.Errorf("expected SIMULATED source trust, got %s", evid.SourceTrust)
	}

	// CRITICAL INVARIANT: Receiver state must NOT transition to VERIFIED_BANK_CREDITED on mock data
	if updated.ReceiverState == domain.ReceiverStateVerifiedBankCredited {
		t.Fatal("CRITICAL INVARIANT VIOLATION: Mock bank transaction transitioned case to VERIFIED_BANK_CREDITED")
	}
}

func TestDepositMatcher_RejectsAmountMismatch(t *testing.T) {
	red := reducer.NewReducer()
	matcher := settlement.NewDepositMatcher(red)
	ctx := context.Background()

	amt := decimal.NewFromFloat(100.00)
	now := time.Now().UTC()

	initialCase := &domain.PaymentCase{
		ID:                uuid.New(),
		TenantID:          uuid.New(),
		ExternalReference: "INV-2026-9999",
		Amount:            amt,
		Currency:          "USD",
		Environment:       domain.EnvPayPalSimulator,
		ProviderState:     domain.ProviderStateCaptureCompleted,
		ReceiverState:     domain.ReceiverStateNotObservable,
	}

	// Deposit amount is off by 1 cent
	deposit := settlement.BankTransaction{
		AccountID:   "acct_1",
		BankName:    "Chase",
		Amount:      decimal.NewFromFloat(99.99),
		Currency:    "USD",
		Narration:   "INV-2026-9999",
		OccurredAt:  now.Add(5 * time.Minute),
		SignedUTR:   "UTR-111",
		TrustLevel:  settlement.TrustAuthenticatedProviderEvidence,
	}

	_, _, err := matcher.MatchAndVerify(ctx, initialCase, now, deposit)
	if err == nil {
		t.Fatal("expected amount mismatch error, got nil")
	}
}

func TestDepositMatcher_RejectsTemporalOrderViolation(t *testing.T) {
	red := reducer.NewReducer()
	matcher := settlement.NewDepositMatcher(red)
	ctx := context.Background()

	amt := decimal.NewFromFloat(200.00)
	captureTime := time.Now().UTC()

	initialCase := &domain.PaymentCase{
		ID:                uuid.New(),
		TenantID:          uuid.New(),
		ExternalReference: "INV-PRE-01",
		Amount:            amt,
		Currency:          "USD",
		Environment:       domain.EnvPayPalSimulator,
		ProviderState:     domain.ProviderStateCaptureCompleted,
		ReceiverState:     domain.ReceiverStateNotObservable,
	}

	// Deposit timestamp is 1 hour BEFORE capture
	deposit := settlement.BankTransaction{
		AccountID:   "acct_2",
		BankName:    "Barclays",
		Amount:      amt,
		Currency:    "USD",
		Narration:   "INV-PRE-01",
		OccurredAt:  captureTime.Add(-1 * time.Hour),
		SignedUTR:   "UTR-PRE-01",
		TrustLevel:  settlement.TrustAuthenticatedProviderEvidence,
	}

	_, _, err := matcher.MatchAndVerify(ctx, initialCase, captureTime, deposit)
	if err == nil {
		t.Fatal("expected temporal ordering violation error, got nil")
	}
}

func TestPlaidClient_LifecycleAndConsentRevocation(t *testing.T) {
	client := settlement.NewPlaidClient("plaid_test_client_id", "plaid_test_secret", "sandbox")
	ctx := context.Background()

	linkToken, err := client.CreateLinkToken(ctx, "user_123", "PayBridge Enterprise")
	if err != nil {
		t.Fatalf("failed to create link token: %v", err)
	}
	if linkToken == "" {
		t.Fatal("expected non-empty link token")
	}

	accessToken, itemID, err := client.ExchangePublicToken(ctx, "public-sandbox-token-123")
	if err != nil {
		t.Fatalf("failed to exchange public token: %v", err)
	}
	if accessToken == "" || itemID == "" {
		t.Fatal("expected valid access token and item ID")
	}

	// Sync transactions
	res, err := client.SyncTransactions(ctx, itemID)
	if err != nil {
		t.Fatalf("transactions sync failed: %v", err)
	}
	if res.Cursor == "" {
		t.Fatal("expected valid cursor")
	}

	// Revoke consent
	err = client.RevokeToken(ctx, itemID)
	if err != nil {
		t.Fatalf("failed to revoke token: %v", err)
	}

	// Post-revocation sync must fail
	_, err = client.SyncTransactions(ctx, itemID)
	if err == nil {
		t.Fatal("expected sync failure on revoked consent, got nil error")
	}
}

func TestSetuAAClient_ConsentLifecycle(t *testing.T) {
	client := settlement.NewSetuAAClient("setu_client_id_test", "setu_secret_test")
	ctx := context.Background()

	consentID, err := client.CreateConsent(ctx, "payer@okhdfcbank")
	if err != nil {
		t.Fatalf("failed to create consent: %v", err)
	}

	status, err := client.GetConsentStatus(ctx, consentID)
	if err != nil {
		t.Fatalf("failed to get consent status: %v", err)
	}
	if status != "ACTIVE" {
		t.Errorf("expected ACTIVE consent, got %s", status)
	}

	err = client.RevokeConsent(ctx, consentID)
	if err != nil {
		t.Fatalf("failed to revoke consent: %v", err)
	}

	status, _ = client.GetConsentStatus(ctx, consentID)
	if status != "REVOKED" {
		t.Errorf("expected REVOKED consent, got %s", status)
	}
}
