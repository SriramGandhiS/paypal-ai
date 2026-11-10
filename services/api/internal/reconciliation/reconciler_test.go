package reconciliation_test

import (
	"testing"

	"github.com/paybridge/paybridge-ai/services/api/internal/reconciliation"
)

func TestReconciliationInvariants(t *testing.T) {
	r := reconciliation.NewReconciler(nil, nil)

	t.Run("matching capture amount results in SYNCHRONISED", func(t *testing.T) {
		truth := reconciliation.CaptureGroundTruth{
			CaptureID:   "cap-12345",
			Status:      "COMPLETED",
			AmountCents: 15000,
			Currency:    "USD",
		}

		outcome := r.EvaluateReconciliation(15000, truth)

		if outcome.ProviderState != "CAPTURE_COMPLETED" {
			t.Errorf("expected provider_state CAPTURE_COMPLETED, got %s", outcome.ProviderState)
		}
		if outcome.ApplicationState != "SYNCHRONISED" {
			t.Errorf("expected application_state SYNCHRONISED, got %s", outcome.ApplicationState)
		}
		if outcome.CaseStatus != "RESOLVED" {
			t.Errorf("expected case_status RESOLVED, got %s", outcome.CaseStatus)
		}
		if outcome.AmountMismatch {
			t.Error("expected AmountMismatch to be false")
		}
	})

	t.Run("amount mismatch preserves CAPTURE_COMPLETED but sets INCONSISTENT", func(t *testing.T) {
		truth := reconciliation.CaptureGroundTruth{
			CaptureID:   "cap-99999",
			Status:      "COMPLETED",
			AmountCents: 12000, // Captured $120.00 instead of expected $150.00
			Currency:    "USD",
		}

		outcome := r.EvaluateReconciliation(15000, truth)

		// Financial Truth Invariant: Ground truth is never falsified
		if outcome.ProviderState != "CAPTURE_COMPLETED" {
			t.Errorf("INVARIANT VIOLATION: provider_state must be CAPTURE_COMPLETED, got %s", outcome.ProviderState)
		}
		if outcome.ApplicationState != "INCONSISTENT" {
			t.Errorf("expected application_state INCONSISTENT, got %s", outcome.ApplicationState)
		}
		if outcome.CaseStatus != "NEEDS_HUMAN_REVIEW" {
			t.Errorf("expected case_status NEEDS_HUMAN_REVIEW, got %s", outcome.CaseStatus)
		}
		if !outcome.AmountMismatch {
			t.Error("expected AmountMismatch to be true")
		}
	})

	t.Run("denied capture sets CAPTURE_DENIED and CLOSED_UNRESOLVED", func(t *testing.T) {
		truth := reconciliation.CaptureGroundTruth{
			CaptureID:   "cap-denied-01",
			Status:      "DENIED",
			AmountCents: 15000,
			Currency:    "USD",
		}

		outcome := r.EvaluateReconciliation(15000, truth)

		if outcome.ProviderState != "CAPTURE_DENIED" {
			t.Errorf("expected provider_state CAPTURE_DENIED, got %s", outcome.ProviderState)
		}
		if outcome.CaseStatus != "CLOSED_UNRESOLVED" {
			t.Errorf("expected case_status CLOSED_UNRESOLVED, got %s", outcome.CaseStatus)
		}
	})
}
