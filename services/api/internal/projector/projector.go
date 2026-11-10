package projector

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	ErrConsentRequired = errors.New("receiver projection requires active consent grant")
	ErrDLPHit           = errors.New("cross-party sensitive field leaked into projected message")
)

// RoleProjection encapsulates data visible to a given actor role
type RoleProjection struct {
	Role            string         `json:"role"` // SENDER, RECEIVER, OPS
	CaseID          string         `json:"case_id"`
	Amount          string         `json:"amount"`
	Currency        string         `json:"currency"`
	PaymentPurpose  string         `json:"payment_purpose"`
	StatusSummary   string         `json:"status_summary"`
	SafeAction      string         `json:"safe_action"`
	VisibleFields   map[string]any `json:"visible_fields"`
	Restricted      bool           `json:"restricted"`
}

type Projector struct {
	bankRegex *regexp.Regexp
	cardRegex *regexp.Regexp
}

func NewProjector() *Projector {
	return &Projector{
		bankRegex: regexp.MustCompile(`\b\d{9,18}\b`),
		cardRegex: regexp.MustCompile(`\b(?:\d{4}[ -]?){3}\d{4}\b`),
	}
}

// BuildSenderProjection builds projection strictly bounded to sender visibility
func (p *Projector) BuildSenderProjection(caseData map[string]any) RoleProjection {
	return RoleProjection{
		Role:           "SENDER",
		CaseID:         fmt.Sprintf("%v", caseData["id"]),
		Amount:         fmt.Sprintf("%v", caseData["amount"]),
		Currency:       fmt.Sprintf("%v", caseData["currency"]),
		PaymentPurpose: fmt.Sprintf("%v", caseData["payment_purpose"]),
		StatusSummary:  fmt.Sprintf("%v", caseData["status_summary"]),
		SafeAction:     fmt.Sprintf("%v", caseData["safe_action"]),
		VisibleFields: map[string]any{
			"sender_country":   caseData["sender_country"],
			"receiver_country": caseData["receiver_country"],
			"order_id":         caseData["order_id"],
		},
		Restricted: false,
	}
}

// BuildReceiverProjection builds projection strictly bounded to receiver visibility (consent-gated)
func (p *Projector) BuildReceiverProjection(caseData map[string]any, hasConsent bool) (RoleProjection, error) {
	if !hasConsent {
		return RoleProjection{
			Role:       "RECEIVER",
			CaseID:     fmt.Sprintf("%v", caseData["id"]),
			Restricted: true,
		}, ErrConsentRequired
	}

	return RoleProjection{
		Role:           "RECEIVER",
		CaseID:         fmt.Sprintf("%v", caseData["id"]),
		Amount:         fmt.Sprintf("%v", caseData["amount"]),
		Currency:       fmt.Sprintf("%v", caseData["currency"]),
		PaymentPurpose: fmt.Sprintf("%v", caseData["payment_purpose"]),
		StatusSummary:  fmt.Sprintf("%v", caseData["status_summary"]),
		SafeAction:     fmt.Sprintf("%v", caseData["safe_action"]),
		VisibleFields: map[string]any{
			"sender_country":   caseData["sender_country"],
			"receiver_country": caseData["receiver_country"],
			// Receiver sees only agreed fields; never sender's funding instrument or internal IDs
		},
		Restricted: false,
	}, nil
}

// ScanAndRedactDLP inspects role-directed text against unauthorized cross-party patterns
func (p *Projector) ScanAndRedactDLP(text string, targetRole string, forbiddenFingerprints []string) (string, error) {
	clean := text
	for _, fp := range forbiddenFingerprints {
		if fp != "" && strings.Contains(clean, fp) {
			clean = strings.ReplaceAll(clean, fp, "[REDACTED-CONFIDENTIAL]")
		}
	}

	// If target is sender view, check for bank account numbers leaking
	if targetRole == "SENDER" {
		if p.bankRegex.MatchString(clean) {
			clean = p.bankRegex.ReplaceAllString(clean, "[REDACTED-BANK-ACCOUNT]")
		}
	}

	// Redact credit card numbers
	if p.cardRegex.MatchString(clean) {
		clean = p.cardRegex.ReplaceAllString(clean, "[REDACTED-CARD]")
	}

	return clean, nil
}
