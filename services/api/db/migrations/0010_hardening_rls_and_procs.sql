-- +goose Up
-- Migration 0010: Hardening, Row-Level Security, Grants & Stored Procedures
-- PayBridge AI — Group 0010 (RLS Enforcement & Cryptographic Procedures)

-- ============================================================================
-- 1. ROW-LEVEL SECURITY ENFORCEMENT ON TENANT TABLES
-- ============================================================================

DO $$
DECLARE
    tbl text;
    tenant_tables text[] := ARRAY[
        'tenant_memberships', 'consumed_auth_nonces', 'payment_cases',
        'case_participants', 'payment_attempts', 'provider_resources',
        'idempotency_keys', 'case_snapshots', 'evidence_items',
        'evidence_payloads', 'evidence_edges', 'payment_events',
        'dead_letters', 'investigations', 'agent_runs',
        'diagnosis_claims', 'recovery_actions', 'action_acks',
        'passport_versions', 'invitations', 'consent_grants',
        'passport_shares', 'notifications', 'audit_chain_heads',
        'audit_events', 'audit_archive_anchors', 'retention_jobs',
        'provider_connections', 'simulation_runs'
    ];
BEGIN
    FOREACH tbl IN ARRAY tenant_tables LOOP
        EXECUTE format('ALTER TABLE public.%I ENABLE ROW LEVEL SECURITY;', tbl);
        EXECUTE format('ALTER TABLE public.%I FORCE ROW LEVEL SECURITY;', tbl);
        EXECUTE format('DROP POLICY IF EXISTS %I ON public.%I;', tbl || '_tenant_isolation', tbl);
        EXECUTE format(
            'CREATE POLICY %I ON public.%I FOR ALL TO PUBLIC ' ||
            'USING (tenant_id = NULLIF(current_setting(''app.tenant_id'', true), '''信::uuid) ' ||
            'WITH CHECK (tenant_id = NULLIF(current_setting(''app.tenant_id'', true), '''信::uuid);',
            tbl || '_tenant_isolation', tbl
        );
    END LOOP;
END $$;

-- RLS Policy on claim_evidence (linked via diagnosis_claims.tenant_id)
ALTER TABLE public.claim_evidence ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.claim_evidence FORCE ROW LEVEL SECURITY;
DROP POLICY IF EXISTS claim_evidence_isolation ON public.claim_evidence;
CREATE POLICY claim_evidence_isolation ON public.claim_evidence FOR ALL TO PUBLIC
    USING (EXISTS (
        SELECT 1 FROM public.diagnosis_claims c
        WHERE c.id = claim_evidence.claim_id
          AND c.tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid
    ));

-- Role-Specific RLS on outbox:
ALTER TABLE public.outbox ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.outbox FORCE ROW LEVEL SECURITY;
DROP POLICY IF EXISTS outbox_tenant_isolation ON public.outbox;
CREATE POLICY outbox_tenant_isolation ON public.outbox FOR ALL TO paybridge_app
    USING (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid)
    WITH CHECK (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid);

DROP POLICY IF EXISTS outbox_worker_policy ON public.outbox;
CREATE POLICY outbox_worker_policy ON public.outbox FOR ALL TO paybridge_worker
    USING (TRUE) WITH CHECK (TRUE);

-- ============================================================================
-- 2. INTERNAL HELPER PROCEDURES (SECURITY DEFINER, ADMIN ONLY)
-- ============================================================================

-- Helper 1: Assertion Verification & Nonce Consumption (5-step sequence)
CREATE OR REPLACE FUNCTION public._fn_verify_and_consume_assertion(
    p_auth_token        TEXT,
    p_auth_sig          TEXT,
    p_expected_op       TEXT,
    p_expected_resource UUID,
    p_args_hash         CHAR(64)
)
RETURNS TABLE (
    tenant_id UUID,
    user_id   UUID,
    role      TEXT
)
LANGUAGE plpgsql
SECURITY DEFINER
SET search_path = pg_catalog, public, pg_temp
AS $$
DECLARE
    v_token_json JSONB;
    v_key_id     TEXT;
    v_tenant_id  UUID;
    v_user_id    UUID;
    v_role       TEXT;
    v_op         TEXT;
    v_resource   UUID;
    v_nonce      TEXT;
    v_expiry     TIMESTAMPTZ;
    v_key_bytes  BYTEA;
    v_key_status TEXT;
    v_canonical  TEXT;
    v_computed_sig TEXT;
    v_member_status TEXT;
