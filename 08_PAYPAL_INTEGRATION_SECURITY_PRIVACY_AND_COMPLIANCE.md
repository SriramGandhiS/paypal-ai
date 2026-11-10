# 08 — PayPal Integration, Security, Privacy and Compliance

| Field | Value |
|---|---|
| Document | `08_PAYPAL_INTEGRATION_SECURITY_PRIVACY_AND_COMPLIANCE.md` |
| Version | 1.0.0 |
| Date | 9 November 2026 |
| Status | **DRAFT — specification. Nothing here has been executed, verified against live PayPal behaviour, or reviewed by counsel/security professionals.** |
| Audience | Backend engineers, security reviewers, owner, coding agents |
| Depends on | 01, 03 (trust boundaries), 04 (adapter §9, webhook §12, RLS §11.4, projections §18) |
| Related | 02 (key handling), 07 (AI privacy), 09 (security tests) |
| Labels | SOURCE · DESIGN DECISION · **VERIFY** · SANDBOX · PRODUCTION |

> **Compliance statement:** this document describes engineering controls and a privacy design. It is **not** legal advice and **never** asserts that PayBridge AI is compliant with any law, regulation, standard, or PayPal program. Regulatory and payment-eligibility questions are marked **VERIFY — current official sources and qualified professional advice**.

