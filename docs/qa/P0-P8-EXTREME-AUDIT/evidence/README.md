# PayBridge AI — QA Evidence Vault Index

**Directory:** `docs/qa/P0-P8-EXTREME-AUDIT/evidence/`  
**Execution Timestamp:** 10 November 2026 20:12 UTC+5:30  
**Environment:** Windows 11 Enterprise (AMD64) / Go 1.24.0 / Node.js v24.15.0 / npm 11.12.1

---

## 1. Machine-Readable Evidence Files

| File Name | Test Target | Generation Command | Key Metric / Verdict |
|---|---|---|:---:|
| [`go-tests.log`](./go-tests.log) | Complete Go backend test suite (21 packages) | `go test -v -count=1 ./services/api/...` | **441 lines, 21/21 packages PASS** |
| [`race-tests.log`](./race-tests.log) | Race detector execution attempt | `go test -race ./services/api/...` | **Documents Windows MinGW/CGO requirement** |
| [`build.log`](./build.log) | Go binaries & Vite production bundle | `go build` (3 binaries) + `npm run build` | **3 binaries built, 0 Vite errors (1.50s)** |
| [`migration-tests.log`](./migration-tests.log) | Goose embedded migrations & 40-table check | `go test -v ./services/api/db/migrations/...` | **40 canonical tables verified** |
| [`security-scans.log`](./security-scans.log) | Secret scan, AI mutation denial, auth bypass | Gitleaks scan + `policy_test.go` + `api_test.go` | **.env ignored, AI mutations blocked** |

---

## 2. Integrity & Zero-Credential Attestation

All evidence files in this directory:
- Contain genuine, unaltered command outputs executed inside the workspace.
- Do NOT contain live API secrets, private keys, session tokens, or unredacted passwords.
- Are committed directly to version control for independent audit reproduction.
