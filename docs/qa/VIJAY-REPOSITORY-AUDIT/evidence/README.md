# PayBridge AI — Vijay Audit Raw Evidence Vault

**Audit Target:** Vijay Kasthuri Final Submitted Repository  
**Source ZIP:** `C:\Users\iamra\Downloads\PayBridge-AI-Vijay-Final-Verified-Repository.zip`  
**Integration Branch:** `integration/vijay-agent-system`  
**Execution Timestamp:** 10 November 2026

---

## Evidence Manifest

| File | Purpose | Source Command / Verification |
| :--- | :--- | :--- |
| [`archive-inspection.log`](file:///D:/Projects/PayBridgeAI/docs/qa/VIJAY-REPOSITORY-AUDIT/evidence/archive-inspection.log) | Archive integrity, path traversal scan, quarantined executable log | Python `zipfile` inspector & SHA-256 hash calculation |
| [`secret-scan.log`](file:///D:/Projects/PayBridgeAI/docs/qa/VIJAY-REPOSITORY-AUDIT/evidence/secret-scan.log) | Entropy, private key, and credential leakage audit | Python regex scanner across all 87 extracted files |
| [`go-tests.log`](file:///D:/Projects/PayBridgeAI/docs/qa/VIJAY-REPOSITORY-AUDIT/evidence/go-tests.log) | Full verbose execution of all 28 Go backend packages post-integration | `go test -v -count=1 ./services/api/...` |
| [`chaos-tests.log`](file:///D:/Projects/PayBridgeAI/docs/qa/VIJAY-REPOSITORY-AUDIT/evidence/chaos-tests.log) | Complete CHAOS-01 through CHAOS-12 fault injection results | `go test -v -count=1 ./services/api/internal/chaos/...` |
| [`dag-scenarios.log`](file:///D:/Projects/PayBridgeAI/docs/qa/VIJAY-REPOSITORY-AUDIT/evidence/dag-scenarios.log) | 20 mandatory agent DAG adversarial scenario test executions | `go test -v -count=1 -run=TestScenario ./services/api/internal/orchestrator/...` |
| [`perf-benchmark.log`](file:///D:/Projects/PayBridgeAI/docs/qa/VIJAY-REPOSITORY-AUDIT/evidence/perf-benchmark.log) | Stress matrix (1, 5, 10, 25 concurrent runs) & 1,000 item payload tests | `go test -v -count=1 -run=TestPerformance ./services/api/internal/orchestrator/...` |
| [`frontend-build.log`](file:///D:/Projects/PayBridgeAI/docs/qa/VIJAY-REPOSITORY-AUDIT/evidence/frontend-build.log) | TypeScript typecheck and Vite production build of all 16 screens | `npm run build` in `apps/web` |
