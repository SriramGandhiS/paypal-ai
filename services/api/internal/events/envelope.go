package events

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

// CanonicalEvent represents the normalized immutable event envelope across all payment rails.
type CanonicalEvent struct {
	ID                 string     `json:"id"`
	TenantID           string     `json:"tenant_id"`
	CaseID             string     `json:"case_id"`
	EvidenceID         string     `json:"evidence_id"`
	Provider           string     `json:"provider"`
	ProviderEventID    string     `json:"provider_event_id,omitempty"`
	Environment        string     `json:"environment"`
	EventType          string     `json:"event_type"`
	StatusObserved     string     `json:"status_observed"`
	SourceTrust        string     `json:"source_trust"`
	Verified           bool       `json:"verified"`
	ProviderOccurredAt *time.Time `json:"provider_occurred_at,omitempty"`
	ReceivedAt         time.Time  `json:"received_at"`
	SchemaVersion      int        `json:"schema_version"`
}

// ComputeEventDigest calculates the deterministic SHA-256 hash of the canonical event.
func (e *CanonicalEvent) ComputeEventDigest() (string, error) {
	bytes, err := json.Marshal(e)
	if err != nil {
		return "", fmt.Errorf("failed to serialize event: %w", err)
	}
	h := sha256.Sum256(bytes)
	return hex.EncodeToString(h[:]), nil
}
