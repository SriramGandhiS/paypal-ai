# 19 — Diny Paul Execution Blueprint & Claude/Astra Master Prompt Deck

| Field | Value |
|---|---|
| Document | `19_DINY_PAUL_CLAUDE_ASTRA_EXECUTION_AND_PROMPT_DECK.md` |
| Assignee | **Diny Paul** |
| Tool / Model | **Anthropic Claude & Astra (Advanced Enterprise Agents)** |
| Role | **Enterprise Infrastructure, Multi-Rail, Banking & Compliance Architect** |
| Primary Specs | `11_ENTERPRISE_PRODUCTION_INFRASTRUCTURE_SCALING_AND_HA.md`, `12_MULTI_RAIL_PAYMENTS_BANK_SETTLEMENT_AND_LIFECYCLE.md`, `13_REGULATORY_COMPLIANCE_KYC_AML_AND_DATA_GOVERNANCE.md`, `14_ENTERPRISE_AI_SLM_GOVERNANCE_AND_SUB_SECOND_FALLBACK.md`, `15_MASTER_PRODUCTION_RELEASE_AND_MIGRATION_PLAYBOOK.md` |
| Secondary Specs | `04_BACKEND_ENGINEERING_SPECIFICATION.md`, `08_PAYPAL_INTEGRATION_SECURITY_PRIVACY_AND_COMPLIANCE.md`, `16_TEAM_WORK_ALLOCATION_AND_TIMELINE_MASTER.md` |

---

