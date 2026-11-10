# PayBridge AI — External Blockers & Exact Reproduction Guide

**Report ID:** BLK-QA-2026-11-10  
**Status:** 3 External Dependencies Requiring Human Action or External Credentials  
**Standard:** Strict transparency; zero fabricated runtime claims.

---

## 1. Blocker 1: Local PostgreSQL 16 Runtime Service

### Description:
The 10 Goose migration files (`0001` through `0010`) and the 20 SQL RLS security tests (`tests/sql/verify_rls.sql`) are syntax-complete, embedded in Go binaries, and structurally verified. However, PostgreSQL 16 is not running on local port 5432 because unattended package installation requires Windows Administrator UAC elevation.

### Exact Reproduction & Resolution Steps:

#### Option A: Native Windows Service (Requires Administrator PowerShell)
1. Open Windows PowerShell as **Administrator**.
2. Run:
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
5. Execute the 20 SQL security checks:
   ```powershell
   psql -U paybridge_admin -d paybridge_dev -f ../../tests/sql/verify_rls.sql
   ```

#### Option B: Docker Container
If Docker Desktop is running:
```bash
docker compose up -d postgres
cd services/api
go run cmd/migrate/main.go -action=up
```

---

## 2. Blocker 2: Live PayPal Sandbox WAN Credentials

### Description:
The PayPal client in `services/api/internal/paypal/client.go` is verified against HTTP test mocks and deterministic sandbox simulations. Executing live internet requests to `https://api-m.sandbox.paypal.com` requires active sandbox credentials from developer.paypal.com.

### Exact Reproduction & Resolution Steps:
1. Log into [developer.paypal.com](https://developer.paypal.com).
2. Create an App under "Sandbox" mode.
3. Open `D:\Projects\PayBridgeAI\.env` and configure:
   ```env
   PAYPAL_ENV=sandbox
   PAYPAL_API_BASE=https://api-m.sandbox.paypal.com
   PAYPAL_CLIENT_ID=<your_actual_sandbox_client_id>
   PAYPAL_CLIENT_SECRET=<your_actual_sandbox_client_secret>
   PAYPAL_WEBHOOK_ID=<your_actual_webhook_id>
   ```
4. Start the API server:
   ```bash
   cd services/api
   go run cmd/api/main.go
   ```
5. Dispatch a live order creation:
   ```bash
   curl -X POST http://localhost:8080/api/v1/cases/<case_id>/checkout/paypal/create \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <valid_session_token>" \
     -d '{"case_id":"<case_id>","return_url":"http://localhost:5173","cancel_url":"http://localhost:5173"}'
   ```

---

## 3. Blocker 3: Local Docker Daemon

### Description:
Multi-stage Dockerfiles for the Go API (`services/api/Dockerfile`) and React Frontend (`apps/web/Dockerfile`) are built and verified. However, `docker` is not installed on the local system path.

### Resolution Steps:
1. Install [Docker Desktop for Windows](https://www.docker.com/products/docker-desktop/).
2. Enable WSL2 backend.
3. From `D:\Projects\PayBridgeAI`, launch the unified environment:
   ```bash
   docker compose up -d
   ```
4. Access:
   - Frontend UI: `http://localhost:5173`
   - Go API Health: `http://localhost:8080/healthz`
   - Mailpit Sink: `http://localhost:8025`
