package reducer

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrSimulatedEventOnRealCase = errors.New("simulated event cannot modify real sandbox or production case")
	ErrMonotonicityViolation    = errors.New("provider state cannot regress backwards without explicit refund/reversal")
	ErrInvalidTransition       = errors.New("invalid state transition")
)

// Reducer is the ONLY component permitted to modify payment state dimensions.
type Reducer struct{}

// NewReducer constructs the deterministic Reducer
func New() *Reducer {
	return NewReducer()
}

func NewReducer() *Reducer {
	return &Reducer{}
}

// Apply deterministically reduces an event against the current state vector.
// Invariant 1: Timeout != Failure. Lost response sets application_state to API_TIMEOUT,
// provider_state to UNKNOWN (preserving last_confirmed_provider_state), NEVER CAPTURE_DENIED.
// Invariant 2: Simulation isolation. E6/SIMULATED cannot modify non-simulated environments.
// Invariant 3: Monotonic lattice. Provider states cannot regress backwards without refund events.
func (r *Reducer) Apply(ctx context.Context, current StateVector, event Event) (StateVector, []StateTransition, error) {
	next := current
	var transitions []StateTransition

	// Environment isolation check
	if (event.Tier == TierE6Simulation || event.SourceTrust == SourceTrustSimulated) &&
		(event.Environment == "REAL_SANDBOX" || event.Environment == "PRODUCTION") {
		return current, nil, ErrSimulatedEventOnRealCase
	}

	switch event.EventType {
	case "CHECKOUT_TIMEOUT", "API_TIMEOUT", "RESPONSE_LOST":
		// TIMEOUT != FAILURE invariant
		if current.ApplicationState != AppStateAPITimeout {
			transitions = append(transitions, StateTransition{
				Dimension: "application_state",
				From:      string(current.ApplicationState),
				To:        string(AppStateAPITimeout),
				Reason:    "Gateway timeout or network response dropped",
				EventID:   event.EventID,
			})
			next.ApplicationState = AppStateAPITimeout
		}

		if current.ProviderState != ProviderStateUnknown {
			transitions = append(transitions, StateTransition{
				Dimension: "provider_state",
				From:      string(current.ProviderState),
				To:        string(ProviderStateUnknown),
				Reason:    "Provider outcome is unknown pending read-only reconciliation",
				EventID:   event.EventID,
			})
			next.ProviderState = ProviderStateUnknown
		}

		if current.CaseStatus == CaseStatusOpen || current.CaseStatus == CaseStatusObserving {
			transitions = append(transitions, StateTransition{
				Dimension: "case_status",
				From:      string(current.CaseStatus),
				To:        string(CaseStatusInvestigating),
				Reason:    "Unknown timeout triggered investigation workflow",
				EventID:   event.EventID,
			})
			next.CaseStatus = CaseStatusInvestigating
		}

	case "PROVIDER_CAPTURE_COMPLETED", "PAYMENT.CAPTURE.COMPLETED":
		// Authoritative completion from E1 (GET capture) or E2 (verified webhook)
		if event.Tier != TierE1DirectProviderRead && event.Tier != TierE2VerifiedWebhook {
			return current, nil, fmt.Errorf("%w: completion requires tier E1 or E2 evidence", ErrInvalidTransition)
		}

		if current.ProviderState != ProviderStateCaptureCompleted {
			transitions = append(transitions, StateTransition{
				Dimension: "provider_state",
				From:      string(current.ProviderState),
				To:        string(ProviderStateCaptureCompleted),
				Reason:    "Provider confirmed capture completion",
				EventID:   event.EventID,
			})
			next.ProviderState = ProviderStateCaptureCompleted
			next.LastConfirmedProviderState = ProviderStateCaptureCompleted
			next.ProviderStateEvidenceID = event.EvidenceID
			next.ProviderStateObservedAt = event.ObservedAt
		}

		// If we were in API_TIMEOUT, RECONCILING, or WAITING_FOR_WEBHOOK, we are now SYNCHRONISED
		if current.ApplicationState != AppStateSynchronised {
			transitions = append(transitions, StateTransition{
				Dimension: "application_state",
				From:      string(current.ApplicationState),
				To:        string(AppStateSynchronised),
				Reason:    "Application synchronized with authoritative provider capture",
				EventID:   event.EventID,
			})
			next.ApplicationState = AppStateSynchronised
		}

	case "PROVIDER_CAPTURE_DENIED", "PAYMENT.CAPTURE.DENIED":
		if current.ProviderState == ProviderStateCaptureCompleted {
			return current, nil, fmt.Errorf("%w: cannot transition completed capture to denied", ErrMonotonicityViolation)
		}
		if event.Tier != TierE1DirectProviderRead && event.Tier != TierE2VerifiedWebhook {
			return current, nil, fmt.Errorf("%w: denial requires verified provider evidence", ErrInvalidTransition)
		}

		if current.ProviderState != ProviderStateCaptureDenied {
			transitions = append(transitions, StateTransition{
				Dimension: "provider_state",
				From:      string(current.ProviderState),
				To:        string(ProviderStateCaptureDenied),
				Reason:    "Provider explicitly declined payment instrument",
				EventID:   event.EventID,
			})
			next.ProviderState = ProviderStateCaptureDenied
			next.LastConfirmedProviderState = ProviderStateCaptureDenied
			next.ProviderStateEvidenceID = event.EvidenceID
			next.ProviderStateObservedAt = event.ObservedAt
		}

		if current.ApplicationState != AppStateSynchronised {
			next.ApplicationState = AppStateSynchronised
		}

	case "PROVIDER_ORDER_CREATED":
		if current.ProviderState == ProviderStateNotCreated {
			transitions = append(transitions, StateTransition{
				Dimension: "provider_state",
				From:      string(current.ProviderState),
				To:        string(ProviderStateCreated),
				Reason:    "Order created on provider",
				EventID:   event.EventID,
			})
			next.ProviderState = ProviderStateCreated
			next.LastConfirmedProviderState = ProviderStateCreated
			next.ProviderStateEvidenceID = event.EvidenceID
			next.ProviderStateObservedAt = event.ObservedAt
		}

	case "BUYER_APPROVED_ORDER":
		if current.ProviderState == ProviderStateCreated {
			transitions = append(transitions, StateTransition{
				Dimension: "provider_state",
				From:      string(current.ProviderState),
				To:        string(ProviderStateApproved),
				Reason:    "Buyer approved checkout order on PayPal",
				EventID:   event.EventID,
			})
			next.ProviderState = ProviderStateApproved
			next.LastConfirmedProviderState = ProviderStateApproved
			next.ProviderStateEvidenceID = event.EvidenceID
			next.ProviderStateObservedAt = event.ObservedAt
		}

	case "RECONCILIATION_STARTED":
		if current.ApplicationState == AppStateAPITimeout || current.ApplicationState == AppStateEventGap {
			transitions = append(transitions, StateTransition{
				Dimension: "application_state",
				From:      string(current.ApplicationState),
				To:        string(AppStateReconciling),
				Reason:    "Worker actively querying provider GET endpoints for authoritative truth",
				EventID:   event.EventID,
			})
			next.ApplicationState = AppStateReconciling
		}

	case "PAYMENT.CAPTURE.REFUNDED":
		if current.ProviderState == ProviderStateCaptureCompleted {
			transitions = append(transitions, StateTransition{
				Dimension: "provider_state",
				From:      string(current.ProviderState),
				To:        string(ProviderStateRefunded),
				Reason:    "Provider confirmed refund event",
				EventID:   event.EventID,
			})
			next.ProviderState = ProviderStateRefunded
			next.LastConfirmedProviderState = ProviderStateRefunded
			next.ProviderStateEvidenceID = event.EvidenceID
			next.ProviderStateObservedAt = event.ObservedAt
		}

	case "BANK_CREDIT_CONFIRMED":
		// Receiver state transition (requires receiver connector evidence)
		if current.ReceiverState != ReceiverStateVerifiedBankCredited {
			transitions = append(transitions, StateTransition{
				Dimension: "receiver_state",
				From:      string(current.ReceiverState),
				To:        string(ReceiverStateVerifiedBankCredited),
				Reason:    "Cryptographic bank deposit receipt verified",
				EventID:   event.EventID,
			})
			next.ReceiverState = ReceiverStateVerifiedBankCredited
		}
	}

	if len(transitions) > 0 {
		next.StateVersion++
	}

	return next, transitions, nil
}
