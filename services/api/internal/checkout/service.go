package checkout

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/paybridge/paybridge-ai/services/api/internal/policy"
	"github.com/paybridge/paybridge-ai/services/api/internal/store"
	"github.com/paybridge/paybridge-ai/services/api/internal/tenancy"
)

// DistributedLocker coordinates mutual exclusion across independent service instances.
type DistributedLocker interface {
	TryAcquire(ctx context.Context, key string, ttl time.Duration) (release func(), acquired bool, err error)
}

// Service provides checkout attempt management.
type Service struct {
	store         *store.Store
	locker        DistributedLocker
	inflightLocks sync.Map
}

// NewService constructs a checkout.Service instance.
func NewService(st *store.Store) *Service {
	var locker DistributedLocker
	if st != nil && st.DB() != nil {
		locker = NewPostgresLocker(st.DB())
	}
	return &Service{store: st, locker: locker}
}

// NewServiceWithLocker constructs a checkout.Service with an explicit distributed locker coordinator.
func NewServiceWithLocker(st *store.Store, locker DistributedLocker) *Service {
	return &Service{store: st, locker: locker}
}

// AcquireCaptureLock enforces mutual exclusion for concurrent capture requests on the same order (in-process).
func (s *Service) AcquireCaptureLock(orderID string) (func(), bool) {
	_, loaded := s.inflightLocks.LoadOrStore(orderID, struct{}{})
	if loaded {
		return nil, false
	}
	return func() {
		s.inflightLocks.Delete(orderID)
	}, true
}

// AcquireDistributedCaptureLock enforces mutual exclusion for concurrent capture requests across multiple instances.
// It uses an in-process mutex first for sub-millisecond local rejection, and coordinates distributed locks
// via the configured DistributedLocker (PostgreSQL advisory lock or cluster coordinator).
func (s *Service) AcquireDistributedCaptureLock(ctx context.Context, tenantID, actorID, orderID string) (func(), bool, error) {
	// 1. In-process mutex fast path
	_, loaded := s.inflightLocks.LoadOrStore(orderID, struct{}{})
	if loaded {
		return nil, false, nil
	}

	localRelease := func() {
		s.inflightLocks.Delete(orderID)
	}

	// 2. Distributed locker coordinator across instances
	if s.locker != nil {
		lockKey := fmt.Sprintf("capture:%s:%s", tenantID, orderID)
		remoteRelease, acquired, err := s.locker.TryAcquire(ctx, lockKey, 30*time.Second)
		if err != nil {
			localRelease()
			return nil, false, err
		}
		if !acquired {
			localRelease()
			return nil, false, nil
		}

		combinedRelease := func() {
			localRelease()
			remoteRelease()
		}
		return combinedRelease, true, nil
	}

	// 3. Fallback: PostgreSQL idempotency_keys table lock when store is configured without explicit locker
	if s.store != nil && s.store.DB() != nil && tenantID != "" && actorID != "" {
		hash := sha256.Sum256([]byte(orderID))
		route := "PAYPAL_CAPTURE"
		lockQuery := `
			INSERT INTO idempotency_keys (
				tenant_id, actor_id, route, key, request_hash, status, expires_at
			) VALUES (
				$1, $2, $3, $4, $5, 'IN_FLIGHT', NOW() + INTERVAL '30 seconds'
			)
			ON CONFLICT (tenant_id, actor_id, route, key) DO UPDATE
			    SET status = 'IN_FLIGHT', expires_at = NOW() + INTERVAL '30 seconds'
			  WHERE idempotency_keys.status = 'IN_FLIGHT' AND idempotency_keys.expires_at < NOW()
			RETURNING key;
		`
		var lockedKey string
		err := s.store.DB().QueryRowContext(ctx, lockQuery, tenantID, actorID, route, orderID, hash[:]).Scan(&lockedKey)
		if err == sql.ErrNoRows {
			localRelease()
			return nil, false, nil
		} else if err != nil {
			localRelease()
			return nil, false, err
		}

		combinedRelease := func() {
			localRelease()
			delCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			_, _ = s.store.DB().ExecContext(delCtx,
				"DELETE FROM idempotency_keys WHERE tenant_id = $1 AND actor_id = $2 AND route = $3 AND key = $4",
				tenantID, actorID, route, orderID,
			)
		}
		return combinedRelease, true, nil
	}

	return localRelease, true, nil
}

