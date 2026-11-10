# PayBridge AI — Audit Blockers, Failures & Reproduction Guide

**Audit Target:** Vijay Kasthuri Submission  
**Evaluator:** Antigravity  
**Audit Branch:** `integration/vijay-agent-system`

---

## 1. Discovered Defects & Remediations

| Defect # | Component | Severity | Description | Remediation Applied |
| :--- | :--- | :---: | :--- | :--- |
| **DEF-01** | `worker.exe` | **HIGH** | Precompiled 5.6 MB executable binary committed into source ZIP archive. | **Quarantined.** Binary artifact was isolated and excluded from the repository. Worker is built exclusively from clean source `cmd/worker/main.go`. |
| **DEF-02** | `go.mod` Toolchain | **HIGH** | `go 1.27.3` specified in `go.mod`. Go toolchain halted attempting to download a non-existent version. | Corrected to `go 1.22`. Ran `go mod tidy` to declare direct dependencies (`chi/v5` and `cors`). |
| **DEF-03** | `go.mod` Namespace | **MEDIUM** | Vijay used module name `github.com/paybridge/api`, causing import failure inside monorepo `github.com/paybridge/paybridge-ai/services/api`. | Re-mapped all imports across 10 packages to canonical monorepo path. |
| **DEF-04** | Migration Drift | **MEDIUM** | Vijay included a partial 10-file migration set with only 10 tables, missing 30 enterprise tables. | Preserved Sriram's canonical 40-table migrations as authoritative. |
| **DEF-05** | Config Access | **LOW** | `cmd/worker/main.go` attempted to access non-existent field `cfg.LLMEnabled`. | Modified worker to read `os.Getenv("LLM_ENABLED") == "true"`, keeping fail-closed config untouched. |

---

## 2. External Infrastructure Blockers (Honest Disclosure)

As mandated by KERNEL rules, tests that could not run due to external environmental dependencies are classified as **BLOCKED BY EXTERNAL DEPENDENCY** rather than passed:

### Blocker 1: Host PostgreSQL 16 Installation Blocked on UAC Elevation
- **Condition:** Windows `winget install PostgreSQL.PostgreSQL.16` requires an interactive Administrator UAC elevation prompt. Standard non-elevated command sessions cannot install system services.
- **Impact:** SQL security tests and Goose migrations were verified structurally in Go tests (`TestEmbeddedMigrationsCount`, `TestCanonical40TablesInventory`, `TestRLSAndForceRLSEnforcement`), but cannot run against a live local port 5432 without administrator elevation.
- **Reproduction Step:**
  ```powershell
  # In an Administrator PowerShell terminal:
  winget install PostgreSQL.PostgreSQL.16 --mode unattended
  cd services/api
  go run cmd/migrate/main.go -action=up
  ```

### Blocker 2: Live PayPal WAN Execution Blocked on Real Credentials
- **Condition:** PayPal sandbox credentials in `.env.example` are placeholders (`sandbox_client_id_placeholder`). Real WAN HTTP calls to `api-m.sandbox.paypal.com` would return `401 Unauthorized`.
- **Impact:** Both Golden Path and Golden Demo B ran against deterministic sandbox emulation and high-fidelity mock transports.
- **Reproduction Step:**
  ```bash
  # Provide real developer sandbox credentials in .env:
  PAYPAL_ENV=sandbox
  PAYPAL_CLIENT_ID=<valid-sandbox-client-id>
  PAYPAL_CLIENT_SECRET=<valid-sandbox-client-secret>
  ```

### Blocker 3: Windows Race Detector Blocked on MinGW GCC
- **Condition:** On Windows, `go test -race` requires CGO and a MinGW-w64 GCC compiler. The host machine lacks a local C compiler.
- **Impact:** Local race tests output `go: -race requires cgo`. Full race detection executes on Linux CI (`ubuntu-latest`) via `.github/workflows/go-tests.yml`.
