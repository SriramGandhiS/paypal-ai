-- +goose Up
-- Migration 0004: Evidence Store, Encrypted Payloads, Provenance Edges, Payment Events & Webhook Inbox
-- PayBridge AI — Group 0004 (Tables 14-18 of 40)

-- 1. Evidence Items (Metadata Table 14/40)
CREATE TABLE evidence_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    case_id UUID NOT NULL REFERENCES payment_cases(id) ON DELETE RESTRICT,
    kind TEXT NOT NULL CHECK (kind IN (
        'PROVIDER_GET_ORDER', 'PROVIDER_GET_CAPTURE', 'PROVIDER_WEBHOOK', 'APP_TRACE',
        'MERCHANT_DB_SNAPSHOT', 'PUBLIC_STATUS', 'USER_REPORT', 'SIMULATION', 'ROUTE_RULE', 'FX_QUOTE'
    )),
    tier SMALLINT NOT NULL CHECK (tier BETWEEN 1 AND 6),
    source_trust TEXT NOT NULL CHECK (source_trust IN (
        'PROVIDER_VERIFIED', 'LOCAL_OBSERVED', 'PUBLIC_SOURCE', 'USER_REPORTED', 'SIMULATED'
    )),
    environment TEXT NOT NULL,
    verified BOOLEAN NOT NULL DEFAULT FALSE,
    untrusted_text BOOLEAN NOT NULL DEFAULT FALSE,
    occurred_at TIMESTAMPTZ,
    observed_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
    redacted_extract JSONB NOT NULL DEFAULT '{}'::jsonb,
    payload_sha256 CHAR(64),
    encrypted_payload_ref TEXT,
    access_class TEXT NOT NULL DEFAULT 'TENANT' CHECK (access_class IN ('TENANT', 'SENDER', 'RECEIVER', 'OPS_ONLY')),
    CHECK (source_trust <> 'SIMULATED' OR tier = 6)
);
ALTER TABLE evidence_items OWNER TO paybridge_admin;
CREATE INDEX ix_evidence_case ON evidence_items(case_id, observed_at);
CREATE INDEX ix_evidence_tenant_kind ON evidence_items(tenant_id, kind);

-- 2. Evidence Payloads (Envelope-Encrypted Raw Payloads Table 15/40)
CREATE TABLE evidence_payloads (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    evidence_id UUID UNIQUE NOT NULL REFERENCES evidence_items(id) ON DELETE CASCADE,
    kek_key_id TEXT NOT NULL,
    encrypted_dek BYTEA NOT NULL,
    encrypted_data BYTEA NOT NULL,
    nonce BYTEA NOT NULL,
    sha256_hash CHAR(64) NOT NULL,
    bytes_length INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
    expires_at TIMESTAMPTZ NOT NULL
);
ALTER TABLE evidence_payloads OWNER TO paybridge_admin;
CREATE INDEX ix_payloads_expires ON evidence_payloads(expires_at);

-- 3. Evidence Edges (Truth Gap Graph Table 16/40)
CREATE TABLE evidence_edges (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    case_id UUID NOT NULL REFERENCES payment_cases(id) ON DELETE CASCADE,
    from_evidence_id UUID NOT NULL REFERENCES evidence_items(id) ON DELETE RESTRICT,
    to_evidence_id UUID NOT NULL REFERENCES evidence_items(id) ON DELETE RESTRICT,
    relation TEXT NOT NULL CHECK (relation IN (
        'CORROBORATES', 'CONTRADICTS', 'PRECEDES', 'DERIVED_FROM', 'REQUIRES', 'RESOLVES', 'INDEPENDENT_OF'
    )),
    rule_id TEXT NOT NULL,
    rule_version INT NOT NULL,
    origin TEXT NOT NULL CHECK (origin IN ('DETERMINISTIC', 'LLM_HYPOTHESIS')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp()
);
ALTER TABLE evidence_edges OWNER TO paybridge_admin;
CREATE INDEX ix_evidence_edges_case ON evidence_edges(case_id);

-- 4. Payment Events (Append-Only Event Store Table 17/40)
CREATE TABLE payment_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    case_id UUID NOT NULL REFERENCES payment_cases(id) ON DELETE RESTRICT,
    evidence_id UUID NOT NULL REFERENCES evidence_items(id) ON DELETE RESTRICT,
    provider TEXT,
    provider_event_id TEXT,
    environment TEXT NOT NULL,
    event_type TEXT NOT NULL,
    status_observed TEXT,
    source_trust TEXT NOT NULL,
    verified BOOLEAN NOT NULL DEFAULT FALSE,
    provider_occurred_at TIMESTAMPTZ,
    received_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
    schema_version INT NOT NULL DEFAULT 1
);
ALTER TABLE payment_events OWNER TO paybridge_admin;
CREATE UNIQUE INDEX uq_verified_provider_events
    ON payment_events(provider, environment, provider_event_id)
    WHERE provider_event_id IS NOT NULL AND verified = TRUE;
CREATE INDEX ix_case_events ON payment_events(case_id, received_at, id);

-- 5. Webhook Inbox (Raw Ingest Buffer Table 18/40)
CREATE TABLE webhook_inbox (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    environment TEXT NOT NULL,
    provider TEXT NOT NULL DEFAULT 'PAYPAL',
    provider_event_id TEXT,
    transmission_id TEXT,
    event_type TEXT,
    signature_verified BOOLEAN NOT NULL DEFAULT FALSE,
    payload_sha256 CHAR(64) NOT NULL,
    raw_payload TEXT NOT NULL,
    processing_status TEXT NOT NULL DEFAULT 'RECEIVED' CHECK (processing_status IN (
        'RECEIVED', 'REJECTED', 'CLAIMED', 'UNCORRELATED', 'PROCESSED', 'DUPLICATE', 'DEAD'
    )),
    duplicate_count INT NOT NULL DEFAULT 0,
    attempts INT NOT NULL DEFAULT 0,
    lease_until TIMESTAMPTZ,
    leased_by TEXT,
    case_id UUID REFERENCES payment_cases(id) ON DELETE SET NULL,
    received_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
    processed_at TIMESTAMPTZ,
    UNIQUE (environment, provider, provider_event_id)
);
ALTER TABLE webhook_inbox OWNER TO paybridge_admin;
CREATE INDEX ix_inbox_ready ON webhook_inbox(processing_status, lease_until, received_at);

-- +goose Down
DROP TABLE IF EXISTS webhook_inbox;
DROP TABLE IF EXISTS payment_events;
DROP TABLE IF EXISTS evidence_edges;
DROP TABLE IF EXISTS evidence_payloads;
DROP TABLE IF EXISTS evidence_items;
