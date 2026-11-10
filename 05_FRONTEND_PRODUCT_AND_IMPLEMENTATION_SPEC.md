# 05 — Frontend Product and Implementation Specification (React + TypeScript + Vite)

| Field | Value |
|---|---|
| Document | `05_FRONTEND_PRODUCT_AND_IMPLEMENTATION_SPEC.md` |
| Version | 1.0.0 |
| Date | 9 November 2026 |
| Status | **DRAFT — specification only. No frontend code exists.** |
| Audience | Frontend engineers, designers, coding agents |
| Depends on | 01 (vocabulary), 03 (SSE ADR-012), 04 (routes §14, error envelope §4, enums) |
| Related | 06 (visual language, per-screen UX — **this doc is the engineering contract, doc 06 is the design contract; both must be satisfied**) |

> **Change note (2026-11-09, Account 2 audit):** §2 lists `/app/lab/mirror` and `/app/demo` for several roles, but doc 04 route R37 (`POST /demo/scenarios/{scenario}`) is `demo_operator`-only. Resolution (doc 01 D-15): *viewing* the lab is open to the listed roles; *run/reset/fault* controls render only when the case/session payload carries `permissions.canRunScenario` (true for `demo_operator` memberships, which the seeded judge/demo account holds alongside `sender`). Doc 06 screens 11 and 16 specify this.

> **Rule of rules:** the UI never computes payment truth. It renders typed API data. A status the API marks `UNKNOWN` is **never** presented as a failure, never offers a retry-checkout, and never turns red.

