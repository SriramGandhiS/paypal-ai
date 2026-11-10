# 09 — Testing, DevOps, Deployment and Release Playbook

| Field | Value |
|---|---|
| Document | `09_TESTING_DEVOPS_DEPLOYMENT_AND_RELEASE_PLAYBOOK.md` |
| Version | 1.0.0 |
| Date | 9 November 2026 |
| Status | **DRAFT — playbook. Every command and config below is a plan; NONE has been run. All test statuses are NOT RUN.** |
| Audience | Engineers, coding agents, the owner (demo operator) |
| Depends on | 01 (acceptance AO-xx), 03 (topology §15), 04 (routes, schema), 05 (FE tests), 07 (AI evals), 08 (security tests, `.env.example`) |
| Related | 10 (phases, evidence pack) |
| Labels | SOURCE · DESIGN DECISION · VERIFY · SANDBOX · SIMULATED · NOT RUN |

> **Honesty rule:** a test is reported only as `PASSED` / `FAILED` after it was actually executed and its output saved. Until then it is **NOT RUN**. Simulator results are always labelled **SIMULATED**; real-PayPal results **SANDBOX**.

## Table of contents
1. [Prerequisites and clean-machine setup](#1-prerequisites-and-clean-machine-setup) · 2. [Repository commands (Makefile)](#2-repository-commands-makefile) · 3. [Docker Compose](#3-docker-compose) · 4. [Environment setup](#4-environment-setup) · 5. [Migrations, seeds, reset](#5-migrations-seeds-reset) · 6. [Running backend and frontend](#6-running-backend-and-frontend) · 7. [Sandbox placeholders and webhook tunnel](#7-sandbox-placeholders-and-webhook-tunnel) · 8. [Test strategy by layer](#8-test-strategy-by-layer) · 9. [Requirements-to-test matrix](#9-requirements-to-test-matrix) · 10. [Deterministic failure simulator](#10-deterministic-failure-simulator) · 11. [Safe local/sandbox chaos testing](#11-safe-localsandbox-chaos-testing) · 12. [CI/CD](#12-cicd) · 13. [Observability: logs, metrics, traces, alerts, health](#13-observability-logs-metrics-traces-alerts-health) · 14. [Deployment (MVP first: Render)](#14-deployment-mvp-first-render) · 15. [Optional AWS reference](#15-optional-aws-reference) · 16. [HTTPS, domain, secrets, migrations, rollback, backups](#16-https-domain-secrets-migrations-rollback-backups) · 17. [Production-readiness checklist (demo-grade)](#17-production-readiness-checklist-demo-grade) · 18. [Runbooks](#18-runbooks) · 19. [Demo reset and reproducibility](#19-demo-reset-and-reproducibility) · 20. [Release gates](#20-release-gates)

---

## 1. Prerequisites and clean-machine setup

**Tooling (versions VERIFY at implementation time; use current stable):**

| Tool | Purpose |
|---|---|
| Git | Source control |
| Go (stable) | Backend |
| Node.js LTS + npm (or pnpm) | Frontend |
| Docker + Docker Compose v2 | Local Postgres/Redis/simulator |
| `make` | Task runner |
| `goose`, `sqlc`, `oapi-codegen` | Migrations / SQL codegen / OpenAPI server types |
| `golangci-lint`, `staticcheck`, `govulncheck` | Go quality/security |
| `gitleaks` | Secret scanning |
| Playwright browsers | E2E/a11y |
| A tunnel tool (any; **VERIFY** choice) | Public HTTPS URL for local webhook testing |
| `jq`, `curl` | Evidence scripts |

### 1.1 Clean-machine bootstrap (macOS/Linux; Windows via WSL2) — PLAN

```bash
# 1. Clone
git clone <REPO_URL> paybridge-ai && cd paybridge-ai

# 2. Verify toolchain
go version && node --version && npm --version && docker --version && docker compose version && make --version

# 3. Install dev tools (versions VERIFY)
make tools            # installs goose, sqlc, oapi-codegen, golangci-lint, govulncheck, gitleaks (via go install / package manager)

# 4. Environment
cp .env.example .env
#   → edit .env: set SESSION_SECRET, ENCRYPTION_KEY, DATABASE_URL password; leave PAYPAL_*/LLM_* placeholders for now
make gen-secrets      # prints generated SESSION_SECRET / ENCRYPTION_KEY to the terminal for YOU to paste (never committed)

# 5. Start infrastructure
make up               # docker compose up -d postgres (+ optional redis, mailsink)

# 6. Database
make migrate          # goose up
make seed             # demo tenant, users, invoice fixtures (labelled demo-only)

# 7. Generate contracts/types
make gen              # openapi → Go server types, TS client types; sqlc

# 8. Run (two terminals + web)
make api              # :8080
make worker
make web              # :5173

# 9. Smoke test
make smoke            # curl /healthz /readyz; logs in as demo sender; creates a SIMULATED case
```
`make gen-secrets` never writes secrets to disk or logs; it prints once.

**Clean-machine acceptance (AO-12):** record the full terminal transcript (secrets redacted) in `docs/evidence/clean-machine-run.md`; any step that fails is fixed in docs/scripts, not worked around silently.

---

## 2. Repository commands (Makefile)

```makefile
.PHONY: tools gen-secrets up down logs migrate migrate-status seed reset-dev gen api worker web \
        lint vet test test-unit test-int test-contract test-arch test-sandbox test-e2e test-a11y \
        test-ai test-security scan demo-reset evidence

tools:            ; ./scripts/install-tools.sh
gen-secrets:      ; ./scripts/gen-secrets.sh          # prints, does not store
up:               ; docker compose -f infra/docker-compose.yml up -d
down:             ; docker compose -f infra/docker-compose.yml down
logs:             ; docker compose -f infra/docker-compose.yml logs -f --tail=100
migrate:          ; cd services/api && go run ./cmd/migrate up
migrate-status:   ; cd services/api && go run ./cmd/migrate status
seed:             ; ./scripts/seed-demo.sh
reset-dev:        ; ./scripts/reset-demo.sh --dev     # drops dev DB volume; REFUSES if APP_ENV=hosted
gen:              ; ./scripts/generate.sh             # oapi-codegen + sqlc + openapi-typescript
api:              ; cd services/api && go run ./cmd/api
worker:           ; cd services/api && go run ./cmd/worker
web:              ; cd apps/web && npm run dev
lint:             ; cd services/api && golangci-lint run ./... && cd ../../apps/web && npm run lint
vet:              ; cd services/api && go vet ./... && staticcheck ./...
test: test-unit test-int test-contract test-arch
test-unit:        ; cd services/api && go test -race -count=1 ./... -short && cd ../../apps/web && npm test -- --run
test-int:         ; cd services/api && go test -race -count=1 -tags=integration ./...      # needs Postgres (Compose)
test-contract:    ; ./scripts/contract-test.sh
test-arch:        ; cd services/api && go test -count=1 ./internal/archtest/...
test-sandbox:     ; cd services/api && go test -count=1 -tags=realsandbox ./tests/real-sandbox/...   # needs real PAYPAL_* ; SANDBOX
test-e2e:         ; cd apps/web && npx playwright test
test-a11y:        ; cd apps/web && npx playwright test --project=a11y
test-ai:          ; cd services/api && go test -count=1 -tags=aieval ./tests/agents/...    # mock provider by default
test-security:    ; ./scripts/security-suite.sh
scan:             ; gitleaks detect --redact && cd services/api && govulncheck ./... && cd ../../apps/web && npm audit --audit-level=high
demo-reset:       ; ./scripts/reset-demo.sh
evidence:         ; ./scripts/validate-evidence.sh    # checks docs/evidence for redaction + completeness
```

---

## 3. Docker Compose

`infra/docker-compose.yml` (PLAN — not tested):

```yaml
name: paybridge
services:
  postgres:
    image: postgres:16        # version VERIFY
    environment:
      POSTGRES_USER: paybridge
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD:?set in .env}
      POSTGRES_DB: paybridge
    ports: ["5432:5432"]
    volumes: ["pgdata:/var/lib/postgresql/data"]
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U paybridge -d paybridge"]
      interval: 5s
      timeout: 3s
      retries: 20
  redis:
    image: redis:7
    profiles: ["redis"]
    ports: ["6379:6379"]
  mailsink:
    image: axllent/mailpit      # any SMTP sink; VERIFY
    profiles: ["dev"]
    ports: ["8025:8025", "1025:1025"]
  fake-paypal:                    # SIMULATED PayPal for offline/CI tests (never used for real proof)
    build: ../tests/fake-paypal
    profiles: ["sim"]
    ports: ["8090:8090"]
  otel-collector:
    image: otel/opentelemetry-collector
    profiles: ["obs"]
    volumes: ["./monitoring/otel.yaml:/etc/otelcol/config.yaml"]
volumes: { pgdata: {} }
```
API/worker/web run on the host via `make` in development (fast reload); a `docker compose --profile app up` variant builds containers from `services/api/Dockerfile` and `apps/web/Dockerfile` for parity tests. **Dockerfiles:** multi-stage, distroless/non-root for Go, `nginx-unprivileged` (or static host) for web.

---

## 4. Environment setup

- Source of truth for names: doc 08 §21 `.env.example`.
- **Three profiles:** `local-sim` (no PayPal/LLM keys: fake PayPal + mock LLM → everything labelled SIMULATED), `local-sandbox` (real PayPal sandbox keys, LLM optional), `hosted-demo` (platform secrets).
- Startup validator prints **names only** for missing config; `LLM_*` missing ⇒ fallback mode; `PAYPAL_*` missing ⇒ `local-sim` only, and the UI shows "Sandbox credentials not configured — running SIMULATED scenarios".
- **Never** paste secret values into chat, issues, or the evidence pack.
- Owner action list (when ready): create PayPal sandbox app/webhook (doc 08 §1) → set `PAYPAL_CLIENT_ID/SECRET/WEBHOOK_ID`; create LLM key → set `LLM_PROVIDER/MODEL/API_KEY`.

---

## 5. Migrations, seeds, reset

| Task | Command | Notes |
|---|---|---|
| Apply | `make migrate` | goose forward-only; embedded SQL in `db/migrations/` |
| Status | `make migrate-status` | CI prints status |
| Fresh DB test | CI job: empty DB → `up` → schema tests (constraints, RLS enumeration) | Required gate |
| Seed | `make seed` | Idempotent upserts: tenant `demo-co`; users `sender@demo.test`, `receiver@demo.test`, `analyst@demo.test`, `admin@demo.test`, `demo-op@demo.test` with **documented demo-only passwords generated at seed time and printed once**; sample invoices (illustrative USD 2,000.00) |
| Reset (dev) | `make reset-dev` | Drops dev volume and re-runs migrate+seed; **refuses** if `APP_ENV=hosted` |
| Reset (hosted demo) | `POST /demo/reset` (R38) | Deletes only `LOCAL_SIMULATION` + demo-tagged data; **never** deletes real-sandbox evidence or audit; confirmation required |
| Rollback | Forward-fix migrations; for emergencies restore from backup (§16) | No down migrations in prod path |

Demo data is flagged `is_demo=true`; the evidence pack must never include seeded data presented as real transactions.

---

## 6. Running backend and frontend

```bash
# Terminal A
make up && make migrate && make seed
make api            # http://localhost:8080  (/healthz /readyz /api/v1/...)
# Terminal B
make worker
# Terminal C
make web            # http://localhost:5173  (proxy /api → :8080)
```
Frontend dev proxy (`vite.config.ts`) forwards `/api` to `http://localhost:8080` so cookies are same-origin; production uses a single origin (web + API behind one host/path) or properly configured CORS + cookies (`SameSite=Lax`, domain-scoped).

**Local-sim quick demo:** `docker compose --profile sim up -d` + `PAYPAL_API_BASE=http://localhost:8090` is **forbidden** by the config guard (must be the sandbox host) — so simulation uses an explicit `PAYPAL_MODE=fake` setting that routes the adapter to the fake server **and stamps environment `PAYPAL_SIMULATOR`**, never `REAL_SANDBOX`. *(DESIGN DECISION; add `PAYPAL_MODE=real|fake` to `.env.example` when implementing; `real` requires the sandbox host.)*

---

## 7. Sandbox placeholders and webhook tunnel

- Real-sandbox tests are tagged and **skipped by default**; they run only when `PAYPAL_CLIENT_ID`, `PAYPAL_CLIENT_SECRET`, and sandbox buyer approval are available.
- **Buyer approval step is interactive** in PayPal's sandbox UI; the `test-sandbox` suite supports (a) a manual-approval mode that prints the approval link and waits, and (b) optional browser automation via Playwright **only if** PayPal's sandbox pages permit automation (**VERIFY**; otherwise stay manual). Never script around anti-automation controls.
- **Local webhook testing:** start a tunnel to `localhost:8080`, set `PUBLIC_BASE_URL`, register `${PUBLIC_BASE_URL}/api/v1/webhooks/paypal` in the sandbox app, copy the webhook id to `PAYPAL_WEBHOOK_ID`. When the tunnel URL changes, update the webhook (record this in the runbook).
- For hosted demos use the Render service URL (stable).

---

## 8. Test strategy by layer

*(Every row: **NOT RUN**.)*

| Layer | Scope | Tools | Run in CI? | Mode |
|---|---|---|---|---|
| **Unit (Go)** | money, validators, classifier, reducer transitions, safe-action property tests, projector allowlists, claim validator, backoff | `testing`, `testify`, property testing | Yes | – |
| **Unit (TS)** | `statusPresentation` exhaustiveness, formatters, idempotency lifecycle, forbidden-copy | Vitest | Yes | – |
| **DB integration** | migrations from empty; all constraints (I1, I2, I7, I11, one-active-attempt); RLS enumeration; append-only triggers | Postgres (Compose/testcontainers) | Yes | – |
| **Handler/contract** | Every route vs OpenAPI; auth matrix; uniform 404; idempotency replay/reuse/in-flight | `httptest`, contract validator | Yes | – |
| **Webhook & reconciliation** | T03–T05; CHAOS-01/02/04/05/11/12 against fake PayPal; assert **0 POSTs** from Reconciler | fake PayPal | Yes | **SIMULATED** |
| **Concurrency** | Two workers/one webhook; double-click capture; stale snapshot | `-race`, harness | Yes | – |
| **Architecture** | Import rules; Reader GET-only; every route authorises; no route without OpenAPI | `archtest` | Yes | – |
| **Real PayPal sandbox** | Create→approve→capture→GET; real webhook delivery; CHAOS-01 with Fault Proxy | tagged tests + scripts | **No** (secrets; manual/opt-in) | **SANDBOX** |
| **Security** | SEC-01…20 (doc 08 §22) | scripts + Playwright | Yes (subset) | mixed |
| **AI evaluation** | T01–T10, adversarial, privacy, consistency (doc 07 §15) | harness (mock provider default; real provider opt-in) | Mock: yes; real: manual | SIMULATED (mock) / real |
| **E2E** | Golden path; receiver consent+revocation; duplicate-click; SSE fallback | Playwright on Compose stack | Yes (sim mode) | SIMULATED; one SANDBOX run manual |
| **Accessibility** | axe on 16 screens; keyboard walkthrough; manual screen-reader pass | axe + manual | axe in CI; manual per release | – |
| **Visual regression** | 16 screens × 3 breakpoints × 2 themes | Playwright screenshots | Yes (baseline review) | – |
| **Load (modest)** | Webhook ingest burst (e.g. N events/s), SSE connections, read APIs | `k6` or `vegeta` | Nightly/manual | SIMULATED |
| **Performance (FE)** | bundle budgets | `size-limit` | Yes | – |
| **Backup/restore drill** | Restore into scratch DB; verify audit chain + counts | script | Manual per release | – |

**Test data rules:** fixtures live in `tests/fixtures/` with a manifest stating `SIMULATED`/`SANDBOX-derived (redacted)`; no real personal data; no secrets.

---

## 9. Requirements-to-test matrix

| Req ID | Requirement | Source | Primary test(s) | Mode | Evidence artefact | Status |
|---|---|---|---|---|---|---|
| REQ-01 | Real sandbox create→approve→capture with provider ids stored; independent GET matches | AO-01, G1 | `tests/real-sandbox/TestCreateApproveCapture` | SANDBOX | `docs/evidence/paypal-sandbox-receipt.md` | NOT RUN |
| REQ-02 | Exactly 1 order + 1 capture after lost response | AO-02, G2, C03 | `TestLostResponseReconcile` (SIM) + sandbox CHAOS-01 run | SIM + SANDBOX | DB ledger export; API trace | NOT RUN |
| REQ-03 | Unknown never rendered as failed; never retry-checkout | AO-03 | FE `statusPresentation` + E2E unknown path; BE transition tests | SIM | Test report | NOT RUN |
| REQ-04 | Forged webhook rejected; duplicate = 1 effect | AO-04, G3 | `TestWebhookForged`, `TestWebhookDuplicate` (T03/T05) | SIM + SANDBOX | Webhook test log | NOT RUN |
| REQ-05 | No state regression on out-of-order events | AO-05 | `TestReducerOutOfOrder` (T04, CHAOS-05) | SIM | Transition-test report | NOT RUN |
| REQ-06 | 100 % confirmed claims have valid evidence | AO-06 | `claims` validator tests; eval harness | SIM/real LLM | Validator report | NOT RUN |
| REQ-07 | 0 unsupported material claims in eval suite (state n) | AO-07 | `tests/agents` T01–T10 + adversarial | mock/real | AI eval report | NOT RUN |
| REQ-08 | 0 role-privacy leaks | AO-08 | Privacy diff E2E; DLP tests; SEC-17 | SIM | Permission diff report | NOT RUN |
| REQ-09 | 0 cross-tenant reads | AO-09 | SEC-03/04, CHAOS-06 | – | Security report | NOT RUN |
| REQ-10 | Truth and guidance work with LLM disabled | AO-10, CHAOS-10 | E2E Demo E; unit fallback tests | SIM/SANDBOX | Recording + test | NOT RUN |
| REQ-11 | Simulated always labelled | AO-11 | UI/API audit script; `environment` assertions | SIM | Audit output | NOT RUN |
| REQ-12 | Fresh clone → running demo via docs | AO-12 | Clean-machine run | – | Terminal transcript | NOT RUN |
| REQ-13 | WCAG 2.2 AA target on 16 screens | AO-13, FE-09 | axe + manual | – | `a11y-report.md` | NOT RUN |
| REQ-14 | 0 secrets in repo/history/logs/bundle | AO-14 | gitleaks, bundle scan, log scan | – | Scan output | NOT RUN |
| REQ-15 | Provider/currency mismatch blocks success label | CHAOS-12, A07 | `TestAmountMismatch` | SIM | Test report | NOT RUN |
| REQ-16 | DB commit failure after provider success is repairable | CHAOS-11 | `TestCommitFailureAfterCapture` | SIM | Test report | NOT RUN |
| REQ-17 | Consent revocation blocks access, stream, notices | CHAOS-09, SEC-16 | E2E + integration | SIM | Test report | NOT RUN |
| REQ-18 | Prompt injection has no tool/state/privacy effect | CHAOS-07, SEC-11 | Adversarial suite | mock/real | AI eval report | NOT RUN |
| REQ-19 | Public status never marks a case failed | T09, CHAOS-08 | A02 tests; Ops UI test | SIM | Test report | NOT RUN |
| REQ-20 | Reader interface GET-only; AI cannot import mutators | G-PAY-1 | `archtest` | – | Test report | NOT RUN |
| REQ-21 | Idempotency: replay returns stored response; key reuse 422 | doc 04 §10 | Handler tests | – | Test report | NOT RUN |
| REQ-22 | Bundle budgets, responsive layouts | FE-10/12 | `size-limit`, visual regression | – | CI output | NOT RUN |

**Traceability rule:** every acceptance outcome (AO-xx), frontend criterion (FE-xx), security test (SEC-xx), and eval ID (T01–T10) must appear in at least one row or be explicitly marked *deferred* with a reason (doc 10 audit).

---

## 10. Deterministic failure simulator

**Purpose:** reproducible, labelled failure injection to prove reconciliation and safety — **without** claiming it equals real PayPal behaviour. Lives in `tests/fake-paypal` + `internal/simulator` + `internal/faultproxy`. All simulator-produced records carry `environment=PAYPAL_SIMULATOR` or `LOCAL_SIMULATION`, evidence tier E6, `source_trust=SIMULATED`, and ids prefixed `SIM-`.

### 10.1 Scenarios (deterministic via seed + scripted timeline)

| Sim ID | Scenario | Mechanism | Maps to | Expected outcome (assertions) |
|---|---|---|---|---|
| SIM-01 | **Delayed response** (slow capture, response arrives after client timeout) | Fake PayPal sleeps > client timeout; completes capture | CHAOS-02 variant | Attempt `CAPTURE_UNKNOWN`; reconcile finds COMPLETED; **0 additional captures** |
| SIM-02 | **Capture timeout, provider success** | Fake PayPal captures then drops response | CHAOS-01 (SIM form) | Same as above; ledger orders=1, captures=1 |
| SIM-03 | **Duplicate webhook** | Deliver same event id N times | CHAOS-04, T03 | 1 `payment_events` row; duplicate_count=N−1; 1 notification |
| SIM-04 | **Out-of-order webhook** | `COMPLETED` then older `PENDING` | CHAOS-05, T04 | State stays COMPLETED; evidence appended |
| SIM-05 | **Provider API unavailable** | Fake returns 503/connection refused for N seconds | – | Reads retry with backoff; case stays UNKNOWN; no mutator retry; recovers on restore |
| SIM-06 | **Unknown outcome** (order-create response lost) | Drop response after order creation | S6 | `ORDER_UNKNOWN`; no second create; REVIEW_REQUIRED if lookup impossible |
| SIM-07 | **Reconciliation** (happy and conflict paths) | Fake returns COMPLETED / PENDING / contradictory / inconclusive | S5 | Matches outcome table doc 04 §15.1; **Reconciler issued zero POSTs** (verb counter) |
| SIM-08 | Forged webhook | Bad signature to the test verifier | CHAOS-05/T05 | Rejected; no state change |
| SIM-09 | Commit failure after capture | Fail DB tx2 after provider 2xx | CHAOS-11 | Persisted intent + original keys; reconcile repairs |
| SIM-10 | Amount/currency mismatch | Fake returns different amount | CHAOS-12 | `AMOUNT_MISMATCH`; success label blocked; review |
| SIM-11 | LLM unavailable/invalid | Mock provider errors | CHAOS-10 | Fallback Passport; Demo E |
| SIM-12 | Receiver consent revoked mid-run | Revoke during investigation | CHAOS-09 | Receiver output discarded; notices cancelled |

### 10.2 Controls
- **Deterministic:** scenario scripts specify event times/ordering; clocks injectable (`clock` package).
- **Verb counter:** fake PayPal records every request (method, path); tests assert mutator counts.
- **Fault Proxy (real sandbox, CHAOS-01):** a small component inside PayBridge's own outbound path used **only** when `FAULT_PROXY_ENABLED=true`, `APP_ENV` is local/sandbox demo, and a demo operator toggles `DROP_CAPTURE_RESPONSE`. It forwards the capture to PayPal, **discards the response**, and surfaces a timeout to the caller. It never alters PayPal or sends extra mutating calls. Disabled and compiled out (`//go:build faultproxy`) in any build that is not demo/local.
- **Labelling:** UI shows `SIMULATED` ribbon for SIM runs and `REAL SANDBOX + FAULT` chip for CHAOS-01 runs.

---

## 11. Safe local/sandbox chaos testing

| Rule | Detail |
|---|---|
| Scope | Chaos limited to PayBridge's own components and the PayPal **sandbox**; never any production system or real credentials |
| Blast radius | Faults injected per case/test, auto-expire; no global kill-switches left on (`FAULT_PROXY_ENABLED` defaults false; CI asserts false in hosted config) |
| Provider etiquette | No load testing against PayPal sandbox beyond normal functional usage (**VERIFY** terms); use fake PayPal for volume |
| Data | Only demo/test data; no real personal data |
| Observability | Every chaos run writes a report (`docs/evidence/chaos-<id>-<date>.md`) with scenario, mode, assertions, result, redacted logs |
| Recovery check | After each chaos test: DB invariants query (one active attempt, no duplicate provider ids, audit chain OK) must pass |
| Game-day | Optional pre-demo rehearsal running SIM-01…12 + sandbox CHAOS-01/04; results recorded **as run** |

---

## 12. CI/CD

### 12.1 Workflows (GitHub Actions — PLAN)

| Workflow | Triggers | Jobs |
|---|---|---|
| `go-tests.yml` | PR, push | `gen-check` (generated code up to date) → `lint/vet/staticcheck` → `unit -race` → `integration` (Postgres service) → `contract` → `archtest` |
| `frontend-tests.yml` | PR, push | `gen-check` → `lint/tsc` → `vitest` → `size-limit` → `build` |
| `security-scan.yml` | PR, push, nightly | `gitleaks` (tree+history on schedule) → `govulncheck` → `npm audit` → `osv-scanner` → container scan |
| `e2e.yml` | PR (label), nightly | Compose stack in **sim mode** → Playwright (golden path, privacy, a11y/axe) → upload traces |
| `ai-eval.yml` | PR touching `agents/`, nightly | Mock-provider eval (T01–T10 + adversarial) → report artefact; **real-provider eval is manual** (`workflow_dispatch` with secret) |
| `deploy-staging.yml` | merge to `main` | Build images → push → deploy to Render (or chosen platform) → run migrations → post-deploy smoke (`/healthz`, `/readyz`, login, SIMULATED scenario) |
| `release.yml` | tag `v*` | Gates (§20) → SBOM → release notes → deploy |

### 12.2 Pipeline rules
- Required checks on `main` branch protection; no direct pushes; reviews required.
- **No secrets in PR workflows from forks**; real-sandbox and real-LLM jobs run only on protected branches/manual dispatch using repository secrets.
- Pin action versions/SHAs; least-privilege `permissions:` per job.
- Artefacts: test reports, coverage (reported only if measured), Playwright traces, eval report, SBOM.

Example job (excerpt):
```yaml
jobs:
  go:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:16
        env: { POSTGRES_PASSWORD: ci-only, POSTGRES_DB: paybridge_test }
        ports: ["5432:5432"]
        options: >-
          --health-cmd "pg_isready -U postgres" --health-interval 5s --health-timeout 3s --health-retries 20
    steps:
      - uses: actions/checkout@<pinned>
      - uses: actions/setup-go@<pinned>
        with: { go-version-file: services/api/go.mod }
      - run: make gen && git diff --exit-code
      - run: make lint vet
      - run: make test
```

---

## 13. Observability: logs, metrics, traces, alerts, health

*(Metric names and log fields: doc 04 §21.)*

| Signal | Design |
|---|---|
| **Logs** | JSON to stdout; allowlisted fields; `correlation_id` everywhere; platform log drain (Render) |
| **Metrics** | Prometheus-style `/metrics` on an internal-only port/path (auth-protected or not exposed publicly); or OTLP to collector in `obs` profile |
| **Traces** | OpenTelemetry; spans for HTTP, DB, PayPal, reducer, agents |
| **Dashboards** (optional) | Compose `obs` profile with a metrics UI; screenshot into evidence pack **only with measured data** |
| **Health** | `/healthz` (liveness), `/readyz` (DB + migrations + config); LLM/PayPal reachability = degraded indicators in Ops → Diagnostics |

### Alerts (initial thresholds are DESIGN DECISION — tune from data)

| Alert | Condition | Severity | Runbook |
|---|---|---|---|
| Webhook signature failure spike | `pb_webhook_signature_failures_total` rate > baseline | S2 | RB-01 |
| Inbox backlog | `pb_inbox_unprocessed` > N for > 5 min | S2 | RB-01 |
| Unknown-state backlog age | `pb_unknown_state_backlog_age_seconds` > T | S2 | RB-02 |
| Dead letters | `pb_dead_letters_total` increase | S2 | RB-03 |
| Provider read failures | high `pb_provider_calls_total{class="unknown"}` | S3 | RB-04 |
| LLM failures/spend | error rate or spend > cap | S3 | RB-05 |
| Cross-tenant denial spike | policy-deny rate anomaly | S1/S2 | RB-06 |
| Audit chain failure | `chain_ok=false` | S1 | RB-06 |
| API 5xx/latency | SLO breach | S3 | RB-07 |
Alert delivery for the hackathon: a single channel (email/chat webhook) owned by the owner.

---

## 14. Deployment (MVP first: Render)

*(Plan features/pricing/sleep behaviour **VERIFY** at deployment time; the API must be always-on for webhook reliability.)*

| Component | Render service type | Notes |
|---|---|---|
| Web App | Static Site (build `npm ci && npm run build`, publish `apps/web/dist`) | SPA rewrite to `/index.html`; headers for CSP/cache; `/api/*` rewritten/proxied to API (same-origin cookies) or CORS configured |
| API | Web Service (Docker or native Go) | Health check `/readyz`; env vars/secrets from dashboard; **no sleep/spin-down** plan |
| Worker | Background Worker | Same image, command `worker`; `LLM_API_KEY` only here |
| PostgreSQL | Managed Postgres | Private connection string; backups enabled (plan-dependent, VERIFY) |
| Redis | Optional | Skip in MVP |

`infra/render.yaml` (Blueprint — sketch, **VERIFY syntax**):
```yaml
services:
  - type: web
    name: paybridge-api
    runtime: docker
    dockerfilePath: services/api/Dockerfile
    healthCheckPath: /readyz
    envVars:
      - { key: APP_ENV, value: hosted }
      - { key: PAYPAL_ENV, value: sandbox }
      - { key: PAYPAL_API_BASE, value: https://api-m.sandbox.paypal.com }
      - { key: DATABASE_URL, fromDatabase: { name: paybridge-db, property: connectionString } }
      - { key: SESSION_SECRET, sync: false }
      - { key: ENCRYPTION_KEY, sync: false }
      - { key: PAYPAL_CLIENT_ID, sync: false }
      - { key: PAYPAL_CLIENT_SECRET, sync: false }
      - { key: PAYPAL_WEBHOOK_ID, sync: false }
      - { key: FAULT_PROXY_ENABLED, value: "false" }
  - type: worker
    name: paybridge-worker
    runtime: docker
    dockerfilePath: services/api/Dockerfile
    dockerCommand: /app/worker
    envVars:
      - { key: LLM_API_KEY, sync: false }
      - { key: LLM_PROVIDER, sync: false }
      - { key: LLM_MODEL, sync: false }
  - type: web
    name: paybridge-web
    runtime: static
    buildCommand: cd apps/web && npm ci && npm run build
    staticPublishPath: apps/web/dist
databases:
  - name: paybridge-db
    plan: basic-256mb   # VERIFY
```
**Deploy procedure:** merge → CI green → image build → **run migrations as a pre-deploy step** (`/app/migrate up`; backward-compatible migrations only) → deploy API → deploy worker → web → smoke tests → register/verify webhook URL = `https://<api-host>/api/v1/webhooks/paypal` → tag release.

---

## 15. Optional AWS reference

*(Reference only — **not provisioned**; see diagram in doc 03 §15.3.)* CloudFront + S3 (web), ALB + ECS Fargate (API, Worker), RDS PostgreSQL Multi-AZ with PITR, Secrets Manager, KMS for evidence encryption, private subnets with controlled egress to PayPal/LLM, WAF on ALB/CloudFront, CloudWatch/OTel. IaC in `infra/terraform/aws-reference/` as documentation until a real need arises.

---

## 16. HTTPS, domain, secrets, migrations, rollback, backups

| Topic | Plan |
|---|---|
| **HTTPS** | Platform-managed TLS; HSTS; redirect HTTP→HTTPS; webhook endpoint must be HTTPS |
| **Domain** | Optional custom domain (DNS CNAME to platform); update `PUBLIC_BASE_URL`, `WEB_ORIGIN`, PayPal webhook URL, cookie domain |
| **Secrets** | Platform secret store; no secrets in `render.yaml`/repo; rotation runbook RB-08 |
| **Migrations** | Pre-deploy, forward-only, expand/contract pattern (add → backfill → switch → remove in later release); never destructive in same release as code change |
| **Rollback** | App: redeploy previous image (migrations are backward-compatible so old code still runs); Data: restore from backup to a **new** instance, verify, then cut over; never `DROP` in emergency |
| **Backups** | Managed daily backups (+PITR if plan supports; VERIFY); **monthly restore drill** (and before the demo): restore into scratch DB, run `audit verify` + row-count checks; record in `docs/evidence/restore-drill.md` |
| **Evidence encryption keys** | `ENCRYPTION_KEY` backed up in the secret store; losing it makes encrypted payloads unrecoverable (documented) |
| **Data residency** | Not claimed; region chosen by the owner and stated |

---

## 17. Production-readiness checklist (demo-grade)

> This is a **demo-deployment checklist**. It does **not** make the system "production ready" and must not be used to claim that.

- [ ] `APP_ENV=hosted`, `PAYPAL_ENV=sandbox`, sandbox host verified by the config guard
- [ ] `FAULT_PROXY_ENABLED=false` by default; toggle restricted to `demo_operator`
- [ ] All secrets in the platform store; none in repo/history (gitleaks clean)
- [ ] API always-on; webhook URL registered and **verified end-to-end** with a real sandbox event
- [ ] Migrations applied; `/readyz` green; seed demo accounts created with rotated passwords
- [ ] Backups enabled; restore drill done
- [ ] Alerts configured to a monitored channel
- [ ] CSP/headers verified; cookies Secure
- [ ] Rate limits on; body limits on
- [ ] Demo banner + disclaimer present; environment badge visible
- [ ] Evidence pack generated from **actual** runs (`make evidence` passes)
- [ ] README distinguishes built/tested/planned; claims checklist signed (doc 10)

---

## 18. Runbooks

*(Each runbook lives in `docs/runbooks/RB-xx-*.md`; summaries below.)*

| ID | Situation | Steps |
|---|---|---|
| **RB-01** | Webhook failures / inbox backlog | Check `/ops/webhooks`; verify `PUBLIC_BASE_URL` and the registered webhook URL; check signature-failure metric (config drift: wrong `PAYPAL_WEBHOOK_ID`?); confirm DB health; scale/restart worker; run reconcile sweeper (polling covers missed events); never ack unverified payloads |
| **RB-02** | Unknown-state backlog growing | List cases in `API_TIMEOUT/RECONCILING/EVENT_GAP` by age; check provider read errors (RB-04); verify worker leases; manually trigger `POST /cases/{id}/reconcile`; escalate cases `NEEDS_HUMAN_REVIEW`; **do not create new orders** |
| **RB-03** | Dead-letter jobs | Inspect `dead_letters` last_error; fix cause; re-queue via admin (audited); if a case-level job, ensure case has human-review flag |
| **RB-04** | PayPal API errors/unavailable | Check public status (advisory); confirm token refresh working; check rate limiting; keep UNKNOWN states; backoff continues; communicate via banner "Verification delayed" (never "failed") |
| **RB-05** | LLM down / spend spike | Set `LLM_ENABLED=false`; verify fallback Passport; inspect spend metric; rotate key if abuse suspected; re-enable after cap review |
| **RB-06** | Suspected cross-tenant/role leak or audit-chain failure | S1: disable affected lens/route; revoke sessions/shares; preserve evidence; run IDOR/DLP suites; review audit chain; notify per incident policy (doc 08 §19) |
| **RB-07** | API latency/5xx | Check DB connections/locks; slow query log; restart; scale replicas; roll back release if regression |
| **RB-08** | Secret rotation | PayPal: new secret in console → update store → redeploy → health check → revoke old. LLM: new key → update worker env → restart. `SESSION_SECRET`: rotation invalidates sessions (announce). `ENCRYPTION_KEY`: introduce new key id; re-encrypt gradually; keep old for decrypt |
| **RB-09** | Wrong statement published in a Passport | Withdraw (new version), publish correction, audit, add regression test |
| **RB-10** | Restore from backup | Provision new DB, restore, run migrations status, `audit verify`, reconcile open cases, cut over, record |
| **RB-11** | Webhook URL changed (tunnel/host) | Update PayPal webhook URL, confirm new `PAYPAL_WEBHOOK_ID` if re-created, send a sandbox event, verify inbox row `PROCESSED` |

---

## 19. Demo reset and reproducibility

**Goal:** the hero demo (Demo B) reproducible from a clean state in under three minutes, every time.

```bash
make demo-reset          # clears LOCAL_SIMULATION + demo-tagged data; preserves real-sandbox evidence & audit
# Demo B steps (also available as one click in /app/demo):
# 1 login as sender → open demo invoice → create case → approve in PayPal sandbox (buyer account)
# 2 toggle Fault Proxy DROP_CAPTURE_RESPONSE → confirm payment → UI shows unknown + Do-not-repeat banner
# 3 worker reconciles original order/capture → Passport v2; webhook arrives (+ replay duplicate)
# 4 open sender lens, receiver lens (accepted invite), Truth Gap Graph, Ledger: Orders 1 · Captures 1
# 5 analyst approves internal reconciliation → case RESOLVED
make evidence            # validates redaction + completeness of docs/evidence/
```
- **Pre-demo checklist:** sandbox buyer/merchant accounts verified; webhook URL reachable; one rehearsal run recorded; LLM key present or Demo E ready; backup of DB; screen resolution/zoom set; secrets hidden in terminal/recording; Fault Proxy off after demo.
- **Fallbacks:** if PayPal sandbox is down → show prior **recorded real run** (labelled recorded, with date) plus live SIMULATED Mirror; if LLM down → Demo E; if network flaky → local-sim profile with banner.
- **Reproducibility evidence:** each rehearsal stored with date, commit hash, environment, result — as measured.

---

## 20. Release gates

A tagged release / demo deployment requires every gate **evidenced** (not asserted):

| Gate | Check | Evidence |
|---|---|---|
| RG-01 | All required CI jobs green on the release commit | CI run link/export |
| RG-02 | Requirements-to-test matrix: every MUST row **PASSED** with artefact, or the release notes list it as NOT RUN/deferred | Filled matrix |
| RG-03 | Real-sandbox proof (REQ-01) executed within this release cycle | Redacted receipt |
| RG-04 | CHAOS-01 (sandbox+fault) and CHAOS-04 executed and recorded | Chaos reports |
| RG-05 | Security gates G-SEC-1…7, G-PAY-1/2 (doc 08 §20) | Security report |
| RG-06 | AI eval executed on the release prompt versions; suite size + failures reported | AI eval report |
| RG-07 | A11y: axe clean + manual keyboard/screen-reader pass recorded | a11y report |
| RG-08 | Clean-machine run succeeds | Transcript |
| RG-09 | Backup restore drill recorded | Restore report |
| RG-10 | Honest-claims checklist signed (doc 10) | Signed checklist |
| RG-11 | Rollback tested (previous image redeploy) | Note |
| RG-12 | README/CHANGELOG accurate: built vs tested vs planned; labels intact | Review |

*End of document 09.*
