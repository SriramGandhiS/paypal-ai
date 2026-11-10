# 17 — Sriram Execution Blueprint & Antigravity Master Prompt Deck

| Field | Value |
|---|---|
| Document | `17_SRIRAM_ANTIGRAVITY_EXECUTION_AND_PROMPT_DECK.md` |
| Assignee | **Sriram** |
| Tool / Model | **Google Antigravity** |
| Role | **Platform Foundations, Infrastructure Core, DB & API Architect** |
| Primary Specs | `04_BACKEND_ENGINEERING_SPECIFICATION.md`, `05_FRONTEND_PRODUCT_AND_IMPLEMENTATION_SPEC.md`, `09_TESTING_DEVOPS_DEPLOYMENT_AND_RELEASE_PLAYBOOK.md` |
| Secondary Specs | `01_PRODUCT_BLUEPRINT_AND_MASTER_VISION.md`, `03_SYSTEM_DESIGN_AND_ARCHITECTURE.md`, `16_TEAM_WORK_ALLOCATION_AND_TIMELINE_MASTER.md` |

---

## Table of Contents
1. [Sriram's Operational Role & Scope](#1-srirams-operational-role--scope)
2. [Task Breakdown & Definition of Done (Tasks S1–S7)](#2-task-breakdown--definition-of-done-tasks-s1s7)
3. [Antigravity Master System Prompt (Copy-Paste Ready)](#3-antigravity-master-system-prompt-copy-paste-ready)
4. [Step-by-Step Task Execution Prompts for Antigravity](#4-step-by-step-task-execution-prompts-for-antigravity)
5. [Verification Commands & Evidence Collection Runbook](#5-verification-commands--evidence-collection-runbook)

---

## 1. Sriram's Operational Role & Scope

As the Foundations and API lead, Sriram is responsible for delivering the rock-solid substrate upon which Vijay's AI/Reducer engine and Diny's Enterprise streaming systems operate.

### Bound Responsibilities
- Monorepo directory setup, Docker Compose local dev environments, and Makefile automation.
- Complete PostgreSQL goose migrations (0001 through 0010), Row-Level Security (RLS) enforcement, and sqlc database queries.
- Authoritative OpenAPI 3.1 specification (`packages/contracts/openapi.yaml`) and dual-codegen pipeline (oapi-codegen for Go server, openapi-typescript for React client).
- Core Go HTTP server (`cmd/api`), Chi router, error envelope middleware, and session management.
- React 18 + TypeScript + Vite base application shell, TanStack Query client, and SSE client hub.

---

## 2. Task Breakdown & Definition of Done (Tasks S1–S7)

| Task ID | Task Title | Primary Deliverable | Acceptance Criteria |
|---|---|---|---|
| **S1** | Monorepo & Dev Rig | `infra/docker-compose.yml`, `Makefile`, `.env.example`, `.gitignore` | `docker compose up -d` starts Postgres 16 & Redis 7.2 healthy; `make test` runs without error. |
| **S2** | Config Guard | `services/api/internal/config/` | Process fails startup if required variables are missing or if `PAYPAL_ENV != sandbox` in MVP mode. Outputs variable names only, never values. |
| **S3** | PostgreSQL Migrations | `services/api/db/migrations/0001_...` to `0010_rls.sql` | `goose postgres $DATABASE_URL up` succeeds from empty DB; `internal/archtest` verifies RLS enabled on all tenant tables. |
| **S4** | OpenAPI & Codegen | `packages/contracts/openapi.yaml`, `scripts/codegen.sh` | Generates Go Chi interfaces via `oapi-codegen` and TS types via `openapi-typescript`; zero schema drift. |
| **S5** | Go Chi HTTP Server | `services/api/cmd/api/main.go`, `services/api/internal/httpx/` | Health endpoints (`/healthz`, `/readyz`), panic recovery, correlation ID middleware, and RFC 7807 error envelopes. |
| **S6** | CI/CD Pipelines | `.github/workflows/{go-tests,lint,security}.yml` | GitHub Actions run `golangci-lint`, `go test -race`, and `gitleaks` clean. |
| **S7** | React Vite Shell | `apps/web/` base scaffolding | React 18 + Vite + Tailwind CSS + TanStack Query provider running cleanly on `http://localhost:5173`; SSE hub connected. |

---

## 3. Antigravity Master System Prompt (Copy-Paste Ready)

*Instruction for Sriram: Copy everything between the block below and paste it into Antigravity at the beginning of your implementation session:*

```markdown
=== BEGIN ANTIGRAVITY SYSTEM PROMPT ===

YOU ARE: Antigravity, implementing the Foundations, Database, and API Core for PayBridge AI.
YOU ARE WORKING FOR: Sriram (Foundations Lead).
YOUR TEAMMATES ARE: Vijay Kasthuri (AI, Reducer, and UX Lead) and Diny Paul (Enterprise Infra, Multi-Rail, and Compliance Lead).

YOUR CORE RESPONSIBILITIES:
1. Deliver the monorepo directory layout, Docker Compose, Makefile, and local developer environment.
2. Build the PostgreSQL database migrations (0001 to 0010) with strict Row-Level Security (RLS) policies and sqlc query interfaces.
3. Author the authoritative OpenAPI 3.1 specification (packages/contracts/openapi.yaml) and establish automated Go and TypeScript type generation.
4. Construct the Go HTTP server (services/api/cmd/api), Chi router, middleware chain, session auth, and SSE hub.
5. Setup the React 18 + Vite + TanStack Query frontend shell.

IMMUTABLE RULES (VIOLATING ANY IS A BLOCKING DEFECT):
1. Spec-First: Update openapi.yaml and JSON schemas BEFORE implementing Go routes or TS clients.
2. Zero Secrets: Never commit, log, print, or echo any secret value. Variable names only.
3. Restricted Boundaries: You own services/api/internal/httpx/, db/, config/, and packages/contracts/. Do NOT touch or overwrite Vijay's agents/ or reducer/ code, or Diny's kafka/ or kms/ code.
4. Complete Work Only: No stubs, no placeholders, no skipped files, no TODOs.
5. Deterministic Error Envelopes: All API errors must conform to doc 04 §4 RFC 7807 problem details: { "type": "...", "title": "...", "status": 4xx, "detail": "...", "instance": "..." }.

ONLINE REFERENCES & DOCUMENTATION TO CONSULT:
- Chi Router: https://github.com/go-chi/chi
- sqlc (Type-safe SQL in Go): https://docs.sqlc.dev/en/latest/
- Goose Migrations: https://pressly.github.io/goose/
- PostgreSQL 16 RLS: https://www.postgresql.org/docs/16/ddl-rowsecurity.html
- OpenAPI 3.1 Specification: https://spec.openapis.org/oas/latest.html
- TanStack Query v5: https://tanstack.com/query/latest/docs/framework/react/overview

PRIMARY BLUEPRINT DOCUMENTS IN WORKSPACE:
Read in full:
- 04_BACKEND_ENGINEERING_SPECIFICATION.md (§2 Layout, §4 Errors, §5 Middleware, §11 Schema/RLS, §14 Routes)
- 05_FRONTEND_PRODUCT_AND_IMPLEMENTATION_SPEC.md (§1-§4 Client Architecture)
- 08_PAYPAL_INTEGRATION_SECURITY_PRIVACY_AND_COMPLIANCE.md (§21 .env.example)
- 09_TESTING_DEVOPS_DEPLOYMENT_AND_RELEASE_PLAYBOOK.md (§1-§4 Setup & CI)
- 16_TEAM_WORK_ALLOCATION_AND_TIMELINE_MASTER.md (Team Handshake Contracts)

When Sriram issues a Task, inspect the repository first, write the implementation incrementally, execute tests using run_command, verify output, and provide a clear milestone report.

=== END ANTIGRAVITY SYSTEM PROMPT ===
```

---

## 4. Step-by-Step Task Execution Prompts for Antigravity

### Task S1 Prompt: Monorepo Skeleton, Makefile & Docker Dev Rig
```text
Task S1: Setup the complete PayBridge AI repository skeleton, Docker Compose dev environment, and Makefile.
Reference: 04_BACKEND_ENGINEERING_SPECIFICATION.md §2 and 09_TESTING_DEVOPS_DEPLOYMENT_AND_RELEASE_PLAYBOOK.md §1.

Actions:
1. Create directory structure:
   - apps/web/
   - services/api/cmd/{api,worker,migrate}/
   - services/api/internal/{config,httpx,auth,policy,tenancy,cases,store}/
   - services/api/db/{migrations,queries}/
   - packages/contracts/
   - infra/docker/
   - scripts/
2. Create infra/docker/docker-compose.yml with:
   - PostgreSQL 16 (port 5432) with healthy pg_isready healthcheck.
   - Redis 7.2 (port 6379).
   - MailHog / Inbucket email sink (port 8025).
3. Create root Makefile with targets: setup, up, down, test, lint, codegen, migrate-up, migrate-status.
4. Create .env.example matching doc 08 §21 exactly. Ensure .gitignore excludes all .env files.
5. Verify: Run `docker compose -f infra/docker/docker-compose.yml up -d` and confirm containers are running.
```

### Task S2 Prompt: Configuration Validator & Production Guard
```text
Task S2: Implement the strict configuration loader and fail-fast startup validator.
Reference: 04_BACKEND_ENGINEERING_SPECIFICATION.md §1 and 08_PAYPAL_INTEGRATION_SECURITY_PRIVACY_AND_COMPLIANCE.md §8 & §21.

Actions:
1. In `services/api/internal/config/config.go`:
   - Parse all environment variables defined in doc 08 §21.
   - Validate required keys: DATABASE_URL, SESSION_SECRET, ENCRYPTION_KEY, PAYPAL_API_BASE, PAYPAL_CLIENT_ID, PAYPAL_CLIENT_SECRET.
   - Hard Guard: If APP_ENV != "production" and PAYPAL_API_BASE contains "paypal.com" without "sandbox", exit immediately with code 1.
   - Security Invariant: Error messages must NEVER print variable values, only variable names (e.g., "missing required configuration: PAYPAL_CLIENT_SECRET").
2. Write unit tests in `services/api/internal/config/config_test.go` covering missing keys, valid sandbox config, and invalid production host triggers.
3. Run `go test -v ./services/api/internal/config/...` and verify clean execution.
```

### Task S3 Prompt: PostgreSQL Goose Migrations & RLS Policies
```text
Task S3: Implement PostgreSQL migrations 0001 through 0010 using goose, and configure sqlc.
Reference: 04_BACKEND_ENGINEERING_SPECIFICATION.md §11.

Actions:
1. In `services/api/db/migrations/`:
   - 0001_extensions_roles.sql: pgcrypto, citext, restricted application role `paybridge_app`.
   - 0002_identity_tenancy.sql: tenants, users, tenant_memberships, sessions.
   - 0003_cases_attempts.sql: payment_cases, case_participants, payment_attempts (with partial unique index uq_one_active_attempt), provider_resources, idempotency_keys.
   - 0004_events_evidence.sql: payment_events, evidence_items, evidence_edges, webhook_inbox.
   - 0005_outbox_jobs.sql: outbox (job queue with FOR UPDATE SKIP LOCKED leases), dead_letters.
   - 0006_ai.sql: investigations, agent_runs, diagnosis_claims, claim_evidence, recovery_actions, action_acks.
   - 0007_passport_consent.sql: passport_versions, passport_shares, invitations, consent_grants.
   - 0008_comms_audit.sql: notifications, audit_events (hash chain), retention_jobs.
   - 0009_simulation_ops.sql: simulation_runs, route_registry, provider_connections.
   - 0010_rls.sql: ALTER TABLE ... ENABLE ROW LEVEL SECURITY; FORCE ROW LEVEL SECURITY; and create tenant isolation policies.
2. Configure `sqlc.yaml` and add core CRUD queries in `services/api/db/queries/*.sql`.
3. Run migrations against local Postgres and execute `sqlc generate`. Verify zero compile errors.
```

### Task S4 Prompt: OpenAPI 3.1 Specification & Dual Codegen
```text
Task S4: Author the authoritative OpenAPI 3.1 contract and configure automated Go and TypeScript type generation.
Reference: 04_BACKEND_ENGINEERING_SPECIFICATION.md §14 (48 Route Inventory) & §22.

Actions:
1. In `packages/contracts/openapi.yaml`, write the complete OpenAPI 3.1 spec covering:
   - Auth & Session routes: POST /api/v1/auth/login, POST /api/v1/auth/logout, GET /api/v1/auth/me.
   - Case & Checkout routes: POST /api/v1/cases, GET /api/v1/cases/{id}, POST /api/v1/cases/{id}/checkout/create-order, POST /api/v1/cases/{id}/checkout/capture.
   - Webhook ingress: POST /api/v1/webhooks/paypal.
   - Passport & Lenses: GET /api/v1/cases/{id}/passport, GET /api/v1/cases/{id}/lenses/sender, GET /api/v1/cases/{id}/lenses/receiver.
   - Real-time SSE: GET /api/v1/cases/{id}/events.
2. In `scripts/codegen.sh`, setup:
   - `oapi-codegen -generate chi-server,types,spec -package httpx -o services/api/internal/httpx/openapi.gen.go packages/contracts/openapi.yaml`
   - `npx openapi-typescript packages/contracts/openapi.yaml -o apps/web/src/types/api.gen.ts`
3. Execute `bash scripts/codegen.sh` and verify that generated Go interfaces and TS definitions compile with zero warnings.
```

### Task S5 Prompt: Go Chi HTTP Server, Middleware & Session Auth
```text
Task S5: Build the core Go HTTP server, Chi routing engine, middleware chain, and session auth.
Reference: 04_BACKEND_ENGINEERING_SPECIFICATION.md §2, §4, §5, §7.

Actions:
1. In `services/api/internal/httpx/`:
   - Router setup using `go-chi/chi/v5`.
   - Middleware: RequestID, RealIP, Slog JSON Logging (with allowlisted fields), Panic Recovery, CORS, and RateLimiter.
   - Error handling: RFC 7807 Problem Details envelope writer.
2. In `services/api/internal/auth/`:
   - Cookie-based session store backed by PostgreSQL `sessions` table.
   - Password hashing via Argon2id.
   - Session validation middleware setting tenant and user context on `r.Context()`.
3. In `services/api/cmd/api/main.go`:
   - Bind Chi server to port from config (default 8080).
   - Wire `/healthz` and `/readyz` probes.
   - Graceful shutdown handling on SIGINT/SIGTERM with 15-second drain window.
4. Test: Start server and verify curl to `http://localhost:8080/healthz` returns `{"status":"ok"}`.
```

### Task S6 Prompt: CI/CD Pipelines & Security Scan
```text
Task S6: Implement GitHub Actions CI/CD workflows for automated linting, test suites, and secret scanning.
Reference: 09_TESTING_DEVOPS_DEPLOYMENT_AND_RELEASE_PLAYBOOK.md §11.

Actions:
1. Create `.github/workflows/go-tests.yml`:
   - Spin up PostgreSQL 16 service container.
   - Run goose migrations.
   - Execute `go test -v -race ./services/api/...`.
2. Create `.github/workflows/lint-and-types.yml`:
   - Run `golangci-lint run ./services/api/...`.
   - Run `cd apps/web && npm ci && npm run typecheck && npm run lint`.
3. Create `.github/workflows/security.yml`:
   - Run `gitleaks detect --verbose`.
   - Run `govulncheck ./services/api/...`.
```

### Task S7 Prompt: React 18 + Vite Frontend Scaffolding & SSE Hub
```text
Task S7: Build the React 18 + TypeScript + Vite frontend shell with TanStack Query and SSE connectivity.
Reference: 05_FRONTEND_PRODUCT_AND_IMPLEMENTATION_SPEC.md §1–§6.

Actions:
1. In `apps/web/`:
   - Initialize React 18 + TypeScript with Vite.
   - Configure Tailwind CSS using design tokens from doc 06.
   - Install and configure TanStack Query v5 (`QueryClientProvider`).
   - Setup React Router with protected route layout.
2. In `apps/web/src/lib/sse.ts`:
   - Create resilient SSE subscriber with automatic reconnection, backoff, and `Last-Event-ID` tracking.
   - Wire SSE events to invalidate relevant TanStack Query keys.
3. Verify: `npm run dev` starts dev server without errors; `npm run build` generates production assets cleanly.
```

---

## 5. Verification Commands & Evidence Collection Runbook

After completing tasks S1 through S7, Sriram executes and captures:

```bash
# Verify Docker Infrastructure
docker compose -f infra/docker/docker-compose.yml ps

# Verify Go Server Compilation & Tests
cd services/api && go test -v -race ./...

# Verify Database Migration Status
goose -dir db/migrations postgres "$DATABASE_URL" status

# Verify Secret Safety
gitleaks detect --source . --verbose

# Verify Frontend Production Build
cd apps/web && npm run build
```
