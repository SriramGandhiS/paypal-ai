package claims

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/paybridge/paybridge-ai/services/api/internal/reducer"
)

var (
	ErrInvalidClaimSchema   = errors.New("claim failed JSON schema validation")
	ErrMissingCitation      = errors.New("confirmed claim lacks required evidence citation")
	ErrEvidenceNotFound     = errors.New("cited evidence ID does not exist in frozen snapshot")
	ErrInsufficientTier     = errors.New("confirmed claim requires Tier E1 or E2 evidence")
	ErrForbiddenTerminology = errors.New("claim contains forbidden misleading terminology")
	ErrPromptInjection      = errors.New("claim text contains suspected prompt injection instruction")
	ErrCrossPartyLeak       = errors.New("claim text leaked cross-party restricted account data")
)

// ClaimKind represents epistemic certainty
type ClaimKind string

const (
	ClaimKindConfirmed ClaimKind = "confirmed"
	ClaimKindInference ClaimKind = "inference"
	ClaimKindUnknown   ClaimKind = "unknown"
)

// Claim represents a structured finding produced by an AI specialist agent
type Claim struct {
	Code            string    `json:"code"`
	Kind            ClaimKind `json:"kind"`
	Statement       string    `json:"statement"`
	EvidenceIDs     []string  `json:"evidence_ids"`
	InferenceRuleID string    `json:"inference_rule_id,omitempty"`
	ViewerScope     string    `json:"viewer_scope"` // SENDER, RECEIVER, OPS, ALL
	Validated       bool      `json:"validated"`
	DowngradeReason string    `json:"downgrade_reason,omitempty"`
}

// SnapshotEvidenceItem represents an evidence item in the frozen snapshot
type SnapshotEvidenceItem struct {
	ID          string               `json:"id"`
	Tier        reducer.EvidenceTier `json:"tier"`
	SourceTrust reducer.SourceTrust  `json:"source_trust"`
	Kind        string               `json:"kind"`
	Verified    bool                 `json:"verified"`
}

// EvidenceSnapshot represents the immutable state against which claims are verified
type EvidenceSnapshot struct {
	SnapshotVersion int64                            `json:"snapshot_version"`
	CaseID          string                           `json:"case_id"`
	EvidenceItems   map[string]SnapshotEvidenceItem `json:"evidence_items"`
}

// Validator validates claims against evidence, safety, and privacy invariants
type Validator struct {
	bankRegex     *regexp.Regexp
	cardRegex     *regexp.Regexp
	injectionPats []*regexp.Regexp
}

func NewValidator() *Validator {
	return &Validator{
		bankRegex: regexp.MustCompile(`\b\d{9,18}\b`),
		cardRegex: regexp.MustCompile(`\b(?:\d{4}[ -]?){3}\d{4}\b`),
		injectionPats: []*regexp.Regexp{
			regexp.MustCompile(`(?i)\bignore\s+(?:all\s+)?previous\s+instructions\b`),
			regexp.MustCompile(`(?i)\bsystem\s+override\b`),
			regexp.MustCompile(`(?i)<script\b`),
			regexp.MustCompile(`(?i)\bdrop\s+table\b`),
			regexp.MustCompile(`(?i)\btransfer\s+(?:all\s+)?funds\s+now\b`),
		},
	}
}

// ValidateClaim validates and enforces citation integrity, tier rules, calm copy rules,
// prompt injection defenses, and cross-party DLP restrictions.
func (v *Validator) ValidateClaim(claim Claim, snapshot *EvidenceSnapshot) (Claim, error) {
	out := claim

	// Schema basic sanity check
	if strings.TrimSpace(out.Code) == "" || strings.TrimSpace(out.Statement) == "" {
		return out, ErrInvalidClaimSchema
	}

	// 1. Prompt Injection Boundary Check
	for _, pat := range v.injectionPats {
		if pat.MatchString(out.Statement) {
			out.Validated = false
			out.Kind = ClaimKindUnknown
			out.DowngradeReason = "Rejected: suspected adversarial prompt injection in statement"
			return out, fmt.Errorf("%w: statement matched injection signature", ErrPromptInjection)
		}
	}

	// 2. Cross-Party DLP / Privacy Check
	// If viewer scope is SENDER, sensitive receiver banking numbers must not appear
	if out.ViewerScope == "SENDER" && v.bankRegex.MatchString(out.Statement) {
		out.Statement = v.bankRegex.ReplaceAllString(out.Statement, "[REDACTED-BANK-ACCOUNT]")
	}
	if v.cardRegex.MatchString(out.Statement) {
		out.Statement = v.cardRegex.ReplaceAllString(out.Statement, "[REDACTED-CARD]")
	}

	// 3. Forbidden language check: Never allow "failed" to describe ambiguous/unknown outcomes
	statementLower := strings.ToLower(out.Statement)
	if strings.Contains(statementLower, "payment failed") || strings.Contains(statementLower, "transaction failed") {
		out.Validated = false
		return out, fmt.Errorf("%w: cannot use 'payment failed' on ambiguous payment outcomes", ErrForbiddenTerminology)
	}

	// 4. Epistemic Gates for Confirmed Claims
	if out.Kind == ClaimKindConfirmed {
		// Citation Gate: Every confirmed claim MUST cite at least one evidence ID
		if len(out.EvidenceIDs) == 0 {
			out.Kind = ClaimKindUnknown
			out.DowngradeReason = "Downgraded to unknown: missing required evidence citation"
			out.Validated = false
			return out, ErrMissingCitation
		}

		// Verify that all cited evidence IDs exist in snapshot and check tiers
		hasE1orE2 := false
		for _, eid := range out.EvidenceIDs {
			item, exists := snapshot.EvidenceItems[eid]
			if !exists {
				out.Kind = ClaimKindUnknown
				out.DowngradeReason = fmt.Sprintf("Downgraded to unknown: cited evidence ID %s does not exist in snapshot", eid)
				out.Validated = false
				return out, ErrEvidenceNotFound
			}

			// Tier Gate: Confirmed claim requires Tier E1 (direct GET) or Tier E2 (verified webhook)
			if item.Tier == reducer.TierE1DirectProviderRead || item.Tier == reducer.TierE2VerifiedWebhook {
				if item.Verified {
					hasE1orE2 = true
				}
			}
		}

		if !hasE1orE2 {
			out.Kind = ClaimKindInference
			out.DowngradeReason = "Downgraded to inference: confirmed claims require verified Tier E1 or E2 evidence"
			out.Validated = false
			return out, ErrInsufficientTier
		}
	}

	out.Validated = true
	return out, nil
}
