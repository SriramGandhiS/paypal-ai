-- +goose Up
-- Migration 0009: Route Registry, Status Telemetry, Provider Connections, Credential Access Logs & Simulation Runs
-- PayBridge AI — Group 0009 (Tables 36-40 of 40)

-- 1. Route Registry (Global Rail Reference Table 36/40)
CREATE TABLE route_registry (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    provider TEXT NOT NULL,
    rail_type TEXT NOT NULL CHECK (rail_type IN ('CARD', 'WALLET', 'ACH', 'RTP', 'FEDNOW', 'SEPA', 'UPI')),
    source_country CHAR(2) NOT NULL,
    destination_country CHAR(2) NOT NULL,
    currency CHAR(3) NOT NULL,
    settlement_speed_hours INT NOT NULL DEFAULT 24,
    cost_basis_bps INT NOT NULL DEFAULT 150,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp()
);
ALTER TABLE route_registry OWNER TO paybridge_admin;
CREATE INDEX ix_routes_lookup ON route_registry(source_country, destination_country, currency, is_active);

-- 2. Provider Status Samples (Global Health Telemetry Table 37/40)
CREATE TABLE provider_status_samples (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    provider TEXT NOT NULL,
    environment TEXT NOT NULL,
    endpoint_name TEXT NOT NULL,
    status_code INT NOT NULL,
    latency_ms INT NOT NULL,
    is_available BOOLEAN NOT NULL DEFAULT TRUE,
    sampled_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp()
);
ALTER TABLE provider_status_samples OWNER TO paybridge_admin;
CREATE INDEX ix_provider_samples_time ON provider_status_samples(provider, sampled_at DESC);

-- 3. Provider Connections (Envelope-Encrypted Secret Store Table 38/40)
CREATE TABLE provider_connections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    provider_type TEXT NOT NULL CHECK (provider_type IN ('PAYPAL', 'STRIPE', 'WISE', 'PLAID', 'SETU')),
    environment TEXT NOT NULL CHECK (environment IN ('REAL_SANDBOX', 'PAYPAL_SIMULATOR', 'LOCAL_SIMULATION', 'PRODUCTION')),
    display_name TEXT NOT NULL,
    kek_key_id TEXT NOT NULL,
    encrypted_dek BYTEA NOT NULL,
    encrypted_credentials BYTEA NOT NULL,
    nonce BYTEA NOT NULL,
    key_version INT NOT NULL DEFAULT 1,
    status TEXT NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'REVOKED', 'DISABLED')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp()
);
ALTER TABLE provider_connections OWNER TO paybridge_admin;
CREATE INDEX ix_connections_tenant ON provider_connections(tenant_id, provider_type, status);

-- Add foreign key constraint to outbox.connection_id
ALTER TABLE outbox
    ADD CONSTRAINT fk_outbox_connection
    FOREIGN KEY (connection_id) REFERENCES provider_connections(id) ON DELETE SET NULL;

-- 4. Connector Credential Access Logs (Operational Audit Index Table 39/40)
CREATE TABLE connector_credential_access_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id UUID NOT NULL REFERENCES outbox(id) ON DELETE CASCADE,
    connection_id UUID NOT NULL REFERENCES provider_connections(id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    accessed_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
    client_addr INET
);
ALTER TABLE connector_credential_access_logs OWNER TO paybridge_admin;
CREATE INDEX ix_connector_logs_time ON connector_credential_access_logs(accessed_at DESC);
CREATE INDEX ix_connector_logs_conn ON connector_credential_access_logs(connection_id, accessed_at DESC);

-- 5. Simulation Runs (Synthetic Payment Sandbox Table 40/40)
CREATE TABLE simulation_runs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    case_id UUID REFERENCES payment_cases(id) ON DELETE CASCADE,
    scenario_id TEXT NOT NULL,
    chaos_injected TEXT,
    success BOOLEAN NOT NULL DEFAULT TRUE,
    run_log JSONB NOT NULL DEFAULT '{}'::jsonb,
    started_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
    completed_at TIMESTAMPTZ
);
ALTER TABLE simulation_runs OWNER TO paybridge_admin;
CREATE INDEX ix_simulation_case ON simulation_runs(case_id);

-- +goose Down
ALTER TABLE outbox DROP CONSTRAINT IF EXISTS fk_outbox_connection;
DROP TABLE IF EXISTS simulation_runs;
DROP TABLE IF EXISTS connector_credential_access_logs;
DROP TABLE IF EXISTS provider_connections;
DROP TABLE IF EXISTS provider_status_samples;
DROP TABLE IF EXISTS route_registry;
