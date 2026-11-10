# PayBridge AI — Database, Schema & Row-Level Security Audit Report

**Report ID:** DB-QA-2026-11-10  
**Target:** PostgreSQL 16 Schema, Migrations 0001–0010 & Row-Level Security  
**Lead:** Sriram (Google Antigravity)  
**Status:** **DDL & SYNTAX VERIFIED; LIVE RUNTIME BLOCKED BY HOST OS UAC**

---

## 1. Executive Summary

The database layer of PayBridge AI was subjected to comprehensive static, structural, and syntactic analysis:
- **Migration Suite:** Exactly 10 Goose migration files (`0001_initial_infra.sql` to `0010_hardening_rls_and_procs.sql`) are embedded and validated via Go tests.
- **Table Inventory:** Exactly 40 canonical tables are created with proper primary keys, foreign keys, indexes, and stored generated columns.
- **Critical Bug Fixed:** A corrupted unicode character (`\u4fe1` / `信`) in migration 0010's dynamic policy definition was discovered and corrected to `''''`.
- **Connection-Pooling Advisory Lock Fixed:** Identified that session-level PostgreSQL advisory locks borrowed from `*sql.DB` return unreleased locks to the pool; refactored `PostgresLocker` to pin a dedicated `*sql.Conn` exclusively until release.

---

## 2. Canonical 40-Table Inventory Verification

Every table was extracted and validated by `TestCanonical40TablesInventory`:

| Group | Tables Declared | Primary Keys & Enforcements |
|---|---|---|
| **0001: Extensions & Base** | `schema_migrations` | Goose version tracking |
| **0002: Identity & Tenancy** | `tenants`, `users`, `tenant_memberships`, `sessions`, `auth_signing_keys`, `consumed_auth_nonces` | UUID PKs, foreign keys, HMAC assertion nonces with 60s TTL |
| **0003: Cases & Idempotency** | `payment_cases`, `case_participants`, `payment_attempts`, `provider_resources`, `idempotency_keys`, `case_snapshots` | 3D state vector checks, `uq_one_active_attempt` blind-retry guard, composite idempotency PK |
| **0004: Evidence & Events** | `evidence_items`, `evidence_payloads`, `evidence_edges`, `payment_events`, `webhook_inbox` | Tier 1–6 validation checks, immutable SHA-256 digests, payload separation |
| **0005: Outbox & Relays** | `outbox`, `dead_letters` | `leased_by`, `leased_until` timestamps for `FOR UPDATE SKIP LOCKED` workers |
| **0006: AI & Investigations** | `investigations`, `agent_runs`, `diagnosis_claims`, `claim_evidence`, `recovery_actions`, `action_acks` | Immutable claim lineage, evidence edges |
| **0007: Passports & Consent** | `passport_versions`, `invitations`, `consent_grants`, `passport_shares` | Cryptographic passport hashes, privacy boundary allowlists |
| **0008: Audit & Retention** | `audit_chain_heads`, `audit_events`, `audit_archive_anchors`, `retention_jobs`, `notifications` | SHA-256 hash chains with transactionally locked head pointers |
| **0009: Connectors & Ops** | `route_registry`, `provider_status_samples`, `provider_connections`, `connector_credential_access_logs`, `simulation_runs` | AES-256-GCM encrypted provider credentials, route latency tracking |
| **0010: Hardening & RLS** | *(RLS policies and stored procedures)* | `ENABLE ROW LEVEL SECURITY; FORCE ROW LEVEL SECURITY;` across 31 tenant tables |

---

## 3. Row-Level Security (RLS) Analysis

In migration `0010_hardening_rls_and_procs.sql`:
1. **Dynamic Enforcement Loop:** 29 core tenant tables are processed via a dynamic PL/pgSQL block that runs:
   ```sql
   ALTER TABLE public.<tbl> ENABLE ROW LEVEL SECURITY;
   ALTER TABLE public.<tbl> FORCE ROW LEVEL SECURITY;
   CREATE POLICY <tbl>_tenant_isolation ON public.<tbl> FOR ALL TO PUBLIC
       USING (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid)
       WITH CHECK (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid);
   ```
2. **Linked Table RLS:** `claim_evidence` does not contain a direct `tenant_id` column; it is secured via a subquery joining `diagnosis_claims.tenant_id`.
3. **Role-Scoped Outbox Policy:** 
   - Application role (`paybridge_app`): restricted strictly to its own `tenant_id`.
   - Worker role (`paybridge_worker`): authorized across all pending outbox jobs for queue dispatching.

---

## 4. Defect Discovered and Repaired: Corrupted Character in 0010

During adversarial scanning of migration 0010:
- **Defect:** Lines 31 and 32 contained a corrupted unicode character:
  ```sql
  -- Before (BROKEN):
  'USING (tenant_id = NULLIF(current_setting(''app.tenant_id'', true), '''信::uuid) '
  ```
- **Root Cause:** IME / UTF-8 conversion corruption during SQL quote escaping.
- **Fix:** Corrected to valid SQL quote escaping:
  ```sql
  -- After (FIXED):
  'USING (tenant_id = NULLIF(current_setting(''app.tenant_id'', true), '''')::uuid) '
  ```
- **Evidence:** Regression test `TestRLSAndForceRLSEnforcement` in `migrations_test.go` now explicitly tests that `\u4fe1` does not exist anywhere in the migration file.

---

## 5. Connection-Pooling Advisory Lock Remediation

- **Issue:** PostgreSQL session advisory locks (`pg_try_advisory_lock`) are bound to the specific physical TCP connection. Using `*sql.DB.QueryRowContext` borrows connection A to acquire the lock and returns it to the pool, while release borrows connection B!
- **Remediation:** [`services/api/internal/checkout/service.go`](file:///D:/Projects/PayBridgeAI/services/api/internal/checkout/service.go#L180-L225) was updated:
  ```go
  conn, err := p.db.Conn(ctx) // Checkout dedicated connection
  ...
  err = conn.QueryRowContext(ctx, "SELECT pg_try_advisory_lock($1)", lockID).Scan(&acquired)
  ...
  release := func() {
      _ = conn.QueryRowContext(relCtx, "SELECT pg_advisory_unlock($1)", lockID).Scan(&unlocked)
      _ = conn.Close() // Safely return unlocked connection to pool
  }
  ```

---

## 6. Runtime Blocker Status

- **Status:** **BLOCKED BY HOST ENVIRONMENT**
- **Root Cause:** PostgreSQL 16 engine is not installed on the local Windows machine. Attempting to install via `winget` requires Administrator UAC elevation (`IsInRole(Administrator) == False`). Docker daemon is also not installed.
- **Reproduction Procedure:** Once an administrator starts PostgreSQL on port 5432, run:
  ```powershell
  cd D:\Projects\PayBridgeAI\services\api
  go run cmd/migrate/main.go -action=up
  psql -U paybridge_admin -d paybridge_dev -f ../../tests/sql/verify_rls.sql
  ```
