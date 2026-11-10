# PayBridge AI — API Contract & Specification Compatibility Report

**Comparing:** Vijay Kasthuri Submission vs Canonical Sriram Platform Baseline  
**OpenAPI Specification:** `packages/contracts/openapi.yaml`  
**Go Module Alignment:** `github.com/paybridge/api` -> `github.com/paybridge/paybridge-ai/services/api`  
**Evaluator:** Antigravity  
**Audit Branch:** `integration/vijay-agent-system`

---

## 1. Specification Compatibility Matrix

| Architectural Interface | Sriram Canonical Contract | Vijay Submitted Contract | Compatibility Status | Adjustment Made |
| :--- | :--- | :--- | :--- | :--- |
| **Go Module Path** | `github.com/paybridge/paybridge-ai/services/api` | `github.com/paybridge/api` | **BREAKING (Namespace)** | Re-mapped all imports to canonical module path during integration. |
| **Go Version in `go.mod`** | `go 1.22` | `go 1.27.3` (Non-existent version) | **BREAKING (Build Blocker)** | Corrected `go.mod` version to `go 1.22` with tidy dependencies. |
| **OpenAPI 3.1 Contract** | 9,805 bytes (Core MVP routes) | 19,554 bytes (Expanded agent routes) | **COMPATIBLE (Superset)** | Vijay added routes for `/v1/cases/{id}/investigate`, `/v1/cases/{id}/timeline`, `/v1/cases/{id}/safe-actions`. |
| **Database Migrations** | 0001–0010 (Canonical 40 tables, RLS) | 0001–0010 (10 partial tables) | **CONFLICT (Schema Drift)** | Preserved Sriram's canonical 40-table migrations. Vijay's agent tables are already represented. |
| **HTTP Router** | Chi router with tenant context | Chi router with tenant context | **COMPATIBLE** | Handlers plug into Chi mux cleanly. |
| **Authentication & RBAC** | Bearer tokens + `X-Tenant-ID` | Bearer tokens + `X-Tenant-ID` | **COMPATIBLE** | Both enforce role-based access control. |
| **AI Role Mutation Guard** | `ai_agent` denied mutation | `ai_agent` outputs claims only | **COMPATIBLE** | Both strictly prevent AI from issuing payment captures or refunds. |
| **PayPal Adapter Interface** | `Client` with `CaptureOrder` & `GetOrder` | `Reader` (GET-only) & `Mutator` (POST) | **COMPATIBLE (Dual Interface)** | Both interfaces coexist cleanly in `internal/paypal`. |
| **State Vector Dimensions** | `provider_state`, `application_state` | 3D: `provider`, `application`, `receiver` | **ENHANCEMENT** | Integrated 3D Reducer enhances state machine fidelity. |
| **Evidence Tiers** | Tier 1 through 6 | Tier E1 through E6 | **100% IDENTICAL** | Exact semantic equivalence across all 6 tiers. |
| **Frontend Component Tree** | React 18 + Vite + Tailwind | React 18 + Vite + CSS Tokens | **COMPATIBLE** | Integrated smoothly; `npm run build` bundles all 16 screens cleanly. |

---

## 2. Detailed Contract Analysis

### A. Go Module Path Mismatch
Vijay built his modules with the package prefix `github.com/paybridge/api`. In our monorepo, the canonical module path defined in `services/api/go.mod` is `github.com/paybridge/paybridge-ai/services/api`.
- **Resolution:** All incoming packages (`internal/agents`, `internal/claims`, `internal/safeaction`, `internal/reducer`, `internal/projector`, `internal/orchestrator`, `cmd/worker`, `cmd/demob`) were programmatically updated to import from `github.com/paybridge/paybridge-ai/services/api/...`.

### B. The `go 1.27.3` Toolchain Glitch
Vijay's `services/api/go.mod` declared `go 1.27.3`. Because Go 1.27 has not been released, standard Go toolchains (1.22, 1.23, 1.26) attempt to download a non-existent toolchain binary from `proxy.golang.org`, stalling local compilation.
- **Resolution:** The `go.mod` version was set to `go 1.22`, matching standard LTS compiler capabilities.

### C. Database Migration Preservation
Vijay's submission included 10 small migration files (`0001_extensions_roles.sql` through `0010_rls.sql`) that created only a subset of tables. Overwriting Sriram's migrations with Vijay's would have destroyed 30 canonical enterprise tables.
- **Resolution:** Sriram's canonical 40-table migrations were preserved as the single source of database truth.
