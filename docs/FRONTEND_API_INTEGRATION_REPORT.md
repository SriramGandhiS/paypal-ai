# PAYBRIDGE AI — FRONTEND API INTEGRATION REPORT

**Document Version:** 2.0  
**Target Branch:** `feature/frontend-complete`  
**API Protocol:** REST / JSON over HTTP/1.1  
**Default Base URL:** `http://127.0.0.1:8080/api/v1`  
**Client Implementation:** [`apps/web/src/api/client.ts`](file:///D:/Projects/PayBridgeAI/apps/web/src/api/client.ts)

---

## 1. INTEGRATION ARCHITECTURE

The PayBridge AI frontend uses a resilient, dual-surface integration architecture:
1. **Live Backend Mode (Default):**
   - Dispatches HTTP requests to the Go API running at `/api/v1`.
   - Sends Bearer authentication headers and client-generated `Idempotency-Key` headers for mutating operations.
   - Respects HTTP response status codes: 200/201 (Success), 401 (Unauthorized), 403 (Forbidden), 404 (Not Found), 500/504 (Server/Gateway Error).
2. **Deterministic Sandbox Simulator Mode:**
   - Explicitly surfaced in the navigation bar with a high-contrast mode indicator pill (`LIVE API CONNECTED` vs `SIMULATOR MODE`).
   - If the Go API is not running on localhost (e.g. during an offline presentation or review), the client catches network exceptions (`TypeError: Failed to fetch`) and gracefully enables simulation mode with high-fidelity seeded cases (`c101`, `c102`, `c103`).
   - **Zero Silent Switching:** The mode indicator explicitly alerts the user, ensuring full transparency.

---

## 2. API ENDPOINT COVERAGE MATRIX

| Route | Method | Go Backend Handler | Frontend Client Method | Live Behavior | Error Handling |
| :--- | :---: | :--- | :--- | :--- | :--- |
| `/healthz` | `GET` | Root Health Check | `api.checkBackendLiveness()` | Returns 200 OK | Falls back to simulator if unreachable |
| `/cases` | `GET` | `cases.ListCasesHandler` | `api.listCases()` | Returns paginated case summaries | Returns 401 if unauthenticated; falls back to seed |
| `/cases/:id` | `GET` | `cases.GetCaseHandler` | `api.getCase(id)` | Returns full case with state vector | Returns 404 if case does not exist |
| `/cases` | `POST` | `cases.CreateCaseHandler` | `api.createCase(input)` | Creates new reconciliation case | Requires `Idempotency-Key`; validates inputs |
| `/cases/:id/attempts` | `GET` | `checkout.ListAttemptsHandler` | `api.listAttempts(caseId)` | Returns provider attempt records | Returns empty array if no attempts |
| `/cases/:id/evidence` | `GET` | `evidence.ListEvidenceHandler` | `api.listEvidence(caseId)` | Returns Tier 1–6 evidence items | Cites SHA-256 hashes |
| `/cases/:id/safe-actions` | `GET` | `safeaction.EvaluateHandler` | `api.getSafeActions(caseId)` | Computes allowed & blocked actions | Invariant: Blocks `RETRY_CHECKOUT` on 504 |
| `/webhooks/paypal` | `POST` | `webhook.PayPalHandler` | Replay trigger in Operations screen | Verifies webhook idempotency | Rejects duplicate webhook IDs with HTTP 200 |

---

## 3. ERROR & SAFETY PROTOCOL VERIFICATION

### 3.1 HTTP 504 Gateway Timeout Handling
- When PayPal or an upstream payment rail fails to respond before the timeout deadline, the frontend **never** interprets this as a payment failure.
- The UI immediately renders the **`DO_NOT_REPEAT_CHECKOUT`** alert banner.
- The "Retry Payment" button is permanently disabled with the message: *"Retry Blocked by Safety Invariant: Active attempt is outcome-unknown"*.

### 3.2 Tenant & Role Isolation
- Sensitive recipient banking numbers (IFSC, SWIFT, IBAN) are strictly redacted on the client when viewed from the **Sender Lens** (`[REDACTED-BANK-ACCOUNT]`).
- The **Receiver Lens** enforces statutory consent verification. If consent is revoked by the contractor, the counterparty payment ledger projection is immediately cleared.

### 3.3 Zero AI Mutations
- No autonomous agent output, LLM suggestion, or client interaction has the capability to execute a payment capture, refund, or retry.
- AI specialist agents (A01–A10) produce structured epistemic findings (`confirmed`, `inference`, `unknown`), leaving all financial actions to deterministic Go backend controllers.
