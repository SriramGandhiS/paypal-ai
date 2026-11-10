package reducer

import (
	"context"
	"testing"
	"time"
)

func TestReducer_TimeoutDoesNotEqualFailure(t *testing.T) {
	r := NewReducer()
	ctx := context.Background()

	initial := StateVector{
		ProviderState:              ProviderStateApproved,
		LastConfirmedProviderState: ProviderStateApproved,
		ApplicationState:           AppStateSynchronised,
		ReceiverState:              ReceiverStateNotObservable,
		CaseStatus:                 CaseStatusOpen,
		StateVersion:               1,
	}

	timeoutEvent := Event{
		EventID:        "evt-001",
		EventType:      "CHECKOUT_TIMEOUT",
		EvidenceID:     "ev-001",
		Tier:           TierE3LocalAppTrace,
		SourceTrust:    SourceTrustLocalObserved,
		Environment:    "REAL_SANDBOX",
		StatusObserved: "TIMEOUT",
		OccurredAt:     time.Now(),
		ObservedAt:     time.Now(),
	}

	next, transitions, err := r.Apply(ctx, initial, timeoutEvent)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Invariant: ProviderState becomes UNKNOWN, NEVER CAPTURE_DENIED
	if next.ProviderState != ProviderStateUnknown {
		t.Errorf("expected provider_state UNKNOWN, got %s", next.ProviderState)
	}
	if next.LastConfirmedProviderState != ProviderStateApproved {
		t.Errorf("expected last_confirmed_provider_state APPROVED, got %s", next.LastConfirmedProviderState)
	}
	if next.ApplicationState != AppStateAPITimeout {
		t.Errorf("expected application_state API_TIMEOUT, got %s", next.ApplicationState)
	}
	if next.CaseStatus != CaseStatusInvestigating {
		t.Errorf("expected case_status INVESTIGATING, got %s", next.CaseStatus)
	}
	if len(transitions) == 0 {
		t.Errorf("expected transitions recorded, got 0")
	}
	if next.StateVersion != 2 {
		t.Errorf("expected state_version incremented to 2, got %d", next.StateVersion)
	}
}

func TestReducer_ReconcileAfterTimeout(t *testing.T) {
	r := NewReducer()
	ctx := context.Background()

	// State after timeout
	state := StateVector{
		ProviderState:              ProviderStateUnknown,
		LastConfirmedProviderState: ProviderStateApproved,
		ApplicationState:           AppStateAPITimeout,
		ReceiverState:              ReceiverStateNotObservable,
		CaseStatus:                 CaseStatusInvestigating,
		StateVersion:               2,
	}

	// 1. Worker begins reconciliation
	reconcileStartEvent := Event{
		EventID:     "evt-rec-start",
		EventType:   "RECONCILIATION_STARTED",
		EvidenceID:  "ev-rec-1",
		Tier:        TierE3LocalAppTrace,
		Environment: "REAL_SANDBOX",
		ObservedAt:  time.Now(),
	}
	state, _, err := r.Apply(ctx, state, reconcileStartEvent)
	if err != nil {
		t.Fatalf("reconcile start error: %v", err)
	}
	if state.ApplicationState != AppStateReconciling {
		t.Errorf("expected RECONCILING, got %s", state.ApplicationState)
	}

	// 2. Direct provider GET confirms capture completed (E1 Tier)
	captureConfirmedEvent := Event{
		EventID:        "evt-capture-get",
		EventType:      "PROVIDER_CAPTURE_COMPLETED",
		EvidenceID:     "ev-paypal-cap-123",
		Tier:           TierE1DirectProviderRead,
		SourceTrust:    SourceTrustProviderVerified,
		Environment:    "REAL_SANDBOX",
		StatusObserved: "COMPLETED",
		ObservedAt:     time.Now(),
	}
	state, transitions, err := r.Apply(ctx, state, captureConfirmedEvent)
	if err != nil {
		t.Fatalf("capture confirmed error: %v", err)
	}

	if state.ProviderState != ProviderStateCaptureCompleted {
		t.Errorf("expected CAPTURE_COMPLETED, got %s", state.ProviderState)
	}
	if state.ApplicationState != AppStateSynchronised {
		t.Errorf("expected SYNCHRONISED, got %s", state.ApplicationState)
	}
	if state.ProviderStateEvidenceID != "ev-paypal-cap-123" {
		t.Errorf("expected evidence ID linked, got %s", state.ProviderStateEvidenceID)
	}
	if len(transitions) == 0 {
		t.Errorf("expected transitions recorded")
	}
}

