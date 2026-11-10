# PayBridge AI — P0–P8 External Blockers & Resolution Runbook

**Date:** 10 November 2026  
**Document:** `P0_P8_BLOCKERS.md`  
**Author:** Sriram (Google Antigravity)  
**Strict Honesty Standard:** In accordance with Kernel Rule 5, all external dependencies and unavailable runtime infrastructure are documented transparently without evasion.

---

## 1. Blocker Inventory

| # | Dependency | Blocker Description | Required Permission / External Resource | Impact |
|:---:|---|---|---|---|
| **B1** | **Host PostgreSQL 16 Engine** | PostgreSQL 16 is not installed on the local Windows OS; `winget` installation requires Administrator UAC elevation. Docker daemon is also not installed. | Windows Administrator Elevation or running Docker Desktop. | Migrations 0001–0010 & SQL RLS tests (`tests/sql/verify_rls.sql`) cannot execute on live host port 5432. |
| **B2** | **Live PayPal Sandbox WAN Access** | PayPal Client ID & Secret currently set to development sandbox defaults. Live WAN roundtrips to PayPal sandbox API require active developer credentials. | PayPal Developer Portal App credentials (`PAYPAL_CLIENT_ID`, `PAYPAL_CLIENT_SECRET`). | Backend operates against unit test HTTP mocks and deterministic simulator fixtures rather than live PayPal WAN servers. |
| **B3** | **Teammate Code Merges (Vijay & Diny)** | Vijay's autonomous AI agent DAG (A01–A10) and Diny's Kafka event bus (>10k TPS) are developed in separate workspaces and not yet merged into `master`. | Git branch merge / PR review when teammates complete their respective modules. | AI orchestrator and Kafka event streaming are isolated to stable interface boundaries. |

---

## 2. Exact Resolution Steps for Blocker B1 (PostgreSQL 16)

When an administrator is available to provision PostgreSQL 16, run either Option A (Native Windows) or Option B (Docker):

### Option A: Native Windows PostgreSQL 16 (Recommended)
1. Open PowerShell **as Administrator**.
2. Execute:
   ```powershell
   winget install PostgreSQL.PostgreSQL.16 --accept-source-agreements --accept-package-agreements --override '--mode unattended --superpassword paybridge_dev_secret'
   ```
3. Start the service:
   ```powershell
   Start-Service postgresql-x64-16
   ```
4. Run migrations from the project root:
   ```powershell
   cd D:\Projects\PayBridgeAI\services\api
   go run cmd/migrate/main.go -action=up
   ```
5. Run the 20 SQL RLS Security Checks:
   ```powershell
   psql -U paybridge_admin -d paybridge_dev -f tests/sql/verify_rls.sql
   ```

### Option B: Docker Suite
If Docker Desktop is installed and running:
```bash
docker compose up -d postgres
cd services/api
go run cmd/migrate/main.go -action=up
```

---

## 3. Exact Resolution Steps for Blocker B2 (PayPal Sandbox WAN)
1. Log in to [developer.paypal.com](https://developer.paypal.com).
2. Create an App under "Sandbox".
3. Update `.env` in `D:\Projects\PayBridgeAI`:
   ```env
   PAYPAL_API_BASE=https://api-m.sandbox.paypal.com
   PAYPAL_CLIENT_ID=your_real_sandbox_client_id
   PAYPAL_CLIENT_SECRET=your_real_sandbox_client_secret
   PAYPAL_WEBHOOK_ID=your_real_webhook_id
   ```
4. Run the live checkout endpoint:
   ```bash
   cd services/api
   go run cmd/api/main.go
   ```
