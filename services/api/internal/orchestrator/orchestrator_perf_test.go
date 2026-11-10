package orchestrator_test

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/paybridge/paybridge-ai/services/api/internal/claims"
	"github.com/paybridge/paybridge-ai/services/api/internal/orchestrator"
	"github.com/paybridge/paybridge-ai/services/api/internal/reducer"
	"github.com/paybridge/paybridge-ai/services/api/internal/safeaction"
)

func runConcurrencyBenchmark(b *testing.B, concurrency int) {
	val := claims.NewValidator()
	orc := orchestrator.NewOrchestrator(orchestrator.Config{TimeoutSeconds: 5}, val)

	var successCount int64
	var failureCount int64

	startMem := runtime.MemStats{}
	runtime.ReadMemStats(&startMem)

	start := time.Now()

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, concurrency)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			snap := buildTestSnapshot(fmt.Sprintf("bench-case-%d", id))
			state := reducer.StateVector{ApplicationState: reducer.AppStateAPITimeout}
			attempt := safeaction.AttemptContext{Status: "CAPTURE_UNKNOWN"}

			res, err := orc.RunDAG(context.Background(), snap, state, attempt)
			if err != nil || res.Status != "COMPLETED" {
				atomic.AddInt64(&failureCount, 1)
			} else {
				atomic.AddInt64(&successCount, 1)
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(start)

	endMem := runtime.MemStats{}
	runtime.ReadMemStats(&endMem)
	allocMB := float64(endMem.TotalAlloc-startMem.TotalAlloc) / 1024 / 1024

	fmt.Printf("[PERF] Concurrency=%2d | Total Time: %v | Avg Per Run: %v | Success: %d | Fail: %d | Alloc: %.2f MB\n",
		concurrency, duration, duration/time.Duration(concurrency), successCount, failureCount, allocMB)
}

func TestPerformance_StressMatrix(t *testing.T) {
	concurrencyLevels := []int{1, 5, 10, 25}
	for _, c := range concurrencyLevels {
		runConcurrencyBenchmark(&testing.B{}, c)
	}
}

// TestLargePayloadStress tests orchestrator performance with 1,000 evidence items
func TestLargePayloadStress(t *testing.T) {
	val := claims.NewValidator()
	orc := orchestrator.NewOrchestrator(orchestrator.Config{TimeoutSeconds: 10}, val)

	snap := &claims.EvidenceSnapshot{
		SnapshotVersion: 1,
		CaseID:          "case-large-payload",
		EvidenceItems:   make(map[string]claims.SnapshotEvidenceItem, 1000),
	}

	for i := 0; i < 1000; i++ {
		eid := fmt.Sprintf("ev-large-%d", i)
		snap.EvidenceItems[eid] = claims.SnapshotEvidenceItem{
			ID:          eid,
			Tier:        reducer.TierE1DirectProviderRead,
			SourceTrust: reducer.SourceTrustProviderVerified,
			Kind:        "PAYLOAD_SAMPLE",
			Verified:    true,
		}
	}

	state := reducer.StateVector{ApplicationState: reducer.AppStateAPITimeout}
	attempt := safeaction.AttemptContext{Status: "CAPTURE_UNKNOWN"}

	start := time.Now()
	res, err := orc.RunDAG(context.Background(), snap, state, attempt)
	dur := time.Since(start)

	if err != nil {
		t.Fatalf("unexpected error with 1,000 evidence items: %v", err)
	}
	if res.Status != "COMPLETED" {
		t.Errorf("expected COMPLETED, got %s", res.Status)
	}
	fmt.Printf("[PERF] 1,000 Evidence Items Processed in: %v\n", dur)
}
