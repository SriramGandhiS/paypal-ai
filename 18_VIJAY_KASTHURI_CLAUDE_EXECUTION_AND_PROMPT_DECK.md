# 18 — Vijay Kasthuri Execution Blueprint & Claude Master Prompt Deck

| Field | Value |
|---|---|
| Document | `18_VIJAY_KASTHURI_CLAUDE_EXECUTION_AND_PROMPT_DECK.md` |
| Assignee | **Vijay Kasthuri** |
| Tool / Model | **Anthropic Claude (Claude 3.5 Sonnet / Claude Opus)** |
| Role | **AI Multi-Agent Lead, Reducer State Machine Architect & Signature UX Engineer** |
| Primary Specs | `01_PRODUCT_BLUEPRINT_AND_MASTER_VISION.md`, `03_SYSTEM_DESIGN_AND_ARCHITECTURE.md`, `06_UI_UX_DESIGN_SYSTEM_AND_SCREEN_BLUEPRINT.md`, `07_AI_MULTI_AGENT_ARCHITECTURE_AND_PROMPT_SPEC.md` |
| Secondary Specs | `04_BACKEND_ENGINEERING_SPECIFICATION.md` (§9, §15, §16, §17), `10_IMPLEMENTATION_ROADMAP_ACCEPTANCE_AND_DEMO.md`, `16_TEAM_WORK_ALLOCATION_AND_TIMELINE_MASTER.md` |

---

