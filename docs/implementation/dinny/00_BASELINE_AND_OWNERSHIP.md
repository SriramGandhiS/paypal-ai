# 00 — Baseline and Ownership Declaration
**Component:** Enterprise Infrastructure, Banking, Multi-Rail, Compliance & Governance  
**Owner:** Dinny Paul Navis  
**Project:** PayBridge AI  
**Branch:** `feat/diny-enterprise-completion`  

## 1. Ownership Charter
Under the tri-pillar team architecture:
- **Sriram (Google Antigravity)** owns Platform Foundations, core schema migrations 0001–0010, the Chi HTTP routing framework, and OpenAPI 3.1 contracts.
- **Vijay Kasthuri (Claude)** owns the Pure State Reducer (`internal/reducer`), Multi-Agent DAG (`internal/agents/dag.go`), Truth Gap Graph, Safe Action Engine, and the React 18 frontend (`apps/web`).
- **Dinny Paul Navis (Enterprise Lead)** owns:
  1. Distributed Event Streaming (`internal/kafka/`, Kafka KRaft topology, Debezium PostgreSQL WAL connector).
  2. Cryptographic Envelope Encryption (`internal/kms/`, AWS KMS v2 SDK, HashiCorp Vault Transit, AES-256-GCM, 90-day rotation).
  3. Multi-Rail Payment Gateway Adapters (`internal/multirail/`, Stripe PaymentIntents & HMAC-SHA256 signature verification, Wise payouts & SWIFT gpi UETR tracking).
  4. Open Banking & Account Aggregator Settlement (`internal/settlement/`, Plaid Link/Sync, Sahamati/Setu AA, deterministic deposit matching).
  5. Immutable Double-Entry SQL Ledger (`internal/ledger/`, Migration `0011_enterprise_ledger.sql`, zero-sum invariant triggers).
  6. Statutory Cross-Border Compliance (`internal/compliance/`, RBI Purpose Code P0802, e-FIRC PDF ingestion, FinCEN Travel Rule, Persona KYC, ComplyAdvantage sanctions).
  7. Enterprise IAM & Sub-Second Private SLM (`internal/iam/`, SCIM 2.0 `/scim/v2/Users`, SAML 2.0, WebAuthn; `internal/slm/`, private vLLM, Redis semantic cache, SR 11-7 model governance).
  8. Enterprise Production HA, Chaos Engineering (CHAOS-01 to CHAOS-12), and Production Release Gates (PRG-01 to PRG-10).

## 2. Invariant Rules
- **AI Never Moves Money**: Mutating operations (`RailMutator`) are strictly restricted to explicit human operator authorization.
- **Provider Evidence is Financial Truth**: Only verified provider receipts / bank records can alter case state through the deterministic Reducer.
- **Timeout != Failure**: Network dropouts and 5xx errors are preserved as `OUTCOME-UNKNOWN`. Blind retries and duplicate transactions are strictly forbidden.
- **Zero-Sum Ledger Balance**: Every financial posting must satisfy $\sum \text{Debits} - \sum \text{Credits} = 0$.

