-- +goose Up
-- Migration 0006: AI Multi-Agent Investigations, Runs, Diagnostic Claims & Safe Actions
-- PayBridge AI — Group 0006 (Tables 21-26 of 40)

-- 1. Investigations (Table 21/40)
CREATE TABLE investigations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    case_id UUID NOT NULL REFERENCES payment_cases(id) ON DELETE CASCADE,
    trigger TEXT NOT NULL,
    snapshot_version BIGINT NOT NULL,
    status TEXT NOT NULL DEFAULT 'QUEUED' CHECK (status IN (
        'QUEUED', 'RUNNING', 'COMPLETED', 'FAILED_FALLBACK', 'STALE_DISCARDED', 'NEEDS_HUMAN_REVIEW'
    )),
    llm_available BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
    completed_at TIMESTAMPTZ
);
ALTER TABLE investigations OWNER TO paybridge_admin;
CREATE INDEX ix_investigations_case ON investigations(case_id, created_at DESC);

-- 2. Agent Runs (Table 22/40)
CREATE TABLE agent_runs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    investigation_id UUID NOT NULL REFERENCES investigations(id) ON DELETE CASCADE,
    agent_id TEXT NOT NULL CHECK (agent_id IN ('A01','A02','A03','A04','A05','A06','A07','A08','A09','A10')),
    model_id TEXT NOT NULL,
    prompt_sha256 CHAR(64) NOT NULL,
    input_evidence_ids UUID[] NOT NULL DEFAULT '{}',
    output JSONB NOT NULL DEFAULT '{}'::jsonb,
    tokens_in INT NOT NULL DEFAULT 0,
    tokens_out INT NOT NULL DEFAULT 0,
    latency_ms INT NOT NULL DEFAULT 0,
    rejection_reasons TEXT[] NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp()
);
ALTER TABLE agent_runs OWNER TO paybridge_admin;
CREATE INDEX ix_agent_runs_investigation ON agent_runs(investigation_id, agent_id);

-- 3. Diagnosis Claims (Table 23/40)
CREATE TABLE diagnosis_claims (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    case_id UUID NOT NULL REFERENCES payment_cases(id) ON DELETE CASCADE,
    investigation_id UUID NOT NULL REFERENCES investigations(id) ON DELETE CASCADE,
    kind TEXT NOT NULL CHECK (kind IN ('confirmed', 'inference', 'unknown')),
    code TEXT NOT NULL,
    statement TEXT NOT NULL,
    inference_rule TEXT,
    viewer_scope TEXT NOT NULL CHECK (viewer_scope IN ('SENDER', 'RECEIVER', 'OPS', 'ALL')),
    validated BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp()
);
ALTER TABLE diagnosis_claims OWNER TO paybridge_admin;
CREATE INDEX ix_claims_case ON diagnosis_claims(case_id, viewer_scope);

-- 4. Claim Evidence (Association Table 24/40)
CREATE TABLE claim_evidence (
    claim_id UUID NOT NULL REFERENCES diagnosis_claims(id) ON DELETE CASCADE,
    evidence_id UUID NOT NULL REFERENCES evidence_items(id) ON DELETE RESTRICT,
    PRIMARY KEY (claim_id, evidence_id)
);
ALTER TABLE claim_evidence OWNER TO paybridge_admin;

-- 5. Recovery Actions (Next Safe Action Engine Proposals Table 25/40)
CREATE TABLE recovery_actions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    case_id UUID NOT NULL REFERENCES payment_cases(id) ON DELETE CASCADE,
    investigation_id UUID REFERENCES investigations(id) ON DELETE SET NULL,
    action_type TEXT NOT NULL CHECK (action_type IN (
        'REQUEST_BUYER_CONSENT', 'SUBMIT_PAYER_NOTE', 'RETRY_CHECKOUT',
        'INITIATE_REFUND', 'CONTACT_SUPPORT', 'CLOSE_CASE', 'ESCALATE_HUMAN'
    )),
    proposed_by TEXT NOT NULL CHECK (proposed_by IN ('SYSTEM_DETERMINISTIC', 'AI_ORCHESTRATOR', 'HUMAN_OPERATOR')),
    action_payload JSONB NOT NULL DEFAULT '{}'::jsonb,
    status TEXT NOT NULL DEFAULT 'PROPOSED' CHECK (status IN (
        'PROPOSED', 'APPROVED', 'REJECTED', 'EXECUTING', 'EXECUTED', 'FAILED'
    )),
    created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp()
);
ALTER TABLE recovery_actions OWNER TO paybridge_admin;
CREATE INDEX ix_recovery_actions_case ON recovery_actions(case_id, status);

-- 6. Action Acknowledgements (Human-in-the-Loop Table 26/40)
CREATE TABLE action_acks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    action_id UUID NOT NULL REFERENCES recovery_actions(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    decision TEXT NOT NULL CHECK (decision IN ('APPROVED', 'REJECTED')),
    notes TEXT,
    acked_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp()
);
ALTER TABLE action_acks OWNER TO paybridge_admin;
CREATE INDEX ix_action_acks_action ON action_acks(action_id);

-- +goose Down
DROP TABLE IF EXISTS action_acks;
DROP TABLE IF EXISTS recovery_actions;
DROP TABLE IF EXISTS claim_evidence;
DROP TABLE IF EXISTS diagnosis_claims;
DROP TABLE IF EXISTS agent_runs;
DROP TABLE IF EXISTS investigations;
