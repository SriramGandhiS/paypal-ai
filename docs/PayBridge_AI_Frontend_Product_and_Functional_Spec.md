# PayBridge AI --- Frontend Product & Functional Specification

**Repository:** `D:\Projects\PayBridgeAI`\
**Frontend:** Existing React + TypeScript + Vite application in
`apps/web`\
**Purpose:** Product definition, screen requirements, button/action
behaviour, API integration rules, design direction, testing and
Antigravity handover.

## 1. Purpose

This document defines what the PayBridge AI frontend must communicate,
which screens it needs, what each control does, which backend operations
it may call, and how it must be tested before acceptance.

PayBridge AI is a **B2B payment incident intelligence and safe recovery
coordination platform** for merchant/payment operations. It is not a
generic fintech dashboard.

Its central promise:

> **Know what happened. See the evidence. Choose the safe next step.**

The user must be able to answer: 1. What actually happened to this
payment? 2. What evidence supports that conclusion? 3. What is the next
safe action?

The frontend must make uncertainty visible. It must never treat a
timeout as proof of payment failure or encourage a blind retry.

## 2. Product definition and non-negotiable rules

### 2.1 Core problem

A payment provider may successfully capture a payment while the merchant
application fails to receive or persist the response. The merchant sees
a timeout or incomplete state and may incorrectly tell the customer
payment failed or attempt another charge.

PayBridge coordinates investigation using provider evidence,
merchant-side evidence, event history and reconciliation.

### 2.2 Core product concepts

-   **Payment case:** Investigation workspace for one payment incident.
-   **Payment attempt:** One attempt to initiate or complete a payment.
-   **Provider truth:** Provider's authoritative payment status and
    evidence.
-   **Evidence:** Timestamped records supporting or contradicting a
    conclusion.
-   **Payment Truth Passport:** Structured, verifiable summary of what
    is known and why.
-   **Truth Gap Graph:** Visual explanation of where the payment
    lifecycle lost agreement.
-   **Next Safe Action:** Policy-aware recommendation; not an automatic
    financial action.
-   **Two-Lens Explanation:** Merchant-facing operational explanation
    and investigator-facing evidence explanation.
-   **Recovery Mirror:** Comparison between evidence-led recovery and
    unsafe blind retry behaviour.
-   **Golden Path:** Repeatable demonstration of timeout, provider
    verification and reconciliation.

### 2.3 Product safety rules

1.  Provider evidence is authoritative for provider payment status.
2.  Timeout does not mean failure.
3.  Unknown is a valid state and must be displayed explicitly.
4.  Never show an unsafe retry action for an uncertain payment.
5.  AI may analyze, classify, explain and recommend. It must not create
    orders, capture, refund, retry or mutate financial state.
6.  Backend authorization is authoritative. Hiding a button is not
    security.
7.  Evidence should retain source, timestamp and verification status
    where available.
8.  Demo data must be visibly labelled **SIMULATED**.
9.  Never invent payment statuses, provider confirmations, financial
    metrics or audit records.
10. Do not expose credentials, tokens, card data or unnecessary personal
    data.

## 3. Users and roles

Support the roles already defined in the backend and canonical project
documents. Do not invent new privileges.

**Merchant operations:** Find cases, understand payment status, inspect
provider confirmation, read the operational explanation, view the
recommended safe next step and access a Passport when authorized.

**Investigator:** Inspect event chronology, evidence, merchant/provider
discrepancies, reconciliation and technical metadata when authorized.

**Administrator / operations lead:** Access operational overview,
queues, audit and configuration information only where actual backend
functionality and permissions exist.

**AI / investigator agents:** System actors, not privileged humans.
Their output is analysis/recommendation with provenance where available.
The UI must never imply they can execute financial mutations.

## 4. Information architecture

Keep the existing six-screen MVP and improve it; do not rebuild it
wholesale without evidence.

Primary navigation: 1. **Overview** --- operational summary and priority
cases. 2. **Cases** --- searchable incident inventory. 3. **Golden
Path** --- interactive, labelled demonstration. 4. **Truth Passport**
--- verified payment summary. 5. **Truth Gap Graph** --- lifecycle
evidence and gaps. 6. **Recovery Mirror** --- safe recovery versus blind
retry.

Inside a case, use contextual tabs only where implemented: - Summary -
Attempts - Timeline - Evidence - Reconciliation - Passport

### Global application shell

**Header:** PayBridge AI identity, environment indicator (Development /
Sandbox / Demo based on actual configuration), authenticated user/role
if available, and supported account controls.

**Navigation:** Active item highlighted, keyboard accessible,
collapsible on small screens.

**Connection/data-source indicator:** Clearly distinguish Live API,
PayPal Sandbox and Demo Simulator. If the app falls back to fixtures,
visibly say so. Never silently show fixture data as live data.

