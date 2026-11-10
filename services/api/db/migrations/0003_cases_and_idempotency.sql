-- +goose Up
-- Migration 0003: Payment Cases, Attempts, Resources, Idempotency & Snapshots
-- PayBridge AI — Group 0003 (Tables 8-13 of 40)

-- 1. Payment Cases (Core Domain Table 8/40)
CREATE TABLE payment_cases (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    external_reference TEXT, -- Merchant order reference (untrusted text)
    sender_user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    receiver_user_id UUID REFERENCES users(id) ON DELETE RESTRICT,
    sender_country CHAR(2) NOT NULL,
    receiver_country CHAR(2) NOT NULL,
    payment_purpose TEXT NOT NULL,
    amount NUMERIC(20,6) NOT NULL CHECK (amount > 0),
    currency CHAR(3) NOT NULL,
    environment TEXT NOT NULL CHECK (environment IN (
        'REAL_SANDBOX', 'PAYPAL_SIMULATOR', 'LOCAL_SIMULATION', 'PRODUCTION'
    )),
    -- 3D State Vector (Written ONLY via Reducer)
    provider_state TEXT NOT NULL DEFAULT 'NOT_CREATED' CHECK (provider_state IN (
        'NOT_CREATED', 'CREATED', 'APPROVED', 'AUTHORISED', 'CAPTURE_PENDING',
        'CAPTURE_COMPLETED', 'CAPTURE_DENIED', 'REFUNDED', 'PARTIALLY_REFUNDED', 'UNKNOWN'
    )),
    last_confirmed_provider_state TEXT NOT NULL DEFAULT 'NOT_CREATED',
    provider_state_evidence_id UUID,
    provider_state_observed_at TIMESTAMPTZ,
    application_state TEXT NOT NULL DEFAULT 'SYNCHRONISED' CHECK (application_state IN (
        'SYNCHRONISED', 'WAITING_FOR_WEBHOOK', 'API_TIMEOUT', 'RECONCILING', 'EVENT_GAP', 'INCONSISTENT'
    )),
    receiver_state TEXT NOT NULL DEFAULT 'NOT_OBSERVABLE' CHECK (receiver_state IN (
        'NOT_APPLICABLE', 'NOT_OBSERVABLE', 'REPORTED_PENDING', 'VERIFIED_AVAILABLE',
        'WITHDRAWAL_IN_PROGRESS', 'VERIFIED_BANK_CREDITED'
    )),
    case_status TEXT NOT NULL DEFAULT 'OPEN' CHECK (case_status IN (
        'OPEN', 'OBSERVING', 'INVESTIGATING', 'GUIDANCE_READY', 'AWAITING_USER_ACTION',
        'VERIFYING_RESOLUTION', 'RESOLVED', 'NEEDS_HUMAN_REVIEW', 'PAUSED_INSUFFICIENT_EVIDENCE',
        'CLOSED_UNRESOLVED', 'ARCHIVED', 'REOPENED'
    )),
    state_version BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp()
);
ALTER TABLE payment_cases OWNER TO paybridge_admin;
CREATE INDEX ix_cases_tenant_status ON payment_cases(tenant_id, case_status, updated_at DESC);
CREATE INDEX ix_cases_external_ref ON payment_cases(tenant_id, external_reference);

