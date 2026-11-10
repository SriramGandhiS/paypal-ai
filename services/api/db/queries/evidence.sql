-- name: ListEvidenceByCase :many
SELECT id, tenant_id, case_id, kind, tier, source_trust, environment,
       verified, untrusted_text, occurred_at, observed_at, redacted_extract,
       payload_sha256, access_class
  FROM evidence_items
 WHERE case_id = $1 AND tenant_id = $2
 ORDER BY observed_at ASC;

-- name: GetEvidenceByID :one
SELECT id, tenant_id, case_id, kind, tier, source_trust, environment,
       verified, untrusted_text, occurred_at, observed_at, redacted_extract,
       payload_sha256, access_class
  FROM evidence_items
 WHERE id = $1 AND tenant_id = $2;