Do not add dead navigation items, fake notifications or unsupported
profile/settings controls.

## 5. Screen specifications and button behaviour

## 5.1 Overview / Dashboard

**Purpose:** Give operations users a fast view of investigation workload
and access to cases needing attention.

**Sections:** - Page title: "Payment Operations Overview". - Brief
product explanation. - Metrics only when backed by actual data: Open
Cases, Needs Review, Confirmed Provider Captures, Unknown / Awaiting
Verification. - Priority/recent case list. - Data-source/environment
indicator. - Create/open case action only if supported and authorized.

**Case list fields (only when supplied by API):** - Case ID. -
Payment/order reference. - Tenant-safe merchant label. - Current
state. - Last updated. - Attention reason. - Provider status and
verification label.

**Controls:** - **Open case:** Navigate to the selected case. - **Create
case:** Open validated form using fields required by existing API;
submit once and display server result. - **Search:** Search supported
case identifiers/fields; disclose local filtering if not server-side. -
**Status filter:** Use actual backend state values. - **Refresh:**
Re-fetch data; this is not a payment retry. - **Clear filters:** Reset
search and filters. - **Retry loading:** Retry the read request only.

**States:** Loading, empty, no search results, API failure, permission
denied, offline simulator. If offline fixtures are used, label them
simulated.

## 5.2 Cases / Case inventory

**Purpose:** Dedicated incident queue.

**Features:** Search, supported status filters, sorting by last update
where supported, pagination/incremental loading where supported,
attention reason and direct case navigation.

**Controls:** - **View details:** Open case. - **Refresh:** Reload
inventory. - **Filter / Reset:** Apply or clear supported filters. -
**New case:** Open creation flow only if API and role permit.

Do not create bulk capture, refund or retry controls.

## 5.3 Case Detail

**Purpose:** Main operational workspace for understanding one incident.

**Layout:** 1. Case header: case ID, order/payment reference, state,
last updated. 2. State summary: merchant observation, provider
observation and reconciliation state. 3. Next Safe Action panel. 4.
Payment attempt history. 5. Event timeline. 6. Evidence viewer. 7.
Reconciliation result. 8. Links to Passport and Truth Gap Graph.

**State labels** may include Confirmed by Provider, Pending at Provider,
Failed at Provider, Unknown / Verification Required, Conflicting
Evidence and Reconciliation Required. Map actual backend values
centrally; do not invent or reinterpret states.

**Next Safe Action panel:** Show recommended step, reason, supporting
evidence, missing preconditions, human approval requirement, whether it
is read-only/operational, and policy/permission result when exposed by
API.

**Controls:** - **Refresh case:** Re-fetch only. - **View attempt:**
Focus/open attempt details. - **View evidence:** Open evidence
details. - **Open Passport:** Navigate to case Passport. - **View Truth
Gap:** Navigate to graph. - **Run reconciliation:** Show only if an
authorized backend operation exists. Explain that it checks provider
truth; it must not capture/refund/retry. - **Copy reference:** Copy
case/order reference with feedback. - **Export:** Only if a real
authorized export capability exists.

Do not show generic "Retry Payment", "Capture Now" or "Refund Now"
controls in an uncertain-payment workflow.

**Evidence viewer:** Show evidence type, source, timestamp,
verification/signature status, related reference, payload summary and
integrity/hash status where available. Raw payload visibility must be
role-scoped and sensitive fields redacted.

## 5.4 Payment Truth Passport

**Purpose:** Evidence-backed explanation of the payment's known status.

**Sections:** - Passport/case/payment reference. - Overall truth
status. - Merchant-side observation. - Provider-side authoritative
observation. - Reconciliation conclusion. - Verification checks
supported by the data model. - Evidence list with sources and
timestamps. - SHA-256 seal/integrity status only if supplied or verified
by backend. - Generated/updated timestamp. - Missing/contradictory
evidence explanation.

**Two-Lens toggle:** - **Merchant Lens:** Plain-language operational
explanation and safe next step. - **Investigator Lens:** Detailed
evidence, event sequence, provider metadata and integrity information,
subject to permissions.

**Controls:** - **Merchant Lens / Investigator Lens:** Switch
explanation only; do not mutate payment data. - **Copy Passport
Reference:** Copy stable reference. - **Verify integrity:** Call a real
verification endpoint if available. Otherwise display returned
verification status; do not hash client-side and call it cryptographic
verification. - **Refresh Passport:** Re-fetch. - **Export Passport:**
Only when a real authorized export exists.

Distinguish "hash present" from "hash verified". Never show a green
verified seal without actual verification.

## 5.5 Truth Gap Graph

**Purpose:** Explain where the payment lifecycle lost agreement or
evidence.

Use actual data and project model. Conceptual participants may include
Customer, Merchant Application, PayBridge and Payment Provider.

