package chat

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/paybridge/paybridge-ai/services/api/internal/domain"
	"github.com/shopspring/decimal"
)

func makeTestCapturedCase() (*domain.PaymentCase, []domain.EvidenceItem) {
	c := &domain.PaymentCase{
		ID:                uuid.New(),
		TenantID:          uuid.New(),
		ExternalReference: "9K201824",
		Amount:            decimal.NewFromFloat(480.00),
		Currency:          "USD",
		ProviderState:     domain.ProviderStateCaptureCompleted,
		ApplicationState:  domain.AppStateAPITimeout,
		ReceiverState:     domain.ReceiverStatePending,
		CreatedAt:         time.Now().UTC(),
	}

	ev := []domain.EvidenceItem{
		{
			ID:          uuid.New(),
			CaseID:      c.ID,
			Kind:        "PAYPAL_ORDER_CAPTURE",
			Tier:        domain.Tier1ProviderDirect,
			Verified:    true,
			ObservedAt:  time.Now().UTC(),
			AccessClass: "PCI_REDACTED",
		},
	}
	return c, ev
}

func makeTestDeniedCase() (*domain.PaymentCase, []domain.EvidenceItem) {
	c := &domain.PaymentCase{
		ID:                uuid.New(),
		TenantID:          uuid.New(),
		ExternalReference: "8K110492",
		Amount:            decimal.NewFromFloat(740.00),
		Currency:          "EUR",
		ProviderState:     domain.ProviderStateCaptureDenied,
		ApplicationState:  domain.AppStateSynchronised,
		ReceiverState:     domain.ReceiverStateNotObservable,
		CreatedAt:         time.Now().UTC(),
	}

	ev := []domain.EvidenceItem{
		{
			ID:          uuid.New(),
			CaseID:      c.ID,
			Kind:        "ISO_8583_DECLINE_CODE_57",
			Tier:        domain.Tier1ProviderDirect,
			Verified:    true,
			ObservedAt:  time.Now().UTC(),
			AccessClass: "PUBLIC",
		},
	}
	return c, ev
}

// 1. "Was I charged?" checks capture evidence
func TestWasIChargedChecksCaptureEvidence(t *testing.T) {
	c, ev := makeTestCapturedCase()
	session := NewInvestigationSession(c, ev, "Gateway Timeout 504", "Do not retry checkout")

	resp, err := session.Respond(context.Background(), "Was I charged for this order?")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(resp.Answer, "Yes, you were charged") {
		t.Errorf("expected confirmation of charge, got: %s", resp.Answer)
	}
	if !strings.Contains(resp.Answer, "480") || !strings.Contains(resp.Answer, "USD") {
		t.Errorf("expected amount to be cited in answer, got: %s", resp.Answer)
	}

	// For denied case
	cDenied, evDenied := makeTestDeniedCase()
	sessionDenied := NewInvestigationSession(cDenied, evDenied, "International Usage Disabled", "Enable international toggle")

	respDenied, err := sessionDenied.Respond(context.Background(), "Was my card charged?")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(respDenied.Answer, "NOT charged") {
		t.Errorf("expected denial confirmation, got: %s", respDenied.Answer)
	}
}

// 2. "Should I retry?" checks payment state and safe-action policy
func TestShouldIRetryChecksSafeActionPolicy(t *testing.T) {
	c, ev := makeTestCapturedCase()
	session := NewInvestigationSession(c, ev, "Gateway Timeout 504", "Do not retry checkout")

	resp, err := session.Respond(context.Background(), "Should I try again and pay?")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(resp.Answer, "Do NOT retry") {
		t.Errorf("expected strict warning against retry, got: %s", resp.Answer)
	}
	if !resp.AttemptLockHeld {
		t.Errorf("expected AttemptLockHeld to be true")
	}
}

