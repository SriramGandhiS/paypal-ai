package main

import (
	"context"
	"io"
	"log/slog"
	"testing"
)

func TestGoldenPath_DeterministicLifecycle(t *testing.T) {
	discardLogger := slog.New(slog.NewTextHandler(io.Discard, nil))
	ctx := context.Background()

	res, err := RunGoldenPath(ctx, discardLogger)
	if err != nil {
		t.Fatalf("RunGoldenPath failed: %v", err)
	}

	if res.CaseID == "" {
		t.Errorf("expected non-empty CaseID")
	}
	if res.OrderID == "" {
		t.Errorf("expected non-empty OrderID")
	}
	if res.CaptureID == "" {
		t.Errorf("expected non-empty CaptureID")
	}
	if res.ReconcileOutcome.ProviderState != "CAPTURE_COMPLETED" {
		t.Errorf("expected ProviderState CAPTURE_COMPLETED, got %s", res.ReconcileOutcome.ProviderState)
	}
	if res.ReconcileOutcome.ApplicationState != "SYNCHRONISED" {
		t.Errorf("expected ApplicationState SYNCHRONISED, got %s", res.ReconcileOutcome.ApplicationState)
	}
	if res.ReconcileOutcome.CaseStatus != "RESOLVED" {
		t.Errorf("expected CaseStatus RESOLVED, got %s", res.ReconcileOutcome.CaseStatus)
	}
}
