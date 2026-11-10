-- name: GetUserByEmail :one
SELECT id, email, display_name, password_hash, status, created_at, updated_at
  FROM users
 WHERE email = $1;

-- name: GetUserByID :one
SELECT id, email, display_name, password_hash, status, created_at, updated_at
  FROM users
 WHERE id = $1;

-- name: GetActiveMembershipsForUser :many
SELECT tenant_id, user_id, role, status
  FROM tenant_memberships
 WHERE user_id = $1 AND status = 'ACTIVE';

-- name: GetActiveSigningKey :one
SELECT key_id, key_material, status, created_at, expires_at
  FROM auth_signing_keys
 WHERE status = 'ACTIVE' AND expires_at > clock_timestamp()
 LIMIT 1;
