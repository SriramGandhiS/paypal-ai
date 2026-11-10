# 04 — Backend Engineering Specification (Go)

| Field | Value |
|---|---|
| Document | `04_BACKEND_ENGINEERING_SPECIFICATION.md` |
| Version | 1.0.0 |
| Date | 9 November 2026 |
| Status | **DRAFT — specification only. No Go code, migration, or test described here exists yet.** |
| Audience | Backend engineers and coding agents |
| Depends on | 01 (vocabulary §12), 03 (architecture, ADRs) |
| Related | 05 Frontend (consumes routes §14) · 07 AI (tools §9 of doc 07) · 08 (security controls, `.env.example`) · 09 (tests) |
| Labels | SOURCE · DESIGN DECISION · VERIFY · SANDBOX · SIMULATED |

> **Change note (2026-11-09, Account 2 audit):** §9.1 previously placed `VerifyWebhookSignature` (an HTTP POST) inside `Reader`, contradicting §2.2 ("Reader has no non-GET method") and doc 03 ("Reader interface has no POST methods"). Corrected by splitting out `WebhookVerifier` (decision D-14, ADR-021). No other content changed.

> **Hard rule restated:** the LLM must never override provider state. The only writer of `provider_state`, `application_state`, `receiver_state` is the **Reducer** (ADR-015).

