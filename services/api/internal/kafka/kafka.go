package kafka

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/fnv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

var (
	ErrDeadLetterLimitExceeded = errors.New("kafka: message exceeded max retry attempts; routed to dead-letter queue")
	ErrBrokerUnavailable      = errors.New("kafka: broker connection unavailable or leader election in progress")
	ErrConsumerGroupClosed     = errors.New("kafka: consumer group context closed")
)

const (
	DefaultPartitionCount = 32
	DefaultEventTopic     = "paybridge.payment.events.v1"
	DefaultDLQTopic       = "paybridge.deadletter.v1"
)

// EventEnvelope is the canonical wire format for all PayBridge Kafka messages
type EventEnvelope struct {
	Topic        string            `json:"topic"`
	Partition    int32             `json:"partition"`     // 0..31
	PartitionKey string            `json:"partition_key"` // tenant_id:case_id
	TenantID     uuid.UUID         `json:"tenant_id"`
	CaseID       uuid.UUID         `json:"case_id"`
	EventType    string            `json:"event_type"`
	Payload      []byte            `json:"payload"`
	Headers      map[string]string `json:"headers"`
	Timestamp    time.Time         `json:"timestamp"`
	Attempts     int               `json:"attempts"`
	Offset       int64             `json:"offset"`
}

// Config specifies Kafka cluster connection and durability settings
type Config struct {
	Brokers            []string      `json:"brokers"`
	ClientID           string        `json:"client_id"`
	Topic              string        `json:"topic"`
	DLQTopic           string        `json:"dlq_topic"`
	Partitions         int32         `json:"partitions"`
	Acks               string        `json:"acks"` // "all"
	Compression        string        `json:"compression"` // "snappy"
	IdempotentProducer bool          `json:"idempotent_producer"`
	MaxRetries         int           `json:"max_retries"`
	RetryBackoff       time.Duration `json:"retry_backoff"`
}

func DefaultConfig() Config {
	return Config{
		Brokers:            []string{"localhost:9092", "localhost:9093", "localhost:9094"},
		ClientID:           "paybridge-core-producer",
		Topic:              DefaultEventTopic,
		DLQTopic:           DefaultDLQTopic,
		Partitions:         DefaultPartitionCount,
		Acks:               "all",
		Compression:        "snappy",
		IdempotentProducer: true,
		MaxRetries:         5,
		RetryBackoff:       50 * time.Millisecond,
	}
}

// Producer defines the publishing interface
type Producer interface {
	Publish(ctx context.Context, topic string, tenantID, caseID uuid.UUID, eventType string, payload []byte) (*EventEnvelope, error)
	Close() error
}

// Consumer defines the subscription interface
type Consumer interface {
	Subscribe(ctx context.Context, groupID string, handler func(env *EventEnvelope) error) error
	GetLag(topic string, groupID string) int64
}

// PartitionRouter calculates deterministic physical partition mapping (0..31)
func PartitionFor(tenantID, caseID uuid.UUID, numPartitions int32) (int32, string) {
	if numPartitions <= 0 {
		numPartitions = DefaultPartitionCount
	}
	partitionKey := fmt.Sprintf("%s:%s", tenantID.String(), caseID.String())
	h := fnv.New32a()
	h.Write([]byte(partitionKey))
	partition := int32(h.Sum32() % uint32(numPartitions))
	return partition, partitionKey
}

// Bus provides high-throughput in-memory enterprise Kafka simulation
// satisfying all 32-partition ordering, DLQ routing, consumer groups, and lag metrics.
type Bus struct {
	mu           sync.RWMutex
	cfg          Config
	partitions   map[int32][]EventEnvelope // partition_id -> ordered slice
	deadLetters  []EventEnvelope
	offsets      map[int32]int64           // current high watermark per partition
	consumerLags map[string]map[int32]int64 // groupID -> partition -> committed offset
	totalEvents  uint64
	closed       bool
}

func NewBus(maxRetries int) *Bus {
	cfg := DefaultConfig()
	if maxRetries > 0 {
		cfg.MaxRetries = maxRetries
	}
	return NewBusWithConfig(cfg)
}

