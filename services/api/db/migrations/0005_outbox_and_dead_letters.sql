-- +goose Up
-- Migration 0005: Outbox Queue & Poison Pill Dead Letters
-- PayBridge AI — Group 0005 (Tables 19-20 of 40)

-- 1. Outbox Queue (Reliable Dispatch Table 19/40)
CREATE TABLE outbox (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    case_id UUID REFERENCES payment_cases(id) ON DELETE SET NULL,
    connection_id UUID, -- Foreign key added after provider_connections in group 0009
    job_type TEXT NOT NULL CHECK (job_type IN (
        'RECONCILE', 'INVESTIGATE', 'NOTIFY', 'PASSPORT_REFRESH',
        'RETENTION', 'SSE_PUBLISH', 'CONNECTOR_DISPATCH'
    )),
    provider_type TEXT,
    payload JSONB NOT NULL,
    dedupe_key TEXT,
    priority INT NOT NULL DEFAULT 100,
    status TEXT NOT NULL DEFAULT 'PENDING' CHECK (status IN (
        'PENDING', 'CLAIMED', 'RUNNING', 'COMPLETED', 'FAILED', 'DEAD'
    )),
    attempts INT NOT NULL DEFAULT 0,
    max_attempts INT NOT NULL DEFAULT 5,
    run_after TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
    locked_until TIMESTAMPTZ,
    locked_by TEXT,
    last_error TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp()
);
ALTER TABLE outbox OWNER TO paybridge_admin;
CREATE INDEX ix_outbox_ready ON outbox(status, run_after, priority DESC) WHERE status IN ('PENDING', 'CLAIMED');
CREATE UNIQUE INDEX uq_outbox_dedupe ON outbox(job_type, dedupe_key)
    WHERE dedupe_key IS NOT NULL AND status IN ('PENDING', 'CLAIMED', 'RUNNING');
CREATE INDEX ix_outbox_tenant ON outbox(tenant_id, created_at DESC);

-- 2. Dead Letters (Poison Pill Sink Table 20/40)
CREATE TABLE dead_letters (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    outbox_id UUID NOT NULL,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    case_id UUID,
    job_type TEXT NOT NULL,
    payload JSONB NOT NULL,
    error_message TEXT NOT NULL,
    attempts_made INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp()
);
ALTER TABLE dead_letters OWNER TO paybridge_admin;
CREATE INDEX ix_dead_letters_tenant ON dead_letters(tenant_id, created_at DESC);

-- +goose Down
DROP TABLE IF EXISTS dead_letters;
DROP TABLE IF EXISTS outbox;
