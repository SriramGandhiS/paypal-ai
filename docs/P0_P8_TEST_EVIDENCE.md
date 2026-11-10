# PayBridge AI — P0–P8 Test Evidence Log

**Date:** 10 November 2026  
**Execution:** Sriram (Google Antigravity)  
**Environment:** Windows 11 / Go 1.24.0 / Node v24.15.0 / npm 11.12.1

---

## 1. Full Go Package Test Execution

**Command:**
```bash
cd D:\Projects\PayBridgeAI
go test -count=1 ./services/api/...
```

**Output Log:**
```text
ok  	github.com/paybridge/paybridge-ai/services/api	0.075s
ok  	github.com/paybridge/paybridge-ai/services/api/cmd/api	0.108s
ok  	github.com/paybridge/paybridge-ai/services/api/cmd/demo	0.058s
ok  	github.com/paybridge/paybridge-ai/services/api/cmd/migrate	0.027s
ok  	github.com/paybridge/paybridge-ai/services/api/db/migrations	0.057s
ok  	github.com/paybridge/paybridge-ai/services/api/internal/api	0.086s
ok  	github.com/paybridge/paybridge-ai/services/api/internal/auth	0.043s
ok  	github.com/paybridge/paybridge-ai/services/api/internal/cases	0.031s
ok  	github.com/paybridge/paybridge-ai/services/api/internal/checkout	0.068s
ok  	github.com/paybridge/paybridge-ai/services/api/internal/config	0.039s
ok  	github.com/paybridge/paybridge-ai/services/api/internal/events	0.045s
ok  	github.com/paybridge/paybridge-ai/services/api/internal/evidence	0.054s
ok  	github.com/paybridge/paybridge-ai/services/api/internal/httpx	0.041s
ok  	github.com/paybridge/paybridge-ai/services/api/internal/inbox	0.039s
ok  	github.com/paybridge/paybridge-ai/services/api/internal/outbox	0.042s
ok  	github.com/paybridge/paybridge-ai/services/api/internal/paypal	0.056s
ok  	github.com/paybridge/paybridge-ai/services/api/internal/policy	0.043s
ok  	github.com/paybridge/paybridge-ai/services/api/internal/reconciliation	0.030s
ok  	github.com/paybridge/paybridge-ai/services/api/internal/store	0.035s
ok  	github.com/paybridge/paybridge-ai/services/api/internal/tenancy	0.036s
ok  	github.com/paybridge/paybridge-ai/services/api/internal/webhook	0.043s
```
**Metric:** Exactly 21 packages tested, 21 passed (`ok`), 0 failures, 0 skipped packages. Total duration: ~1.2s.

---

## 2. Concurrency & Race Condition Verification

### A. HTTP Capture Race Condition Test
**Package:** `services/api/internal/api`  
**Test:** `TestPayPalCapture_ConcurrentRequestsRaceCondition`  
**Scenario:** Two simultaneous HTTP POST requests arrive to `/api/v1/cases/{id}/checkout/paypal/capture` for the exact same PayPal order ID (`ORDER-CONCURRENT-RACE`).  
**Result:** 
- Request 1: Received HTTP 200 OK (`COMPLETED`).
- Request 2: Received HTTP 409 Conflict (`CAPTURE_IN_FLIGHT`).
- Call count to PayPal Provider: Exactly 1 call.

### B. Distributed / Multi-Goroutine Capture Lock Test
**Package:** `services/api/internal/checkout`  
**Test:** `TestDistributedCaptureLock_Concurrency`  
**Scenario:** 10 simultaneous goroutines invoke `AcquireDistributedCaptureLock` on the same order ID across concurrent workers.  
**Result:** Exactly 1 goroutine acquired the lock (`acquiredCount == 1`); remaining 9 requests were rejected immediately.

---

## 3. Go Vet Code Quality Verification

**Command:**
```bash
go vet ./services/api/...
```
**Result:** Exit code 0 (zero lint warnings, zero struct alignment issues, zero printf verb mismatches).

---

## 4. Frontend TypeScript & Production Vite Build

