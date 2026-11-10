package safeaction

import (
	"testing"

	"github.com/paybridge/paybridge-ai/services/api/internal/reducer"
)

func TestSafeAction_BlocksRetryOnUnknownCapture(t *testing.T) {
	testCases := []struct {
		name    string
		state   reducer.StateVector
		attempt AttemptContext
	}{
		{
			name: "API Timeout with Approved provider state",
			state: reducer.StateVector{
				ProviderState:    reducer.ProviderStateApproved,
				ApplicationState: reducer.AppStateAPITimeout,
			},
			attempt: AttemptContext{Status: "CAPTURE_REQUESTED", AttemptNo: 1},
		},
		{
			name: "Unknown provider state",
			state: reducer.StateVector{
				ProviderState:    reducer.ProviderStateUnknown,
				ApplicationState: reducer.AppStateAPITimeout,
			},
			attempt: AttemptContext{Status: "CAPTURE_UNKNOWN", AttemptNo: 1},
		},
		{
			name: "Attempt status CAPTURE_UNKNOWN",
			state: reducer.StateVector{
				ProviderState:    reducer.ProviderStateUnknown,
				ApplicationState: reducer.AppStateReconciling,
			},
			attempt: AttemptContext{Status: "CAPTURE_UNKNOWN", AttemptNo: 1},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := Evaluate(tc.state, tc.attempt)

			// 1. Primary action MUST be DO_NOT_REPEAT_CHECKOUT
			if res.PrimaryAction.Code != ActionDoNotRepeatCheckout {
				t.Errorf("expected primary action DO_NOT_REPEAT_CHECKOUT, got %s", res.PrimaryAction.Code)
			}

			// 2. RETRY_CHECKOUT MUST be in blocked list
			retryBlocked := false
			for _, b := range res.Blocked {
				if b.Code == ActionRetryCheckout {
					retryBlocked = true
					break
				}
			}
			if !retryBlocked {
				t.Errorf("expected RETRY_CHECKOUT to be in Blocked actions list")
			}

			// 3. RETRY_CHECKOUT MUST NEVER appear in Allowed list
			for _, a := range res.Allowed {
				if a.Code == ActionRetryCheckout {
					t.Fatalf("VIOLATION: RETRY_CHECKOUT appeared in Allowed actions list during ambiguous state!")
				}
			}
		})
	}
}

func TestSafeAction_CompletedCapture(t *testing.T) {
	state := reducer.StateVector{
		ProviderState:    reducer.ProviderStateCaptureCompleted,
		ApplicationState: reducer.AppStateSynchronised,
	}
	attempt := AttemptContext{Status: "CAPTURE_COMPLETED", AttemptNo: 1, Terminal: false}

	res := Evaluate(state, attempt)
	if res.PrimaryAction.Code != ActionNoneRequired {
		t.Errorf("expected NONE_REQUIRED, got %s", res.PrimaryAction.Code)
	}
}

func TestSafeAction_OrderUnknownBlocksRetry(t *testing.T) {
	state := reducer.StateVector{
		ProviderState:    reducer.ProviderStateUnknown,
		ApplicationState: reducer.AppStateAPITimeout,
	}
	attempt := AttemptContext{Status: "ORDER_UNKNOWN", AttemptNo: 1}

	res := Evaluate(state, attempt)
	if res.PrimaryAction.Code != ActionDoNotRepeatCheckout {
		t.Errorf("expected DO_NOT_REPEAT_CHECKOUT on ORDER_UNKNOWN, got %s", res.PrimaryAction.Code)
	}

	for _, a := range res.Allowed {
		if a.Code == ActionRetryCheckout {
			t.Fatalf("VIOLATION: retry allowed on ORDER_UNKNOWN")
		}
	}
}

func TestSafeAction_DeniedCaptureRequestsReview(t *testing.T) {
	state := reducer.StateVector{
		ProviderState:    reducer.ProviderStateCaptureDenied,
		ApplicationState: reducer.AppStateSynchronised,
	}
	attempt := AttemptContext{Status: "CAPTURE_DENIED", AttemptNo: 1, Terminal: true}

	res := Evaluate(state, attempt)
	if res.PrimaryAction.Code != ActionRequestHumanReview {
		t.Errorf("expected REQUEST_HUMAN_REVIEW, got %s", res.PrimaryAction.Code)
	}
}

