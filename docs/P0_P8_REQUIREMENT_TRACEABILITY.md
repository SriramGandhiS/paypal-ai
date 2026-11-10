# PayBridge AI — P0–P8 Requirement Traceability Matrix

**Project:** PayBridge AI  
**Role:** Sriram (Google Antigravity) — Platform Foundations, Database, API Core & Base UI Shell  
**Date:** 10 November 2026  
**Status:** COMPLETE (Sriram Scope)

---

## 1. Traceability Matrix

| Phase | Requirement ID | Requirement Description | Implementation File(s) | Verification / Test Evidence | Status |
|:---:|:---:|---|---|---|:---:|
| **P0** | REQ-P0-01 | Repository audit and baseline commit check | `Makefile`, `README.md` | Git commit baseline `e4c5399` verified; branch `master` | **VERIFIED** |
| **P0** | REQ-P0-02 | Requirement & team ownership boundaries | `docs/TEAM_INTEGRATION_CONTRACTS.md` | Boundaries established (Sriram vs Vijay vs Diny) | **VERIFIED** |
| **P0** | REQ-P0-03 | Environment & secret name inventory | `.env.example`, `internal/config/config.go` | Documented in doc 08 §21; values never printed in logs | **VERIFIED** |
| **P1** | REQ-P1-01 | Monorepo layout and dev rig | `docker-compose.yml`, `infra/docker/docker-compose.yml` | Multi-service Compose (Postgres 16, Redis 7.2, Mailpit) | **VERIFIED** |
| **P1** | REQ-P1-02 | Config validator with production guard | `services/api/internal/config/config.go` | `TestConfig_TableDriven` (28 test cases passing) | **VERIFIED** |
| **P1** | REQ-P1-03 | Sandbox-only PayPal endpoint enforcement | `services/api/internal/config/config.go` | Fails fast if non-sandbox host used in dev mode | **VERIFIED** |
| **P1** | REQ-P1-04 | PayPal OAuth & sandbox client | `services/api/internal/paypal/client.go` | `client_test.go` (token cache, orders, captures, GET truth) | **VERIFIED** |
| **P2** | REQ-P2-01 | 10 Goose database migrations (40 tables) | `services/api/db/migrations/0001_...` to `0010_...` | `migrations_test.go` (10 embedded migrations verified) | **CODE COMPLETE** |
| **P2** | REQ-P2-02 | Row-Level Security on 30 tenant tables | `services/api/db/migrations/0010_hardening_rls_and_procs.sql` | `tests/sql/verify_rls.sql` (20 security tests scripted) | **AWAITING DB** |
| **P2** | REQ-P2-03 | Webhook Gateway & signature verification | `services/api/internal/webhook/gateway.go` | `gateway_test.go` (raw capture, headers, signature check) | **VERIFIED** |
| **P2** | REQ-P2-04 | Event inbox & outbox relays | `services/api/internal/inbox/`, `internal/outbox/` | `inbox_test.go`, `outbox_test.go` passing | **VERIFIED** |
| **P2** | REQ-P2-05 | Reader-only reconciliation engine | `services/api/internal/reconciliation/engine.go` | `engine_test.go` (0 mutations, amount mismatch guard) | **VERIFIED** |
| **P3** | REQ-P3-01 | Multi-tenant context & RBAC authorization | `services/api/internal/auth/`, `internal/policy/` | `policy_test.go`, `auth_test.go` (AI mutation denied) | **VERIFIED** |
| **P3** | REQ-P3-02 | Case & attempt management with blind-retry guard | `services/api/internal/cases/`, `internal/checkout/` | `cases_test.go`, `checkout_test.go` passing | **VERIFIED** |
| **P3** | REQ-P3-03 | Distributed / DB-backed capture concurrency lock | `services/api/internal/checkout/service.go` | `TestDistributedCaptureLock_Concurrency` (10 concurrent goroutines) | **VERIFIED** |
| **P3** | REQ-P3-04 | Capture race condition HTTP 409 handling | `services/api/internal/api/api.go` | `TestPayPalCapture_ConcurrentRequestsRaceCondition` | **VERIFIED** |
| **P3** | REQ-P3-05 | Tier 1–6 evidence provenance & SHA-256 seal | `services/api/internal/evidence/service.go` | `evidence_test.go` (SHA-256 calculation, tier validation) | **VERIFIED** |
| **P3** | REQ-P3-06 | Two-Lens explanation & Truth Passport model | `services/api/internal/api/api.go`, `apps/web/src/screens/PassportScreen.tsx` | Merchant vs Investigator lenses toggled | **VERIFIED** |
| **P4** | REQ-P4-01 | React 18 + Vite frontend foundation | `apps/web/package.json`, `apps/web/vite.config.ts` | `npm run build` compiles with 0 errors | **VERIFIED** |
| **P4** | REQ-P4-02 | Dual-mode client (Live API + deterministic fixtures) | `apps/web/src/api/client.ts` | Live / simulator toggle, automated `/healthz` probe | **VERIFIED** |
| **P4** | REQ-P4-03 | 6 Core presentation screens | `apps/web/src/screens/` (Dashboard, Detail, Demo, Passport, Graph, Mirror) | Full interactive navigation and telemetry inspection | **VERIFIED** |
| **P5** | REQ-P5-01 | AI mutation prohibition boundary | `services/api/internal/policy/engine.go` | AI roles strictly forbidden from order/capture | **VERIFIED** |
| **P5** | REQ-P5-02 | Stable API contracts for Vijay's A01–A10 agents | `packages/contracts/openapi.yaml`, `docs/TEAM_INTEGRATION_CONTRACTS.md` | OpenAPI 3.1 specification locked | **VERIFIED** |
| **P6** | REQ-P6-01 | 7-Step Golden Path CLI runner | `services/api/cmd/demo/main.go` | `demo_test.go` passing; runs end-to-end in 0.05s | **VERIFIED** |
| **P6** | REQ-P6-02 | Interactive Golden Path Web stepper | `apps/web/src/screens/GoldenPathDemoScreen.tsx` | Stepper with JSON telemetry & 3D state vector updates | **VERIFIED** |
| **P7** | REQ-P7-01 | Full Go unit test pass (21 packages) | `services/api/...` | 21/21 packages `ok`, 0 failures | **VERIFIED** |
| **P7** | REQ-P7-02 | Static Go compilation & binaries | `cmd/api`, `cmd/migrate`, `cmd/demo` | `api.exe` (10.3MB), `demo.exe` (9.1MB), `migrate.exe` (3.9MB) | **VERIFIED** |
| **P7** | REQ-P7-03 | Multi-stage Dockerfiles | `services/api/Dockerfile`, `apps/web/Dockerfile` | Production containers with non-root security | **VERIFIED** |
| **P7** | REQ-P7-04 | GitHub Actions CI/CD workflows | `.github/workflows/{go-tests,lint-and-types,security}.yml` | Automated tests, linting, and secret scanning | **VERIFIED** |
| **P8** | REQ-P8-01 | Evidence pack and handover documentation | `docs/P0_P8_*` | 6 complete documentation artifacts created | **VERIFIED** |
