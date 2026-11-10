# PayBridge AI — S3 to S4 Engineering Handoff Guide

## 1. Milestone Transition Overview
* **Milestone S3 (Completed):** PostgreSQL 16 Database Architecture, 40-Table Inventory, 10 Goose Migration Groups, Hardened RLS, Stored Procedures, and Go Database Connection Management.
* **Milestone S4 (Current/Next):** Core REST/gRPC API Implementation, OpenAPI Codegen, Session Authentication, Case Lifecycle Handlers, and Webhook Receiver Pipeline.

---

## 2. Artifacts Produced in S3

1. **Migration Files:** Located in `services/api/db/migrations/`:
   - `0001_initial_infra.sql`
   - `0002_identity_and_auth.sql`
   - `0003_cases_and_idempotency.sql`
   - `0004_evidence_and_events.sql`
   - `0005_outbox_and_dead_letters.sql`
   - `0006_investigations_and_actions.sql`
   - `0007_passports_and_sharing.sql`
   - `0008_audit_and_retention.sql`
   - `0009_connectors_and_routing.sql`
   - `0010_hardening_rls_and_procs.sql`
   - `migrations.go` (embedded filesystem)
2. **Go Store Package:** Located in `services/api/internal/store/`:
   - Connection pool management (`pgx`/`database/sql`).
   - Transaction manager with rollback-on-error (`WithTx`).
   - Transactional tenant context injector (`SetTenantContext`).
3. **HTTP Server & Error Framework:** Located in `services/api/internal/httpx/`:
   - Panic recovery and request correlation ID middleware.
   - RFC 7807 problem details error envelopes.
   - `/healthz` and `/readyz` health endpoints.
4. **Standalone Migration CLI:** Located in `services/api/cmd/migrate/main.go`.
5. **OpenAPI 3.1 Contract:** Located in `packages/contracts/openapi.yaml`.

---

## 3. How S4 Implements Against S3 Contracts

1. **API Router Wiring:** `services/api/cmd/api/main.go` uses `httpx.Router` to register routes from the route inventory (`04_BACKEND_ENGINEERING_SPECIFICATION.md` §14).
2. **Procedural Calls:** Instead of issuing raw `INSERT/UPDATE/DELETE` queries on core financial tables, Go handlers invoke the verified stored procedures (`fn_cases_create_case`, `fn_reducer_apply_transition`, etc.) using assertion tokens signed by the Auth service.
3. **Outbox Polling:** The background worker daemon (`services/api/cmd/worker`) polls `outbox` via `ClaimOutboxBatch` using `FOR UPDATE SKIP LOCKED`, preventing race conditions.
4. **PayPal Webhook Ingest:** `POST /api/v1/webhooks/paypal` writes raw payloads into `webhook_inbox` under the dedicated `paybridge_webhook` role, preventing direct business data tampering.
