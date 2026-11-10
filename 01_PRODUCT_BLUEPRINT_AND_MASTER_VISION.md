# 01 — PayBridge AI: Product Blueprint and Master Vision

| Field | Value |
|---|---|
| Document | `01_PRODUCT_BLUEPRINT_AND_MASTER_VISION.md` |
| Version | 1.0.0 (blueprint draft for approval) |
| Date | 9 November 2026 |
| Status | **DRAFT — awaiting owner approval. No code exists. Nothing in this document is evidence of a working system.** |
| Audience | Product owner, Claude Code / Antigravity coding agents, human developers, hackathon reviewers |
| Primary source | `PayBridge_AI_Complete_End_to_End_M20_Master_Plan.md` (v3.0, 9 Nov 2026), referred to below as **the M20 plan** |
| Depends on | Nothing (this is the root document) |
| Related | 02 Implementation Prompt · 03 Architecture · 04 Backend · 05 Frontend · 06 UI/UX · 07 AI Agents · 08 PayPal/Security · 09 Testing/DevOps · 10 Roadmap/Demo |

> **Change note (2026-11-09, Account 2 audit):** decisions D-11…D-16 and conflicts C-09…C-14 added to §19; §24 addendum added. No other text in this document was altered. Support files: `INDEX.md`, `BLUEPRINT_QUALITY_AUDIT.md`.

> **One payment. Two people. One verified story.**
> North-star: **"Know the last confirmed step, understand what is still unknown, and never retry blindly."**

---

## Table of contents

