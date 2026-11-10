# PayBridge AI — P0–P8 Handover & Team Integration Guide

**Author:** Sriram (Google Antigravity) — Foundations, Database, API Core & Base Shell  
**Recipients:** 
- **Vijay Kasthuri** (Claude) — AI Orchestrator DAG (A01–A10), Safe Action Rules & 16-Screen Frontend
- **Diny Paul** (Claude & Astra) — Enterprise Infrastructure, Multi-Rail, Kafka & Regulatory Compliance  
**Date:** 10 November 2026

---

## 1. System Architecture State Handover

Sriram has delivered the foundational platform layer for PayBridge AI. The substrate is tested, hardened, and locked:

1. **Go Backend Monolith (`services/api`)**:
   - Standardized `go-chi` HTTP router with RFC 7807 problem details error handling.
   - 21 out of 21 packages pass unit tests (`go test -count=1 ./services/api/...`).
   - Distributed capture concurrency guard via `AcquireDistributedCaptureLock`.
   - Security-hardened middleware rejecting unauthorized development bypasses.
2. **Database Migrations (`services/api/db/migrations`)**:
   - 10 Goose migration files covering all 40 canonical tables.
   - `FORCE ROW LEVEL SECURITY` policies defined on all 30 tenant-scoped tables.
   - Hardened `SECURITY DEFINER` stored procedures for assertion verification, audit hashing, and connector credentials.
3. **OpenAPI 3.1 Contract (`packages/contracts/openapi.yaml`)**:
   - Complete, authoritative API definitions for all 48 routes.
4. **Base Frontend Application (`apps/web`)**:
   - React 18 + Vite 5 + TypeScript build compiling with 0 errors.
   - 6 functional presentation screens with dual-mode simulation fallbacks.

---

## 2. Integration Instructions for Vijay Kasthuri (Claude)

### A. AI Role Security Boundaries
* **Prime Invariant:** AI agents must **NEVER** initiate financial captures, orders, or refunds.
* **Implementation:** Verified in [`services/api/internal/policy/engine.go`](file:///D:/Projects/PayBridgeAI/services/api/internal/policy/engine.go).
* The roles `ai_agent` and `investigator_agent` are strictly restricted to:
  - `ActionReadCase` (`GET /api/v1/cases/{id}`)
  - `ActionReadEvidence` (`GET /api/v1/cases/{id}/evidence`)
  - `ActionTriggerInvestigate` (`POST /api/v1/cases/{id}/investigate`)
  - Any call to `ActionCreateOrder` or `ActionCapturePayment` by an AI role is rejected with HTTP 403 Forbidden (`policy.ErrForbidden`).

### B. Connecting the A01–A10 Multi-Agent DAG
* Ingest evidence by consuming the Tier 1–6 evidence store via `GET /api/v1/cases/{id}/evidence`.
* Submit diagnosis claims and suggested recovery actions to:
  - `POST /api/v1/cases/{id}/diagnose`
  - `POST /api/v1/cases/{id}/actions`
* All claims submitted by your agents will automatically link to the immutable evidence ledger via `evidence_edges`.

### C. Expanding Frontend Screens
* Sriram has built the core shell (`apps/web/src/App.tsx`, `Navbar.tsx`, and 6 primary screens).
* To add your remaining 10 screens (e.g. settings, detailed audit explorer, webhook manager):
  1. Add screen components into `apps/web/src/screens/`.
  2. Add navigation tabs in `apps/web/src/components/Navbar.tsx`.
  3. Reference typed client methods in `apps/web/src/api/client.ts`.

---

## 3. Integration Instructions for Diny Paul (Claude & Astra)

### A. High-Throughput Event Streaming (>10,000 TPS)
* The Outbox relay table is fully provisioned in `0005_outbox_and_dead_letters.sql`:
  - `outbox` table with `leased_by`, `leased_until`, and `FOR UPDATE SKIP LOCKED` claim logic.
  - Connect your Kafka producer / Debezium CDC listener to the `outbox` table or the Outbox dispatcher in `services/api/internal/outbox/`.

### B. Multi-Rail Adapters (Stripe, Wise)
* Implement multi-rail connectors implementing the standard interface in `packages/contracts/openapi.yaml`.
* Store provider transaction references in `provider_resources` (Table 11/40).

### C. Statutory Compliance (RBI Purpose Codes & e-FIRC)
* Use the existing Tier 1–6 evidence ingestion pipeline:
  - Ingest e-FIRC certificates as Tier 3 (`BANK_SWIFT_CONFIRMATION`) evidence items.
  - Set `access_class = 'COMPLIANCE'` to enforce tenant data isolation.

---

## 4. Master Health & Verification Commands
To verify the entire foundational stack before starting feature integration:
```bash
# 1. Verify Go Backend Tests (21/21 ok)
go test -count=1 ./services/api/...

# 2. Verify Frontend Build (0 errors)
cd apps/web && npm run build

# 3. Verify Golden Path Simulation
cd services/api && go run cmd/demo/main.go
```
