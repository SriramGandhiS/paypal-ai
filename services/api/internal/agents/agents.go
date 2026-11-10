package agents

import (
	"fmt"
	"strings"

	"github.com/paybridge/paybridge-ai/services/api/internal/claims"
	"github.com/paybridge/paybridge-ai/services/api/internal/safeaction"
)

type AgentID string

const (
	AgentA01PaymentReadiness  AgentID = "A01"
	AgentA02TrafficService    AgentID = "A02"
	AgentA03SenderAnalysis    AgentID = "A03"
	AgentA04ReceiverAnalysis  AgentID = "A04"
	AgentA05JourneyTracker    AgentID = "A05"
	AgentA06RootCauseAnalyst  AgentID = "A06"
	AgentA07CurrencyIntel     AgentID = "A07"
	AgentA08RecoveryPlanner   AgentID = "A08"
	AgentA09DualCommunicator  AgentID = "A09"
	AgentA10ResolutionTracker AgentID = "A10"
)

// AgentOutput models standard structured output from specialist agents
type AgentOutput struct {
	AgentID     AgentID             `json:"agent_id"`
	AgentName   string              `json:"agent_name"`
	Summary     string              `json:"summary"`
	Claims      []claims.Claim      `json:"claims"`
	NextActions []safeaction.Action `json:"next_actions,omitempty"`
	Metadata    map[string]any      `json:"metadata,omitempty"`
}

// WrapUntrustedText applies prompt injection defense delimiters
func WrapUntrustedText(label, text string) string {
	sanitized := strings.ReplaceAll(text, "</untrusted_data>", "")
	return fmt.Sprintf("<untrusted_data label=%q>\n%s\n</untrusted_data>", label, sanitized)
}

