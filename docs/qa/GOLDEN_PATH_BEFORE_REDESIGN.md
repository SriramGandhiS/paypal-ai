# Golden Path Demo — Pre-Redesign Inspection & Baseline Audit

**Generated:** 2026-11-11  
**Repository:** `D:\Projects\PayBridgeAI`  
**Dedicated Branch:** `feature/golden-path-redesign`  
**Baseline Git Commit Checkpoint:** `921d7916c6c80d8bb257b4dc4a6361d87f888433`  

---

## 1. Original Component File Paths

| Component / Layer | File Path | Description |
| :--- | :--- | :--- |
| **Vanilla Web Application (Active Local Server)** | `D:\Projects\PayBridgeAI\apps\web\index.html` | Contains `#workbench` container, `switchView('demo')`, and `renderDemoView()` which renders the Golden Path demo component. |
| **React Application Component** | `D:\Projects\PayBridgeAI\apps\web\src\screens\GoldenPathDemoScreen.tsx` | 7-step React interactive wizard component with JSON payloads and state vector pills. |
| **React App Parent Router** | `D:\Projects\PayBridgeAI\apps\web\src\App.tsx` | Routes `currentTab === 'demo'` to `<GoldenPathDemoScreen />`. |
| **API Client Service** | `D:\Projects\PayBridgeAI\apps\web\src\api\client.ts` | Dual-mode API client (live backend at `/api/v1` with fallback to deterministic local simulation). |
| **Go Backend Routes** | `D:\Projects\PayBridgeAI\services\api\internal\api\api.go` | Authoritative REST API (`/api/v1/cases`, `/checkout/paypal`, `/reconcile`, etc.). |
| **Go PayPal Integration** | `D:\Projects\PayBridgeAI\services\api\internal\paypal\client.go` | Server-side PayPal REST v2 Sandbox client with OAuth token acquisition and capture verification. |

---

## 2. Baseline Git State & Protection Status

- **Working Tree:** Clean (`nothing to commit, working tree clean`).
- **Baseline Commit:** `921d791` (`fix(conveyor): eliminate card text overlays and adjust typography scale`).
- **Dedicated Branch Created:** `feature/golden-path-redesign`.
- **Protected Branches Untouched:** `master` and `integration/vijay-agent-system` remain untouched.
- **Rollback Snapshot:** Full file backups and git references preserved.

---

## 3. Existing Component Behaviour & Visual Structure

### 3.1 Visual Structure (Before Redesign)
- **Container:** Contained within `#workbench` section under the interactive verification suite.
- **Navigation:** Used 3 numbered navigation sub-tabs:
  1. `1. Customer Payment Box`
  2. `2. AI Autonomous Recovery`
  3. `3. Resolution Roadmap`
- **Buttons & Stepping:** Used Next / Back buttons (`setDemoStage(1|2|3)`) to switch between isolated containers.
- **Introductory Section:** Used technical kicker `Golden Path. Cross-Border 504 Recovery Engine` with sequential telemetry subtitle.
- **Screenshot Evidence:** Saved to `docs/qa/golden_path_before_redesign.png` and artifact storage.

### 3.2 Simulation Behaviour
- **State Machine:** Governed by `demoState` variable with states: `'idle'`, `'loading'`, `'timeout'`, `'locked'`, `'reconciled'`.
- **Timeout Interception:** Simulates HTTP 504 Gateway drop after 1200ms delay.
- **Attempt Lock:** Engages `DO_NOT_REPEAT_CHECKOUT` mutex lock.
- **Reconciliation:** Polls PayPal provider REST query `/v2/checkout/orders/9K201824`, verifies capture, and unlocks safe refund path.

---

## 4. Existing API & Backend Dependencies

1. **Backend Liveness Endpoint:** `GET /healthz` (detects live Go backend vs. in-memory simulator).
2. **Cases Endpoint:** `POST /api/v1/cases` (creates payment case container).
3. **PayPal Checkout Order:** `POST /api/v1/cases/{id}/checkout/paypal` (initiates PayPal Sandbox v2 order).
4. **PayPal Capture:** `POST /api/v1/cases/{id}/checkout/paypal/capture` (authoritative capture).
5. **Reconciliation:** `POST /api/v1/cases/{id}/reconcile` (computes discrepancy against provider ledger).
6. **AI Copilot Endpoint:** `/api/v1/chat` or direct Groq/OpenAI client with financial governance constraints.

---

## 5. Authorized Files for Modification in Golden Path Redesign

In strict adherence to **Section 0: Absolute Change Restrictions**, only the Golden Path Demo component files are authorized for modification:

1. **`D:\Projects\PayBridgeAI\apps\web\index.html`**  
   *Scope restricted strictly to:*
   - `#workbench` Golden Path component markup (`renderDemoView` function and dedicated `.golden-*` scoped CSS).
   - Absolute zero changes to website background, gradients, wave graphics, header, navigation, conveyor cards, or footer.
2. **`D:\Projects\PayBridgeAI\apps\web\src\screens\GoldenPathDemoScreen.tsx`**  
   *Scope restricted strictly to:*
   - Golden Path demo UI implementation matching the approved reference design.
3. **`D:\Projects\PayBridgeAI\docs\qa\*`**  
   *Scope restricted strictly to:*
   - Audit documentation, changelogs, test reports, and rollback guides.

**Strictly Prohibited Files (Zero Changes):**
- Global typography or base styles (`base.css`, `tokens.css`, `index.css`).
- Navigation bar (`Navbar.tsx` or `.distill-nav`).
- Conveyor belt cards (`.k-card`, `.k-nav`, `.k-case-*`, etc.).
- Footer (`.footer`, Meadow Footer).
- Other screens (`DashboardScreen.tsx`, `PassportScreen.tsx`, `CaseDetailScreen.tsx`, etc.).
- Go backend database migrations or architectural schemas.

---

## 6. Checkpoint 1 Sign-Off

- [x] Codebase inspected and component boundaries strictly isolated.
- [x] Baseline git commit and dedicated branch `feature/golden-path-redesign` established.
- [x] Pre-redesign screenshot captured and archived.
- [x] Authorized file list locked down.
- [x] Ready to proceed to **Checkpoint 2 — UI Redesign**.
