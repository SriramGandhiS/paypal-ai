# 08 — Enterprise Production Infrastructure, Scaling & High Availability Runbook
**Phase:** 8 — Enterprise Infrastructure & HA  
**Owner:** Dinny Paul Navis  
**Branch:** `feat/diny-enterprise-completion`  
**Date:** 2026-11-11  

---

## 1. Multi-Region Active-Passive Architecture

| Region | Designation | Role | Database Tier | Kafka Cluster |
|---|---|---|---|---|
| **us-east-1** (Primary) | Primary Operational Hub | Live Writes, API ingress, Worker processing | Aurora PostgreSQL 16 (Primary Read/Write) | 3-Broker KRaft (32 Partitions) |
| **us-west-2** (Secondary) | Standby Disaster Recovery | Read-Only Replicas, Standby Workers | Aurora Global Database (< 1.0s RPO) | Standby Mirror Cluster |

### High Availability Metrics & Targets
- **Recovery Point Objective (RPO)**: $< 1.0\text{ second}$ (synchronous / streaming logical replication).
- **Recovery Time Objective (RTO)**: $< 60\text{ seconds}$ (automated Route 53 health check DNS failover).
- **Availability Target**: $99.99\%$ annual availability ($< 52.6\text{ minutes}$ downtime/year).

---

## 2. Local Container Topology (`infra/docker-compose.enterprise.yml`)

1. **PostgreSQL 16**: Configured with `wal_level=logical`, `max_replication_slots=10`, and `max_wal_senders=10`.
2. **Kafka 3-Broker Cluster**: Apache Kafka 7.6 (KRaft mode, zero ZooKeeper dependency) running across ports 9092, 9093, and 9094.
3. **Kafka Init Container**: Automatically provisions 32 partitions on `paybridge.payment.events.v1` (`min.insync.replicas=2`) and 8 partitions on `paybridge.deadletter.v1`.
4. **Debezium Connect**: Ingests PostgreSQL outbox WAL records into Kafka with zero application latency.
5. **Redis 7**: Persistent vector cache on port 6379.

---

## 3. Telemetry & Observability Pipeline

- **Prometheus Alert Rules** (`infra/observability/prometheus-alerts.yaml`):
  - `KafkaConsumerLagHigh`: Alerts if lag exceeds 500 unconsumed events.
  - `KMSEncryptionLatencyHigh`: Alerts if envelope encryption P99 exceeds 50ms.
  - `UnbalancedLedgerAttemptDetected`: Pages on any ledger zero-sum violation.
  - `SanctionsHitDetected`: Tracks compliance holds.
  - `SLMTriageP99LatencyBreach`: Triggers if sub-second triage exceeds 450ms SLA.
- **OpenTelemetry Collector** (`infra/observability/otel-collector.yaml`):
  - Strips PII, raw secrets, and authorization headers before metric/trace export.