## Table of Contents
1. [Diny Paul's Operational Role & Scope](#1-diny-pauls-operational-role--scope)
2. [Task Breakdown & Definition of Done (Tasks D1–D7)](#2-task-breakdown--definition-of-done-tasks-d1d7)
3. [Claude & Astra Master System Prompt (Copy-Paste Ready)](#3-claude--astra-master-system-prompt-copy-paste-ready)
4. [Step-by-Step Task Execution Prompts for Claude & Astra](#4-step-by-step-task-execution-prompts-for-claude--astra)
5. [Verification Commands & Evidence Collection Runbook](#5-verification-commands--evidence-collection-runbook)

---

## 1. Diny Paul's Operational Role & Scope

As the Enterprise Infrastructure, Multi-Rail, and Compliance Lead, Diny Paul upgrades PayBridge AI from a single-provider sandbox prototype into an institutional banking-grade platform:
- Distributed Event Streaming: Replacing database worker polling with Apache Kafka (32 partitions) and Debezium Change Data Capture (CDC) processing >10,000 TPS.
- Cryptographic Vault & KMS: Replacing static `.env` keys with FIPS 140-3 HSM / KMS Envelope Encryption (AES-256-GCM) and automated 90-day rotation.
- Multi-Rail Adapters: PayPal Live, Stripe PaymentIntents, and Wise Multi-Currency Payouts.
- Bank Settlement Observability: Integrating Open Banking (Plaid) and RBI-regulated Account Aggregators (Sahamati / Setu) to achieve real-time proof of deposit (`receiver_state = VERIFIED_BANK_CREDITED`).
- Statutory Cross-Border Compliance: RBI Purpose Codes (P0802), automated e-FIRC PDF ingestion, Persona KYC/KYB, and ComplyAdvantage OFAC/PEP screening.
- Institutional IAM & Sub-Second AI: Okta SAML 2.0 / OIDC, SCIM 2.0 automated provisioning, and private vLLM Small Language Model clusters (<450 ms P99 SLA).

---

## 2. Task Breakdown & Definition of Done (Tasks D1–D7)

| Task ID | Task Title | Primary Deliverable | Acceptance Criteria |
|---|---|---|---|
| **D1** | Distributed Streaming Bus | `services/api/internal/kafka/`, Debezium CDC config | Kafka handles 10,000 synthetic msgs/sec with in-order case partitioning; Debezium streams PostgreSQL WAL outbox events with <15ms latency. |
| **D2** | KMS Envelope Encryption | `services/api/internal/kms/` | Generates 256-bit DEKs via AWS KMS / Vault; encrypts raw payloads via AES-256-GCM; zero plaintext DEKs stored; 90-day key rotation tested. |
| **D3** | Multi-Rail Adapter Suite | `services/api/internal/multirail/` | Standardized Go adapters for Stripe and Wise; full refund and dispute state machines operational. |
| **D4** | Open Banking & Settlement | `services/api/internal/settlement/` | Plaid and Sahamati AA adapters retrieve signed bank transactions; matches deposit amount/date/ref and transitions `receiver_state` to `VERIFIED_BANK_CREDITED`. |
| **D5** | Double-Entry SQL Ledger | `services/api/db/migrations/0011_ledger.sql` | Immutable append-only double-entry ledger; database trigger enforces that every transaction sums to exactly ZERO. |
| **D6** | Cross-Border Compliance | `services/api/internal/compliance/` | RBI Purpose Code P0802 validation; e-FIRC SHA-256 PDF ingestion; Persona KYC and ComplyAdvantage OFAC/PEP screening integration. |
| **D7** | Enterprise IAM & Sub-Sec AI | `services/api/internal/iam/`, `internal/slm/` | Okta SAML 2.0 / OIDC and SCIM 2.0 provisioning; private vLLM 8B model cluster + Redis semantic vector cache delivers <450ms P99 triage. |

---

## 3. Claude & Astra Master System Prompt (Copy-Paste Ready)

*Instruction for Diny: Copy everything between the block below and paste it into Claude or Astra at the beginning of your implementation session:*

```markdown
=== BEGIN CLAUDE & ASTRA SYSTEM PROMPT ===

YOU ARE: Claude & Astra, implementing the Institutional Enterprise Infrastructure, Multi-Rail, Banking Rails, and Compliance Engine for PayBridge AI.
YOU ARE WORKING FOR: Diny Paul (Enterprise Infrastructure & Compliance Lead).
YOUR TEAMMATES ARE: Sriram (Foundations Lead) and Vijay Kasthuri (AI, Reducer & UX Lead).

YOUR CORE RESPONSIBILITIES:
1. Build the high-throughput distributed event streaming infrastructure (Apache Kafka + Debezium CDC) scaling beyond 10,000 TPS.
2. Implement AWS KMS / HashiCorp Vault Hardware Security Module (HSM) Envelope Encryption for all sensitive evidence.
3. Deliver the multi-rail payment adapters (PayPal Live, Stripe PaymentIntents, Wise Payouts) and full refund/dispute lifecycles.
4. Solve receiver bank settlement observability by integrating Open Banking (Plaid) and RBI Account Aggregators (Sahamati / Setu).
5. Enforce statutory cross-border compliance: RBI Purpose Codes (P0802), e-FIRC ingestion, FinCEN Travel Rule, and Persona KYC/AML.
6. Implement Enterprise Identity (Okta SAML 2.0 / OIDC, SCIM 2.0) and Sub-Second Private SLMs (vLLM on L40S) with Redis Semantic Caching.

ENTERPRISE INVARIANTS (STRICT ENFORCEMENT):
1. Zero Data Loss: RPO < 1.0 second; all payment events partitioned strictly by tenant_id:case_id in Kafka.
2. Balanced Ledger: Every financial entry in the SQL ledger MUST balance to zero (debits == credits).
3. Cryptographic Proof: Bank deposits are only marked VERIFIED_BANK_CREDITED when backed by signed Open Banking / AA transaction receipts.
4. Directory Ownership: You own services/api/internal/kafka/, internal/kms/, internal/multirail/, internal/settlement/, internal/compliance/, and internal/iam/. Do NOT overwrite Sriram's db/ or httpx/ code, or Vijay's reducer/ or agents/ code.

ONLINE REFERENCES & DOCUMENTATION TO CONSULT:
- Debezium PostgreSQL CDC: https://debezium.io/documentation/reference/stable/connectors/postgresql.html
- AWS KMS Envelope Encryption: https://docs.aws.amazon.com/kms/latest/developerguide/concepts.html#enveloping
- Stripe API Reference: https://stripe.com/docs/api
- Wise Platform API: https://docs.wise.com/api-reference
- Plaid Transactions API: https://plaid.com/docs/api/products/transactions/
- Sahamati Account Aggregator Architecture: https://sahamati.org.in/technical-framework/
- Okta SCIM 2.0 Specification: https://developer.okta.com/docs/concepts/scim/
- vLLM High-Throughput Inference: https://docs.vllm.ai/

PRIMARY BLUEPRINT DOCUMENTS IN WORKSPACE:
Read in full:
- 11_ENTERPRISE_PRODUCTION_INFRASTRUCTURE_SCALING_AND_HA.md (Kafka, Multi-Region, KMS)
- 12_MULTI_RAIL_PAYMENTS_BANK_SETTLEMENT_AND_LIFECYCLE.md (Multi-Rail, Disputes, Open Banking)
- 13_REGULATORY_COMPLIANCE_KYC_AML_AND_DATA_GOVERNANCE.md (RBI, e-FIRC, Persona, Okta SCIM)
- 14_ENTERPRISE_AI_SLM_GOVERNANCE_AND_SUB_SECOND_FALLBACK.md (vLLM SLM, Semantic Cache, SR 11-7)
- 15_MASTER_PRODUCTION_RELEASE_AND_MIGRATION_PLAYBOOK.md (Release Gates PRG-01 to PRG-10)
- 16_TEAM_WORK_ALLOCATION_AND_TIMELINE_MASTER.md (Team Handshake Contracts)

When Diny issues a Task, inspect existing interfaces, implement with institutional precision, run integration tests, and provide a comprehensive milestone report.

=== END CLAUDE & ASTRA SYSTEM PROMPT ===
```

---

## 4. Step-by-Step Task Execution Prompts for Claude & Astra

### Task D1 Prompt: Distributed Event Streaming (Kafka + Debezium CDC)
```text
Task D1: Implement the distributed event streaming pipeline using Apache Kafka and Debezium Change Data Capture.
Reference: 11_ENTERPRISE_PRODUCTION_INFRASTRUCTURE_SCALING_AND_HA.md §3.

Actions:
1. In `services/api/internal/kafka/`:
   - Producer client with snappy compression, idempotent writes (acks=all), and retry backoff.
   - Consumer client with consumer groups, automatic offset management, and dead-letter queue routing (`paybridge.deadletter.v1`).
   - Partition Key Rule: Every event MUST use partition key `tenant_id:case_id` to guarantee strict in-order processing per case.
2. Configure Debezium PostgreSQL connector (`infra/kafka/debezium-postgres.json`):
   - Capture changes on table `outbox_events` from PostgreSQL WAL.
   - Stream directly to topic `paybridge.payment.events.v1` with <15ms latency.
3. Write high-throughput benchmark test simulating 10,000 events/sec. Verify zero out-of-order deliveries.
```

### Task D2 Prompt: AWS KMS / Vault Envelope Encryption
```text
Task D2: Implement the FIPS 140-3 HSM / KMS Envelope Encryption engine.
Reference: 11_ENTERPRISE_PRODUCTION_INFRASTRUCTURE_SCALING_AND_HA.md §5.

Actions:
1. In `services/api/internal/kms/`:
   - Implement `VaultClient` supporting AWS KMS and HashiCorp Vault.
   - `GenerateDataKey(ctx, keyID)` -> (plaintextDEK, encryptedDEK, error).
   - `EncryptPayload(ctx, rawBytes)` -> encrypts payload via AES-256-GCM using unique 96-bit nonce; wipes plaintext DEK from memory; returns encrypted ciphertext and encrypted DEK reference.
   - `DecryptPayload(ctx, ciphertext, encryptedDEK)` -> decrypts DEK via KMS, then decrypts payload.
   - Support automated 90-day key rotation by reading key version headers from ciphertext.
2. Write unit tests asserting that memory is securely zeroed and encrypted payloads cannot be read without KMS access.
```

### Task D3 Prompt: Multi-Rail Payment Adapter Suite (Stripe & Wise)
```text
Task D3: Implement multi-rail provider adapters for Stripe and Wise, with full refund and dispute lifecycles.
Reference: 12_MULTI_RAIL_PAYMENTS_BANK_SETTLEMENT_AND_LIFECYCLE.md §1–§4 & §6.

Actions:
1. In `services/api/internal/multirail/`:
   - `stripe/`: Stripe PaymentIntents client, webhook signature verification, chargeback handler (`charge.dispute.created`), and refund handler.
   - `wise/`: Wise Multi-Currency Payouts client, wire transfer tracking (`outgoing_payment_sent`), and UETR/IMAD banking clearing numbers ingestion.
   - Both adapters MUST implement the unified Go `Reader` and `Mutator` interfaces from Doc 12 §2.
2. Implement dispute evidence dossier generator (Agent A08 extension) assembling verified timestamps, checkout auth tokens, and bank receipts into a standard dispute packet.
3. Run integration tests with simulated Stripe and Wise webhooks.
```

### Task D4 Prompt: Open Banking & Account Aggregator Settlement Observability
```text
Task D4: Build the Open Banking and Account Aggregator settlement observation engine to verify direct recipient bank deposits.
Reference: 12_MULTI_RAIL_PAYMENTS_BANK_SETTLEMENT_AND_LIFECYCLE.md §5.

Actions:
1. In `services/api/internal/settlement/`:
   - `plaid/`: Plaid Link token exchange and `/transactions/sync` client for US/UK bank account monitoring.
   - `sahamati/`: India Account Aggregator (AA) client connecting via Setu / Sahamati FIP framework.
   - `matcher.go`: Deterministic deposit match verification algorithm:
     * deposit.amount == canonical_invoice.net_amount
     * deposit.occurred_at >= capture.occurred_at
     * deposit.narration matches payment reference
   - Upon verified match, emit `BANK_DEPOSIT_VERIFIED` event, transitioning `receiver_state` from `REPORTED_PENDING` to `VERIFIED_BANK_CREDITED` backed by Tier E1 bank evidence.
2. Write mock tests verifying the state transition on valid bank transaction match, and rejection on amount mismatch.
```

### Task D5 Prompt: Immutable Double-Entry SQL Ledger Engine
```text
Task D5: Implement the immutable append-only double-entry financial ledger schema and balancing triggers.
Reference: 12_MULTI_RAIL_PAYMENTS_BANK_SETTLEMENT_AND_LIFECYCLE.md §7.

Actions:
1. Create `services/api/db/migrations/0011_ledger.sql`:
   - Tables: `ledger_accounts`, `ledger_entries`, `ledger_postings`.
   - Constraints: Every posting must have non-zero amount (Positive for Debit, Negative for Credit).
   - Invariant Trigger: `verify_ledger_entry_balanced()` PostgreSQL trigger asserting that the sum of amounts across all postings for a given `entry_id` is EXACTLY ZERO.
2. In `services/api/internal/ledger/`:
   - Transactional posting helper ensuring balanced entries for checkout capture, provider fees, FX conversion, and payouts.
3. Write unit tests trying to insert an unbalanced ledger entry; verify that PostgreSQL rejects it with a transaction abort.
```

### Task D6 Prompt: Statutory Cross-Border Compliance & Sanctions Engine
```text
Task D6: Implement the cross-border statutory compliance engine: RBI Purpose Codes, e-FIRC ingestion, FinCEN Travel Rule, and Persona KYC/AML.
Reference: 13_REGULATORY_COMPLIANCE_KYC_AML_AND_DATA_GOVERNANCE.md.

Actions:
1. In `services/api/internal/compliance/`:
   - `purpose_codes.go`: Mandatory RBI Purpose Code validator (P0802 for software export, P0803 for data processing, etc.). Rejects India-bound cases lacking valid purpose code.
   - `efirc.go`: Ingests AD Category-I Bank electronic FIRC PDFs; extracts Unique Transaction Reference (UTR) and computes cryptographic SHA-256 hash.
   - `travel_rule.go`: Checks if transaction exceeds $3,000 USD; enforces collection of verified originator and beneficiary identifiers.
   - `sanctions.go`: Real-time screening integration with ComplyAdvantage / Dow Jones API checking OFAC SDN, PEP, and international sanctions lists. Freezes case immediately on match.
2. Test: Verify simulated OFAC hit blocks payment capture within 100 milliseconds.
```

### Task D7 Prompt: Enterprise IAM (Okta SAML/SCIM) & Sub-Second Private SLMs
```text
Task D7: Implement enterprise IAM via Okta SAML 2.0 / SCIM 2.0, and deploy private sub-second SLMs with Redis semantic caching.
Reference: 13_REGULATORY_COMPLIANCE_KYC_AML_AND_DATA_GOVERNANCE.md §4 and 14_ENTERPRISE_AI_SLM_GOVERNANCE_AND_SUB_SECOND_FALLBACK.md.

Actions:
1. In `services/api/internal/iam/`:
   - SAML 2.0 / OIDC middleware supporting Okta, Azure AD, and WorkOS with mandatory WebAuthn FIDO2 MFA.
   - RFC 7644 SCIM 2.0 endpoints (/scim/v2/Users) for automated user provisioning and instant session revocation on employee offboarding.
2. In `services/api/internal/slm/`:
   - Private vLLM client connecting to quantized 8B model cluster for Tier 1 triage (A01, A05).
   - Redis semantic vector cache using cosine similarity >= 0.96 for recurring triage patterns (12ms latency).
   - NeMo Guardrails and Canary Token injection to detect prompt injection and prevent data exfiltration.
   - SR 11-7 Model Risk logging: records immutable model weights hash, prompt version, and citation scores.
3. Verify P99 triage latency is <450ms across 500 concurrent synthetic requests.
```

---

## 5. Verification Commands & Evidence Collection Runbook

After completing tasks D1 through D7, Diny executes and captures:

```bash
# Verify Kafka Streaming & Consumer Lag
kafka-consumer-groups --bootstrap-server localhost:9092 --describe --group paybridge-workers

# Verify KMS Envelope Encryption
cd services/api && go test -v ./internal/kms/...

# Verify Multi-Rail & Open Banking Settlement
go test -v ./internal/multirail/... ./internal/settlement/...

# Verify Double-Entry Ledger Zero-Sum Constraint
go test -v ./internal/ledger/...

# Verify Compliance, KYC & Sanctions Engine
go test -v ./internal/compliance/...

# Verify Sub-Second AI Latency & Semantic Cache Hit Rate
go test -v ./internal/slm/... -bench=.
```
