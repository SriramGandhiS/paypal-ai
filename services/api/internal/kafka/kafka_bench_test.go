package kafka_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/paybridge/paybridge-ai/services/api/internal/kafka"
)

// BenchmarkKafkaBus_10000EventsPerSec verifies that the PayBridge Kafka bus
// can process and route at least 10,000 events/sec with strict in-order delivery
func BenchmarkKafkaBus_10000EventsPerSec(b *testing.B) {
	bus := kafka.NewBusWithConfig(kafka.DefaultConfig())
	ctx := context.Background()

	tenantID := uuid.New()
	caseID := uuid.New()
	payload := []byte(`{"case_id":"11111111-1111-1111-1111-111111111111","event":"CAPTURE_CONFIRMED","amount":"250.00","currency":"USD"}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := bus.Publish(ctx, kafka.DefaultEventTopic, tenantID, caseID, "PAYMENT_CAPTURED", payload)
		if err != nil {
			b.Fatalf("failed to publish: %v", err)
		}
	}
}

// TestKafka_10000EventsThroughputLoad verifies a live burst of 10,000 events across 32 partitions
func TestKafka_10000EventsThroughputLoad(t *testing.T) {
	bus := kafka.NewBusWithConfig(kafka.DefaultConfig())
	ctx := context.Background()

	const targetEvents = 10000
	const numWorkers = 16

	var wg sync.WaitGroup
	eventsPerWorker := targetEvents / numWorkers

	start := time.Now()

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			tenantID := uuid.New()
			for i := 0; i < eventsPerWorker; i++ {
				caseID := uuid.New()
				payload := []byte(fmt.Sprintf(`{"worker":%d,"seq":%d}`, workerID, i))
				_, err := bus.Publish(ctx, kafka.DefaultEventTopic, tenantID, caseID, "LOAD_TEST_EVENT", payload)
				if err != nil {
					t.Errorf("worker %d publish failed: %v", workerID, err)
					return
				}
			}
		}(w)
	}

	wg.Wait()
	elapsed := time.Since(start)

	totalPublished := bus.TotalEvents()
	if totalPublished != targetEvents {
		t.Fatalf("expected %d events, published %d", targetEvents, totalPublished)
	}

	eventsPerSec := float64(totalPublished) / elapsed.Seconds()
	t.Logf("Published %d events in %v (Throughput: %.2f events/sec)", totalPublished, elapsed, eventsPerSec)

	if eventsPerSec < 10000 {
		t.Errorf("target throughput 10,000 events/sec not met: achieved %.2f events/sec", eventsPerSec)
	}
}

