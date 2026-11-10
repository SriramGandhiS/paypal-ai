# PayBridge AI — Database Security & RLS Test Catalog

## 1. Scope & Objective
This document defines the 20 automated database security tests verifying multi-tenant isolation, cryptographic non-repudiation, replay defense, and job lease integrity for Milestone S3.

---

## 2. Test Catalog

| Test ID | Test Name | Target Invariant | Execution Rule / SQL Assertion | Expected Outcome |
|---|---|---|---|---|
| **SEC-01** | Empty DB Migration | Repeatability | Execute migrations 0001–0010 on clean database. | All 10 groups apply with 0 errors. |
| **SEC-02** | Version Tracking | Goose Integrity | Inspect `schema_migrations` table. | Exactly 10 version rows present with `is_applied = true`. |
| **SEC-03** | Exact Table Count | Inventory Parity | Query `information_schema.tables WHERE table_schema = 'public'`. | Count is **exactly 40**. |
| **SEC-04** | Foreign Key Rigor | Relational Integrity | Attempt orphaned row insert into `payment_attempts`. | Rejected with `foreign_key_violation` (`23503`). |
| **SEC-05** | App Direct SELECT Denial | Procedural Isolation | As `paybridge_app`, execute `SELECT * FROM payment_cases;`. | Rejected with `permission denied for table payment_cases` (`42501`). |
| **SEC-06** | Cross-Tenant Denial | Tenant Boundary | Tenant B token queries Tenant A case via `fn_cases_get_case`. | Rejected with `ERR_CASE_NOT_FOUND` / 0 rows. |
| **SEC-07** | HMAC Signature Check | Cryptographic Binding | Pass invalid HMAC signature to `_fn_verify_and_consume_assertion`. | Rejected with `ERR_INVALID_SIGNATURE` (`PB016`). |
| **SEC-08** | Nonce Replay Defense | Replay Prevention | Re-submit same token nonce in second transaction. | Rejected with `ERR_ASSERTION_REPLAYED` (`PB001`). |
| **SEC-09** | Expired Token Defense | Temporal Validity | Submit token with `expiry <= clock_timestamp()`. | Rejected with `ERR_ASSERTION_EXPIRED` (`PB011`). |
| **SEC-10** | Audit Monotonicity | Chain Integrity | Inspect `audit_events.seq` for Tenant A. | Sequence is strictly contiguous ($1, 2, 3 \dots$) with zero gaps. |
| **SEC-11** | Audit Tamper Detection | Hash Chain Non-repudiation | Verify `event_hash == digest(prev_hash || payload)`. | All hashes match linear chain back to genesis. |
| **SEC-12** | Snapshot Immutability | Historical Provenance | Attempt `UPDATE` or `DELETE` on `case_snapshots`. | Denied by role permissions and RLS. |
| **SEC-13** | Idempotency Concurrency | Race Condition Defense | Two concurrent requests with identical `(tenant_id, actor_id, route, key)`. | Exactly one succeeds (`200/201`), second returns `409` in-flight. |
| **SEC-14** | Outbox SKIP LOCKED | Concurrency | Two worker connections call `ClaimOutboxBatch` simultaneously. | Zero double-claims; distinct disjoint job sets returned. |
| **SEC-15** | Connector Unrelated Job Denial | Job Association Check | Connector presents `job_id` but requests different `connection_id`. | Rejected with `ERR_CONNECTION_MISMATCH` (`PB046`). |
| **SEC-16** | Expired Lease Rejection | Lease Expiry | Connector calls credentials on job with `locked_until < now()`. | Rejected with `ERR_JOB_LEASE_INVALID` (`PB045`). |
| **SEC-17** | Credential Access Audit | Audit Non-repudiation | Connector fetches credentials via `fn_connector_get_job_credentials`. | Event `CONNECTOR_CREDENTIAL_ACCESSED` logged in `audit_events`. |
| **SEC-18** | Retention Authorization | Data Minimization | `paybridge_retention` attempts purge on unapproved job. | Rejected; only `status = 'APPROVED'` allows deletion. |
| **SEC-19** | Capture Truth Preservation | Financial Invariant | Amount mismatch between invoice and capture response. | `provider_state` = `CAPTURE_COMPLETED`, `application_state` = `INCONSISTENT`. |
| **SEC-20** | Conflicting Webhook Reconciliation | Reconciliation | Out-of-order `CAPTURE_DENIED` webhook followed by capture. | Outbox `RECONCILE` job enqueued to query live provider status. |

---

## 3. Automated Test Execution
Run the full SQL test suite:
```powershell
psql -U paybridge_admin -d paybridge_dev -f tests/sql/verify_rls.sql
```