func TestReducer_MonotonicityGuards(t *testing.T) {
	r := NewReducer()
	ctx := context.Background()

	completedState := StateVector{
		ProviderState:              ProviderStateCaptureCompleted,
		LastConfirmedProviderState: ProviderStateCaptureCompleted,
		ApplicationState:           AppStateSynchronised,
		StateVersion:               5,
	}

	// Attempting to deny an already completed capture violates monotonicity
	denialEvent := Event{
		EventID:     "evt-deny",
		EventType:   "PROVIDER_CAPTURE_DENIED",
		EvidenceID:  "ev-denial",
		Tier:        TierE1DirectProviderRead,
		Environment: "REAL_SANDBOX",
		ObservedAt:  time.Now(),
	}

	_, _, err := r.Apply(ctx, completedState, denialEvent)
	if err == nil {
		t.Errorf("expected error on attempting to deny a completed capture, got nil")
	}
}

func TestReducer_SimulatedEvidenceCannotModifyRealCase(t *testing.T) {
	r := NewReducer()
	ctx := context.Background()

	realState := StateVector{
		ProviderState:              ProviderStateApproved,
		LastConfirmedProviderState: ProviderStateApproved,
		ApplicationState:           AppStateSynchronised,
		StateVersion:               3,
	}

	simulatedEvent := Event{
		EventID:     "evt-sim",
		EventType:   "PROVIDER_CAPTURE_COMPLETED",
		EvidenceID:  "ev-sim-1",
		Tier:        TierE6Simulation,
		SourceTrust: SourceTrustSimulated,
		Environment: "REAL_SANDBOX", // Illegitimate mix
		ObservedAt:  time.Now(),
	}

	_, _, err := r.Apply(ctx, realState, simulatedEvent)
	if err != ErrSimulatedEventOnRealCase {
		t.Errorf("expected ErrSimulatedEventOnRealCase, got %v", err)
	}
}

func TestReducer_OutOfOrderEvents(t *testing.T) {
	r := NewReducer()
	ctx := context.Background()

	// State is already CAPTURE_COMPLETED
	completedState := StateVector{
		ProviderState:              ProviderStateCaptureCompleted,
		LastConfirmedProviderState: ProviderStateCaptureCompleted,
		ApplicationState:           AppStateSynchronised,
		StateVersion:               5,
	}

	// An out-of-order stale event arrives: BUYER_APPROVED_ORDER
	staleApprovalEvent := Event{
		EventID:     "evt-stale-approval",
		EventType:   "BUYER_APPROVED_ORDER",
		EvidenceID:  "ev-stale-appr",
		Tier:        TierE1DirectProviderRead,
		Environment: "REAL_SANDBOX",
		ObservedAt:  time.Now(),
	}

	next, transitions, err := r.Apply(ctx, completedState, staleApprovalEvent)
	if err != nil {
		t.Fatalf("unexpected error on out-of-order event: %v", err)
	}

	// State MUST NOT regress from CAPTURE_COMPLETED to APPROVED
	if next.ProviderState != ProviderStateCaptureCompleted {
		t.Errorf("expected provider_state to remain CAPTURE_COMPLETED, got: %s", next.ProviderState)
	}
	if len(transitions) != 0 {
		t.Errorf("expected 0 transitions for stale event, got %d", len(transitions))
	}
	if next.StateVersion != completedState.StateVersion {
		t.Errorf("expected state version to remain unchanged, got %d", next.StateVersion)
	}
}

func TestReducer_DuplicateEvents(t *testing.T) {
	r := NewReducer()
	ctx := context.Background()

	completedState := StateVector{
		ProviderState:              ProviderStateCaptureCompleted,
		LastConfirmedProviderState: ProviderStateCaptureCompleted,
		ApplicationState:           AppStateSynchronised,
		StateVersion:               5,
	}

	duplicateCaptureEvent := Event{
		EventID:     "evt-dup-cap",
		EventType:   "PAYMENT.CAPTURE.COMPLETED",
		EvidenceID:  "ev-dup-cap",
		Tier:        TierE2VerifiedWebhook,
		Environment: "REAL_SANDBOX",
		ObservedAt:  time.Now(),
	}

	next, transitions, err := r.Apply(ctx, completedState, duplicateCaptureEvent)
	if err != nil {
		t.Fatalf("unexpected error on duplicate event: %v", err)
	}

	if next.ProviderState != ProviderStateCaptureCompleted {
		t.Errorf("expected provider_state to remain CAPTURE_COMPLETED, got: %s", next.ProviderState)
	}
	if len(transitions) != 0 {
		t.Errorf("expected 0 transitions for duplicate event, got %d", len(transitions))
	}
}