Show nodes, transitions, timestamps, evidence references,
missing/delayed acknowledgements, contradictions and current
investigation focus.

**Controls:** - **Select node:** Show node details/evidence. - **Select
edge:** Show event, timestamp and source. - **Reset view:** Restore
graph viewport. - **Fit to screen:** Fit graph. - **Open evidence:**
Navigate to linked evidence. - **Back to case:** Return to case detail.

Provide an accessible textual event list; graph cannot be the only way
to understand the incident. Never invent events or evidence.

## 5.6 Recovery Mirror

**Purpose:** Explain why evidence-led recovery is safer than blind
retry.

**With PayBridge:** Detect uncertainty → query provider status → inspect
evidence → reconcile → recommend safe next step.

**Blind-retry failure pattern:** Merchant sees timeout → assumes failure
→ retries without provider check → duplicate-charge risk.

**Controls:** - **Select scenario:** Choose documented simulator
scenario. - **Run simulation:** Execute deterministic, clearly labelled
simulation. - **Reset simulation:** Reset simulator state. - **View
evidence:** Show evidence associated with scenario. - **Open Golden
Path:** Navigate to full demo.

All hypothetical results must show **SIMULATED**. Never present
hypothetical duplicate charges as real transactions.

## 5.7 Golden Path Demo

**Purpose:** Repeatable demonstration of the core PayBridge workflow.

Seven steps: 1. Create payment case. 2. Create checkout/order in
simulator or actual sandbox, clearly identified. 3. Simulate or receive
customer/provider approval. 4. Simulate merchant response timeout. 5.
Prevent blind retry. 6. Query authoritative provider status. 7.
Reconcile and generate Payment Truth Passport.

**UI:** Stepper, current/completed states, short step explanation, event
timeline, JSON telemetry inspector, current case/payment state, mode
indicator and reset/restart control.

**Controls:** - **Start demo:** Begin deterministic simulation. - **Next
step:** Execute next permitted step. - **Previous step:** Navigate
display history only; never reverse actual financial operations. -
**Reset demo:** Reset simulator only. Confirm if persisted demo records
would be deleted. - **Inspect JSON:** Expand/collapse telemetry. -
**Copy JSON:** Copy visible non-sensitive demo telemetry. - **Open
Passport:** Navigate after Passport exists.

Disable Next when prerequisites are missing. Never imply a real provider
capture unless it occurred in PayPal Sandbox.

## 6. Common interaction requirements

Every button must have a clear label, real handler, loading/disabled
state, success/error feedback, duplicate-submit protection, permission
handling and keyboard accessibility. No inert decorative controls.

Forms must label fields, explain required values, validate before
submission, show server errors, preserve values after recoverable
errors, prevent duplicate submit and provide Cancel/Close.

Status must never rely on colour alone. Use text labels for Confirmed,
Pending, Unknown, Needs Review, Failed and Simulated.

All data-backed screens must handle: - Initial loading and refresh. -
Empty result. - API/network error. - Permission denied. - Stale data /
last updated where relevant. - Offline simulator mode where supported.

Do not show "Payment successful" merely because a case was created or an
API request was accepted.

## 7. Visual design direction

The interface should feel like a premium, trustworthy payment-operations
intelligence workspace: precise, calm, technical and easy to scan.

Principles: - Strong information hierarchy. - High readability for dense
operational data. - Clear separation of payment status, evidence and
recommendation. - Avoid excessive nested cards; use whitespace and
dividers. - Restrained accent colours for active states and meaningful
statuses. - Consistent typography, spacing, radii and iconography. -
Charts/graphs must be meaningful, labelled and based on actual data. -
Avoid decorative financial charts and invented metrics. - Avoid flashy
effects that reduce readability or make the product look like a crypto
trading dashboard.

A refined dark operational interface or consistent premium light/dark
hybrid is acceptable. First inspect the current app and select one
coherent system; do not change the entire visual language without
review.

### Responsive requirements

Test 320, 375, 390, 768, 1024, 1440 and 1920px widths.

On mobile: - Sidebar becomes a usable drawer/compact navigation. -
Tables become readable stacked rows or clearly scrollable regions. -
Graph offers fit/reset controls and textual alternative. - Modals fit
the viewport. - No horizontal page overflow. - Primary actions remain
reachable. - Respect `prefers-reduced-motion`; essential information
cannot depend on animation.

## 8. API and data integration contract

Existing Go API routes, types and `packages/contracts/openapi.yaml` are
authoritative.

Before implementing a control: 1. Inspect OpenAPI. 2. Inspect actual Go
route and handler. 3. Inspect request/response types. 4. Inspect
authentication and role enforcement. 5. Inspect existing API client. 6.
Confirm the operation is truly supported.

Do not invent endpoints, payload fields or permissions.

