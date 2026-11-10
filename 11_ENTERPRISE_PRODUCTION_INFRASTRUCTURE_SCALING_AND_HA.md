# 11 — Enterprise Production Infrastructure, Scaling, and High Availability

| Field | Value |
|---|---|
| Document | `11_ENTERPRISE_PRODUCTION_INFRASTRUCTURE_SCALING_AND_HA.md` |
| Version | 2.0.0-ENTERPRISE |
| Date | 10 November 2026 |
| Status | **SPECIFICATION — Enterprise Production Extension** |
| Audience | Enterprise Architects, Cloud/DevOps Engineers, SREs, Security Officers |
| Extends | `03_SYSTEM_DESIGN_AND_ARCHITECTURE.md`, `04_BACKEND_ENGINEERING_SPECIFICATION.md`, `09_TESTING_DEVOPS_DEPLOYMENT_AND_RELEASE_PLAYBOOK.md` |
| Purpose | Upgrades PayBridge AI from single-region demo/sandbox hosting to an institutional-grade, multi-region, high-throughput distributed system processing >10,000 TPS with zero data loss. |

---

## Table of Contents
1. [Target Architecture & SLA Tiers](#1-target-architecture--sla-tiers)
2. [Multi-Region Cloud Topology](#2-multi-region-cloud-topology)
3. [Distributed Streaming Event Bus (Kafka / NATS)](#3-distributed-streaming-event-bus-kafka--nats)
4. [Enterprise Database Sharding, Partitioning & Failover](#4-enterprise-database-sharding-partitioning--failover)
5. [Hardware Security Modules (HSM) & KMS Envelope Encryption](#5-hardware-security-modules-hsm--kms-envelope-encryption)
6. [Distributed Caching, Lock Coordination & Edge Protection](#6-distributed-caching-lock-coordination--edge-protection)
7. [Observability, Telemetry & Distributed Tracing at Scale](#7-observability-telemetry--distributed-tracing-at-scale)
8. [Disaster Recovery, Chaos Engineering & Failover Runbooks](#8-disaster-recovery-chaos-engineering--failover-runbooks)

---

## 1. Target Architecture & SLA Tiers

For high-end enterprise banking, commerce, and fintech deployment, the transactional and informational boundaries operate under strict contractual Service Level Agreements (SLAs).

| Metric | MVP / Sandbox Target (Docs 01–10) | Enterprise Tier-1 Production (This Spec) |
|---|---|---|
| **System Availability** | 99.0% (~7.2 hours downtime/month) | **99.99% (< 4.38 minutes downtime/month)** |
| **Peak Throughput** | 50 operations/second | **10,000+ operations/second** |
| **Recovery Point Objective (RPO)** | 1 hour | **< 1.0 second (Zero financial data loss)** |
| **Recovery Time Objective (RTO)** | 4 hours | **< 60 seconds (Automated regional failover)** |
| **Read P99 Latency** | < 400 ms | **< 45 ms globally (Edge-cached projections)** |
| **Reconciliation Processing** | Worker DB polling (`SKIP LOCKED`) | **Sub-second distributed stream processing** |
| **Secrets Management** | Local `.env` + static AES key | **Cloud HSM / KMS envelope encryption with 90-day automated rotation** |

---

## 2. Multi-Region Cloud Topology

```
                                  [GLOBAL EDGE LAYER]
                     Cloudflare Enterprise / AWS CloudFront
             - DDoS Protection (L3/L4/L7) + Anycast IP Routing
             - Web Application Firewall (OWASP CRS, Bot Management)
             - Edge TLS 1.3 Termination & Mutual TLS (mTLS) for Webhooks
                                      │
            ┌─────────────────────────┴─────────────────────────┐
            ▼                                                   ▼
   [PRIMARY REGION (US-EAST)]                          [SECONDARY REGION (EU-CENTRAL)]
   Active Traffic Zone                                 Warm Standby / In-Region Data Zone
┌───────────────────────────────────────────┐       ┌───────────────────────────────────────────┐
│ Kubernetes (EKS/GKE) Cluster              │       │ Kubernetes (EKS/GKE) Cluster              │
│  - API Pods (Horizontal Auto Scaler)      │       │  - API Pods (Standby Scaled Down)         │
│  - Worker Pods (Keda Stream Consumers)    │       │  - Worker Pods (Ready to Scale)           │
│  - Envoy Service Mesh (Strict mTLS)       │       │  - Envoy Service Mesh (Strict mTLS)       │
│                                           │       │                                           │
│ Data Infrastructure:                      │       │ Data Infrastructure:                      │
│  - Managed Apache Kafka Cluster (3 AZs)   │◀─────▶│  - MirrorMaker 2 Cross-Region Streaming   │
│  - PostgreSQL Aurora Multi-AZ (Write Lead)│──────▶│  - Aurora Cross-Region Replica (Read Only)│
│  - Redis Enterprise Cluster (3 AZs)       │       │  - Redis Replica Cluster                  │
│  - AWS KMS / Vault HSM                    │       │  - In-Region KMS Keys                     │
└───────────────────────────────────────────┘       └───────────────────────────────────────────┘
```

### 2.1 Traffic Routing & Ingress Rules
1. **Anycast Routing with Health Probes**: Global traffic enters via Cloudflare Enterprise with health checks monitoring `/healthz` and `/readyz` every 5 seconds.
2. **Automated DNS Failover**: If the Primary Region health drops below 98% healthy backends, Cloudflare automatically reroutes 100% of ingress to the Secondary Region.
3. **mTLS for Inbound Webhooks**: Payment provider webhooks (PayPal, Stripe, Banking webhooks) terminate at dedicated edge API gateways enforcing client certificate validation where supported, or dedicated webhook validation proxies.

---

## 3. Distributed Streaming Event Bus (Kafka / NATS)

While PostgreSQL `FOR UPDATE SKIP LOCKED` (Doc 04 §13) is ideal for initial workloads up to 200 jobs/sec, enterprise scale requires decoupling event ingress, outbox emission, and worker execution onto a high-performance distributed commit log.

### 3.1 Topic Topology & Partitioning Strategy

| Topic Name | Partitions | Partition Key | Retention | Compaction |
|---|---|---|---|---|
| `paybridge.payment.events.v1` | 32 | `tenant_id:case_id` | 7 years (Audit tier) | Disabled |
| `paybridge.webhooks.inbound.v1` | 32 | `provider:resource_id` | 30 days | Disabled |
| `paybridge.reconciliation.jobs.v1` | 16 | `case_id` | 7 days | Disabled |
| `paybridge.passport.updates.v1` | 32 | `case_id` | 90 days | Enabled |
| `paybridge.agent.triggers.v1` | 16 | `investigation_id` | 7 days | Disabled |
| `paybridge.deadletter.v1` | 8 | `source_topic` | 365 days | Disabled |

### 3.2 Ordering and Exactly-Once Guarantees
1. **Partition Key Invariant**: Every event for a specific `case_id` is assigned the same partition key (`tenant_id:case_id`). This guarantees **strict in-order delivery** to worker consumers for any single payment case.
2. **Transactional Outbox via Debezium CDC**:
   - Rather than backend workers polling the SQL database, the Go API writes domain events atomically to a SQL `outbox_events` table inside the business transaction.
   - Debezium (Change Data Capture) captures PostgreSQL WAL (Write-Ahead Log) records and streams them to Kafka with sub-10ms latency.
3. **Consumer Idempotency**: Consumers maintain a local bloom filter backed by Redis and SQL idempotency tracking (`idempotency_keys` table) guaranteeing that redelivered Kafka messages produce zero side-effects.

---

## 4. Enterprise Database Sharding, Partitioning & Failover

### 4.1 Table Partitioning (PostgreSQL)
High-volume financial ledgers and event tables must use declarative table partitioning by time range (`RANGE (occurred_at)`) or composite hash (`HASH (tenant_id)`):

```sql
-- Partitioned Payment Events Table
CREATE TABLE payment_events (
  id UUID NOT NULL DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL,
  case_id UUID NOT NULL,
  attempt_id UUID,
  event_type TEXT NOT NULL,
  environment TEXT NOT NULL,
  payload_redacted JSONB NOT NULL,
  encrypted_payload_ref TEXT,
  occurred_at TIMESTAMPTZ NOT NULL,
  received_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (occurred_at, tenant_id, id)
) PARTITION BY RANGE (occurred_at);

-- Automated Monthly Partition Creation (pg_partman)
-- Partitions older than 24 months are detached and archived to Parquet format in S3 cold storage.
```

### 4.2 Connection Pool Sizing & Multiplexing
- **PgBouncer / Supavisor Architecture**: Deployed as sidecars or a dedicated tier. Transaction-level pooling (`pool_mode = transaction`) enables 10,000 concurrent client connections to map onto 120 actual PostgreSQL server backends without thread contention.
- **Failover SLA**: AWS Aurora PostgreSQL configured with Multi-AZ cluster and global database replication. Unplanned leader crash initiates automated failover to replica in **< 18 seconds**.

---

## 5. Hardware Security Modules (HSM) & KMS Envelope Encryption

In an institutional environment, static symmetric keys stored in `.env` (`ENCRYPTION_KEY`) violate PCI-DSS (Req 3.5/3.6) and SOC 2 (CC6.1). Enterprise PayBridge uses **KMS Envelope Encryption**.

```
                           [AWS KMS / HashiCorp Vault]
                                        │
                       Generates 256-bit Data Encryption Key (DEK)
                                        │
                   ┌────────────────────┴────────────────────┐
                   ▼                                         ▼
         [Plaintext DEK]                           [Encrypted DEK (Ciphertext)]
        (Kept in RAM only;                         (Stored safely in PostgreSQL
      wiped after encryption)                       alongside the ciphertext)
                   │                                         │
                   ▼                                         ▼
   AES-256-GCM Payload Encryption             Database Column: `encrypted_payload`
```

### 5.1 Key Lifecycle & Rotation
1. **Master Key (CMK)**: Stored in FIPS 140-3 Level 3 validated Hardware Security Modules (HSMs).
2. **Envelope Scheme**: For every sensitive evidence payload (raw bank wires, raw provider webhook dumps, unredacted financial tokens), PayBridge requests a unique Data Encryption Key (DEK). The raw DEK encrypts the payload via AES-256-GCM with a unique 96-bit initialization vector (IV); the plaintext DEK is securely zeroed in memory, and the encrypted DEK is stored in the database.
3. **Automated Key Rotation**: Root CMKs rotate automatically every 90 days. Decryption reads the key version from the payload header and uses historical versions transparently.

---

## 6. Distributed Caching, Lock Coordination & Edge Protection

### 6.1 Redis Cluster Configuration
- **Topology**: 6-node Redis 7.2+ Cluster (3 primaries, 3 secondaries across 3 Availability Zones).
- **Redlock Distributed Mutex**:
  ```go
  // Invariant: Only one capture or reconciliation worker may operate on an attempt across all pods.
  mutex := redsync.New(pools).NewMutex(
      fmt.Sprintf("lock:attempt:%s", attemptID),
      redsync.WithExpiry(15 * time.Second),
      redsync.WithTries(3),
  )
  if err := mutex.LockContext(ctx); err != nil {
      return ErrAttemptConcurrentModification
  }
  defer mutex.UnlockContext(ctx)
  ```

### 6.2 Edge Rate Limiting & Bot Defense
- **Distributed Token Bucket**: Enforced at the Cloudflare edge worker tier and backed by Redis Enterprise at the ingress layer.
- **Credential Stuffing & DDoS Defense**: Webhook endpoints enforce strict IP whitelisting based on published provider IP CIDR blocks (e.g., PayPal / Stripe IP ranges), with strict rate limiting (max 5,000 req/min per webhook source).

---

## 7. Observability, Telemetry & Distributed Tracing at Scale

### 7.1 Distributed Tracing Standards (OpenTelemetry + W3C Trace Context)
Every HTTP call, Kafka message, DB query, and outbound provider call propagates standard W3C `traceparent` headers.

```
Client Checkout -> API (trace_id: 4bf92f3577b34da6a3ce929d0e0e4736)
  ├── PostgreSQL Write (attempt_id, status=ORDER_CREATING) [1.2ms]
  ├── Outbox Publish -> Kafka (topic: paybridge.payment.events.v1) [2.1ms]
  └── PayPal Outbound POST /v2/checkout/orders [142.3ms]
       └── Async Webhook Handler (parent_trace: 4bf92f3577b34da6a3ce929d0e0e4736) [4.5ms]
```

### 7.2 Key Production Metrics & Alert Thresholds

| Metric Name | Type | Target | Alert Threshold (PagerDuty P1) |
|---|---|---|---|
| `paybridge_webhook_inbox_latency_seconds` | Histogram | P99 < 500ms | > 2.0s for 2 consecutive minutes |
| `paybridge_reconciliation_lag_seconds` | Gauge | < 10s | > 60s |
| `paybridge_duplicate_payment_prevention_total`| Counter | N/A (informational)| Spike > 50/hour (indicates upstream chaos)|
| `paybridge_kms_decryption_errors_total` | Counter | 0 | > 0 (Immediate investigation) |
| `paybridge_kafka_consumer_lag` | Gauge | < 100 msgs | > 5,000 msgs for 5 minutes |

---

## 8. Disaster Recovery, Chaos Engineering & Failover Runbooks

### 8.1 Automated Regional Failover Protocol (RTO < 60s)
1. **Detection**: Route53 / Cloudflare synthetic monitoring detects Primary Region ingress degradation (< 95% availability for 30s).
2. **Standby Promotion**:
   - Aurora PostgreSQL global database in Secondary Region initiates managed failover via AWS Lambda or Terraform operator (`rds:failover-global-cluster`).
   - Secondary region consumer pods scale up via KEDA based on Kafka partition assignments.
3. **Traffic Shift**: Edge DNS reroutes all live traffic to the Secondary Region.
4. **Post-Failover Audit**: Hash-chained audit logs run cryptographic verification between regions to verify zero data loss.

### 8.2 Production Chaos Engineering (Game Days)
- **Chaos Mesh / Gremlin Scenarios**:
  - `CHAOS-PROD-01`: Kill primary PostgreSQL master during active capture load.
  - `CHAOS-PROD-02`: Partition Kafka broker 1 of 3 during webhook burst.
  - `CHAOS-PROD-03`: Induce 4,000ms artificial latency on PayPal OAuth endpoints.
- **Pass Criteria**: Zero duplicate charges, zero un-reconciled transactions left in `ORDER_CREATING`, zero silent data corruption.
