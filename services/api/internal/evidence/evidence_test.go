package evidence_test

import (
	"testing"

	"github.com/paybridge/paybridge-ai/services/api/internal/evidence"
)

func TestEvidenceValidation(t *testing.T) {
	t.Run("tier out of range error", func(t *testing.T) {
		in := evidence.IngestInput{
			Tier:        7,
			SourceTrust: "PROVIDER_VERIFIED",
			RawPayload:  []byte("test"),
		}
		if err := in.Validate(); err != evidence.ErrInvalidTier {
			t.Errorf("expected ErrInvalidTier, got %v", err)
		}
	})

	t.Run("simulated source trust must be tier 6", func(t *testing.T) {
		in := evidence.IngestInput{
			Tier:        1,
			SourceTrust: "SIMULATED",
			RawPayload:  []byte("test"),
		}
		if err := in.Validate(); err != evidence.ErrSimulatedTierRule {
			t.Errorf("expected ErrSimulatedTierRule, got %v", err)
		}
	})

	t.Run("valid tier 6 simulation passes", func(t *testing.T) {
		in := evidence.IngestInput{
			Tier:        6,
			SourceTrust: "SIMULATED",
			RawPayload:  []byte("test"),
		}
		if err := in.Validate(); err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
	})

	t.Run("compute payload SHA256", func(t *testing.T) {
		data := []byte("paypal-order-payload-12345")
		h := evidence.ComputePayloadHash(data)
		if len(h) != 64 {
			t.Errorf("expected 64 hex chars, got %d", len(h))
		}
	})
}
