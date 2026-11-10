# PayBridge AI — Vijay Repository Inventory and Ownership Audit

**Audit Target:** Vijay Kasthuri Final Submitted Repository  
**Source ZIP:** `C:\Users\iamra\Downloads\PayBridge-AI-Vijay-Final-Verified-Repository.zip`  
**Total Source Entries:** 87 files  
**Canonical Spec Reference:** `16_TEAM_WORK_ALLOCATION_AND_TIMELINE_MASTER.md`

---

## 1. Ownership Boundaries (Doc 16)

According to canonical Document 16:
- **Vijay Kasthuri owns:**
  1. Autonomous Multi-Agent DAG A01–A10 (`internal/agents`, `internal/orchestrator`).
  2. The 3D State Reducer (`internal/reducer`).
  3. Server-Side Claim Validator (`internal/claims`).
  4. Deterministic Next Safe Action Engine (`internal/safeaction`).
  5. The Remaining 10 Frontend Screens (`features/screens/Screens.tsx`).
  6. Frontend Signature Interfaces (`LifecycleTimeline.tsx`, `TruthGapGraph.tsx`, `PaymentTruthPassport.tsx`).
  7. Automated Golden Demo B runner (`cmd/demob`, `scripts/demo-b.ps1`, `scripts/demo-b.sh`) and Chaos regression test suite (`internal/chaos`).
- **Sriram owns (Preserved Baseline):**
  1. Monorepo skeleton, Docker Compose, CI/CD pipelines, root configuration.
  2. Database Migrations 0001–0010 (Canonical 40 tables with full RLS).
  3. Core HTTP API server (`cmd/api`), Chi router, tenancy context, error envelopes, and distributed capture lock.
  4. Base React 18 shell, API client, and MVP 6 core screens.
- **Diny Paul owns (External Scope):**
  1. Kafka event streaming (>10,000 TPS).
  2. AWS KMS envelope encryption.
  3. Stripe / Wise multi-rail bank settlement and RBI e-FIRC statutory ingestion.

---

## 2. File-by-File Classification Matrix

