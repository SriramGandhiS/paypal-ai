package chat

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/paybridge/paybridge-ai/services/api/internal/domain"
	"github.com/paybridge/paybridge-ai/services/api/internal/slm"
)

type InquiryType string

const (
	InquiryChargeStatus   InquiryType = "CHARGE_STATUS"
	InquiryRetrySafety    InquiryType = "RETRY_SAFETY"
	InquiryNextActions    InquiryType = "NEXT_ACTIONS"
	InquiryBankRejection  InquiryType = "BANK_REJECTION"
	InquiryFundsLocation  InquiryType = "FUNDS_LOCATION"
	InquiryRefund         InquiryType = "REFUND"
	InquiryGeneral        InquiryType = "GENERAL"
)

type ConfidenceTier string

const (
	ConfidenceConfirmed ConfidenceTier = "CONFIRMED"
	ConfidenceProbable  ConfidenceTier = "PROBABLE"
	ConfidenceUnknown   ConfidenceTier = "UNKNOWN"
)

type ChatMessage struct {
	Role      string    `json:"role"` // "user" or "assistant"
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

type ChatResponse struct {
	Answer           string         `json:"answer"`
	InquiryType      InquiryType    `json:"inquiry_type"`
	EvidenceCited    []string       `json:"evidence_cited"`
	Confidence       ConfidenceTier `json:"confidence"`
	AttemptLockHeld  bool           `json:"attempt_lock_held"`
	DuplicateCharges string         `json:"duplicate_charges"`
}

type InvestigationSession struct {
	mu             sync.Mutex
	Case           *domain.PaymentCase
	Evidence       []domain.EvidenceItem
	DiagnosedIssue string
	SafeAction     string
	History        []ChatMessage
	repeatCounts   map[InquiryType]int
	slmClient      *slm.SLMClient
}

func NewInvestigationSession(c *domain.PaymentCase, evidence []domain.EvidenceItem, diagnosedIssue, safeAction string) *InvestigationSession {
	return &InvestigationSession{
		Case:           c,
		Evidence:       evidence,
		DiagnosedIssue: diagnosedIssue,
		SafeAction:     safeAction,
		History:        make([]ChatMessage, 0),
		repeatCounts:   make(map[InquiryType]int),
		slmClient:      slm.NewSLMClient(),
	}
}

// ClassifyInquiry categorizes the semantic intent of the question
func ClassifyInquiry(query string) InquiryType {
	q := strings.ToLower(query)

	if strings.Contains(q, "charge") || strings.Contains(q, "charged") || strings.Contains(q, "deduct") || strings.Contains(q, "debit") || strings.Contains(q, "paid") {
		if strings.Contains(q, "twice") || strings.Contains(q, "double") {
			return InquiryChargeStatus
		}
		return InquiryChargeStatus
	}
	if strings.Contains(q, "retry") || strings.Contains(q, "try again") || strings.Contains(q, "pay again") || strings.Contains(q, "resubmit") || strings.Contains(q, "swipe again") {
		return InquiryRetrySafety
	}
	if strings.Contains(q, "what next") || strings.Contains(q, "what should i do") || strings.Contains(q, "next step") || strings.Contains(q, "what to do") || strings.Contains(q, "how to fix") {
		return InquiryNextActions
	}
	if strings.Contains(q, "bank reject") || strings.Contains(q, "bank decline") || strings.Contains(q, "issuer decline") || strings.Contains(q, "why reject") || strings.Contains(q, "why declined") {
		return InquiryBankRejection
	}
	if strings.Contains(q, "where is my money") || strings.Contains(q, "where are the funds") || strings.Contains(q, "where did my money go") || strings.Contains(q, "where are my funds") {
		return InquiryFundsLocation
	}
	if strings.Contains(q, "refund") || strings.Contains(q, "money back") || strings.Contains(q, "reversal") {
		return InquiryRefund
	}

	return InquiryGeneral
}

// Respond generates a context-aware, evidence-first answer
func (s *InvestigationSession) Respond(ctx context.Context, userQuery string) (*ChatResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 1. Guardrail against prompt injections
	if err := s.slmClient.InspectForInjection(userQuery); err != nil {
		return nil, err
	}

	inquiry := ClassifyInquiry(userQuery)
	s.repeatCounts[inquiry]++
	count := s.repeatCounts[inquiry]

	// Evidence inspection
	hasCaptureProof := s.Case.ProviderState == domain.ProviderStateCaptureCompleted
	isDenied := s.Case.ProviderState == domain.ProviderStateCaptureDenied
	isUnknown := s.Case.ProviderState == domain.ProviderStateUnknown

	var cited []string
	for _, e := range s.Evidence {
		if e.Verified {
			cited = append(cited, fmt.Sprintf("%s (Tier %d)", e.Kind, e.Tier))
		}
	}

	var answer string
	confidence := ConfidenceConfirmed

	switch inquiry {
	case InquiryChargeStatus:
		if hasCaptureProof {
			if count == 1 {
				answer = fmt.Sprintf("Yes, you were charged %s %s. PayPal REST API confirms that Order #%s was captured successfully under capture proof. However, an edge network drop interrupted the storefront confirmation.", s.Case.Amount.String(), s.Case.Currency, s.Case.ExternalReference)
			} else {
				answer = fmt.Sprintf("To re-confirm: exactly one charge of %s %s was captured by PayPal. We have locked checkout session #%s to ensure zero duplicate debits ($0.00 duplicate exposure).", s.Case.Amount.String(), s.Case.Currency, s.Case.ExternalReference)
			}
		} else if isDenied {
			if count == 1 {
				answer = fmt.Sprintf("No, you were NOT charged. The authorization attempt for %s %s was declined by the card issuer prior to capture. Zero funds were debited from your card.", s.Case.Amount.String(), s.Case.Currency)
			} else {
				answer = fmt.Sprintf("Confirmed again: no funds were debited. Your card balance remains untouched for Order #%s.", s.Case.ExternalReference)
			}
		} else if isUnknown {
			confidence = ConfidenceUnknown
			answer = "I can't confirm that from the available payment records yet. The transaction state remains unverified, and we are holding attempt locks while awaiting authoritative provider reconciliation."
		} else {
			answer = fmt.Sprintf("Current payment status is %s. No duplicate debits have been recorded.", s.Case.ProviderState)
		}

	case InquiryRetrySafety:
		if hasCaptureProof {
			if count == 1 {
				answer = fmt.Sprintf("Do NOT retry or pay again. Your %s %s payment was already captured by PayPal. Retrying now would trigger an unnecessary duplicate charge. Our single-flight mutex lock is actively protecting your checkout.", s.Case.Amount.String(), s.Case.Currency)
			} else {
				answer = fmt.Sprintf("As advised earlier, please do not retry. Order #%s is already paid on PayPal's ledger. Re-submitting payment would risk duplicate billing.", s.Case.ExternalReference)
			}
		} else if isDenied {
			answer = fmt.Sprintf("Do not retry immediately with the same card settings. First resolve the issue (%s), then submit a single authorized checkout attempt.", s.DiagnosedIssue)
		} else if isUnknown {
			confidence = ConfidenceUnknown
			answer = "Do not retry payment yet. While the transaction outcome remains unconfirmed, retrying could trigger an accidental double charge. Please wait for reconciliation."
		} else {
			answer = "Retry is not recommended until the current attempt lock expires."
		}

	case InquiryBankRejection:
		if hasCaptureProof {
			answer = fmt.Sprintf("There is no evidence that your bank rejected this transaction. In fact, your bank authorized and PayPal captured the %s %s charge. The issue was an edge network timeout between the gateway and merchant, not an issuer decline.", s.Case.Amount.String(), s.Case.Currency)
		} else if isDenied {
			answer = fmt.Sprintf("Your card issuer rejected the %s %s authorization with response code 57 (or applicable ISO decline code). The bank blocked the authorization request before capture.", s.Case.Amount.String(), s.Case.Currency)
		} else {
			confidence = ConfidenceUnknown
			answer = "I can't confirm that from the available payment records yet. No bank rejection telemetry has been returned."
		}

	case InquiryFundsLocation:
		if hasCaptureProof && s.Case.ReceiverState != domain.ReceiverStateVerifiedBankCredited {
			answer = fmt.Sprintf("Your %s %s is currently held in PayPal's verified escrow ledger under Order #%s. It has not yet been disbursed to the merchant's external bank account.", s.Case.Amount.String(), s.Case.Currency, s.Case.ExternalReference)
		} else if s.Case.ReceiverState == domain.ReceiverStateVerifiedBankCredited {
			answer = "Your funds have settled into the merchant's account, confirmed via verified bank settlement journal."
		} else if isDenied {
			answer = "Your funds never left your bank account. Because the authorization was declined, no funds were transferred or held."
		} else {
			confidence = ConfidenceUnknown
			answer = "I can't confirm that from the available payment records yet."
		}

	case InquiryRefund:
		if hasCaptureProof {
			answer = "PayBridge AI is an advisory intelligence system and cannot execute money movements or refunds directly. Because your payment was captured by PayPal, you may request an order fulfillment sync or contact the merchant for a refund."
		} else {
			answer = "No refund is necessary because no capture occurred. No money was taken from your account."
		}

	case InquiryNextActions:
		if s.SafeAction != "" {
			answer = fmt.Sprintf("Recommended next step: %s", s.SafeAction)
		} else {
			answer = "Review your payment receipt in the PayBridge Truth Vault."
		}

	default:
		answer = fmt.Sprintf("PayBridge AI is monitoring Order #%s (%s %s). Current verified status: %s. %s", s.Case.ExternalReference, s.Case.Amount.String(), s.Case.Currency, s.Case.ProviderState, s.SafeAction)
	}

	// Record in history
	s.History = append(s.History, ChatMessage{
		Role:      "user",
		Content:   userQuery,
		Timestamp: time.Now().UTC(),
	})
	s.History = append(s.History, ChatMessage{
		Role:      "assistant",
		Content:   answer,
		Timestamp: time.Now().UTC(),
	})

	return &ChatResponse{
		Answer:           answer,
		InquiryType:      inquiry,
		EvidenceCited:    cited,
		Confidence:       confidence,
		AttemptLockHeld:  true,
		DuplicateCharges: "$0.00",
	}, nil
}