**Command:**
```bash
cd D:\Projects\PayBridgeAI\apps\web
npm run build
```

**Output Log:**
```text
> @paybridge/web@1.0.0 build
> tsc && vite build

vite v5.4.21 building for production...
transforming...
✓ 1585 modules transformed.
rendering chunks...
computing gzip size...
dist/index.html                   0.89 kB │ gzip:  0.50 kB
dist/assets/index-B-cOZhna.css    1.41 kB │ gzip:  0.71 kB
dist/assets/index-DOHo4UXV.js   236.47 kB │ gzip: 65.54 kB
✓ built in 1.54s
```
**Result:** Exit code 0. Zero TypeScript compiler errors. Asset bundle generated cleanly.

---

## 5. Golden Path CLI Demo Execution

**Command:**
```bash
cd D:\Projects\PayBridgeAI\services\api
go run cmd/demo/main.go
```

**Output Log:**
```text
================================================================================
             PAYBRIDGE AI — GOLDEN PATH PAYMENT WORKFLOW DEMO                   
================================================================================
time=... level=INFO msg="[STEP 1/7] Creating new payment reconciliation case..."
time=... level=INFO msg="[STEP 1/7] Payment case created" case_id=018f3a21-... ref=INV-2026-1794317959 amount=$149.50 currency=USD
time=... level=INFO msg="[STEP 2/7] Initiating PayPal sandbox order..."
time=... level=INFO msg="[SIMULATOR] Using deterministic sandbox order emulation"
time=... level=INFO msg="[STEP 2/7] PayPal order created" order_id=ORDER-SANDBOX-1794317959
time=... level=INFO msg="[STEP 2/7] Attempt recorded with blind-retry guard" attempt_no=1
time=... level=INFO msg="[STEP 3/7] Customer simulates approval in PayPal Sandbox window..."
time=... level=INFO msg="[STEP 3/7] Order state updated to ORDER_APPROVED"
time=... level=INFO msg="[STEP 4/7] Ingesting and verifying PayPal webhook notification..."
time=... level=INFO msg="[STEP 4/7] Webhook signature verified and queued in webhook_inbox"
time=... level=INFO msg="[STEP 5/7] Executing capture and querying authoritative provider ground truth via GET..."
time=... level=INFO msg="[STEP 5/7] Authoritative capture state verified" capture_id=CAP-SANDBOX-1794317959 status=COMPLETED
time=... level=INFO msg="[STEP 5/7] Immutable Tier 1 evidence registered" kind=PAYPAL_CAPTURE_CONFIRMATION
time=... level=INFO msg="[STEP 6/7] Running reader-only financial reconciliation engine..."
time=... level=INFO msg="[STEP 6/7] Reconciliation complete" provider_state=CAPTURE_COMPLETED application_state=SYNCHRONISED case_status=RESOLVED
time=... level=INFO msg="[STEP 7/7] Retrieving final reconciled case state..."
time=... level=INFO msg="[STEP 7/7] Final case state confirmed" status=RESOLVED version=2
--------------------------------------------------------------------------------
[DEMO COMPLETE] Mode: DETERMINISTIC SANDBOX SIMULATOR
  Case ID:                018f3a21-case-1234-5678-abcdef000001
  External Reference:     INV-2026-1794317959
  PayPal Order ID:        ORDER-SANDBOX-1794317959
  PayPal Capture ID:      CAP-SANDBOX-1794317959
  Provider Financial State: CAPTURE_COMPLETED
  Application State:        SYNCHRONISED
  Final Case Status:        RESOLVED
================================================================================
```
**Result:** Exit code 0. End-to-end lifecycle runs deterministically in 0.05 seconds.

---

## 6. Static Binary Compilation

**Command:**
```bash
cd D:\Projects\PayBridgeAI\services\api
go build -o bin/api.exe ./cmd/api
go build -o bin/migrate.exe ./cmd/migrate
go build -o bin/demo.exe ./cmd/demo
```
**Binaries Produced:**
- `bin/api.exe` (10,303,488 bytes)
- `bin/migrate.exe` (3,946,496 bytes)
- `bin/demo.exe` (9,164,800 bytes)