Use the existing typed API client where possible. Centralize base URL
and auth handling. Handle non-2xx responses, timeouts and network errors
consistently. Never log secrets or sensitive payloads.

Do not automatically retry financial mutations. Read-only requests may
use bounded safe retries. Mutations must follow backend idempotency
contracts.

Map backend states centrally. Do not scatter ad hoc string comparisons
throughout components.

Clearly distinguish: - Live Go API data. - PayPal Sandbox evidence. -
Deterministic Demo Simulator fixtures.

If a backend capability does not exist, do not fake it. Mark it
unavailable and record the missing contract.

## 9. Accessibility, privacy and security

-   Use semantic HTML and keyboard-accessible controls.
-   Provide visible focus states and accessible names.
-   Label inputs and icon-only buttons.
-   Announce asynchronous updates where appropriate.
-   Maintain readable contrast.
-   Redact sensitive provider payload fields.
-   Do not expose tokens, secrets, full payment credentials or
    unnecessary personal data.
-   Do not store sensitive payment data in browser local storage.
-   Never bundle secrets into frontend assets.
-   Frontend role-based visibility is usability, not the security
    boundary; backend authorization remains mandatory.

## 10. Frontend test and acceptance criteria

### Functional

-   Dashboard shows real or explicitly simulated data.
-   Case search/filter works.
-   Case creation validates/submits only where supported.
-   Correct case detail loads.
-   Evidence source/timestamps display.
-   Safe Action explains recommendation and permissions.
-   Passport lenses switch correctly.
-   Passport verification is honest.
-   Graph node/edge selection works.
-   Recovery Mirror scenarios are labelled simulated.
-   Golden Path follows correct sequence.
-   Navigation works.
-   Every button performs its advertised action.
-   API errors are visible and recoverable.

### Security

-   Unauthorized users cannot access restricted operations.
-   Frontend cannot bypass backend policy.
-   AI actors have no financial mutation controls.
-   Cross-tenant data is not displayed.
-   Sensitive evidence is role-scoped/redacted.
-   No secrets appear in built assets.

### Build and regression

Run `npm run build`, available frontend tests, lint/type checks, and
existing Go tests after integration. Add browser end-to-end tests if a
suitable framework is configured.

Acceptance requires: - All screens have defined purpose. - All visible
controls work. - Data is real or explicitly labelled fixture. - Unknown
remains unknown. - No blind retry is encouraged. - AI cannot perform
financial mutations. - Permissions are respected. -
Loading/empty/error/offline states are covered. - Responsive layouts
pass target widths. - Production build passes. - No backend contract is
silently changed. - A screen-by-screen test report is saved.

## 11. Antigravity implementation handover

Work inside `D:\Projects\PayBridgeAI`.

### Step 1 --- Inspect

Read this document. Inspect `apps/web`, all six current screens, shared
components, API client, OpenAPI contract, actual Go routes/handlers,
auth and policy.

Produce a concise gap matrix before editing.

### Step 2 --- Preserve

Do not replace the app wholesale. Do not delete the six MVP screens. Do
not duplicate working components. Do not change backend contracts to
suit an invented UI. Do not modify Vijay/Diny-owned modules unless
integration is explicitly authorized.

### Step 3 --- Implement

Implement navigation, information hierarchy, controls and states screen
by screen. Prioritize functional correctness over decoration.

Build reusable components for: - Status indicators. -
Loading/error/empty states. - Evidence rows/viewer. - Case references. -
Confirmation dialogs. - Permission-aware actions. -
Environment/data-source indicator.

### Step 4 --- Wire every action

For every button document: - User intent. - Handler. - API endpoint or
local-only behaviour. - Permission requirement. - Loading state. -
Success state. - Failure state. - Test covering it.

If backend support is missing, do not fake it. Record the missing
contract.

### Step 5 --- Test and report

Create/update: - `docs/FRONTEND_PRODUCT_AND_FUNCTIONAL_SPEC.md` -
`docs/FRONTEND_IMPLEMENTATION_GAP_MATRIX.md` -
`docs/FRONTEND_BUTTON_ACTION_MAP.md` - `docs/FRONTEND_TEST_REPORT.md`

Include actual test results and evidence paths.

### Step 6 --- Integration safety

Work on a separate branch such as `feature/frontend-product-alignment`.
Preserve the baseline branch. Inspect `git diff` and `git status`. Never
commit `.env`, credentials, build secrets or unrelated files.

Do not claim completion while controls are placeholders or API
operations are mocked.

## 12. Final product principle

Build PayBridge AI as a trustworthy payment incident intelligence
product, not a generic admin dashboard.

The user must be able to open a case, understand the merchant/provider
discrepancy, inspect evidence, see why a safe next step is recommended,
and inspect a Payment Truth Passport.

**Show what is known. Show what is unknown. Show the evidence. Never
pretend a payment failed just because a request timed out.**
