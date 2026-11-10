-- name: ClaimOutboxBatch :many
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

-- name: CompleteOutboxJob :exec
UPDATE outbox
   SET status = 'COMPLETED',
       locked_until = NULL,
       updated_at = clock_timestamp()
 WHERE id = $1;

-- name: FailOutboxJob :exec
UPDATE outbox
   SET status = CASE WHEN attempts >= max_attempts THEN 'FAILED' ELSE 'PENDING' END,
       locked_until = NULL,
       last_error = $2,
       updated_at = clock_timestamp()
 WHERE id = $1;