func NewBusWithConfig(cfg Config) *Bus {
	if cfg.Partitions <= 0 {
		cfg.Partitions = DefaultPartitionCount
	}
	b := &Bus{
		cfg:          cfg,
		partitions:   make(map[int32][]EventEnvelope),
		deadLetters:  make([]EventEnvelope, 0),
		offsets:      make(map[int32]int64),
		consumerLags: make(map[string]map[int32]int64),
	}
	for i := int32(0); i < cfg.Partitions; i++ {
		b.partitions[i] = make([]EventEnvelope, 0, 1024)
		b.offsets[i] = 0
	}
	return b
}

// Publish publishes an event to a topic, strictly enforcing partition key `tenant_id:case_id`
// and routing to the designated partition [0..31]
func (b *Bus) Publish(ctx context.Context, topic string, tenantID, caseID uuid.UUID, eventType string, payload []byte) (*EventEnvelope, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return nil, ErrBrokerUnavailable
	}

	partition, partitionKey := PartitionFor(tenantID, caseID, b.cfg.Partitions)
	offset := b.offsets[partition]
	b.offsets[partition] = offset + 1

	now := time.Now().UTC()
	env := EventEnvelope{
		Topic:        topic,
		Partition:    partition,
		PartitionKey: partitionKey,
		TenantID:     tenantID,
		CaseID:       caseID,
		EventType:    eventType,
		Payload:      payload,
		Headers: map[string]string{
			"producer_id": b.cfg.ClientID,
			"compression": b.cfg.Compression,
			"acks":        b.cfg.Acks,
		},
		Timestamp: now,
		Attempts:  1,
		Offset:    offset,
	}

	b.partitions[partition] = append(b.partitions[partition], env)
	atomic.AddUint64(&b.totalEvents, 1)

	return &env, nil
}

// GetPartitionStream returns the ordered event log for a given case
func (b *Bus) GetPartitionStream(tenantID, caseID uuid.UUID) []EventEnvelope {
	b.mu.RLock()
	defer b.mu.RUnlock()

	partition, _ := PartitionFor(tenantID, caseID, b.cfg.Partitions)
	partitionEvents := b.partitions[partition]

	res := make([]EventEnvelope, 0, len(partitionEvents))
	for _, env := range partitionEvents {
		if env.TenantID == tenantID && env.CaseID == caseID {
			res = append(res, env)
		}
	}
	return res
}

// GetPartitionEvents returns all raw events on a physical partition [0..31]
func (b *Bus) GetPartitionEvents(partition int32) []EventEnvelope {
	b.mu.RLock()
	defer b.mu.RUnlock()

	items := b.partitions[partition]
	res := make([]EventEnvelope, len(items))
	copy(res, items)
	return res
}

// RouteToDeadLetter routes a repeatedly failing message to the dead-letter queue
func (b *Bus) RouteToDeadLetter(env EventEnvelope) {
	b.mu.Lock()
	defer b.mu.Unlock()

	env.Topic = b.cfg.DLQTopic
	b.deadLetters = append(b.deadLetters, env)
}

func (b *Bus) GetDeadLetters() []EventEnvelope {
	b.mu.RLock()
	defer b.mu.RUnlock()

	res := make([]EventEnvelope, len(b.deadLetters))
	copy(res, b.deadLetters)
	return res
}

// CommitOffset updates the committed offset for a consumer group on a given partition
func (b *Bus) CommitOffset(groupID string, partition int32, offset int64) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, ok := b.consumerLags[groupID]; !ok {
		b.consumerLags[groupID] = make(map[int32]int64)
	}
	b.consumerLags[groupID][partition] = offset
}

// GetConsumerLag calculates total unconsumed events across all 32 partitions for a consumer group
func (b *Bus) GetConsumerLag(groupID string) int64 {
	b.mu.RLock()
	defer b.mu.RUnlock()

	var totalLag int64
	committedGroup, exists := b.consumerLags[groupID]

	for p := int32(0); p < b.cfg.Partitions; p++ {
		highWatermark := b.offsets[p]
		var committed int64 = 0
		if exists {
			committed = committedGroup[p]
		}
		if highWatermark > committed {
			totalLag += (highWatermark - committed)
		}
	}
	return totalLag
}

func (b *Bus) TotalEvents() uint64 {
	return atomic.LoadUint64(&b.totalEvents)
}

func (b *Bus) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.closed = true
	return nil
}

// Checksum computes the SHA-256 integrity hash of an event payload
func Checksum(data []byte, offset int64) [32]byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(offset))
	return sha256.Sum256(append(buf, data...))
}
