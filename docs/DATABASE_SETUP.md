# PayBridge AI — PostgreSQL 16 Database Setup Runbook

## 1. Overview
PayBridge AI requires PostgreSQL 16 with Row-Level Security (`FORCE ROW LEVEL SECURITY`) and cryptographic extensions (`pgcrypto`, `uuid-ossp`, `citext`).

Database connection string format:
```text
postgres://paybridge_admin:paybridge_dev_secret@localhost:5432/paybridge_dev?sslmode=disable
```

---

## 2. Setup Options

### Option A: Docker Desktop (Recommended)
If Docker Desktop is installed, start the isolated dev rig from the repository root:
```powershell
docker compose -f infra/docker/docker-compose.yml up -d
```
Verify container health:
```powershell
docker compose -f infra/docker/docker-compose.yml ps
```

### Option B: Native Windows PostgreSQL 16 Installation
If running natively without Docker:
1. Download and install PostgreSQL 16 for Windows from [EnterpriseDB](https://www.enterprisedb.com/downloads/postgres-postgresql-downloads).
2. Set administrator password during setup.
3. Open `psql` as superuser and create the development database and administrative user:
```sql
CREATE USER paybridge_admin WITH PASSWORD 'paybridge_dev_secret' NOBYPASSRLS;
CREATE DATABASE paybridge_dev OWNER paybridge_admin;
GRANT ALL PRIVILEGES ON DATABASE paybridge_dev TO paybridge_admin;
```

---

## 3. Running Database Migrations (0001–0010)

### Method 1: Using the Built-in Go Migration Runner
The project includes a standalone migration CLI that embeds all forward migrations. It can be run from the repository root or from `services/api`:

**From Repository Root (`D:\Projects\PayBridgeAI`):**
```powershell
# Applies all forward migrations 0001-0010
go run ./services/api/cmd/migrate -action=up

# Check migration status
go run ./services/api/cmd/migrate -action=status

# With custom database URL (optional):
go run ./services/api/cmd/migrate -db-url="postgres://paybridge_admin:paybridge_dev_secret@localhost:5432/paybridge_dev?sslmode=disable" -action=up
```

**From `services/api` Directory (`D:\Projects\PayBridgeAI\services\api`):**
```powershell
# Applies all forward migrations 0001-0010
go run ./cmd/migrate -action=up

# Check migration status
go run ./cmd/migrate -action=status
```

### Method 2: Using Goose CLI
If `goose` is installed on PATH:
```powershell
goose -dir services/api/db/migrations postgres "postgres://paybridge_admin:paybridge_dev_secret@localhost:5432/paybridge_dev?sslmode=disable" up
```

---

## 4. Verification Check
Run the automated verification suite to confirm that exactly 40 application tables, row-level security policies, and hardened stored procedures are active:
```powershell
psql -U paybridge_admin -d paybridge_dev -f tests/sql/verify_rls.sql
```
Expected output:
```text
NOTICE:  SUCCESS: Exactly 40 application tables verified.
NOTICE:  SUCCESS: FORCE ROW LEVEL SECURITY verified on core tables.
NOTICE:  SUCCESS: Direct SELECT revoked from paybridge_app.
NOTICE:  SUCCESS: All SECURITY DEFINER functions have hardened search_path.
test_banner: === ALL S3 DATABASE SECURITY VERIFICATIONS PASSED ===
```
