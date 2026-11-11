package multirail

import (
	"context"
	"errors"
	"time"

	"github.com/shopspring/decimal"
)

var (
	ErrInvalidSignature       = errors.New("multirail: webhook signature verification failed")
	ErrWebhookReplay          = errors.New("multirail: duplicate or replayed webhook event rejected")
	ErrWebhookTimestampExpired= errors.New("multirail: webhook timestamp outside valid tolerance window")
	ErrUnauthorizedMutation   = errors.New("multirail: mutation rejected; requires explicit human authorization")
	ErrIdempotencyViolation   = errors.New("multirail: idempotency key conflict with differing request parameters")
	ErrProviderTimeout        = errors.New("multirail: downstream payment gateway timed out (outcome-unknown)")
)

type RailType string

const (
	RailPayPal RailType = "PAYPAL"
	RailStripe RailType = "STRIPE"
	RailWise   RailType = "WISE"
)

type DisputeStatus string

const (
	DisputeStatusNeedsEvidence DisputeStatus = "NEEDS_EVIDENCE"
	DisputeStatusUnderReview   DisputeStatus = "UNDER_REVIEW"
	DisputeStatusWon           DisputeStatus = "WON"
	DisputeStatusLost          DisputeStatus = "LOST"
)

type RailTransaction struct {
	Rail           RailType        `json:"rail"`
	TransactionID  string          `json:"transaction_id"`
	Amount         decimal.Decimal `json:"amount"`
	Currency       string          `json:"currency"`
	Status         string          `json:"status"`
	BankClearRef   string          `json:"bank_clear_ref,omitempty"` // Real UETR / IMAD / Fedwire
	OccurredAt     time.Time       `json:"occurred_at"`
	RawPayloadHash string          `json:"raw_payload_hash"`
	IdempotencyKey string          `json:"idempotency_key,omitempty"`
}

type DisputeRecord struct {
	DisputeID     string          `json:"dispute_id"`
	TransactionID string          `json:"transaction_id"`
	Rail          RailType        `json:"rail"`
	Amount        decimal.Decimal `json:"amount"`
	Currency      string          `json:"currency"`
	Reason        string          `json:"reason"`
	Status        DisputeStatus   `json:"status"`
	EvidenceDueBy time.Time       `json:"evidence_due_by"`
}

// HumanAuthorization carries operator approval tokens
type HumanAuthorization struct {
	OperatorID string
	Reason     string
	Timestamp  time.Time
}

// RailReader provides strictly read-only inspection across all payment rails
// AI agents, background reconciliation, and ops telemetry must only access this interface
type RailReader interface {
	GetTransaction(ctx context.Context, txID string) (*RailTransaction, error)
	GetDispute(ctx context.Context, disputeID string) (*DisputeRecord, error)
}

// RailMutator handles financial mutations (Charge, Capture, Refund, Payout)
// STRICT NON-NEGOTIABLE RULE: AI MUST NEVER ACCESS THIS INTERFACE.
// Mutator calls require an explicit human operator authorization token.
type RailMutator interface {
	InitiatePayment(ctx context.Context, auth HumanAuthorization, amount decimal.Decimal, currency, reference string) (*RailTransaction, error)
	CapturePayment(ctx context.Context, auth HumanAuthorization, txID string, amount decimal.Decimal) (*RailTransaction, error)
	IssueRefund(ctx context.Context, auth HumanAuthorization, txID string, amount decimal.Decimal, reason string) (*RailTransaction, error)
}

// WebhookVerifier authenticates incoming raw webhooks before JSON decoding
type WebhookVerifier interface {
	VerifyWebhook(rawBody []byte, signatureHeader string) (providerEventID string, err error)
}