// GenerateDeterministicOutput produces reliable rule-based output when LLM_ENABLED=false (Demo E)
// or as the ground-truth baseline for all 10 specialist agents (A01–A10).
func GenerateDeterministicOutput(agentID AgentID, contextData map[string]any) AgentOutput {
	orderEvID := "ev-order-1"
	if v, ok := contextData["order_evidence_id"]; ok && fmt.Sprintf("%v", v) != "" {
		orderEvID = fmt.Sprintf("%v", v)
	}

	timeoutEvID := "ev-timeout-1"
	if v, ok := contextData["timeout_evidence_id"]; ok && fmt.Sprintf("%v", v) != "" {
		timeoutEvID = fmt.Sprintf("%v", v)
	}

	routeEvID := "ev-route-rule-1"
	if v, ok := contextData["route_evidence_id"]; ok && fmt.Sprintf("%v", v) != "" {
		routeEvID = fmt.Sprintf("%v", v)
	}

	switch agentID {
	case AgentA01PaymentReadiness:
		return AgentOutput{
			AgentID:   AgentA01PaymentReadiness,
			AgentName: "Payment Readiness",
			Summary:   "Verified cross-border corridor US -> IN eligibility and statutory purpose code P0802.",
			Claims: []claims.Claim{
				{
					Code:        "ROUTE_READY_VERIFIED",
					Kind:        claims.ClaimKindConfirmed,
					Statement:   "Corridor US to IN is supported under RBI purpose code P0802 (Software Consulting Services).",
					EvidenceIDs: []string{routeEvID},
					ViewerScope: "ALL",
				},
			},
			Metadata: map[string]any{
				"corridor":     "US-IN",
				"purpose_code": "P0802",
				"status":       "ELIGIBLE",
			},
		}

	case AgentA02TrafficService:
		return AgentOutput{
			AgentID:   AgentA02TrafficService,
			AgentName: "Traffic & Service Context",
			Summary:   "Telemetry shows upstream PayPal API responding normally; isolated timeout occurred at client transport layer.",
			Claims: []claims.Claim{
				{
					Code:        "TELEMETRY_ISOLATED_TIMEOUT",
					Kind:        claims.ClaimKindInference,
					Statement:   "Observed HTTP 504 gateway timeout is localized to this individual capture call without systemic provider degradation.",
					EvidenceIDs: []string{timeoutEvID},
					InferenceRuleID: "RULE_T02_ISOLATED_INCIDENT",
					ViewerScope: "OPS",
				},
			},
			Metadata: map[string]any{
				"p99_latency_ms": 1420,
				"provider_health": "NOMINAL",
			},
		}

	case AgentA03SenderAnalysis:
		return AgentOutput{
			AgentID:   AgentA03SenderAnalysis,
			AgentName: "Sender Analysis",
			Summary:   "Sender checkout attempt initiated with verified authorization.",
			Claims: []claims.Claim{
				{
					Code:        "SENDER_CAPTURE_IN_FLIGHT",
					Kind:        claims.ClaimKindConfirmed,
					Statement:   "The sender submitted capture request with stable idempotency key.",
					EvidenceIDs: []string{orderEvID},
					ViewerScope: "SENDER",
				},
			},
		}

	case AgentA04ReceiverAnalysis:
		hasConsent := false
		if v, ok := contextData["receiver_consent"]; ok {
			if b, ok := v.(bool); ok && b {
				hasConsent = true
			}
		}

		if !hasConsent {
			return AgentOutput{
				AgentID:   AgentA04ReceiverAnalysis,
				AgentName: "Receiver Analysis",
				Summary:   "Receiver analysis suspended pending active consent grant.",
				Metadata: map[string]any{
					"consent_required": true,
				},
			}
		}

		return AgentOutput{
			AgentID:   AgentA04ReceiverAnalysis,
			AgentName: "Receiver Analysis",
			Summary:   "Receiver business entity and domestic clearing bank readiness verified.",
			Claims: []claims.Claim{
				{
					Code:        "RECEIVER_ACCOUNT_ACTIVE",
					Kind:        claims.ClaimKindConfirmed,
					Statement:   "Receiver account is verified and eligible to receive foreign inward remittance.",
					EvidenceIDs: []string{orderEvID},
					ViewerScope: "RECEIVER",
				},
			},
		}

	case AgentA05JourneyTracker:
		return AgentOutput{
			AgentID:   AgentA05JourneyTracker,
			AgentName: "Journey Tracker",
			Summary:   "Chronological lifecycle timeline reconstructed from local app traces and provider status updates.",
			Claims: []claims.Claim{
				{
					Code:        "JOURNEY_TIMEOUT_DETECTED",
					Kind:        claims.ClaimKindConfirmed,
					Statement:   "Network gateway timeout detected during PayPal capture call. The response was not received before deadline.",
					EvidenceIDs: []string{timeoutEvID},
					ViewerScope: "ALL",
				},
			},
		}

	case AgentA06RootCauseAnalyst:
		return AgentOutput{
			AgentID:   AgentA06RootCauseAnalyst,
			AgentName: "Root Cause Analyst",
			Summary:   "Primary cause: Network communication drop between PayBridge client and PayPal gateway during capture execution (F02 Transient Network Drop).",
			Claims: []claims.Claim{
				{
					Code:            "ROOT_CAUSE_TRANSIENT_DROP",
					Kind:            claims.ClaimKindInference,
					Statement:       "Client experienced HTTP 504 Gateway Timeout while PayPal order status remains pending confirmation.",
					EvidenceIDs:     []string{timeoutEvID},
					InferenceRuleID: "RULE_F02_TIMEOUT_DISCONNECT",
					ViewerScope:     "ALL",
				},
			},
		}

	case AgentA07CurrencyIntel:
		return AgentOutput{
			AgentID:   AgentA07CurrencyIntel,
			AgentName: "Currency & Fee Intelligence",
			Summary:   "FX exchange rate quote locked at 84.35 INR/USD; guaranteed transparent fee schedule.",
			Claims: []claims.Claim{
				{
					Code:        "FX_RATE_LOCKED",
					Kind:        claims.ClaimKindConfirmed,
					Statement:   "Exchange rate of 84.35 INR per USD verified with zero hidden spreads.",
					EvidenceIDs: []string{routeEvID},
					ViewerScope: "ALL",
				},
			},
			Metadata: map[string]any{
				"quote_rate":  "84.35",
				"fee_amount":  "2.50",
				"currency_in": "USD",
				"currency_out": "INR",
			},
		}

	case AgentA08RecoveryPlanner:
		return AgentOutput{
			AgentID:   AgentA08RecoveryPlanner,
			AgentName: "Recovery Planner",
			Summary:   "Strictly blocking blind retry checkout. Safe recovery initiated via read-only GET query on original order.",
			NextActions: []safeaction.Action{
				{
					Code:             safeaction.ActionDoNotRepeatCheckout,
					Label:            "Do Not Repeat Checkout",
					Rationale:        "Prevent duplicate payment capture.",
					PreconditionsMet: true,
					DesignatedActor:  "SENDER",
				},
				{
					Code:             safeaction.ActionAwaitProviderConfirmation,
					Label:            "Reconcile Original Order (Read-Only)",
					Rationale:        "Reconciler polling original capture reference.",
					PreconditionsMet: true,
					DesignatedActor:  "SYSTEM",
				},
			},
		}

	case AgentA09DualCommunicator:
		return AgentOutput{
			AgentID:   AgentA09DualCommunicator,
			AgentName: "Dual-Side Communicator",
			Summary:   "Role-specific calm communications formulated without misleading 'failed' terminology.",
			Claims: []claims.Claim{
				{
					Code:        "COMM_SENDER_NOTICE",
					Kind:        claims.ClaimKindConfirmed,
					Statement:   "Your payment is currently being verified with PayPal. Please do not submit another payment while verification is underway.",
					EvidenceIDs: []string{orderEvID},
					ViewerScope: "SENDER",
				},
				{
					Code:        "COMM_RECEIVER_NOTICE",
					Kind:        claims.ClaimKindConfirmed,
					Statement:   "An incoming transfer is undergoing provider confirmation. Updates will reflect once PayPal confirms completion.",
					EvidenceIDs: []string{orderEvID},
					ViewerScope: "RECEIVER",
				},
			},
		}

	case AgentA10ResolutionTracker:
		return AgentOutput{
			AgentID:   AgentA10ResolutionTracker,
			AgentName: "Resolution Tracker",
			Summary:   "SLA tracking active. Case resolution deadline set to 120s from timeout event.",
			Claims: []claims.Claim{
				{
					Code:        "SLA_ACTIVE",
					Kind:        claims.ClaimKindConfirmed,
					Statement:   "Automated reconciliation SLA is on schedule; expected resolution within 45 seconds.",
					EvidenceIDs: []string{timeoutEvID},
					ViewerScope: "OPS",
				},
			},
			Metadata: map[string]any{
				"sla_deadline_seconds": 120,
				"reconciliation_active": true,
			},
		}

	default:
		return AgentOutput{
			AgentID:   agentID,
			AgentName: string(agentID),
			Summary:   "Deterministic analysis complete.",
		}
	}
}