## Table of Contents
1. [Vijay Kasthuri's Operational Role & Scope](#1-vijay-kasthuris-operational-role--scope)
2. [Task Breakdown & Definition of Done (Tasks V1–V7)](#2-task-breakdown--definition-of-done-tasks-v1v7)
3. [Claude Master System Prompt (Copy-Paste Ready)](#3-claude-master-system-prompt-copy-paste-ready)
4. [Step-by-Step Task Execution Prompts for Claude](#4-step-by-step-task-execution-prompts-for-claude)
5. [Verification Commands & Evidence Collection Runbook](#5-verification-commands--evidence-collection-runbook)

---

## 1. Vijay Kasthuri's Operational Role & Scope

As the AI, State Machine, and UX Lead, Vijay Kasthuri is responsible for the brain and signature experience of PayBridge AI:
- The 3-Dimensional State Reducer: Financial truth modeling ensuring deterministic transitions and the rule that **TIMEOUT ≠ FAILURE**.
- PayPal Sandbox Adapter: Implementing the separated `Reader` (GET-only) and `Mutator` (Checkout-only) interfaces.
- The 10 Specialist Agents (A01–A10) and Fixed DAG Orchestrator: Multi-agent analysis strictly bounded by server-side JSON schemas.
- Server-Side Claim Validator: Enforcing that 100% of confirmed AI claims cite verified provider evidence.
- Signature UX: The Payment Truth Passport, SVG Truth Gap Graph, and the 16 calm, information-dense screens.

---

## 2. Task Breakdown & Definition of Done (Tasks V1–V7)

| Task ID | Task Title | Primary Deliverable | Acceptance Criteria |
|---|---|---|---|
| **V1** | 3D State Reducer | `services/api/internal/reducer/` | Unit tests verify: Reducer is the ONLY state writer; timeout transitions to `API_TIMEOUT` and `UNKNOWN` (never `FAILED`); evidence tiers E1–E6 strictly enforced. |
| **V2** | PayPal Sandbox Adapter | `services/api/internal/paypal/` | `Reader` has ZERO POST methods; `Mutator` sends stable `PayPal-Request-Id`; mock and live sandbox tests pass. |
| **V3** | Next Safe Action Engine | `services/api/internal/safeaction/` | Deterministic action engine blocks duplicate checkouts (`DO_NOT_REPEAT_CHECKOUT`) on any in-flight or unknown capture attempt. |
| **V4** | A01–A10 Orchestrator DAG | `services/api/internal/orchestrator/`, `internal/agents/` | Orchestrator runs fixed DAG from immutable evidence snapshot; all 10 agents output typed JSON schemas; fallback works with LLM disabled (Demo E). |
| **V5** | Claim Validator & Privacy Filter | `services/api/internal/claims/` | Validator rejects any confirmed claim lacking valid evidence citation; DLP filter strips cross-party data before generation. |
| **V6** | 16 React Screens & Truth Gap Graph | `apps/web/src/features/` | Hand-built SVG Truth Gap Graph; 16 screens implemented matching doc 06wireframes; calm copy rules honoured (no "failed" on timeout). |
| **V7** | Golden Demo B & Chaos Suite | `scripts/demo-b.sh`, `tests/chaos/` | Automated script runs Golden Demo B: PayPal capture succeeds, merchant confirmation drops, PayBridge reconciles without blind retry. |

---

## 3. Claude Master System Prompt (Copy-Paste Ready)

*Instruction for Vijay: Copy everything between the block below and paste it into Claude (or set as `CLAUDE.md`) at the beginning of your implementation session:*

```markdown
=== BEGIN CLAUDE SYSTEM PROMPT ===

YOU ARE: Claude, implementing the AI Multi-Agent Core, Reducer State Machine, and Signature UX for PayBridge AI.
YOU ARE WORKING FOR: Vijay Kasthuri (AI & UX Lead).
YOUR TEAMMATES ARE: Sriram (Foundations Lead) and Diny Paul (Enterprise Infra & Compliance Lead).

NORTH STAR:
"Know the last confirmed step, understand what is still unknown, and never retry blindly."

PRIME DIRECTIVES (NEVER VIOLATE):
1. PayBridge never moves money on its own initiative, and AI NEVER moves money at all. No agent, prompt, or tool reachable from an LLM may initiate, capture, refund, or retry a payment.
2. Provider evidence is financial truth; AI is NOT. An LLM output cannot write provider_state, application_state, receiver_state, or case_status. Only deterministic Go code in the Reducer can.
3. Timeout != Failure. An HTTP timeout or connection drop is OUTCOME-UNKNOWN until reconciled. Never render "failed" from a timeout.
4. Spec-First & Strict Citations: Every confirmed claim in an agent output must cite a valid evidence_id from the frozen snapshot.
5. Directory Ownership: You own services/api/internal/reducer/, internal/paypal/, internal/orchestrator/, internal/agents/, internal/claims/, internal/safeaction/, and apps/web/src/features/. Do NOT overwrite Sriram's httpx/ or db/ migrations, or Diny's kafka/ or kms/ packages.

ONLINE REFERENCES & DOCUMENTATION TO CONSULT:
- PayPal Orders v2 REST API: https://developer.paypal.com/docs/api/orders/v2/
- PayPal Payments v2 REST API: https://developer.paypal.com/docs/api/payments/v2/
- Anthropic Claude Structured Outputs & Tool Calling: https://docs.anthropic.com/en/docs/build-with-claude/tool-use
- JSON Schema Specification: https://json-schema.org/
- Radix UI Primitives: https://www.radix-ui.com/primitives
- Lucide React Icons: https://lucide.dev/guide/packages/lucide-react

PRIMARY BLUEPRINT DOCUMENTS IN WORKSPACE:
Read in full:
- 01_PRODUCT_BLUEPRINT_AND_MASTER_VISION.md (§5 Golden Use Case, §12 Canonical Vocabulary)
- 03_SYSTEM_DESIGN_AND_ARCHITECTURE.md (§10 Sequence Diagrams S1–S10, §11 State Machines, ADRs)
- 06_UI_UX_DESIGN_SYSTEM_AND_SCREEN_BLUEPRINT.md (Tokens, Grammar, 16 Screen Blueprints, Copy Deck)
- 07_AI_MULTI_AGENT_ARCHITECTURE_AND_PROMPT_SPEC.md (A01–A10 Roster, Schemas, Prompts, Validator)
- 16_TEAM_WORK_ALLOCATION_AND_TIMELINE_MASTER.md (Team Handshake Contracts)

When Vijay issues a Task, inspect existing interfaces, implement with architectural rigor, run unit tests, and provide a comprehensive milestone report.

=== END CLAUDE SYSTEM PROMPT ===
```

---

## 4. Step-by-Step Task Execution Prompts for Claude

### Task V1 Prompt: The 3D State Reducer & State Machines
```text
Task V1: Implement the deterministic Go Reducer and 3-dimensional state vector transition engine.
Reference: 03_SYSTEM_DESIGN_AND_ARCHITECTURE.md §11–§12 and 04_BACKEND_ENGINEERING_SPECIFICATION.md §15.

Actions:
1. In `services/api/internal/reducer/`:
   - Define state dimensions:
     * provider_state: NOT_CREATED, CREATED, APPROVED, AUTHORISED, CAPTURE_PENDING, CAPTURE_COMPLETED, CAPTURE_DENIED, REFUNDED, UNKNOWN.
     * application_state: SYNCHRONISED, WAITING_FOR_WEBHOOK, API_TIMEOUT, RECONCILING, EVENT_GAP, INCONSISTENT.
     * receiver_state: NOT_APPLICABLE, NOT_OBSERVABLE, REPORTED_PENDING, VERIFIED_AVAILABLE, VERIFIED_BANK_CREDITED.
     * case_status: OPEN, OBSERVING, INVESTIGATING, GUIDANCE_READY, AWAITING_USER_ACTION, VERIFYING_RESOLUTION, RESOLVED.
   - Implement `Reducer.Apply(ctx, event) -> (NewStateVector, []StateTransition, error)`.
   - Invariant: LLMs cannot call this package. Only verified provider events or deterministic timeouts can trigger state changes.
   - Invariant: A timeout or dropped response transitions attempt to CAPTURE_UNKNOWN, application_state to API_TIMEOUT, provider_state to UNKNOWN (preserving last_confirmed_provider_state). It NEVER transitions to CAPTURE_DENIED.
   - Invariant: Monotonicity. Verified provider states cannot regress backwards without an explicit reversal or refund event.
2. Write comprehensive table-driven tests in `services/api/internal/reducer/reducer_test.go` covering all valid and invalid transitions.
3. Run `go test -v ./services/api/internal/reducer/...` and verify 100% pass rate.
```

### Task V2 Prompt: PayPal Sandbox Adapter (Reader vs Mutator)
```text
Task V2: Implement the PayPal REST adapter with strict separation between Reader and Mutator interfaces.
Reference: 04_BACKEND_ENGINEERING_SPECIFICATION.md §9 and 08_PAYPAL_INTEGRATION_SECURITY_PRIVACY_AND_COMPLIANCE.md §1–§6.

Actions:
1. In `services/api/internal/paypal/`:
   - `oauth.go`: Client-credentials OAuth with token caching, singleflight refresh, and 60s pre-expiration renewal.
   - `reader.go`: GET-only Reader interface:
     * ShowOrder(ctx, env, orderID) -> (Order, RawEvidence, error)
     * ShowCapture(ctx, env, captureID) -> (Capture, RawEvidence, error)
     * Compile-time assertion: Reader MUST NOT contain any method with HTTP POST/PUT/DELETE.
   - `mutator.go`: Checkout-only Mutator interface:
     * CreateOrder(ctx, in) -> (Order, RawEvidence, error)
     * CaptureOrder(ctx, in) -> (Capture, RawEvidence, error)
     * Requires stable PayPal-Request-Id persisted BEFORE calling PayPal.
   - `verify.go`: WebhookVerifier implementing POST /v1/notifications/verify-webhook-signature (separated from Reader per ADR-021).
   - `classify.go`: Map HTTP status and body to { success, explicit_decline, pending, unknown_missing_response }.
2. Create fake PayPal test server (`fake_test.go`) simulating network drops (504 Gateway Timeout) and capture completions.
3. Run `go test -v ./services/api/internal/paypal/...`.
```

### Task V3 Prompt: Next Safe Action Engine
```text
Task V3: Implement the deterministic Next Safe Action Engine.
Reference: 04_BACKEND_ENGINEERING_SPECIFICATION.md §17.

Actions:
1. In `services/api/internal/safeaction/`:
   - Implement `Evaluate(case, attempt, stateVector) -> (allowedActions, blockedActions)`.
   - Action Catalogue: DO_NOT_REPEAT_CHECKOUT, AWAIT_PROVIDER_CONFIRMATION, REFRESH_ORIGINAL_ORDER, RECONCILE_MERCHANT_ORDER, CHECK_OFFICIAL_ACCOUNT_ACTIVITY, REQUEST_HUMAN_REVIEW, NONE_REQUIRED.
   - Safety Invariant: If attempt status is CAPTURE_UNKNOWN or application_state is API_TIMEOUT, DO_NOT_REPEAT_CHECKOUT is mandatory and retry checkout is BLOCKED.
   - Precondition evaluation: Actions require explicit precondition verification before moving from blocked to allowed.
2. Write unit tests asserting that no combination of inputs can result in a retry recommendation while a capture is outcome-unknown.
```

### Task V4 Prompt: Multi-Agent Orchestrator DAG & Agents A01–A10
```text
Task V4: Implement the Go Orchestrator, fixed DAG scheduler, and specialist agents A01 through A10.
Reference: 07_AI_MULTI_AGENT_ARCHITECTURE_AND_PROMPT_SPEC.md (§4 Orchestrator, §5 Roster, §8 Schemas, §9 Prompts).

Actions:
1. In `services/api/internal/orchestrator/`:
   - Orchestrator takes an immutable evidence snapshot at snapshot_version.
   - Policy router selects agent execution path (e.g. checkout timeout triggers: A05 -> {A03, A02} -> A06 -> A08 -> A09 -> A10).
   - Enforce token budgets, agent timeouts (30s max), and concurrency limits.
   - Implement Kill Switch: If LLM_ENABLED=false, orchestrator uses deterministic templates (Demo E).
2. In `services/api/internal/agents/`:
   - Implement prompts and typed JSON Schemas for:
     * A01 Payment Readiness
     * A02 Traffic & Service Intelligence
     * A03 Sender Analysis
     * A04 Receiver Analysis
     * A05 Journey Tracker
     * A06 Root Cause Analyst
     * A07 Currency Intelligence
     * A08 Recovery Planner
     * A09 Dual-Side Communicator
     * A10 Resolution Tracker
   - Prompt Injection Defense: Delimit all user and untrusted invoice text in labelled blocks (<untrusted_data>); never put untrusted text in instruction position.
3. Write mock-based DAG execution test in `services/api/internal/orchestrator/orchestrator_test.go`.
```

### Task V5 Prompt: Server-Side Claim Validator & Privacy DLP
```text
Task V5: Build the Server-Side Claim Validator and Cross-Party Privacy Projection Filter.
Reference: 04_BACKEND_ENGINEERING_SPECIFICATION.md §18 and 07_AI_MULTI_AGENT_ARCHITECTURE_AND_PROMPT_SPEC.md §10 & §14.

Actions:
1. In `services/api/internal/claims/`:
   - `ValidateClaim(claim, snapshot) -> (ValidClaim, error)`:
     * Schema validation: strict JSON Schema check.
     * Citation gate: every confirmed claim MUST cite valid evidence_ids present in the frozen snapshot.
     * Evidence tier gate: confirmed claims require Tier E1 or E2 evidence; Tier E5/E6 evidence can only support inference or hypothesis.
     * Downgrade policy: If a claim lacks valid evidence, downgrade to unknown or reject entirely.
2. In `services/api/internal/projector/`:
   - Projection-before-generation: Project separate data views for Sender, Receiver, and Merchant Ops BEFORE passing to LLM.
   - DLP filter: Scans generated text for unauthorized cross-party fields (e.g., receiver's private bank account number in sender's view) and redacts immediately.
3. Test with forged citations and cross-party leakage attempts; assert validator rejects them.
```

### Task V6 Prompt: Frontend 16-Screen UX, Truth Gap Graph & Passport
```text
Task V6: Implement the 16 calm, high-density React screens and signature SVG data visualizations.
Reference: 06_UI_UX_DESIGN_SYSTEM_AND_SCREEN_BLUEPRINT.md.

Actions:
1. In `apps/web/src/features/signature/`:
   - `TruthGapGraph.tsx`: Hand-built interactive SVG visualization showing contradictions among Provider, Merchant DB, Browser, and Recipient systems. Contradictions highlighted in amber/blue; never collapsed into one misleading status.
   - `PaymentTruthPassport.tsx`: Versioned, cryptographic summary showing confirmed facts, uncertain facts, evidence tier badges (E1–E6), and next safe action.
   - `LifecycleTimeline.tsx`: Chronological causal event timeline with source trust markers.
2. In `apps/web/src/features/screens/`:
   - Implement all 16 screens matching doc 06 blueprints:
     * Screen 1: Public Landing & 60-second explanation
     * Screen 2: Case Discovery & Multi-System Matrix
     * Screen 3: New Payment & Readiness Check
     * Screen 4: Checkout In-Flight & Timeout Hero Screen (with DO-NOT-RETRY banner)
     * Screen 5: Payment Truth Passport & Case Shell
     * Screen 6: Dual-Lens Sender View
     * Screen 7: Dual-Lens Receiver View (Consent-gated)
     * Screens 8–16: Truth Gap Workbench, Ops Console, Recovery Mirror Lab, etc.
   - Strict Copy Enforcement: "Failed" NEVER appears for timeout outcomes. State uses calm language: "Outcome unknown — verifying original transaction with PayPal".
```

### Task V7 Prompt: Golden Demo B Automated Runner & Chaos Suite
```text
Task V7: Build the end-to-end Golden Demo B automated scenario runner and chaos test suite.
Reference: 10_IMPLEMENTATION_ROADMAP_ACCEPTANCE_AND_DEMO.md §9 and 09_TESTING_DEVOPS_DEPLOYMENT_AND_RELEASE_PLAYBOOK.md §10.

Actions:
1. In `tests/chaos/`:
   - `CHAOS-01`: Capture response dropped by network proxy; assert Reconciler queries original capture and updates state to SYNCHRONISED with 0 duplicate checkouts.
   - `CHAOS-02`: Duplicate webhook delivery; assert idempotent single processing.
   - `CHAOS-04`: Forged webhook signature; assert immediate rejection.
   - `CHAOS-10`: LLM service unavailable; assert deterministic Passport fallback works seamlessly.
2. In `scripts/demo-b.sh`:
   - Fully automated shell script executing the Golden Demo end-to-end and generating the verified evidence report in `docs/evidence/demo-b-receipt.md`.
```

---

## 5. Verification Commands & Evidence Collection Runbook

After completing tasks V1 through V7, Vijay executes and captures:

```bash
# Run Reducer and State Transition Matrix Tests
cd services/api && go test -v ./internal/reducer/...

# Run PayPal Reader/Mutator Isolation Archtest
go test -v ./internal/archtest/...

# Run Multi-Agent DAG & Claim Validator Tests
go test -v ./internal/orchestrator/... ./internal/claims/...

# Run Frontend Component & Rendering Tests
cd apps/web && npm test

# Run End-to-End Golden Demo B
bash scripts/demo-b.sh
```