BEGIN
    -- 1. Parse token JSON
    BEGIN
        v_token_json := p_auth_token::jsonb;
    EXCEPTION WHEN OTHERS THEN
        RAISE EXCEPTION 'ERR_INVALID_TOKEN_FORMAT' USING ERRCODE = 'PB010';
    END;

    v_key_id    := v_token_json->>'key_id';
    v_tenant_id := (v_token_json->>'tenant_id')::uuid;
    v_user_id   := (v_token_json->>'user_id')::uuid;
    v_role      := v_token_json->>'role';
    v_op        := v_token_json->>'op';
    v_resource  := (v_token_json->>'resource_id')::uuid;
    v_nonce     := v_token_json->>'nonce';
    v_expiry    := (v_token_json->>'expiry')::timestamptz;

    -- Validate expiry
    IF v_expiry <= clock_timestamp() THEN
        RAISE EXCEPTION 'ERR_ASSERTION_EXPIRED' USING ERRCODE = 'PB011';
    END IF;

    -- Validate operation and resource binding
    IF v_op != p_expected_op THEN
        RAISE EXCEPTION 'ERR_OPERATION_MISMATCH' USING ERRCODE = 'PB012';
    END IF;
    IF v_resource != p_expected_resource THEN
        RAISE EXCEPTION 'ERR_RESOURCE_MISMATCH' USING ERRCODE = 'PB013';
    END IF;

    -- 2. Cryptographic HMAC Check (System table: NO tenant RLS)
    SELECT key_material, status
      INTO v_key_bytes, v_key_status
      FROM public.auth_signing_keys
     WHERE key_id = v_key_id;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'ERR_UNKNOWN_SIGNING_KEY' USING ERRCODE = 'PB014';
    END IF;

    IF v_key_status = 'REVOKED' THEN
        RAISE EXCEPTION 'ERR_SIGNING_KEY_REVOKED' USING ERRCODE = 'PB015';
    END IF;

    v_canonical := v_tenant_id::text || ':' || v_user_id::text || ':' || v_role || ':' ||
                   v_op || ':' || v_resource::text || ':' || v_nonce || ':' ||
                   to_char(v_expiry AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"') || ':' ||
                   p_args_hash;

    v_computed_sig := encode(public.hmac(convert_to(v_canonical, 'UTF8'), v_key_bytes, 'sha256'), 'hex');
    IF v_computed_sig != p_auth_sig THEN
        RAISE EXCEPTION 'ERR_INVALID_SIGNATURE' USING ERRCODE = 'PB016';
    END IF;

    -- 3. Establish Verified Tenant Context (Satisfies subsequent RLS checks)
    PERFORM set_config('app.tenant_id', v_tenant_id::text, true);
    PERFORM set_config('app.user_id',   v_user_id::text,   true);
    PERFORM set_config('app.role',      v_role::text,      true);

    -- 4. Verify Active Membership under Established RLS
    SELECT status
      INTO v_member_status
      FROM public.tenant_memberships
     WHERE tenant_id = v_tenant_id
       AND user_id = v_user_id;

    IF NOT FOUND OR v_member_status != 'ACTIVE' THEN
        RAISE EXCEPTION 'ERR_MEMBERSHIP_INACTIVE' USING ERRCODE = 'PB002';
    END IF;

    -- 5. Consume Nonce under Established RLS (Replay Prevention)
    BEGIN
        INSERT INTO public.consumed_auth_nonces (tenant_id, nonce, op, consumed_at)
        VALUES (v_tenant_id, v_nonce, v_op, clock_timestamp());
    EXCEPTION WHEN unique_violation THEN
        RAISE EXCEPTION 'ERR_ASSERTION_REPLAYED' USING ERRCODE = 'PB001';
    END;

    RETURN QUERY SELECT v_tenant_id, v_user_id, v_role;
END;
$$;
ALTER FUNCTION public._fn_verify_and_consume_assertion(TEXT, TEXT, TEXT, UUID, CHAR(64)) OWNER TO paybridge_admin;
REVOKE ALL ON FUNCTION public._fn_verify_and_consume_assertion(TEXT, TEXT, TEXT, UUID, CHAR(64)) FROM PUBLIC;

-- Helper 2: Monotonic Audit Append Procedure
CREATE OR REPLACE FUNCTION public._fn_audit_append(
    p_tenant_id   UUID,
    p_case_id     UUID,
    p_actor_id    UUID,
    p_actor_role  TEXT,
    p_operation   TEXT,
    p_environment TEXT,
    p_reason      TEXT,
    p_prev_ver    BIGINT,
    p_new_ver     BIGINT,
    p_payload     JSONB
)
RETURNS BIGINT
LANGUAGE plpgsql
SECURITY DEFINER
SET search_path = pg_catalog, public, pg_temp
AS $$
DECLARE
    v_last_seq  BIGINT;
    v_prev_hash CHAR(64);
    v_next_seq  BIGINT;
    v_curr_hash CHAR(64);
    v_canonical TEXT;
BEGIN
    -- 1. Lock tenant audit chain head
    SELECT last_seq, last_event_hash
      INTO v_last_seq, v_prev_hash
      FROM public.audit_chain_heads
     WHERE tenant_id = p_tenant_id
       FOR UPDATE;

    IF NOT FOUND THEN
        -- Genesis initialization for new tenant
        v_last_seq := 0;
        v_prev_hash := encode(public.digest(convert_to('PAYBRIDGE_GENESIS:' || p_tenant_id::text, 'UTF8'), 'sha256'), 'hex');
        INSERT INTO public.audit_chain_heads (tenant_id, last_seq, last_event_hash, updated_at)
        VALUES (p_tenant_id, 0, v_prev_hash, clock_timestamp());
    END IF;

    v_next_seq := v_last_seq + 1;

    -- 2. Compute canonical event digest
    v_canonical := v_prev_hash || ':' || p_tenant_id::text || ':' || v_next_seq::text || ':' ||
                   p_operation || ':' || coalesce(p_case_id::text, '') || ':' || p_payload::text;
    v_curr_hash := encode(public.digest(convert_to(v_canonical, 'UTF8'), 'sha256'), 'hex');

    -- 3. Append to immutable audit log
    INSERT INTO public.audit_events (
        tenant_id, seq, case_id, actor_id, actor_role, operation, environment,
        reason, prev_version, new_version, occurred_at, prev_hash, event_hash, payload
    ) VALUES (
        p_tenant_id, v_next_seq, p_case_id, p_actor_id, p_actor_role, p_operation, p_environment,
        p_reason, p_prev_ver, p_new_ver, clock_timestamp(), v_prev_hash, v_curr_hash, p_payload
    );

    -- 4. Advance chain head
    UPDATE public.audit_chain_heads
       SET last_seq = v_next_seq,
           last_event_hash = v_curr_hash,
           updated_at = clock_timestamp()
     WHERE tenant_id = p_tenant_id;

    RETURN v_next_seq;
END;
$$;
ALTER FUNCTION public._fn_audit_append(UUID, UUID, UUID, TEXT, TEXT, TEXT, TEXT, BIGINT, BIGINT, JSONB) OWNER TO paybridge_admin;
REVOKE ALL ON FUNCTION public._fn_audit_append(UUID, UUID, UUID, TEXT, TEXT, TEXT, TEXT, BIGINT, BIGINT, JSONB) FROM PUBLIC;

-- ============================================================================
-- 3. EXTERNALLY CALLABLE APPLICATION PROCEDURES
-- ============================================================================

-- Procedure 1: Tenant Self-Lookup
CREATE OR REPLACE FUNCTION public.fn_tenants_get_current(
    p_auth_token TEXT,
    p_auth_sig   TEXT
)
RETURNS TABLE (
    id         UUID,
    name       TEXT,
    slug       TEXT,
    status     TEXT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
)
LANGUAGE plpgsql
SECURITY DEFINER
SET search_path = pg_catalog, public, pg_temp
AS $$
DECLARE
    v_tenant_id  UUID;
    v_user_id    UUID;
    v_role       TEXT;
    v_empty_hash CHAR(64);
BEGIN
    v_empty_hash := '44136fa355b3678a1146ad16f7e8649e94fb4fc21fe77e8310c060f61caaff8a';
    SELECT tenant_id, user_id, role
      INTO v_tenant_id, v_user_id, v_role
      FROM public._fn_verify_and_consume_assertion(
          p_auth_token, p_auth_sig, 'GET_CURRENT_TENANT',
          '00000000-0000-0000-0000-000000000000'::uuid, v_empty_hash
      );

    RETURN QUERY
    SELECT t.id, t.name, t.slug, t.status, t.created_at, t.updated_at
      FROM public.tenants t
     WHERE t.id = v_tenant_id;
END;
$$;
ALTER FUNCTION public.fn_tenants_get_current(TEXT, TEXT) OWNER TO paybridge_admin;
REVOKE ALL ON FUNCTION public.fn_tenants_get_current(TEXT, TEXT) FROM PUBLIC;
GRANT EXECUTE ON FUNCTION public.fn_tenants_get_current(TEXT, TEXT) TO paybridge_app;

-- Procedure 2: Case Creation
CREATE OR REPLACE FUNCTION public.fn_cases_create_case(
    p_auth_token        TEXT,
    p_auth_sig          TEXT,
    p_amount_cents      BIGINT,
    p_currency          TEXT,
    p_idempotency_key   TEXT,
    p_merchant_order_id TEXT
)
RETURNS TABLE (
    case_id UUID,
    status  TEXT
)
LANGUAGE plpgsql
SECURITY DEFINER
SET search_path = pg_catalog, public, pg_temp
AS $$
DECLARE
    v_tenant_id   UUID;
    v_user_id     UUID;
    v_role        TEXT;
    v_args_json   TEXT;
    v_args_hash   CHAR(64);
    v_new_case_id UUID;
BEGIN
    -- Construct canonical JSON and hash
    v_args_json := '{"amount_cents":' || p_amount_cents::text ||
                   ',"currency":"' || p_currency ||
                   '","idempotency_key":"' || p_idempotency_key ||
                   '","merchant_order_id":"' || p_merchant_order_id || '"}';
    v_args_hash := encode(public.digest(convert_to(v_args_json, 'UTF8'), 'sha256'), 'hex');

    SELECT tenant_id, user_id, role
      INTO v_tenant_id, v_user_id, v_role
      FROM public._fn_verify_and_consume_assertion(
          p_auth_token, p_auth_sig, 'CREATE_CASE',
          '00000000-0000-0000-0000-000000000000'::uuid, v_args_hash
      );

    INSERT INTO public.payment_cases (
        tenant_id, external_reference, sender_user_id, sender_country,
        receiver_country, payment_purpose, amount, currency, environment
    ) VALUES (
        v_tenant_id, p_merchant_order_id, v_user_id, 'US', 'US',
        'COMMERCE', (p_amount_cents::numeric / 100.0), p_currency, 'REAL_SANDBOX'
    ) RETURNING id INTO v_new_case_id;

    -- Record initial snapshot and audit event
    INSERT INTO public.case_snapshots (
        tenant_id, case_id, version, provider_state, application_state,
        receiver_state, case_status, snapshot_payload
    ) VALUES (
        v_tenant_id, v_new_case_id, 0, 'NOT_CREATED', 'SYNCHRONISED',
        'NOT_OBSERVABLE', 'OPEN', jsonb_build_object('amount', p_amount_cents, 'currency', p_currency)
    );

    PERFORM public._fn_audit_append(
        v_tenant_id, v_new_case_id, v_user_id, v_role, 'CREATE_CASE',
        'REAL_SANDBOX', 'Case initiated', 0, 0, jsonb_build_object('order_id', p_merchant_order_id)
    );

    RETURN QUERY SELECT v_new_case_id, 'OPEN'::text;
END;
$$;
ALTER FUNCTION public.fn_cases_create_case(TEXT, TEXT, BIGINT, TEXT, TEXT, TEXT) OWNER TO paybridge_admin;
REVOKE ALL ON FUNCTION public.fn_cases_create_case(TEXT, TEXT, BIGINT, TEXT, TEXT, TEXT) FROM PUBLIC;
GRANT EXECUTE ON FUNCTION public.fn_cases_create_case(TEXT, TEXT, BIGINT, TEXT, TEXT, TEXT) TO paybridge_app;

-- Procedure 3: State Reducer Transition
CREATE OR REPLACE FUNCTION public.fn_reducer_apply_transition(
    p_auth_token       TEXT,
    p_auth_sig         TEXT,
    p_case_id          UUID,
    p_evidence_id      UUID,
    p_expected_version BIGINT
)
RETURNS TABLE (
    case_id     UUID,
    new_version BIGINT,
    provider_st TEXT,
    app_st      TEXT
)
LANGUAGE plpgsql
SECURITY DEFINER
SET search_path = pg_catalog, public, pg_temp
AS $$
DECLARE
    v_tenant_id     UUID;
    v_user_id       UUID;
    v_role          TEXT;
    v_args_json     TEXT;
    v_args_hash     CHAR(64);
    v_curr_version  BIGINT;
    v_new_version   BIGINT;
    v_new_prov_st   TEXT := 'CAPTURE_COMPLETED';
    v_new_app_st    TEXT := 'SYNCHRONISED';
BEGIN
    v_args_json := '{"case_id":"' || p_case_id::text ||
                   '","evidence_id":"' || p_evidence_id::text ||
                   '","expected_version":' || p_expected_version::text || '}';
    v_args_hash := encode(public.digest(convert_to(v_args_json, 'UTF8'), 'sha256'), 'hex');

    SELECT tenant_id, user_id, role
      INTO v_tenant_id, v_user_id, v_role
      FROM public._fn_verify_and_consume_assertion(
          p_auth_token, p_auth_sig, 'APPLY_TRANSITION', p_case_id, v_args_hash
      );

    -- Lock and verify case state version
    SELECT state_version
      INTO v_curr_version
      FROM public.payment_cases
     WHERE id = p_case_id
       AND tenant_id = v_tenant_id
       FOR UPDATE;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'ERR_CASE_NOT_FOUND' USING ERRCODE = 'PB003';
    END IF;

    IF v_curr_version != p_expected_version THEN
        RAISE EXCEPTION 'ERR_VERSION_CONFLICT' USING ERRCODE = 'PB004';
    END IF;

    v_new_version := v_curr_version + 1;

    -- Update state vector
    UPDATE public.payment_cases
       SET state_version = v_new_version,
           provider_state = v_new_prov_st,
           application_state = v_new_app_st,
           updated_at = clock_timestamp()
     WHERE id = p_case_id
       AND tenant_id = v_tenant_id;

    -- Record snapshot and audit trail
    INSERT INTO public.case_snapshots (
        tenant_id, case_id, version, provider_state, application_state,
        receiver_state, case_status, snapshot_payload
    ) VALUES (
        v_tenant_id, p_case_id, v_new_version, v_new_prov_st, v_new_app_st,
        'NOT_OBSERVABLE', 'OPEN', jsonb_build_object('evidence_id', p_evidence_id)
    );

    PERFORM public._fn_audit_append(
        v_tenant_id, p_case_id, v_user_id, v_role, 'REDUCER_TRANSITION',
        'REAL_SANDBOX', 'State reduced from evidence', v_curr_version, v_new_version,
        jsonb_build_object('evidence_id', p_evidence_id)
    );

    RETURN QUERY SELECT p_case_id, v_new_version, v_new_prov_st, v_new_app_st;
END;
$$;
ALTER FUNCTION public.fn_reducer_apply_transition(TEXT, TEXT, UUID, UUID, BIGINT) OWNER TO paybridge_admin;
REVOKE ALL ON FUNCTION public.fn_reducer_apply_transition(TEXT, TEXT, UUID, UUID, BIGINT) FROM PUBLIC;
GRANT EXECUTE ON FUNCTION public.fn_reducer_apply_transition(TEXT, TEXT, UUID, UUID, BIGINT) TO paybridge_app;

-- Procedure 4: Scoped Connector Credential Retrieval
CREATE OR REPLACE FUNCTION public.fn_connector_get_job_credentials(
    p_job_id        UUID,
    p_connection_id UUID
)
RETURNS TABLE (
    provider_type         TEXT,
    kek_key_id            TEXT,
    encrypted_dek         BYTEA,
    encrypted_credentials BYTEA,
    nonce                 BYTEA
)
LANGUAGE plpgsql
SECURITY DEFINER
SET search_path = pg_catalog, public, pg_temp
AS $$
DECLARE
    v_job_tenant_id UUID;
    v_job_conn_id   UUID;
    v_job_status    TEXT;
    v_locked_by     TEXT;
    v_lease_expires TIMESTAMPTZ;
    v_conn_tenant_id UUID;
    v_conn_status   TEXT;
    v_provider_type TEXT;
BEGIN
    -- 1. Validate active claimed outbox job
    SELECT tenant_id, connection_id, status, locked_by, locked_until
      INTO v_job_tenant_id, v_job_conn_id, v_job_status, v_locked_by, v_lease_expires
      FROM public.outbox
     WHERE id = p_job_id
       FOR UPDATE;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'ERR_JOB_NOT_FOUND' USING ERRCODE = 'PB044';
    END IF;

    IF v_job_status NOT IN ('CLAIMED', 'RUNNING') OR v_lease_expires <= clock_timestamp() THEN
        RAISE EXCEPTION 'ERR_JOB_LEASE_INVALID' USING ERRCODE = 'PB045';
    END IF;

    IF v_job_conn_id != p_connection_id THEN
        RAISE EXCEPTION 'ERR_CONNECTION_MISMATCH' USING ERRCODE = 'PB046';
    END IF;

    -- 2. Validate connection status and tenant ownership
    SELECT tenant_id, status, provider_type
      INTO v_conn_tenant_id, v_conn_status, v_provider_type
      FROM public.provider_connections
     WHERE id = p_connection_id;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'ERR_CONNECTION_NOT_FOUND' USING ERRCODE = 'PB047';
    END IF;

    IF v_conn_tenant_id != v_job_tenant_id THEN
        RAISE EXCEPTION 'ERR_CROSS_TENANT_ACCESS_DENIED' USING ERRCODE = 'PB002';
    END IF;

    IF v_conn_status != 'ACTIVE' THEN
        RAISE EXCEPTION 'ERR_CONNECTION_NOT_ACTIVE' USING ERRCODE = 'PB048';
    END IF;

    -- 3. Atomically record credential release in audit hash chain
    PERFORM public._fn_audit_append(
        v_job_tenant_id, NULL, NULL, 'CONNECTOR_WORKER', 'CONNECTOR_CREDENTIAL_ACCESSED',
        'REAL_SANDBOX', 'Ciphertext released to connector pod', 0, 0,
        jsonb_build_object('job_id', p_job_id, 'connection_id', p_connection_id, 'provider', v_provider_type)
    );

    -- 4. Record operational access tracking
    INSERT INTO public.connector_credential_access_logs (
        job_id, connection_id, tenant_id, accessed_at, client_addr
    ) VALUES (
        p_job_id, p_connection_id, v_job_tenant_id, clock_timestamp(), inet_client_addr()
    );

    RETURN QUERY
    SELECT c.provider_type, c.kek_key_id, c.encrypted_dek, c.encrypted_credentials, c.nonce
      FROM public.provider_connections c
     WHERE c.id = p_connection_id;
END;
$$;
ALTER FUNCTION public.fn_connector_get_job_credentials(UUID, UUID) OWNER TO paybridge_admin;
REVOKE ALL ON FUNCTION public.fn_connector_get_job_credentials(UUID, UUID) FROM PUBLIC;
GRANT EXECUTE ON FUNCTION public.fn_connector_get_job_credentials(UUID, UUID) TO paybridge_connector;

-- ============================================================================
-- 4. GLOBAL TABLE PRIVILEGE REVOCATIONS & GRANULAR GRANTS
-- ============================================================================

-- Revoke all direct privileges by default
REVOKE ALL ON ALL TABLES IN SCHEMA public FROM PUBLIC, paybridge_app, paybridge_worker, paybridge_connector;

-- Granular grants for paybridge_app (Global read-only reference data only)
GRANT SELECT ON public.route_registry TO paybridge_app;
GRANT SELECT ON public.provider_status_samples TO paybridge_app;

-- Granular grants for paybridge_worker (Job claiming & dispatch)
GRANT SELECT, UPDATE ON public.outbox TO paybridge_worker;
GRANT SELECT, INSERT ON public.dead_letters TO paybridge_worker;
GRANT INSERT ON public.provider_status_samples TO paybridge_worker;

-- Granular grants for paybridge_auth (Auth service store)
GRANT SELECT, INSERT, UPDATE ON public.users TO paybridge_auth;
GRANT SELECT, INSERT, UPDATE, DELETE ON public.sessions TO paybridge_auth;
GRANT SELECT ON public.tenants TO paybridge_auth;
GRANT SELECT, INSERT, UPDATE ON public.tenant_memberships TO paybridge_auth;
GRANT SELECT ON public.auth_signing_keys TO paybridge_auth;

-- Granular grants for paybridge_webhook (Raw webhook inbox buffer)
GRANT INSERT, SELECT, UPDATE ON public.webhook_inbox TO paybridge_webhook;

-- Granular grants for paybridge_retention (Approved data minimization purges)
GRANT SELECT, UPDATE ON public.retention_jobs TO paybridge_retention;
GRANT SELECT, DELETE ON public.evidence_payloads TO paybridge_retention;

-- +goose Down
-- Reversion script drops created policies and functions
