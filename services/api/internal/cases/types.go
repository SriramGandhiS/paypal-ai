package cases

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalidAmount   = errors.New("amount must be greater than zero")
	ErrInvalidCurrency = errors.New("currency must be a valid 3-letter ISO code")
	ErrInvalidCountry  = errors.New("country must be a valid 2-letter ISO code")
	ErrCaseNotFound    = errors.New("payment case not found")
)

// CreateCaseInput represents the parameters required to initialize a payment case.
type CreateCaseInput struct {
	ExternalReference string `json:"external_reference"`
	AmountCents       int64  `json:"amount_cents"`
	Currency          string `json:"currency"`
	SenderCountry     string `json:"sender_country"`
	ReceiverCountry   string `json:"receiver_country"`
	PaymentPurpose    string `json:"payment_purpose"`
	IdempotencyKey    string `json:"idempotency_key"`
}

// Validate performs fail-fast domain validation on creation input.
func (in *CreateCaseInput) Validate() error {
	if in.AmountCents <= 0 {
		return ErrInvalidAmount
	}
	in.Currency = strings.ToUpper(strings.TrimSpace(in.Currency))
	if len(in.Currency) != 3 {
		return ErrInvalidCurrency
	}
	in.SenderCountry = strings.ToUpper(strings.TrimSpace(in.SenderCountry))
	if len(in.SenderCountry) != 2 {
		return ErrInvalidCountry
	}
	in.ReceiverCountry = strings.ToUpper(strings.TrimSpace(in.ReceiverCountry))
	if len(in.ReceiverCountry) != 2 {
		return ErrInvalidCountry
	}
	return nil
}

// Case represents a full payment reconciliation case record.
type Case struct {
	ID                         string    `json:"id"`
	TenantID                   string    `json:"tenant_id"`
	ExternalReference          string    `json:"external_reference"`
	SenderUserID               string    `json:"sender_user_id"`
	ReceiverUserID             *string   `json:"receiver_user_id,omitempty"`
	SenderCountry              string    `json:"sender_country"`
	ReceiverCountry            string    `json:"receiver_country"`
	PaymentPurpose             string    `json:"payment_purpose"`
	Amount                     string    `json:"amount"`
	Currency                   string    `json:"currency"`
	Environment                string    `json:"environment"`
	ProviderState              string    `json:"provider_state"`
	LastConfirmedProviderState string    `json:"last_confirmed_provider_state"`
	ApplicationState           string    `json:"application_state"`
	ReceiverState              string    `json:"receiver_state"`
	CaseStatus                 string    `json:"case_status"`
	StateVersion               int64     `json:"state_version"`
	CreatedAt                  time.Time `json:"created_at"`
	UpdatedAt                  time.Time `json:"updated_at"`
}

// CaseSummary is a compact representation for list views and timelines.
type CaseSummary struct {
	ID                string    `json:"id"`
	ExternalReference string    `json:"external_reference"`
	Amount            string    `json:"amount"`
	Currency          string    `json:"currency"`
	ProviderState     string    `json:"provider_state"`
	ApplicationState  string    `json:"application_state"`
	CaseStatus        string    `json:"case_status"`
	StateVersion      int64     `json:"state_version"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// ListFilter specifies filtering and pagination criteria.
type ListFilter struct {
	Status string
	Limit  int
	Offset int
}
