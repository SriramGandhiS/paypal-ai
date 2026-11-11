# 09 — Chaos Engineering Execution Report (CHAOS-01 to CHAOS-12)
**Phase:** 9 — Chaos Engineering  
**Owner:** Dinny Paul Navis  
**Branch:** `feat/diny-enterprise-completion`  
**Test Suite:** `services/api/tests/chaos/chaos_test.go`  
**Date:** 2026-11-11  

---

## 1. Chaos Test Results Summary

| Scenario ID | Test Name | Injected Failure Mode | Verified System Invariant | Result |
|---|---|---|---|---|
| **CHAOS-01** | `TestChaos_01_CaptureResponseDrop` | 504 Gateway Timeout on capture response | Raises `DO_NOT_REPEAT_CHECKOUT`; `ProceedToCheckout` blocked | **PASS** |
| **CHAOS-02** | `TestChaos_02_DelayedWebhook` | Out-of-order webhook arriving 10m late | Deterministic Reducer rejects regression from `CaptureCompleted` | **PASS** |
| **CHAOS-03** | `TestChaos_03_WebhookVerificationOutage` | Corrupted signature header on webhook | Verification fails safely; webhook rejected before processing | **PASS** |
| **CHAOS-04** | `TestChaos_04_ForgedWebhook` | Attacker forged signature hash | HMAC-SHA256 signature mismatch rejected immediately | **PASS** |
| **CHAOS-05** | `TestChaos_05_DuplicateWebhookFlood` | Repeated submission of identical event payload | Replay defense rejects duplicate event SHA-256 | **PASS** |
| **CHAOS-06** | `TestChaos_06_CrossTenantIDOR` | Attacker session attempts cross-tenant query | `AuthorizeTenantAction` blocks access with `ErrForbidden` | **PASS** |
| **CHAOS-07** | `TestChaos_07_PromptInjection` | Prompt injection payload attempting refund override | Regex & pattern guardrail blocks prompt before SLM inference | **PASS** |
| **CHAOS-08** | `TestChaos_08_MultiRegionPrimaryPartition` | Network partition leaves outcome unknown | Safe action engine blocks checkout until reconciliation succeeds | **PASS** |
| **CHAOS-09** | `TestChaos_09_ConsentRevocation` | Open Banking consent token revoked by user | Plaid `/transactions/sync` fails closed with `ErrConsentExpiredOrRevoked` | **PASS** |
| **CHAOS-10** | `TestChaos_10_AIClusterUnavailable` | Private vLLM port dead / connection refused | Instant deterministic fallback (`GAP_01_LOST_CONFIRMATION`), 0ms hallucination | **PASS** |
| **CHAOS-11** | `TestChaos_11_BankSettlementDiscrepancy` | Bank statement amount differs from invoice | `DepositMatcher` rejects discrepancy; `receiver_state` un-transitioned | **PASS** |
| **CHAOS-12** | `TestChaos_12_KMSDecryptionFailure` | 1-bit ciphertext bit-flip tampering | AES-256-GCM authentication tag check fails closed | **PASS** |
| **BONUS-01** | `TestChaos_Bonus_UnbalancedLedgerRejected` | Debit != Credit journal entry | Double-entry ledger rejects unbalanced transaction before commit | **PASS** |
| **BONUS-02** | `TestChaos_Bonus_SanctionsHoldEnforced` | OFAC SDN matched entity in pre-auth | `PreAuthorizationGate` freezes transaction; AI override blocked | **PASS** |

---

## 2. Test Execution Output

```text
=== RUN   TestChaos_01_CaptureResponseDrop
--- PASS: TestChaos_01_CaptureResponseDrop (0.00s)
=== RUN   TestChaos_02_DelayedWebhook
--- PASS: TestChaos_02_DelayedWebhook (0.00s)
=== RUN   TestChaos_03_WebhookVerificationOutage
--- PASS: TestChaos_03_WebhookVerificationOutage (0.00s)
=== RUN   TestChaos_04_ForgedWebhook
--- PASS: TestChaos_04_ForgedWebhook (0.00s)
=== RUN   TestChaos_05_DuplicateWebhookFlood
--- PASS: TestChaos_05_DuplicateWebhookFlood (0.00s)
=== RUN   TestChaos_06_CrossTenantIDOR
--- PASS: TestChaos_06_CrossTenantIDOR (0.00s)
=== RUN   TestChaos_07_PromptInjection
--- PASS: TestChaos_07_PromptInjection (0.00s)
=== RUN   TestChaos_08_MultiRegionPrimaryPartition
--- PASS: TestChaos_08_MultiRegionPrimaryPartition (0.00s)
=== RUN   TestChaos_09_ConsentRevocation
--- PASS: TestChaos_09_ConsentRevocation (0.00s)
=== RUN   TestChaos_10_AIClusterUnavailable
--- PASS: TestChaos_10_AIClusterUnavailable (0.01s)
=== RUN   TestChaos_11_BankSettlementDiscrepancy
--- PASS: TestChaos_11_BankSettlementDiscrepancy (0.00s)
=== RUN   TestChaos_12_KMSDecryptionFailure
--- PASS: TestChaos_12_KMSDecryptionFailure (0.00s)
=== RUN   TestChaos_Bonus_UnbalancedLedgerRejected
--- PASS: TestChaos_Bonus_UnbalancedLedgerRejected (0.00s)
=== RUN   TestChaos_Bonus_SanctionsHoldEnforced
--- PASS: TestChaos_Bonus_SanctionsHoldEnforced (0.00s)
PASS
ok  	paybridge-ai/services/api/tests/chaos	3.267s
```

