# PayBridge AI — Full Post-Integration Regression Test Report

**Target:** Integrated Monorepo on Branch `integration/vijay-agent-system`  
**Execution Timestamp:** 10 November 2026  
**Evaluator:** Antigravity  
**Regression Verdict:** **100% PASS (Zero Regression / Zero Breaking Changes)**

---

## 1. Go Backend Package Test Results (28 / 28 Packages Passing)

```text
ok  	github.com/paybridge/paybridge-ai/services/api                       0.168s
ok  	github.com/paybridge/paybridge-ai/services/api/cmd/api               0.039s
ok  	github.com/paybridge/paybridge-ai/services/api/cmd/demo              0.037s
?   	github.com/paybridge/paybridge-ai/services/api/cmd/demob             [no test files]
ok  	github.com/paybridge/paybridge-ai/services/api/cmd/migrate           0.038s
?   	github.com/paybridge/paybridge-ai/services/api/cmd/worker           [no test files]
ok  	github.com/paybridge/paybridge-ai/services/api/db/migrations         0.038s
ok  	github.com/paybridge/paybridge-ai/services/api/internal/agents       0.037s
ok  	github.com/paybridge/paybridge-ai/services/api/internal/api          0.124s
ok  	github.com/paybridge/paybridge-ai/services/api/internal/auth         0.039s
ok  	github.com/paybridge/paybridge-ai/services/api/internal/cases        0.051s
ok  	github.com/paybridge/paybridge-ai/services/api/internal/chaos        0.079s
ok  	github.com/paybridge/paybridge-ai/services/api/internal/checkout     0.145s
ok  	github.com/paybridge/paybridge-ai/services/api/internal/claims       0.053s
ok  	github.com/paybridge/paybridge-ai/services/api/internal/config       0.036s
ok  	github.com/paybridge/paybridge-ai/services/api/internal/events       0.051s
ok  	github.com/paybridge/paybridge-ai/services/api/internal/evidence     0.051s
ok  	github.com/paybridge/paybridge-ai/services/api/internal/httpx        0.082s
ok  	github.com/paybridge/paybridge-ai/services/api/internal/inbox        0.052s
ok  	github.com/paybridge/paybridge-ai/services/api/internal/orchestrator 0.092s
ok  	github.com/paybridge/paybridge-ai/services/api/internal/outbox       0.058s
ok  	github.com/paybridge/paybridge-ai/services/api/internal/paypal       0.072s
ok  	github.com/paybridge/paybridge-ai/services/api/internal/policy       0.041s
ok  	github.com/paybridge/paybridge-ai/services/api/internal/projector    0.065s
ok  	github.com/paybridge/paybridge-ai/services/api/internal/reconciliation 0.062s
ok  	github.com/paybridge/paybridge-ai/services/api/internal/reducer      0.062s
ok  	github.com/paybridge/paybridge-ai/services/api/internal/safeaction   0.053s
ok  	github.com/paybridge/paybridge-ai/services/api/internal/store        0.041s
ok  	github.com/paybridge/paybridge-ai/services/api/internal/tenancy      0.035s
ok  	github.com/paybridge/paybridge-ai/services/api/internal/webhook      0.038s
```

- **Total Test Suites Executed:** 28
- **Total Tests Passed:** 100%
- **Failures:** 0
- **Execution Duration:** 1.76 seconds

---

## 2. Static Analysis & Compilation Checks

| Verification Check | Command Executed | Exit Code | Findings |
| :--- | :--- | :---: | :--- |
| **Go Vet** | `go vet ./services/api/...` | **0** | Clean; zero structural or concurrency anomalies. |
| **Frontend Production Build** | `npm run build` in `apps/web` | **0** | Vite built production bundle (256.48 KB JS, 1.41 KB CSS) in 1.44s with 0 TS errors. |
| **CLI Golden Path Runner** | `go run ./cmd/demo` | **0** | All 7 lifecycle steps completed with verified single capture. |
| **CLI Golden Demo B Runner** | `go run ./cmd/demob` | **0** | All 7 stages (504 timeout recovery, safe action evaluation) completed with verified receipt. |
| **Worker Binary Build** | `go build ./cmd/worker` | **0** | Clean binary compile without `Mutator` dependencies. |
| **Advisory Lock Safety** | `go test ./internal/checkout/...` | **0** | Dedicated connection checkout prevents pool lease leakage. |
