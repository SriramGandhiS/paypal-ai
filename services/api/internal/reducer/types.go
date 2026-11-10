package reducer

import (
	"time"
)

// ProviderState represents PayPal or provider authoritative status
type ProviderState string

const (
	ProviderStateNotCreated       ProviderState = "NOT_CREATED"
	ProviderStateCreated          ProviderState = "CREATED"
	ProviderStateApproved         ProviderState = "APPROVED"
	ProviderStateAuthorised       ProviderState = "AUTHORISED"
	ProviderStateCapturePending   ProviderState = "CAPTURE_PENDING"
	ProviderStateCaptureCompleted ProviderState = "CAPTURE_COMPLETED"
	ProviderStateCaptureDenied    ProviderState = "CAPTURE_DENIED"
	ProviderStateRefunded         ProviderState = "REFUNDED"
	ProviderStatePartiallyRefunded ProviderState = "PARTIALLY_REFUNDED"
	ProviderStateUnknown          ProviderState = "UNKNOWN"
)

// ApplicationState represents PayBridge application and merchant observational state
type ApplicationState string

const (
	AppStateSynchronised     ApplicationState = "SYNCHRONISED"
	AppStateWaitingWebhook   ApplicationState = "WAITING_FOR_WEBHOOK"
	AppStateAPITimeout       ApplicationState = "API_TIMEOUT"
	AppStateReconciling      ApplicationState = "RECONCILING"
	AppStateEventGap         ApplicationState = "EVENT_GAP"
	AppStateInconsistent     ApplicationState = "INCONSISTENT"
)

// ReceiverState represents bank settlement and recipient funds observability
type ReceiverState string

const (
	ReceiverStateNotApplicable       ReceiverState = "NOT_APPLICABLE"
	ReceiverStateNotObservable       ReceiverState = "NOT_OBSERVABLE"
	ReceiverStateReportedPending     ReceiverState = "REPORTED_PENDING"
	ReceiverStateVerifiedAvailable   ReceiverState = "VERIFIED_AVAILABLE"
	ReceiverStateWithdrawalInProgress ReceiverState = "WITHDRAWAL_IN_PROGRESS"
	ReceiverStateVerifiedBankCredited ReceiverState = "VERIFIED_BANK_CREDITED"
)

// CaseStatus represents PayBridge workflow state
type CaseStatus string

const (
	CaseStatusOpen                    CaseStatus = "OPEN"
	CaseStatusObserving               CaseStatus = "OBSERVING"
	CaseStatusInvestigating           CaseStatus = "INVESTIGATING"
	CaseStatusGuidanceReady           CaseStatus = "GUIDANCE_READY"
	CaseStatusAwaitingUserAction      CaseStatus = "AWAITING_USER_ACTION"
	CaseStatusVerifyingResolution     CaseStatus = "VERIFYING_RESOLUTION"
	CaseStatusResolved                CaseStatus = "RESOLVED"
	CaseStatusNeedsHumanReview        CaseStatus = "NEEDS_HUMAN_REVIEW"
	CaseStatusPausedEvidence          CaseStatus = "PAUSED_INSUFFICIENT_EVIDENCE"
	CaseStatusClosedUnresolved        CaseStatus = "CLOSED_UNRESOLVED"
	CaseStatusArchived                CaseStatus = "ARCHIVED"
	CaseStatusReopened                CaseStatus = "REOPENED"
)

// EvidenceTier represents the hierarchy of evidence reliability (E1 to E6)
type EvidenceTier int

const (
	TierE1DirectProviderRead EvidenceTier = 1 // Provider GET (TLS verified, direct request)
	TierE2VerifiedWebhook    EvidenceTier = 2 // Cryptographically signed provider webhook
	TierE3LocalAppTrace      EvidenceTier = 3 // Local HTTP client/server request/response traces
	TierE4MerchantRecord     EvidenceTier = 4 // Merchant DB order/invoice records
	TierE5UserReport         EvidenceTier = 5 // Free-text statement, buyer report (untrusted)
	TierE6Simulation         EvidenceTier = 6 // Synthetic/local simulation fixture
)

// SourceTrust represents source provenance
type SourceTrust string

const (
	SourceTrustProviderVerified SourceTrust = "PROVIDER_VERIFIED"
	SourceTrustLocalObserved    SourceTrust = "LOCAL_OBSERVED"
	SourceTrustPublicSource     SourceTrust = "PUBLIC_SOURCE"
	SourceTrustUserReported     SourceTrust = "USER_REPORTED"
	SourceTrustSimulated        SourceTrust = "SIMULATED"
)

// StateVector encapsulates the authoritative 3D financial state plus case workflow
type StateVector struct {
	ProviderState              ProviderState    `json:"provider_state"`
	LastConfirmedProviderState ProviderState    `json:"last_confirmed_provider_state"`
	ProviderStateEvidenceID    string           `json:"provider_state_evidence_id,omitempty"`
	ProviderStateObservedAt    time.Time        `json:"provider_state_observed_at,omitempty"`
	ApplicationState           ApplicationState `json:"application_state"`
	ReceiverState              ReceiverState    `json:"receiver_state"`
	CaseStatus                 CaseStatus       `json:"case_status"`
	StateVersion               int64            `json:"state_version"`
}

// Event represents an observation ingested into the Reducer
type Event struct {
	EventID        string       `json:"event_id"`
	EventType      string       `json:"event_type"`
	EvidenceID     string       `json:"evidence_id"`
	Tier           EvidenceTier `json:"tier"`
	SourceTrust    SourceTrust  `json:"source_trust"`
	Environment    string       `json:"environment"`
	StatusObserved string       `json:"status_observed"`
	OccurredAt     time.Time    `json:"occurred_at"`
	ObservedAt     time.Time    `json:"observed_at"`
	Payload        map[string]any `json:"payload,omitempty"`
}

// StateTransition records a deterministic state change
type StateTransition struct {
	Dimension string `json:"dimension"`
	From      string `json:"from"`
	To        string `json:"to"`
	Reason    string `json:"reason"`
	EventID   string `json:"event_id"`
}
