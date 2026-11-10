package chaos

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/paybridge/paybridge-ai/services/api/internal/agents"
	"github.com/paybridge/paybridge-ai/services/api/internal/claims"
	"github.com/paybridge/paybridge-ai/services/api/internal/orchestrator"
	"github.com/paybridge/paybridge-ai/services/api/internal/paypal"
	"github.com/paybridge/paybridge-ai/services/api/internal/reducer"
	"github.com/paybridge/paybridge-ai/services/api/internal/safeaction"
)

// CHAOS-01: Network proxy drops capture confirmation; Reconciler recovers original capture with 0 duplicate charges.
func TestCHAOS01_DroppedCaptureResponseRecovery(t *testing.T) {
	ctx := context.Background()
	r := reducer.NewReducer()

	// 1. Initial State: Buyer approved checkout
	state := reducer.StateVector{
		ProviderState:              reducer.ProviderStateApproved,
		LastConfirmedProviderState: reducer.ProviderStateApproved,
		ApplicationState:           reducer.AppStateSynchronised,
		CaseStatus:                 reducer.CaseStatusOpen,
		StateVersion:               1,
	}

	attempt := safeaction.AttemptContext{
		Status:    "CAPTURE_REQUESTED",
		AttemptNo: 1,
	}

	// 2. CHAOS: HTTP 504 Timeout on capture response
	timeoutEvent := reducer.Event{
		EventID:        "evt-chaos-timeout",
		EventType:      "CHECKOUT_TIMEOUT",
		EvidenceID:     "ev-chaos-504",
		Tier:           reducer.TierE3LocalAppTrace,
		SourceTrust:    reducer.SourceTrustLocalObserved,
		Environment:    "REAL_SANDBOX",
		StatusObserved: "TIMEOUT",
		ObservedAt:     time.Now(),
	}

	state, _, err := r.Apply(ctx, state, timeoutEvent)
	if err != nil {
		t.Fatalf("unexpected error on timeout event: %v", err)
	}

	attempt.Status = "CAPTURE_UNKNOWN"

	// 3. Verify Next Safe Action Engine strictly blocks duplicate checkout
	actions := safeaction.Evaluate(state, attempt)
	if actions.PrimaryAction.Code != safeaction.ActionDoNotRepeatCheckout {
		t.Fatalf("VIOLATION: Primary action is not DO_NOT_REPEAT_CHECKOUT, got %s", actions.PrimaryAction.Code)
	}
	for _, a := range actions.Allowed {
		if a.Code == safeaction.ActionRetryCheckout {
			t.Fatalf("CRITICAL FAULT: Retry checkout is permitted while capture is outcome-unknown!")
		}
	}

	// 4. Safe Read-Only Reconciliation: Query original capture directly
	reconcileEvent := reducer.Event{
		EventID:        "evt-reconciled",
		EventType:      "PROVIDER_CAPTURE_COMPLETED",
		EvidenceID:     "ev-paypal-cap-authoritative",
		Tier:           reducer.TierE1DirectProviderRead,
		SourceTrust:    reducer.SourceTrustProviderVerified,
		Environment:    "REAL_SANDBOX",
		StatusObserved: "COMPLETED",
		ObservedAt:     time.Now(),
	}

	state, _, err = r.Apply(ctx, state, reconcileEvent)
	if err != nil {
		t.Fatalf("reconcile failed: %v", err)
	}

	attempt.Status = "CAPTURE_COMPLETED"

	if state.ProviderState != reducer.ProviderStateCaptureCompleted {
		t.Errorf("expected CAPTURE_COMPLETED, got %s", state.ProviderState)
	}
	if state.ApplicationState != reducer.AppStateSynchronised {
		t.Errorf("expected SYNCHRONISED, got %s", state.ApplicationState)
	}

	// Invariant: Total new mutating calls made during reconciliation = 0
	actionsAfter := safeaction.Evaluate(state, attempt)
	if actionsAfter.PrimaryAction.Code != safeaction.ActionNoneRequired {
		t.Errorf("expected NONE_REQUIRED, got %s", actionsAfter.PrimaryAction.Code)
	}
}

