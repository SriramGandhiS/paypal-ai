-- +goose Up
-- Migration 0002: Identity, Tenants, Memberships & Authentication
-- PayBridge AI — Group 0002 (Tables 2-7 of 40)

-- 1. Tenants (Root Directory Table 2/40)
CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    slug TEXT UNIQUE NOT NULL,
    status TEXT NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'SUSPENDED', 'DELETED')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp()
);
ALTER TABLE tenants OWNER TO paybridge_admin;

-- 2. Users (Global Identity Table 3/40)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email CITEXT UNIQUE NOT NULL,
    display_name TEXT NOT NULL,
    password_hash TEXT, -- Argon2id hash for demo/local auth
    status TEXT NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'LOCKED', 'DEACTIVATED')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp()
);
ALTER TABLE users OWNER TO paybridge_admin;

-- 3. Tenant Memberships (RBAC Store Table 4/40)
CREATE TABLE tenant_memberships (
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    role TEXT NOT NULL CHECK (role IN (
        'sender', 'receiver', 'merchant_analyst', 'merchant_admin', 'platform_operator', 'demo_operator'
    )),
    status TEXT NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'INVITED', 'SUSPENDED')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
    PRIMARY KEY (tenant_id, user_id, role)
);
ALTER TABLE tenant_memberships OWNER TO paybridge_admin;
CREATE INDEX ix_memberships_user ON tenant_memberships(user_id);

-- 4. Sessions (Web Session Store Table 5/40)
CREATE TABLE sessions (
    id_hash BYTEA PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
    last_seen_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    ip_hash BYTEA,
    ua_hash BYTEA
);
ALTER TABLE sessions OWNER TO paybridge_admin;
CREATE INDEX ix_sessions_user_tenant ON sessions(user_id, tenant_id);
CREATE INDEX ix_sessions_expires ON sessions(expires_at) WHERE revoked_at IS NULL;

-- 5. Auth Signing Keys (HMAC Secret Store Table 6/40)
CREATE TABLE auth_signing_keys (
    key_id TEXT PRIMARY KEY,
    key_material BYTEA NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('ACTIVE', 'RETIRED', 'REVOKED')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ
);
ALTER TABLE auth_signing_keys OWNER TO paybridge_admin;
REVOKE ALL ON auth_signing_keys FROM PUBLIC, paybridge_app, paybridge_worker, paybridge_connector, paybridge_webhook, paybridge_retention;
GRANT SELECT ON auth_signing_keys TO paybridge_auth;

-- 6. Consumed Auth Nonces (Replay Protection Table 7/40)
CREATE TABLE consumed_auth_nonces (
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    nonce TEXT NOT NULL,
    op TEXT NOT NULL,
    consumed_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
    PRIMARY KEY (tenant_id, nonce)
);
ALTER TABLE consumed_auth_nonces OWNER TO paybridge_admin;
CREATE INDEX ix_consumed_nonces_timestamp ON consumed_auth_nonces(consumed_at);

-- +goose Down
DROP TABLE IF EXISTS consumed_auth_nonces;
DROP TABLE IF EXISTS auth_signing_keys;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS tenant_memberships;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS tenants;