## Table of contents
1. [PayPal sandbox setup](#1-paypal-sandbox-setup) · 2. [Official documentation to re-check](#2-official-documentation-to-re-check) · 3. [Supported APIs and eligibility caveats](#3-supported-apis-and-eligibility-caveats) · 4. [OAuth and token lifecycle](#4-oauth-and-token-lifecycle) · 5. [Orders and capture — only as permitted](#5-orders-and-capture--only-as-permitted) · 6. [Webhooks: subscription and verification](#6-webhooks-subscription-and-verification) · 7. [Provider authority](#7-provider-authority) · 8. [Sandbox–production isolation](#8-sandboxproduction-isolation) · 9. [Security architecture](#9-security-architecture) · 10. [Threat model (STRIDE)](#10-threat-model-stride) · 11. [OWASP API and web risk mapping](#11-owasp-api-and-web-risk-mapping) · 12. [Specific attack classes](#12-specific-attack-classes) · 13. [Secrets and dependency management](#13-secrets-and-dependency-management) · 14. [Encryption and data minimisation](#14-encryption-and-data-minimisation) · 15. [Privacy matrix](#15-privacy-matrix) · 16. [Consent, retention, deletion](#16-consent-retention-deletion) · 17. [India payment eligibility and compliance (VERIFY)](#17-india-payment-eligibility-and-compliance-verify) · 18. [Unobservable provider and bank data](#18-unobservable-provider-and-bank-data) · 19. [Incident response](#19-incident-response) · 20. [Release gates](#20-release-gates) · 21. [`.env.example`](#21-envexample) · 22. [Security tests](#22-security-tests)

---

## 1. PayPal sandbox setup

*Console labels and flows change — **VERIFY** each step in the current PayPal Developer Dashboard.*

| Step | Action | Output (record **names/presence only** in docs) |
|---|---|---|
| 1 | Create/sign in to a PayPal Developer account | – |
| 2 | Dashboard → **Apps & Credentials** → **Sandbox** → create a REST app (name e.g. `paybridge-sandbox`) | Client ID, Client Secret → env vars `PAYPAL_CLIENT_ID`, `PAYPAL_CLIENT_SECRET` |
| 3 | **Sandbox accounts**: ensure a **Business** (merchant/receiver) and a **Personal** (buyer) test account exist; note their sandbox emails privately | Used for approval step; **never** commit passwords |
| 4 | Confirm the REST app is linked to the sandbox **Business** account (merchant) and which features are enabled (VERIFY available products: Orders/Payments/Webhooks) | – |
| 5 | Decide the **public HTTPS webhook URL**: `${PUBLIC_BASE_URL}/api/v1/webhooks/paypal` (Render service URL or a tunnel for local) | `PUBLIC_BASE_URL` |
| 6 | App settings → **Webhooks** → add the URL; subscribe to event types in §6.2 (VERIFY names) | Webhook ID → `PAYPAL_WEBHOOK_ID` |
| 7 | Run the **P1 proof**: create → approve (buyer sandbox account) → capture → GET order → GET capture; store a **redacted** receipt in `docs/evidence/paypal-sandbox-receipt.md` | AO-01 evidence |
| 8 | (Optional) Try PayPal's documented negative-testing mechanism for declines **only if supported** (CHAOS-03; VERIFY) | – |

**Cross-border configuration (AS-03):** attempt to configure a US-payer / non-US-receiver business route in the sandbox. If the sandbox cannot represent it, use a supported sandbox configuration for the **real** proof and keep the US→IN readiness narrative in the labelled simulation layer. Document the outcome either way in `docs/evidence/`.

**Webhook simulator caveat (SOURCE M12):** the developer-dashboard webhook simulator produces mock events that do not correspond to real transactions and cannot be verified through the real postback. They must route to the isolated `PAYPAL_SIMULATOR` path and be labelled.

---

## 2. Official documentation to re-check

*(Locations are described by topic; locate current URLs on developer.paypal.com. Do not code from memory — **VERIFY** all.)*

| Topic | What to confirm |
|---|---|
| Orders v2 API (create, show, capture) | Request/response schemas; `intent=CAPTURE`; statuses; `PayPal-Request-Id` behaviour; error shapes; whether `PayPal-Request-Id` is supported/required on capture; idempotency retention window |
| Payments v2 (captures) | `GET /v2/payments/captures/{id}` fields; capture statuses (`COMPLETED`, `PENDING`, `DECLINED`, …); `status_details` reasons; refund endpoints (out of scope) |
| Webhooks | Event type names; payload shape; delivery/retry policy; signature headers; **verify-webhook-signature** request/response; webhook ID semantics; sandbox vs live |
| Authentication | OAuth2 client-credentials; token lifetime; scope behaviour |
| Idempotency guidance | Which endpoints honour `PayPal-Request-Id`; retention |
| Sandbox testing | Test accounts; negative testing/simulation headers; limits and known sandbox differences |
| Regional/receiving requirements | India (and any other route) receiving/eligibility requirements; account-type rules |
| Developer terms / acceptable use / brand guidelines | Permitted usage; **no implied endorsement**; logo/brand use |
| Hackathon rules | Eligibility, required stack/sponsor tech, submission artefacts, licence, video limits, deadline (**VERIFY on the official hackathon page**) |

**Rule:** every PayPal fact that appears in code comments, README, or pitch must cite *document + retrieved date* in `docs/evidence/paypal-docs-register.md`.

---

## 3. Supported APIs and eligibility caveats

| Capability | MVP use | Caveat (VERIFY) |
|---|---|---|
| Orders v2 create/show/capture | **Yes** (human-triggered checkout only) | Sandbox features can differ from live; eligibility depends on app/account configuration |
| Payments v2 show capture | **Yes** (Reader) | Field availability varies |
| Webhooks + verify signature | **Yes** | Delivery to a public HTTPS endpoint required; ordering/exactly-once **not** guaranteed |
| Refund/void/reverse | **No** (no interface exposes them) | Future, human-authorised only |
| Payouts/Disputes/Reporting/Transaction Search | **No** | Possible later with permissions; each has its own eligibility |
| Receiver availability/withdrawal/bank settlement | **Not observable** | No API in scope exposes this (§18) |
| Cross-border (US→IN) | Readiness narrative + best-effort sandbox | Account-type and product rules vary; may not be configurable in sandbox |

**Never claimed:** that any account is eligible; that a route is supported; that sandbox behaviour equals production.

---

## 4. OAuth and token lifecycle

Authoritative design in doc 04 §9.3. Security requirements:

- Client credentials from env only; secret never leaves the server process; never logged, stored in evidence, or sent to the browser or LLM.
- Token cached in memory; refresh at `expires_at − 60 s`; singleflight; on `401` refresh once.
- Use `Authorization: Bearer` over TLS only; verify certificates; no custom CA bypass.
- Distinct credentials per environment; sandbox credentials only in the MVP; config validator **refuses to start** if `PAYPAL_API_BASE` is not the sandbox host or if `PAYPAL_ENV != sandbox`.
- Rotation: documented in doc 09 runbook (regenerate secret in console → update platform secret → restart → verify with health probe → revoke old).

---

## 5. Orders and capture — only as permitted

| Rule | Detail |
|---|---|
| Initiation | Only the human-triggered `checkout` module; routes R09/R10; sender role on own case |
| Amount/currency | From canonical case; verified against provider response (`AMOUNT_MISMATCH` blocks success label) |
| Buyer approval | Required before capture; PayPal hosted approval flow (PayBridge never sees buyer credentials) |
| Capture | Only on the original approved order; row-locked attempt; stable `PayPal-Request-Id`; **one** capture request per attempt |
| Retry | No blind retry; unknown ⇒ reconcile (read-only) |
| Refund/void/reverse | Not implemented |
| AI | Zero access to mutators (import rules + process separation) |
| Brand | No PayPal logos or buttons imitating PayPal beyond what PayPal's official SDK/brand rules permit (**VERIFY**); neutral text "PayPal" only; disclaimer present |

*(If the official PayPal JS SDK buttons are used for approval, load it from PayPal's official host only, with a strict CSP allowance — **VERIFY** current recommended integration; a redirect-based approval flow is the simpler MVP default.)*

---

## 6. Webhooks: subscription and verification

### 6.1 Subscription
Register the webhook in the sandbox app; store the resulting **Webhook ID** as `PAYPAL_WEBHOOK_ID`. One webhook per environment. Re-register if the public URL changes.

### 6.2 Event types to subscribe (VERIFY exact names)
`CHECKOUT.ORDER.APPROVED`, `PAYMENT.CAPTURE.COMPLETED`, `PAYMENT.CAPTURE.PENDING`, `PAYMENT.CAPTURE.DENIED`, `PAYMENT.CAPTURE.REFUNDED`, `PAYMENT.CAPTURE.REVERSED` (plus any others the docs recommend for Orders v2).

### 6.3 Verification procedure (doc 04 §12 is authoritative)
1. Read **raw bytes**; do not re-serialise.
2. Extract the transmission headers (`PAYPAL-TRANSMISSION-ID`, `-TIME`, `-SIG`, `PAYPAL-CERT-URL`, `PAYPAL-AUTH-ALGO`) — **names VERIFY**.
3. Call PayPal's verify-webhook-signature endpoint with `webhook_id` and the parsed event; require `verification_status = SUCCESS` (**VERIFY** response field).
4. Failure ⇒ `REJECTED`, 400, metric + alert on spikes. Verification service unavailable ⇒ **503** (let PayPal retry); never acknowledge unverified data.
5. Correlate resource ids to a known case/tenant **after** verification (a valid signature is not tenant authorisation).
6. Dedupe on event id; reducer is monotonic (replay-harmless).
7. Optional hardening (LATER): offline certificate-chain verification with pinned/validated cert URLs restricted to PayPal domains and CRC32+RSA check per docs — only if the docs recommend it (VERIFY); never fetch arbitrary `PAYPAL-CERT-URL` hosts (SSRF).

---

## 7. Provider authority

*(SOURCE ADR-001, I4, I6.)*
- Verified provider evidence (E1/E2) outranks local state, user claims, and AI.
- A `COMPLETED` capture is **not** bank settlement. `receiver_state` stays `NOT_OBSERVABLE` without separately authorised evidence.
- A webhook is a *notification*, not proof of delivery order; a fresh authenticated GET (E1) overrides an older webhook for the same field.
- Contradictory E1 evidence ⇒ human review; AI explanations for that field are withheld.
- PayBridge is an **observer and coordinator**; PayPal is the processor of record.

---

## 8. Sandbox–production isolation

| Control | Detail |
|---|---|
| Environment enum | `REAL_SANDBOX`, `PAYPAL_SIMULATOR`, `LOCAL_SIMULATION`, `PRODUCTION` on **every** row; `PRODUCTION` hard-disabled in MVP |
| Config guard | Startup fails if `PAYPAL_ENV != sandbox` or host ≠ sandbox host or if any variable named `*_LIVE_*` / live client id detected |
| Separate secrets | No shared secret stores/paths between sandbox and any future production |
| Data separation | Production would use a separate database, project, and keys; sandbox data never promoted |
| UI | Persistent environment badge; sandbox banner on checkout (“No real money”) |
| Namespacing | Simulator ids prefixed (`SIM-…`); never accepted as provider ids by `provider_resources` for `REAL_SANDBOX` cases |
| Promotion | Moving to production is a **separate project** requiring counsel review, PayPal approval/eligibility checks, new threat model, penetration test — **not** part of this blueprint |

---

## 9. Security architecture

```mermaid
flowchart LR
  U["Browser"] -->|"TLS 1.2+"| EDGE["Edge: TLS, WAF or rate limit"]
  EDGE --> API["API: authn, CSRF, Policy Engine"]
  PP["PayPal webhooks"] -->|"signed"| EDGE
  API --> DB[("PostgreSQL: RLS, restricted role")]
  API --> ENC[("Encrypted evidence store: AES-GCM")]
  WK["Worker"] --> DB
  WK --> ENC
  WK -->|"redacted only"| LLM["LLM provider"]
  API -->|"mutating, human-triggered"| PPAPI["PayPal API"]
  WK -->|"read-only GET"| PPAPI
  SEC["Secrets: platform secret store"] --> API
  SEC --> WK
```

**Principles (SOURCE M16):** default deny; least privilege; defence in depth (app **and** DB tenant checks); secure defaults; fail closed on verification; no secrets in agents; minimise, redact, expire; auditability.

**Controls summary**

| Area | Control |
|---|---|
| Transport | HTTPS only; HSTS; secure cookies; internal service-to-service TLS where supported |
| AuthN | Server-side sessions (MVP); MFA for ops roles in any non-demo deployment; OIDC adapter path |
| AuthZ | Policy Engine at every handler; RLS; uniform 404; field-level projections |
| Input | OpenAPI strict validation; size limits; Unicode normalisation; allowlists |
| Output | Context-aware encoding (React escapes by default); no `dangerouslySetInnerHTML`; markdown disabled for untrusted text |
| CSRF | Double-submit token + SameSite=Lax; state-changing routes require header |
| CSP | `default-src 'self'`; script/style from self (hashed/nonce where needed); `connect-src 'self'`; `frame-ancestors 'none'`; PayPal SDK host only if used (VERIFY) |
| Headers | `X-Content-Type-Options`, `Referrer-Policy: no-referrer`, `Permissions-Policy`, `Cache-Control: no-store` on financial API responses; share routes add `X-Robots-Tag: noindex` |
| Logging | Allowlist; no PII/tokens/bodies |
| Backups | Encrypted; restore drilled (doc 09) |

---

## 10. Threat model (STRIDE)

*(Assets: provider credentials, session cookies, case data, evidence, consent state, LLM key, audit log. Actors: anonymous internet, authenticated sender/receiver, malicious insider/analyst, compromised dependency, malicious webhook sender, prompt-injection author.)*

| STRIDE | Threat | Example | Mitigations | Test |
|---|---|---|---|---|
| **S**poofing | Forged webhook | Attacker POSTs fake `PAYMENT.CAPTURE.COMPLETED` | Signature verification via PayPal; webhook ID check; 503 when unverifiable; reject-by-default | CHAOS-05 (T05) |
| S | Session theft/fixation | Stolen cookie | HttpOnly/Secure/SameSite; rotation on login; idle/absolute timeout; IP/UA hash anomaly flag | Session tests |
| S | Receiver identity confusion | Someone else uses invite link | Authenticated intended-identity match (S9) | Invite tests |
| **T**ampering | Alter amounts in checkout | Browser edits amount | Server-authoritative amount; provider response compare | `AMOUNT_MISMATCH` test |
| T | Evidence/audit tampering | DB insider edits events | Append-only triggers, REVOKE UPDATE/DELETE, hash-chained audit | Audit-chain verify |
| T | Replay of old webhook | Re-send valid signed event | Dedupe + monotonic reducer | CHAOS-04/05 |
| **R**epudiation | "I never approved" | Dispute | Audit with actor/time/version/request id; evidence hashes | Audit tests |
| **I**nformation disclosure | Cross-tenant read (IDOR) | Guess case UUID | Policy Engine + RLS + uniform 404 | CHAOS-06 |
| I | Receiver sees sender-private data | Prompt/projection leak | Projection-before-generation; DLP | Privacy diff tests |
| I | Secrets in logs/repo/LLM | Key printed | Redaction, allowlist logging, gitleaks, key isolation | Secret scans |
| I | Share-link leakage | Link forwarded | Expiry, revocation, noindex, minimal fields | Share tests |
| **D**enial of service | Webhook flood | Spam endpoint | Size limit, rate limit, cheap reject, async processing | Load test |
| D | LLM cost exhaustion | Trigger many investigations | Dedupe per snapshot, daily cap, rate limits | Cost-cap test |
| D | Slow provider/DB | Hang | Timeouts everywhere; circuit breaker; readiness | Chaos |
| **E**levation of privilege | Receiver calls capture | Role confusion | Role matrix; server checks | RBAC matrix tests |
| E | LLM tool escalation | Injected instructions call tools | No write tools; allowlist; schema confinement | CHAOS-07 |
| E | Mass assignment | Extra JSON fields set state | Strict schema; DTOs ≠ DB models | Validation tests |

---

## 11. OWASP API and web risk mapping

*(OWASP API Security Top 10 — 2023 edition names; **VERIFY** current list.)*

| Risk | Application |
|---|---|
| API1 Broken Object Level Authorization | Policy Engine on every object route; RLS; IDOR test matrix across all R-routes |
| API2 Broken Authentication | Sessions, login throttling, MFA for ops, no tokens in URLs/localStorage |
| API3 Broken Object Property Level Authorization | Role projections; strict response DTOs; no mass assignment |
| API4 Unrestricted Resource Consumption | Body/size limits; pagination caps; rate limits; LLM budgets; SSE connection caps per user |
| API5 Broken Function Level Authorization | Role × action matrix; explicit `Authorize`; admin routes separated |
| API6 Unrestricted Access to Sensitive Business Flows | Checkout rate limits + one-active-attempt; idempotency; capture lock |
| API7 Server Side Request Forgery | No user-controlled URLs; allowlisted outbound hosts (PayPal API, LLM provider, status page); cert URL restricted |
| API8 Security Misconfiguration | Config validator; secure headers; CORS allowlist; no debug in hosted; IaC review |
| API9 Improper Inventory Management | OpenAPI as inventory; route-registry check in CI; no undocumented routes |
| API10 Unsafe Consumption of APIs | Treat PayPal/LLM/status responses as untrusted: schema-validate, size-cap, timeouts, injection-safe handling |
| Web: XSS | React escaping; no raw HTML; CSP; untrusted text flagged |
| Web: CSRF | Token + SameSite |
| Web: Clickjacking | `frame-ancestors 'none'` |
| Web: Injection (SQL) | sqlc parameterised queries only; no string-built SQL |
| Web: Vulnerable components | Dependency scanning (§13) |
| Web: Logging failures | Structured audit + alerts |

---

## 12. Specific attack classes

| Class | Controls |
|---|---|
| **SSRF** | No endpoint fetches user-supplied URLs; `PAYPAL-CERT-URL` never fetched directly (PayPal verify API used); status-page fetch restricted to configured host; HTTP client denies redirects to other hosts and private IP ranges |
| **Injection (SQL/command/template)** | Parameterised SQL; no shell-outs with user input; no server-side templates with user input |
| **XSS** | React default escaping; sanitise nothing into HTML; JSON viewer renders text nodes; CSP |
| **CSRF** | Double-submit + SameSite=Lax; `Origin`/`Referer` check on unsafe methods |
| **Auth** | Argon2id for demo passwords; generic errors; throttling; session revocation on logout/consent change |
| **IDOR** | UUIDs are not secrets; every access is policy-checked; 404 uniform; automated cross-tenant/cross-role tests (CHAOS-06) |
| **Replay** | Webhook dedupe + monotonic reducer; idempotency keys; CSRF tokens per session |
| **Webhook spoofing** | See §6; reject/503 semantics |
| **Rate limits** | §10 of doc 04; per IP/session/identity/case; `Retry-After` |
| **Secrets** | §13 |
| **Dependencies** | §13 |
| **Prompt injection** | doc 07 §11 |
| **Open redirect** | PayPal return/cancel URLs fixed server-side, not from client input; `next` params validated against same-origin allowlist |
| **File upload** | None in MVP |
| **Timing/user enumeration** | Constant-time compare; uniform responses for login/invite/share |

---

## 13. Secrets and dependency management

**Secrets**
- Env/secret-store only; `.env*` git-ignored (`.env.example` committed with placeholders only).
- gitleaks pre-commit + CI (and scan git history once at project start).
- Browser bundle contains only `VITE_*` public values (none secret). A CI check greps `dist/` for known secret **names** and high-entropy strings.
- Separate keys for dev/demo; rotate on suspected exposure; document rotation (doc 09).
- Never place secrets in screenshots, demo videos, issue trackers, or chat. If a secret is exposed, treat as compromised and rotate immediately.

**Dependencies**
- Lockfiles committed (`go.sum`, `package-lock.json`); Dependabot/Renovate with review.
- `govulncheck`, `npm audit` (or equivalent), `osv-scanner`; container image scan (e.g. Trivy); SBOM generation (e.g. `syft`) as release artefact.
- Pin CI actions by version/SHA; minimal base images; non-root containers; read-only root filesystem where possible.
- New dependency checklist: licence compatible with the chosen open-source licence (**VERIFY** hackathon requirements), maintenance status, transitive footprint.

---

## 14. Encryption and data minimisation

| Topic | Design |
|---|---|
| In transit | TLS for all external traffic; TLS to managed Postgres where available |
| At rest | Managed disk/DB encryption (platform); **application-level AES-256-GCM** (`ENCRYPTION_KEY`) for raw provider payloads/webhook bodies and share artefacts; key rotation plan with key id stored per ciphertext |
| Hashing | SHA-256 for payload integrity; session ids and share tokens stored **hashed**; Argon2id for passwords |
| Minimisation | Never store card PAN/CVV, bank credentials, OTPs, buyer passwords; do not store PayPal buyer PII beyond what the case needs; store redacted extracts by default; raw payload optional and short-retention |
| Masking | Ids masked in UI/exports (`3AB…9F2`); emails masked for the other party |
| LLM | Only redacted projections; no names/emails by default (role labels) |
| Telemetry | Aggregates; no per-user identifiers in metrics labels |

---

## 15. Privacy matrix

*(Default deny. ✔ = may see; ◐ = masked/summary; ✘ = never.)*

| Data element | Sender | Receiver (consented) | Merchant analyst/admin | Platform operator (JIT) | AI agents |
|---|:-:|:-:|:-:|:-:|:-:|
| Case reference, amount, currency | ✔ | ✔ | ✔ | ◐ | ✔ (role-scoped) |
| Provider capture status | ✔ | ✔ | ✔ | ✔ | ✔ |
| Provider ids (order/capture) | ◐ masked | ◐ masked | ✔ | ✔ | ◐ masked |
| App timeout/trace summary | ✔ plain | ✘ (seller-side summary only) | ✔ | ✔ | sender/ops snapshot |
| Webhook delivery diagnostics | ✘ | ✘ | ✔ | ✔ | ops only |
| Sender funding/issuer data | ✘ (not stored) | ✘ | ✘ | ✘ | ✘ |
| Receiver account verification/limits | ✘ | ✔ (own) | restricted | ✘ | receiver-only |
| Bank/withdrawal/settlement | n/a (`NOT_OBSERVABLE`) | n/a | n/a | n/a | n/a |
| Other party's identity | ◐ | ◐ | ✔ | ✔ | ◐ (role labels) |
| Evidence items | role-class subset | `RECEIVER`-class subset | tenant | tenant (JIT) | per-role subset |
| Raw webhook/provider payloads | ✘ | ✘ | ✘/restricted | restricted + logged | ✘ |
| Audit trail | ✘ | ✘ | admin | JIT + logged | ✘ |
| Agent prompts/outputs | ✘ | ✘ | ops (redacted) | ✔ (redacted) | n/a |

---

## 16. Consent, retention, deletion

**Consent (design):** receiver invitation lists exact fields; acceptance creates a versioned `consent_grant`; notifications are a *separate* consent; revocation is immediate for access, tokens, SSE streams, and pending notices; delivered messages cannot be recalled (stated plainly). Sender consent to share is implied by creating/inviting but fields remain bounded by the allowlist. Consent records are audited.

**Retention (proposed defaults — subject to legal review, **VERIFY**):** case metadata per merchant policy; raw payloads short (e.g. days–weeks, configurable); redacted extracts with the case; agent artefacts sanitised; share links expiring; telemetry aggregated. Defaults are **DESIGN DECISION** placeholders, not legal determinations.

**Deletion by design:** user request → validation → delete/anonymise eligible data + encrypted artefacts + caches → tombstone; some audit/legal records may need retention — only counsel can define; UI never promises instant erasure of all logs/backups. Backups age out per backup retention (stated in privacy notice).

**Data subject requests:** export (redacted case data) and deletion request flows (R43/R44) produce an auditable job; the process, identity verification, and statutory timelines are **VERIFY with counsel** for any real deployment.

**Notices:** the privacy notice must disclose LLM vendor processing (redacted), hosting provider, retention, and that PayBridge is independent of PayPal. Drafting a compliant notice is **not** provided here.

---

## 17. India payment eligibility and compliance (VERIFY)

> **Not legal advice. Nothing here asserts that any party is eligible or compliant.**

- **SOURCE fact (M20 §12.6, from PayPal India help at time of research):** international receivables for Indian accounts require identity confirmation, PAN, an Indian bank account and a purpose code, with additional information depending on business type; account-type rules and cross-border products vary. **VERIFY current PayPal India documentation** before any statement.
- Open verification list (each needs a current official source and professional advice before any claim): whether the owner's intended use (service-export receipts via PayPal) is permitted for the relevant account type; purpose-code requirements; limits; documentation/retention duties; payment-aggregator or cross-border payment regulations and licensing questions for **PayBridge itself** (an observer that does not move or hold funds — **VERIFY with counsel**); data-protection law obligations (e.g. India's DPDP Act and rules — **VERIFY current status and obligations**); US/other-jurisdiction privacy and sanctions considerations for multi-country users.
- Product behaviour: A01 returns `UNKNOWN_UNTIL_VERIFIED` with the documented prerequisites and source/time, never "allowed for everyone" (doc 07 A01).
- **PayBridge does not collect KYC documents, PAN, or bank details.**
- Hackathon use is sandbox-only; **production use requires a separate compliance assessment** (not part of this blueprint).

---

## 18. Unobservable provider and bank data

| Not observable by this project | UI/API behaviour |
|---|---|
| Recipient bank settlement/posting | `receiver_state=NOT_OBSERVABLE`; "Bank settlement not verified" |
| PayPal internal queues, risk models, global traffic | Never claimed; A02 limited to own telemetry + public status (advisory) |
| Issuer/card decline rationale beyond provider-exposed codes | Not inferred (T07) |
| Receiver KYC outcome | Not inferred; guidance to official account |
| Withdrawal status without a supported, consented connector | `NOT_OBSERVABLE` |
| Sandbox ≠ production behaviour | Always labelled |

---

## 19. Incident response

| Phase | Actions |
|---|---|
| **Prepare** | Runbooks (doc 09); on-call contact (owner); alert routes; secrets inventory; rotation procedures |
| **Detect** | Alerts: signature-failure spike, cross-tenant denial spike, secret-scan hit, unusual LLM spend, dead-letter growth, audit-chain failure |
| **Triage** | Severity: S1 suspected secret/data exposure or wrong financial claim; S2 integrity failure; S3 availability |
| **Contain** | Revoke/rotate keys; disable LLM (`LLM_ENABLED=false`); disable Fault Proxy/demo endpoints; revoke sessions/shares; block IPs |
| **Eradicate/recover** | Patch; restore from backup if needed; re-run reconciliation for affected cases |
| **Notify** | Affected users/partners as required by law/contract (**VERIFY with counsel**); PayPal if their credentials were exposed (rotate in console) |
| **Learn** | Blameless post-mortem; add regression test; update threat model |
Specific playbooks: *secret exposed*, *wrong Passport statement published* (withdraw version, publish correction, audit), *cross-role leak* (disable lens, purge notifications, review DLP), *webhook forgery attempt*, *provider incident*.

---

## 20. Release gates

A release (including the hackathon demo deployment) requires **all** of:

| Gate | Evidence |
|---|---|
| G-SEC-1 | gitleaks clean on tree **and** history |
| G-SEC-2 | `govulncheck`/`npm audit`/`osv-scanner` — no unaddressed high/critical |
| G-SEC-3 | IDOR/RBAC/tenant matrix passes (CHAOS-06) |
| G-SEC-4 | Forged/duplicate/out-of-order webhook tests pass (T03–T05) |
| G-SEC-5 | Role-privacy diff tests pass (0 leaks, AO-08) |
| G-SEC-6 | CSRF, headers, CSP checks pass |
| G-SEC-7 | Config guard proves `PRODUCTION` disabled |
| G-PRV-1 | Privacy notice reviewed by owner; limitations stated; counsel review flagged as not done |
| G-PAY-1 | Zero mutating PayPal paths reachable from AI/worker (arch test) |
| G-PAY-2 | AO-02 (one order, one capture) evidenced |
| G-CLAIM-1 | Honest-claims checklist (doc 10) signed |
| G-OPS-1 | Backup restore drill performed and recorded |

---

## 21. `.env.example`

```dotenv
# ===== PayBridge AI — .env.example =====
# Copy to .env (git-ignored). NEVER commit real values. Placeholders only.

# --- Application ---
APP_ENV=local                         # local | hosted
LOG_LEVEL=info
PUBLIC_BASE_URL=http://localhost:8080 # Public HTTPS URL when hosted (webhook endpoint base)
WEB_ORIGIN=http://localhost:5173      # Allowed CORS origin for the web app
DEMO_MODE=true                        # enables demo/scenario routes (never in non-sandbox)
FAULT_PROXY_ENABLED=false             # CHAOS-01 response-drop; sandbox/dev only

# --- Database / cache ---
DATABASE_URL=postgres://paybridge:__SET_ME__@localhost:5432/paybridge?sslmode=disable   # SECRET
REDIS_URL=                            # optional; rate limiting only (never financial truth)

# --- Security ---
SESSION_SECRET=__SET_ME_LONG_RANDOM__       # SECRET: cookie/session signing
ENCRYPTION_KEY=__SET_ME_BASE64_32_BYTES__   # SECRET: AES-256-GCM key for encrypted evidence payloads
ENCRYPTION_KEY_ID=k1                         # key rotation identifier

# --- PayPal (SANDBOX ONLY) ---
PAYPAL_MODE=real                      # real | fake. "fake" routes to the local fake PayPal and stamps PAYPAL_SIMULATOR (never REAL_SANDBOX)
PAYPAL_ENV=sandbox                    # must equal "sandbox" — app refuses to start otherwise
PAYPAL_API_BASE=https://api-m.sandbox.paypal.com
PAYPAL_CLIENT_ID=__SET_ME__           # sandbox REST app client id
PAYPAL_CLIENT_SECRET=__SET_ME__       # SECRET: server-side only
PAYPAL_WEBHOOK_ID=__SET_ME__          # sandbox webhook id for signature verification
PAYPAL_HTTP_TIMEOUT_MS=15000

# --- LLM (worker only) ---
LLM_ENABLED=true                      # false => deterministic fallback everywhere
LLM_PROVIDER=__SET_ME__               # e.g. anthropic | openai | gemini (adapter must exist)
LLM_MODEL=__SET_ME__
LLM_API_KEY=__SET_ME__                # SECRET: worker only
LLM_DAILY_USD_CAP=5                   # spend ceiling for the demo

# --- Frontend (public values only; NEVER put secrets in VITE_*) ---
VITE_API_BASE=/api/v1
VITE_ENV_LABEL=REAL SANDBOX
```
**Missing-variable protocol:** the config validator prints **only names** of missing/invalid required variables, e.g. `missing required config: PAYPAL_CLIENT_SECRET, PAYPAL_WEBHOOK_ID`, and exits (API/worker). If only `LLM_*` is missing, the app starts with `LLM_ENABLED=false` and a visible "AI explanation unavailable" state.

---

## 22. Security tests

*(All **NOT RUN**; mapped to doc 09 requirement IDs.)*

| ID | Test | Method |
|---|---|---|
| SEC-01 | Forged webhook rejected; no state change; metric increments | Integration with fake verifier + real-sandbox check |
| SEC-02 | Verification outage ⇒ 503, nothing acknowledged | Fault injection |
| SEC-03 | IDOR matrix: every `{id}` route × {other tenant, other role, anonymous} ⇒ uniform 404 | Generated test from OpenAPI |
| SEC-04 | RLS enumeration: every tenant table has FORCE RLS + policies | SQL test |
| SEC-05 | CSRF: unsafe routes fail without token; cross-site form POST fails | HTTP tests |
| SEC-06 | Session fixation/rotation/timeout/logout revoke | Auth tests |
| SEC-07 | Role matrix (doc 04 §8.2) exhaustively asserted | Table-driven |
| SEC-08 | Strict schema: mass-assignment attempts (extra fields) rejected | Contract fuzz |
| SEC-09 | SSRF: no route accepts URLs; outbound client denies private ranges/redirects | Unit + egress test |
| SEC-10 | XSS payload in memo renders inert (DOM + CSP) | Playwright |
| SEC-11 | Prompt injection suite (doc 07 §15.3) | Eval harness |
| SEC-12 | Secrets: gitleaks tree+history; bundle scan; log scan for tokens | CI |
| SEC-13 | Rate limits: login, checkout, webhook buckets | Load/unit |
| SEC-14 | Config guard: refuses non-sandbox host/env | Unit |
| SEC-15 | Audit chain tamper detection | Integration |
| SEC-16 | Consent revocation: access, stream, share tokens, pending notices (CHAOS-09) | E2E |
| SEC-17 | Privacy diff: sender vs receiver payload/DOM have zero cross-party fields | E2E diff |
| SEC-18 | Dependency/container scans | CI |
| SEC-19 | Headers/CSP/CORS conformance | HTTP tests |
| SEC-20 | Reader interface has no mutating methods; AI packages cannot import mutators | Arch test |

*End of document 08.*