## Table of contents
1. [Scope and stack](#1-scope-and-stack) · 2. [Go module and package layout](#2-go-module-and-package-layout) · 3. [Domain model](#3-domain-model) · 4. [API conventions and errors](#4-api-conventions-and-errors) · 5. [Middleware chain](#5-middleware-chain) · 6. [Validation](#6-validation) · 7. [Auth and sessions](#7-auth-and-sessions) · 8. [RBAC, ABAC and case access](#8-rbac-abac-and-case-access) · 9. [PayPal adapter and OAuth lifecycle](#9-paypal-adapter-and-oauth-lifecycle) · 10. [Idempotency, timeouts, rate limits](#10-idempotency-timeouts-rate-limits) · 11. [PostgreSQL schema](#11-postgresql-schema) · 12. [Webhook pipeline](#12-webhook-pipeline) · 13. [Inbox/outbox workers, retry, dead-letter](#13-inboxoutbox-workers-retry-dead-letter) · 14. [Route inventory](#14-route-inventory) · 15. [Reconciliation and transition rules](#15-reconciliation-and-transition-rules) · 16. [Evidence, provenance, Truth Gap Graph, Passport](#16-evidence-provenance-truth-gap-graph-passport) · 17. [Next Safe Action Engine](#17-next-safe-action-engine) · 18. [Two-Lens privacy projections](#18-two-lens-privacy-projections) · 19. [Audit, retention, deletion, consent](#19-audit-retention-deletion-consent) · 20. [Transactions, locking, concurrency](#20-transactions-locking-concurrency) · 21. [Observability](#21-observability) · 22. [OpenAPI](#22-openapi) · 23. [Tests](#23-tests) · 24. [Implementation sequence](#24-implementation-sequence)

---

## 1. Scope and stack

| Item | Decision |
|---|---|
| Language/runtime | Go (version **VERIFY**; use the latest stable supported release at implementation time) |
| Processes | `cmd/api`, `cmd/worker`, `cmd/migrate` (ADR-005) |
| Router | `chi` (ADR-011) |
| DB access | `pgx/v5` + `sqlc`; migrations with `goose` (ADR-014) |
| Money | decimal string at API; `NUMERIC(20,6)` in SQL; Go decimal type, never `float64` (D-09) |
| IDs | UUIDv4 for PayBridge ids; opaque strings for PayPal ids, unique on `(environment, resource_type, id)` |
| Logging | `log/slog` JSON with allowlisted fields |
| Telemetry | OpenTelemetry traces/metrics (SDK versions **VERIFY**) |
| Contract | `packages/contracts/openapi.yaml` is the source of truth (ADR-020) |
| Lint/safety | `go vet`, `staticcheck`, `govulncheck`, `golangci-lint`, `depguard` import rules, `-race` tests |

---

## 2. Go module and package layout

```text
services/api/
├── cmd/
│   ├── api/main.go          # HTTP server, wiring, graceful shutdown
│   ├── worker/main.go       # outbox relay, inbox processor, reconciler, orchestrator, notifier, retention
│   └── migrate/main.go      # goose runner (up / status / reset-dev-only)
├── internal/
│   ├── config/              # env parsing + fail-fast validator (names only in errors)
│   ├── httpx/               # router assembly, middleware, error envelope, pagination, SSE helpers
│   ├── auth/                # sessions, password/demo-login, CSRF, OIDC interface (later)
│   ├── policy/              # Policy Engine: RBAC + ABAC + consent decisions (single choke point)
│   ├── tenancy/             # tenant context, SET LOCAL app.tenant_id helper
│   ├── cases/               # case commands/queries, invitations, workflow transitions
│   ├── consent/             # grants, field allowlists, revocation, share tokens
│   ├── checkout/            # ONLY caller of paypal.Mutator: create order, capture, attempt locking
│   ├── paypal/              # adapter: oauth.go, reader.go, mutator.go, classify.go, types.go, verify.go
│   ├── webhook/             # Gateway: raw capture, verify, inbox insert, 2xx
│   ├── inbox/               # inbox processor (lease, correlate, normalise)
│   ├── outbox/              # relay + job leasing + dead-letter
│   ├── events/              # canonical event envelope, schema version, hashing
│   ├── evidence/            # evidence items, provenance edges, tiers
│   ├── reducer/             # state vector transition function (only state writer)
│   ├── reconciliation/      # Reconciler (Reader-only)
│   ├── safeaction/          # Next Safe Action Engine
│   ├── projector/           # role projections (sender/receiver/ops/ai_snapshot)
│   ├── passport/            # versioned Passport builder + hashing
│   ├── truthgap/            # Truth Gap Graph builder + gap detector
│   ├── orchestrator/        # fixed DAG, router, budgets, snapshots
│   ├── agents/              # A01..A10 definitions (prompts, schemas, tool allowlists)
│   ├── tools/               # Tool Gateway: typed read-only tools
│   ├── llmgw/               # provider-neutral LLM interface + adapters
│   ├── claims/              # Claim Validator (schema, citation, tier, privacy, temporal)
│   ├── notifications/       # Notifier, delivery records
│   ├── audit/               # hash-chained audit writer/verifier
│   ├── health/              # /healthz, /readyz, provider-status advisory, own telemetry
│   ├── simulator/           # scenarios, fake PayPal, SIMULATED event factory
│   ├── faultproxy/          # CHAOS-01 response-drop (build-flag gated)
│   ├── money/               # decimal helpers, ISO-4217 exponents
│   ├── redact/              # PII/secret redaction, log allowlist
│   ├── sse/                 # per-case SSE hub (Last-Event-ID replay)
│   └── store/               # sqlc-generated code, tx helpers, RLS session setup
├── db/migrations/           # 0001_*.sql … (goose)
├── db/queries/              # *.sql for sqlc
└── go.mod
```

### 2.1 Import rules (enforced by `depguard` + an `internal/archtest` test)

| Package | MAY import | MUST NOT import |
|---|---|---|
| `orchestrator`, `agents`, `tools`, `llmgw`, `claims` | `store` (read via tools only), `evidence`, `projector`, `events`, `redact`, `money` | `checkout`, `paypal` (any), `webhook`, `reducer` write APIs |
| `reconciliation` | `paypal` **Reader interface only**, `reducer`, `evidence`, `store` | `paypal.Mutator`, `checkout` |
| `checkout` | `paypal.Mutator` + `paypal.Reader`, `store`, `audit` | `orchestrator`, `agents`, `tools`, `llmgw` |
| `webhook` | `paypal` (verify only), `store` (inbox writes), `events` | `reducer`, `cases` mutation |
| `projector` | `store` read | any generation package (it builds *inputs*, never text) |
| `reducer` | `events`, `evidence`, `store`, `audit` | `llmgw`, `agents` |

### 2.2 Wiring principle
Constructors take interfaces; `main.go` is the only place concrete types meet. Tests inject fakes. A **compile-time assertion** (`var _ paypal.Reader = (*readerClient)(nil)`) and a reflection test confirm `Reader` has **no** method whose HTTP verb is not GET.

---

## 3. Domain model

### 3.1 Aggregates (SOURCE annex §8)

| Aggregate | Meaning | Key invariants |
|---|---|---|
| `PaymentCase` | Business contract: payer, payee, expected amount/currency, purpose | Not a bank account; carries the state vector + `state_version` |
| `PaymentAttempt` | Exactly one initiated provider order/operation with idempotency key | Never created implicitly by a timeout |
| `ProviderResource` | Order / authorization / capture / refund ids and statuses | Unique `(environment, resource_type, provider_resource_id)` |
| `PaymentEvent` (VerifiedEvent) | Immutable observation with provenance and time | Append-only |
| `EvidenceItem` | Source-backed fact with tier, trust, hash | Immutable; redacted payload |
| `Investigation` / `AgentRun` / `DiagnosisClaim` | Time-bound run over a frozen snapshot | Claims cite evidence |
| `RecoveryAction` | Suggested non-money-moving action + preconditions | Never money-moving |
| `Passport` / `PassportVersion` | Derived, versioned, role-scoped proof projection | Private by default; append-only |
| `ConsentGrant` / `PassportShare` | Receiver/field permissions, expiring tokens | Revocable |

### 3.2 Twelve invariants (SOURCE I1–I12; implement in Go **and** PostgreSQL)

| ID | Invariant | Enforcement |
|---|---|---|
| I1 | A provider order belongs to exactly one tenant-owned attempt | `UNIQUE(environment, resource_type, provider_resource_id)` + FK |
| I2 | A mutation has one unique internal operation id and one stable provider request id | `UNIQUE(tenant_id, provider, environment, operation_request_id)` |
| I3 | Timeout/5xx = unknown until reconciled, never failure | Classifier + transition table + UI test |
| I4 | Verified provider events beat assumptions, user claims, LLM output; contradictions go to review | Reducer ranking |
| I5 | Capture only on the buyer-approved original order via the authorised integration | `checkout` precondition + row lock |
| I6 | `COMPLETED` capture ≠ bank `SETTLED` | `receiver_state` never derived from provider capture |
| I7 | A provider event id is processed at most once logically | Partial unique index on verified events |
| I8 | Sender/receiver payloads come from separate projections | `projector` + tests |
| I9 | Every AI conclusion has valid evidence refs or is withheld | `claims.Validator` |
| I10 | No agent tool has payment-write credentials | Import rules + process separation |
| I11 | Simulated events cannot elevate trust | Reducer: SIMULATED ⇒ tier E6, never touches `REAL_SANDBOX` cases |
| I12 | Closure is deterministic, evidence-based — not an LLM confidence | `cases.ResolutionPredicate` |

### 3.3 State vectors and enums
Use the canonical enums from doc 01 §12 verbatim. Go representation: typed string enums with `Valid()` methods and exhaustive-switch lint; mirrored as SQL `CHECK` constraints and OpenAPI enums (single generation source is OpenAPI; SQL CHECKs are asserted by a schema test).

---

## 4. API conventions and errors

### 4.1 Conventions
- Base path `/api/v1`; JSON (`application/json; charset=utf-8`); SSE for streams; money as decimal **strings**; timestamps RFC 3339 UTC.
- IDs in paths are UUIDs (cases) or opaque tokens (invites, shares).
- **Idempotency:** every state-mutating `POST` that a user can retry requires `Idempotency-Key` (UUID). Behaviour in §10.
- **Pagination:** cursor-based: `?limit=` (default 25, max 100) & `?cursor=`; response `{ "items": [...], "next_cursor": "..." }`.
- **Versioned reads:** case/passport responses carry `version` (state_version / passport_version) so the UI can invalidate on change.
- **Correlation:** every response has `X-Correlation-Id`; echoed in error bodies.
- **Uniform denial:** non-participants get `404 NOT_FOUND` (no existence disclosure). Authenticated participants lacking permission for a *field* get a `restricted` marker rather than an error.

### 4.2 Error envelope (SOURCE M15 §15.7)
```json
{
  "error": {
    "code": "CASE_NOT_FOUND",
    "message": "We couldn't find that case.",
    "correlation_id": "6f13fbee-57d0-41ee-8c5f-793bfb4f0021",
    "details": { "field": "case_id" }
  }
}
```
`details` carries **safe fields only**. Never include credentials, raw webhook bodies, stack traces, or another user's data.

### 4.3 Error code catalogue

| HTTP | `code` | Meaning |
|---|---|---|
| 400 | `VALIDATION_FAILED` | Schema/field validation (details.fields[]) |
| 401 | `UNAUTHENTICATED` | No/expired session |
| 403 | `FORBIDDEN_ACTION` | Authenticated participant, action not permitted (e.g. receiver calling capture) |
| 404 | `NOT_FOUND` | Missing **or** not visible |
| 409 | `ACTIVE_ATTEMPT_EXISTS` | A prior attempt is not terminal |
| 409 | `OUTCOME_UNKNOWN_RECONCILE_FIRST` | Blind retry blocked (I3) |
| 409 | `STATE_VERSION_CONFLICT` | Optimistic concurrency |
| 409 | `INVALID_TRANSITION` | Forbidden transition (§15) |
| 412 | `PRECONDITION_FAILED` | e.g. order not approved |
| 422 | `IDEMPOTENCY_KEY_REUSED` | Same key, different request hash |
| 422 | `AMOUNT_MISMATCH` | Provider amount/currency ≠ case |
| 429 | `RATE_LIMITED` | With `Retry-After` |
| 502 | `PROVIDER_UNAVAILABLE` | PayPal unreachable (outcome may be unknown → see body `outcome` field) |
| 503 | `SERVICE_UNAVAILABLE` | DB/readiness failure |
| 202 | *(not an error)* | `outcome: "UNKNOWN"` body for lost responses (see S2) |

**Critical:** a lost provider response on a mutation returns **`202 Accepted`** with `{ "outcome": "UNKNOWN", "safe_action": "DO_NOT_REPEAT_CHECKOUT" }`, never 5xx-as-failure semantics to the UI.

---

## 5. Middleware chain

Order (outermost first):

1. `RequestID` — generate/propagate `X-Correlation-Id`; attach to `slog` and OTel.
2. `Recoverer` — panic → 500 envelope (no stack leak); alert.
3. `SecurityHeaders` — `Strict-Transport-Security`, `X-Content-Type-Options: nosniff`, `Referrer-Policy: no-referrer`, `Permissions-Policy`, CSP for API docs routes.
4. `CORS` — allowlist from `PUBLIC_BASE_URL`/`WEB_ORIGIN`; credentials true; no wildcard.
5. `BodyLimit` — default 1 MiB; webhook route separate limit (256 KiB, **VERIFY** PayPal payload sizes).
6. `RateLimit` — token bucket per IP and per session (§10).
7. `Logging` — structured access log via allowlist (no bodies, no cookies, no auth headers).
8. `Metrics/Tracing`.
9. `Auth` (session) — except `/webhooks/paypal`, `/healthz`, `/readyz`, `/auth/login`, public share token route.
10. `CSRF` — double-submit token on cookie-authenticated unsafe methods.
11. `TenantContext` — resolves `tenant_id` from **verified membership/participation** (never from a client-supplied value) and opens the request transaction with `SET LOCAL app.tenant_id`.
12. `Idempotency` — for annotated routes.
13. `Handler` → calls `policy.Authorize(...)` explicitly (middleware only authenticates; **authorisation is explicit per handler** so nothing is forgotten silently — a route without an `Authorize` call fails an architecture test).

---

## 6. Validation

- Request bodies validated against the generated OpenAPI schema (strict: unknown fields rejected) **and** a domain validator.
- Amount: positive, decimal string `^\d{1,14}(\.\d{1,6})?$`, scale ≤ ISO-4217 exponent for the currency (`money` package table); reject `NaN`, exponent forms.
- Currency: ISO-4217 allowlist (`USD` at minimum for MVP).
- Country: ISO-3166-1 alpha-2.
- Strings: length caps, Unicode normalisation (NFC), control-character rejection; free-text (memo, notes) flagged `untrusted=true` in the evidence layer (prompt-injection surface — doc 07).
- IDs: UUID parse; opaque tokens: length + charset bounds.
- **Server-authoritative fields:** amount, currency, payee/merchant identity come from the canonical case/invoice — **never** from client request bodies for checkout/capture.

---

## 7. Auth and sessions

*(ADR-013; production IdP per doc 08.)*

- **MVP:** server-side sessions. `POST /auth/login` accepts seeded demo identities (`SESSION_SECRET`-signed opaque session id; row in `sessions`). Cookie: `HttpOnly; Secure; SameSite=Lax; Path=/; Max-Age` short (e.g. 8 h demo, idle timeout 30 min — values **DESIGN DECISION**, tune).
- Session rows: `id (random 256-bit, stored hashed), user_id, tenant_id (active), created_at, last_seen_at, expires_at, ip_hash, ua_hash, revoked_at`.
- **CSRF:** double-submit cookie + `X-CSRF-Token` header on unsafe methods.
- **Login throttling:** per-identity + per-IP; constant-time comparison; generic failure messages.
- **MFA:** required for `merchant_admin`/`merchant_analyst` in any non-demo deployment (doc 08); demo accounts documented as demo-only.
- **OIDC adapter interface** (`auth.IdentityProvider`) prepared; not implemented in MVP.
- **Receiver identity:** invitation token alone is **not** authentication. The receiver must be signed in as the **intended identity** (email match against the invite) to accept (S9).

---

## 8. RBAC, ABAC and case access

### 8.1 Decision inputs
`policy.Authorize(ctx, Request{Actor, Tenant, Case, Action, Field?, Purpose, Environment, ConsentID?})` → `Decision{Allow, Reason, RestrictedFields}`. Inputs per SOURCE M16 §16.2: actor, tenant, case, requested field, purpose, consent, action sensitivity, environment.

### 8.2 Role × capability matrix (default deny)

| Capability | sender | receiver | merchant_analyst | merchant_admin | platform_operator | demo_operator |
|---|:-:|:-:|:-:|:-:|:-:|:-:|
| Create case | ✔ | – | ✔ | ✔ | – | ✔ (sim) |
| Read case (own scope) | ✔ | ✔ (consent) | ✔ (tenant) | ✔ (tenant) | JIT only | sim only |
| Create order / capture | ✔ (own case) | – | – | – | – | sim only |
| Trigger reconcile (read-only) | ✔ | – | ✔ | ✔ | – | ✔ |
| Request investigation | ✔ | ✔ (own view) | ✔ | ✔ | – | ✔ |
| Read Sender Passport | ✔ | – | ✔ | ✔ | – | ✔ |
| Read Receiver Passport | – | ✔ | ✔ | ✔ | – | ✔ |
| Invite receiver / manage shares | ✔ (owner) | – | ✔ | ✔ | – | – |
| Ack action | designated actor | designated actor | ✔ | ✔ | – | – |
| Approve internal reconciliation | – | – | ✔ | ✔ | – | – |
| Ops queue / webhook health | – | – | ✔ | ✔ | JIT | – |
| Audit export | – | – | – | ✔ | JIT + logged | – |
| Run simulator / fault toggle | – | – | – | – | – | ✔ (never in non-sandbox) |

### 8.3 ABAC rules (all must hold)
1. `actor` is a member of the case's `tenant_id` **or** a verified, consented participant of the case.
2. Environment guard: `demo_operator` limited to `LOCAL_SIMULATION`/`REAL_SANDBOX`; **no** role may act on `PRODUCTION` (hard-disabled).
3. Field-level: receiver access limited to the **consent grant's field allowlist** and the receiver projection (§18).
4. Purpose limitation: operations such as audit export require `purpose` and are themselves audited.
5. **Platform operator** access is just-in-time (time-boxed grant row), logged separately.
6. Every deny returns uniform `404` to non-participants and writes an `audit_events` row with reason code (not the denied resource's content).

### 8.4 Defence in depth
PostgreSQL **RLS** on every tenant-scoped table (§11.4) with a non-bypass app role. A test enumerates all tables with a `tenant_id` column and asserts RLS is `FORCE`d with `SELECT/INSERT/UPDATE/DELETE` policies (SOURCE M11 §11.4).

---

## 9. PayPal adapter and OAuth lifecycle

> **VERIFY before coding:** exact request/response schemas, headers, error shapes, and endpoint semantics in the official PayPal docs (Orders v2, Payments v2, Webhooks, idempotency). Endpoint concepts below are from the M20 plan §12.3.

### 9.1 Interfaces (ADR-016)

```go
package paypal

// Reader: GET-only. The ONLY PayPal capability available to reconciliation.
type Reader interface {
    ShowOrder(ctx context.Context, env Environment, orderID string) (Order, RawEvidence, error)
    ShowCapture(ctx context.Context, env Environment, captureID string) (Capture, RawEvidence, error)
}

// WebhookVerifier: performs an HTTP POST to PayPal's verify endpoint, which has no financial side effect.
// Deliberately NOT part of Reader so that Reader can be proven GET-only. Injected only into webhook.Gateway.
type WebhookVerifier interface {
    VerifyWebhookSignature(ctx context.Context, in WebhookVerifyInput) (VerifyResult, error)
}

// Mutator: money-affecting. Constructed ONLY in cmd/api and injected ONLY into checkout.
type Mutator interface {
    CreateOrder(ctx context.Context, in CreateOrderInput) (Order, RawEvidence, error)   // requires RequestID
    CaptureOrder(ctx context.Context, in CaptureInput) (Capture, RawEvidence, error)    // requires RequestID
}
```
`WebhookVerifier.VerifyWebhookSignature` performs an HTTP POST but is read-only with respect to money; it lives in `verify.go`, is a separate interface from `Reader` (decision D-14), and is explicitly allow-listed in the architecture test as the only non-GET, non-Mutator provider call.

### 9.2 Endpoint map (concept; **VERIFY**)

| Operation | HTTP | Side effect | Interface |
|---|---|---|---|
| OAuth token | `POST /v1/oauth2/token` | none | internal |
| Create order | `POST /v2/checkout/orders` (`intent=CAPTURE`) | creates provider order | Mutator |
| Show order | `GET /v2/checkout/orders/{id}` | none | Reader |
| Capture order | `POST /v2/checkout/orders/{id}/capture` | may move sandbox money | Mutator |
| Show capture | `GET /v2/payments/captures/{id}` | none | Reader |
| Verify webhook | `POST /v1/notifications/verify-webhook-signature` | none | WebhookVerifier |
| Refund capture | `POST /v2/payments/captures/{id}/refund` | **excluded from MVP; no interface exposes it** | – |

Base host: `PAYPAL_API_BASE` must equal `https://api-m.sandbox.paypal.com`; config validator **refuses to start otherwise** (PRODUCTION hard-disabled).

### 9.3 OAuth/token lifecycle
- Client-credentials grant with `PAYPAL_CLIENT_ID`/`PAYPAL_CLIENT_SECRET` (server-only). Never in browser, logs, or evidence.
- In-memory token cache per process; refresh when `now >= expires_at − 60s`; **singleflight** to avoid thundering herd.
- On `401` from an API call: invalidate token, refresh once, retry the *same call only if it is a Reader call*. A Mutator call that gets 401 **before** any processing may be retried once with the **same** `PayPal-Request-Id`; any other ambiguous Mutator failure → `unknown` (never a new key).
- Token never logged; log only `token_age_seconds` and `token_refresh_count`.

### 9.4 Response classification (SOURCE M16 §16.3)
`classify(resp, err) → { success | explicit_decline | pending | unknown_missing_response }`

| Observation | Class |
|---|---|
| 2xx with `COMPLETED` capture | success |
| 2xx with `PENDING` | pending |
| 4xx with documented denial (e.g. instrument declined) | explicit_decline |
| context deadline / EOF / connection reset / 5xx / unparseable body after request sent | **unknown_missing_response** |
| 4xx validation error *before* processing (malformed request) | request_error (not a payment outcome) |

### 9.5 Request-id discipline
- `PayPal-Request-Id` = `pb-{operation}-{operation_request_id}`; **different** per operation type; **stable** across retries of the same operation; persisted **before** the call.
- Provider idempotency duration/semantics vary (**VERIFY**), therefore PayBridge also keeps its own unique constraint and attempt lock.

### 9.6 Raw evidence capture
Every provider response is stored as a **redacted JSON extract + SHA-256 of original bytes**; the original bytes (if retained) go to the encrypted evidence store (AES-GCM, `ENCRYPTION_KEY`), referenced by `encrypted_payload_ref`.

### 9.7 Outbound safety
Allowlisted hosts only; no redirects followed to other hosts; TLS verification on; response size caps; per-call timeouts (§10).

---

## 10. Idempotency, timeouts, rate limits

### 10.1 Application idempotency (`Idempotency-Key`)
Table `idempotency_keys(tenant_id, actor_id, route, key, request_hash, status, response_code, response_body, created_at, expires_at)`; `UNIQUE(tenant_id, actor_id, route, key)`.
- Same key + same hash → return the **stored response** (replay).
- Same key + different hash → `422 IDEMPOTENCY_KEY_REUSED`.
- In-flight duplicate → `409` with `Retry-After` (never executes twice).
- TTL e.g. 24 h (**DESIGN DECISION**).
This is **in addition to** the attempt-level lock and PayPal request ids — three layers.

### 10.2 Timeouts (all configurable; defaults DESIGN DECISION)

| Call | Connect | Total | Retries |
|---|---|---|---|
| PayPal Reader GET | 3 s | 10 s | Backoff+jitter, max attempts 5, max elapsed 10 min, honour `Retry-After` |
| PayPal create order | 3 s | 15 s | **None** under a new key; a single replay with same key only for a provably pre-send failure |
| PayPal capture | 3 s | 25 s | Same as above — **never** blind |
| Webhook verify | 3 s | 10 s | Bounded; if verification itself is unavailable → return **5xx** so PayPal retries (do not ack unverified) |
| LLM call | 3 s | 30 s per agent | 1 retry on transient; budget-bounded |
| DB statement | – | 5 s default; 30 s for migrations/maintenance | – |
| HTTP server | read header 5 s, read 15 s, write 30 s, idle 60 s; SSE exempt from write timeout with heartbeat | – | – |

### 10.3 Rate limits (token bucket)
- Per IP (unauthenticated): e.g. 60/min; per session: 300/min; login: 10/min/IP + 5/min/identity; mutations on checkout: 10/min/case; webhook route: separate generous-but-bounded bucket (provider retries must not be starved).
- Storage: Redis if configured, else in-process. **Rate-limit state is never financial truth.** Values are **DESIGN DECISION**, tune from measurement.
- Responses: `429` + `Retry-After`.

---

## 11. PostgreSQL schema

> Migrations are forward-only goose files. Starter based on SOURCE annex §9, **corrected** for conflicts C-01/C-02/C-06 and extended. Types: `UUID`, `NUMERIC(20,6)`, `TIMESTAMPTZ`, `JSONB`. Execute under a **restricted role** `paybridge_app` (no `BYPASSRLS`, no `UPDATE/DELETE` on append-only tables).

### 11.1 Migration plan

| File | Contents |
|---|---|
| `0001_extensions_roles.sql` | extensions (`pgcrypto`, `citext`), roles, default privileges |
| `0002_identity_tenancy.sql` | tenants, users, memberships, sessions |
| `0003_cases_attempts.sql` | payment_cases, case_participants, payment_attempts, provider_resources, idempotency_keys |
| `0004_events_evidence.sql` | payment_events, evidence_items, evidence_edges, webhook_inbox |
| `0005_outbox_jobs.sql` | outbox (jobs with leases), dead_letters |
| `0006_ai.sql` | investigations, agent_runs, diagnosis_claims, claim_evidence, recovery_actions, action_acks |
| `0007_passport_consent.sql` | passport_versions, passport_shares, invitations, consent_grants |
| `0008_comms_audit.sql` | notifications, audit_events (hash chain), retention_jobs |
| `0009_simulation_ops.sql` | simulation_runs, route_registry, provider_connections, provider_status_samples |
| `0010_rls.sql` | RLS enable/force + policies + privilege revokes + append-only triggers |

### 11.2 Core DDL (excerpt — implement fully in migrations)

```sql
-- 0002
CREATE TABLE tenants (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email CITEXT UNIQUE NOT NULL,
  display_name TEXT NOT NULL,
  password_hash TEXT,                     -- demo login only; argon2id
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE tenant_memberships (          -- resolves conflict C-01: six roles
  tenant_id UUID NOT NULL REFERENCES tenants(id),
  user_id   UUID NOT NULL REFERENCES users(id),
  role TEXT NOT NULL CHECK (role IN
    ('sender','receiver','merchant_analyst','merchant_admin','platform_operator','demo_operator')),
  PRIMARY KEY (tenant_id, user_id, role)
);
CREATE TABLE sessions (
  id_hash BYTEA PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(id),
  tenant_id UUID REFERENCES tenants(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  last_seen_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  expires_at TIMESTAMPTZ NOT NULL,
  revoked_at TIMESTAMPTZ,
  ip_hash BYTEA, ua_hash BYTEA
);

-- 0003
CREATE TABLE payment_cases (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL REFERENCES tenants(id),
  external_reference TEXT,                       -- invoice ref (untrusted text)
  sender_user_id UUID NOT NULL REFERENCES users(id),
  receiver_user_id UUID REFERENCES users(id),
  sender_country CHAR(2) NOT NULL,
  receiver_country CHAR(2) NOT NULL,
  payment_purpose TEXT NOT NULL,
  amount NUMERIC(20,6) NOT NULL CHECK (amount > 0),
  currency CHAR(3) NOT NULL,
  environment TEXT NOT NULL CHECK (environment IN
    ('REAL_SANDBOX','PAYPAL_SIMULATOR','LOCAL_SIMULATION','PRODUCTION')),
  -- state vector (written ONLY by the Reducer)
  provider_state TEXT NOT NULL DEFAULT 'NOT_CREATED' CHECK (provider_state IN
    ('NOT_CREATED','CREATED','APPROVED','AUTHORISED','CAPTURE_PENDING','CAPTURE_COMPLETED',
     'CAPTURE_DENIED','REFUNDED','PARTIALLY_REFUNDED','UNKNOWN')),
  last_confirmed_provider_state TEXT NOT NULL DEFAULT 'NOT_CREATED',
  provider_state_evidence_id UUID,
  provider_state_observed_at TIMESTAMPTZ,
  application_state TEXT NOT NULL DEFAULT 'SYNCHRONISED' CHECK (application_state IN
    ('SYNCHRONISED','WAITING_FOR_WEBHOOK','API_TIMEOUT','RECONCILING','EVENT_GAP','INCONSISTENT')),
  receiver_state TEXT NOT NULL DEFAULT 'NOT_OBSERVABLE' CHECK (receiver_state IN
    ('NOT_APPLICABLE','NOT_OBSERVABLE','REPORTED_PENDING','VERIFIED_AVAILABLE',
     'WITHDRAWAL_IN_PROGRESS','VERIFIED_BANK_CREDITED')),
  case_status TEXT NOT NULL DEFAULT 'OPEN' CHECK (case_status IN
    ('OPEN','OBSERVING','INVESTIGATING','GUIDANCE_READY','AWAITING_USER_ACTION',
     'VERIFYING_RESOLUTION','RESOLVED','NEEDS_HUMAN_REVIEW','PAUSED_INSUFFICIENT_EVIDENCE',
     'CLOSED_UNRESOLVED','ARCHIVED','REOPENED')),
  state_version BIGINT NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
  -- Guards that cannot be expressed as CHECKs (simulated evidence never moves a REAL_SANDBOX case;
  -- receiver_state never VERIFIED_* without receiver-connector evidence) are TRIGGERS in 0010 + Reducer rules.
);
CREATE INDEX ix_cases_tenant_status ON payment_cases(tenant_id, case_status, updated_at DESC);

CREATE TABLE payment_attempts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL REFERENCES tenants(id),
  case_id UUID NOT NULL REFERENCES payment_cases(id),
  attempt_no INT NOT NULL CHECK (attempt_no > 0),
  provider TEXT NOT NULL DEFAULT 'PAYPAL',
  environment TEXT NOT NULL,
  order_operation_request_id TEXT NOT NULL,
  capture_operation_request_id TEXT,
  status TEXT NOT NULL DEFAULT 'ORDER_CREATING' CHECK (status IN
    ('ORDER_CREATING','ORDER_CREATED','ORDER_UNKNOWN','APPROVAL_PENDING','APPROVED','CANCELLED',
     'CAPTURE_REQUESTED','CAPTURE_COMPLETED','CAPTURE_PENDING','CAPTURE_DENIED','CAPTURE_UNKNOWN',
     'REVIEW_REQUIRED','APPLICATION_RECONCILED','APPLICATION_MISMATCH')),
  terminal BOOLEAN GENERATED ALWAYS AS (status IN
    ('CANCELLED','CAPTURE_DENIED','APPLICATION_RECONCILED')) STORED,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (tenant_id, case_id, attempt_no),
  UNIQUE (tenant_id, provider, environment, order_operation_request_id)
);
-- At most ONE non-terminal attempt per case (blind-retry guard, I3/I5):
CREATE UNIQUE INDEX uq_one_active_attempt ON payment_attempts(case_id) WHERE terminal = FALSE;

CREATE TABLE provider_resources (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL REFERENCES tenants(id),
  attempt_id UUID NOT NULL REFERENCES payment_attempts(id),
  resource_type TEXT NOT NULL CHECK (resource_type IN ('ORDER','AUTHORIZATION','CAPTURE','REFUND')),
  provider_resource_id TEXT NOT NULL,
  environment TEXT NOT NULL,
  provider_status TEXT,
  last_verified_at TIMESTAMPTZ,
  UNIQUE (environment, resource_type, provider_resource_id)       -- I1
);
CREATE TABLE idempotency_keys (
  tenant_id UUID NOT NULL, actor_id UUID NOT NULL, route TEXT NOT NULL, key TEXT NOT NULL,
  request_hash BYTEA NOT NULL, status TEXT NOT NULL CHECK (status IN ('IN_FLIGHT','DONE')),
  response_code INT, response_body JSONB,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(), expires_at TIMESTAMPTZ NOT NULL,
  PRIMARY KEY (tenant_id, actor_id, route, key)
);

-- 0004
CREATE TABLE evidence_items (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL REFERENCES tenants(id),
  case_id UUID NOT NULL REFERENCES payment_cases(id),
  kind TEXT NOT NULL CHECK (kind IN
    ('PROVIDER_GET_ORDER','PROVIDER_GET_CAPTURE','PROVIDER_WEBHOOK','APP_TRACE','MERCHANT_DB_SNAPSHOT',
     'PUBLIC_STATUS','USER_REPORT','SIMULATION','ROUTE_RULE','FX_QUOTE')),
  tier SMALLINT NOT NULL CHECK (tier BETWEEN 1 AND 6),                 -- E1..E6
  source_trust TEXT NOT NULL CHECK (source_trust IN
    ('PROVIDER_VERIFIED','LOCAL_OBSERVED','PUBLIC_SOURCE','USER_REPORTED','SIMULATED')),
  environment TEXT NOT NULL,
  verified BOOLEAN NOT NULL DEFAULT FALSE,
  untrusted_text BOOLEAN NOT NULL DEFAULT FALSE,
  occurred_at TIMESTAMPTZ,
  observed_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  redacted_extract JSONB NOT NULL DEFAULT '{}'::jsonb,
  payload_sha256 CHAR(64),
  encrypted_payload_ref TEXT,
  access_class TEXT NOT NULL DEFAULT 'TENANT' CHECK (access_class IN ('TENANT','SENDER','RECEIVER','OPS_ONLY')),
  CHECK (source_trust <> 'SIMULATED' OR tier = 6)                      -- I11
);
CREATE INDEX ix_evidence_case ON evidence_items(case_id, observed_at);

CREATE TABLE evidence_edges (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL, case_id UUID NOT NULL REFERENCES payment_cases(id),
  from_evidence_id UUID NOT NULL REFERENCES evidence_items(id),
  to_evidence_id   UUID NOT NULL REFERENCES evidence_items(id),
  relation TEXT NOT NULL CHECK (relation IN
    ('CORROBORATES','CONTRADICTS','PRECEDES','DERIVED_FROM','REQUIRES','RESOLVES','INDEPENDENT_OF')),
  rule_id TEXT NOT NULL, rule_version INT NOT NULL,
  origin TEXT NOT NULL CHECK (origin IN ('DETERMINISTIC','LLM_HYPOTHESIS')),   -- LLM edges stay hypotheses
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE payment_events (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL REFERENCES tenants(id),
  case_id UUID NOT NULL REFERENCES payment_cases(id),
  evidence_id UUID NOT NULL REFERENCES evidence_items(id),
  provider TEXT, provider_event_id TEXT,
  environment TEXT NOT NULL,
  event_type TEXT NOT NULL,
  status_observed TEXT,
  source_trust TEXT NOT NULL,
  verified BOOLEAN NOT NULL DEFAULT FALSE,
  provider_occurred_at TIMESTAMPTZ,
  received_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  schema_version INT NOT NULL DEFAULT 1
);
CREATE UNIQUE INDEX uq_verified_provider_events
  ON payment_events(provider, environment, provider_event_id)
  WHERE provider_event_id IS NOT NULL AND verified = TRUE;           -- I7
CREATE INDEX ix_case_events ON payment_events(case_id, received_at, id);

CREATE TABLE webhook_inbox (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  environment TEXT NOT NULL,
  provider_event_id TEXT,
  transmission_id TEXT,
  event_type TEXT,
  signature_verified BOOLEAN NOT NULL DEFAULT FALSE,
  payload_sha256 CHAR(64) NOT NULL,
  encrypted_payload_ref TEXT,
  processing_status TEXT NOT NULL DEFAULT 'RECEIVED' CHECK (processing_status IN
    ('RECEIVED','REJECTED','CLAIMED','UNCORRELATED','PROCESSED','DUPLICATE','DEAD')),
  duplicate_count INT NOT NULL DEFAULT 0,
  attempts INT NOT NULL DEFAULT 0,
  lease_until TIMESTAMPTZ, leased_by TEXT,
  case_id UUID,
  received_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  processed_at TIMESTAMPTZ,
  UNIQUE (environment, provider_event_id)
);
CREATE INDEX ix_inbox_ready ON webhook_inbox(processing_status, lease_until, received_at);

-- 0005
CREATE TABLE outbox (
  id BIGSERIAL PRIMARY KEY,
  tenant_id UUID, case_id UUID,
  kind TEXT NOT NULL CHECK (kind IN
    ('RECONCILE','INVESTIGATE','NOTIFY','PASSPORT_REFRESH','RETENTION','SSE_PUBLISH')),
  payload JSONB NOT NULL,
  dedupe_key TEXT,
  status TEXT NOT NULL DEFAULT 'READY' CHECK (status IN ('READY','LEASED','DONE','DEAD')),
  attempts INT NOT NULL DEFAULT 0, max_attempts INT NOT NULL DEFAULT 8,
  run_after TIMESTAMPTZ NOT NULL DEFAULT now(),
  lease_until TIMESTAMPTZ, leased_by TEXT, last_error TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX uq_outbox_dedupe ON outbox(kind, dedupe_key) WHERE dedupe_key IS NOT NULL AND status IN ('READY','LEASED');
CREATE INDEX ix_outbox_ready ON outbox(status, run_after);

-- 0006 (excerpt)
CREATE TABLE investigations (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL, case_id UUID NOT NULL REFERENCES payment_cases(id),
  trigger TEXT NOT NULL, snapshot_version BIGINT NOT NULL,
  status TEXT NOT NULL CHECK (status IN ('QUEUED','RUNNING','COMPLETED','FAILED_FALLBACK','STALE_DISCARDED','NEEDS_HUMAN_REVIEW')),
  llm_available BOOLEAN, created_at TIMESTAMPTZ NOT NULL DEFAULT now(), completed_at TIMESTAMPTZ
);
CREATE TABLE agent_runs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  investigation_id UUID NOT NULL REFERENCES investigations(id),
  agent_id TEXT NOT NULL CHECK (agent_id IN ('A01','A02','A03','A04','A05','A06','A07','A08','A09','A10')),
  model_id TEXT, prompt_sha256 CHAR(64), input_evidence_ids UUID[] NOT NULL DEFAULT '{}',
  output JSONB, tokens_in INT, tokens_out INT, latency_ms INT,
  rejection_reasons TEXT[] NOT NULL DEFAULT '{}', created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE diagnosis_claims (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL, case_id UUID NOT NULL,
  investigation_id UUID NOT NULL REFERENCES investigations(id),
  kind TEXT NOT NULL CHECK (kind IN ('confirmed','inference','unknown')),
  code TEXT NOT NULL, statement TEXT NOT NULL, inference_rule TEXT,
  viewer_scope TEXT NOT NULL CHECK (viewer_scope IN ('SENDER','RECEIVER','OPS','ALL')),
  validated BOOLEAN NOT NULL DEFAULT FALSE
);
CREATE TABLE claim_evidence (
  claim_id UUID NOT NULL REFERENCES diagnosis_claims(id),
  evidence_id UUID NOT NULL REFERENCES evidence_items(id),
  PRIMARY KEY (claim_id, evidence_id)
);

-- 0007
CREATE TABLE passport_versions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL, case_id UUID NOT NULL REFERENCES payment_cases(id),
  viewer_role TEXT NOT NULL CHECK (viewer_role IN ('sender','receiver','ops')),
  version INT NOT NULL CHECK (version > 0),
  snapshot_version BIGINT NOT NULL,
  content JSONB NOT NULL,
  content_sha256 CHAR(64) NOT NULL,
  prev_sha256 CHAR(64),
  evidence_ids UUID[] NOT NULL,
  policy_version INT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (case_id, viewer_role, version)
);
CREATE TABLE invitations (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL, case_id UUID NOT NULL REFERENCES payment_cases(id),
  invitee_email CITEXT NOT NULL, token_hash BYTEA NOT NULL UNIQUE,
  proposed_fields TEXT[] NOT NULL, expires_at TIMESTAMPTZ NOT NULL,
  status TEXT NOT NULL CHECK (status IN ('PENDING','ACCEPTED','DECLINED','EXPIRED','REVOKED')),
  created_by UUID NOT NULL REFERENCES users(id), created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE consent_grants (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL, case_id UUID NOT NULL REFERENCES payment_cases(id),
  grantee_user_id UUID NOT NULL REFERENCES users(id),
  field_allowlist TEXT[] NOT NULL, notification_allowed BOOLEAN NOT NULL DEFAULT FALSE,
  granted_at TIMESTAMPTZ NOT NULL DEFAULT now(), revoked_at TIMESTAMPTZ,
  version INT NOT NULL DEFAULT 1
);
CREATE TABLE passport_shares (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL, case_id UUID NOT NULL REFERENCES payment_cases(id),
  viewer_role TEXT NOT NULL, token_hash BYTEA NOT NULL UNIQUE,
  expires_at TIMESTAMPTZ NOT NULL, revoked_at TIMESTAMPTZ,
  created_by UUID NOT NULL, created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- 0008
CREATE TABLE audit_events (
  seq BIGSERIAL PRIMARY KEY,
  tenant_id UUID, case_id UUID, actor_id UUID, actor_role TEXT,
  operation TEXT NOT NULL, environment TEXT, reason TEXT,
  prev_version BIGINT, new_version BIGINT, evidence_ids UUID[],
  request_id TEXT, policy_version INT,
  occurred_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  prev_hash CHAR(64), row_hash CHAR(64) NOT NULL        -- sha256(prev_hash || canonical row)
);
```

> **Implementation note (important):** two guards live in `0010` as **triggers**, not CHECKs: (1) reject any `UPDATE` of `receiver_state` to `VERIFIED_AVAILABLE`/`WITHDRAWAL_IN_PROGRESS`/`VERIFIED_BANK_CREDITED` unless the same transaction inserted receiver-connector evidence of tier ≤ 2 (no such connector exists in the MVP ⇒ never satisfiable); (2) reject any state-vector update on a `REAL_SANDBOX` case whose causing evidence has `source_trust='SIMULATED'`. If a new `kind` value for receiver-connector evidence is needed, extend the `evidence_items.kind` CHECK in a later migration.

### 11.3 Indexes and constraints (summary)
- Partial unique `uq_one_active_attempt` (blind-retry guard).
- Partial unique `uq_verified_provider_events` (I7) and `webhook_inbox UNIQUE(environment, provider_event_id)`.
- `provider_resources UNIQUE(environment, resource_type, provider_resource_id)` (I1).
- `evidence_items CHECK (source_trust <> 'SIMULATED' OR tier = 6)` (I11).
- Foreign keys everywhere; `ON DELETE RESTRICT` default (retention jobs delete explicitly, in order).

### 11.4 Row-level security (SOURCE M11 §11.4, extended)

```sql
-- applied to EVERY table that has tenant_id
ALTER TABLE payment_cases ENABLE ROW LEVEL SECURITY;
ALTER TABLE payment_cases FORCE  ROW LEVEL SECURITY;
CREATE POLICY pc_select ON payment_cases FOR SELECT
  USING (tenant_id = nullif(current_setting('app.tenant_id', true), '')::uuid);
CREATE POLICY pc_insert ON payment_cases FOR INSERT
  WITH CHECK (tenant_id = nullif(current_setting('app.tenant_id', true), '')::uuid);
CREATE POLICY pc_update ON payment_cases FOR UPDATE
  USING (tenant_id = nullif(current_setting('app.tenant_id', true), '')::uuid)
  WITH CHECK (tenant_id = nullif(current_setting('app.tenant_id', true), '')::uuid);
-- tenant context: `SET LOCAL app.tenant_id = '<uuid>'` inside each authenticated transaction.
-- the Worker sets the tenant from the job row; the webhook gateway uses a dedicated narrow role limited to webhook_inbox INSERT/UPDATE-duplicate.
```
Append-only tables (`payment_events`, `evidence_items`, `audit_events`, `passport_versions`): `REVOKE UPDATE, DELETE` from `paybridge_app` plus a `BEFORE UPDATE OR DELETE` trigger raising an exception (retention jobs run under a separate `paybridge_retention` role).

---

## 12. Webhook pipeline

*(SOURCE M12 §12.5, M13 §13.5.)*

### 12.1 Receive path (`POST /api/v1/webhooks/paypal`)
```text
1. TLS terminated upstream; handler reads the RAW body bytes (do not JSON-decode first).
2. Enforce size limit; capture headers: PAYPAL-TRANSMISSION-ID, PAYPAL-TRANSMISSION-TIME,
   PAYPAL-TRANSMISSION-SIG, PAYPAL-CERT-URL, PAYPAL-AUTH-ALGO  (names VERIFY in docs)
3. Classify: real (REAL_SANDBOX) vs simulator-origin. Simulator/mock events go to an ISOLATED route
   and environment PAYPAL_SIMULATOR — they cannot be verified through the real postback.
4. Verify via PayPal verify-webhook-signature using PAYPAL_WEBHOOK_ID and the raw event.
   - verification FAILS → insert inbox row REJECTED (hash only), return 400; metric signature_failures++.
   - verification UNAVAILABLE (network) → return 503 so PayPal retries. NEVER ack unverified data.
5. Compute payload_sha256; encrypt + store raw bytes (optional, retention-controlled).
6. INSERT inbox row (UNIQUE environment+provider_event_id). On conflict: duplicate_count++, status stays; return 200.
7. Return 200 ONLY after the insert committed. If the DB insert fails → 5xx.
8. No case mutation happens in this handler.
```

### 12.2 Replay protection (DESIGN DECISION + **VERIFY**)
- Primary: dedupe on `provider_event_id`; record `transmission_id` for diagnostics.
- Timestamp window: **do not** reject genuine provider retries by age (PayPal retries can arrive long after the first attempt — **VERIFY** retry policy). Instead protect against *state regression* with the monotonic Reducer (a replayed old signed webhook is harmless: it appends evidence but cannot regress state — SOURCE §13.6).
- Optional: log and alert when `PAYPAL-TRANSMISSION-TIME` is older than a configurable threshold, without rejecting.

### 12.3 Cross-tenant safety
A verified signature is **not** tenant authorisation. The processor must confirm `resource.id` (order/capture) maps to a `provider_resources` row owned by a case in a tenant bound to the configured merchant (SOURCE §12.5). Unknown resource ⇒ `UNCORRELATED` (bounded re-check, then `DEAD` with alert).

### 12.4 Processing (Inbox processor)
```text
claim row WHERE processing_status IN ('RECEIVED','UNCORRELATED') AND (lease_until IS NULL OR lease_until < now())
  ORDER BY received_at FOR UPDATE SKIP LOCKED LIMIT n   -- set lease_until = now()+lease, leased_by
parse → normalise → correlate resource → case
BEGIN tx
  insert evidence_items (E2, PROVIDER_VERIFIED, verified=true)
  insert payment_events (unique verified provider_event_id; ON CONFLICT DO NOTHING → duplicate path)
  advisory lock on case; Reducer.Apply(case, evidence) with optimistic state_version
  insert outbox rows (PASSPORT_REFRESH, SSE_PUBLISH, NOTIFY, maybe INVESTIGATE)
  insert audit_events
  mark inbox PROCESSED
COMMIT
if stale/conflicting webhook → enqueue RECONCILE (read-through) before trusting it
```

---

## 13. Inbox/outbox workers, retry, dead-letter

### 13.1 Worker loops
`outbox.Relay` polls `outbox WHERE status='READY' AND run_after<=now()` using `FOR UPDATE SKIP LOCKED`, leases (`lease_until`), dispatches by `kind` to a handler registry. Every handler is **idempotent** (keyed by `dedupe_key`/case snapshot).

### 13.2 Retry and backoff
`delay = min(cap, base · 2^attempt) ± jitter(0–30 %)`; defaults: base 5 s, cap 15 min, max_attempts 8 (**DESIGN DECISION**). Read-only provider calls honour `Retry-After`. **Never** reschedule a financial mutation under a new operation id.

### 13.3 Dead-letter
After `max_attempts`: `status='DEAD'`, insert `dead_letters` row (job, last_error, case), emit metric + alert, show in Ops → Webhooks/Jobs view, and open `NEEDS_HUMAN_REVIEW` on the case where relevant. Dead jobs can be **manually re-queued** by `merchant_admin` (audited), never auto.

### 13.4 Queue depth & health
Gauges: `outbox_ready`, `outbox_oldest_age_seconds`, `inbox_unprocessed`, `unknown_state_backlog_age_seconds`, `dead_letters_total`.

### 13.5 Scheduled jobs
Reconcile sweeper (cases in `API_TIMEOUT/RECONCILING/EVENT_GAP` older than threshold) · stale-lease reaper · retention/deletion · provider status sampler (A02 advisory) · idempotency key expiry · session expiry.

---

## 14. Route inventory

Auth legend: **S** session cookie; **P** provider signature; **T** share token; **—** none. Idempotency: **K** = `Idempotency-Key` required; **N** = natural (read/idempotent); **U** = unique constraint-based. All routes under `/api/v1` unless noted. Errors reference §4.3 plus `401/404/429/503` implicitly.

| ID | Method | Path | Auth / roles | Request | Response | Specific errors | Idem. |
|---|---|---|---|---|---|---|---|
| R01 | POST | `/auth/login` | — | `{email, password}` | `200` session cookie + `{user, memberships}` | 401, 429 | N |
| R02 | POST | `/auth/logout` | S | — | `204` | – | N |
| R03 | GET | `/auth/me` | S | — | `{user, memberships, active_tenant}` | – | N |
| R04 | GET | `/auth/csrf` | S | — | `{token}` | – | N |
| R05 | POST | `/cases` | S sender/analyst/admin | `{external_reference, receiver_email?, receiver_country, sender_country, purpose, amount, currency}` | `201 Case` | `VALIDATION_FAILED` | K |
| R06 | GET | `/cases` | S any (scoped) | `?status&role&env&cursor&limit` | `{items:[CaseSummary], next_cursor}` | – | N |
| R07 | GET | `/cases/{id}` | S participant | — | `Case` (role-projected) | `NOT_FOUND` | N |
| R08 | POST | `/cases/{id}/readiness` | S sender/analyst | `{}` | `Readiness{support_result, requirements[], rate_quote, provenance}` | – | K |
| R09 | POST | `/cases/{id}/orders` | S sender (own) | `{}` (amount from case) | `201 {attempt, order_approval_context}` or `202 {outcome:"UNKNOWN"}` | `ACTIVE_ATTEMPT_EXISTS`, `OUTCOME_UNKNOWN_RECONCILE_FIRST`, `PROVIDER_UNAVAILABLE` | K |
| R10 | POST | `/cases/{id}/capture` | S sender (own) | `{attempt_id}` | `200 Case` or `202 {outcome:"UNKNOWN", safe_action}` | `PRECONDITION_FAILED`, `AMOUNT_MISMATCH`, `INVALID_TRANSITION` | K |
| R11 | GET | `/cases/{id}/timeline` | S participant | `?since` | `{events:[TimelineEvent]}` role-filtered | – | N |
| R12 | POST | `/cases/{id}/reconcile` | S sender/analyst/admin | `{}` | `202 {job_id}` | `RATE_LIMITED` | K |
| R13 | POST | `/cases/{id}/investigations` | S participant | `{trigger?}` | `202 {investigation_id}` | – | K |
| R14 | GET | `/cases/{id}/investigations/{run}` | S role-scoped | — | `Investigation{claims, agents[], status}` | – | N |
| R15 | GET | `/cases/{id}/passport` | S participant | `?version&viewer?` | `Passport` (viewer-filtered) | – | N |
| R16 | GET | `/cases/{id}/passport/versions` | S participant | — | `{versions:[{version, created_at, snapshot_version, sha256}]}` | – | N |
| R17 | POST | `/cases/{id}/passport/shares` | S owner/analyst | `{viewer_role, expires_in}` | `201 {share_id, url_token}` | `FORBIDDEN_ACTION` | K |
| R18 | DELETE | `/shares/{id}` | S share owner | — | `204` | – | N |
| R19 | GET | `/shares/{token}/passport` | T | — | `Passport` (share-scoped, no-index headers) | `NOT_FOUND` | N |
| R20 | POST | `/cases/{id}/receiver-invite` | S sender/analyst | `{email, fields[]}` | `201 {invite_id}` | `VALIDATION_FAILED` | K |
| R21 | POST | `/invites/{token}/accept` | S intended receiver | `{fields_accepted[], notifications}` | `200 {consent_id}` | `FORBIDDEN_ACTION` | K |
| R22 | POST | `/invites/{token}/decline` | S intended receiver | — | `204` | – | N |
| R23 | GET | `/cases/{id}/consents` | S owner/receiver | — | `{consents:[...]}` | – | N |
| R24 | POST | `/consents/{id}/revoke` | S grantee/owner | — | `200` | – | K |
| R25 | POST | `/cases/{id}/actions/{action}/ack` | S designated actor | `{note?}` | `200 Case` | `FORBIDDEN_ACTION` | K |
| R26 | GET | `/cases/{id}/safe-actions` | S participant | — | `SafeActionSet{allowed[], blocked[]}` | – | N |
| R27 | GET | `/cases/{id}/stream` | S participant | `Last-Event-ID` | SSE events `case.version`, `passport.version`, `timeline.append`, `heartbeat` | – | N |
| R28 | GET | `/cases/{id}/truth-gap-graph` | S participant | — | `TruthGapGraph{lanes, nodes, edges, gaps}` | – | N |
| R29 | GET | `/cases/{id}/evidence` | S analyst/admin (+role-scoped subset) | `?kind&tier&cursor` | `{items:[EvidenceSummary]}` | – | N |
| R30 | GET | `/evidence/{id}` | S role-scoped | — | `Evidence` (redacted) + `provenance` | `NOT_FOUND` | N |
| R31 | POST | `/cases/{id}/internal-reconciliation` | S analyst/admin | `{note, evidence_ids[]}` | `200 Case` (application state only) | `INVALID_TRANSITION` | K |
| R32 | GET | `/operations/cases` | S analyst/admin | `?sort=uncertainty&filters` | `{items, next_cursor}` | – | N |
| R33 | GET | `/operations/health` | S analyst/admin | — | `{app_telemetry, provider_status_advisory{source,as_of}, queues}` | – | N |
| R34 | GET | `/operations/webhooks` | S analyst/admin | `?status&cursor` | `{items:[InboxSummary]}` | – | N |
| R35 | POST | `/webhooks/paypal` | P | raw PayPal event | `200` / `400` / `503` | – | U |
| R36 | GET | `/demo/scenarios` | S demo_operator/any (demo mode) | — | `{scenarios:[Scenario{id,label,mode}]}` | – | N |
| R37 | POST | `/demo/scenarios/{scenario}` | S demo_operator | `{mode}` | `202 {run_id, case_id, environment}` | `FORBIDDEN_ACTION` | K |
| R38 | POST | `/demo/reset` | S demo_operator | `{}` | `200` | – | K |
| R39 | POST | `/demo/fault-proxy` | S demo_operator | `{enabled, mode:"DROP_CAPTURE_RESPONSE"}` | `200 {enabled}` | `FORBIDDEN_ACTION` (non-sandbox) | N |
| R40 | GET | `/notifications` | S any | `?cursor` | `{items}` | – | N |
| R41 | POST | `/notifications/{id}/read` | S owner | — | `204` | – | N |
| R42 | GET | `/settings/privacy` | S any | — | `{consents, retention, share_links}` | – | N |
| R43 | POST | `/settings/privacy/export` | S self | — | `202 {job_id}` | – | K |
| R44 | POST | `/settings/privacy/deletion-request` | S self | `{scope}` | `202` | – | K |
| R45 | GET | `/settings/provider` | S admin | — | `{environment:"REAL_SANDBOX", client_id_present, secret_present, webhook_id_present}` (booleans only) | – | N |
| R46 | GET | `/audit` | S admin | `?cursor&purpose` | `{items, chain_ok}` | – | N |
| R47 | GET | `/status/provider` | S | — | `ProviderStatusAdvisory{source,as_of,state}` | – | N |
| R48 | GET | `/healthz`, `/readyz` (**no `/api/v1` prefix**) | — | — | `200`/`503` | – | N |

**Streaming semantics (R27):** events carry an incrementing `id`; reconnect with `Last-Event-ID` replays missed events from `outbox`-derived `sse_events`; 25 s heartbeat; authorisation rechecked on connect and on consent revocation (stream closed).

---

## 15. Reconciliation and transition rules

### 15.1 Reconciler (Reader-only) — SOURCE M13 §13.4

```go
func (r *Reconciler) Reconcile(ctx context.Context, caseID uuid.UUID) error {
    // 1 authenticated, case-scoped read permission (job carries tenant + case)
    att := r.store.OriginalAttempt(ctx, caseID)              // stored order_id, request ids, capture ids
    if att.OrderID == "" { return r.review(ctx, caseID, "NO_ORIGINAL_IDENTIFIERS") }   // no arbitrary search

    ord, rawO, err := r.pp.ShowOrder(ctx, att.Env, att.OrderID)                           // GET only
    if err != nil { return r.recordUnknownAndRetry(ctx, caseID, err) }                    // stays UNKNOWN
    var capID = ord.CaptureRef()
    var cap paypal.Capture
    if capID != "" { cap, rawC, err = r.pp.ShowCapture(ctx, att.Env, capID) }             // independent confirmation

    // compare amount/currency/merchant identifiers against the canonical case
    if mismatch(ord, cap, case) { return r.review(ctx, caseID, "AMOUNT_OR_MERCHANT_MISMATCH") }

    return r.store.Tx(ctx, func(tx Tx) error {
        ev := tx.InsertEvidence(E1, PROVIDER_GET_ORDER/CAPTURE, verified=true, redactedExtract, sha)
        tx.Reducer(ev)                       // the ONLY state writer
        tx.Outbox(INVESTIGATE{trigger: "reconciliation.completed"}, dedupe=snapshot)
        tx.Audit("reconcile", ...)
        return nil
    })
    // NEVER: create order, capture, retry mutation. Interface cannot.
}
```

Outcome table: capture `COMPLETED` & application unknown/failed → `application_state=INCONSISTENT`, open `APPLICATION_DIVERGENCE` investigation, recommend internal reconciliation (not a second charge) · `PENDING` → schedule bounded follow-up · contradictory → `NEEDS_HUMAN_REVIEW` · inconclusive → preserve `UNKNOWN` with last verification timestamp.

### 15.2 Reducer evidence ranking
Order evidence by `(tier ASC, provider_occurred_at DESC, received_at DESC)`; E1 > E2 > E3 > E4 > E5 > E6. A **fresh E1** overrides an older E2 for the same field (SOURCE §13.6: webhook says pending, fresh GET says completed → completed, history kept).

### 15.3 Provider-state transition table (legal moves)

| From \ To | CREATED | APPROVED | AUTHORISED | CAPTURE_PENDING | CAPTURE_COMPLETED | CAPTURE_DENIED | REFUNDED/PARTIALLY_REFUNDED | UNKNOWN |
|---|:-:|:-:|:-:|:-:|:-:|:-:|:-:|:-:|
| NOT_CREATED | ✔ | – | – | – | – | – | – | ✔ (lost) |
| CREATED | – | ✔ | – | – | – | – | – | ✔ |
| APPROVED | – | – | ✔ | ✔ | ✔ | ✔ | – | ✔ |
| AUTHORISED | – | – | – | ✔ | ✔ | ✔ | – | ✔ |
| CAPTURE_PENDING | – | – | – | – | ✔ | ✔ | – | ✔ |
| CAPTURE_COMPLETED | – | – | – | **✘ (ignore stale)** | – | **✘ conflict→review** | ✔ (verified refund/reversal evidence only) | **✘ never reverts to UNKNOWN** |
| CAPTURE_DENIED | – | – | – | – | **✘ conflict→review** | – | – | – |
| UNKNOWN | via E1/E2 | via E1/E2 | via E1/E2 | via E1/E2 | via E1/E2 | via E1/E2 | – | – |

`UNKNOWN` is **leaving-only by evidence**: it can resolve to any state a verified read supports; it **cannot** transition to `CAPTURE_DENIED` from a timeout, only from an E1/E2 provider observation.

### 15.4 Forbidden transitions (tested)
1. Any timeout/5xx/EOF-derived write of `CAPTURE_DENIED`, `CANCELLED`, or case `CLOSED_UNRESOLVED`.
2. `CAPTURE_COMPLETED →` `CREATED/APPROVED/AUTHORISED/CAPTURE_PENDING/UNKNOWN`.
3. Any state change sourced from `SIMULATED` evidence on a `REAL_SANDBOX` case.
4. `receiver_state` to `VERIFIED_*`/`WITHDRAWAL_IN_PROGRESS` from provider capture evidence.
5. Creating a second attempt while a non-terminal attempt exists (DB-enforced).
6. Setting `case_status=RESOLVED` unless the resolution predicate holds (§16.6).

### 15.5 Resolution predicate (deterministic, I12)
`RESOLVED` iff: (provider_state is terminal-known by E1/E2 **or** explicitly `NOT_APPLICABLE` objective) ∧ (application_state = `SYNCHRONISED`) ∧ (all required `recovery_actions` acknowledged or waived by an analyst) ∧ (latest Passport version validated at current `snapshot_version`) ∧ (no open contradiction). Never an LLM confidence value.

---

## 16. Evidence, provenance, Truth Gap Graph, Passport

### 16.1 Event envelope schema (`packages/contracts/event.schema.json`)
```json
{
  "$id": "paybridge.event.v1",
  "type": "object",
  "required": ["event_id","event_schema","case_id","source","source_trust","environment","event_type","observed_at","verified","payload_sha256","redacted"],
  "properties": {
    "event_id": {"type":"string","format":"uuid"},
    "event_schema": {"const":"paybridge.event.v1"},
    "case_id": {"type":"string","format":"uuid"},
    "source": {"enum":["PAYPAL_WEBHOOK","PAYPAL_GET","APP_TRACE","MERCHANT_DB","PUBLIC_STATUS","USER_REPORT","SIMULATOR"]},
    "source_trust": {"enum":["PROVIDER_VERIFIED","LOCAL_OBSERVED","PUBLIC_SOURCE","USER_REPORTED","SIMULATED"]},
    "environment": {"enum":["REAL_SANDBOX","PAYPAL_SIMULATOR","LOCAL_SIMULATION","PRODUCTION"]},
    "event_type": {"type":"string"},
    "provider_resource": {"type":"object","properties":{"type":{"enum":["ORDER","AUTHORIZATION","CAPTURE","REFUND"]},"id":{"type":"string"}}},
    "occurred_at": {"type":["string","null"],"format":"date-time"},
    "observed_at": {"type":"string","format":"date-time"},
    "verified": {"type":"boolean"},
    "payload_sha256": {"type":"string","pattern":"^[a-f0-9]{64}$"},
    "redacted": {"type":"boolean"}
  },
  "additionalProperties": false
}
```

### 16.2 Evidence tiers → tier numbers
E1 = provider GET (verified, authenticated) · E2 = verified signed webhook · E3 = merchant-controlled app log/DB snapshot · E4 = official docs/public status · E5 = consented user testimony · E6 = simulation or model output. Strength is **field-specific**; a "browser timed out" fact is E3, a "capture completed" fact must be E1/E2.

### 16.3 Truth Gap Graph (SOURCE Feature B)
**Lanes:** `BROWSER_CLIENT`, `MERCHANT_APP`, `PROVIDER` (PayPal), `WEBHOOK`, `RECEIVER`. Each lane holds the latest **claim** about the payment status with `{claim_code, value, evidence_ids, tier, observed_at, trust, environment}`.

**Gap detector (deterministic):** compares lane claims via a rule table; every rule has `rule_id`, `version`.

| Rule | Condition | Edge | Gap label |
|---|---|---|---|
| G01 | BROWSER=`UNKNOWN/timeout` ∧ PROVIDER=`CAPTURE_COMPLETED` | CONTRADICTS | "Browser uncertainty vs provider completion" |
| G02 | MERCHANT_APP=`NO_PAYMENT` ∧ PROVIDER=`CAPTURE_COMPLETED` | CONTRADICTS | "Merchant record behind provider" |
| G03 | WEBHOOK=`PENDING` ∧ PROVIDER GET=`COMPLETED` (fresher) | PRECEDES/SUPERSEDED | "Webhook stale vs provider read" |
| G04 | PROVIDER=`CAPTURE_COMPLETED` ∧ RECEIVER=`NOT_OBSERVABLE` | INDEPENDENT_OF | "Bank stage not observable" (informational, not a contradiction) |
| G05 | Amount/currency mismatch provider vs case | CONTRADICTS | "Amount differs from invoice" |
| G06 | WEBHOOK duplicate | CORROBORATES | "Same event received N times" |

`GET /cases/{id}/truth-gap-graph` returns `{lanes[], nodes[], edges[], gaps[{rule_id, severity, lanes, evidence_ids, plain_text}]}`; `plain_text` is **template-generated** (not LLM). LLM-suggested relations stay `origin=LLM_HYPOTHESIS` and are rendered as hypotheses.

### 16.4 Causal chain (example, SOURCE M11 §11.6)
`provider capture completed` (E1) → `merchant response timeout` (E3) → `merchant order stale` (E3) → `case inconsistent` (deterministic) → "payment succeeded at provider, application confirmation failed" (validated claim).

### 16.5 Passport versioning (SOURCE §11.7)
- Passport per `(case, viewer_role)`; **append-only** `passport_versions`.
- A new version is created when: a new verified event/evidence changes any claim; the validated claim set changes; policy_version changes; consent changes the visible field set.
- `content_sha256 = SHA-256(canonical JSON)`; `prev_sha256` chains versions.
- **Contradiction-resistant publish gate** before release: (1) no statement contradicts the current verified snapshot; (2) stale evidence marked; (3) role-based redaction occurred; (4) every `confirmed` claim has valid evidence; (5) snapshot_version still current.
- Schema (`packages/contracts/passport.schema.json`):

```json
{
  "passport_schema": "1.0",
  "case_id": "uuid",
  "viewer_role": "sender|receiver|ops",
  "environment": "REAL_SANDBOX|PAYPAL_SIMULATOR|LOCAL_SIMULATION",
  "version": 2,
  "snapshot_version": 7,
  "financial_summary": {
    "amount": {"value": "2000.00", "currency": "USD"},
    "provider_capture_status": "COMPLETED",
    "application_state": "API_TIMEOUT",
    "receiver_state": "NOT_OBSERVABLE"
  },
  "last_verified_at": "2026-11-09T07:00:04Z",
  "confirmed_claims": [{"code":"CAPTURE_CONFIRMED_APP_RESPONSE_LOST","text":"...","evidence_ids":["..."],"tier":1}],
  "uncertain_claims": [],
  "unknowns": [{"code":"BANK_CREDIT_NOT_OBSERVABLE","text":"..."}],
  "next_safe_action": {"code":"DO_NOT_REPEAT_CHECKOUT","text":"..."},
  "share_policy": {"expires_at": null, "can_revoke": true},
  "labels": {"simulated": false}
}
```
(All values shown are **examples**, not real results.)

### 16.6 Resolution
See §15.5.

---

## 17. Next Safe Action Engine

**Purpose (SOURCE Feature D):** a deterministic action catalogue that blocks dangerous options based on evidence class and state. The AI may *request* an allowed recommendation; it cannot add, unblock or reorder safety policy.

### 17.1 Input facts
`provider_state`, `last_confirmed_provider_state`, `application_state`, `receiver_state`, `active_attempt.status`, `evidence_age`, `has_E1_after_unknown`, `viewer_role`, `consent`, `environment`.

### 17.2 Action catalogue (all non-money-moving)

| Code | Meaning | Allowed when | Owner | Needs human approval |
|---|---|---|---|---|
| `DO_NOT_REPEAT_CHECKOUT` | Banner: do not pay again | `provider_state ∈ {UNKNOWN, CAPTURE_PENDING, CAPTURE_COMPLETED}` ∨ active non-terminal attempt | system | no |
| `AWAIT_PROVIDER_CONFIRMATION` | Wait; follow-up scheduled | `CAPTURE_PENDING` or `UNKNOWN` with reconcile queued | system | no |
| `REFRESH_ORIGINAL_ORDER` | Re-read the same order/capture | any unknown/stale evidence | sender/analyst | no (read-only) |
| `RECONCILE_MERCHANT_ORDER` | Fix the **internal** merchant record | `CAPTURE_COMPLETED` ∧ `application_state ∈ {API_TIMEOUT, INCONSISTENT}` | merchant_analyst/admin | **yes** |
| `CHECK_OFFICIAL_ACCOUNT_ACTIVITY` | Receiver checks PayPal account | `CAPTURE_COMPLETED` ∧ receiver consent | receiver | no |
| `COMPLETE_OFFICIAL_VERIFICATION` | Provider-requested user action | provider reports user-action requirement | receiver/sender | no |
| `CONTACT_PROVIDER_SUPPORT` | Escalate with evidence pack | contradictions persist / provider fault suspected | analyst | yes |
| `REQUEST_HUMAN_REVIEW` | Open analyst review | any unresolved conflict/invalid AI output | system/analyst | no |
| `NONE_REQUIRED` | Nothing to do | resolved | system | no |

### 17.3 Blocked actions (always returned with reasons)

| Blocked code | Blocked when | Reason code |
|---|---|---|
| `CREATE_NEW_ORDER` | active non-terminal attempt exists | `ACTIVE_ATTEMPT_EXISTS` |
| `CREATE_NEW_ORDER` | `provider_state ∈ {UNKNOWN, CAPTURE_PENDING}` | `OUTCOME_UNKNOWN_RECONCILE_FIRST` |
| `CREATE_NEW_ORDER` | `provider_state = CAPTURE_COMPLETED` | `PAYMENT_ALREADY_CAPTURED` |
| `CAPTURE_AGAIN` | any prior capture requested | `CAPTURE_ALREADY_REQUESTED` |
| `REFUND`, `VOID`, `REVERSE` | always in MVP | `NOT_SUPPORTED_NO_MONEY_ACTIONS` |
| `MARK_FAILED` | provider state unknown | `UNKNOWN_IS_NOT_FAILED` |
| `MARK_BANK_SETTLED` | no authorised settlement signal | `NOT_OBSERVABLE` |

`START_NEW_ATTEMPT` (a *user-confirmed* new commercial purchase) is allowed **only** when the prior attempt is terminal (`CANCELLED`/`CAPTURE_DENIED` by **E1/E2**), and requires an explicit confirmation step with a fresh `Idempotency-Key`. It is executed by `checkout`, not by the engine.

### 17.4 Properties (tested)
- Pure function `Decide(facts) → {allowed[], blocked[]}`; exhaustive property tests over the enum cross-product assert: no blocked code is ever in `allowed`; unknown never yields `MARK_FAILED`; `CREATE_NEW_ORDER` is blocked in all non-terminal states.
- Output copy keys are stable codes; **text is template-rendered** from `projector` (doc 06 copy rules), never LLM-authored.

---

## 18. Two-Lens privacy projections

**Principle (SOURCE I8 / M16 §16.5):** two-stage control — (1) **role-scoped input retrieval**; (2) **role-scoped output policy validation**. A prompt saying "don't reveal" is insufficient.

### 18.1 Field allowlists (default deny)

| Field | Sender | Receiver (with consent) | Ops/Analyst | AI snapshot (per role) |
|---|:-:|:-:|:-:|:-:|
| Case id, invoice reference | ✔ | ✔ (reference only) | ✔ | ✔ |
| Amount + currency | ✔ | ✔ | ✔ | ✔ |
| Provider capture status | ✔ | ✔ | ✔ | ✔ |
| Provider order/capture ids (masked) | ✔ (masked, copyable) | masked (last 4) | ✔ full | masked |
| Application state / timeout trace | ✔ plain language | ✘ (seller-side summary only) | ✔ | ✔ per role |
| Webhook delivery diagnostics | ✘ | ✘ | ✔ | ops only |
| Sender funding instrument / issuer response | ✘ never stored (opaque last-4 only if needed) | ✘ | ✘/restricted | ✘ |
| Receiver account verification/limits | ✘ | ✔ (own) | restricted | receiver-only |
| Receiver bank/withdrawal | n/a (`NOT_OBSERVABLE`) | n/a | n/a | n/a |
| Other party's email/name | masked | masked | ✔ | masked |
| Evidence items | role-class subset | `RECEIVER`-class subset | all tenant | per-role subset |
| Internal audit trail | ✘ | ✘ | admin | ✘ |

### 18.2 Mechanism
`projector.Build(ctx, role, caseSnapshot) → Projection` selects **only** allowlisted fields and *pre-redacts* identifiers; agents A03/A04/A09 receive **only** their role's projection (S8). The Claim Validator then **scans every generated message** against a forbidden-field detector (value-level fingerprints of fields not in the role's projection, e.g. the other party's masked identifiers). A hit ⇒ redact or block + audit + `needs_human_review`.

### 18.3 Copy contracts (excerpt; full set in doc 06)
- Sender (hero): *"PayPal confirmed your payment capture. Our website didn't receive the response. We're reconciling your original payment. Please don't pay again."*
- Receiver (hero): *"A PayPal capture was recorded on the seller's side. Use your official PayPal account activity to see availability and withdrawal. We can't confirm bank credit from the data available here."*

---

## 19. Audit, retention, deletion, consent

### 19.1 Audit (SOURCE M16 §16.6)
Each sensitive action writes: actor, tenant, case, operation, previous/new version, evidence refs, environment, reason, request id, policy version, coarse client metadata, timestamp, reviewer. **Hash chain:** `row_hash = SHA256(prev_hash ‖ canonical(row))`; periodic digest export; `GET /audit` verifies `chain_ok`. This detects tampering under stated assumptions; it is **not** legal admissibility or regulatory certification.

### 19.2 Retention (proposed defaults, **subject to legal review**)

| Data class | Handling |
|---|---|
| Case metadata | Retain per merchant policy; configurable |
| Raw verified webhook | Encrypted, restricted, minimised, deletion policy |
| Redacted extracts | Keep with case for retention period |
| App traces | Redact secrets/PII/tokens/full payloads |
| Agent prompts/outputs | Store sanitised prompt hash, model/version, structured findings |
| Share links | Expire; revocable; `X-Robots-Tag: noindex` |
| Telemetry raw | Aggregate, then expire |

### 19.3 Deletion by design
Request → policy validation → erase/anonymise eligible rows + encrypted artefacts + caches → record tombstone; legally required audit retained in restricted form **only after counsel's direction**. Never promise instant deletion of all logs. Jobs run as `paybridge_retention` role.

### 19.4 Consent
- Invitation shows exactly which fields will be shared; accept ⇒ `consent_grants` row with `field_allowlist` and separate `notification_allowed`.
- Revocation (S9): grant revoked; share tokens revoked; SSE streams closed; pending notices cancelled; future access `404`; delivered messages cannot be recalled (stated in UI).

---

## 20. Transactions, locking, concurrency

| Concern | Mechanism |
|---|---|
| Duplicate checkout click | UI disables; `Idempotency-Key`; `uq_one_active_attempt`; `SELECT … FOR UPDATE` on attempt row in capture |
| Two workers, same webhook | Unique inbox key + `SKIP LOCKED` lease + state-version check |
| Reducer serialisation | `pg_advisory_xact_lock(hashtext(case_id::text))` at the start of every reducer tx |
| Optimistic concurrency | `UPDATE payment_cases … WHERE id=$1 AND state_version=$2 RETURNING` → on 0 rows retry or `STATE_VERSION_CONFLICT` |
| New event during LLM run | `investigations.snapshot_version` compared at publish; mismatch ⇒ discard (`STALE_DISCARDED`) and re-enqueue |
| Isolation level | `READ COMMITTED` default with explicit locks; `SERIALIZABLE` not required if the above hold (**DESIGN DECISION**; revisit with tests) |
| Provider call inside a tx | **Never.** Pattern: tx1 persist intent + keys → commit → provider call → tx2 persist outcome. A crash between leaves a persisted intent that reconciliation picks up (CHAOS-11) |
| Long transactions | Forbidden around LLM/provider/network calls |

---

## 21. Observability

**Logs:** structured JSON; allowlisted keys (`correlation_id, case_id, tenant_id_hash, route, op, status, latency_ms, env, provider_resource_type, operation_key_hash`). No bodies, tokens, secrets, cookies, emails, or raw webhooks.

**Metrics (Prometheus/OTel names):** `pb_http_requests_total`, `pb_http_latency_seconds`, `pb_webhook_received_total{result}`, `pb_webhook_signature_failures_total`, `pb_inbox_unprocessed`, `pb_outbox_ready`, `pb_outbox_oldest_age_seconds`, `pb_unknown_state_backlog_age_seconds`, `pb_reconcile_duration_seconds`, `pb_provider_calls_total{op,class}`, `pb_agent_tokens_total{agent}`, `pb_agent_rejections_total{reason}`, `pb_claim_validator_rejections_total`, `pb_sse_connections`, `pb_dead_letters_total`.

**Traces:** one trace per request/job; spans for provider calls, reducer, agents; attributes use the safe-key set.

**Health:** `/healthz` (process up) and `/readyz` (DB reachable, migrations applied, required config present; LLM/PayPal reachability reported as *degraded*, not failing readiness).

**Alert runbooks:** doc 09 §Runbooks (webhook failure spike, unknown backlog, provider incident, LLM down, wrong transition, suspected cross-tenant access).

---

## 22. OpenAPI

- File: `packages/contracts/openapi.yaml` (OpenAPI 3.1). Components: all schemas above; security schemes `sessionCookie`, `csrfHeader`, `shareToken`, `paypalSignature`.
- Codegen: `oapi-codegen` → `internal/httpx/gen` (server stubs + types); `openapi-typescript` → `apps/web/src/api/schema.d.ts`.
- CI checks: (a) spec lints (Spectral); (b) generated code is up to date (`git diff --exit-code` after generate); (c) handler registry matches spec operations 1:1; (d) contract tests validate responses against the spec.
- Every operation declares: auth, roles, idempotency requirement, error responses. An architecture test fails if a route lacks an `x-authorize` annotation.

Skeleton:
```yaml
openapi: 3.1.0
info: { title: PayBridge AI API, version: 1.0.0 }
servers: [{ url: /api/v1 }]
paths:
  /cases/{id}/capture:
    post:
      operationId: captureCase
      security: [{ sessionCookie: [], csrfHeader: [] }]
      x-authorize: { roles: [sender], action: capture, scope: own_case }
      parameters:
        - { name: Idempotency-Key, in: header, required: true, schema: { type: string, format: uuid } }
      responses:
        "200": { description: Captured, content: { application/json: { schema: { $ref: "#/components/schemas/Case" } } } }
        "202": { description: Outcome unknown, content: { application/json: { schema: { $ref: "#/components/schemas/UnknownOutcome" } } } }
        "409": { $ref: "#/components/responses/Conflict" }
```

---

## 23. Tests

*All are **NOT RUN** — this is the plan.*

| Layer | What | Tools |
|---|---|---|
| Unit | money/decimal, validators, classifier, reducer transition table, safe-action property tests, projector allowlists, claim validator, backoff | `testing`, `testify`, `rapid` (property) |
| DB integration | migrations up from empty; every constraint (I1, I2, I7, I11, one-active-attempt) violated deliberately; RLS enumeration test | real Postgres (testcontainers/Compose) |
| Handler/contract | each route vs OpenAPI; auth matrix; uniform 404; idempotency replay/reuse/in-flight | `httptest`, contract validator |
| Webhook | forged rejected; duplicate; out-of-order; unverifiable→503; uncorrelated | fake PayPal + fixtures (`SIMULATED`) |
| Reconciliation | CHAOS-01/02/11/12 with fake PayPal; assert **zero POSTs** from the Reconciler (fake counts verbs) | fake PayPal |
| Concurrency | two workers/one webhook; double capture click; stale snapshot | `-race`, goroutine harness |
| Architecture | import rules; Reader has only GET; every route authorises | `archtest` |
| Real sandbox | tagged `realsandbox`: create→approve→capture→GET; webhook delivery | real creds, manual/CI-secret |
| Security | IDOR matrix, CSRF, header checks, secret scan | doc 08/09 |

Traceability IDs map to doc 09's requirements-to-test matrix.

---

## 24. Implementation sequence

*(Maps to doc 02 phases; each step = one bounded task.)*

1. **P1.1** repo skeleton, Makefile, Compose (Postgres), config validator (+ `.env.example`), health endpoints, slog/OTel scaffolding.
2. **P1.2** migrations 0001–0003; `store` + sqlc; tx helper with RLS session setup; DB constraint tests.
3. **P1.3** `money`, `ids`, `redact`; PayPal adapter: OAuth + Reader + Mutator + classifier + request-id discipline; fake PayPal; unit tests.
4. **P1.4** `checkout`: create order → approval context; capture with lock; unknown handling; **real sandbox receipt** (manual, tagged test).
5. **P2.1** migrations 0004–0005; evidence/events; Reducer + transition tests.
6. **P2.2** Webhook Gateway + verify + inbox; Inbox processor; duplicate/out-of-order tests.
7. **P2.3** Outbox relay + worker framework + dead-letter; Reconciler; CHAOS-01/02/04/05/11/12 on fake PayPal.
8. **P3.1** auth/sessions/CSRF; policy engine; RBAC/ABAC; RLS policies (0010); IDOR/tenant tests.
9. **P3.2** consent/invitations/shares; projector; passport builder/versioning; safe-action engine; truth-gap builder.
10. **P3.3** audit hash chain; SSE hub; ops endpoints; idempotency middleware complete.
11. **P5.x** (after frontend core) `llmgw`, `tools`, `agents`, `orchestrator`, `claims`, evals — specified in doc 07.
12. **P6/P7** simulator + Fault Proxy + demo endpoints; retention jobs; hardening.

*End of document 04.*
