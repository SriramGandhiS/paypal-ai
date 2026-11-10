# PayBridge AI — QA Execution Command Log

**Audit Session:** P0–P8 Extreme Adversarial Audit  
**Workspace:** `D:\Projects\PayBridgeAI`  
**Date:** 10 November 2026

---

## Chronological Command Execution Journal

### Step 1: Directory Setup
```bash
mkdir -p "D:\Projects\PayBridgeAI\docs\qa\P0-P8-EXTREME-AUDIT\evidence"
# Exit Code: 0
```

### Step 2: Checkout Concurrency & Multi-Instance Lock Tests
```bash
cd D:\Projects\PayBridgeAI
go test -v ./services/api/internal/checkout/...
# Exit Code: 0 | Duration: 0.116s
# Verdict: PASS (All 4 multi-instance and failure tests passed)
```

### Step 3: Config Loader & Security Parameter Tests
```bash
cd D:\Projects\PayBridgeAI
go test -v ./services/api/internal/config/...
# Exit Code: 0 | Duration: 0.021s
# Verdict: PASS (All 28 table-driven tests passed)
```

### Step 4: Python Character Scan & Migration Corrupted Character Detection
```bash
python -c "content = open('D:/Projects/PayBridgeAI/services/api/db/migrations/0010_hardening_rls_and_procs.sql', encoding='utf-8', errors='ignore').readlines(); ..."
# Exit Code: 0
# Result: Discovered corrupted character \u4fe1 on lines 31-32; repaired to ''''
```

### Step 5: Database Migrations Structural & 40-Table Inventory Test
```bash
cd D:\Projects\PayBridgeAI
go test -v ./services/api/db/migrations/...
# Exit Code: 0 | Duration: 0.024s
# Verdict: PASS (All 4 migration tests passed; 40 canonical tables verified)
```

### Step 6: API Handlers, Capture Idempotency & Concurrency Race Tests
```bash
cd D:\Projects\PayBridgeAI
go test -v ./services/api/internal/api/...
# Exit Code: 0 | Duration: 0.058s
# Verdict: PASS (All 10 API tests passed)
```

### Step 7: Race Detector Check on Windows
```bash
cd D:\Projects\PayBridgeAI
go test -race ./services/api/internal/checkout/...
# Exit Code: 1
# Output: go: -race requires cgo; enable cgo by setting CGO_ENABLED=1
# Note: Documented Windows CGO/MinGW requirement in race-tests.log
```

### Step 8: Capture Full Go Test Log (Evidence)
```bash
cd D:\Projects\PayBridgeAI
go test -v -count=1 ./services/api/... | Out-File -FilePath "docs\qa\P0-P8-EXTREME-AUDIT\evidence\go-tests.log" -Encoding utf8
# Exit Code: 0 | Lines: 441
```

### Step 9: Capture Build Evidence
```bash
cd services/api
go build -v -o bin/api.exe ./cmd/api
go build -v -o bin/migrate.exe ./cmd/migrate
go build -v -o bin/demo.exe ./cmd/demo
cd ../../apps/web
npm run build
# Exit Code: 0 | Output written to evidence/build.log
```

### Step 10: Capture Migration & Security Scan Logs
```bash
go test -v ./services/api/db/migrations/... > evidence/migration-tests.log
git check-ignore .env > evidence/security-scans.log
go test -v ./services/api/internal/policy/... >> evidence/security-scans.log
# Exit Code: 0
```
