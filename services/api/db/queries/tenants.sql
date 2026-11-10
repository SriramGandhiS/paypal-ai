-- name: GetTenantByID :one
SELECT id, name, slug, status, created_at, updated_at
  FROM tenants
 WHERE id = $1;

-- name: GetTenantBySlug :one
SELECT id, name, slug, status, created_at, updated_at
  FROM tenants
 WHERE slug = $1;