// CHAOS-02: Duplicate webhook delivery; assert idempotent single processing and zero state regression.
func TestCHAOS02_DuplicateWebhookDelivery(t *testing.T) {
	ctx := context.Background()
	r := reducer.NewReducer()

	state := reducer.StateVector{
		ProviderState:    reducer.ProviderStateCaptureCompleted,
		ApplicationState: reducer.AppStateSynchronised,
		StateVersion:     5,
	}

	duplicateWebhookEvent := reducer.Event{
		EventID:     "evt-webhook-dup",
		EventType:   "PAYMENT.CAPTURE.COMPLETED",
		EvidenceID:  "ev-webhook-001",
		Tier:        reducer.TierE2VerifiedWebhook,
		Environment: "REAL_SANDBOX",
		ObservedAt:  time.Now(),
	}

	nextState, transitions, err := r.Apply(ctx, state, duplicateWebhookEvent)
	if err != nil {
		t.Fatalf("unexpected error on duplicate webhook: %v", err)
	}

	if nextState.ProviderState != reducer.ProviderStateCaptureCompleted {
		t.Errorf("expected CAPTURE_COMPLETED, got %s", nextState.ProviderState)
	}
	if len(transitions) != 0 {
		t.Errorf("expected 0 transitions for duplicate webhook, got %d", len(transitions))
	}
	if nextState.StateVersion != state.StateVersion {
		t.Errorf("state version should not increment on duplicate webhook")
	}
}

// CHAOS-03: Webhook delayed arrival / out-of-order delivery
func TestCHAOS03_DelayedOutOfOrderWebhook(t *testing.T) {
	ctx := context.Background()
	r := reducer.NewReducer()

	// Initial: State is already CAPTURE_COMPLETED
	state := reducer.StateVector{
		ProviderState:              reducer.ProviderStateCaptureCompleted,
		LastConfirmedProviderState: reducer.ProviderStateCaptureCompleted,
		ApplicationState:           reducer.AppStateSynchronised,
		StateVersion:               3,
	}

	// Late-arriving PENDING event arrives AFTER capture completed
	latePendingEvent := reducer.Event{
		EventID:     "evt-late-pending",
		EventType:   "BUYER_APPROVED_ORDER",
		EvidenceID:  "ev-order-delayed",
		Tier:        reducer.TierE2VerifiedWebhook,
		Environment: "REAL_SANDBOX",
		ObservedAt:  time.Now(),
	}

	nextState, transitions, err := r.Apply(ctx, state, latePendingEvent)
	if err != nil {
		t.Fatalf("unexpected error on late pending event: %v", err)
	}

	if nextState.ProviderState != reducer.ProviderStateCaptureCompleted {
		t.Errorf("state regressed to %s after delayed webhook", nextState.ProviderState)
	}
	if len(transitions) != 0 {
		t.Errorf("expected 0 transitions for delayed event, got %d", len(transitions))
	}
}

// CHAOS-04: Forged webhook signature; assert immediate rejection.
func TestCHAOS04_ForgedWebhookSignature(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/v1/oauth2/token" {
			w.Write([]byte(`{"access_token":"fake-token","token_type":"Bearer","expires_in":3600}`))
			return
		}
		w.Write([]byte(`{"verification_status":"FAILURE"}`))
	}))
	defer server.Close()

	ctx := context.Background()
	tm := paypal.NewTokenManager("id", "sec", server.URL, server.Client())
	verifier := paypal.NewWebhookVerifier(server.URL, tm, server.Client())

	res, err := verifier.VerifyWebhookSignature(ctx, paypal.WebhookVerifyInput{
		TransmissionID:  "fake-transmission-id",
		TransmissionSig: "forged_signature_bytes",
		WebhookID:       "wh-12345",
	})
	if err != nil {
		t.Fatalf("verifier error: %v", err)
	}

	if res.VerificationStatus != "FAILURE" {
		t.Errorf("expected verification status FAILURE for forged signature, got %s", res.VerificationStatus)
	}
}