| File Path | Classification | Ownership | Action Taken | Rationale |
| :--- | :--- | :--- | :--- | :--- |
| `services/api/internal/agents/agents.go` | **NEW** | Vijay | **INTEGRATED** | Core A01–A10 specialist agent deterministic implementation. Imports aligned. |
| `services/api/internal/agents/agents_test.go` | **NEW** | Vijay | **INTEGRATED** | Unit tests for A01–A10 deterministic outputs and untrusted wrapping. |
| `services/api/internal/claims/validator.go` | **NEW** | Vijay | **INTEGRATED** | Server-side claim validator, prompt injection filter, and DLP account redaction. |
| `services/api/internal/claims/validator_test.go` | **NEW** | Vijay | **INTEGRATED** | Tests for tier downgrades, citation gates, and prompt injection rejection. |
| `services/api/internal/safeaction/safeaction.go` | **NEW** | Vijay | **INTEGRATED** | Next Safe Action Engine enforcing `DO_NOT_REPEAT_CHECKOUT` invariants. |
| `services/api/internal/safeaction/safeaction_test.go` | **NEW** | Vijay | **INTEGRATED** | Tests verifying blind retry blocking on unknown capture state. |
| `services/api/internal/reducer/reducer.go` | **NEW** | Vijay | **INTEGRATED** | 3D State Reducer (`provider_state`, `application_state`, `receiver_state`). |
| `services/api/internal/reducer/reducer_test.go` | **NEW** | Vijay | **INTEGRATED** | Tests for timeout != failure and monotonicity guards. |
| `services/api/internal/reducer/types.go` | **NEW** | Vijay | **INTEGRATED** | Canonical state vector, evidence tier E1–E6, and event definitions. |
| `services/api/internal/projector/projector.go` | **NEW** | Vijay | **INTEGRATED** | Dual-lens projection engine (Sender vs Receiver lens with DLP). |
| `services/api/internal/projector/projector_test.go` | **NEW** | Vijay | **INTEGRATED** | Tests verifying receiver consent gating and cross-party redaction. |
| `services/api/internal/orchestrator/orchestrator.go` | **NEW** | Vijay | **INTEGRATED** | Multi-agent DAG orchestrator coordinating 6 sequential/parallel steps. |
| `services/api/internal/orchestrator/orchestrator_test.go` | **NEW** | Vijay | **INTEGRATED** | Tests verifying DAG timeout and deterministic fallback execution. |
| `services/api/internal/chaos/chaos_test.go` | **NEW** | Vijay | **INTEGRATED** | CHAOS-01 to CHAOS-12 adversarial resilience test suite. |
| `services/api/cmd/worker/main.go` | **NEW** | Vijay | **INTEGRATED** | Background worker hosting reconciler and orchestrator (mutator-free). |
| `services/api/cmd/demob/main.go` | **NEW** | Vijay | **INTEGRATED** | Golden Demo B scenario runner (504 timeout recovery). |
| `scripts/demo-b.ps1` | **NEW** | Vijay | **INTEGRATED** | Safe PowerShell runner for Golden Demo B. |
| `scripts/demo-b.sh` | **NEW** | Vijay | **INTEGRATED** | Safe Bash runner for Golden Demo B. |
| `apps/web/src/features/screens/Screens.tsx` | **NEW** | Vijay | **INTEGRATED** | All 16 screens, providing the remaining 10 enterprise screens. |
| `apps/web/src/features/signature/TruthGapGraph.tsx` | **NEW** | Vijay | **INTEGRATED** | Interactive SVG Truth Gap visualization. |
| `apps/web/src/features/signature/LifecycleTimeline.tsx` | **NEW** | Vijay | **INTEGRATED** | Chronological lifecycle timeline with dual application/provider states. |
| `apps/web/src/features/signature/PaymentTruthPassport.tsx` | **NEW** | Vijay | **INTEGRATED** | Tamper-sealed Payment Truth Passport viewer. |
| `apps/web/src/types/index.ts` | **NEW** | Vijay | **INTEGRATED** | Frontend TypeScript types for case state vectors and safe actions. |
| `apps/web/src/styles/base.css` | **NEW** | Vijay | **INTEGRATED** | Tokenized CSS styles. |
| `apps/web/src/styles/tokens.css` | **NEW** | Vijay | **INTEGRATED** | Design tokens for theme and color spaces. |
| `services/api/internal/paypal/classify.go` | **NEW** | Vijay | **INTEGRATED** | PayPal response classification engine into canonical states. |
| `services/api/internal/paypal/reader.go` | **NEW** | Vijay | **INTEGRATED** | GET-only PayPal Reader client interface. |
| `services/api/internal/paypal/mutator.go` | **NEW** | Vijay | **INTEGRATED** | Isolated PayPal Mutator client interface. |
| `services/api/internal/paypal/oauth.go` | **NEW** | Vijay | **INTEGRATED** | PayPal TokenManager with cached token invalidation. |
| `services/api/internal/paypal/types.go` | **NEW** | Vijay | **INTEGRATED** | PayPal payload and response definitions. |
| `services/api/worker.exe` | **CONFLICT** | Vijay | **QUARANTINED** | Precompiled 5.6MB binary committed into archive. Quarantined for security. |
| `services/api/db/migrations/0001–0010` | **CONFLICT** | Sriram | **PRESERVED SRIRAM** | Vijay's migrations defined only 10 partial tables. Sriram's canonical 40-table migrations preserved. |
| `services/api/cmd/api/main.go` | **DUPLICATE** | Sriram | **PRESERVED SRIRAM** | Sriram's HTTP API server is canonical, includes Chi router, auth, and tenancy. |
| `services/api/internal/config/config.go` | **DUPLICATE** | Sriram | **PRESERVED SRIRAM** | Sriram's fail-closed config engine with credential masking and 28 tests preserved. |
| `apps/web/package.json` | **CONFLICT** | Sriram | **PRESERVED SRIRAM** | Vijay's package.json added extraneous `@radix-ui` dependencies not required by code. |
| `docker-compose.yml` | **DUPLICATE** | Sriram | **PRESERVED SRIRAM** | Sriram's multi-service Docker Compose preserved. |