// 3. "Why did my bank reject it?" does not fabricate a decline when captured
func TestWhyDidBankRejectDoesNotFabricateDeclineWhenCaptured(t *testing.T) {
	c, ev := makeTestCapturedCase()
	session := NewInvestigationSession(c, ev, "Gateway Timeout 504", "Do not retry checkout")

	resp, err := session.Respond(context.Background(), "Why did my bank reject it?")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(resp.Answer, "Your bank declined") {
		t.Errorf("fabricated bank decline when transaction was actually captured: %s", resp.Answer)
	}
	if !strings.Contains(resp.Answer, "no evidence that your bank rejected") {
		t.Errorf("expected clear explanation that bank did not reject, got: %s", resp.Answer)
	}
}

// 4. "Where is my money?" does not invent settlement information
func TestWhereIsMyMoneyDoesNotInventSettlement(t *testing.T) {
	c, ev := makeTestCapturedCase() // ReceiverState is Pending (not settled)
	session := NewInvestigationSession(c, ev, "Gateway Timeout 504", "Do not retry checkout")

	resp, err := session.Respond(context.Background(), "Where is my money right now?")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(resp.Answer, "settled into the merchant bank account") {
		t.Errorf("fabricated bank settlement when evidence was only escrow pending: %s", resp.Answer)
	}
	if !strings.Contains(resp.Answer, "PayPal's verified escrow ledger") {
		t.Errorf("expected verified escrow explanation, got: %s", resp.Answer)
	}
}

// 5. Repeated questions do not return blindly duplicated canned responses
func TestRepeatedQuestionsDoNotReturnDuplicateCanned(t *testing.T) {
	c, ev := makeTestCapturedCase()
	session := NewInvestigationSession(c, ev, "Gateway Timeout 504", "Do not retry checkout")

	resp1, _ := session.Respond(context.Background(), "Was I charged?")
	resp2, _ := session.Respond(context.Background(), "Was I charged?")

	if resp1.Answer == resp2.Answer {
		t.Errorf("expected natural variation on repeated question, but got exact duplicate: %s", resp1.Answer)
	}
	if !strings.Contains(resp2.Answer, "re-confirm") {
		t.Errorf("expected conversational acknowledgment on repeat question, got: %s", resp2.Answer)
	}
}

// 6. Follow-up answers use the correct case context
func TestFollowUpAnswersUseCorrectCaseContext(t *testing.T) {
	c, ev := makeTestCapturedCase()
	session := NewInvestigationSession(c, ev, "Gateway Timeout 504", "Reconcile original order reference")

	session.Respond(context.Background(), "Was I charged?")
	respNext, err := session.Respond(context.Background(), "What should I do next?")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(respNext.Answer, "Reconcile original order reference") {
		t.Errorf("expected case safe action in answer, got: %s", respNext.Answer)
	}
	if len(session.History) != 4 {
		t.Errorf("expected 4 history messages (2 turns), got: %d", len(session.History))
	}
}

// 7. Conflicting / unknown evidence produces an uncertainty explanation
func TestConflictingUnknownEvidenceProducesUncertainty(t *testing.T) {
	cUnknown := &domain.PaymentCase{
		ID:                uuid.New(),
		TenantID:          uuid.New(),
		ExternalReference: "UNK-9912",
		Amount:            decimal.NewFromFloat(100.00),
		Currency:          "USD",
		ProviderState:     domain.ProviderStateUnknown,
		ApplicationState:  domain.AppStateInconsistent,
		ReceiverState:     domain.ReceiverStateNotObservable,
		CreatedAt:         time.Now().UTC(),
	}

	session := NewInvestigationSession(cUnknown, nil, "Inconsistent Telemetry", "Hold mutex lock")

	resp, err := session.Respond(context.Background(), "Was I charged for this transaction?")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Confidence != ConfidenceUnknown {
		t.Errorf("expected ConfidenceUnknown, got: %s", resp.Confidence)
	}
	if !strings.Contains(resp.Answer, "can't confirm that from the available payment records yet") {
		t.Errorf("expected uncertainty disclosure, got: %s", resp.Answer)
	}
}