## Table of contents
1. [Stack and project structure](#1-stack-and-project-structure) · 2. [Route map](#2-route-map) · 3. [Feature modules](#3-feature-modules) · 4. [Component hierarchy](#4-component-hierarchy) · 5. [API client and generated types](#5-api-client-and-generated-types) · 6. [Query, cache and state decisions](#6-query-cache-and-state-decisions) · 7. [Experiences by role](#7-experiences-by-role) · 8. [Route guards](#8-route-guards) · 9. [Status presentation layer](#9-status-presentation-layer) · 10. [Loading, error, empty, restricted, unknown states](#10-loading-error-empty-restricted-unknown-states) · 11. [Slow and offline handling](#11-slow-and-offline-handling) · 12. [Responsive behaviour](#12-responsive-behaviour) · 13. [Accessibility](#13-accessibility) · 14. [Design-system integration](#14-design-system-integration) · 15. [Real-time updates: SSE vs polling](#15-real-time-updates-sse-vs-polling) · 16. [Forms, validation, notifications](#16-forms-validation-notifications) · 17. [Timelines, graphs, evidence visualisation](#17-timelines-graphs-evidence-visualisation) · 18. [Performance and code splitting](#18-performance-and-code-splitting) · 19. [Real API screens vs labelled demo fixtures](#19-real-api-screens-vs-labelled-demo-fixtures) · 20. [Testing](#20-testing) · 21. [Acceptance criteria](#21-acceptance-criteria) · 22. [Implementation sequence](#22-implementation-sequence)

---

## 1. Stack and project structure

| Concern | Decision | Rationale |
|---|---|---|
| Framework | React (current stable, **VERIFY**) + TypeScript `strict` | SOURCE |
| Build | Vite | SOURCE; fast HMR; route-level code splitting |
| Routing | React Router (data router) | Nested layouts, guards via loaders |
| Server state | **TanStack Query** | SOURCE M15 §15.9; cache + invalidation by `version` |
| Client/UI state | React state + small context providers; `zustand` only if a cross-route UI store is genuinely needed (e.g. graph selection) | Avoid shadow payment state |
| API types | `openapi-typescript` (generated) + `openapi-fetch` client | OpenAPI is the contract (ADR-020) |
| Forms | `react-hook-form` + `zod` schemas (hand-aligned with OpenAPI; a codegen check compares enums) | Typed validation |
| UI primitives | Radix UI (unstyled, accessible) | Focus management, dialogs, tabs, tooltips |
| Styling | CSS custom-property tokens (doc 06) + CSS Modules; no runtime CSS-in-JS | Theming (light/dark), performance |
| Visualisation | Hand-built SVG (+ `d3-scale`/`d3-shape` utilities only) for timeline, Truth Gap Graph, mirror | Full control, a11y, no heavy graph lib |
| Icons | Single SVG icon set via inline components (tree-shaken) | Consistency |
| Dates/money | `Intl.DateTimeFormat`/`Intl.NumberFormat`; money stays a **decimal string**; no float math | SOURCE M15 §15.9 |
| Tests | Vitest + Testing Library; Playwright (e2e + a11y via axe); MSW for fixture mode | See §20 |
| Lint | ESLint (`jsx-a11y`), Prettier, `tsc --noEmit` | – |

### 1.1 Directory structure

```text
apps/web/
├── index.html
├── vite.config.ts   tsconfig.json   .eslintrc.cjs   playwright.config.ts
├── public/                      # favicon, og image (no provider logos)
└── src/
    ├── main.tsx   App.tsx   router.tsx
    ├── api/
    │   ├── schema.d.ts          # GENERATED from openapi.yaml — never hand-edited
    │   ├── client.ts            # openapi-fetch instance: credentials include, CSRF header, correlation id capture
    │   ├── errors.ts            # ApiError (code, message, correlationId, details)
    │   ├── keys.ts              # query key factory
    │   └── idempotency.ts       # crypto.randomUUID() per user intent (stable across retries of that intent)
    ├── app/                     # providers: QueryProvider, ThemeProvider, AuthProvider, ToastProvider, OfflineProvider
    ├── features/
    │   ├── landing/             # public intro (S1)
    │   ├── auth/                # login, onboarding, consent invite accept
    │   ├── cases/               # dashboard, new payment, case shell
    │   ├── payment/             # payment detail, Passport, safe-action panel
    │   ├── timeline/            # lifecycle timeline
    │   ├── truthgap/            # Truth Gap Graph workbench
    │   ├── lenses/              # sender lens, receiver lens
    │   ├── mirror/              # Counterfactual Recovery Mirror lab
    │   ├── evidence/            # evidence explorer + provenance drawer
    │   ├── notifications/
    │   ├── ops/                 # queue, diagnostics, webhooks
    │   ├── settings/            # privacy, consent, integration
    │   └── demo/                # walkthrough + scenario selector
    ├── components/              # design-system components (doc 06 §Component inventory)
    ├── hooks/                   # useCaseStream, useOnlineStatus, useReducedMotion, useRoleView …
    ├── lib/
    │   ├── status/              # statusPresentation() — the ONE mapper (§9)
    │   ├── format/              # money, time, ids (mask/copy)
    │   ├── a11y/                # live-region helper, focus utilities
    │   └── fixtures/            # demo fixtures (labelled; §19)
    ├── styles/                  # tokens.css, themes.css, base.css, motion.css
    └── test/                    # setup, MSW handlers, axe helpers
```

---

## 2. Route map

*(Prefix `/app/...` per decision D-06; M08's shorter sitemap is superseded.)*

| Path | Screen (doc 06 #) | Roles | Data source | Mode |
|---|---|---|---|---|
| `/` | 1 Landing | public | static + optional `GET /demo/scenarios` | static |
| `/how-it-works`, `/security`, `/developers` | 1 (sections) | public | static | static |
| `/login` | 2 Auth | public | R01–R04 | **real** |
| `/invite/:token` | 2 Onboarding/consent | receiver | R21–R22 | **real** |
| `/app` | redirect by role | any | R03 | real |
| `/app/cases` | 3 Workspace dashboard | sender/analyst/admin | R06 | **real** |
| `/app/cases/new` | 4 New payment | sender/analyst | R05, R08, R09 | **real** |
| `/app/cases/:id` | 5 Payment detail + Passport (shell) | participant | R07, R15, R26, R27 | **real** |
| `/app/cases/:id/timeline` | 6 Lifecycle timeline | participant | R11 | **real** |
| `/app/cases/:id/truth-gap` | 7 Truth Gap workbench | sender/analyst/admin (receiver: reduced) | R28, R29 | **real** |
| `/app/cases/:id/sender` | 8 Sender lens | sender/analyst | R15 (viewer=sender) | **real** |
| `/app/cases/:id/receiver` | 9 Receiver lens | receiver/analyst | R15 (viewer=receiver) | **real** |
| `/app/cases/:id/actions` | 10 Next Safe Action panel (also embedded in shell) | participant | R26, R25 | **real** |
| `/app/lab/mirror` | 11 Recovery Mirror | demo_operator/analyst/sender | R36–R37 | **SIMULATED (real API, simulated data)** |
| `/app/cases/:id/evidence` | 12 Evidence explorer | analyst/admin (+role subset) | R29, R30 | **real** |
| `/app/notifications` | 13 Notifications | any | R40–R41 | **real** |
| `/app/ops` | 14 Ops queue | analyst/admin | R32 | **real** |
| `/app/ops/health` | 14 Diagnostics | analyst/admin | R33, R47 | **real** |
| `/app/ops/webhooks` | 14 Webhook status | analyst/admin | R34 | **real** |
| `/app/settings/privacy` | 15 Privacy/consent | any | R42–R44, R23–R24 | **real** |
| `/app/settings/integration` | 15 Integration | admin | R45 | **real** |
| `/app/settings/team` | 15 (later) | admin | – | LATER |
| `/app/demo` | 16 Walkthrough + scenario selector | any (demo mode) | R36–R39 | real API driving SANDBOX or SIMULATED runs |
| `/share/:token` | share view (Passport only) | token | R19 | **real** |
| `/docs` | public docs | public | static | static |

Deep links inside a case use plain `#anchor` only (no hash-routing). Unknown routes → friendly 404 (also used for restricted case ids — **uniform** with API).

---

## 3. Feature modules

Each feature exports: `routes`, `queries` (hooks wrapping generated client), `components`, `copy` (all user-facing strings centralised for review against doc 06 copy rules), and `tests`. Cross-feature imports go through `components/` or `lib/` only (enforced with `eslint-plugin-boundaries`).

| Module | Key components | Key hooks |
|---|---|---|
| `cases` | `CaseTable`, `CaseCard`, `UncertaintyMeter` (non-completing), `NewPaymentWizard` | `useCases`, `useCreateCase` |
| `payment` | `CaseShell`, `PaymentStateTriad`, `PassportView`, `SafeActionPanel`, `DoNotRepeatBanner` | `useCase`, `usePassport`, `useSafeActions` |
| `timeline` | `LifecycleTimeline`, `TimelineEventRow`, `LaneStrip` | `useTimeline` |
| `truthgap` | `TruthGapCanvas`, `LaneColumn`, `ClaimNode`, `GapEdge`, `GapList`, `EvidenceDrawer` | `useTruthGap` |
| `lenses` | `SenderLens`, `ReceiverLens`, `RoleSwitcher` | `useRoleView` |
| `mirror` | `MirrorLab`, `ScenarioStage`, `ReplayScrubber`, `LedgerCounter`, `SimulatedBanner` | `useScenarioRun` |
| `evidence` | `EvidenceTable`, `ProvenanceDrawer`, `TierBadge`, `RedactedPayload` | `useEvidence` |
| `ops` | `OpsQueue`, `ThreeColumnInvestigator`, `WebhookInboxTable`, `ServiceHealthStrip` | `useOpsCases`, `useWebhooks` |
| `settings` | `ConsentCard`, `ShareLinkManager`, `IntegrationStatus`, `DataRequestForm` | `useConsents` |
| `demo` | `Walkthrough`, `ScenarioSelector`, `StepRail` | `useDemoScenarios` |

---

## 4. Component hierarchy

```text
<App>
 ├─ <QueryProvider/> <ThemeProvider/> <AuthProvider/> <OfflineProvider/> <ToastProvider/>
 └─ <RouterProvider>
     ├─ <PublicLayout>   Landing, HowItWorks, Security, Developers, Docs
     ├─ <AuthLayout>     Login, InviteAccept
     └─ <AppLayout>      (guard: authenticated)
         ├─ <TopBar>  (EnvironmentBadge, RoleSwitcher*, NotificationBell, UserMenu)
         ├─ <NavRail> (role-dependent)
         ├─ <OfflineBanner/> <Toasts/>
         └─ <Outlet>
             ├─ CasesDashboard
             ├─ NewPayment (Wizard → ReadinessStep → CheckoutStep)
             ├─ <CaseShell id>                    ← loads Case, subscribes SSE
             │    ├─ <CaseHeader>  (amount, EnvironmentBadge, version, last verified)
             │    ├─ <DoNotRepeatBanner/>         ← from safe-actions, not from local logic
             │    ├─ <PaymentStateTriad/>         ← provider | application | receiver lanes
             │    └─ <Outlet>  Passport | Timeline | TruthGap | Lens | Evidence | Actions
             ├─ OpsQueue → OpsCaseInvestigator
             ├─ MirrorLab
             └─ Settings…
```
`EnvironmentBadge` renders on **every** screen showing a case (`REAL SANDBOX` / `SIMULATED` / `PAYPAL SIMULATOR`), never omitted (ADR-008).

---

## 5. API client and generated types

- `npm run gen:api` → `openapi-typescript ../../packages/contracts/openapi.yaml -o src/api/schema.d.ts`. CI fails if the committed file is stale.
- `client.ts`: base URL `/api/v1`; `credentials: 'include'`; adds `X-CSRF-Token` (from R04, cached) on unsafe methods; reads `X-Correlation-Id` from every response and attaches it to errors/logs; **no tokens in `localStorage`**.
- `ApiError` mirrors the error envelope; `isUnknownOutcome(res)` is true for `202` + `outcome:"UNKNOWN"` — such responses are **successful** from an HTTP view and must **not** reach the error boundary.
- **Idempotency:** the client generates one `Idempotency-Key` per *user intent* (e.g. one click on "Confirm payment" → one key, reused on automatic network retry of that same intent, replaced only by a new deliberate user action). The key is stored in component state/ref for the duration of the intent; checkout buttons are disabled while an intent is in flight **and** while the case reports an unresolved attempt.
- Enum drift guard: a unit test imports generated enums and asserts `statusPresentation()` covers every enum member exhaustively (TypeScript `never` check + runtime test).

---

## 6. Query, cache and state decisions

**DESIGN DECISION:** TanStack Query is the **only** store for server-derived data.

| Data | Query key | Stale time | Invalidation |
|---|---|---|---|
| Case | `['case', id]` | 0 (always revalidate on focus/SSE) | SSE `case.version`; mutation success |
| Passport | `['passport', id, viewer]` | 30 s | SSE `passport.version` |
| Timeline | `['timeline', id]` | 0 | SSE `timeline.append` (append, not refetch, when possible) |
| Safe actions | `['safe', id]` | 0 | any case/passport change |
| Truth gap | `['truthgap', id]` | 15 s | SSE `case.version` |
| Evidence list | `['evidence', id, filters]` | 60 s | case version bump |
| Cases list | `['cases', filters]` | 10 s | window focus; SSE (ops) |
| Ops health | `['ops','health']` | 15 s polling | interval |

**Versioned invalidation:** every Case/Passport payload carries `version`; the stream/poller compares with cache and invalidates only when greater. **Never optimistic-update payment state** — mutations on checkout/capture are *pessimistic*; the UI shows "Confirming…" until the server answers (`200`, or `202 UNKNOWN`).

`gcTime` modest; no persistence of financial payloads to disk/localStorage (privacy). On logout: `queryClient.clear()`.

---

## 7. Experiences by role

### 7.1 Sender (`/app/cases`, `/app/cases/:id`)
Create case → readiness → approve → "Confirming payment" → result. In the hero state: persistent **DoNotRepeatBanner**, state triad showing PayPal lane `CAPTURE COMPLETED (verified)`, merchant lane `CONFIRMATION LOST`, receiver lane `NOT OBSERVABLE`; sender lens text per doc 04 §18.3. The checkout button is **absent** (not merely disabled) when `safe-actions.blocked` contains `CREATE_NEW_ORDER`, with the blocked reason shown in words.

### 7.2 Receiver (`/app/cases/:id/receiver`, via `/invite/:token`)
Consent card first (exact fields to be shared; separate notification toggle) → receiver lens: *PayPal capture* and *Bank settlement* shown as **separate** items; official-account guidance; "ask for help" and "revoke participation". No sender private fields; masked ids only.

### 7.3 Operations (`/app/ops/*`)
Queue sorted by uncertainty/verification age/event-delivery health; three-column investigator (`Evidence | Causal Timeline | Diagnosis & Actions`); **Compare provider vs app** toggle (first-class); webhook inbox table with verified vs simulated filter; approve internal reconciliation (R31) behind a confirmation dialog stating plainly: *"This updates your own merchant record. It does not contact PayPal or move money."*

### 7.4 Developer / Demo
Scenario selector (R36–R39) with a persistent **mode chip** per scenario (`REAL SANDBOX`, `SIMULATED`, `REAL SANDBOX + FAULT`), run history, replay scrubber, ledger counters (`orders: 1 · captures: 1`).

---

## 8. Route guards

**UX guard only** — authorisation is always server-side (doc 04 §8). Guards live in React Router loaders:

1. `requireAuth`: calls R03; on `401` → `/login?next=…`.
2. `requireRole(roles)`: if role not permitted → render the **same** friendly "not found" as a missing resource (uniform denial; no enumeration).
3. `requireConsentForReceiver`: receiver case routes check consent state from R07 `viewer.consent_active`; inactive → consent screen or restricted state.
4. `demoMode`: demo/lab routes appear only if `GET /demo/scenarios` succeeds (server-controlled).
5. A revoked consent mid-session produces `404` on next fetch and a closed SSE stream → UI shows the restricted state (§10) and clears case caches.

Sensitive affordances (capture button, approve-internal-reconciliation) are rendered from **server-provided permission flags** on the case payload (`permissions: {canCapture, canReconcileInternal, …}`), never inferred from role strings alone.

---

## 9. Status presentation layer

`lib/status/statusPresentation.ts` is the **single** place that turns enums into label + icon + tone + copy. Everything renders through it. It is exhaustively typed so adding an enum value breaks the build until mapped.

### 9.1 Tone vocabulary (colour is never the only channel — icon + text always)

| Tone | Used for | Never used for |
|---|---|---|
| `verified` (green) | Provider-verified facts (`CAPTURE_COMPLETED` by E1/E2) | Unverified "sent" stages |
| `neutral` | `CREATED`, `APPROVED`, informational | – |
| `pending` (blue-grey) | `CAPTURE_PENDING`, `WAITING_FOR_WEBHOOK`, `RECONCILING` | – |
| `unknown` (amber) | `UNKNOWN`, `API_TIMEOUT`, `EVENT_GAP`, `NOT_OBSERVABLE`, `CAPTURE_UNKNOWN`, `ORDER_UNKNOWN`, `REVIEW_REQUIRED` | **Never red, never "failed"** |
| `attention` (amber-deep) | `INCONSISTENT`, `APPLICATION_MISMATCH`, `NEEDS_HUMAN_REVIEW` | – |
| `error` (red) | **Confirmed** provider denial (`CAPTURE_DENIED` via E1/E2) or confirmed security risk | any timeout/unknown |
| `simulated` (violet, hatched) | Anything with environment ≠ `REAL_SANDBOX` | Real events |

### 9.2 Mandatory copy mapping (excerpt — full in doc 06 §Copy)

| Enum | Label | Helper text |
|---|---|---|
| `provider_state=UNKNOWN` | "Outcome not yet confirmed" | "We're checking the original PayPal record. This does not mean PayPal declined your payment." |
| `application_state=API_TIMEOUT` | "Our confirmation was lost" | "We didn't hear back from our server. Your payment may have gone through." |
| `provider_state=CAPTURE_COMPLETED` | "PayPal capture confirmed" | "Verified from the original PayPal record at {time}." |
| `receiver_state=NOT_OBSERVABLE` | "Bank settlement not verified" | "We can't see this stage with the data available here." |
| `provider_state=CAPTURE_DENIED` | "PayPal declined this payment" | "Confirmed by PayPal at {time}." (only with E1/E2 evidence) |
| `application_state=INCONSISTENT` | "Records disagree" | "PayPal and our merchant record differ. We're reconciling — no new payment is needed." |

**Forbidden strings (lint rule `no-forbidden-copy` + test):** "Payment failed" when the underlying state is any `*UNKNOWN*`/`API_TIMEOUT`/`EVENT_GAP`; "Try again" next to an unresolved attempt; "Transferred to bank"; "100%"/progress-complete indicators for stages beyond last verified; "PayPal is down" without `PUBLIC_SOURCE` evidence attached.

---

## 10. Loading, error, empty, restricted, unknown states

*Every data surface implements all six; they are visually and semantically distinct (SOURCE M15 §15.9).*

| State | Meaning | Presentation | Never |
|---|---|---|---|
| **Loading** | Request in flight | Skeletons matching final layout; `aria-busy` | Spinners blocking the whole screen |
| **Empty** | No data exists | Illustration-free explanatory text + next action | Implying an error |
| **Error (known)** | API error envelope | Message + `correlation_id` (copyable) + retry for **reads** | Retry for checkout mutations |
| **Unknown outcome** | `202 UNKNOWN` / unknown state | Amber "checking" with live timestamp of last check; banner | Red; "failed"; retry-checkout |
| **Restricted** | Role lacks field/section (consent) | Inline lock panel: *"This part isn't shared with you."* + what the viewer *can* do | Blank/missing without explanation |
| **Offline/slow** | Connectivity | See §11 | Converting stale to failure |

Error boundaries per feature route; a boundary fallback never offers "retry payment". Mutation errors on checkout show: *"We couldn't confirm what happened. Please don't repeat this payment — we're checking."* when the error class is network/timeout (the client cannot know), and the exact server message for explicit `4xx` validation errors that are **provably pre-processing**.

---

## 11. Slow and offline handling

- `useOnlineStatus` (navigator + failed-request heuristic) drives `OfflineBanner`: *"You're offline. What you see may be out of date. Nothing has been changed."*
- While offline: reads show cached data with a **stale badge** and "last updated" time; **all mutating controls are disabled** with an explanation; SSE shows "Live updates paused" and auto-resumes.
- **Slow request (>3 s):** progressive disclosure — "Still checking…" ; (>10 s on checkout mutation) switch to *outcome-unknown* messaging and begin polling case state — **never** tell the user to retry.
- **Tab restore / refresh during PayPal approval:** on return, the case shell re-fetches; the attempt state decides the screen (e.g. `APPROVED` → "Confirm payment"; `CAPTURE_UNKNOWN` → unknown banner). The client holds no authoritative state.
- Service-worker/PWA caching of API responses is **not** used in MVP (privacy + staleness).

---

## 12. Responsive behaviour

Breakpoints (SOURCE M15 §15.6): **1440**, **1024**, **390** px; fluid in between (`clamp()` typography and spacing from tokens).

| Region | ≥1280 | 768–1279 | <768 |
|---|---|---|---|
| Nav | Persistent rail | Collapsible rail | Bottom tab bar + sheet |
| Case shell | Header + triad; right **Evidence/Safe-action panel** pinned | Right panel becomes tabbed drawer | Single column; sticky action bar with safe-action headline |
| Truth Gap | 5-lane canvas | 5 lanes with horizontal scroll inside own container | Lane **stack** accordion with gap summary on top |
| Timeline | Vertical rail + detail pane | Rail + modal detail | Compact rail; detail as full-screen sheet |
| Tables | Dense | Medium | Card list with key fields; horizontal scroll never on `<body>` |

Touch targets ≥ 44 px where possible; no hover-only information; wide content scrolls in its **own** `overflow-x:auto` container.

---

## 13. Accessibility

**Target: WCAG 2.2 AA** (SOURCE M08 §8.5). Detail in doc 06 §Accessibility; engineering requirements:

- Semantic landmarks (`header/nav/main/aside`), one `h1` per route, logical heading order.
- All status chips = **icon + text + (optional) colour**; `aria-label` includes tone meaning ("Unknown — outcome not yet confirmed").
- **Live regions:** a single polite `role="status"` region announces version changes ("Passport updated: provider capture confirmed"); `role="alert"` only for confirmed errors/security.
- **Keyboard:** every interactive element reachable; visible 2px focus ring using `--focus`; roving tabindex for graph nodes and timeline items (arrow keys), `Esc` closes drawers/dialogs, focus trap + return focus in Radix dialogs; skip-to-content link.
- **Graph/timeline accessibility:** each has an equivalent **data table / list view** toggle ("View as table"); nodes expose `aria-describedby` evidence summaries.
- **Reduced motion:** all non-essential animation behind `@media (prefers-reduced-motion: no-preference)`; reduced mode shows end-states instantly.
- Locale-aware date/currency; UTC shown in evidence/audit details.
- Forms: associated labels, error summaries with focus management, `aria-invalid`/`aria-describedby`.
- Automated: `eslint-plugin-jsx-a11y`; axe in Vitest component tests and Playwright for all 16 screens; **manual** keyboard + screen-reader pass (NVDA/VoiceOver) is required evidence (AO-13).

---

## 14. Design-system integration

Doc 06 defines tokens and components. Engineering contract:

- `styles/tokens.css` defines CSS custom properties exactly as named in doc 06 §Tokens (`--surface-0…`, `--text-1…`, `--verified`, `--unknown`, `--attention`, `--error`, `--simulated`, `--focus`, spacing/radius/shadow/motion scales). Dark theme via `@media (prefers-color-scheme: dark)` **and** `[data-theme="dark"]` override with an in-app toggle (persisted in `localStorage` — theme preference only; no financial data).
- No hard-coded colours in components; lint rule `no-raw-color`.
- Components in `src/components` match the doc 06 inventory names 1:1 and ship a Storybook-free **gallery route** `/app/dev/gallery` (dev builds only) showing every component in every state for review and visual regression.
- Fonts: system stack + one variable display/text font self-hosted (licence **VERIFY**); `font-display: swap`; tabular numerals for money.

---

## 15. Real-time updates: SSE vs polling

**Decision: SSE primary, adaptive polling fallback** (ADR-012).

| Option | Pros | Cons | Verdict |
|---|---|---|---|
| **SSE** (`EventSource`) | One-way server→client matches our need; plain HTTP; cookie auth works; built-in reconnect + `Last-Event-ID`; proxies/Render support it | No custom headers on native `EventSource`; some corporate proxies buffer | **Primary** |
| WebSocket | Bi-directional | We have no client→server realtime need; extra auth/ops complexity | Rejected |
| Polling | Simplest; robust | Latency/load | **Fallback** + ops health |

Implementation:
- `useCaseStream(caseId)` opens `GET /cases/{id}/stream` (R27). Same-origin cookie auth (no tokens in URL).
- Events: `case.version {version}`, `passport.version {version, viewer}`, `timeline.append {event}`, `heartbeat`. Handler compares versions and invalidates queries; on `timeline.append` it appends to the cache.
- Reconnect: browser auto-reconnect with `Last-Event-ID`; if the stream errors 3× or is blocked, fall back to polling `['case', id]` every 5 s (10 s when tab hidden; paused when offline) and show a subtle "Live updates paused — refreshing periodically" note.
- **Revocation/consent change:** server closes the stream; the next fetch `404`s → restricted state.
- Accessibility: updates announced through the single live region, throttled (≤ 1 announcement / 5 s, coalesced).

---

## 16. Forms, validation, notifications

### 16.1 Forms
| Form | Library | Rules |
|---|---|---|
| Login | RHF + zod | Generic failure copy; no credential hints |
| New payment wizard (5 steps: participants, invoice, route, amount/currency, consent) | RHF + zod | Amount string regex `^\d{1,14}(\.\d{1,6})?$` + currency exponent check; currency USD (MVP); country ISO-2; free-text memo flagged as untrusted (shown escaped, never rendered as HTML/markdown with links auto-activated) |
| Receiver invite | RHF + zod | Email validation; field allowlist checkboxes with plain descriptions |
| Consent accept | – | Shows exact fields; notification consent separate and unchecked by default |
| Share link | – | Role + expiry select; copy link; revoke |
| Privacy export/deletion | RHF | Confirmation + explanation of audit-retention limits |

Client validation mirrors server rules but the **server is authoritative**; server `VALIDATION_FAILED.details.fields[]` maps onto fields. The **amount sent for checkout comes from the case**, never the browser; the wizard's amount only creates the case.

### 16.2 Notifications
- In-app toasts for low-risk confirmations; **persistent banners** (not toasts) for anything affecting payment safety (Do-not-repeat, restricted, stale).
- Notification centre (`/app/notifications`) lists server notifications (R40) with read state; empty/restricted/error states per §10.
- Notification copy is server-generated from validated Passport claims; the client never composes payment-status text from raw enums except through `statusPresentation()`.

---

## 17. Timelines, graphs, evidence visualisation

*(Visual design in doc 06; engineering contracts here.)*

### 17.1 Lifecycle timeline (`LifecycleTimeline`)
- Input: `TimelineEvent[]` (`id, kind, label, lane, tier, trust, environment, occurred_at, observed_at, evidence_id, verified`).
- Renders **lanes** (Browser/Client, Merchant App, PayPal, Webhook, Receiver) on a shared time axis with **observed vs occurred** marks; events with `environment ≠ REAL_SANDBOX` get the simulated pattern.
- A **"last verified stage" marker** sits on the axis; everything after it is drawn as *unobserved* (dashed), **never** as progress toward completion.
- Interactions: hover/focus → tooltip (tier, trust, times); click → `EvidenceDrawer`; keyboard: ←/→ move between events, `Enter` opens evidence, `T` toggles table view.
- Virtualise beyond ~200 events (simple windowing).

### 17.2 Truth Gap Graph (`TruthGapCanvas`)
- Data: R28 → lanes, nodes (`claim_code, value, tier, trust, observed_at, evidence_ids`), edges (`relation, rule_id, origin`), gaps.
- Layout: **deterministic columnar layout** (5 lanes × time order) computed in TS — no force simulation (stable, accessible, testable).
- Encoding: lane columns; node = claim chip with tier badge; `CONTRADICTS` = zig-zag amber connector with a labelled gap pill; `CORROBORATES` = solid thin; `INDEPENDENT_OF` (bank stage) = dotted; `LLM_HYPOTHESIS` edges = dashed with "AI hypothesis" label.
- Interactions: select node → highlight its edges, open provenance; select gap → focus both claims; "Compare provider vs app" toggle collapses to two columns; table view toggle.
- Plain text of every gap comes from the API (`gaps[].plain_text`), not generated client-side.

### 17.3 Evidence explorer (`EvidenceTable` + `ProvenanceDrawer`)
- Filters: tier (E1–E6), kind, trust, environment, verified.
- Provenance drawer: id (copy), kind, tier badge with tooltip definition, trust, times (UTC + local), `payload_sha256` (short + copy), **redacted payload** (JSON viewer with collapsed sensitive-key markers), signature status for webhooks, **"used by"** list (claims/passport versions citing it), and its edges.
- Never renders unredacted payloads; fields absent from the API response are not requested.

### 17.4 Passport (`PassportView`)
Three sections — **Confirmed**, **Unknown**, **Next safe action** — with evidence chips, observed times, privacy scope, share expiry, version selector (diff between versions via R15/R16), and "Why is this unknown?" help. **No** "100 %"/progress-complete bar.

---

## 18. Performance and code splitting

| Technique | Detail |
|---|---|
| Route-level splitting | `React.lazy` per feature route; landing is separate chunk with no app code |
| Heavy parts lazy | `TruthGapCanvas`, `MirrorLab`, JSON viewer load on demand |
| Budgets (targets, **not measured**) | Landing JS ≤ ~100 KB gz; app shell ≤ ~180 KB gz; each feature chunk ≤ ~80 KB gz — enforce with `size-limit` in CI |
| Rendering | Memoise table rows; virtualise long lists/timelines; avoid layout thrash in SVG (transform-only animations) |
| Network | HTTP caching for static assets; TanStack Query dedupe; no polling when tab hidden (except SSE) |
| Fonts/images | Subset fonts; `font-display: swap`; SVG illustrations; preload critical font |
| Metrics | Web Vitals captured in dev/staging to console + `/ops/health` optional beacon (no PII) |

---

## 19. Real API screens vs labelled demo fixtures

**Policy (owner prompt):** no disconnected mock UI presented as complete.

| Category | Rule |
|---|---|
| **Real API screen** | Bound to the generated client; works against the backend; appears in the milestone report as "real" |
| **Demo fixture** | Static data from `lib/fixtures/` rendered **only** (a) in the public landing's illustrative animation, (b) in Storybook-free gallery, (c) in tests (MSW). Every fixture-driven view shows a persistent **`DEMO FIXTURE`** ribbon and its test IDs are prefixed `fixture-` |
| **Simulated run** | Real API + `LOCAL_SIMULATION` data (Mirror Lab). Shows `SIMULATED` banner and never shares ids/namespaces with real cases |

**Enforcement:** a CI script greps for imports of `lib/fixtures` outside the allowed folders; a Playwright check asserts any route marked `real` in §2 issues network calls to `/api/v1`; the milestone report must list each screen with `real | fixture | simulated`.

**Phase rule:** a screen is "done" only when it is **real** (except Landing and the gallery). A fixture-only screen is at most **STUB** status.

---

## 20. Testing

*(All **NOT RUN** — plan.)*

| Level | Tooling | What |
|---|---|---|
| Unit | Vitest | `statusPresentation` exhaustive (every enum), money/date formatters, idempotency key lifecycle, forbidden-copy lint test |
| Component | Testing Library + axe | All states (loading/empty/error/unknown/restricted/offline) for each major component; keyboard navigation of timeline and graph; focus management |
| Integration | Vitest + MSW | Hooks against OpenAPI-validated mocks: `202 UNKNOWN` path never renders error/retry; SSE → cache invalidation; fallback-to-polling |
| Contract | Schema checks | MSW handlers are generated/validated from `openapi.yaml` so mocks cannot drift |
| E2E | Playwright | Golden path on the real stack (Compose): login → case → (fault) unknown → banner → reconcile → Passport v2; receiver consent + revocation (restricted state); duplicate-click prevention |
| Visual regression | Playwright screenshots | 16 screens × 3 breakpoints × light/dark (baseline reviewed by humans) |
| Accessibility | axe + manual | Zero serious/critical axe violations; manual screen-reader/keyboard checklist stored in `docs/evidence/a11y-report.md` |
| Safety | Dedicated tests | "Unknown never red", "no retry for unresolved attempt", "no forbidden strings", "simulated ribbon present" |

---

## 21. Acceptance criteria

| ID | Criterion | Verified by |
|---|---|---|
| FE-01 | All 16 screens exist, bound to the real API except Landing/gallery | Playwright network assertions |
| FE-02 | Types generated from OpenAPI; CI fails on drift | CI |
| FE-03 | `statusPresentation()` is exhaustive; unknown/timeout states never render error tone or failure copy | Unit + component tests |
| FE-04 | No checkout control is shown when `CREATE_NEW_ORDER` is blocked; reason shown | Component/E2E |
| FE-05 | SSE updates a Passport version bump without reload; polling fallback works | E2E |
| FE-06 | Receiver view never contains sender-private fields (DOM + network) | Privacy E2E (permission diff) |
| FE-07 | Every case screen shows an `EnvironmentBadge`; simulated data always bannered | E2E |
| FE-08 | Keyboard-only completion of golden path; visible focus; reduced motion honoured | Manual + Playwright |
| FE-09 | WCAG 2.2 AA target: no serious/critical axe issues; manual pass recorded | AO-13 evidence |
| FE-10 | Layouts verified at 1440/1024/390 | Visual regression |
| FE-11 | No tokens/secrets/financial payloads in `localStorage`/`sessionStorage` | E2E storage inspection |
| FE-12 | Bundle budgets enforced | CI `size-limit` |

---

## 22. Implementation sequence

*(Maps to doc 02 phases P4–P6.)*
1. **P4.1** Vite scaffold, tokens/themes, base components, `api/` client + codegen, auth + layouts + guards, gallery route.
2. **P4.2** `statusPresentation` + forbidden-copy lint; Cases dashboard; New payment wizard (create case, readiness, order, approval hand-off, capture).
3. **P4.3** Case shell, state triad, Do-not-repeat banner, Passport view, safe-action panel, SSE + polling fallback.
4. **P4.4** Lifecycle timeline; sender/receiver lenses; invite/consent/revocation flows; notifications.
5. **P6.1** Truth Gap workbench; evidence explorer + provenance drawer.
6. **P6.2** Ops queue + three-column investigator + webhook status + diagnostics.
7. **P6.3** Mirror lab + demo walkthrough/scenario selector; settings/privacy/integration.
8. **P6.4** Landing page (S1), polish, responsive and a11y passes, visual-regression baselines.

*End of document 05.*
