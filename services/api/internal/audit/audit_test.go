package audit_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/paybridge/paybridge-ai/services/api/internal/audit"
)

func TestAuditChain_AppendAndVerifyIntegrity(t *testing.T) {
	chain := audit.NewChain()
	ctx := context.Background()

	tenantID := uuid.New()
	caseID := uuid.New()
	actorID := uuid.New()

	// Append 10 audit records
	for i := 1; i <= 10; i++ {
		_, err := chain.Append(
			ctx,
			tenantID, caseID, actorID,
			"platform_operator",
			"CAPTURE_AUTHORIZED",
			"SANDBOX",
			"User manual approval",
			int64(i-1), int64(i),
			"REQ-101",
		)
		if err != nil {
			t.Fatalf("failed to append audit record %d: %v", i, err)
		}
	}

	valid, count, err := chain.VerifyChain(ctx)
	if err != nil || !valid {
		t.Fatalf("chain verification failed: %v", err)
	}

	if count != 10 {
		t.Errorf("expected 10 verified events, got %d", count)
	}
}

