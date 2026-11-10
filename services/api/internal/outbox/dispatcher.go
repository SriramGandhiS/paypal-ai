package outbox

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/paybridge/paybridge-ai/services/api/internal/store"
)

// Job represents a background task dispatched via the outbox pattern.
type Job struct {
	ID           string
	TenantID     string
	CaseID       *string
	ConnectionID *string
	JobType      string
	ProviderType *string
	Payload      json.RawMessage
	Attempts     int
}

// Handler is an interface for processing specific outbox job types.
type Handler interface {
	HandleJob(ctx context.Context, job Job) error
}

// Dispatcher manages outbox polling and task execution.
type Dispatcher struct {
	store    *store.Store
	handlers map[string]Handler
	logger   *slog.Logger
}

// NewDispatcher constructs a new outbox Dispatcher.
func NewDispatcher(st *store.Store, logger *slog.Logger) *Dispatcher {
	if logger == nil {
		logger = slog.Default()
	}
	return &Dispatcher{
		store:    st,
		handlers: make(map[string]Handler),
		logger:   logger,
	}
}

// RegisterHandler binds an execution handler to a job type.
func (d *Dispatcher) RegisterHandler(jobType string, h Handler) {
	d.handlers[jobType] = h
}

// ProcessNextBatch claims and executes pending jobs using SKIP LOCKED.
func (d *Dispatcher) ProcessNextBatch(ctx context.Context, workerID string, batchSize int) (int, error) {
	if d.store == nil || d.store.DB() == nil {
		return 0, nil
	}

	var processed int
	err := d.store.WithTx(ctx, func(tx *sql.Tx) error {
		query := `
			UPDATE outbox
			   SET status = 'CLAIMED',
			       locked_until = clock_timestamp() + interval '5 minutes',
			       locked_by = $1,
			       attempts = attempts + 1,
			       updated_at = clock_timestamp()
			 WHERE id IN (
			     SELECT id FROM outbox
			      WHERE status = 'PENDING'
			        AND (locked_until IS NULL OR locked_until < clock_timestamp())
			      ORDER BY priority DESC, created_at ASC
			      LIMIT $2
			      FOR UPDATE SKIP LOCKED
			 )
			 RETURNING id, tenant_id, case_id, connection_id, job_type, provider_type, payload, attempts;
		`
		rows, err := tx.QueryContext(ctx, query, workerID, batchSize)
		if err != nil {
			return fmt.Errorf("failed to claim outbox batch: %w", err)
		}
		defer rows.Close()

		var jobs []Job
		for rows.Next() {
			var j Job
			if err := rows.Scan(&j.ID, &j.TenantID, &j.CaseID, &j.ConnectionID, &j.JobType, &j.ProviderType, &j.Payload, &j.Attempts); err != nil {
				return err
			}
			jobs = append(jobs, j)
		}
		if err := rows.Err(); err != nil {
			return err
		}

		for _, job := range jobs {
			handler, exists := d.handlers[job.JobType]
			if !exists {
				d.logger.Warn("no handler registered for job type", "type", job.JobType, "id", job.ID)
				continue
			}

			if err := handler.HandleJob(ctx, job); err != nil {
				d.logger.Error("job execution failed", "id", job.ID, "type", job.JobType, "error", err)
				_, _ = tx.ExecContext(ctx, `
					UPDATE outbox
					   SET status = CASE WHEN attempts >= max_attempts THEN 'FAILED' ELSE 'PENDING' END,
					       locked_until = NULL,
					       last_error = $2,
					       updated_at = clock_timestamp()
					 WHERE id = $1;
				`, job.ID, err.Error())
				continue
			}

			// Job completed successfully
			_, _ = tx.ExecContext(ctx, `
				UPDATE outbox
				   SET status = 'COMPLETED',
				       locked_until = NULL,
				       updated_at = clock_timestamp()
				 WHERE id = $1;
			`, job.ID)
			processed++
		}
		return nil
	})

	return processed, err
}
