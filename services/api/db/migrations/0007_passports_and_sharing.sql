-- +goose Up
-- Migration 0007: Merchant Payment Passports, Invitations, Consent & Sharing
-- PayBridge AI — Group 0007 (Tables 27-30 of 40)

-- 1. Passport Versions (Merchant Payment Passport Table 27/40)
CREATE TABLE passport_versions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    case_id UUID NOT NULL REFERENCES payment_cases(id) ON DELETE CASCADE,
    viewer_role TEXT NOT NULL CHECK (viewer_role IN ('sender', 'receiver', 'ops')),
    version INT NOT NULL CHECK (version > 0),
    snapshot_version BIGINT NOT NULL,
    content JSONB NOT NULL,
    content_sha256 CHAR(64) NOT NULL,
    prev_sha256 CHAR(64),
    evidence_ids UUID[] NOT NULL DEFAULT '{}',
    policy_version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
    UNIQUE (case_id, viewer_role, version)
);
ALTER TABLE passport_versions OWNER TO paybridge_admin;
CREATE INDEX ix_passport_case_version ON passport_versions(case_id, viewer_role, version DESC);

-- 2. Invitations (Onboarding Invitations Table 28/40)
CREATE TABLE invitations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    case_id UUID NOT NULL REFERENCES payment_cases(id) ON DELETE CASCADE,
    invitee_email CITEXT NOT NULL,
    token_hash BYTEA NOT NULL UNIQUE,
    proposed_fields TEXT[] NOT NULL DEFAULT '{}',
    expires_at TIMESTAMPTZ NOT NULL,
    status TEXT NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'ACCEPTED', 'DECLINED', 'EXPIRED', 'REVOKED')),
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp()
);
ALTER TABLE invitations OWNER TO paybridge_admin;
CREATE INDEX ix_invitations_email ON invitations(invitee_email);
CREATE INDEX ix_invitations_case ON invitations(case_id);

-- 3. Consent Grants (Explicit Privacy Allowlist Table 29/40)
CREATE TABLE consent_grants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    case_id UUID NOT NULL REFERENCES payment_cases(id) ON DELETE CASCADE,
    grantee_user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    field_allowlist TEXT[] NOT NULL DEFAULT '{}',
    notification_allowed BOOLEAN NOT NULL DEFAULT FALSE,
    granted_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
    revoked_at TIMESTAMPTZ,
    version INT NOT NULL DEFAULT 1
);
ALTER TABLE consent_grants OWNER TO paybridge_admin;
CREATE INDEX ix_consent_grants_case ON consent_grants(case_id, grantee_user_id);

-- 4. Passport Shares (Tokenized Capability Shares Table 30/40)
CREATE TABLE passport_shares (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    case_id UUID NOT NULL REFERENCES payment_cases(id) ON DELETE CASCADE,
    viewer_role TEXT NOT NULL,
    token_hash BYTEA NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp()
);
ALTER TABLE passport_shares OWNER TO paybridge_admin;
CREATE INDEX ix_shares_case ON passport_shares(case_id);

-- +goose Down
DROP TABLE IF EXISTS passport_shares;
DROP TABLE IF EXISTS consent_grants;
DROP TABLE IF EXISTS invitations;
DROP TABLE IF EXISTS passport_versions;