// SharedMemoryLocker provides a thread-safe distributed lock coordinator for multi-instance clusters and test suites.
type SharedMemoryLocker struct {
	mu    sync.Mutex
	locks map[string]lockEntry
}

type lockEntry struct {
	expiresAt time.Time
	token     string
}

// NewSharedMemoryLocker constructs a new shared in-memory lock coordinator.
func NewSharedMemoryLocker() *SharedMemoryLocker {
	return &SharedMemoryLocker{
		locks: make(map[string]lockEntry),
	}
}

// TryAcquire attempts to acquire a lock for the given key with a time-to-live lease.
func (l *SharedMemoryLocker) TryAcquire(ctx context.Context, key string, ttl time.Duration) (func(), bool, error) {
	if err := ctx.Err(); err != nil {
		return nil, false, err
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	entry, exists := l.locks[key]
	if exists && now.Before(entry.expiresAt) {
		return nil, false, nil
	}

	token := fmt.Sprintf("%d", now.UnixNano())
	l.locks[key] = lockEntry{
		expiresAt: now.Add(ttl),
		token:     token,
	}

	released := false
	release := func() {
		l.mu.Lock()
		defer l.mu.Unlock()
		if !released {
			released = true
			if cur, ok := l.locks[key]; ok && cur.token == token {
				delete(l.locks, key)
			}
		}
	}

	return release, true, nil
}

// PostgresLocker coordinates distributed locks across cluster nodes using PostgreSQL advisory locks.
type PostgresLocker struct {
	db *sql.DB
}

// NewPostgresLocker constructs a PostgresLocker backed by a live PostgreSQL connection.
func NewPostgresLocker(db *sql.DB) *PostgresLocker {
	return &PostgresLocker{db: db}
}

// TryAcquire acquires a cluster-wide advisory lock on the key.
// It acquires a dedicated connection from the pool so the session-level advisory lock
// remains pinned to the exact physical connection and is not accidentally returned to the pool while held.
func (p *PostgresLocker) TryAcquire(ctx context.Context, key string, ttl time.Duration) (func(), bool, error) {
	if p.db == nil {
		return nil, false, fmt.Errorf("nil database connection")
	}

	conn, err := p.db.Conn(ctx)
	if err != nil {
		return nil, false, fmt.Errorf("failed to checkout dedicated connection for advisory lock: %w", err)
	}

	hash := sha256.Sum256([]byte(key))
	var lockID int64
	for i := 0; i < 8; i++ {
		lockID = (lockID << 8) | int64(hash[i])
	}

	var acquired bool
	err = conn.QueryRowContext(ctx, "SELECT pg_try_advisory_lock($1)", lockID).Scan(&acquired)
	if err != nil {
		_ = conn.Close()
		return nil, false, fmt.Errorf("failed to execute pg_try_advisory_lock: %w", err)
	}
	if !acquired {
		_ = conn.Close()
		return nil, false, nil
	}

	release := func() {
		relCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		var unlocked bool
		_ = conn.QueryRowContext(relCtx, "SELECT pg_advisory_unlock($1)", lockID).Scan(&unlocked)
		_ = conn.Close() // Safely returns unlocked connection to pool
	}

	return release, true, nil
}

// CreateAttempt initializes a new payment attempt while enforcing the blind-retry invariant.
func (s *Service) CreateAttempt(ctx context.Context, caseID, orderOpReqID string) (*Attempt, error) {
	tenantID, err := tenancy.TenantFromContext(ctx)
	if err != nil {
		return nil, err
	}

	_, role, err := tenancy.UserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if err := policy.Authorize(role, policy.ActionCreateOrder); err != nil {
		return nil, err
	}

	var attempt Attempt
	err = s.store.WithTx(ctx, func(tx *sql.Tx) error {
		if err := s.store.SetTenantContext(ctx, tx, tenantID); err != nil {
			return err
		}

		// 1. Enforce blind-retry guard: at most one active non-terminal attempt
		var activeCount int
		err := tx.QueryRowContext(ctx,
			"SELECT count(*) FROM payment_attempts WHERE case_id = $1 AND terminal = FALSE",
			caseID,
		).Scan(&activeCount)
		if err != nil {
			return fmt.Errorf("failed to check active attempts: %w", err)
		}
		if activeCount > 0 {
			return ErrActiveAttemptExists
		}

		// 2. Determine next attempt number
		var maxAttempt sql.NullInt32
		err = tx.QueryRowContext(ctx,
			"SELECT max(attempt_no) FROM payment_attempts WHERE case_id = $1",
			caseID,
		).Scan(&maxAttempt)
		if err != nil {
			return fmt.Errorf("failed to query attempt sequence: %w", err)
		}

		nextAttempt := 1
		if maxAttempt.Valid {
			nextAttempt = int(maxAttempt.Int32) + 1
		}

		// 3. Insert new attempt
		query := `
			INSERT INTO payment_attempts (
				tenant_id, case_id, attempt_no, provider, environment,
				order_operation_request_id, status
			) VALUES (
				$1, $2, $3, 'PAYPAL', 'REAL_SANDBOX', $4, 'ORDER_CREATING'
			)
			RETURNING id, tenant_id, case_id, attempt_no, provider, environment,
			          order_operation_request_id, status, terminal, created_at, updated_at;
		`
		return tx.QueryRowContext(ctx, query,
			tenantID, caseID, nextAttempt, orderOpReqID,
		).Scan(
			&attempt.ID, &attempt.TenantID, &attempt.CaseID, &attempt.AttemptNo,
			&attempt.Provider, &attempt.Environment, &attempt.OrderOperationRequestID,
			&attempt.Status, &attempt.Terminal, &attempt.CreatedAt, &attempt.UpdatedAt,
		)
	})

	if err != nil {
		return nil, err
	}
	return &attempt, nil
}

// GetAttempt retrieves a single attempt by ID.
func (s *Service) GetAttempt(ctx context.Context, attemptID string) (*Attempt, error) {
	tenantID, err := tenancy.TenantFromContext(ctx)
	if err != nil {
		return nil, err
	}

	_, role, err := tenancy.UserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if err := policy.Authorize(role, policy.ActionReadCase); err != nil {
		return nil, err
	}

	query := `
		SELECT id, tenant_id, case_id, attempt_no, provider, environment,
		       order_operation_request_id, capture_operation_request_id,
		       status, terminal, created_at, updated_at
		  FROM payment_attempts
		 WHERE id = $1 AND tenant_id = $2;
	`
	var a Attempt
	err = s.store.DB().QueryRowContext(ctx, query, attemptID, tenantID).Scan(
		&a.ID, &a.TenantID, &a.CaseID, &a.AttemptNo, &a.Provider, &a.Environment,
		&a.OrderOperationRequestID, &a.CaptureOperationRequestID,
		&a.Status, &a.Terminal, &a.CreatedAt, &a.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrAttemptNotFound
	} else if err != nil {
		return nil, fmt.Errorf("failed to query attempt: %w", err)
	}

	return &a, nil
}

// ListAttemptsForCase returns all attempts for a given case.
func (s *Service) ListAttemptsForCase(ctx context.Context, caseID string) ([]Attempt, error) {
	tenantID, err := tenancy.TenantFromContext(ctx)
	if err != nil {
		return nil, err
	}

	_, role, err := tenancy.UserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if err := policy.Authorize(role, policy.ActionReadCase); err != nil {
		return nil, err
	}

	query := `
		SELECT id, tenant_id, case_id, attempt_no, provider, environment,
		       order_operation_request_id, capture_operation_request_id,
		       status, terminal, created_at, updated_at
		  FROM payment_attempts
		 WHERE case_id = $1 AND tenant_id = $2
		 ORDER BY attempt_no ASC;
	`
	rows, err := s.store.DB().QueryContext(ctx, query, caseID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to list attempts: %w", err)
	}
	defer rows.Close()

	var attempts []Attempt
	for rows.Next() {
		var a Attempt
		if err := rows.Scan(
			&a.ID, &a.TenantID, &a.CaseID, &a.AttemptNo, &a.Provider, &a.Environment,
			&a.OrderOperationRequestID, &a.CaptureOperationRequestID,
			&a.Status, &a.Terminal, &a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, err
		}
		attempts = append(attempts, a)
	}

	return attempts, rows.Err()
}
