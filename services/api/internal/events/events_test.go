package events_test

import (
	"testing"
	"time"

	"github.com/paybridge/paybridge-ai/services/api/internal/events"
)

func TestCanonicalEventDigest(t *testing.T) {
	evt := events.CanonicalEvent{
		ID:              "evt-01",
		TenantID:        "018f3a21-9b34-7a12-89cd-456789000000",
		CaseID:          "018f3a21-9b34-7a12-89cd-456789111111",
		EvidenceID:      "018f3a21-9b34-7a12-89cd-456789222222",
		Provider:        "PAYPAL",
		ProviderEventID: "paypal-event-12345",
		Environment:     "REAL_SANDBOX",
		EventType:       "PAYMENT.CAPTURE.COMPLETED",
		StatusObserved:  "CAPTURE_COMPLETED",
		SourceTrust:     "PROVIDER_VERIFIED",
		Verified:        true,
		ReceivedAt:      time.Date(2026, 11, 10, 12, 0, 0, 0, time.UTC),
		SchemaVersion:   1,
	}

	digest, err := evt.ComputeEventDigest()
	if err != nil {
		t.Fatalf("unexpected error computing digest: %v", err)
	}
	if len(digest) != 64 {
		t.Errorf("expected 64 hex chars, got %d", len(digest))
	}
}
