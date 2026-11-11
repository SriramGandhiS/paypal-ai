package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// EvidenceTier represents the 6-tier evidence hierarchy
type EvidenceTier int

const (
	Tier1ProviderDirect   EvidenceTier = 1
	Tier2ProviderWebhook  EvidenceTier = 2
	Tier3NetworkTelemetry EvidenceTier = 3
	Tier4AppLedger        EvidenceTier = 4
	Tier5UserReported     EvidenceTier = 5
	Tier6Simulated        EvidenceTier = 6
)

// ProviderState represents provider authoritative status
type ProviderState string

const (
	ProviderStateNotCreated        ProviderState = "NOT_CREATED"
	ProviderStateCreated           ProviderState = "CREATED"
	ProviderStateApproved          ProviderState = "APPROVED"
	ProviderStateAuthorised        ProviderState = "AUTHORISED"
	ProviderStateCapturePending    ProviderState = "CAPTURE_PENDING"
	ProviderStateCaptureCompleted  ProviderState = "CAPTURE_COMPLETED"
	ProviderStateCaptureDenied     ProviderState = "CAPTURE_DENIED"
	ProviderStateRefunded          ProviderState = "REFUNDED"
	ProviderStatePartiallyRefunded ProviderState = "PARTIALLY_REFUNDED"
	ProviderStateUnknown           ProviderState = "UNKNOWN"
)

// ApplicationState represents PayBridge application and merchant observational state
type ApplicationState string

const (
	AppStateSynchronised       ApplicationState = "SYNCHRONISED"
	AppStateWaitingWebhook     ApplicationState = "WAITING_FOR_WEBHOOK"
	AppStateAPITimeout         ApplicationState = "API_TIMEOUT"
	AppStateAppCrashedPreCommit ApplicationState = "CRASHED_PRE_COMMIT"
	AppStateDuplicateAttempt   ApplicationState = "DUPLICATE_ATTEMPT_DETECTED"
	AppStateNetworkPartition   ApplicationState = "NETWORK_PARTITION"
	AppStateInconsistent       ApplicationState = "INCONSISTENT"
)

// ReceiverState represents receiver ledger status
type ReceiverState string

const (
	ReceiverStateCredited      ReceiverState = "CREDITED"
	ReceiverStatePending       ReceiverState = "PENDING"
	ReceiverStateNotObservable ReceiverState = "NOT_OBSERVABLE"
	ReceiverStateDisputed      ReceiverState = "DISPUTED"
	ReceiverStateVerifiedBankCredited ReceiverState = "VERIFIED_BANK_CREDITED"
)

// Common environment and action constants
const (
	EnvPayPalSimulator       = "SIMULATED_SANDBOX"
	EnvRealSandbox           = "REAL_SANDBOX"
	EnvProduction            = "PRODUCTION"
	ActionProceedToCheckout = "PROCEED_TO_CHECKOUT"
)

// PaymentCase represents a cross-border payment investigation entity
type PaymentCase struct {
	ID                  uuid.UUID        `json:"id"`
	TenantID            uuid.UUID        `json:"tenant_id"`
	ExternalReference   string           `json:"external_reference"`
	SenderUserID        string           `json:"sender_user_id"`
	SenderCountry       string           `json:"sender_country"`
	ReceiverCountry     string           `json:"receiver_country"`
	PaymentPurpose      string           `json:"payment_purpose"`
	Amount              decimal.Decimal  `json:"amount"`
	Currency            string           `json:"currency"`
	Environment         string           `json:"environment"`
	ProviderState       ProviderState    `json:"provider_state"`
	ApplicationState    ApplicationState `json:"application_state"`
	ReceiverState       ReceiverState    `json:"receiver_state"`
	CaseStatus          string           `json:"case_status"`
	StateVersion        int              `json:"state_version"`
	CreatedAt           time.Time        `json:"created_at"`
	UpdatedAt           time.Time        `json:"updated_at"`
}

// PaymentAttempt represents an individual attempt to execute checkout
type PaymentAttempt struct {
	ID             uuid.UUID `json:"id"`
	CaseID         uuid.UUID `json:"case_id"`
	AttemptNumber  int       `json:"attempt_number"`
	IdempotencyKey string    `json:"idempotency_key"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}

// EvidenceItem represents an immutable piece of financial provenance
type EvidenceItem struct {
	ID              uuid.UUID      `json:"id"`
	TenantID        uuid.UUID      `json:"tenant_id"`
	CaseID          uuid.UUID      `json:"case_id"`
	Kind            string         `json:"kind"`
	Tier            EvidenceTier   `json:"tier"`
	SourceTrust     string         `json:"source_trust"`
	Environment     string         `json:"environment"`
	Verified        bool           `json:"verified"`
	ObservedAt      time.Time      `json:"observed_at"`
	OccurredAt      *time.Time     `json:"occurred_at,omitempty"`
	RawPayloadSHA   string         `json:"raw_payload_sha"`
	PayloadSHA256   string         `json:"payload_sha256"`
	AccessClass     string         `json:"access_class"`
	RedactedExtract map[string]any `json:"redacted_extract,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
}
