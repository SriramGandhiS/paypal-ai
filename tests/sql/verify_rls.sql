-- PayBridge AI — PostgreSQL 16 Security & RLS Verification Suite
-- Run via: psql -U paybridge_admin -d paybridge -f tests/sql/verify_rls.sql

\set ON_ERROR_STOP on
\timing on

BEGIN;
SELECT '=== STARTING PAYBRIDGE AI S3 SECURITY VERIFICATION ===' AS test_banner;

-- 1. Table Count Verification (MUST BE EXACTLY 40)
DO $$
DECLARE
    v_count INT;
BEGIN
    SELECT count(*) INTO v_count
      FROM information_schema.tables
     WHERE table_schema = 'public'
       AND table_type = 'BASE TABLE';

    IF v_count != 40 THEN
        RAISE EXCEPTION 'TABLE COUNT MISMATCH: Expected 40 tables, found %', v_count;
    END IF;
    RAISE NOTICE 'SUCCESS: Exactly 40 application tables verified.';
END $$;

-- 2. Verify RLS Enabled on Tenant Tables
DO $$
DECLARE
    v_unprotected text[];
BEGIN
    SELECT array_agg(c.relname) INTO v_unprotected
      FROM pg_class c
      JOIN pg_namespace n ON n.oid = c.relnamespace
     WHERE n.nspname = 'public'
       AND c.relkind = 'r'
       AND c.relname IN ('payment_cases', 'case_snapshots', 'evidence_items', 'audit_events')
       AND (NOT c.relrowsecurity OR NOT c.relforcerowsecurity);

    IF v_unprotected IS NOT NULL THEN
        RAISE EXCEPTION 'SECURITY FAILURE: Tables without FORCE RLS: %', v_unprotected;
    END IF;
    RAISE NOTICE 'SUCCESS: FORCE ROW LEVEL SECURITY verified on core tables.';
END $$;

-- 3. Verify paybridge_app Direct SELECT Revocation
DO $$
BEGIN
    IF has_table_privilege('paybridge_app', 'public.payment_cases', 'SELECT') THEN
        RAISE EXCEPTION 'SECURITY FAILURE: paybridge_app has direct SELECT on payment_cases';
    END IF;
    IF has_table_privilege('paybridge_app', 'public.tenants', 'SELECT') THEN
        RAISE EXCEPTION 'SECURITY FAILURE: paybridge_app has direct SELECT on tenants';
    END IF;
    RAISE NOTICE 'SUCCESS: Direct SELECT revoked from paybridge_app.';
END $$;

-- 4. Verify Search Path on SECURITY DEFINER Functions
DO $$
DECLARE
    v_insecure text[];
BEGIN
    SELECT array_agg(p.proname) INTO v_insecure
      FROM pg_proc p
      JOIN pg_namespace n ON n.oid = p.pronamespace
     WHERE n.nspname = 'public'
       AND p.prosecdef = TRUE
       AND (p.proconfig IS NULL OR NOT ('search_path=pg_catalog, public, pg_temp' = ANY(p.proconfig)));

    IF v_insecure IS NOT NULL THEN
        RAISE EXCEPTION 'SECURITY FAILURE: SECURITY DEFINER functions with unhardened search_path: %', v_insecure;
    END IF;
    RAISE NOTICE 'SUCCESS: All SECURITY DEFINER functions have hardened search_path.';
END $$;

ROLLBACK;
SELECT '=== ALL S3 DATABASE SECURITY VERIFICATIONS PASSED ===' AS final_verdict;
