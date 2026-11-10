# 16 — Team Work Allocation, Dependency Matrix, and Master Timeline

| Field | Value |
|---|---|
| Document | `16_TEAM_WORK_ALLOCATION_AND_TIMELINE_MASTER.md` |
| Version | 2.0.0-ENTERPRISE |
| Date | 10 November 2026 |
| Status | **SPECIFICATION — Multi-Engineer Coordination & Execution Blueprint** |
| Audience | Sriram, Vijay Kasthuri, Diny Paul, Technical Leads, AI Autonomous Agents |
| Bridges | Docs 01–10 (Core Platform) & Docs 11–15 (Enterprise Production Suite) |
| Dedicated Prompt Decks | `17_SRIRAM_ANTIGRAVITY_EXECUTION_AND_PROMPT_DECK.md`<br>`18_VIJAY_KASTHURI_CLAUDE_EXECUTION_AND_PROMPT_DECK.md`<br>`19_DINY_PAUL_CLAUDE_ASTRA_EXECUTION_AND_PROMPT_DECK.md` |

---

## Table of Contents
1. [Team Roster, Tooling & Responsibility Matrix](#1-team-roster-tooling--responsibility-matrix)
2. [Workload Distribution & Technical Specialization](#2-workload-distribution--technical-specialization)
3. [Master Sprint Timeline & Dependency Graph](#3-master-sprint-timeline--dependency-graph)
4. [Cross-Team Handshake Contracts & Interface Boundaries](#4-cross-team-handshake-contracts--interface-boundaries)
5. [Git Workflow, Branching Strategy & Zero-Conflict Protocol](#5-git-workflow-branching-strategy--zero-conflict-protocol)
6. [Daily Sync, Milestone Exit Gates & Sign-Off Checklist](#6-daily-sync-milestone-exit-gates--sign-off-checklist)

---

## 1. Team Roster, Tooling & Responsibility Matrix

To achieve institutional-grade delivery of PayBridge AI without duplication or architectural drift, the platform is divided into three equal, highly specialized operational pillars matched directly to model capabilities and developer tooling.

```
┌────────────────────────────────────────────────────────────────────────────────────────┐
│                                  PAYBRIDGE AI PLATFORM                                 │
└───────────────────────────────────────────┬────────────────────────────────────────────┘
                                            │
        ┌───────────────────────────────────┼───────────────────────────────────┐
        ▼                                   ▼                                   ▼
 [MEMBER 1: SRIRAM]               [MEMBER 2: VIJAY KASTHURI]             [MEMBER 3: DINY PAUL]
 Tool: Google Antigravity         Tool: Anthropic Claude (Sonnet/Opus)   Tool: Claude & Astra
 Pillar: PLATFORM FOUNDATIONS,    Pillar: AI ORCHESTRATION,              Pillar: ENTERPRISE INFRA,
 INGESTION, DB & CORE API         REDUCER, UX & SIGNATURE LABS           MULTI-RAIL, BANK & COMPLIANCE
 ─────────────────────────────    ─────────────────────────────          ─────────────────────────────
 • Repo Skeleton & Docker Dev     • Reducer & State Machines             • Kafka & Debezium Streaming
 • DB Migrations (0001–0010)      • PayPal Sandbox Adapter               • Multi-Rail (Stripe/Wise)
 • OpenAPI 3.1 & Dual Codegen     • A01–A10 Agents & DAG Orchestrator    • Open Banking & Bank Deposit
 • Chi Router, Auth & Middleware  • Claim Validator & Privacy Filter     • AWS KMS Envelope Encryption
 • CI/CD Pipelines & Test Rig     • 16 React Screens & Truth Gap Graph   • RBI Purpose Codes & e-FIRC
 • React Base Shell & SSE Hub     • Next Safe Action Engine              • Okta SAML/SCIM & Sub-sec AI
```

---

## 2. Workload Distribution & Technical Specialization

### Member 1: Sriram (Google Antigravity)
* **Core Specialty**: Concrete workspace execution, filesystem manipulation, local runtime compilation, containerization, schema integrity, and automated developer tooling.
* **Assigned Documents**: `04_BACKEND_ENGINEERING_SPECIFICATION.md` (Layout, DDL, Routes), `05_FRONTEND_PRODUCT_AND_IMPLEMENTATION_SPEC.md` (Base setup), `09_TESTING_DEVOPS_DEPLOYMENT_AND_RELEASE_PLAYBOOK.md` (CI/CD, Docker Compose).
* **Key Deliverables**:
  1. Monorepo directory structure, Makefile, `.env.example`, and Docker Compose (Postgres 16, Redis 7).
  2. Complete PostgreSQL goose migrations (0001 to 0010), RLS policies, and sqlc query compile.
  3. OpenAPI 3.1 contract specification (`packages/contracts/openapi.yaml`) with Go server and TypeScript client codegen.
  4. Core Go HTTP server (`cmd/api`), Chi router, error envelope middleware, and session management.
  5. React 18 + Vite frontend shell, TanStack Query provider, Tailwind CSS token system, and Server-Sent Events (SSE) subscriber client.

### Member 2: Vijay Kasthuri (Anthropic Claude)
* **Core Specialty**: Complex state machine modeling, multi-agent AI prompt engineering, deterministic safety rules, rich UX data visualization, and end-to-end chaos testing.
* **Assigned Documents**: `01_PRODUCT_BLUEPRINT_AND_MASTER_VISION.md`, `03_SYSTEM_DESIGN_AND_ARCHITECTURE.md`, `06_UI_UX_DESIGN_SYSTEM_AND_SCREEN_BLUEPRINT.md`, `07_AI_MULTI_AGENT_ARCHITECTURE_AND_PROMPT_SPEC.md`.
* **Key Deliverables**:
  1. The 3D State Reducer (`provider_state`, `application_state`, `receiver_state`) and checkout lifecycle engine.
  2. PayPal Sandbox Adapter (`Reader` GET-only interface, `Mutator` checkout flow, `WebhookVerifier`).
  3. Multi-Agent DAG Orchestrator and complete prompts/schemas for specialist agents A01 through A10.
  4. Server-Side Claim Validator (evidence tier verification, DLP privacy screening, untrusted text delimiters).
  5. Deterministic Next Safe Action Engine (`DO_NOT_REPEAT_CHECKOUT` invariants).
  6. Frontend Signature Interfaces: Payment Truth Passport, SVG Truth Gap Graph, Dual-Lens Views, and 16 complete screen wireframes.
  7. Automated Golden Demo B runner and Chaos regression test suite (CHAOS-01 to CHAOS-12).

### Member 3: Diny Paul (Claude & Astra)
* **Core Specialty**: Institutional scaling, high-throughput streaming, cryptographic security, multi-rail financial connectivity, banking API integrations, and international compliance.
* **Assigned Documents**: `08_PAYPAL_INTEGRATION_SECURITY_PRIVACY_AND_COMPLIANCE.md`, `11_ENTERPRISE_PRODUCTION_INFRASTRUCTURE_SCALING_AND_HA.md`, `12_MULTI_RAIL_PAYMENTS_BANK_SETTLEMENT_AND_LIFECYCLE.md`, `13_REGULATORY_COMPLIANCE_KYC_AML_AND_DATA_GOVERNANCE.md`, `14_ENTERPRISE_AI_SLM_GOVERNANCE_AND_SUB_SECOND_FALLBACK.md`, `15_MASTER_PRODUCTION_RELEASE_AND_MIGRATION_PLAYBOOK.md`.
* **Key Deliverables**:
  1. Distributed Event Streaming Pipeline: Apache Kafka (32 partitions) + Debezium CDC on PostgreSQL WAL (>10,000 TPS).
  2. AWS KMS / HashiCorp Vault Envelope Encryption client (DEK generation, AES-256-GCM data encryption, 90-day rotation).
  3. Multi-Rail Provider Adapters: PayPal Live, Stripe PaymentIntents, Wise Multi-Currency Payouts.
  4. Receiver Bank Settlement Observability: Plaid (US) and Sahamati / Setu Account Aggregator (India) for cryptographic bank deposit verification.
  5. Double-Entry SQL Ledger Engine with zero-sum invariant triggers.
  6. Cross-Border Statutory Engine: RBI Purpose Codes (P0802), automated e-FIRC PDF ingestion, Persona KYC/KYB, and ComplyAdvantage OFAC/PEP real-time screening.
  7. Enterprise IAM (Okta SAML 2.0 / OIDC, SCIM 2.0) and Sub-Second Private SLM cluster (vLLM on L40S) with Redis Semantic Cache.

---

## 3. Master Sprint Timeline & Dependency Graph

The project is structured across **5 sequential milestones**. No member begins a task until its upstream handshake dependency is committed to the shared repository.

```mermaid
gantt
    title PayBridge AI Master 5-Milestone Timeline
    dateFormat  YYYY-MM-DD
    section Sriram (Foundations & Shell)
    M1: Repo Skeleton, Docker & DB Migrations       :s1, 2026-11-10, 3d
    M2: OpenAPI Spec & Go Chi API Core            :s2, after s1, 3d
    M3: React Base Shell & Design Token System     :s3, after s2, 3d
    M4: SSE Real-Time Hub & Client Wiring          :s4, after s3, 2d
    M5: Production Container Build & CI Release    :s5, after s4, 3d

    section Vijay (AI, Reducer & UX)
    M1: Reducer State Vector & PayPal Sandbox      :v1, after s1, 3d
    M2: Multi-Agent Orchestrator DAG (A01-A10)     :v2, after s2, 4d
    M3: Safe Action Engine & Claim Validator       :v3, after v2, 3d
    M4: 16 React Screens & Truth Gap Graph         :v4, after s3, 4d
    M5: Golden Demo B & Chaos Test Verification    :v5, after v4, 2d

    section Diny (Enterprise, Rails & Scale)
    M1: KMS Envelope Client & Double-Entry Ledger  :d1, after s1, 3d
    M2: Multi-Rail Adapters (Stripe & Wise)        :d2, after v1, 3d
    M3: Kafka Streaming & Debezium CDC             :d3, after s2, 4d
    M4: Open Banking Ingestion (Plaid & AA)        :d4, after d2, 3d
    M5: Okta SAML/SCIM, KYC & Private SLM Cluster  :d5, after d4, 3d
```

---

## 4. Cross-Team Handshake Contracts & Interface Boundaries

To prevent team members from blocking each other or causing breaking changes:

### Handshake 1: The API & Schema Contract (Sriram ──▶ Vijay & Diny)
* **Contract File**: `packages/contracts/openapi.yaml` and `services/api/internal/store/`
* **Protocol**: Sriram defines the data models and OpenAPI routes first. Vijay and Diny bind their Go implementations directly to the sqlc-generated store interfaces and generated OpenAPI request/response types.
* **Rule**: Neither Vijay nor Diny may modify database table definitions without submitting an explicit goose migration PR (`0011_*.sql`).

### Handshake 2: The Reducer & State Vector (Vijay ──▶ Sriram & Diny)
* **Contract Package**: `services/api/internal/reducer/`
* **Protocol**: Vijay owns the Reducer package. Sriram's HTTP checkout handlers and Diny's multi-rail webhook handlers pass raw event envelopes into `reducer.Apply(ctx, event)`.
* **Rule**: Neither Sriram nor Diny may write to `provider_state`, `application_state`, or `receiver_state` directly.

### Handshake 3: The Streaming & KMS Vault (Diny ──▶ Sriram & Vijay)
* **Contract Packages**: `services/api/internal/kms/` and `services/api/internal/events/`
* **Protocol**: Diny provides the encrypted payload helper:
  `kms.EncryptPayload(ctx, rawBytes) -> (encryptedRef, error)`
  Vijay and Sriram call this helper whenever storing raw provider webhook bytes.

---

## 5. Git Workflow, Branching Strategy & Zero-Conflict Protocol

### 5.1 Branch Topology
* `main`: Protected. Commits require passing CI tests, gitleaks scan clean, and 1 peer review.
* `develop`: Integration branch. Merged daily via automated integration runs.
* Developer Working Branches:
  - Sriram: `feat/sriram-foundations`, `feat/sriram-openapi`, `feat/sriram-react-shell`
  - Vijay: `feat/vijay-reducer`, `feat/vijay-agents-dag`, `feat/vijay-truth-gap-ui`
  - Diny: `feat/diny-kms-ledger`, `feat/diny-kafka-streaming`, `feat/diny-multi-rail-open-banking`

### 5.2 Directory Ownership (Depguard Isolation)
To guarantee zero git merge conflicts, code packages are partitioned by owner:

| Directory Path | Primary Owner | Allowed Editors | Prohibited Editors |
|---|---|---|---|
| `infra/docker/`, `services/api/db/`, `services/api/internal/httpx/` | **Sriram** | Sriram | Vijay, Diny |
| `services/api/internal/agents/`, `services/api/internal/orchestrator/`, `apps/web/src/features/` | **Vijay** | Vijay | Sriram, Diny |
| `services/api/internal/kafka/`, `services/api/internal/kms/`, `services/api/internal/compliance/` | **Diny** | Diny | Sriram, Vijay |
| `packages/contracts/` | **Shared** | Requires all 3 sign-offs | Individual commits |

---

## 6. Daily Sync, Milestone Exit Gates & Sign-Off Checklist

Before any Milestone is declared complete, the three team members run the following joint verification command:

```bash
# 1. Full Go Compilation & Unit Test Matrix (Zero Race Conditions)
go test -v -race ./services/api/...

# 2. Architecture & Boundary Import Verification
go test -v ./services/api/internal/archtest/...

# 3. Security & Secret Exposure Verification
gitleaks detect --source . --verbose

# 4. Frontend Type Checking & Production Build
cd apps/web && tsc --noEmit && npm run build

# 5. Contract Drift Check
git diff --exit-code packages/contracts/openapi.yaml
```

*When all 5 checks exit with code `0`, the milestone release report is signed.*
