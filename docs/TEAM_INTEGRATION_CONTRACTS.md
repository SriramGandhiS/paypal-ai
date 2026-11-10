# PayBridge AI — Team Handshake Contracts & Interface Boundaries

## 1. Operational Pillar Ownership
Per `16_TEAM_WORK_ALLOCATION_AND_TIMELINE_MASTER.md`:

```
┌────────────────────────────────────────────────────────────────────────────────────────┐
│                                  PAYBRIDGE AI PLATFORM                                 │
└───────────────────────────────────────────┬────────────────────────────────────────────┘
                                            │
        ┌───────────────────────────────────┼───────────────────────────────────┐
        ▼                                   ▼                                   ▼
 [MEMBER 1: SRIRAM]               [MEMBER 2: VIJAY KASTHURI]             [MEMBER 3: DINY PAUL]
 Tool: Google Antigravity         Tool: Anthropic Claude                 Tool: Claude & Astra
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

## 2. Shared Handshake Contracts

### Contract 1: Session Authentication & HMAC Assertion Tokens
* **Owner:** Sriram  
* **Consumers:** Vijay (React UI, Agent Calls), Diny (Enterprise API)  
* **Schema:**
  ```json
  {
    "key_id": "key-2026-prod-01",
    "tenant_id": "018f3a21-9b34-7a12-89cd-456789000000",
    "user_id": "018f3a21-9b34-7a12-89cd-456789111111",
    "role": "merchant_admin",
    "op": "CREATE_CASE",
    "resource_id": "00000000-0000-0000-0000-000000000000",
    "nonce": "n-78a9c2",
    "expiry": "2026-11-10T16:00:00Z"
  }
  ```
* **Status:** **LOCKED (S3)**

### Contract 2: 3D State Reducer Input/Output
* **Owner:** Sriram (Storage & Locks) / Vijay (Reducer Rules)  
* **Canonical Input (RFC 8785):**
  ```json
  {"case_id":"018f3a21-9b34-7a12-89cd-456789111111","evidence_id":"018f3a21-9b34-7a12-89cd-456789222222","expected_version":1}
  ```
* **Verified SHA-256 Digest:**
  `4fc20f0e89261a829b68f023fcc582934dae262520ded004aa12ec660191ac6b`
* **Status:** **LOCKED (S3)**

### Contract 3: AI Investigation & Diagnostic Claims
* **Owner:** Vijay  
* **Consumer:** Sriram (Case Views, Notifications)  
* **Contract Interface:** Rows in `investigations` and `diagnosis_claims` linked to `claim_evidence`.
* **Status:** **PENDING INTEGRATION** (Awaiting Vijay's orchestrator module)

### Contract 4: Safe Action Engine Proposals & Human Acks
* **Owner:** Vijay  
* **Consumer:** Sriram (Checkout Executor, Ops Console)  
* **Contract Interface:** Rows in `recovery_actions` with `action_payload`; human acknowledges in `action_acks`.
* **Status:** **PENDING INTEGRATION** (Awaiting Vijay's Safe Action Engine)

### Contract 5: Merchant Payment Passport
* **Owner:** Sriram (Storage & Hashing) / Diny (Multi-Rail Aggregation)  
* **Contract Interface:** JSON snapshot stored in `passport_versions`, capability tokens in `passport_shares`.
* **Status:** **LOCKED (S3)**

### Contract 6: Data Retention & Purge Worker
* **Owner:** Sriram (Purge Procedures) / Diny (Retention Policies)  
* **Contract Interface:** Procedure `fn_retention_execute_purge(p_retention_job_id)` deleting from `evidence_payloads`.
* **Status:** **LOCKED (S3)**

### Contract 7: Cloud KMS Envelope Encryption
* **Owner:** Diny  
* **Consumer:** Sriram (Connector Ingestion)  
* **Contract Interface:** Envelope format in `provider_connections` (`kek_key_id`, `encrypted_dek`, `encrypted_credentials`, `nonce`).
* **Status:** **LOCKED (S3)**

### Contract 8: Core REST API & PayPal Checkout Golden Path
* **Owner:** Sriram  
* **Consumers:** Vijay (React UI & Action Engine), Diny (Integration Tests)  
* **Endpoints:**
  - `POST /api/v1/cases` & `GET /api/v1/cases`: Case creation with idempotency and paginated listing.
  - `GET /api/v1/cases/{id}`: Case retrieval.
  - `POST /api/v1/cases/{id}/attempts` & `GET /api/v1/cases/{id}/attempts`: Payment attempt initiation (with blind-retry guard) and history.
  - `POST /api/v1/cases/{id}/evidence` & `GET /api/v1/cases/{id}/evidence`: Provenance evidence ingestion (with I11 simulated tier 6 rule) and retrieval.
  - `POST /api/v1/cases/{id}/reconcile`: Reader-only financial state reconciliation preserving provider ground truth.
  - `POST /api/v1/cases/{id}/checkout/paypal`: Golden path PayPal sandbox order creation with approval URL.
  - `POST /api/v1/cases/{id}/checkout/paypal/capture`: Golden path PayPal sandbox capture and automatic reconciliation.
  - `POST /api/v1/webhooks/paypal`: Cryptographically verified webhook ingestion gateway into `webhook_inbox`.
* **Status:** **IMPLEMENTED & UNIT TESTED (100% PASS)**