-- 2. Case Participants (Table 9/40)
CREATE TABLE case_participants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    case_id UUID NOT NULL REFERENCES payment_cases(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    participant_type TEXT NOT NULL CHECK (participant_type IN ('SENDER', 'RECEIVER', 'MERCHANT', 'OPERATOR')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
    UNIQUE (tenant_id, case_id, user_id, participant_type)
);
ALTER TABLE case_participants OWNER TO paybridge_admin;
CREATE INDEX ix_participants_user ON case_participants(user_id);

-- 3. Payment Attempts (Table 10/40)
CREATE TABLE payment_attempts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    case_id UUID NOT NULL REFERENCES payment_cases(id) ON DELETE RESTRICT,
    attempt_no INT NOT NULL CHECK (attempt_no > 0),
    provider TEXT NOT NULL DEFAULT 'PAYPAL',
    environment TEXT NOT NULL,
    order_operation_request_id TEXT NOT NULL,
    capture_operation_request_id TEXT,
    status TEXT NOT NULL DEFAULT 'ORDER_CREATING' CHECK (status IN (
        'ORDER_CREATING', 'ORDER_CREATED', 'ORDER_UNKNOWN', 'APPROVAL_PENDING', 'APPROVED',
        'CANCELLED', 'CAPTURE_REQUESTED', 'CAPTURE_COMPLETED', 'CAPTURE_PENDING',
        'CAPTURE_DENIED', 'CAPTURE_UNKNOWN', 'REVIEW_REQUIRED', 'APPLICATION_RECONCILED', 'APPLICATION_MISMATCH'
    )),
    terminal BOOLEAN GENERATED ALWAYS AS (status IN (
        'CANCELLED', 'CAPTURE_DENIED', 'APPLICATION_RECONCILED'
    )) STORED,
    created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
    UNIQUE (tenant_id, case_id, attempt_no),
    UNIQUE (tenant_id, provider, environment, order_operation_request_id)
);
ALTER TABLE payment_attempts OWNER TO paybridge_admin;
-- At most ONE non-terminal attempt per case (Blind Retry Guard)
CREATE UNIQUE INDEX uq_one_active_attempt ON payment_attempts(case_id) WHERE terminal = FALSE;
CREATE INDEX ix_attempts_case ON payment_attempts(case_id, attempt_no);

-- 4. Provider Resources (Provider Mapping Table 11/40)
CREATE TABLE provider_resources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    attempt_id UUID NOT NULL REFERENCES payment_attempts(id) ON DELETE RESTRICT,
    resource_type TEXT NOT NULL CHECK (resource_type IN ('ORDER', 'AUTHORIZATION', 'CAPTURE', 'REFUND')),
    provider_resource_id TEXT NOT NULL,
    environment TEXT NOT NULL,
    provider_status TEXT,
    last_verified_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
    UNIQUE (environment, resource_type, provider_resource_id)
);
ALTER TABLE provider_resources OWNER TO paybridge_admin;
CREATE INDEX ix_provider_resources_attempt ON provider_resources(attempt_id);

-- 5. Idempotency Keys (API Deduplication Table 12/40)
CREATE TABLE idempotency_keys (
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    actor_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    route TEXT NOT NULL,
    key TEXT NOT NULL,
    request_hash BYTEA NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('IN_FLIGHT', 'DONE')),
    response_code INT,
    response_body JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
    expires_at TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (tenant_id, actor_id, route, key)
);
ALTER TABLE idempotency_keys OWNER TO paybridge_admin;
CREATE INDEX ix_idempotency_expires ON idempotency_keys(expires_at);

-- 6. Case Snapshots (Immutable Historical State Table 13/40)
CREATE TABLE case_snapshots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    case_id UUID NOT NULL REFERENCES payment_cases(id) ON DELETE RESTRICT,
    version BIGINT NOT NULL CHECK (version >= 0),
    provider_state TEXT NOT NULL,
    application_state TEXT NOT NULL,
    receiver_state TEXT NOT NULL,
    case_status TEXT NOT NULL,
    causing_event_id UUID,
    snapshot_payload JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
    UNIQUE (case_id, version)
);
ALTER TABLE case_snapshots OWNER TO paybridge_admin;
CREATE INDEX ix_snapshots_case_version ON case_snapshots(case_id, version DESC);

-- +goose Down
DROP TABLE IF EXISTS case_snapshots;
DROP TABLE IF EXISTS idempotency_keys;
DROP TABLE IF EXISTS provider_resources;
DROP TABLE IF EXISTS payment_attempts;
DROP TABLE IF EXISTS case_participants;
DROP TABLE IF EXISTS payment_cases;