1. [Purpose, source hierarchy and how to read this blueprint set](#1-purpose-source-hierarchy-and-how-to-read-this-blueprint-set)
2. [Label legend](#2-label-legend)
3. [Mission, vision, pitch](#3-mission-vision-pitch)
4. [The problem](#4-the-problem)
5. [Golden use case](#5-golden-use-case)
6. [Personas](#6-personas)
7. [Jobs-to-be-done](#7-jobs-to-be-done)
8. [Stakeholders](#8-stakeholders)
9. [Product principles](#9-product-principles)
10. [Goals and non-goals](#10-goals-and-non-goals)
11. [Assumptions and constraints](#11-assumptions-and-constraints)
12. [Canonical vocabulary (shared by all ten documents)](#12-canonical-vocabulary-shared-by-all-ten-documents)
13. [Differentiators and honest prior art](#13-differentiators-and-honest-prior-art)
14. [Capability inventory: MVP vs later, with truth labels](#14-capability-inventory-mvp-vs-later-with-truth-labels)
15. [User journeys](#15-user-journeys)
16. [Scope boundaries](#16-scope-boundaries)
17. [Measurable acceptance outcomes](#17-measurable-acceptance-outcomes)
18. [Risks](#18-risks)
19. [Decision log and source conflicts](#19-decision-log-and-source-conflicts)
20. [Phase roadmap (summary)](#20-phase-roadmap-summary)
21. [The full end-to-end product story](#21-the-full-end-to-end-product-story)
22. [Claims policy](#22-claims-policy)
23. [Hackathon fit](#23-hackathon-fit)
24. [Repository and environment audit finding](#24-repository-and-environment-audit-finding)

---

## 1. Purpose, source hierarchy and how to read this blueprint set

**Purpose.** This document fixes *product truth*: what PayBridge AI is, what it is not, who it serves, which words mean what, what the MVP contains, and how success is measured. Every later document must remain consistent with it. When another document appears to contradict this one, this one wins unless a decision-log entry (§19) says otherwise.

**Source hierarchy** (highest wins):

1. Explicit instructions from the project owner in the current session.
2. The M20 plan (Parts I–II). Part III (the archived original specification) is lineage only.
3. This blueprint set (documents 01–10).
4. Current official PayPal documentation, which must be re-checked before implementation (**VERIFY**).
5. Everything else (library defaults, general practice).

Where the M20 plan conflicts with itself, the conflict is **recorded in §19.2**, not silently resolved.

**The ten documents form one system:**

| # | Document | Owns |
|---|---|---|
| 01 | Product Blueprint (this) | Product truth, vocabulary, scope, acceptance outcomes |
| 02 | Master Implementation Prompt | The reusable coding-agent operating contract |
| 03 | System Design and Architecture | C4, sequences, state machines, ADRs |
| 04 | Backend Engineering Spec | Go layout, schema, APIs, reconciliation, safe-action engine |
| 05 | Frontend Spec | React/TS structure, data flow, real vs fixture screens |
| 06 | UI/UX Design System | Visual language, 16 screens, custom visualisations |
| 07 | AI Multi-Agent Spec | A01–A10, orchestrator, schemas, prompts, evals |
| 08 | PayPal, Security, Privacy, Compliance | Integration contract, threat model, privacy matrix |
| 09 | Testing, DevOps, Deployment | Setup, simulator, CI/CD, runbooks, release gates |
| 10 | Roadmap, Acceptance, Demo | Phases, epics, storyboard, claims checklist |

---

## 2. Label legend

Every non-obvious statement in this blueprint set carries a label.

| Label | Meaning |
|---|---|
| **SOURCE** | Retained from the M20 plan or the owner's master prompt. |
| **DESIGN DECISION** | A choice made by this blueprint. Reversible via the decision log. |
| **ASSUMPTION** | Believed but unproven; needs testing. |
| **VERIFY** | Depends on current external documentation or account eligibility; must be re-checked before coding. |
| **SANDBOX** | Intended to be genuine PayPal sandbox behaviour — *only counts once actually executed and evidenced*. |
| **SIMULATED** | Local fixture/simulator behaviour. Never presented as a real provider event. |
| **PRODUCTION** | Future work. Not in the MVP. Never claimed. |

Status vocabulary used for capabilities: **PLANNED** (designed, not built), **BUILT** (code exists), **TESTED** (automated test executed and result recorded), **DEMONSTRATED** (shown live with evidence). At the time of writing every capability is **PLANNED**.

---

## 3. Mission, vision, pitch

**Product category (SOURCE).** PayBridge AI is an **evidence-first cross-border payment intelligence, investigation and safe-recovery-coordination platform**. It is **not** a payment processor, a bank, an escrow, a wallet, a generic payment tracker, or a chatbot.

**Mission (SOURCE).** Reduce preventable repeated payments, reduce confusion and support effort, and coordinate the next compliant action — without letting AI invent transaction facts or control money.

**Vision (SOURCE).** Every international commercial payment should have a trustworthy, privacy-respecting explanation from initiation through the last observable provider stage, even when systems fail or disagree.

**Brand values (SOURCE).** Precision over optimism. Privacy over oversharing. Verification over speculation. Helpfulness without unauthorised financial decisions. Clarity for non-technical users.

**Pitch variants (SOURCE, lightly edited).**

- *10 seconds:* "PayBridge AI explains a confusing international payment to both parties, backed by actual payment evidence — and prevents dangerous blind retries."
- *30 seconds:* "A US business pays a freelancer abroad. PayPal captures the money, but the merchant app times out. The sender thinks it failed; the receiver is unsure; paying again could duplicate the payment. PayBridge retrieves the original provider record, reconstructs the event sequence, tells each side what is confirmed and what is unknown, and recommends the next safe step. Every claim links to evidence."
- *Enterprise:* "An evidence-first payment incident workspace that connects provider observations, webhook quality, case ownership and role-aware AI explanation, designed to reduce support ambiguity and improve merchant operational reliability."

**Tagline:** *One payment. Two people. One verified story.*
**Honest footer (SOURCE):** *Independent project using PayPal developer APIs; not affiliated with or endorsed by PayPal.*

---

## 4. The problem

A US customer pays a service business abroad. At least five different systems each hold a *different* answer to "what is the payment status?":

1. The buyer's checkout confirmation (browser).
2. The seller's order page (merchant database).
3. PayPal's transaction record (the provider).
4. The seller's account availability (receiver-side).
5. The recipient bank's settlement (usually **not observable** through any API available to this project).

People treat these as one "status". They are five. Failure happens when one system knows more than another, or when an event is delayed, duplicated, or lost.

**Failure examples (SOURCE, M03 §3.2):**

| Situation | What may really have happened | Dangerous reading | What PayBridge does |
|---|---|---|---|
| Server times out during capture | Capture succeeded; response lost | "Pay again" | Reconcile original order/capture; block retry |
| Webhook delayed or duplicated | Payment fine; notification late | "Provider lost funds" | Verify, dedupe, compare timestamps |
| Payment pending | Provider has not finalised | "Receiver got paid" | Preserve pending, show authoritative state |
| Merchant DB says failed | Local write crashed after provider success | "Refund now" | Compare local vs provider before advising |
| Receiver sees no bank credit | Bank stage is outside API visibility | "Capture failed" | Mark bank settlement NOT_OBSERVABLE |
| Public incident reported | Some component degraded | "Every payment failed" | Separate incident advisory from case proof |
| AI hallucinated a cause | Insufficient evidence | "Provider rejected receiver" | Block unsupported claims; escalate |

**Systemic root cause (SOURCE, five-whys):** the application models *uncertainty* and *transaction identity* incorrectly — it equates "my request timed out" with "the provider declined", and on retry it creates a **new** order instead of reconciling the **old** one. The fix is correct modelling plus evidence, not a faster provider.

**Failure families (SOURCE, M03 §3.3)** — used as a taxonomy across the platform: F01 Client · F02 Merchant application · F03 Webhook delivery · F04 Provider state · F05 User action · F06 Currency/accounting · F07 Data provenance · F08 External rails · F09 Privacy/security · F10 Provider health. A case can have several contributing families; never assume a single cause without authoritative evidence.

**What research validates (SOURCE):** [VERIFIED in the M20 plan's register] PayPal provides Orders/Payments APIs, webhooks, idempotency headers and sandbox testing, and documents pending/asynchronous outcomes. **[ASSUMPTION]** Merchants lose measurable time and money to exactly this gap — *no market sizing or quantified impact exists*, and this blueprint set invents none.

---

## 5. Golden use case

> **Capture succeeds at the provider; the merchant/app confirmation times out; PayBridge reconciles the original transaction, explains the evidence to both parties, and prevents a blind duplicate payment.**

This is the single flagship workflow. Everything in the MVP exists to make it real, demonstrable, and defensible. Protect it above all else.

### 5.1 The scripted timeline (SOURCE M07 §7.1, normalised to blueprint vocabulary)

| Step | Actor | Event | State after |
|---|---|---|---|
| 1 | Sender | Logs in; selects an open invoice (illustrative USD 2,000.00) | — |
| 2 | PayBridge | Creates `payment_case` (UUID, tenant, invoice ref, amount, currency, purpose, route) | `case_status=OPEN` |
| 3 | A01 / route registry | Readiness: `supported`, `unsupported` or `unknown` with provenance | provider_state `NOT_CREATED` |
| 4 | Go API | Persists attempt + `operation_request_id`; calls PayPal create order; persists `provider_order_id` | provider `CREATED`; attempt `ORDER_CREATED` |
| 5 | Buyer | Approves in PayPal's hosted flow | provider `APPROVED`; attempt `APPROVED` |
| 6 | Go API | Calls capture **once** under a stable request id | attempt `CAPTURE_REQUESTED` |
| 7 | Fault proxy **(SANDBOX + controlled fault)** | Drops *only PayBridge's own capture response* after PayPal has captured | application `API_TIMEOUT`; attempt `CAPTURE_UNKNOWN`; provider **UNKNOWN (last confirmed = APPROVED)** |
| 8 | UI | Shows "Payment outcome unknown. Please do not pay again." | — |
| 9 | Reconciliation worker | `GET` original order, then `GET` capture. **Never** creates a new order/capture | provider `CAPTURE_COMPLETED` (verified GET) |
| 10 | Webhook path | Signed `PAYMENT.CAPTURE.COMPLETED` arrives before/after/twice | deduped; no regression |
| 11 | Orchestrator | Routes A05 → A03 → A06 → A08 → A09 (A04 only with receiver consent) | investigation complete |
| 12 | Passport | Sender: "Provider capture completed; our confirmation was lost; do not pay again." Receiver: "Seller-side capture recorded; check official account activity for funds." | passport v2 |
| 13 | Merchant analyst | Approves *internal* reconciliation of the merchant order (no PayPal call) | application `SYNCHRONISED` |
| 14 | System | Verifies resolution predicate; case `RESOLVED` | receiver state remains `NOT_OBSERVABLE` |

**Required visible proofs:** a real SANDBOX order id and capture id (redacted as needed); the provider `GET` result; the app timeout trace; a ledger showing **exactly one** order and **exactly one** capture; sender and receiver messages that are visibly different; zero money-moving AI actions.

### 5.2 The five demonstrations (SOURCE M17 §17.3)

| Demo | Story | Mode |
|---|---|---|
| A | Verified normal sandbox checkout; green "PayPal capture confirmed"; bank lane "not observed" | SANDBOX |
| **B** | **Hero: payment succeeded, our response disappeared** | **SANDBOX + app-side fault** |
| C | Real webhook, then duplicate replay (labelled replay); one business effect | SANDBOX + LOCAL_SIMULATION |
| D | Cross-border readiness uncertainty (US→IN) shown as `Unknown until account verification` | SIMULATED readiness |
| E | LLM disabled; payment truth and deterministic guidance still work | SANDBOX + degraded AI |

---

## 6. Personas

| ID | Persona | Situation | Success looks like |
|---|---|---|---|
| **S** | US business payer (accounts payable) | Paying a USD 2,000 invoice for approved design services | No duplicate charge; knows what happened; clear receipt |
| **R** | Indian export-service freelancer | Completed work; wants to know what PayPal actually confirms | Clarity without seeing payer-private data; no false "arrived" claims |
| **O** | Merchant support operator | Handles many cases (illustrative 20–200/day — **not measured**) | Fast triage; reproducible explanation; fewer escalations |
| **D** | Integration developer | Debugging lost responses, webhook signatures, env mismatches | Deterministic replay; evidence ids; simulator |
| **A** | Security/compliance reviewer | Needs consent, minimisation, tenant isolation, audit | No unapproved money-moving actions; no cross-party leaks |
| **J** | Hackathon judge (project-specific) | Has minutes; must run or watch the product | Understands the problem in under a minute; sees real PayPal sandbox evidence |

*Persona J is a blueprint addition (DESIGN DECISION) because the immediate success condition is a judged demonstration.*

**Buyer vs user (SOURCE M05).** Economic buyer: merchant founder / ops manager / payment engineering lead. Daily users: S, R, O. Technical buyer: lead developer. Influencers: customer success, compliance, finance reconciliation.

---

## 7. Jobs-to-be-done

1. *When checkout times out,* help me find the original outcome so I don't pay twice.
2. *When my client says they paid,* help me see the last confirmed stage without misleading bank-arrival claims.
3. *When a ticket mentions an international payment,* help me reconstruct what our app and the provider recorded.
4. *When the payment is truly denied or needs action,* tell the right person which official step to take.
5. *When an AI explains a payment,* show the exact evidence and uncertainty so I can trust or challenge it.

---

## 8. Stakeholders

| Stakeholder | Interest | Constraint they impose |
|---|---|---|
| Project owner (Gandhi) | Win the PayPal hackathon; ship a credible product | Honest claims; reproducible demo |
| PayPal (platform) | Developer-platform usage | Respect API terms; no implied endorsement; no brand misuse |
| Hackathon judges | Technology, design, impact, innovation, presentation (equal weight per M20 §20.4 — **VERIFY** against official rules) | Working demo, public repo with OSS licence, short video |
| Senders / receivers | Truth and privacy | Consent, role-scoped views |
| Merchant ops / developers | Operational reliability | Auditability, simulators |
| Privacy/legal reviewers | Compliance posture | No legal-compliance claims; counsel for production |
| LLM vendor | Data processing | Minimised, redacted prompts |

---

## 9. Product principles

*Ten non-negotiables (SOURCE M02 §2.2), each with its enforcement locus.*

| # | Principle | Enforced by |
|---|---|---|
| 1 | **Truth precedence** — verified provider state beats local cache; a user report is not provider evidence | Reducer + evidence tiers (doc 04) |
| 2 | **Unknown is a valid state** — a timeout is not evidence of failure | State machine + UI copy rules (docs 03, 06) |
| 3 | **No blind retry** — reconcile the original request before any new financial op | Attempt lock + safe-action engine (doc 04) |
| 4 | **One transaction, different authorised views** | Projection layer + consent (docs 04, 08) |
| 5 | **Evidence per claim** | Claim validator (docs 04, 07) |
| 6 | **Deterministic control plane** — models interpret; Go decides | Architecture boundary (doc 03) |
| 7 | **Security as architecture** — no secrets to agents; webhook verification; tenant boundaries first | Doc 08 |
| 8 | **Real vs simulated separation** | Environment enum, distinct UI treatment, distinct id namespaces |
| 9 | **Measure what matters** — controlled test cases only, no extrapolation to PayPal production | Doc 09/10 |
| 10 | **Modular monolith first** | ADR-005 (doc 03) |

**Product-design corollaries (DESIGN DECISION):**

- The centre of the product is the **evidence-based payment case**, never a generic finance dashboard.
- The UI never shows a progress bar that implies completion beyond the last observable stage.
- *Verified* green is reserved for provider-verified facts. Unverified "sent" stages never get optimistic green.
- Red is reserved for **confirmed** errors or security risk. Unknown is amber, never red.
- Financial verbs are never conflated: *approve, capture, refund, cancel, confirm, reconcile* each mean one thing.

---

## 10. Goals and non-goals

### 10.1 Goals (MVP)

| ID | Goal |
|---|---|
| G1 | A **real PayPal sandbox** Orders v2 create → approve → capture flow with stored provider ids |
| G2 | Reproduce the golden failure (provider success + lost app response) and **reconcile without a second financial effect** |
| G3 | A webhook inbox with signature verification, dedupe, replay protection, out-of-order safety |
| G4 | The **Payment Truth Passport**, versioned and role-scoped, every confirmed claim linked to evidence |
| G5 | The **Truth Gap Graph** visualising contradictions between browser, merchant app, provider and webhook |
| G6 | **Two-Lens** sender and receiver explanations with consent and access control |
| G7 | A deterministic **Next Safe Action Engine** that blocks dangerous actions |
| G8 | The **Counterfactual Recovery Mirror** — a clearly labelled simulation lab |
| G9 | Ten bounded agents (A01–A10) under a Go orchestrator, read-only, schema-validated, with deterministic fallback |
| G10 | A premium, accessible UI (doc 06) and a hosted, reproducible demo |

### 10.2 Non-goals (SOURCE M06 §6.3, M16 §16.4)

Consumer remittances; bypassing KYC/AML; universal bank-transfer tracking; internal PayPal queue or traffic data; guaranteeing FX, fees or ETAs; private account lookups without permission; escrow; **autonomous money movement**; compliance eligibility decisions; presenting synthetic records as genuine provider records; storing card PAN/CVV or banking credentials; legal or financial advice; market validation claims.

---

## 11. Assumptions and constraints

### 11.1 Assumptions

| ID | Assumption | Test |
|---|---|---|
| AS-01 | A PayPal sandbox business + personal account pair can run Orders v2 create/approve/capture | Executed in Phase 1 (**VERIFY**) |
| AS-02 | Sandbox webhooks can reach a publicly reachable HTTPS endpoint during the demo | Executed in Phase 2 |
| AS-03 | A US→IN business route may **not** be configurable in the sandbox | If unsupported, real proof uses a supported sandbox config; the US→IN story lives in a labelled simulation layer (SOURCE M12 §12.6) |
| AS-04 | A hosted LLM with structured/JSON output is reachable via owner-provided API key | Executed in Phase 5 |
| AS-05 | Merchants will pay for this | **Unvalidated.** No claim is made |
| AS-06 | Sender/receiver sharing creates enough value to justify consent complexity | **Unvalidated** |

### 11.2 Constraints

- **No invented credentials.** If a key is missing the agent reports the exact variable name and setup steps.
- **No secret values** in documents, logs, screenshots, commits, or prompts.
- Sandbox only; the `PRODUCTION` environment value is reserved and **hard-disabled** in the MVP (doc 08).
- Idempotency semantics differ by endpoint and duration (**VERIFY**) — PayBridge persists its own keys regardless.
- Provider API schemas must be re-checked in official docs before coding (**VERIFY**).
- Stack baseline (SOURCE + owner prompt): Go modular monolith; React + TypeScript + Vite; PostgreSQL; Redis optional and never financial truth; PayPal REST + sandbox; provider-neutral hosted LLM; Docker Compose locally; OpenAPI as contract source of truth; no Python in the MVP.

---

## 12. Canonical vocabulary (shared by all ten documents)

This section is **normative**. Other documents reference it by name rather than redefining it.

### 12.1 Environments (`environment`)
`REAL_SANDBOX` · `PAYPAL_SIMULATOR` · `LOCAL_SIMULATION` · `PRODUCTION` (reserved, disabled).

### 12.2 The three orthogonal state dimensions
A single global status is dangerously lossy (SOURCE ADR-002). Track a vector `state = (provider, application, receiver)` with a timestamp and source per member.

| Dimension | Values |
|---|---|
| **provider_state** | `NOT_CREATED`, `CREATED`, `APPROVED`, `AUTHORISED`, `CAPTURE_PENDING`, `CAPTURE_COMPLETED`, `CAPTURE_DENIED`, `REFUNDED`, `PARTIALLY_REFUNDED`, `UNKNOWN` (+ a stored `last_confirmed_provider_state`) |
| **application_state** | `SYNCHRONISED`, `WAITING_FOR_WEBHOOK`, `API_TIMEOUT`, `RECONCILING`, `EVENT_GAP`, `INCONSISTENT` |
| **receiver_state** | `NOT_APPLICABLE`, `NOT_OBSERVABLE`, `REPORTED_PENDING`, `VERIFIED_AVAILABLE`, `WITHDRAWAL_IN_PROGRESS`, `VERIFIED_BANK_CREDITED` (the last three need separately authorised evidence; none is reachable in the MVP) |

### 12.3 Case workflow (`case_status`)
`OPEN → OBSERVING → INVESTIGATING → GUIDANCE_READY → AWAITING_USER_ACTION → VERIFYING_RESOLUTION → RESOLVED`; exception states `NEEDS_HUMAN_REVIEW`, `PAUSED_INSUFFICIENT_EVIDENCE`, `CLOSED_UNRESOLVED`, `ARCHIVED`, and `REOPENED` (see decision D-04). `RESOLVED` means the case's *stated objective* was verified — not that every downstream stage completed.

### 12.4 Attempt status (checkout machine, SOURCE M13 §13.2)
`ORDER_CREATING`, `ORDER_CREATED`, `ORDER_UNKNOWN`, `APPROVAL_PENDING`, `APPROVED`, `CANCELLED`, `CAPTURE_REQUESTED`, `CAPTURE_COMPLETED`, `CAPTURE_PENDING`, `CAPTURE_DENIED`, `CAPTURE_UNKNOWN`, `REVIEW_REQUIRED`, `APPLICATION_RECONCILED`, `APPLICATION_MISMATCH`.

### 12.5 Evidence
- `source_trust`: `PROVIDER_VERIFIED`, `LOCAL_OBSERVED`, `PUBLIC_SOURCE`, `USER_REPORTED`, `SIMULATED`.
- **Evidence tiers** (SOURCE M14 §14.5, precedence high → low): **E1** verified authorised provider GET · **E2** verified signed provider webhook · **E3** merchant-controlled application logs · **E4** official published rules and public status · **E5** consented user testimony · **E6** simulation / model output.
- Strength is **field-specific**, never one magic confidence number.
- **Claim kinds:** `confirmed` · `inference` (hypothesis) · `unknown`.

### 12.6 Roles
`sender`, `receiver`, `merchant_analyst`, `merchant_admin`, `platform_operator`, `demo_operator`. (See conflict C-01 for the annex's shorter list.)

### 12.7 Core signature artefacts
**Payment Truth Passport (PTP)** · **Truth Gap Graph** · **Two-Lens Explanation** · **Next Safe Action Engine** · **Counterfactual Recovery Mirror**.

### 12.8 Agents (exact source roster — do not rename, add, or remove)
| ID | Name |
|---|---|
| A01 | Payment Readiness |
| A02 | Traffic & Service Intelligence |
| A03 | Sender Analysis |
| A04 | Receiver Analysis |
| A05 | Journey Tracker |
| A06 | Root Cause Analyst |
| A07 | Currency Intelligence |
| A08 | Recovery Planner |
| A09 | Dual-Side Communicator |
| A10 | Resolution Tracker |

### 12.9 Safe action codes (DESIGN DECISION; catalogue owned by doc 04)
`DO_NOT_REPEAT_CHECKOUT`, `AWAIT_PROVIDER_CONFIRMATION`, `REFRESH_ORIGINAL_ORDER`, `RECONCILE_MERCHANT_ORDER`, `CHECK_OFFICIAL_ACCOUNT_ACTIVITY`, `COMPLETE_OFFICIAL_VERIFICATION`, `CONTACT_PROVIDER_SUPPORT`, `REQUEST_HUMAN_REVIEW`, `NONE_REQUIRED`.

### 12.10 Scenario and test identifiers
Chaos `CHAOS-01…12` (SOURCE M13) · AI evals `T01…T10` (SOURCE M14) · Demos `A…E` · Acceptance `AO-xx` (§17) · Requirements `REQ-xx` (doc 09 traceability).

---

## 13. Differentiators and honest prior art

### 13.1 The five signature inventions (SOURCE M04 §4.4 — *proposed, not patented*)

| | Feature | One-line definition |
|---|---|---|
| A | **Payment Truth Passport** | Versioned, role-scoped bundle: source, timestamp, last verified stage, unknowns, next safe action. A PayBridge artefact — *not* a PayPal document |
| B | **Truth Gap Graph** | What the browser thinks, what the merchant DB thinks, what PayPal proves, what the recipient verified — contradictions highlighted, never collapsed into one status |
| C | **Two-Lens Explanation** | One investigation, two privacy-filtered narratives from *separate data projections*, not "ask the LLM to hide" |
| D | **Next Safe Action Engine** | Deterministic action catalogue that blocks dangerous options by evidence class and state; AI can request but not override |
| E | **Counterfactual Recovery Mirror** | Labelled simulation replaying a lost response or duplicate webhook to prove the app would reconcile, not recharge — never impersonating a real provider event |

### 13.2 Prior art (SOURCE M04 §4.2 — honest comparison)
PayPal's own merchant activity and APIs; Wise transfer tracker; Airwallex monitoring; and several PayPal-hackathon entries (a cross-border reconciliation agent, Cleared, Clearhold, SettleSight, REVERSA). **PayBridge does not claim novelty for tracking, reconciliation, chatbots, or cross-border payments.** The proposed gap is the *combination*: a consented cross-party case with explicit uncertainty, evidence provenance, and two private explanations, centred on unknown-outcome reconciliation.

**Approved phrasing until proved (SOURCE M04 §4.5):**

| Do not say | Say instead |
|---|---|
| "The first platform" | "A new approach we designed" |
| "Prevents duplicate charges" | "Designed to prevent blind retries" (until the invariant test is run) |
| "Tracks money to the Indian bank" | "Tracks to the last stage we can verify" |
| "Predicts payment failure" | "Flags observed degradation and uncertain outcomes" |
| "Cuts support time 80%" | "Aims to reduce support investigation time" |
| "Live PayPal data" | "PayPal sandbox data" (after a test) or "simulation" |

---

## 14. Capability inventory: MVP vs later, with truth labels

Priority: **MUST** (golden path) · **SHOULD** (strengthens judging) · **LATER**.
Mode: SANDBOX = real PayPal sandbox; SIM = local simulation; BOTH.
Status at time of writing: **every row is PLANNED**.

| # | Capability | Priority | Mode | Acceptance signal |
|---|---|---|---|---|
| C01 | PayPal OAuth + Orders v2 create/capture | MUST | SANDBOX | Real order id + capture id recorded |
| C02 | Case + attempt + event spine | MUST | BOTH | One case ↔ one original provider operation |
| C03 | Unknown-state reconciliation | MUST | SANDBOX + fault | No new order/capture after lost response |
| C04 | Webhook inbox: verify, dedupe, replay-safe | MUST | SANDBOX | Forged rejected; duplicate = one effect |
| C05 | Payment Truth Passport (versioned) | MUST | BOTH | Every confirmed claim has evidence ids |
| C06 | Sender and receiver role views + consent | MUST | BOTH | No cross-party fields in either response |
| C07 | Next Safe Action Engine | MUST | BOTH | Retry unavailable until preconditions met |
| C08 | Bounded agentic diagnosis (A01–A10) | MUST | BOTH | Cited, validated output; fallback works |
| C09 | Truth Gap Graph workbench | MUST | BOTH | Contradictions visible with provenance drill-down |
| C10 | Counterfactual Recovery Mirror lab | SHOULD | SIM | Visibly labelled; never mixed with real ids |
| C11 | Ops console (queue, webhook health, audit) | SHOULD | BOTH | Tenant-bound queue; evidence drill-down |
| C12 | Public provider status advisory (A02) | SHOULD | Public source | Timestamped; `unknown` when unavailable |
| C13 | FX/fee explanation (A07) | SHOULD | Indicative | `unavailable` rather than invented |
| C14 | Realtime updates (SSE) | SHOULD | BOTH | Passport version bump appears without reload |
| C15 | Demo walkthrough / scenario selector | MUST | BOTH | One-click reproducible Demo B |
| C16 | Bank-settlement proof | LATER | needs authorised data | Otherwise `NOT_OBSERVABLE` |
| C17 | Multi-tenant production isolation (RLS everywhere) | SHOULD (RLS) / LATER (scale) | BOTH | Cross-tenant tests pass |
| C18 | Helpdesk/CRM/accounting integrations | LATER | needs permissions | — |
| C19 | Predictive risk models | LATER | needs dataset | Calibrated lift over rules |
| C20 | PWA/push, native mobile | LATER | — | After usage data |
| C21 | Refund/void flows | LATER, human-authorised only | — | Separate audited controller; never agent-driven |

---

## 15. User journeys

### 15.1 Journey matrix (SOURCE M07 §7.2)

| Journey | Initiator | Trigger | User-visible outcome | Safe end state |
|---|---|---|---|---|
| Readiness | Sender | Before checkout | Prerequisites, advisory | Proceed only via provider-authorised flow |
| Normal completion | Buyer | Approved checkout | "PayPal capture confirmed" | Resolved after sync |
| Unknown capture | App/worker | Timeout | "Checking original transaction" | Reconciled or escalated |
| Pending capture | Provider | `PENDING` | Last confirmed state; no fulfilment conclusion | Await terminal state |
| Decline | Provider | `DENIED` | Sender guidance without private rationale | Resolved as denied |
| Duplicate webhook | Provider | Same event twice | No duplicate effect | One effective processing |
| Out-of-order events | Provider | Later state first | Monotonic verified state | No regression without valid reversal |
| Receiver action | Receiver | Declared requirement | Role-specific checklist | Evidence of action or escalation |
| Dispute | Merchant/PayPal | Dispute status | Human-approved official workflow | No autonomous dispute decision |
| Unknown bank settlement | Receiver | Withdrawal question | `NOT_OBSERVABLE` + official guidance | Outside our evidence |

### 15.2 Per-persona journey sketches

**Sender (S).** Log in → pick invoice → readiness check → approve in PayPal → see "Confirming payment" → (hero) sees "unknown — do not pay again" with the Do-not-retry banner → Passport updates to "provider capture confirmed" → acknowledges.

**Receiver (R).** Receives a consent invitation showing exactly what will be shared → accepts → sees "Incoming payment status": *PayPal capture confirmed* separately from *Bank settlement not verified* → follows official account link; can ask for help or revoke.

**Operator (O).** Opens the queue sorted by uncertainty/verification age → opens the three-column investigator (Evidence | Causal Timeline | Diagnosis/Actions) → toggles *Compare provider vs app state* → approves internal reconciliation → closes with evidence.

**Developer (D).** Opens the Sandbox Lab → chooses a scenario (e.g., CHAOS-01) → runs → inspects the verified-vs-simulated webhook table and API traces → exports a reproducible report.

**Judge (J).** Lands on the public landing → watches the 60-second story → opens the scenario selector → runs Demo B → reads the evidence pack.

### 15.3 Exceptional cases that must be designed up front (SOURCE M07 §7.4)
Unauthorised user guessing a case URL · expired session · sandbox token expiry · stale webhook certificate · refresh during approval · currency exponent edge cases · multiple captures per order · pending-then-completed · missing original order id · inconclusive provider GET · webhook before local commit · provider completed but merchant DB failed · duplicate case import · LLM timeout · prompt injection in invoice memo · withdrawal status inaccessible · timezone/currency display ambiguity · premature case closure by support. **Each needs a documented safe fallback and a test** (doc 09).

---

## 16. Scope boundaries

| In scope | Out of scope |
|---|---|
| Commercial checkout observation | Consumer remittances |
| Evidence-first incident handling | Bypassing KYC/AML or sanctions screening |
| Consented sender/receiver case sharing | Universal bank-transfer tracking |
| Deterministic safe guidance | PayPal internal queues, risk, or global traffic |
| Developer observability + sandbox lab | Escrow, custody, wallets |
| Reproducible labelled demos | Autonomous refund/capture/retry |
| Role-scoped explanations | Compliance eligibility decisions |
| Honest unknown/NOT_OBSERVABLE states | Fabricated provider records or bank settlement |

**Scope-drift rule (DESIGN DECISION).** Any proposed feature not in §14 requires a decision-log entry. Any feature that touches money movement is rejected by default.

---

## 17. Measurable acceptance outcomes

*These are acceptance **thresholds**, not results. A value becomes a result only when a test is actually executed and its output saved.*

| ID | Outcome | Measure | Threshold | Evidence artefact |
|---|---|---|---|---|
| AO-01 | Real sandbox integration | A real order id and capture id exist and an independent `GET` confirms them | 1 complete run | `docs/evidence/paypal-sandbox-receipt.md` (redacted) |
| AO-02 | No blind repeat | Orders + captures per case after hero scenario | exactly 1 / 1 | DB ledger export + API trace |
| AO-03 | Unknown ≠ failed | UI/API tests asserting `API_TIMEOUT`/`CAPTURE_UNKNOWN` never render "failed" | 0 violations | Test report |
| AO-04 | Webhook safety | Forged webhook accepted / duplicate with >1 effect | 0 / 0 | Webhook test log |
| AO-05 | Temporal consistency | State regression from `CAPTURE_COMPLETED` on out-of-order events | 0 | State-transition tests |
| AO-06 | Citation coverage | Confirmed claims with valid evidence ids ÷ all confirmed claims | 100 % | Validator report |
| AO-07 | Unsupported material claims | Count in the controlled eval suite | 0 | AI eval report (with suite size stated) |
| AO-08 | Role privacy | Leakage in privacy challenge tests | 0 | Permission diff report |
| AO-09 | Tenant isolation | Cross-tenant read success | 0 | Security tests |
| AO-10 | AI independence | Core case truth and deterministic guidance with LLM disabled | works | Demo E recording |
| AO-11 | Simulation labelling | Simulated events lacking a SIMULATED marker anywhere in UI/API | 0 | UI + API audit |
| AO-12 | Reproducibility | A fresh clone reaches a running demo using documented commands | works | Clean-machine log |
| AO-13 | Accessibility | Automated + manual checks on the 16 screens against WCAG 2.2 AA target | no blocking failures | A11y report |
| AO-14 | Secret hygiene | Secrets found in repo/history/logs by scanner | 0 | Scan output |

**Honesty rule.** A 100 % result on three hand-picked cases is not proof of general reliability; always report suite size and failure count (SOURCE M17 §17.4).

---

## 18. Risks

*Condensed from M18 §18.6 with owners; the full register is in doc 10.*

| Risk | Sev. | Mitigation | Contingency |
|---|---|---|---|
| Cross-border sandbox route unsupported | High | Validate account config first (Phase 1) | Real supported sandbox checkout + labelled US→IN simulation |
| Hosted webhook endpoint unreliable | High | Stable HTTPS URL; durable inbox; polling reconcile | Local harness clearly labelled |
| LLM invents a cause | High | Schema + citation validator + abstention + evals | Deterministic known/unknown explanation |
| Receiver consent complexity | High | Narrow invitations, field allowlist | Ship ops + sender first |
| Credentials committed | Critical | `.env` discipline, secret scanning, rotation drill | Revoke and audit |
| Judge cannot run app | High | Hosted demo + local runner + recorded evidence | Evidence pack |
| Prior-art overlap | Med | Honest narrative; prove differentiators | Emphasise UX + dual-side trust |
| Merchant WTP unproven | High | Interviews, pilots | Reposition as open-source reliability toolkit |
| Scope creep (the plan is huge) | High | MUST/SHOULD/LATER gating (doc 10) | Cut SHOULD items, never MUST |

---

## 19. Decision log and source conflicts

### 19.1 Decisions

| ID | Decision | Rationale | Revisit if |
|---|---|---|---|
| D-01 | Go modular monolith: one API process + one worker process | Clear boundaries, low ops cost (SOURCE ADR-005) | Measured throughput needs |
| D-02 | PostgreSQL authoritative for PayBridge records; Redis optional | ACID + uniqueness constraints (SOURCE) | — |
| D-03 | SQL outbox, not a broker | No lost steps between commit and job (SOURCE) | Measured need → NATS/SQS |
| D-04 | Add `REOPENED` to the M13 case workflow | Annex had it; resolution must reopen on contrary evidence (SOURCE I-resolution) | — |
| D-05 | `/api/v1` REST + OpenAPI contract; SSE for live updates | SOURCE M15 | — |
| D-06 | Frontend route prefix `/app/...` (M15) | Matches API-era nav; M08 sitemap superseded | — |
| D-07 | Repository layout per M18 §18.4 (`apps/web`, `services/api`, `packages/contracts`, `infra`, `tests`, `docs`) | Newest, most complete | — |
| D-08 | MVP authentication: server-side session cookies with seeded demo accounts; OIDC as a pluggable later path | Fast to build; M16 requires a *mature* IdP for production (doc 08) | Moving beyond a demo |
| D-09 | Money stored as `NUMERIC(20,6)` with currency-exponent validation; API carries decimal strings | SOURCE annex §9 / M15.9 | — |
| D-10 | Render for the hosted demo; AWS reference only | SOURCE M10 | — |
| D-11 | Annex golden scenarios `T01…T20` (M20 annex §17) are renamed **`GS-01…GS-20`**; the ID prefix `T01…T10` is reserved for the **AI evaluation set** (M14 §14.9) | Same prefix meant two different suites (conflict C-09) | — |
| D-12 | The MVP Orchestrator is **deterministic Go** (fixed DAG + rule-table router). An LLM "planner/advisor" prompt is specified in doc 07 but is **OFF by default** and can only choose a subset of the Go-allowed agents; it can never add agents, tools or actions | Reconciles M20 §14 "Orchestrator LLM" with doc 03's Go policy router (conflict C-12) | A measured need for LLM routing |
| D-13 | Environment variable names are defined authoritatively in doc 08 `.env.example`; `PAYPAL_API_BASE` is kept (docs 02/04) and M20's `PAYPAL_BASE_URL` is not used | Conflict C-10 | — |
| D-14 | `paypal.Reader` is GET-only (`ShowOrder`, `ShowCapture`); webhook verification is a separate `paypal.WebhookVerifier` interface | Doc 04 §9.1 contradicted §2.2/doc 03 (conflict C-11) | — |
| D-15 | Mirror Lab / Demo *run* controls require `demo_operator`; the seeded judge/demo account holds `sender` + `demo_operator` in the demo tenant; other roles may view runs | Doc 05 §2 vs doc 04 R37 (conflict C-13) | Real multi-tenant onboarding |
| D-16 | Visual direction for doc 06 is **"Evidence Ledger"** (warm paper surfaces, ink text, serif display + sans UI + mono ids, seals for evidence tiers, a two-lens motif) | Premium, distinctive, calm; avoids generic SaaS look | Design validation findings |

### 19.2 Conflicts found in the source (documented, not silently resolved)

| ID | Conflict | Where | Resolution adopted |
|---|---|---|---|
| C-01 | Roles: M16 lists `sender, receiver, merchant_analyst, merchant_admin, platform_operator, demo_operator`; annex `tenant_memberships` CHECK allows only `sender, receiver, analyst, admin` | M16 §16.2 vs annex §9 | Use the M16 six. Schema migrated accordingly (doc 04) |
| C-02 | Case workflow: M13 `OPEN→OBSERVING→…→RESOLVED` vs annex `OPEN\|NEEDS_ACTION\|UNDER_INVESTIGATION\|HUMAN_REVIEW\|GUIDANCE_SENT\|RESOLVED\|REOPENED` | M13 §13.3 vs annex §8 | Use M13 + `REOPENED` (D-04) |
| C-03 | Application-state naming: passport example `TIMED_OUT`, API example `CONFIRMATION_LOST`, annex enum `API_TIMEOUT` | M11 §11.7, M15 §15.8, annex §8 | Canonical enum = `API_TIMEOUT`. `CONFIRMATION_LOST` is a claim code / display label only |
| C-04 | Sitemap prefixes: M08 `/sender/...` vs M15 `/app/sender/...` | M08 §8.2 vs M15 §15.2 | M15 (D-06) |
| C-05 | Repo layout: archived spec `services/backend` vs M18 `services/api` | Part III vs M18 | M18 (D-07) |
| C-06 | Environments: annex checks four values; M11 uses `REAL_SANDBOX` only in examples | annex §9 vs M11 | Four-value enum (§12.1) |
| C-07 | Annex says Python is not required; M10 mentions "Python only isolated model work" as an alternative | M02 vs annex | No Python in the MVP |
| C-08 | Source says "Render" for demo and also lists a "sponsor ecosystem" reason; no sponsor requirement is documented | M02 §2.3 | Render chosen on technical merit; **VERIFY** any hackathon sponsor-stack requirements against official rules |
| C-09 | Test-ID clash: M14 `T01–T10` (AI evals) vs annex §17 `T01–T20` (golden scenarios) | M14 §14.9 vs annex §17 | `GS-xx` for golden scenarios (D-11) |
| C-10 | Env var naming: `PAYPAL_API_BASE` (docs 02/04) vs `PAYPAL_BASE_URL` (M20 annex sample) | doc 02 §6 vs annex §20 | `PAYPAL_API_BASE` (D-13) |
| C-11 | Doc 04 §9.1 put a POST method inside `Reader` while §2.2/doc 03 require Reader to be GET-only | doc 04 | Split `WebhookVerifier` (D-14); doc 04 corrected with change note |
| C-12 | M20 §14 says an "Orchestrator LLM" chooses agents; doc 03 §5/S7 say a Go policy router does | M20 §14.3 vs doc 03 | Deterministic Go orchestrator; optional advisory prompt OFF (D-12) |
| C-13 | Doc 05 §2 lets sender/analyst use the Mirror Lab; doc 04 R37 is `demo_operator`-only | doc 05 vs 04 | D-15; doc 05 note added |
| C-14 | Capability IDs `C01–C21` (§14) look like conflict IDs `C-01…` | doc 01 | Clarified here: capability IDs have no hyphen; no rename needed |

---

## 20. Phase roadmap (summary)

*Dependency-led, no artificial time estimates (SOURCE M18). Full epics/tasks in doc 10.*

| Phase | Name | Deliverable | Gate |
|---|---|---|---|
| P0 | Audit and truth definitions | Repo audit, state matrix, env/key inventory | Owner approves blueprint |
| P1 | Foundations + PayPal sandbox proof | Repo skeleton, Compose, migrations, OAuth, create/capture, show-order/show-capture | A real sandbox capture recorded and independently re-read |
| P2 | Reliable event pipeline | Webhook inbox, verification, dedupe, outbox, workers, reconciliation | Timeout + replay tests pass |
| P3 | Cases, roles, consent, Passport | Auth, RBAC/ABAC, invites, projections, Passport versioning, safe-action engine | Privacy + tenant tests pass |
| P4 | Frontend core | Landing, auth, dashboard, payment detail, Passport, timeline, sender/receiver lenses | Golden path usable end-to-end on real API |
| P5 | Agentic investigation | Orchestrator, A01–A10, tools, validator, fallback, eval suite | Eval thresholds + Demo E |
| P6 | Signature visualisations | Truth Gap Graph, Recovery Mirror, evidence explorer, ops console | Demo B end-to-end |
| P7 | Hardening and deployment | Security tests, observability, Render deploy, backup/restore, runbooks | Release gate (doc 09) |
| P8 | Evidence pack and submission | Receipts, eval report, video, README, licence | Honest-claims checklist signed off |

---

## 21. The full end-to-end product story

*Narrative form — usable as the basis for the README, the pitch, and the demo voice-over.*

**Act 1 — The invoice.** Maya runs accounts payable at a small US studio. A freelance designer in India, Arjun, has delivered a brand system; the invoice is USD 2,000. Maya opens PayBridge and creates a payment case. The platform asks the one question most tools never ask: *do we actually know this route is supported?* It answers honestly — "unknown until account verification" — and lists the documented prerequisites without pretending to check Arjun's private KYC.

**Act 2 — The checkout.** Maya proceeds. PayBridge's Go service persists a case, an attempt, and a stable operation key *before* calling PayPal. It creates the order, stores the opaque id, and sends Maya into PayPal's own approval flow. She approves.

**Act 3 — The disappearance.** The merchant backend calls capture. PayPal captures the money. And then — the response never reaches the application. (In the demo this is a deliberate fault in PayBridge's own response path; PayPal did nothing wrong, and PayBridge says so.) Maya's screen does not say "Payment failed". It says: **"We didn't hear back from our server. This does not mean PayPal declined your payment. Please don't pay again."**

**Act 4 — The investigation.** A reconciliation worker reads the *original* order and capture from PayPal, never creating another. A signed webhook arrives — and then arrives again; the inbox absorbs the duplicate. The orchestrator wakes only the agents it needs: the Journey Tracker assembles the timeline, Sender Analysis interprets Maya's side, the Root Cause Analyst weighs evidence, the Recovery Planner proposes only non-money-moving steps, and the Dual-Side Communicator writes two different messages. A validator rejects any sentence without evidence.

**Act 5 — Two lenses, one truth.** Maya's Passport reads: *Provider capture confirmed. Our confirmation was lost. Do not pay again. Recipient bank credit: not observable.* Arjun, who consented to the case invitation, sees only what he is permitted: *PayPal capture recorded on the seller side — check your official account activity for availability and withdrawal.* Neither sees the other's private data. The Truth Gap Graph shows the contradiction that caused the panic: the browser said "unknown", the merchant database said "no payment", PayPal said "completed".

**Act 6 — The close.** A merchant analyst approves reconciling the *internal* order record. The ledger shows one order, one capture. The case resolves, the Passport version is recorded, and the Counterfactual Recovery Mirror lets any viewer replay the failure in a clearly labelled simulation and watch the system reconcile instead of recharge.

**Epilogue — What PayBridge will not claim.** It will not say the money reached Arjun's bank. It will not say it is the first of its kind. It will not call itself production-ready. It will say precisely what it can prove, and show where each proof came from.

---

## 22. Claims policy

Every public claim (README, pitch, video, UI copy) maps to exactly one category (SOURCE M20 §20.5):

| Category | Meaning | Use |
|---|---|---|
| VERIFIED DOCUMENTATION | Current official doc, URL + date + scope | May cite |
| REAL SANDBOX OBSERVATION | Executed call with redacted receipt | May state "PayPal sandbox" |
| SIMULATED BEHAVIOUR | Local fixture | Must say "simulation" |
| IMPLEMENTED PRODUCT | Working code + repeatable test | Only after verified |
| VALIDATED USER INSIGHT | Interview/experiment with method and consent | Only if it exists |
| PROPOSED DESIGN | Not built | Must say "planned" |
| UNKNOWN | Production eligibility, bank settlement, PayPal internals, ROI, worldwide originality, adoption | **Never claimed** |

Forbidden phrases unless proven: *first ever, fully compliant, bank-level secure, production ready, guaranteed winner, prevents all duplicate payments, tracks to the bank, live PayPal data (for simulated flows).*

---

## 23. Hackathon fit

*From M20 §20.4. **VERIFY** all rules against the official PayPal AI Hackathon page and submission form before the deadline; this document makes no claim about deadlines or prizes.*

| Judging dimension (equal weight per the plan) | PayBridge strategy | Anti-pattern to avoid |
|---|---|---|
| Technology | Real Orders v2 capture, signed webhooks, idempotency, reconciliation, typed agents | Mock PayPal success with no real API evidence |
| Design | Complete sender/receiver/ops experience; premium, calm uncertainty states (doc 06) | Generic admin table |
| Impact | One specific, painful, demonstrable incident with a measured safe outcome | Fabricated TAM or savings |
| Innovation | Passport + Truth Gap Graph + Two-Lens + Safe Action + Mirror | Claiming trackers never existed |
| Presentation | 3-minute conflict → proof → explanation → outcome | Reading agent names without showing value |

**Submission checklist (VERIFY):** working demo (not a mockup alone) · public source repository with open-source licence and runnable instructions · short public video (the plan states under three minutes on YouTube) · meaningful use of PayPal developer platform and AI.

---

## 24. Repository and environment audit finding

*Result of the mandatory "inspect first" step, performed 9 Nov 2026 in the session sandbox.*

| Item | Finding |
|---|---|
| M20 plan | **Found and read** (4,722 lines). |
| Application repository | **Not found.** No `.git`, `go.mod`, `package.json`, README, source, tests, or build config were present in the session. |
| `.env` / secret config | **None present.** No key names or key values were inspected, because no file exists. |
| Existing functioning code | None found, so nothing was overwritten or preserved. |

**Consequences (DESIGN DECISION):**

1. The blueprint assumes a **greenfield** repository laid out per D-07. If a repository exists on the owner's machine, the implementation agent **must run the audit in doc 02 §4 first** and reconcile before writing code.
2. Required environment variable *names* are defined in doc 08 §`.env.example`. The owner must supply values; no credential has been created, guessed, or stored.
3. Nothing in these ten documents is evidence of a working integration, a passing test, or a measured result. **All tests are NOT RUN. All integrations are UNVERIFIED.**

---

**Addendum (Account 2 session, 9 Nov 2026).** The workspace was re-inspected in a second Claude session: still **no application repository, no `.git`, no `go.mod`/`package.json`, no `.env*` files**, and none of the project environment-variable names set. The greenfield assumption stands. Documents 06–10 were added in this session; the whole pack remains **DRAFT** and every test/integration remains **NOT RUN / UNVERIFIED**.

*End of document 01. Next: `02_CLAUDE_MASTER_IMPLEMENTATION_PROMPT.md`.*
