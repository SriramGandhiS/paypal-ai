package projector

import (
	"testing"
)

func TestProjector_ReceiverConsentGating(t *testing.T) {
	p := NewProjector()
	caseData := map[string]any{
		"id":              "case-xyz",
		"amount":          "500.00",
		"currency":        "USD",
		"payment_purpose": "Software development services",
	}

	// Without consent, must fail with ErrConsentRequired
	proj, err := p.BuildReceiverProjection(caseData, false)
	if err != ErrConsentRequired {
		t.Fatalf("expected ErrConsentRequired, got %v", err)
	}
	if !proj.Restricted {
		t.Errorf("expected Restricted=true")
	}

	// With consent, succeeds
	proj, err = p.BuildReceiverProjection(caseData, true)
	if err != nil {
		t.Fatalf("unexpected error with consent: %v", err)
	}
	if proj.Restricted {
		t.Errorf("expected Restricted=false")
	}
	if proj.Amount != "500.00" {
		t.Errorf("expected amount 500.00, got %s", proj.Amount)
	}
}

func TestProjector_DLPRedaction(t *testing.T) {
	p := NewProjector()

	rawMessage := "Hello sender, the receiver's private account 9876543210987 has not yet received funds."
	cleaned, err := p.ScanAndRedactDLP(rawMessage, "SENDER", []string{"receiver@privatecorp.internal"})
	if err != nil {
		t.Fatalf("DLP error: %v", err)
	}

	if cleaned == rawMessage {
		t.Errorf("expected bank account number to be redacted")
	}
}
