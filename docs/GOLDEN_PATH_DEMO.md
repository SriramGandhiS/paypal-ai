# PayBridge AI — Golden Path Payment Lifecycle & Demonstration Guide

**Milestone Target:** S4 Backend & Payment Workflow  
**Demonstration Runner:** [`services/api/cmd/demo/main.go`](file:///D:/Projects/PayBridgeAI/services/api/cmd/demo/main.go)  
**Automated Verification:** [`services/api/cmd/demo/demo_test.go`](file:///D:/Projects/PayBridgeAI/services/api/cmd/demo/demo_test.go)

---

## 1. Executive Summary

The **PayBridge AI Golden Path Demonstration** validates the complete end-to-end payment reconciliation lifecycle:
```
┌─────────────────┐       ┌─────────────────┐       ┌─────────────────┐       ┌─────────────────┐
│     STEP 1      │       │     STEP 2      │       │     STEP 3      │       │     STEP 4      │
│  Create Payment │ ────> │ Initiate PayPal │ ────> │ Customer Sample │ ────> │ Verified PayPal │
│  Reconcil. Case │       │ Sandbox Order   │       │ Approval Link   │       │ Webhook Ingest  │
└─────────────────┘       └─────────────────┘       └─────────────────┘       └─────────────────┘
                                                                                       │
                                                                                       ▼
┌─────────────────┐       ┌─────────────────┐       ┌─────────────────┐       ┌─────────────────┐
│     STEP 7      │       │     STEP 6      │       │     STEP 5      │       │   (Continued)   │
│ Final Reconciled│ <──── │   Reader-Only   │ <──── │ Authoritative   │ <─────┘                 │
│ State & Passport│       │ Reconciliation  │       │ Capture Fetch   │                         │
└─────────────────┘       └─────────────────┘       └─────────────────┘
```

The system operates in two distinct, clearly labelled modes:
1. **Deterministic Sandbox Simulation Mode:** Runs instantly with zero external dependencies or internet connectivity, producing guaranteed mathematical and state vector transitions.
2. **Live PayPal Sandbox Mode:** Automatically engages when valid `PAYPAL_CLIENT_ID` and `PAYPAL_CLIENT_SECRET` environment variables are configured, communicating directly with PayPal's sandbox endpoint (`api-m.sandbox.paypal.com`).

---

## 2. The 7-Step Payment Lifecycle

### Step 1: Merchant Creates Payment Case
* **Endpoint:** `POST /api/v1/cases`
* **Security Context:** Authenticated tenant (`018f3a21-...`), sender user, role `sender` or `merchant_admin`.
* **Payload:**
  ```json
  {
    "external_reference": "INV-2026-DEMO-001",
    "amount_cents": 14950,
    "currency": "USD",
    "sender_country": "US",
    "receiver_country": "DE",
    "payment_purpose": "CONSULTING_SERVICES",
    "idempotency_key": "idem-key-991"
  }
  ```
* **Initial State:** `provider_state = UNKNOWN`, `application_state = INCONSISTENT`, `case_status = OBSERVING`.

### Step 2: Customer Initiates PayPal Sandbox Checkout
* **Endpoint:** `POST /api/v1/cases/{id}/checkout/paypal`
* **Actions:**
  - Creates PayPal Order with `intent: CAPTURE` and purchase unit `$149.50 USD`.
  - Obtains `order_id` and buyer `approval_url`.
  - Atomically creates `payment_attempts` record (attempt #1) while enforcing the **blind-retry invariant** (rejecting duplicate active attempts with HTTP 409 Conflict).

### Step 3: Customer Approves Sandbox Payment
* **Action:** Buyer approves payment in PayPal sandbox interface (`PAYER_ACTION_REQUIRED` $\to$ `APPROVED`).
* **Attempt State:** Updated to `ORDER_APPROVED`.

### Step 4: PayPal Emits Cryptographically Verified Webhook
* **Endpoint:** `POST /api/v1/webhooks/paypal`
* **Security & Ingestion:**
  - Raw payload buffered with 1MB size limit.
  - Validates `PAYPAL-TRANSMISSION-SIG`, `PAYPAL-CERT-URL`, and `PAYPAL-AUTH-ALGO` via PayPal API or deterministic cryptographic verifier.
  - Queued into `webhook_inbox` with status `RECEIVED`, deduplicated on `(environment, provider, provider_event_id)`.

### Step 5: Authoritative Capture & Ground Truth Query
* **Endpoint:** `POST /api/v1/cases/{id}/checkout/paypal/capture`
* **Idempotency & Safety Protections:**
  - **Inspection Before Mutation:** Queries `GET /v2/checkout/orders/{id}` first. If PayPal already reports `COMPLETED`, capture is not re-attempted.
  - **Timeout Recovery:** If the capture call encounters a network timeout (504), queries `GET /v2/checkout/orders/{id}` to verify whether the capture actually completed at PayPal during the failure.
  - **Direct Ground Truth Fetch:** Calls `GET /v2/payments/captures/{id}` to obtain provider-signed amount and status.
  - **Immutable Evidence Ingestion:** Registers Tier 1 `PAYPAL_CAPTURE_CONFIRMATION` evidence item with hex SHA-256 payload digest.

### Step 6: Reader-Only Reconciliation Engine
* **Endpoint:** `POST /api/v1/cases/{id}/reconcile`
* **Invariant Enforced:**
  - **Financial Ground Truth Invariant:** If provider captured funds, `provider_state` is set to `CAPTURE_COMPLETED` unconditionally.
  - Any merchant invoice discrepancy surfaces in `application_state = INCONSISTENT` and `case_status = NEEDS_HUMAN_REVIEW` without altering provider financial reality.
  - On matched amounts: `application_state = SYNCHRONISED`, `case_status = RESOLVED`.

### Step 7: Final Reconciled Case State
* **Endpoint:** `GET /api/v1/cases/{id}`
* **Final Vector:**
  - `provider_state`: `CAPTURE_COMPLETED`
  - `application_state`: `SYNCHRONISED`
  - `case_status`: `RESOLVED`
  - `state_version`: incremented atomically (`2`)

---

## 3. How to Run the Demonstration

### A. Run via CLI Demo Runner
From the repository root:
```powershell
go run ./services/api/cmd/demo
```
Expected output:
```text
================================================================================
             PAYBRIDGE AI — GOLDEN PATH PAYMENT WORKFLOW DEMO                   
================================================================================
[STEP 1/7] Creating new payment reconciliation case...
[STEP 1/7] Payment case created: case_id=018f3a21-... ref=INV-2026-... amount=$149.50
[STEP 2/7] Initiating PayPal sandbox order...
[STEP 2/7] PayPal order created: order_id=ORDER-SANDBOX-... approval_url=...
[STEP 2/7] Attempt recorded with blind-retry guard: attempt_no=1
[STEP 3/7] Customer simulates approval in PayPal Sandbox window...
[STEP 3/7] Order state updated to ORDER_APPROVED
[STEP 4/7] Ingesting and verifying PayPal webhook notification...
[STEP 4/7] Webhook signature verified and queued in webhook_inbox: tx=tx-paypal-...
[STEP 5/7] Executing capture and querying authoritative provider ground truth via GET...
[STEP 5/7] Authoritative capture state verified: cap_id=CAP-SANDBOX-... status=COMPLETED
[STEP 5/7] Immutable Tier 1 evidence registered: sha256=...
[STEP 6/7] Running reader-only financial reconciliation engine...
[STEP 6/7] Reconciliation complete: provider_state=CAPTURE_COMPLETED app_state=SYNCHRONISED case_status=RESOLVED
[STEP 7/7] Retrieving final reconciled case state...
[STEP 7/7] Final case state confirmed: status=RESOLVED version=2
--------------------------------------------------------------------------------
[DEMO COMPLETE] Mode: DETERMINISTIC SANDBOX SIMULATOR
  Case ID:                018f3a21-case-1234-5678-abcdef000001
  External Reference:     INV-2026-1794316132
  PayPal Order ID:        ORDER-SANDBOX-1794316132
  Customer Approval Link: https://www.sandbox.paypal.com/checkoutnow?token=...
  PayPal Capture ID:      CAP-SANDBOX-1794316132
  Webhook Transmission:   tx-paypal-... (Verified)
  Provider Financial State: CAPTURE_COMPLETED
  Application State:        SYNCHRONISED
  Final Case Status:        RESOLVED
================================================================================
```

### B. Run Automated Verification Test
```powershell
go test -v ./services/api/cmd/demo
```

---

## 4. Hardened Security & Safety Invariants Verified

| Invariant / Guard | Implementation Verification | Test Coverage |
|---|---|---|
| **Dev Header Bypass Prohibition** | `X-Dev-Tenant-ID` is strictly disabled unless `ALLOW_DEV_AUTH=true` is explicitly set in non-production. Attempts return HTTP 403 Forbidden. | `TestDevAuthBypassDisabled_RejectsXDevHeaders` |
| **Cross-Tenant Access Denial** | All database queries and procedure calls enforce tenant scoping via context and RLS. | `TestCasesContextRequirement` |
| **AI Payment Mutation Guard** | AI roles (`ai_agent`, `investigator_agent`) have zero payment write permissions. Attempts return HTTP 403 Forbidden. | `TestNoAIInitiationMutation_Forbidden` |
| **Capture Idempotency** | Orders already reported `COMPLETED` by provider are not re-captured. | `TestPayPalCapture_IdempotentDuplicateCapture` |
| **Timeout Recovery via GET** | Upstream 504 gateway timeouts trigger provider state recovery via `GET /v2/checkout/orders/{id}`. | `TestPayPalCapture_TimeoutRecoveryViaGetOrder` |
| **Financial Ground Truth Preservation** | Provider `CAPTURE_COMPLETED` is never demoted, even if merchant invoice amount differs. | `TestReconciliationInvariants` |
| **Blind-Retry Invariant** | At most one active non-terminal attempt per case. | `TestCheckoutRoleEnforcement` |
