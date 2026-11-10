-- name: GetAuditChainHead :one
SELECT tenant_id, last_seq, last_event_hash, updated_at
  FROM audit_chain_heads
 WHERE tenant_id = $1;

-- name: ListAuditEventsForCase :many
SELECT seq, tenant_id, case_id, actor_id, actor_role, operation,
       environment, reason, prev_version, new_version, occurred_at,
       prev_hash, event_hash, payload
  FROM audit_events
 WHERE tenant_id = $1 AND case_id = $2
 ORDER BY seq ASC;

-- name: ListAuditEventsFromSeq :many
SELECT seq, tenant_id, case_id, actor_id, actor_role, operation,
       environment, reason, prev_version, new_version, occurred_at,
       prev_hash, event_hash, payload
  FROM audit_events
 WHERE tenant_id = $1 AND seq >= $2
 ORDER BY seq ASC
 LIMIT $3;