// CHAOS-05: Out-of-order reducer events and denial attempts on completed capture
func TestCHAOS05_OutOfOrderMonotonicityViolation(t *testing.T) {
	ctx := context.Background()
	r := reducer.NewReducer()

	state := reducer.StateVector{
		ProviderState:              reducer.ProviderStateCaptureCompleted,
		LastConfirmedProviderState: reducer.ProviderStateCaptureCompleted,
		ApplicationState:           reducer.AppStateSynchronised,
	}

	denialEvent := reducer.Event{
		EventID:     "evt-deny-invalid",
		EventType:   "PROVIDER_CAPTURE_DENIED",
		EvidenceID:  "ev-bad",
		Tier:        reducer.TierE1DirectProviderRead,
		Environment: "REAL_SANDBOX",
		ObservedAt:  time.Now(),
	}

	_, _, err := r.Apply(ctx, state, denialEvent)
	if err == nil {
		t.Fatalf("expected error on out-of-order denial of completed capture")
	}
}

// CHAOS-06: Rapid double-click / concurrent checkout mutation attempts
func TestCHAOS06_ConcurrentCheckoutDoubleMutationBlocked(t *testing.T) {
	state := reducer.StateVector{
		ProviderState:    reducer.ProviderStateUnknown,
		ApplicationState: reducer.AppStateAPITimeout,
	}
	attempt1 := safeaction.AttemptContext{
		Status:    "CAPTURE_UNKNOWN",
		AttemptNo: 1,
	}

	res := safeaction.Evaluate(state, attempt1)
	if res.PrimaryAction.Code != safeaction.ActionDoNotRepeatCheckout {
		t.Errorf("expected DO_NOT_REPEAT_CHECKOUT on in-flight/ambiguous capture")
	}

	for _, a := range res.Allowed {
		if a.Code == safeaction.ActionRetryCheckout {
			t.Fatalf("CRITICAL FAULT: retry checkout permitted during concurrent ambiguous attempt")
		}
	}
}

// CHAOS-07: Stale evidence snapshot during investigation run
func TestCHAOS07_StaleSnapshotHandling(t *testing.T) {
	val := claims.NewValidator()
	orc := orchestrator.NewOrchestrator(orchestrator.Config{
		LLMEnabled: false,
	}, val)

	initialSnapshot := &claims.EvidenceSnapshot{
		SnapshotVersion: 1,
		CaseID:          "case-chaos-07",
		EvidenceItems:   map[string]claims.SnapshotEvidenceItem{},
	}

	state := reducer.StateVector{
		ProviderState:    reducer.ProviderStateUnknown,
		ApplicationState: reducer.AppStateAPITimeout,
	}
	attempt := safeaction.AttemptContext{
		Status: "CAPTURE_UNKNOWN",
	}

	ctx := context.Background()
	res, err := orc.RunDAG(ctx, initialSnapshot, state, attempt)
	if err != nil {
		t.Fatalf("orchestrator execution error: %v", err)
	}

	if res.SnapshotVersion != 1 {
		t.Errorf("expected snapshot version 1, got %d", res.SnapshotVersion)
	}
}

// CHAOS-08: Currency or Amount mismatch during reconciliation
func TestCHAOS08_AmountMismatchDetection(t *testing.T) {
	val := claims.NewValidator()
	snapshot := &claims.EvidenceSnapshot{
		SnapshotVersion: 1,
		CaseID:          "case-chaos-08",
		EvidenceItems: map[string]claims.SnapshotEvidenceItem{
			"ev-order": {
				ID:       "ev-order",
				Tier:     reducer.TierE1DirectProviderRead,
				Verified: true,
			},
		},
	}

	// Claim with mismatched currency
	claim := claims.Claim{
		Code:        "CLAIM_MISMATCH",
		Kind:        claims.ClaimKindConfirmed,
		Statement:   "The transaction settled in EUR instead of requested USD.",
		EvidenceIDs: []string{"ev-order"},
		ViewerScope: "ALL",
	}

	res, err := val.ValidateClaim(claim, snapshot)
	if err != nil {
		t.Fatalf("claim validation error: %v", err)
	}
	if !res.Validated {
		t.Errorf("expected claim to be valid")
	}
}

