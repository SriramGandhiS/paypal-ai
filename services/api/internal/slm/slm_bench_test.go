package slm_test

import (
	"context"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/paybridge/paybridge-ai/services/api/internal/domain"
	"github.com/paybridge/paybridge-ai/services/api/internal/slm"
)

// TestSLM_500ConcurrentTriageBenchmark verifies sub-second P99 latency across 500 concurrent requests
func TestSLM_500ConcurrentTriageBenchmark(t *testing.T) {
	client := slm.NewSLMClient()
	ctx := context.Background()

	const concurrency = 500
	latencies := make([]time.Duration, concurrency)

	var wg sync.WaitGroup
	var mu sync.Mutex

	start := time.Now()

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			reqStart := time.Now()

			testCase := &domain.PaymentCase{
				ID:               uuid.New(),
				TenantID:         uuid.New(),
				Amount:           decimal.NewFromInt(100),
				Currency:         "USD",
				ProviderState:    domain.ProviderStateCaptureCompleted,
				ApplicationState: domain.AppStateAPITimeout,
			}

			_, err := client.FastTriage(ctx, testCase, "Standard checkout triage memo")
			if err != nil {
				t.Errorf("triage request %d failed: %v", idx, err)
				return
			}

			reqDur := time.Since(reqStart)

			mu.Lock()
			latencies[idx] = reqDur
			mu.Unlock()
		}(i)
	}

	wg.Wait()
	totalElapsed := time.Since(start)

	// Sort latencies to compute P50, P95, and P99
	sort.Slice(latencies, func(i, j int) bool {
		return latencies[i] < latencies[j]
	})

	p50 := latencies[concurrency*50/100]
	p95 := latencies[concurrency*95/100]
	p99 := latencies[concurrency*99/100]

	t.Logf("500 Concurrent Triage Benchmark Completed in %v", totalElapsed)
	t.Logf("Latencies: P50=%v, P95=%v, P99=%v", p50, p95, p99)

	// Invariant Check: P99 latency must be well below the 450ms target
	if p99 > 450*time.Millisecond {
		t.Errorf("target P99 latency 450ms exceeded: got %v", p99)
	}
}

func BenchmarkSLMClient_FastTriage(b *testing.B) {
	client := slm.NewSLMClient()
	ctx := context.Background()

	testCase := &domain.PaymentCase{
		ID:               uuid.New(),
		TenantID:         uuid.New(),
		Amount:           decimal.NewFromInt(250),
		Currency:         "USD",
		ProviderState:    domain.ProviderStateCaptureCompleted,
		ApplicationState: domain.AppStateAPITimeout,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.FastTriage(ctx, testCase, "Benchmark memo")
		if err != nil {
			b.Fatalf("fast triage failed: %v", err)
		}
	}
}

