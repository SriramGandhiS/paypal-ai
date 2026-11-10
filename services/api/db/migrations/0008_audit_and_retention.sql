-- +goose Up
-- Migration 0008: Cryptographic Audit Hash Chains, Merkle Archive Anchors, Retention & Notifications
-- PayBridge AI — Group 0008 (Tables 31-35 of 40)

-- 1. Audit Chain Heads (Transactional Lock Table 31/40)
CREATE TABLE audit_chain_heads (
    tenant_id UUID PRIMARY KEY REFERENCES tenants(id) ON DELETE RESTRICT,
    last_seq BIGINT NOT NULL DEFAULT 0,
    last_event_hash CHAR(64) NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp()
);
ALTER TABLE audit_chain_heads OWNER TO paybridge_admin;

-- 2. Audit Events (Append-Only Linear Hash Chain Table 32/40)
CREATE TABLE audit_events (
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    seq BIGINT NOT NULL CHECK (seq > 0),
    case_id UUID REFERENCES payment_cases(id) ON DELETE SET NULL,
    actor_id UUID,
    actor_role TEXT,
    operation TEXT NOT NULL,
    environment TEXT,
    reason TEXT,
    prev_version BIGINT,
    new_version BIGINT,
    evidence_ids UUID[] NOT NULL DEFAULT '{}',
    request_id TEXT,
    policy_version INT NOT NULL DEFAULT 1,
    occurred_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
    prev_hash CHAR(64) NOT NULL,
    event_hash CHAR(64) NOT NULL,
    payload JSONB NOT NULL DEFAULT '{}'::jsonb,
    PRIMARY KEY (tenant_id, seq)
);
ALTER TABLE audit_events OWNER TO paybridge_admin;
CREATE INDEX ix_audit_case ON audit_events(tenant_id, case_id, seq DESC);
CREATE INDEX ix_audit_occurred ON audit_events(tenant_id, occurred_at);

-- 3. Audit Archive Anchors (Cold Merkle Manifest Anchors Table 33/40)
CREATE TABLE audit_archive_anchors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    from_seq BIGINT NOT NULL,
    to_seq BIGINT NOT NULL,
    events_count BIGINT NOT NULL CHECK (events_count > 0),
    merkle_root CHAR(64) NOT NULL,
    manifest_sha256 CHAR(64) NOT NULL,
    ed25519_key_id TEXT NOT NULL,
    signature_hex TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
    UNIQUE (tenant_id, from_seq, to_seq)
);
ALTER TABLE audit_archive_anchors OWNER TO paybridge_admin;
CREATE INDEX ix_archive_anchors_tenant ON audit_archive_anchors(tenant_id, created_at DESC);

-- 4. Retention Jobs (Data Minimization Purge Table 34/40)
CREATE TABLE retention_jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    policy_type TEXT NOT NULL CHECK (policy_type IN ('PAYLOAD_ERASURE', 'LEGAL_PURGE', 'TTL_EXPIRY')),
    status TEXT NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'APPROVED', 'RUNNING', 'COMPLETED', 'FAILED')),
    target_evidence_id UUID REFERENCES evidence_items(id) ON DELETE CASCADE,
    approved_by UUID REFERENCES users(id) ON DELETE SET NULL,
    approved_at TIMESTAMPTZ,
    purged_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp()
);
ALTER TABLE retention_jobs OWNER TO paybridge_admin;
CREATE INDEX ix_retention_jobs_status ON retention_jobs(status, created_at);

-- 5. Notifications (Tenant User Alerts Table 35/40)
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    case_id UUID REFERENCES payment_cases(id) ON DELETE CASCADE,
    channel TEXT NOT NULL CHECK (channel IN ('IN_APP', 'EMAIL', 'WEBHOOK')),
    title TEXT NOT NULL,
    body TEXT NOT NULL,
    read_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp()
);
ALTER TABLE notifications OWNER TO paybridge_admin;
CREATE INDEX ix_notifications_user ON notifications(user_id, read_at, created_at DESC);

-- +goose Down
DROP TABLE IF EXISTS notifications;
DROP TABLE IF EXISTS retention_jobs;
DROP TABLE IF EXISTS audit_archive_anchors;
DROP TABLE IF EXISTS audit_events;
DROP TABLE IF EXISTS audit_chain_heads;