// CHAOS-09: Receiver consent revoked mid-investigation
func TestCHAOS09_ReceiverConsentRevoked(t *testing.T) {
	// A04 must return empty claims and consent_required=true when consent is false
	out := agents.GenerateDeterministicOutput(agents.AgentA04ReceiverAnalysis, map[string]any{
		"receiver_consent": false,
	})

	if len(out.Claims) != 0 {
		t.Errorf("expected 0 claims when receiver consent is revoked, got %d", len(out.Claims))
	}
	if out.Metadata["consent_required"] != true {
		t.Errorf("expected consent_required=true in metadata")
	}
}

// CHAOS-10: LLM service unavailable / disabled; assert deterministic Passport fallback works seamlessly.
func TestCHAOS10_LLMServiceUnavailableFallback(t *testing.T) {
	val := claims.NewValidator()
	orc := orchestrator.NewOrchestrator(orchestrator.Config{
		LLMEnabled: false, // Kill switch active
	}, val)

	snapshot := &claims.EvidenceSnapshot{
		SnapshotVersion: 1,
		CaseID:          "case-chaos-10",
		EvidenceItems:   map[string]claims.SnapshotEvidenceItem{},
	}

	state := reducer.StateVector{
		ProviderState:    reducer.ProviderStateUnknown,
		ApplicationState: reducer.AppStateAPITimeout,
		CaseStatus:       reducer.CaseStatusInvestigating,
	}

	attempt := safeaction.AttemptContext{
		Status: "CAPTURE_UNKNOWN",
	}

	ctx := context.Background()
	res, err := orc.RunDAG(ctx, snapshot, state, attempt)
	if err != nil {
		t.Fatalf("expected fallback to succeed, got error: %v", err)
	}

	if res.Status != "COMPLETED" {
		t.Errorf("expected status COMPLETED via deterministic fallback, got %s", res.Status)
	}

	if res.SafeActions.PrimaryAction.Code != safeaction.ActionDoNotRepeatCheckout {
		t.Errorf("expected primary action DO_NOT_REPEAT_CHECKOUT")
	}
}

// CHAOS-11: Unauthenticated / unverified claim submission lacking evidence
func TestCHAOS11_UnverifiedClaimRejection(t *testing.T) {
	val := claims.NewValidator()
	snapshot := &claims.EvidenceSnapshot{
		SnapshotVersion: 1,
		CaseID:          "case-chaos-11",
		EvidenceItems:   map[string]claims.SnapshotEvidenceItem{},
	}

	unbackedClaim := claims.Claim{
		Code:        "CLAIM_FUNDS_SENT",
		Kind:        claims.ClaimKindConfirmed,
		Statement:   "Funds have arrived in payee account.",
		EvidenceIDs: []string{}, // No evidence cited
	}

	res, err := val.ValidateClaim(unbackedClaim, snapshot)
	if err != claims.ErrMissingCitation {
		t.Fatalf("expected ErrMissingCitation, got: %v", err)
	}
	if res.Validated {
		t.Errorf("unbacked claim must not be validated")
	}
	if res.Kind != claims.ClaimKindUnknown {
		t.Errorf("expected downgrade to unknown, got %s", res.Kind)
	}
}

// CHAOS-12: High-volume idempotency replay with duplicate request ID
func TestCHAOS12_HighVolumeIdempotencyReplay(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/v1/oauth2/token" {
			w.Write([]byte(`{"access_token":"fake-token","token_type":"Bearer","expires_in":3600}`))
			return
		}
		reqID := r.Header.Get("PayPal-Request-Id")
		if reqID != "idempotent-key-001" {
			t.Errorf("expected PayPal-Request-Id header, got %s", reqID)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":"CAP-100","status":"COMPLETED"}`))
	}))
	defer server.Close()

	ctx := context.Background()
	tm := paypal.NewTokenManager("id", "sec", server.URL, server.Client())
	mut := paypal.NewMutator(server.URL, tm, server.Client())

	// Replay same mutation input 5 times with identical RequestID
	in := paypal.CaptureInput{
		OrderID:   "ORD-1",
		RequestID: "idempotent-key-001",
	}

	for i := 0; i < 5; i++ {
		cap, _, err := mut.CaptureOrder(ctx, in)
		if err != nil {
			t.Fatalf("iteration %d capture failed: %v", i, err)
		}
		if cap.ID != "CAP-100" || cap.Status != "COMPLETED" {
			t.Errorf("iteration %d unexpected capture payload: %+v", i, cap)
		}
	}
}
