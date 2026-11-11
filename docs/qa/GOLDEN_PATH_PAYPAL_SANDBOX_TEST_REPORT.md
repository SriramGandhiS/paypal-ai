# Golden Path Demo — PayPal Sandbox Integration & Safety Report

**Date:** 2026-11-11  
**Authoritative Backend:** `services/api/internal/paypal` & `services/api/internal/checkout`  

---

## 1. PayPal Sandbox Architecture & Integration Status

1. **Authentication:** Server-side OAuth 2.0 token acquisition via PayPal Sandbox endpoints (`https://api-m.sandbox.paypal.com/v1/oauth2/token`). Client secrets remain securely on the backend.
2. **Order Creation:** `POST /api/v1/cases/{id}/checkout/paypal` constructs authoritative v2 checkout orders with `EXTERNAL_REFERENCE` and idempotent case context.
3. **Capture Execution:** Authoritative provider capture at `POST /api/v1/cases/{id}/checkout/paypal/capture` with server-side validation.
4. **Webhook Gateway:** Cryptographically verifies incoming PayPal webhooks (`CHECKOUT.ORDER.APPROVED`, `PAYMENT.CAPTURE.COMPLETED`).

---

## 2. Mandatory Financial Safety Invariants

| Safety Invariant | Implementation Status | Verification |
| :--- | :--- | :--- |
| **HTTP 504 Is Not Failure** | **ENFORCED** | Timeout does not trigger cancellation or blind retries. State remains `UNKNOWN_IN_FLIGHT` pending provider ledger read. |
| **Attempt Mutex Locking** | **ENFORCED** | `DO_NOT_REPEAT_CHECKOUT` mutex lock blocks customer retry clicks during flight. |
| **Zero Blind Retries** | **ENFORCED** | Frontend buttons disable re-checkout while verification is underway. |
| **Authoritative Reconciliation** | **ENFORCED** | Provider query `/v2/checkout/orders/{id}` verifies capture before any refund or resolution path is authorized. |
| **AI Action Boundary** | **ENFORCED** | AI Copilot is strictly read-only and advisory; agents cannot execute money movement or capture operations. |
