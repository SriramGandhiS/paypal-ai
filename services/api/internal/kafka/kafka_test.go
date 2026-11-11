package kafka_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/paybridge/paybridge-ai/services/api/internal/kafka"
)

func TestKafkaBus_StrictInOrderPartitioning(t *testing.T) {
	bus := kafka.NewBus(3)
	ctx := context.Background()

	tenantID := uuid.New()
	caseID := uuid.New()

	// Publish 10 sequential events for the same case
	for i := 1; i <= 10; i++ {
		payload := []byte(fmt.Sprintf("step_%d", i))
		_, err := bus.Publish(ctx, kafka.DefaultEventTopic, tenantID, caseID, "STATE_UPDATED", payload)
		if err != nil {
			t.Fatalf("publish failed at step %d: %v", i, err)
		}
	}

	stream := bus.GetPartitionStream(tenantID, caseID)
	if len(stream) != 10 {
		t.Fatalf("expected 10 ordered events, got %d", len(stream))
	}

	// Verify exact sequence preservation
	for i := 0; i < 10; i++ {
		expectedPayload := fmt.Sprintf("step_%d", i+1)
		if string(stream[i].Payload) != expectedPayload {
			t.Errorf("out-of-order delivery at index %d: got %s, want %s", i, string(stream[i].Payload), expectedPayload)
		}
	}
}

func TestKafkaBus_DeadLetterRouting(t *testing.T) {
	bus := kafka.NewBus(3)
	tenantID := uuid.New()
	caseID := uuid.New()

	failedEnv := kafka.EventEnvelope{
		Topic:        kafka.DefaultEventTopic,
		TenantID:     tenantID,
		CaseID:       caseID,
		EventType:    "WEBHOOK_PARSE_FAILED",
		Payload:      []byte("corrupted_payload"),
		Attempts:     4,
	}

	bus.RouteToDeadLetter(failedEnv)

	deadLetters := bus.GetDeadLetters()
	if len(deadLetters) != 1 {
		t.Fatalf("expected 1 dead letter, got %d", len(deadLetters))
	}

	if deadLetters[0].Topic != kafka.DefaultDLQTopic {
		t.Errorf("expected deadletter topic, got %s", deadLetters[0].Topic)
	}
}

func TestKafkaBus_32PartitionDistribution(t *testing.T) {
	bus := kafka.NewBusWithConfig(kafka.DefaultConfig())
	ctx := context.Background()
	tenantID := uuid.New()

	seenPartitions := make(map[int32]int)

	// Publish across 128 different cases to verify 32-partition distribution
	for i := 0; i < 128; i++ {
		caseID := uuid.New()
		env, err := bus.Publish(ctx, kafka.DefaultEventTopic, tenantID, caseID, "CASE_CREATED", []byte("{}"))
		if err != nil {
			t.Fatalf("publish failed: %v", err)
		}
		if env.Partition < 0 || env.Partition >= 32 {
			t.Fatalf("invalid partition index %d (must be between 0 and 31)", env.Partition)
		}
		seenPartitions[env.Partition]++
	}

	// With 128 items across 32 partitions, we expect at least 16 partitions to be populated
	if len(seenPartitions) < 16 {
		t.Errorf("poor partition distribution: only %d / 32 partitions utilized", len(seenPartitions))
	}
}

func TestKafkaBus_ConsumerLagTracking(t *testing.T) {
	bus := kafka.NewBusWithConfig(kafka.DefaultConfig())
	ctx := context.Background()
	tenantID := uuid.New()
	caseID := uuid.New()

	// Publish 5 events
	for i := 0; i < 5; i++ {
		_, err := bus.Publish(ctx, kafka.DefaultEventTopic, tenantID, caseID, "EVENT", []byte("payload"))
		if err != nil {
			t.Fatalf("publish failed: %v", err)
		}
	}

	groupID := "reconciliation-worker-group"
	lag := bus.GetConsumerLag(groupID)
	if lag != 5 {
		t.Fatalf("expected consumer lag of 5, got %d", lag)
	}

	// Commit offset for partition
	p, _ := kafka.PartitionFor(tenantID, caseID, 32)
	bus.CommitOffset(groupID, p, 5)

	lagAfterCommit := bus.GetConsumerLag(groupID)
	if lagAfterCommit != 0 {
		t.Fatalf("expected 0 lag after committing offset, got %d", lagAfterCommit)
	}
}
