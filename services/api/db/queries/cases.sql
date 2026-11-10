-- name: GetCaseByID :one
SELECT id, tenant_id, external_reference, sender_user_id, receiver_user_id,
       sender_country, receiver_country, payment_purpose, amount, currency,
       environment, provider_state, last_confirmed_provider_state,
       provider_state_evidence_id, provider_state_observed_at, application_state,
       receiver_state, case_status, state_version, created_at, updated_at
  FROM payment_cases
 WHERE id = $1 AND tenant_id = $2;

-- name: ListCasesByTenant :many
SELECT id, tenant_id, external_reference, sender_user_id, receiver_user_id,
       amount, currency, provider_state, application_state, receiver_state,
       case_status, state_version, created_at, updated_at
  FROM payment_cases
 WHERE tenant_id = $1
 ORDER BY updated_at DESC
 LIMIT $2 OFFSET $3;

-- name: GetCaseSnapshotByVersion :one
SELECT id, tenant_id, case_id, version, provider_state, application_state,
       receiver_state, case_status, snapshot_payload, created_at
  FROM case_snapshots
 WHERE case_id = $1 AND version = $2 AND tenant_id = $3;
