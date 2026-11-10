package orchestrator

import (
	"context"
	"testing"

	"github.com/paybridge/paybridge-ai/services/api/internal/claims"
	"github.com/paybridge/paybridge-ai/services/api/internal/reducer"
	"github.com/paybridge/paybridge-ai/services/api/internal/safeaction"
)

func TestOrchestrator_RunDAGTimeoutScenario(t *testing.T) {
	val := claims.NewValidator()
	orc := NewOrchestrator(Config{
		LLMEnabled:     false, // Kill switch active (Demo E deterministic path)
		TimeoutSeconds: 5,
	}, val)

	snapshot := &claims.EvidenceSnapshot{
		SnapshotVersion: 3,
		CaseID:          "case-test-timeout",
		EvidenceItems: map[string]claims.SnapshotEvidenceItem{
			"ev-timeout-1": {
				ID:          "ev-timeout-1",
				Tier:        reducer.TierE1DirectProviderRead,
				SourceTrust: reducer.SourceTrustProviderVerified,
				Verified:    true,
			},
			"ev-order-1": {
				ID:          "ev-order-1",
				Tier:        reducer.TierE1DirectProviderRead,
				SourceTrust: reducer.SourceTrustProviderVerified,
				Verified:    true,
			},
		},
	}

	state := reducer.StateVector{
		ProviderState:    reducer.ProviderStateUnknown,
		ApplicationState: reducer.AppStateAPITimeout,
		CaseStatus:       reducer.CaseStatusInvestigating,
	}

	attempt := safeaction.AttemptContext{
		Status:    "CAPTURE_UNKNOWN",
		AttemptNo: 1,
	}

	ctx := context.Background()
	res, err := orc.RunDAG(ctx, snapshot, state, attempt)
	if err != nil {
		t.Fatalf("RunDAG failed: %v", err)
	}

	if res.Status != "COMPLETED" {
		t.Errorf("expected status COMPLETED, got %s", res.Status)
	}

	if len(res.AgentOutputs) < 5 {
		t.Errorf("expected at least 5 agent outputs in DAG, got %d", len(res.AgentOutputs))
	}

	// Invariant: DO_NOT_REPEAT_CHECKOUT must be present
	if res.SafeActions.PrimaryAction.Code != safeaction.ActionDoNotRepeatCheckout {
		t.Errorf("expected primary safe action DO_NOT_REPEAT_CHECKOUT, got %s", res.SafeActions.PrimaryAction.Code)
	}

	// Invariant: Validated claims count > 0
	if len(res.ValidatedClaims) == 0 {
		t.Errorf("expected validated claims > 0")
	}
}
