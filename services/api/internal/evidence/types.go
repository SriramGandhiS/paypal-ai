package evidence

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"
)

var (
	ErrInvalidTier       = errors.New("tier must be between 1 and 6")
	ErrSimulatedTierRule = errors.New("simulated source trust must always be assigned Tier 6 (I11)")
	ErrMissingPayload    = errors.New("raw payload is required for SHA-256 calculation")
	ErrEvidenceNotFound  = errors.New("evidence item not found")
)

// EvidenceItem represents an immutable piece of financial provenance.
type EvidenceItem struct {
	ID            string          `json:"id"`
	TenantID      string          `json:"tenant_id"`
	CaseID        string          `json:"case_id"`
	Kind          string          `json:"kind"`
	Tier          int             `json:"tier"`
	SourceTrust   string          `json:"source_trust"`
	Environment   string          `json:"environment"`
	Verified      bool            `json:"verified"`
	UntrustedText bool            `json:"untrusted_text"`
	OccurredAt    *time.Time      `json:"occurred_at,omitempty"`
	ObservedAt    time.Time       `json:"observed_at"`
	RedactedJSON  json.RawMessage `json:"redacted_extract"`
	PayloadSHA256 string          `json:"payload_sha256"`
	AccessClass   string          `json:"access_class"`
}

// IngestInput provides parameters for evidence registration.
type IngestInput struct {
	CaseID        string          `json:"case_id"`
	Kind          string          `json:"kind"`
	Tier          int             `json:"tier"`
	SourceTrust   string          `json:"source_trust"`
	OccurredAt    *time.Time      `json:"occurred_at,omitempty"`
	RawPayload    []byte          `json:"-"`
	RedactedJSON  json.RawMessage `json:"redacted_extract"`
	UntrustedText bool            `json:"untrusted_text"`
	AccessClass   string          `json:"access_class"`
}

// Validate checks evidence constraints including I11.
func (in *IngestInput) Validate() error {
	if in.Tier < 1 || in.Tier > 6 {
		return ErrInvalidTier
	}
	if in.SourceTrust == "SIMULATED" && in.Tier != 6 {
		return ErrSimulatedTierRule
	}
	if len(in.RawPayload) == 0 {
		return ErrMissingPayload
	}
	if in.AccessClass == "" {
		in.AccessClass = "TENANT"
	}
	return nil
}

// ComputePayloadHash returns the hex-encoded SHA-256 digest of the raw bytes.
func ComputePayloadHash(payload []byte) string {
	h := sha256.Sum256(payload)
	return hex.EncodeToString(h[:])
}
