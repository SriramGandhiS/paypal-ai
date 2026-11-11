package audit

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	ErrTamperedAuditChain = errors.New("audit invariant violation: cryptographic hash chain verification failed")
)

type Event struct {
	Seq           int64     `json:"seq"`
	TenantID      uuid.UUID `json:"tenant_id"`
	CaseID        uuid.UUID `json:"case_id"`
	ActorID       uuid.UUID `json:"actor_id"`
	ActorRole     string    `json:"actor_role"`
	Operation     string    `json:"operation"`
	Environment   string    `json:"environment"`
	Reason        string    `json:"reason"`
	PrevVersion   int64     `json:"prev_version"`
	NewVersion    int64     `json:"new_version"`
	EvidenceIDs   []string  `json:"evidence_ids"`
	RequestID     string    `json:"request_id"`
	PolicyVersion int       `json:"policy_version"`
	OccurredAt    time.Time `json:"occurred_at"`
	PrevHash      string    `json:"prev_hash"`
	RowHash       string    `json:"row_hash"`
}

type Chain struct {
	mu       sync.RWMutex
	events   []Event
	lastHash string
}

func NewChain() *Chain {
	genesisHash := fmt.Sprintf("%x", sha256.Sum256([]byte("paybridge_genesis_audit_block_2026")))
	return &Chain{
		events:   make([]Event, 0),
		lastHash: genesisHash,
	}
}

// computeHash calculates the cryptographic SHA-256 hash of (prev_hash || canonical_row)
func computeHash(prevHash string, seq int64, tenantID, caseID uuid.UUID, op string, occurredAt time.Time) string {
	raw := fmt.Sprintf("%s|%d|%s|%s|%s|%d",
		prevHash, seq, tenantID.String(), caseID.String(), op, occurredAt.UnixNano())
	h := sha256.Sum256([]byte(raw))
	return fmt.Sprintf("%x", h)
}

// Append records a new immutable audit record into the hash chain
func (c *Chain) Append(
	ctx context.Context,
	tenantID, caseID, actorID uuid.UUID,
	actorRole, op, env, reason string,
	prevVersion, newVersion int64,
	requestID string,
) (*Event, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	seq := int64(len(c.events) + 1)
	now := time.Now().UTC()
	prevHash := c.lastHash

	rowHash := computeHash(prevHash, seq, tenantID, caseID, op, now)

	event := Event{
		Seq:           seq,
		TenantID:      tenantID,
		CaseID:        caseID,
		ActorID:       actorID,
		ActorRole:     actorRole,
		Operation:     op,
		Environment:   env,
		Reason:        reason,
		PrevVersion:   prevVersion,
		NewVersion:    newVersion,
		RequestID:     requestID,
		PolicyVersion: 1,
		OccurredAt:    now,
		PrevHash:      prevHash,
		RowHash:       rowHash,
	}

	c.events = append(c.events, event)
	c.lastHash = rowHash

	return &event, nil
}

// VerifyChain verifies the integrity of the entire audit chain from genesis to head
func (c *Chain) VerifyChain(ctx context.Context) (bool, int, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	expectedPrev := fmt.Sprintf("%x", sha256.Sum256([]byte("paybridge_genesis_audit_block_2026")))

	for i, e := range c.events {
		if e.PrevHash != expectedPrev {
			return false, i, fmt.Errorf("%w: broken link at sequence %d", ErrTamperedAuditChain, e.Seq)
		}

		expectedHash := computeHash(e.PrevHash, e.Seq, e.TenantID, e.CaseID, e.Operation, e.OccurredAt)
		if e.RowHash != expectedHash {
			return false, i, fmt.Errorf("%w: hash mismatch at sequence %d", ErrTamperedAuditChain, e.Seq)
		}

		expectedPrev = e.RowHash
	}

	return true, len(c.events), nil
}

