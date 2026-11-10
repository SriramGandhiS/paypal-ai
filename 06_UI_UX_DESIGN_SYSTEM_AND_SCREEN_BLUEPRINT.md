# 06 — UI/UX Design System and Screen Blueprint

| Field | Value |
|---|---|
| Document | `06_UI_UX_DESIGN_SYSTEM_AND_SCREEN_BLUEPRINT.md` |
| Version | 1.0.0 |
| Date | 9 November 2026 |
| Status | **DRAFT — design specification only. No screen, token file or component described here has been built, rendered or user-tested.** |
| Audience | Frontend coding agent, designers, reviewers, hackathon judges reading the design rationale |
| Depends on | 01 (vocabulary §12, principles §9, decisions D-16), 03 (components/states), 04 (route IDs R01–R48, safe actions §17, projections §18), 05 (route map §2, `statusPresentation` §9, SSE §15, fixtures policy §19) |
| Related | 07 (what the AI panels may show), 08 (privacy matrix, consent), 09 (visual-regression and a11y tests), 10 (demo storyboard) |
| Labels | SOURCE · DESIGN DECISION · ASSUMPTION · VERIFY · SIMULATED · SANDBOX |

> **Contract split.** Doc 05 is the *engineering* contract (routes, state, data flow). **This document is the *design* contract** (look, feel, composition, copy, interaction detail). Both must be satisfied. Where this document says "Screen N", N is exactly the numbering of doc 05 §2 and of the owner's 16-screen list. **No new routes are introduced here.**

> **Rule of rules (restated from 01/05):** unknown is never red and never "failed"; colour is never the only signal; a simulation is never mistaken for a real event; the UI never computes payment truth.

## Table of contents
1. [Design intent and traceability](#1-design-intent-and-traceability) · 2. [Brand narrative and art direction](#2-brand-narrative-and-art-direction) · 3. [Design tokens](#3-design-tokens) · 4. [Typography](#4-typography) · 5. [Layout, grid, density, elevation](#5-layout-grid-density-elevation) · 6. [Motion](#6-motion) · 7. [Iconography and illustration](#7-iconography-and-illustration) · 8. [Visual grammar for evidence products](#8-visual-grammar-for-evidence-products) · 9. [Component state matrix](#9-component-state-matrix) · 10. [Accessibility](#10-accessibility) · 11. [Component inventory](#11-component-inventory-and-api-expectations) · 12. [Copy system](#12-copy-system) · 13. [Screens 1–16](#13-screen-blueprints) · 14. [Visual QA checklist](#14-visual-qa-checklist) · 15. [Acceptance and open items](#15-acceptance-and-open-items)

---

## 1. Design intent and traceability

**The design problem (DESIGN DECISION).** A payment investigation UI has an unusual job: it must make *uncertainty legible and calm* while making *proof one click away*. Most fintech UIs optimise for the feeling of completion (green ticks, progress bars, big balances). PayBridge optimises for the opposite virtue: **an honest boundary between what is confirmed and what is not.** The design therefore needs (a) a visual language where "unknown" is a dignified, first-class place; (b) a way of showing *two people's views* of one payment without leaking either; (c) a way of showing *evidence strength* that does not rely on colour or a fake confidence percentage.

### 1.1 Screen ↔ route ↔ doc 05 mapping (no new routes)

| # | Screen | Route(s) from doc 05 §2 | Mode | Primary API route IDs (doc 04 §14) |
|---|---|---|---|---|
| 1 | Landing / product introduction | `/`, `/how-it-works`, `/security`, `/developers`, `/docs` | static (+ labelled illustration) | R36 (optional) |
| 2 | Authentication and onboarding | `/login`, `/invite/:token` | real | R01–R04, R21–R22 |
| 3 | Workspace / case dashboard | `/app/cases` | real | R06, R03, R40, R47 |
| 4 | New payment / connect sandbox / create investigation | `/app/cases/new` | real | R05, R08, R09, R10, R45 (admin only) |
| 5 | Payment detail + Payment Truth Passport | `/app/cases/:id` | real | R07, R15, R16, R17, R26, R27 |
| 6 | Payment lifecycle timeline | `/app/cases/:id/timeline` | real | R11, R27 |
| 7 | Truth Gap Graph workbench | `/app/cases/:id/truth-gap` | real | R28, R29, R30 |
| 8 | Sender Two-Lens | `/app/cases/:id/sender` | real | R15 (viewer=sender), R23 |
| 9 | Receiver Two-Lens | `/app/cases/:id/receiver` | real | R15 (viewer=receiver), R20–R24 |
| 10 | Next Safe Action panel | `/app/cases/:id/actions` (+ embedded in shell) | real | R26, R25, R12, R31 |
| 11 | Counterfactual Recovery Mirror lab | `/app/lab/mirror` | **SIMULATED** (real API, simulated data) | R36–R39 |
| 12 | Evidence explorer + provenance | `/app/cases/:id/evidence` | real | R29, R30 |
| 13 | Notification center | `/app/notifications` | real | R40, R41 |
| 14 | Operations / sandbox diagnostics / webhook status | `/app/ops`, `/app/ops/health`, `/app/ops/webhooks` | real | R32–R34, R47, R12, R31 |
| 15 | Settings / privacy / consent / integrations | `/app/settings/privacy`, `/app/settings/integration` | real | R42–R46, R17–R18, R23–R24 |
| 16 | Demo walkthrough / scenario selector | `/app/demo` | real API driving SANDBOX or SIMULATED runs | R36–R39 |

---

## 2. Brand narrative and art direction

### 2.1 Visual concept — "Evidence Ledger" (D-16)

**Narrative (SOURCE tagline, DESIGN expansion).** *One payment. Two people. One verified story.* The brand metaphor is a **ledger kept by a careful clerk who will only write down what they have seen** — and who leaves a clearly ruled blank where they haven't. Every surface is a page of that ledger: warm paper, ink-dark text, hairline rules, marginal notes (evidence chips), and a small seal beside every fact telling you how it was witnessed.

The second metaphor is **two lenses**. The sender and the receiver look at the same payment through different glass. The logo, the Passport header and the lens switcher all reuse one device: **two overlapping circles whose intersection (the vesica) is the part both may see.** That intersection shape recurs as the "shared facts" marker.

| Idea | How it shows up visually |
|---|---|
| Careful clerk, honest blanks | Ruled rows; unobserved stretches drawn as *dashed rules*, never as an empty progress bar |
| Witnessed facts | **Evidence seals** (E1–E6) — shape-coded, not colour-coded (§8.2) |
| Two lenses | Overlapping-circle motif; lens switcher; "frosted silhouette" of fields a viewer may not see (§8.8) |
| Calm under uncertainty | Warm neutrals, amber (not red) for unknown, generous whitespace, slow settle motion |
| Not PayPal | Neutral integration wording only; no PayPal logo or colours; no fake provider UI (brand-rule risk, doc 08) |

### 2.2 Design principles (ordered; earlier wins conflicts)

1. **Last confirmed step first.** Every case screen answers "what do we know for sure, and as of when?" before anything else.
2. **Unknown is a place, not an error.** It gets its own icon, tone, copy, and layout position. It is never red, never "failed", never a spinner that implies imminent completion.
3. **Proof on demand.** Any statement is ≤ 1 interaction from its evidence (seal → provenance drawer).
4. **Two lenses, one ledger.** Role views are built from different projections; the UI makes the *difference* visible and trustworthy, not hidden.
5. **Calm density.** Information-rich, but with strong hierarchy: one dominant sentence, one dominant state, one dominant action per screen.
6. **Honest decoration.** Ornament exists only if it encodes meaning (a rule, a seal, a lane). No stock people, no decorative gradients, no glow, no fake charts.
7. **Accessible by construction.** Contrast, focus, keyboard, reduced motion, table alternatives are part of every component's definition of done.

### 2.3 Anti-patterns (blocking in design review)

Generic card grid + left sidebar as the whole product · KPI tiles with invented numbers · rainbow status colours · green for anything unverified · progress bars/rings that imply completion · skeuomorphic bank/credit-card art · glassmorphism, neon glows or mesh gradients · "AI sparkle" iconography implying authority · confidence percentages ("97 % sure") · dark-pattern urgency ("Pay now!") · tooltips as the only carrier of essential meaning.

### 2.4 Imagery and illustration direction

Line diagrams drawn with the same 1.5 px stroke as icons: lanes, lenses, rules, seals. Two-tone (ink + one tone colour). **No photography of people or money.** Landing illustrations are fixtures and carry a persistent `ILLUSTRATION — NOT LIVE DATA` caption (doc 05 §19). Real screenshots are used only after implementation (M20 §8.3).

---

## 3. Design tokens

*Token names are the contract with doc 05 §14 (`styles/tokens.css`). Contrast ratios below were **computed** from the proposed hex values using the WCAG 2.x relative-luminance formula; they must be re-verified in the browser after implementation (**VERIFY**, AO-13).*

### 3.1 Colour — light theme ("Paper")

| Token | Hex | Role | Notes |
|---|---|---|---|
| `--surface-0` | `#F5F3EE` | App canvas | Warm paper |
| `--surface-1` | `#FBFAF7` | Raised page/section | |
| `--surface-2` | `#FFFFFF` | Cards, inputs, popovers | |
| `--surface-sunken` | `#EDEAE2` | Wells, table headers, code blocks | |
| `--text-1` | `#101A2E` | Primary text (ink) | 15.66:1 on surface-0 · 17.37:1 on surface-2 |
| `--text-2` | `#434E63` | Secondary text | 7.55:1 / 8.37:1 |
| `--text-3` | `#5B6578` | Tertiary/captions (still body-legible) | 5.29:1 / 5.87:1 |
| `--border` | `#D8D3C5` | Decorative hairlines only | 1.35:1 — **never** carries meaning alone |
| `--border-strong` | `#8E8878` | Input borders, meaningful separators | 3.53:1 on surface-2 (≥ 3:1 non-text) |
| `--accent` | `#1F4FC4` | Links, primary buttons, selection | 6.37:1 / 7.06:1; white-on-accent 7.06:1 |
| `--accent-hover` | `#183E9B` | Hover/pressed | |
| `--focus` | `#1F4FC4` + 2 px white offset ring | Focus ring | Visible on all surfaces incl. dark (see §10) |

### 3.2 Colour — semantic status tones (light)

| Tone (doc 05 §9.1) | Fg token | Fg hex | Bg token | Bg hex | Contrast fg/bg | Icon | Used for |
|---|---|---|---|---|---|---|---|
| `verified` | `--verified` | `#06683F` | `--verified-bg` | `#E3F2EA` | 5.93 | seal-check | Provider-verified facts only (E1/E2) |
| `neutral` | `--text-2` | `#434E63` | `--surface-sunken` | `#EDEAE2` | ≥ 5 | dot | `CREATED`, `APPROVED`, info |
| `pending` | `--pending` | `#345A7A` | `--pending-bg` | `#E6EDF3` | 6.15 | clock-arc | `CAPTURE_PENDING`, `WAITING_FOR_WEBHOOK`, `RECONCILING` |
| `unknown` | `--unknown` | `#8A5200` | `--unknown-bg` | `#FBEFD9` | 5.61 | open-ring | `UNKNOWN`, `API_TIMEOUT`, `EVENT_GAP`, `NOT_OBSERVABLE`, `*_UNKNOWN`, `REVIEW_REQUIRED` |
| `attention` | `--attention` | `#97400A` | `--attention-bg` | `#FBE9DC` | 5.83 | offset-arrows | `INCONSISTENT`, `APPLICATION_MISMATCH`, `NEEDS_HUMAN_REVIEW` |
| `error` | `--error` | `#B42318` | `--error-bg` | `#FCE8E6` | 5.58 | octagon-x | **Confirmed** `CAPTURE_DENIED` (E1/E2) or confirmed security risk **only** |
| `simulated` | `--simulated` | `#5B34B0` | `--simulated-bg` | `#ECE6F8` | 6.73 | flask | Any `environment ≠ REAL_SANDBOX` |

> **Hue-separation rule (DESIGN DECISION).** `unknown` (amber-brown) and `attention` (burnt orange) are deliberately far from `error` (red) and from each other in lightness *and* carry different icons and words. A colour-blind reader can still tell them apart because **icon + label are mandatory** on every status (WCAG 1.4.1).

### 3.3 Colour — dark theme ("Night ledger")

| Token | Hex | Contrast notes (computed) |
|---|---|---|
| `--surface-0` / `-1` / `-2` / `-sunken` | `#0D1320` / `#131B2B` / `#1A2438` / `#090E18` | |
| `--text-1` / `-2` / `-3` | `#EEF1F7` / `#BCC5D6` / `#97A3B8` | 16.41 / 10.69 / 7.29 on surface-0 |
| `--border` / `--border-strong` | `#2A3650` / `#6C7A96` | strong: 3.59 on surface-2 |
| `--accent` / on-accent | `#7FA6FF` / `#0A1430` | 7.78 on surface-0 · 7.62 on-accent |
| `--verified` / bg | `#5DD19A` / `#10291E` | 8.14 |
| `--pending` / bg | `#8FB7DA` / `#15263A` | 7.26 |
| `--unknown` / bg | `#F2B955` / `#33260D` | 8.32 |
| `--attention` / bg | `#F59A62` / `#35200F` | 7.09 |
| `--error` / bg | `#FF8A7E` / `#3A1511` | 7.09 |
| `--simulated` / bg | `#B79CFF` / `#251A45` | 7.02 |

### 3.4 Theme strategy and CSS variable skeleton

**DESIGN DECISION.** Light "Paper" is the default and the theme used for the demo video (best print/screenshot fidelity). Dark is fully specified, honours `prefers-color-scheme`, and has an in-app toggle that persists **theme only** (doc 05 §14). `color-scheme: light dark` is declared so native form controls adapt. A `forced-colors: active` block maps tones to system colours and relies on icon + text.

```css
/* styles/tokens.css — PROPOSED; names are the contract with doc 05 §14 */
:root {
  color-scheme: light dark;
  --surface-0:#F5F3EE; --surface-1:#FBFAF7; --surface-2:#FFFFFF; --surface-sunken:#EDEAE2;
  --text-1:#101A2E; --text-2:#434E63; --text-3:#5B6578;
  --border:#D8D3C5; --border-strong:#8E8878;
  --accent:#1F4FC4; --accent-hover:#183E9B; --on-accent:#FFFFFF; --focus:#1F4FC4;
  --verified:#06683F; --verified-bg:#E3F2EA;
  --pending:#345A7A;  --pending-bg:#E6EDF3;
  --unknown:#8A5200;  --unknown-bg:#FBEFD9;
  --attention:#97400A; --attention-bg:#FBE9DC;
  --error:#B42318;    --error-bg:#FCE8E6;
  --simulated:#5B34B0; --simulated-bg:#ECE6F8;
  /* spacing (4 px base) */
  --space-0:0; --space-1:2px; --space-2:4px; --space-3:8px; --space-4:12px; --space-5:16px;
  --space-6:24px; --space-7:32px; --space-8:48px; --space-9:64px; --space-10:96px;
  /* radii */ --radius-chip:4px; --radius-input:8px; --radius-card:12px; --radius-sheet:20px; --radius-full:999px;
  /* elevation */ --shadow-1:0 1px 2px rgb(16 26 46 / .06); --shadow-2:0 6px 20px rgb(16 26 46 / .10); --shadow-3:0 16px 48px rgb(16 26 46 / .16);
  /* motion */ --dur-instant:80ms; --dur-fast:160ms; --dur-base:240ms; --dur-slow:400ms; --ease-settle:cubic-bezier(.2,.7,.2,1); --ease-standard:cubic-bezier(.4,0,.2,1);
}
@media (prefers-color-scheme: dark){ :root:not([data-theme="light"]){ /* dark values from §3.3 */ } }
:root[data-theme="dark"]{ /* dark values from §3.3 */ }
@media (prefers-reduced-motion: reduce){ :root{ --dur-fast:0ms; --dur-base:0ms; --dur-slow:0ms; } }
```

### 3.5 Simulated pattern (non-colour channel)

`--simulated-hatch: repeating-linear-gradient(135deg, transparent 0 6px, color-mix(in srgb, var(--simulated) 22%, transparent) 6px 7px)`. Applied as background on every simulated card/lane/row, plus the `SIMULATED` ribbon (§8.3). In forced-colors mode the hatch is dropped and the ribbon text + flask icon carry the signal.

---

## 4. Typography

| Role | Family (recommendation) | Fallback stack | Use |
|---|---|---|---|
| **Display / narrative** | *Newsreader* (variable serif, optical sizes) — **VERIFY** licence (SIL OFL expected) and self-host | `Georgia, 'Times New Roman', serif` | Landing headlines, Passport "verdict line", lens headlines. **Never** for data or UI chrome |
| **UI / body** | *Inter* (variable) — **VERIFY** licence, self-host | `system-ui, -apple-system, 'Segoe UI', Roboto, sans-serif` | All UI, tables, forms |
| **Mono / identifiers** | *IBM Plex Mono* — **VERIFY** licence, self-host | `ui-monospace, SFMono-Regular, Menlo, Consolas, monospace` | Order/capture ids, hashes, JSON, timestamps |

`font-display: swap`; subset Latin; preload the UI and display variable files only (doc 05 §18 budgets).

### 4.1 Scale (fluid with `clamp()`; px at 1440)

| Token | Size / line-height | Weight | Use |
|---|---|---|---|
| `display-xl` | 64/68 | Newsreader 500 | Landing hero |
| `display-l` | 44/50 | Newsreader 500 | Section heroes, Passport verdict line (desktop) |
| `h1` | 32/38 | Inter 650 | Page titles |
| `h2` | 24/30 | Inter 600 | Section heads |
| `h3` | 18/26 | Inter 600 | Card titles |
| `body-l` | 18/28 | Inter 400 | Lead paragraphs, lens narrative |
| `body` | 16/24 | Inter 400 | Default |
| `body-s` | 14/20 | Inter 400/500 | Dense tables, chips |
| `caption` | 12/16 | Inter 500, `letter-spacing:.02em` | Timestamps, seal labels (never below 12 px) |
| `mono` | 13/20 | Plex Mono 450 | Ids, hashes |
| `money-xl` | 40/44 | Inter 600 tabular | Case amount |

### 4.2 Numerals, money, time, ids

* Always `font-variant-numeric: tabular-nums slashed-zero`. Money is **currency code + amount** (`USD 2,000.00`), formatted from the decimal **string** via `Intl.NumberFormat` without float parsing (doc 05 §1). Never symbol-only; never an INR/settlement figure unless a provider-evidenced field exists (M20 §12.7).
* Time: local time with a visible UTC tooltip; evidence/audit views show UTC first. Relative time ("14 min ago") always paired with an absolute timestamp in `title` and screen-reader text.
* Ids: masked display `ORD-…7F2A` (last 4 emphasised), copy button with "Copied" live-region announcement; full value only where the role projection allows (doc 04 §18.1).

---

## 5. Layout, grid, density, elevation

### 5.1 Breakpoints and grid

| Name | Width | Columns / gutter | Page margin | Notes |
|---|---|---|---|---|
| `xl` | ≥ 1280 (design at 1440) | 12 / 24 | 32 (max content 1360) | Persistent nav rail + pinned evidence dock |
| `lg` | 1024–1279 | 12 / 20 | 24 | Collapsible rail; dock becomes tabbed drawer |
| `md` | 768–1023 | 8 / 16 | 20 | |
| `sm` | 390–767 (design at 390) | 4 / 16 | 16 | Bottom tab bar; sticky safe-action bar |

Wide content (tables, graph) scrolls inside its own `overflow-x:auto` container; the page body **never** scrolls sideways (doc 05 §12).

### 5.2 App shell zones (screens 3–16)

```
┌────────────────────────────────────────────────────────────────────────────┐
│ TopBar 56: wordmark · EnvironmentBadge · (RoleSwitcher*) · 🔔 · user menu     │
├──────┬─────────────────────────────────────────────────────┬───────────────┤
│ Nav  │ MAIN STAGE (case content; max 880–960 reading width    │ EVIDENCE DOCK │
│ rail │ for narrative, full width for graph/timeline)          │ 360 px pinned │
│ 72/  │                                                        │ (≥1280)       │
│ 232  │                                                        │               │
├──────┴─────────────────────────────────────────────────────┴───────────────┤
│ (sm only) Sticky Safe-Action bar + bottom tab bar                           │
└────────────────────────────────────────────────────────────────────────────┘
```

### 5.3 Density modes (DESIGN DECISION)

`comfortable` (default for sender/receiver: row height 56, body 16) and `compact` (ops/evidence tables: row height 36, body-s 14). Density is a per-view setting, not a theme; it never changes colour or copy.

### 5.4 Borders, radii, elevation

Flat-first: surfaces separate with **hairline rules and tonal steps**, not shadows. `--shadow-1` only on cards that are interactive; `--shadow-2` popovers/menus; `--shadow-3` modals/drawers. Radii per tokens. Card = `surface-2` + `--border` + `--radius-card`; *ledger row* = no radius, bottom rule only.

---

## 6. Motion

**Philosophy.** Motion communicates *state change in evidence*, never excitement. Nothing bounces; nothing loops indefinitely except a quiet "checking" pulse that is replaced by static text under reduced motion.

| Interaction | Behaviour (default) | Duration/easing | Reduced-motion alternative |
|---|---|---|---|
| **Verify-settle** (a fact becomes verified) | Pending/unknown chip cross-fades to verified; seal "stamps" (scale 1.06→1, 1 shadow pulse) once | 240 ms `--ease-settle` | Instant swap; live-region announcement still fires |
| **Lane tick** (state triad updates) | The last-verified marker slides along the rail to its new position | 400 ms | Instant reposition |
| **Banner enter** (Do-not-repeat) | Height expand + fade; no shake, no bounce | 240 ms | Instant |
| **Passport version bump** | Changed lines get a 1.5 s left-rule highlight (`--accent`), then fade | 160/1500 ms | Static "Updated" tag for 10 s |
| **Graph select** | Selected node outlines; connected edges thicken; others dim to 55 % | 160 ms | Same, no transition |
| **Checking pulse** | Open-ring icon's dotted arc rotates slowly | 2.4 s linear loop | Static icon + "Last checked 14:02:11" text |
| **Drawer/dialog** | Slide 16 px + fade | 240 ms | Instant |
| **Mirror replay** | Scrubber moves events in order; "naive" vs "PayBridge" ledgers tick | step-driven | Step buttons only; no autoplay |

Rules: no parallax, no scroll-jacking, no autoplay video, only `transform`/`opacity` animations (doc 05 §18). `prefers-reduced-motion: reduce` removes all non-essential animation; essential feedback becomes text.

---

## 7. Iconography and illustration

20 px grid, 1.5 px stroke, rounded joins, single inline-SVG set (tree-shaken). Each status has **one** icon that is reused everywhere:

| Meaning | Icon | Shape logic |
|---|---|---|
| Verified | `seal-check` | Check inside scalloped seal |
| Pending / reconciling | `clock-arc` | Ring with moving arc |
| Unknown / not observable | `open-ring` | Ring with a *gap* and dotted arc (the gap = "not seen") |
| Attention / records disagree | `offset-arrows` | Two arrows pointing past each other |
| Error (confirmed) | `octagon-x` | Octagon with X |
| Simulated | `flask` | Lab flask |
| Restricted | `lock-lens` | Lock overlapping a lens circle |
| Do-not-repeat | `shield-pause` | Shield with pause bars |
| Evidence | `seal` (tier variants §8.2) | |
| Compare lenses | `two-lenses` | Overlapping circles, vesica filled |

Empty-state illustrations are drawn from the same lens/lane/rule vocabulary; they contain **no** people, coins, bank buildings, or provider marks.

---

## 8. Visual grammar for evidence products

### 8.1 Status chip anatomy
`[icon] LABEL · helper (optional)` — min height 24 px; icon 16 px; label `body-s` 600; `aria-label` = tone meaning + label + helper ("Unknown — Outcome not yet confirmed"). Chips are never colour-only and never abbreviations of enums (always the §12 label).

### 8.2 Evidence seals (tier badges) — shape-coded

| Tier | Seal | Pattern | Tooltip/definition (also in drawer) |
|---|---|---|---|
| **E1** | Double-ring seal, solid | ◉ | Authorised provider GET, verified |
| **E2** | Single-ring seal, solid | ● | Signed provider webhook, verified |
| **E3** | Rounded-square seal | ■ | Merchant-controlled application record |
| **E4** | Open ring | ○ | Official documentation / public status |
| **E5** | Dotted ring | ◌ | Consented user testimony |
| **E6** | Hatched ring (+ `SIMULATED` or `AI` suffix) | ⌀ | Simulation or model output — never proof |

A seal always carries its text (`E1`) and, on hover/focus, its plain-language definition. **No percentage, no star rating, no "confidence" number anywhere** (M20 §14.5, doc 07).

### 8.3 Environment and trust markers
* **EnvironmentBadge** (every case screen, ADR-008): `REAL SANDBOX` (ink outline, seal-check-less neutral), `PAYPAL SIMULATOR` and `SIMULATED` (violet fill + hatch + flask). Sits in TopBar *and* CaseHeader.
* **SimulatedRibbon**: diagonal-cornered band reading `SIMULATED — not a real PayPal event` on any simulated surface; ids rendered `SIM-…` in violet mono.
* **Verified-vs-simulated filter chips** wherever events are listed (screens 12, 14).

### 8.4 State Triad ("three rails")
Three horizontal rails stacked: **PayPal (provider_state)**, **Merchant app (application_state)**, **Receiver (receiver_state)**. Each rail: label, current-state chip, last-verified marker (a seal on the rail at the verified time), solid rule *up to* the marker and **dashed rule beyond it** (unobserved). If `provider_state=UNKNOWN`, the marker sits at `last_confirmed_provider_state` and the dashed rule carries an `open-ring` at the right end captioned "Not yet confirmed". The receiver rail for `NOT_OBSERVABLE` is a dotted rule with `eye-slash` and "Bank settlement not verified". **No rail ever shows a filled bar or percentage.**

### 8.5 Lifecycle timeline — "ledger tape"
Vertical tape, newest at top (default; toggle to oldest-first). Left gutter = time ticks (UTC + local); centre = five lane columns' dots joined by a spine only where events are consecutive; right = detail. Each event row: seal (tier) · label · lane glyph · "occurred" vs "observed" times (two small ticks; the gap between them is drawn if > 1 s). A **"last verified stage" rule** (solid ink, labelled) crosses the tape; events after it are dashed outlines labelled "not yet verified".

### 8.6 Truth Gap Graph encoding (extends doc 05 §17.2)
Five lane columns (Browser/Client · Merchant App · PayPal · Webhook · Receiver) on a deterministic time-ordered grid. **Claim chips** = status chip + seal + time. Edge grammar (shape, not just colour): `CONTRADICTS` = amber **zig-zag "tear"** connector with a labelled *gap pill* mid-edge; `CORROBORATES` = solid double hairline; `PRECEDES/SUPERSEDED` = arrow with strike on the older node; `INDEPENDENT_OF` = dotted rule with `eye-slash` (bank stage); `LLM_HYPOTHESIS` = dashed edge with `AI hypothesis` tag and E6 hatched seal. The hero state reads left-to-right as three short sentences: *Browser: unknown → Merchant: no payment → PayPal: captured*, with two tears.

### 8.7 Provenance drawer
Header (claim text) → seal + tier definition → source/trust/env → times (UTC/local, occurred vs observed) → `payload_sha256` short+copy → **redacted payload** viewer (sensitive keys shown as `▒▒ redacted ▒▒` with reason) → signature status (webhooks: `Verified / Rejected / Not applicable`) → **Used by** (claims, Passport versions) → edges. Never shows unredacted payloads.

### 8.8 The "frosted silhouette" (restricted fields)
When a viewer lacks a field, the *shape* of the section appears as a pale hatched block with `lock-lens` and the sentence *"This part isn't shared with you."* — plus what they **can** do (doc 05 §10). The silhouette never contains the hidden content, a length hint, or a blurred rendition of it (no blur-based "peek"): it is a placeholder drawn from the schema, not from the data.

### 8.9 Do-not-repeat banner anatomy
Persistent (not a toast), amber `unknown` surface, `shield-pause` icon, headline in `body-l` 600, one sentence of reason, one live line ("Last checked 14:02:11 · next check in ~30 s"), and a **disclosure** "Why can't I pay again?" (expands the safe-action reasons from R26). It has **no dismiss** control while `DO_NOT_REPEAT_CHECKOUT` is in `allowed[]`.

### 8.10 Charts and graphs policy
1. Charts exist only where the **API returns the series** (e.g., ops queue age distribution, reconcile latency). No decorative or invented data.
2. Every chart: title, unit, time window, **sample size (n)**, source label (`measured / reported / simulated`), and a "View as table" toggle.
3. Axes start at zero for counts; no dual axes; no smoothing that invents points; gaps drawn as gaps.
4. Colour encodes category only with an additional shape/pattern and direct labels (no legend-only decoding). `unknown/no data` is hatched, never zero.
5. Simulated series are hatched violet with `SIMULATED`. Never label simulated metrics "live".
6. Forbidden: gauges/speedometers, "success rate %" tiles without n and window, completion rings.

### 8.11 Case history
A compact audit-style list of case-status transitions (`OPEN → … → RESOLVED`) as a **horizontal path of nodes** with the current node outlined; each node shows cause, actor, evidence ref (from R11/R46-derived, role-filtered). `RESOLVED` is annotated "objective verified — not a bank-settlement claim".

---

## 9. Component state matrix

States: **idle · hover · active · focus · disabled · loading · success · unknown · warning · restricted · simulated**.

| State | Button | Input | Status chip | Card/Row | Banner |
|---|---|---|---|---|---|
| idle | Solid accent (primary) / outline (secondary) | `surface-2`, 1 px `--border-strong` | tone fg/bg | `surface-2` | tone surface |
| hover | `accent-hover`; underline on text buttons | border → `--text-2` | none | rule darkens | none |
| active | Pressed: 1 px inset, no movement | – | – | – | – |
| focus | 2 px `--focus` ring + 2 px offset (never removed) | same ring | same | same | same |
| disabled | 45 % opacity **+ visible reason text** (never silent) | read-only style + helper | – | – | – |
| loading | Label stays; inline `clock-arc`; `aria-busy` | skeleton | skeleton | skeleton in final layout | – |
| success | Brief `seal-check` + "Done" (live region) | `verified` tick + text | verified | verified left rule | verified (rare) |
| unknown | n/a (buttons never inherit unknown) | – | amber + open-ring | amber dashed left rule | amber banner |
| warning (attention) | – | inline attention text | attention tone | attention rule | attention banner |
| restricted | Absent (not disabled) when action isn't permitted | Hidden or frosted silhouette | – | frosted silhouette | – |
| simulated | Violet outline + flask | Violet focus + SIM label | violet hatch | hatched + ribbon | violet ribbon |

Rule: **a control that would be dangerous is *absent*, not disabled** (doc 05 §7.1 — no checkout button when `CREATE_NEW_ORDER` is blocked).

---

## 10. Accessibility

**Target: WCAG 2.2 AA** (SOURCE M08 §8.5). This section lists design obligations; evidence (axe + manual SR/keyboard) is required to claim conformance (AO-13) — **no conformance claim is made here.**

| Area | Obligation |
|---|---|
| Contrast | Text ≥ 4.5:1 (all token pairs in §3 meet it, computed); UI boundaries/icons ≥ 3:1 (`--border-strong`); `--border` is decorative only |
| Colour independence | Every status = icon + label; seals shape-coded; edges shape-coded; charts pattern + direct labels (1.4.1) |
| Focus | 2 px ring + 2 px offset, never clipped; not obscured by sticky bars (2.4.11/2.4.12) — scroll-margin on anchored sections |
| Keyboard | Full operation; roving tabindex in timeline and graph (arrows), `Enter` opens evidence, `T` toggles table, `Esc` closes drawers/dialogs, `?` opens shortcut sheet; skip-to-content; no keyboard traps |
| Targets | ≥ 44 px preferred, ≥ 24 px minimum (2.5.8) |
| Reflow/zoom | No horizontal page scroll at 320 CSS px / 400 % zoom; graph falls back to stacked lanes (1.4.10) |
| Live regions | One polite `role="status"` region (throttled ≤ 1/5 s, coalesced); `role="alert"` only for confirmed errors/security |
| Names | Seals: "Evidence tier E1 — verified provider read"; ids: "Order id ending 7F2A, copy"; icons decorative unless sole content |
| Tables | Graph, timeline, charts each have a table equivalent |
| Motion | §6 reduced-motion table; no flashing > 3/s |
| Forms | Visible labels; error summary with focus management; `aria-invalid`/`aria-describedby`; accessible authentication (no cognitive tests; paste allowed — 3.3.8); redundant entry avoided in wizard (3.3.7) |
| Language | Plain-language copy (§12); target ≈ grade 8 reading level for sender/receiver copy (**ASSUMPTION**, check with a readability tool); `lang` attribute; no jargon without a "What's this?" disclosure |
| Forced colors / high contrast | Tones fall back to `CanvasText/Highlight` + icon + text |
| Screen-reader order | DOM order = visual order; landmarks `header/nav/main/aside`; one `h1` per route |

---

## 11. Component inventory and API expectations

*Names match doc 05 §3/§4 and M20 §15.6 and ship in `/app/dev/gallery` (dev only). `Props` list the essential contract; all components accept `className`, `data-testid`.*

| Component | Purpose | Key props / data | States | Screens |
|---|---|---|---|---|
| `AppShell`, `TopBar`, `NavRail`, `BottomTabs` | Frame | `role`, `env` | loading/offline | 3–16 |
| `EnvironmentBadge` | Always-on mode marker | `environment` (4-value enum) | real/simulated | 3–16 |
| `SimulatedRibbon` | Simulated marker | `label?` | – | 11, 14, 16, any SIM case |
| `StatusChip` | Enum → label+icon+tone | `dimension`, `value` (via `statusPresentation`) | all tones | all |
| `EvidenceSeal` | Tier marker | `tier 1–6`, `trust`, `verified` | tooltip/focus | 5–7, 12 |
| `ProvenanceDrawer` | Evidence detail | `evidenceId` (R30) | loading/restricted/error | 5–7, 12 |
| `RedactedPayload` | Safe JSON viewer | `extract`, `redactions[]` | collapsed/expanded | 12 |
| `CaseHeader` | Amount, ref, env, version, last verified | `case`, `version` | stale/offline | 5–10 |
| `PaymentStateTriad` | Three rails | `provider`, `application`, `receiver`, `lastConfirmed` | unknown/simulated | 3 (mini), 5, 6, 8, 9 |
| `DoNotRepeatBanner` | Safety banner | `safeActions` (R26) | live-check line | 5, 6, 8, 10 |
| `VerdictLine` | One-sentence headline (serif) | `claim` (server text) | unknown/pending | 5, 8, 9 |
| `PassportView` | Confirmed/Unknown/Next action | `passport` (R15) | versioned/diff | 5, 8, 9 |
| `PassportVersionSelector` | Version + diff | R16 | – | 5 |
| `SafeActionPanel` | Allowed/blocked actions | `SafeActionSet` | blocked-with-reason | 5, 10 |
| `LifecycleTimeline`, `TimelineEventRow`, `LaneStrip` | Ledger tape | R11 | virtualised | 6 |
| `TruthGapCanvas`, `LaneColumn`, `ClaimNode`, `GapEdge`, `GapList`, `CompareToggle` | Graph | R28 | table view | 7 |
| `RoleSwitcher`, `LensSwitch`, `LensCompare` | Lenses + permission diff | R15 | restricted | 8, 9 |
| `FrostedSilhouette` | Restricted placeholder | `schemaPath` | – | 8, 9, any |
| `ConsentCard` | Exact fields + notification toggle | R21 | – | 2, 9, 15 |
| `ShareLinkManager` | Create/copy/revoke | R17–R19 | expiry | 5, 15 |
| `MirrorLab`, `ScenarioStage`, `ReplayScrubber`, `LedgerCounter`, `NaiveVsPayBridge` | Simulation lab | R36–R37 | simulated | 11, 16 |
| `EvidenceTable`, `EvidenceFilters` | Evidence explorer | R29 | compact | 12, 14 |
| `OpsQueue`, `ThreeColumnInvestigator`, `WebhookInboxTable`, `ServiceHealthStrip`, `QueueAgeDistribution` | Ops | R32–R34, R47 | measured/simulated | 14 |
| `NotificationList` | Inbox | R40–R41 | restricted | 13 |
| `PrivacyCenter`, `DataRequestForm`, `IntegrationStatus` | Settings | R42–R46 | – | 15 |
| `Walkthrough`, `StepRail`, `ScenarioSelector`, `ModeChip`, `RunHistory` | Demo | R36–R39 | – | 16 |
| `Skeleton*`, `EmptyState`, `ErrorState` (with `CorrelationId`), `RestrictedState`, `OfflineBanner`, `StaleBadge`, `Toast` | System states | – | – | all |
| `Dialog`, `Drawer`, `Tabs`, `Tooltip`, `Disclosure` (Radix-based) | Primitives | – | focus-managed | all |

---

## 12. Copy system

### 12.1 Voice
Calm, specific, plain. State what is known, what is not, and what to do — in that order. No exclamation marks, no blame, no jargon without a disclosure. Short sentences. The product never says it *believes* money moved; it says what a source **confirmed**, and when.

### 12.2 Financial verbs (distinct, never conflated — M20 §8.5)
**Approve** (buyer authorises in PayPal) · **Capture** (merchant collects an approved order) · **Refund** (not in MVP) · **Cancel** (abandon before capture) · **Confirm** (a source verified a fact) · **Reconcile** (compare records; *never moves money*). UI buttons use these verbs exactly; "Retry payment" and "Pay again" do not exist.

### 12.3 Approved vs forbidden phrasing

| Say | Never say |
|---|---|
| "PayPal capture confirmed" | "Payment successful" (unqualified), "Transferred to bank" |
| "Outcome not yet confirmed" | "Payment failed" (for any unknown/timeout) |
| "We're checking the original PayPal record" | "Try again" next to an unresolved attempt |
| "Bank settlement not verified" | "Funds arrived" / "Funds will arrive tomorrow" |
| "Our confirmation was lost" | "PayPal is down" (without `PUBLIC_SOURCE` evidence) |
| "AI hypothesis" (always with E6 seal) | "AI confirmed", any percentage confidence |
| "Simulation — not a real PayPal event" | "Live", "Real", for simulated flows |

### 12.4 Canonical state copy (extends doc 05 §9.2; rendered only via `statusPresentation()`)

| Enum | Label | Helper |
|---|---|---|
| `provider_state=NOT_CREATED` | "No PayPal order yet" | "Nothing has been sent to PayPal." |
| `provider_state=CREATED` | "PayPal order created" | "Waiting for the buyer to approve in PayPal." |
| `provider_state=APPROVED` | "Buyer approved" | "Approval is not a payment. The merchant still has to capture it." |
| `provider_state=CAPTURE_PENDING` | "PayPal is still processing" | "PayPal hasn't finalised this capture. Nothing is final yet." |
| `provider_state=CAPTURE_COMPLETED` | "PayPal capture confirmed" | "Verified from the original PayPal record at {time}." |
| `provider_state=CAPTURE_DENIED` | "PayPal declined this payment" | "Confirmed by PayPal at {time}." (E1/E2 only) |
| `provider_state=REFUNDED / PARTIALLY_REFUNDED` | "Refund recorded by PayPal" | "Verified at {time}. This is separate from the original capture." |
| `provider_state=UNKNOWN` | "Outcome not yet confirmed" | "We're checking the original PayPal record. This does not mean PayPal declined your payment." |
| `application_state=SYNCHRONISED` | "Records agree" | "Our record matches PayPal's." |
| `application_state=WAITING_FOR_WEBHOOK` | "Waiting for PayPal's notice" | "We'll also check the record directly." |
| `application_state=API_TIMEOUT` | "Our confirmation was lost" | "We didn't hear back from our server. Your payment may have gone through." |
| `application_state=RECONCILING` | "Checking the original record" | "Re-reading the same PayPal order. No new payment is created." |
| `application_state=EVENT_GAP` | "A notice is overdue" | "PayPal's notice hasn't arrived. We're checking directly." |
| `application_state=INCONSISTENT` | "Records disagree" | "PayPal and our merchant record differ. We're reconciling — no new payment is needed." |
| `receiver_state=NOT_OBSERVABLE` | "Bank settlement not verified" | "We can't see this stage with the data available here." |
| `receiver_state=NOT_APPLICABLE` | "Not applicable" | "" |
| `case_status=RESOLVED` | "Objective verified" | "The goal of this case was confirmed by evidence. This is not a bank-settlement claim." |
| `case_status=NEEDS_HUMAN_REVIEW` | "A person is reviewing" | "The evidence needs an analyst's decision." |
| attempt `CAPTURE_UNKNOWN` | "Capture outcome not yet confirmed" | "Please don't repeat this payment — we're checking." |
| attempt `REVIEW_REQUIRED` | "Needs a second look" | "We couldn't confirm this automatically." |

### 12.5 Hero messages (SOURCE doc 04 §18.3 + M20 §15.3)
* **Sender:** "PayPal confirmed your payment capture. Our website didn't receive the response. We're reconciling your original payment. Please don't pay again."
* **Receiver:** "A PayPal capture was recorded on the seller's side. Use your official PayPal account activity to see availability and withdrawal. We can't confirm bank credit from the data available here."
* **Analyst:** "Provider says COMPLETED (E1, 14:02:11 UTC). Merchant app recorded API_TIMEOUT (E3, 14:02:02). One order, one capture. Recommended: reconcile the internal merchant order. This does not contact PayPal or move money."

### 12.6 System-state copy
* Offline: "You're offline. What you see may be out of date. Nothing has been changed."
* Slow (> 3 s): "Still checking…" · (> 10 s on a mutation) → outcome-unknown copy + polling; never "retry".
* Mutation network error: "We couldn't confirm what happened. Please don't repeat this payment — we're checking."
* Known read error: "{message} · Reference {correlation_id} [Copy]" + **Try again** (reads only).
* Restricted: "This part isn't shared with you." + what the viewer can do.
* Empty cases: "No payment cases yet. Create one to see how PayBridge tracks what's confirmed and what isn't." [Create a case]
* Live paused: "Live updates paused — refreshing periodically."
* Consent revoked: "Access to this case has been withdrawn. Messages already delivered can't be recalled."
* Internal reconcile dialog: "This updates your own merchant record. It does not contact PayPal or move money."


---

## 13. Screen blueprints

**Template used for every screen.** *Route/roles/mode/API* → *Purpose & user goal* → *Regions (IA)* → *Composition at three sizes* → *Above-the-fold hierarchy* → *Components* → *Actions* → *Interactions* → *States (loading · empty · slow · error · unknown · restricted · simulated)* → *Accessibility* → *Representative copy* → *Acceptance (testable)*. Acceptance IDs are `UX-<screen>-<n>` and feed doc 09's traceability matrix. API dependencies are doc 04 route IDs. **None of this is built.**

---

### Screen 1 — Landing / product introduction

**Route:** `/` (+ `/how-it-works`, `/security`, `/developers`, `/docs` as scroll/section pages) · **Roles:** public · **Mode:** static; illustrations are labelled fixtures · **API:** none required; optional R36 to show the demo CTA only when demo mode is on.

**Purpose & goal.** A judge or merchant understands, in under 60 seconds and without technical background, *why a payment can look failed while money moved* and how PayBridge answers it — then enters the demo with one click.

**Regions (top→bottom).**
1. **Hero** — headline `display-xl` serif: *"Where is your payment — really?"* (SOURCE M08 S1); sub: *"The last confirmed step. What's still unknown. Never retry blindly."*; primary CTA **Explore the verified demo** → `/login` (or `/app/demo` if signed in); secondary **See how it works** (anchor). Right: **the Ledger Split** illustration (below).
2. **Five systems, five answers** — scroll-revealed lanes (Browser · Merchant app · PayPal · Webhook · Bank) each with a different status chip for the *same* illustrative payment; the bank lane is an `open-ring` "Not observable".
3. **One verified story** — a static example **Passport** (Confirmed / Unknown / Next safe action) captioned `EXAMPLE — ILLUSTRATIVE DATA`.
4. **Two lenses** — sender and receiver cards side by side showing *different* sentences for the same payment; the vesica motif between them.
5. **Never retry blindly** — three-step explainer: capture succeeded → our response was lost → we re-read the original; "no second order, no second capture".
6. **Proof on demand** — evidence seal legend E1–E6 with plain definitions.
7. **What we don't claim** — an honest panel: no bank-settlement claim, no production readiness, sandbox-only, "independent project, not affiliated with or endorsed by PayPal".
8. **For developers** — sandbox lab, simulator, OpenAPI link, repo link (placeholders until real), licence badge once it exists.
9. **Footer** — disclaimer, links, theme toggle, reduced-motion note.

**The Ledger Split (hero illustration).** One payment glyph on the left splits into five thin lanes ending in five chips. Chips settle one by one (verify-settle, §6) and the final lane stays an open-ring. Under reduced motion all five are shown at once. Caption: `ILLUSTRATION — NOT LIVE DATA` (doc 05 §19; `fixture-` test ids). It uses **no** PayPal marks or provider UI.

**Composition.** *Desktop (1440):* hero 7/5 columns, text left, illustration right; sections alternate full-bleed paper bands with ruled dividers; max text width 640. *Tablet:* illustration drops under hero text, lanes scroll inside their own container. *Mobile:* single column; lanes become a vertical stack; CTA is full-width and sticky at bottom after the hero scrolls out.

**Above the fold.** (1) Headline, (2) one-line promise, (3) CTA, (4) illustration caption proving honesty. Nothing else.

**Interactions.** Anchored sections with `scroll-margin-top` for sticky header; lane reveal is IntersectionObserver-driven but static when reduced motion; CTA focus ring visible; "See how it works" is a real link to `#how`.

**States.** Loading: static shell renders instantly (no spinner). Error: n/a (static). Demo CTA hidden if R36 fails (server-controlled). Simulated: illustration caption is the marker.

**Accessibility.** One `h1`; landmarks; illustration has a text equivalent (`<figure>` + `figcaption` describing the five answers); no autoplay; contrast per §3; skip link.

**Copy.** CTA: "Explore the verified demo" · honesty panel: "PayBridge reports what a source confirmed, and when. It does not see your bank, and it never moves money." · footer: "Independent project using PayPal developer APIs; not affiliated with or endorsed by PayPal." (SOURCE doc 01 §3).

**Acceptance.** UX-S1-1 page renders with JS disabled for text content · UX-S1-2 illustration carries the visible caption and `fixture-` test ids · UX-S1-3 no PayPal logo/colours anywhere · UX-S1-4 forbidden-copy lint passes (no "failed" for unknown) · UX-S1-5 landing chunk ≤ budget (doc 05 §18) · UX-S1-6 axe: zero serious/critical.

---

### Screen 2 — Authentication and onboarding

**Routes:** `/login`, `/invite/:token` · **Roles:** public → any; receiver invite requires the intended identity · **Mode:** real · **API:** R01 login, R02 logout, R03 me, R04 csrf, R21 accept, R22 decline.

**Purpose & goal.** Sign in with minimal friction; understand *which role you are, which environment you are in, and what you can see*; for receivers, give **informed, revocable consent** before seeing anything.

**2A. Login.** Two-column split: left = calm brand panel (vesica motif + the line "One payment. Two people. One verified story."); right = form (email, password, **Sign in**). Below: when demo mode is on (`GET /demo/scenarios` succeeds) a **Demo identities** panel lists role cards (Sender · Receiver · Analyst · Admin · Demo operator) that **pre-fill the email only** and label "Demo identity — demo data only". The demo password is *not* embedded in the UI or this document; it is supplied in the README/submission notes (**ASSUMPTION**, doc 08).
*Errors:* generic "We couldn't sign you in. Check your details and try again." (no hints, constant time); rate-limited → "Too many attempts. Try again in a few minutes." (R01 429). Session expiry → returned to `/login?next=…` with "You were signed out for your security."

**2B. First-run orientation (3 short steps, skippable, once per account).** (1) *Your role* — what you can see/do, rendered from `GET /auth/me` memberships. (2) *Environments* — side-by-side `REAL SANDBOX` vs `SIMULATED` markers with the sentence "Sandbox is a test environment with no real money. Simulated means PayBridge made the events up to show you what would happen." (3) *Privacy* — one line on consent and revocation. Each step has a progress dots indicator (step count only, not "completion").

**2C. Receiver invitation (`/invite/:token`).** *Gate:* user must be signed in as the **intended identity**; a token alone never authenticates (doc 04 §7). Mismatch → "This invitation was sent to a different account. Sign in with the email it was sent to." (no existence leak of the case). *Consent card* (`ConsentCard`): title "Review what will be shared with you"; a table listing **each field** to be shared (from `proposed_fields`) with a plain-language description and a mini lens preview; **separate**, unchecked-by-default toggle "Send me notifications about this payment"; text "You can withdraw access at any time. Messages already delivered can't be recalled."; buttons **Accept and view** / **Decline**. Expired/revoked → "This invitation is no longer valid." with a "Ask the sender for a new one" hint.

**Composition.** Desktop: split layout; tablet: stacked with brand panel as a header band; mobile: form first, brand panel hidden, demo identities collapse into a disclosure.

**States.** Loading: skeleton form; slow (>3 s): "Still checking…" under the button; unknown: n/a; restricted: invite mismatch (above); simulated: demo identities carry the `DEMO` chip.

**Accessibility.** Labels, `autocomplete` attributes, paste allowed in password, error summary with focus, no CAPTCHA/cognitive test (3.3.8), consent card is a `fieldset` with legend; toggle is a real checkbox with description.

**Acceptance.** UX-S2-1 failed login never reveals which field was wrong · UX-S2-2 notification toggle is unchecked by default and independent of data-sharing · UX-S2-3 accept calls R21 only after explicit click · UX-S2-4 wrong-identity visit shows the mismatch state and no case metadata · UX-S2-5 no token/credential in `localStorage` (FE-11) · UX-S2-6 keyboard-only login and accept.

---

### Screen 3 — Workspace / case dashboard ("Case Ledger")

**Route:** `/app/cases` · **Roles:** sender, merchant_analyst, merchant_admin (demo_operator in sim) · **Mode:** real · **API:** R06 list, R03, R40 (bell count), R47 (provider advisory chip).

**Purpose & goal.** See *which cases need attention because of uncertainty* — not a revenue dashboard. One glance answers: "Is anything unconfirmed, and for how long?"

**Regions.**
1. **Page header** — `h1` "Cases"; **New payment** button (primary).
2. **Attention strip** — a single ruled band: "**3 cases have an outcome not yet confirmed · oldest 14 min**" (computed from R06 data only; hidden if 0). Each item is a link chip. This is *the* hero element; it replaces KPI tiles.
3. **Ledger** — the main list as **ledger rows** (not cards): `[EnvironmentBadge] Reference · amount (tabular) · mini State Triad (three pips with icons) · case_status chip · last verified (relative+absolute) · updated`. Sorted by *uncertainty age* by default; alternate sorts: updated, amount.
4. **Filter rail** — status (case_status), environment, role, "has unknown", search by reference. Filters are URL query params.
5. **Right column (≥ 1280)** — **Live pulse**: provider status advisory chip with `as_of` (R47; `STATUS_UNKNOWN` when unavailable — never "all systems operational"), notification count, and "Demo mode" shortcut if enabled.
6. **View toggle** — *Ledger* (default) | *By status* (kanban-like columns by `case_status`, read-only).

**Composition.** Desktop: ledger 8 cols + pulse 4; tablet: pulse collapses to a strip under the header; mobile: rows become two-line cards (reference + amount / triad pips + status), filters in a sheet.

**Mini State Triad.** Three 12 px pips on one rail: filled seal = verified, open-ring = unknown, clock-arc = pending, offset-arrows = attention; always with an accessible name ("PayPal: confirmed; Merchant app: confirmation lost; Receiver: not observable").

**Interactions.** Row click/Enter → `/app/cases/:id`; arrow keys move between rows (roving); `/` focuses search; SSE-derived highlight when a row's `version` increases (left-rule highlight, §6); filters update URL.

**States.** Loading: row skeletons at final height. Empty: "No payment cases yet…" + **Create a case**. Empty-filtered: "No cases match these filters." + **Clear filters**. Slow: "Still loading cases…". Error: `ErrorState` with correlation id + **Try again**. Unknown: shown per row (amber), and in the attention strip. Restricted: a role without cases sees only those it participates in; there is no "hidden cases" count (no existence disclosure). Simulated: simulated cases are hatched and ribboned; mixed lists keep sim and real visually distinct.

**Accessibility.** The ledger is a real `<table>` (or `role="grid"`) with column headers; the attention strip is `role="status"` (polite); row actions reachable by keyboard; sort controls announce state.

**Copy.** Attention strip: "3 cases have an outcome not yet confirmed. Oldest: 14 min." · none: "Nothing is waiting on confirmation." (only when R06 returns zero unknowns).

**Acceptance.** UX-S3-1 no invented KPI/number; every figure derivable from R06 · UX-S3-2 an `UNKNOWN`/`API_TIMEOUT` case never renders red/"failed" · UX-S3-3 default sort puts oldest unconfirmed first · UX-S3-4 env badge on every row · UX-S3-5 provider advisory shows `as_of` or `STATUS_UNKNOWN` · UX-S3-6 390 px has no horizontal body scroll.

---

### Screen 4 — New payment / connect sandbox / create investigation

**Route:** `/app/cases/new` · **Roles:** sender, merchant_analyst (merchant_admin for provider status) · **Mode:** real · **API:** R05 create case, R08 readiness, R09 create order, R10 capture (after return), R45 provider config booleans (admin), R03.

**Purpose & goal.** Create a payment case, learn honestly whether the route is *known to be supported*, and hand off to PayPal's own approval — with the amount **locked to the case**.

**Structure — a 6-step guided flow with a sticky "Case draft" summary rail.**
1. **Participants** — your org (read-only), receiver email (optional; invitation sent only after case creation) and country.
2. **Invoice** — external reference, purpose. Free text is shown escaped and flagged "treated as untrusted text" in a small help disclosure (prompt-injection surface, doc 07).
3. **Route** — sender country, receiver country. A route chip shows `Route not yet verified` until readiness runs.
4. **Amount & currency** — amount string with live format check; currency `USD` (MVP). Helper: "This amount is saved to the case. The checkout will use this saved amount — not what the browser sends." Illustrative default `USD 2,000.00` is a *placeholder example*, not pre-filled.
5. **Consent** — "Invite the receiver to see their view" toggle + exact field checkboxes (maps to R20 after creation). Off by default.
6. **Readiness & checkout** — runs R08 and shows the **Readiness card**: `Supported / Unsupported / Unknown until account verification` with requirements list (each `status` as chips like `Self-attested`, `Not checked`), `provenance` (source, last checked), and "This is readiness advice, not a provider acceptance decision." Then **Continue to PayPal** (R09). On `201`, show the **Approval hand-off** screen: "You'll approve this payment in PayPal. We never see your PayPal password or card details." → redirect. On return the case shell loads; if `APPROVED`, a **Confirm payment** primary action appears there (R10).

**Connect sandbox panel (admin/analyst only).** `IntegrationStatus` shows booleans from R45: *client id present · secret present · webhook id present* with `seal-check` / `open-ring` and setup guidance pointing at doc 08 (no values, ever). For a non-admin when unconfigured: "Checkout isn't available yet. Ask an administrator to finish the sandbox setup." Environment is always `REAL SANDBOX`; a visible line states "No real money moves in the sandbox."

**Create investigation.** Investigations start automatically on triggers; this screen only offers a quiet link after creation: "What happens if something goes wrong?" opening a one-screen explainer of the unknown-outcome path (non-modal). Manual investigation (R13) lives on Screen 5.

**Composition.** Desktop: form 7 cols + sticky summary rail 5; tablet: summary above; mobile: one step per screen with a bottom action bar and step title (no percent bar; "Step 3 of 6" text only).

**Key interactions.** Inline validation on blur; `Idempotency-Key` generated once per user intent and reused on retry of the same click (doc 05 §5); **Continue to PayPal** disabled while an intent is in flight *and* while the case reports an unresolved attempt; if `ACTIVE_ATTEMPT_EXISTS` is returned: "A payment for this case is already in progress. Open the case to see where it stands." (link to Screen 5; no second-order button).

**States.** Loading: skeleton step; slow checkout create (>10 s): switches to outcome-unknown copy ("We couldn't confirm what happened. Please don't repeat this payment — we're checking.") and navigates to the case; error `VALIDATION_FAILED`: field-level messages + summary; error `PROVIDER_UNAVAILABLE`: show body `outcome` — if `UNKNOWN`, treat as unknown, never "failed"; restricted: receiver role sees "Not available for your role"; simulated: if `demo_operator` creates a sim case, a `SIMULATED` ribbon wraps the whole flow and the PayPal hand-off is replaced by "Simulated checkout".

**Accessibility.** Stepper is an ordered list with `aria-current="step"`; each step is a `fieldset`; error summary focus; amount has `inputmode="decimal"`; the summary rail is `aside` with `aria-label`.

**Acceptance.** UX-S4-1 browser never sends amount to R09/R10 (verified by network assertion) · UX-S4-2 readiness `UNKNOWN` renders amber/"Unknown until account verification", never "Unsupported" · UX-S4-3 no checkout control when an unresolved attempt exists · UX-S4-4 consent toggle off by default · UX-S4-5 R45 shows booleans only · UX-S4-6 mobile step flow completes keyboard/screen-reader only.

---

### Screen 5 — Payment detail and Payment Truth Passport (the centrepiece)

**Route:** `/app/cases/:id` (shell) with tabs to Screens 6–10, 12 · **Roles:** participants (role-projected) · **Mode:** real · **API:** R07 case, R15 passport, R16 versions, R17 share, R26 safe actions, R27 SSE, R13 investigate, R12 reconcile.

**Purpose & goal.** In five seconds: *what is the last confirmed step, what is unknown, and what is safe to do next* — then drill into proof.

**Regions (desktop wireframe).**

```
┌ CaseHeader ──────────────────────────────────────────────────────────────────┐
│ REAL SANDBOX │ INV-2041 · USD 2,000.00 (money-xl)      v7 · verified 14:02:11 │
├──────────────────────────────────────────────────────────────────────────────┤
│ ⚠ DoNotRepeatBanner  "Please don't pay again…"  Last checked 14:02:41  [Why?]  │
├──────────────────────────────────────────────────────────────────────────────┤
│ VerdictLine (serif, display-l):                                               │
│ "PayPal confirmed your payment. Our confirmation was lost."                    │
├──────────────────────────────────────────────────────────────────┬────────────┤
│ STATE TRIAD (three rails)                                         │ EVIDENCE   │
│ PayPal   ──◉ Capture confirmed (E1) ───────────────                │ DOCK       │
│ Merchant ──● … ▒▒ Our confirmation was lost (E3)                   │ 360 px     │
│ Receiver ┄┄┄┄┄┄ ◌ Bank settlement not verified                      │ (selected  │
├──────────────────────────────────────────────────────────────────┤ claim's    │
│ PASSPORT  v7 ▾  [Share] [Compare versions]                         │ seal,      │
│ ┌ Confirmed ──────┐ ┌ Unknown ─────────┐ ┌ Next safe action ─────┐ │ sources,   │
│ │ claims + seals  │ │ bank credit …    │ │ DO_NOT_REPEAT_CHECKOUT │ │ edges)     │
│ └─────────────────┘ └──────────────────┘ └───────────────────────┘ │            │
├──────────────────────────────────────────────────────────────────┴────────────┤
│ Tabs: Passport · Timeline · Truth Gap · Lens · Evidence · Actions                  │
└──────────────────────────────────────────────────────────────────────────────┘
```

**Hierarchy.** (1) VerdictLine (the single dominant sentence; text comes from the server's top confirmed claim or template — never composed client-side) → (2) State Triad → (3) Do-not-repeat banner (visually *above* the verdict whenever active) → (4) Passport trio → (5) tabs.

**Passport trio.** **Confirmed** (claims with `EvidenceSeal`, observed time, privacy-scope chip, "Why?" link → drawer) · **Unknown** (each unknown with "Why is this unknown?" disclosure; bank credit always present for completed captures) · **Next safe action** (headline action from R26 `allowed[0]` with owner and whether human approval is needed; blocked actions listed beneath as "Not available — {reason in words}"). A version selector (`v7 ▾`) loads R15 `?version=`; **Compare versions** shows a line-level diff (added/removed/changed claims) with the version-bump highlight (§6). **Share** opens `ShareLinkManager` (R17): role, expiry, copy link (shown once), "Anyone with the link can see this {role} view until it expires or you revoke it."

**Evidence dock (≥ 1280).** Shows the *selected* claim/chip's seal, tier definition, sources, observed/occurred times, edges, and an **Open in Evidence explorer** link; empty state: "Select a claim to see where it comes from."

**Interactions.** Selecting a claim chip updates the dock (and focuses it only on `Enter`); SSE `passport.version` triggers a polite announcement "Passport updated: provider capture confirmed" and a left-rule highlight; **Refresh original order** (`REFRESH_ORIGINAL_ORDER` → R12) is available to sender/analyst as a read-only action with explanatory text "Re-reads the same PayPal order. Creates nothing."; **Ask for an investigation** (R13) available to participants.

**Hero-state content (illustrative, SOURCE).** Rails: PayPal `Capture confirmed` (E1, 14:02:11) · Merchant `Our confirmation was lost` (E3, 14:02:02) · Receiver `Bank settlement not verified`. Passport: Confirmed → "PayPal confirmed the original capture." + "Our application did not receive the capture response."; Unknown → "Recipient bank credit isn't observable."; Next → "Please don't pay again. Reconciling your original payment."

**States.** Loading: skeletons laid out as verdict + triad + trio. Empty Passport (version 0): "No verified facts yet. We'll add them as PayPal confirms them." Error: ErrorState (reads only). **Unknown:** amber verdict ("Outcome not yet confirmed") with live "Last checked"; the VerdictLine never says "failed". Pending: blue-grey with `clock-arc`. Restricted: sections the role cannot see render `FrostedSilhouette`; a revoked receiver gets the uniform not-found screen and caches are cleared (doc 05 §8). Offline: `StaleBadge` + disabled mutating controls. Simulated: whole shell hatched, ribbon, ids `SIM-`.

**Responsive.** *Tablet:* dock becomes a right-side drawer opened by a "Evidence" button; Passport trio stacks to two columns. *Mobile:* single column in this order: banner → verdict → triad (stacked rails) → **sticky safe-action bar** (headline action) → Passport sections as accordions (Confirmed open by default) → tabs as a bottom sheet list.

**Accessibility.** One `h1` (reference + amount); Passport sections are `h2`s; triad exposes `role="list"` with each rail's meaning; the banner is `role="status"`; version diff available as a text list; the dock is `aside` with `aria-live="off"` (updates announced via the single live region).

**Acceptance.** UX-S5-1 `provider_state=UNKNOWN` ⇒ no red, no "failed", banner present · UX-S5-2 `CAPTURE_COMPLETED` with `receiver_state=NOT_OBSERVABLE` renders both, separately · UX-S5-3 no progress bar/percentage anywhere · UX-S5-4 every confirmed claim has a clickable seal opening provenance · UX-S5-5 version bump via SSE appears without reload (FE-05) · UX-S5-6 checkout control absent when `CREATE_NEW_ORDER` blocked (FE-04) · UX-S5-7 VerdictLine text equals the server field (snapshot test) · UX-S5-8 EnvironmentBadge present (FE-07).

---

### Screen 6 — Payment lifecycle timeline ("ledger tape")

**Route:** `/app/cases/:id/timeline` · **Roles:** participants (role-filtered) · **Mode:** real · **API:** R11 (`?since`), R27 (`timeline.append`).

**Purpose & goal.** Reconstruct *what happened in what order, according to whom* — and see exactly where observation stops.

**Regions.** (1) Toolbar: lane filter chips (Browser · Merchant app · PayPal · Webhook · Receiver), tier filter (E1–E6), environment filter, order toggle (newest/oldest first), **View as table**. (2) The **tape**: time ticks gutter, lane dots, event rows (§8.5). (3) A **Last verified stage** rule across the tape. (4) Detail pane (≥ 1280) showing the focused event: label, lane, seal, trust, environment, occurred vs observed (UTC + local), evidence link.

**Event row.** `[seal] Label · lane glyph · occurred 14:02:11 · observed 14:02:12 · [SIMULATED?]`. Webhook duplicates collapse into one row with a "received 3 times" tag (G06) that expands to show each receipt; out-of-order arrivals show an "arrived late" tag with both times.

**Interactions.** ←/→ (or ↑/↓ on mobile list) move between events (roving tabindex); `Enter` opens `ProvenanceDrawer`; `T` toggles table; hover/focus tooltip (tier, trust, times); virtualised beyond ~200 events (doc 05 §17.1); new SSE events slide in at the top with a polite announcement (throttled).

**Composition.** Desktop: tape 7 cols + detail pane 5. Tablet: tape full-width; detail as modal. Mobile: compact rail (lane glyph + label), detail as full-screen sheet; lane filter in a sheet.

**States.** Loading: tape skeleton. Empty: "No events recorded yet." Error: reads retry. **Unknown:** after the last verified stage, an unfilled dashed segment with "Not yet verified" — never an arrow toward "complete". Restricted: events not in the viewer's projection are absent with a single line "Some events are only visible to the merchant team." (no count). Simulated: hatched rows + ribbon.

**Accessibility.** Tape is a list with `aria-label` "Payment events, newest first"; table view is a true `<table>`; each row's accessible name reads "PayPal capture completed. Evidence tier E1, verified. Occurred 14:02:11 UTC, observed 14:02:12 UTC."

**Copy.** Marker: "Last verified stage: PayPal capture confirmed (14:02:11 UTC)." · After marker: "Not yet verified."

**Acceptance.** UX-S6-1 events after the last verified stage are never drawn as progress · UX-S6-2 duplicate webhooks render as one row with a receipt count · UX-S6-3 occurred vs observed shown for every event · UX-S6-4 table toggle equals tape content · UX-S6-5 keyboard-only traversal of 50 events · UX-S6-6 simulated rows hatched.

---

### Screen 7 — Truth Gap Graph investigation workspace

**Route:** `/app/cases/:id/truth-gap` · **Roles:** sender (reduced), merchant_analyst, merchant_admin; receiver (reduced) · **Mode:** real · **API:** R28 graph, R29 evidence list, R30 evidence detail.

**Purpose & goal.** See *where systems disagree*, *why*, and *which disagreement is a real problem versus an expected blind spot*.

**Layout (desktop wireframe).**

```
┌ Toolbar: [Compare provider vs app] [Hide hypotheses] [View as table] [Zoom −/+] ┐
├───────────────────────────────────────────────────────────────┬────────────────┤
│ Browser   │ Merchant app │ PayPal        │ Webhook      │ Receiver│ GAP LIST      │
│ ◌ UNKNOWN │ ○ NO PAYMENT │ ◉ CAPTURE     │ ● COMPLETED  │ ┄ NOT   │ G01 Browser    │
│  14:02:02 │  14:02:02    │   COMPLETED   │   14:02:13   │  OBSERV.│ uncertainty vs │
│   E3      │   E3         │   14:02:11 E1 │   E2         │  (G04)  │ provider  ~~    │
│     ⟿ tear (G01) ⟿ tear (G02) ⟿                              │ G02 Merchant   │
│                  ═══ corroborates (G06: 2 receipts) ═══       │ record behind… │
├───────────────────────────────────────────────────────────────┴────────────────┤
│ Selected: gap G02 → plain_text from API (template, not AI) + [Open evidence]       │
└───────────────────────────────────────────────────────────────────────────────┘
```

**Regions.** Toolbar; **canvas** (five lane columns, time-ordered rows; deterministic layout, no force simulation); **gap list** (right, from `gaps[]`: rule id, severity, lanes, plain text); **detail strip** (bottom) showing the selected node or gap with evidence chips and **Open in Evidence explorer**.

**Semantics.** Claim chips and edge grammar per §8.6. Gaps carry a **severity** chip: `informational` (G04 bank stage — dotted, neutral), `attention` (G01/G02/G03/G05). The bank gap is *explicitly informational*: "Not a contradiction — this stage can't be observed here."

**Interactions.** Select node → highlights connected edges, dims others, opens provenance; select gap → focuses both claims and scrolls the list; **Compare provider vs app** collapses to two columns with a diff line (first-class toggle, M15 §15.5); **Hide hypotheses** hides `LLM_HYPOTHESIS` edges (default: shown, labelled "AI hypothesis" + E6 hatched seal); arrow keys navigate nodes within and across lanes; `T` table view; `Esc` clears selection.

**Composition.** Desktop: 5 columns + gap list. Tablet: 5 columns scroll horizontally *inside* the canvas container; gap list below. Mobile: **lane stack accordion** with a gap summary on top ("2 disagreements · 1 blind spot"); tap a gap to expand the two involved claims side by side.

**Receiver/sender reduced view.** Senders see their own lane claims plus provider and (summarised) merchant claims; receivers see only provider and receiver lanes. Hidden lanes are shown as a `FrostedSilhouette` column titled "Not shared with you".

**States.** Loading: lane skeletons. Empty: "No comparisons yet — we need at least two sources to compare." Error: reads retry. **Unknown:** lane claims with `UNKNOWN` are amber with open-ring; unknown≠gap by itself. Restricted: silhouettes. Simulated: whole canvas hatched with ribbon; edges tagged `SIMULATED`.

**Accessibility.** Canvas is an SVG with `role="img"` + a full **data table equivalent** (Lane · Claim · Value · Tier · Observed · Edges); nodes are focusable with `aria-describedby` evidence summaries; gap text is real text; zoom ≥ 200 % reflows to the lane stack.

**Copy.** G01: "Browser uncertainty vs provider completion." · G02: "Merchant record behind provider." · G03: "Webhook stale vs provider read." · G04: "Bank stage not observable." · G05: "Amount differs from invoice." · G06: "Same event received {n} times." (templates from API).

**Acceptance.** UX-S7-1 gap `plain_text` rendered verbatim from API · UX-S7-2 hypothesis edges always labelled and E6-sealed · UX-S7-3 bank stage never drawn as a contradiction · UX-S7-4 table view contains every node/edge · UX-S7-5 at 390 px lanes stack with no horizontal body scroll · UX-S7-6 selecting a gap focuses both claims by keyboard.

---

### Screen 8 — Sender Two-Lens explanation

**Route:** `/app/cases/:id/sender` · **Roles:** sender (own), merchant_analyst/admin, demo_operator · **Mode:** real · **API:** R15 (`viewer=sender`), R23 (consent summary), R20 (invite).

**Purpose & goal.** Give the payer *a short, calm, actionable explanation in their own terms* and make the privacy boundary visible: what the other side sees is *different by design*.

**Regions.** (1) **Lens header** — `LensSwitch` showing "You are viewing: Sender lens" with the two-lens mark; (2) **Verdict** (serif): the sender hero message (§12.5); (3) **What we know / What we don't know / What to do** — three short blocks, each claim with a seal; (4) **What your receiver can see** — a compact *shared-facts vesica*: the list of facts that appear in *both* lenses (amount, reference, capture status) and an explicit line "They do **not** see your funding details or your application errors." — drawn from the projection allowlists (doc 04 §18.1), not from the receiver's actual content; (5) **Receiver access** panel: invitation status/consent (R23), **Invite receiver** (R20) when none.

**Analyst/demo extra — "Compare lenses".** For `merchant_analyst`, `merchant_admin`, `demo_operator`: a toggle opens `LensCompare`: sender and receiver Passports **side by side** with a **permission diff** (rows = fields; columns = Sender / Receiver / Analyst; cells = ✔ shown · masked · ✘ absent). This is the visual proof of AO-08 for judges. Senders and receivers never see this view of each other's content.

**Composition.** Desktop: narrative column (max 640) + right rail (receiver access + shared facts); tablet: rail below; mobile: narrative first, sticky safe-action bar, rail as accordion.

**Interactions.** `Why?` on any claim → provenance drawer; **Invite receiver** opens a dialog with field checkboxes and expiry; **Copy hero message** is intentionally *not* offered (messages are generated, not authored here); acknowledging an action (R25) uses a confirm dialog that restates what it does.

**States.** Loading skeleton; empty (no passport yet): "We'll explain this as soon as PayPal confirms something."; **unknown** hero: "We're not sure yet — and that's okay. PayPal's original record is being checked. Please don't pay again."; restricted: n/a for own lens; error: reads retry; simulated: ribbon.

**Accessibility.** The three blocks are labelled regions; compare view is a real table with `scope` headers; ✔/masked/✘ have text, not symbols alone.

**Acceptance.** UX-S8-1 sender DOM/network contains **no** receiver-private fields (FE-06) · UX-S8-2 "what your receiver can see" derives from allowlists, not receiver content · UX-S8-3 hero never says "failed" for unknown/timeout · UX-S8-4 Compare-lenses shown only to analyst/admin/demo_operator · UX-S8-5 permission diff matches the projection tables in doc 04 §18.1 (golden test) · UX-S8-6 invite dialog lists exact fields.

---

### Screen 9 — Receiver Two-Lens explanation

**Route:** `/app/cases/:id/receiver` (entered via `/invite/:token`) · **Roles:** receiver (with active consent), merchant_analyst/admin, demo_operator · **Mode:** real · **API:** R15 (`viewer=receiver`), R23 consent, R24 revoke, R25 ack, R20–R22 (invite lifecycle).

**Purpose & goal.** Tell the receiver *what PayPal has confirmed on the seller side*, keep **PayPal capture** and **bank settlement** visibly separate, show what they can do next in *their official account*, and make revocation easy — without a single sender-private field.

**Regions.** (1) **Lens header** — "Receiver lens" + two-lens mark; (2) **Verdict** (serif): receiver hero (§12.5); (3) **Two separate facts** — a pair of equal-weight panels: **PayPal capture** (seal, status, time) and **Bank settlement** (`open-ring`, "Not verified — we can't see this stage with the data available here"); (4) **Your next step** — checklist from the Safe Action set for the receiver (`CHECK_OFFICIAL_ACCOUNT_ACTIVITY`, `COMPLETE_OFFICIAL_VERIFICATION` when the provider reports a requirement) with an **official PayPal account** external link (opens in new tab; link target is the provider's own site, never a deep link into PayBridge's data); (5) **Privacy & access** — consent summary (fields shared, notifications on/off), **Ask for help**, **Withdraw access** (R24); (6) a quiet "What you can't see here" disclosure listing *categories* withheld (sender funding details, merchant application errors) — categories only, no values or counts.

**Interactions.** **Withdraw access** opens a dialog: "You'll stop seeing this case immediately. Messages already delivered can't be recalled." → confirms → R24 → screen switches to the **uniform not-found / access-withdrawn** state, SSE closed, caches cleared (S9, doc 05 §8). **Ask for help** opens a simple message to the case owner (server route not in inventory → **LATER**; MVP shows "Contact the sender with your case reference" with a copyable reference — no new route invented). Notification toggle edits consent (separate from data sharing).

**Composition.** Desktop: single reading column (max 680) centred, with the privacy panel as a right rail; tablet: rail below; mobile: verdict → two facts stacked → next step → privacy accordion; sticky bar shows the primary next step.

**Fields shown (allowlist, doc 04 §18.1).** Case id/reference (reference only), amount + currency, provider capture status, masked order/capture ids (last 4), receiver's own account verification state (if provided by a supported connector — none in MVP, so `Not available`), seller-side summary. **Absent:** sender funding, issuer response, webhook diagnostics, application timeout trace, other party's email/name beyond masked.

**States.** Loading skeleton; **pre-consent**: only the consent card (Screen 2C) — never the case; unknown: "PayPal's record is still being confirmed. We'll update this page." (no "failed"); revoked: access-withdrawn state; restricted fields: `FrostedSilhouette` with categories; error: reads retry; simulated: ribbon + `SIM-` ids.

**Accessibility.** The two facts are a `dl`/two labelled regions; external link announces "opens in a new tab"; dialog traps focus and returns it; no information conveyed by position alone.

**Copy.** Bank panel: "Bank settlement not verified. We can't see this stage with the data available here." · Next step: "Check your PayPal account activity for availability and withdrawal options." · Privacy: "You can withdraw access at any time."

**Acceptance.** UX-S9-1 receiver DOM/network contains no sender-private or application-trace fields (FE-06) · UX-S9-2 capture and bank settlement are always two separate items · UX-S9-3 revocation removes access within one request cycle and closes SSE · UX-S9-4 pre-consent state shows no case data · UX-S9-5 "What you can't see" lists categories only · UX-S9-6 never says funds "arrived" or "will arrive".

---

### Screen 10 — Next Safe Action panel

**Route:** `/app/cases/:id/actions` (also embedded in Screen 5 as the third Passport block and the mobile sticky bar) · **Roles:** participants; ack by designated actor; internal approval by analyst/admin · **Mode:** real · **API:** R26 safe actions, R25 ack, R12 reconcile, R31 internal reconciliation, R13 investigate.

**Purpose & goal.** Show *exactly what is safe to do* and *why dangerous options are unavailable* — in plain words. The panel is the UI of the deterministic Next Safe Action Engine (doc 04 §17); the AI can request, never unlock.

**Regions.** (1) **Headline action** (largest): the first `allowed[]` item with owner chip and "needs human approval" tag; (2) **Allowed** list: each = icon, title, one-line meaning, owner, preconditions, button *only if the viewer is the designated actor*; (3) **Not available** list — each blocked code shown as a locked row with its **reason in words** (below); (4) **Acknowledgement log** (who acknowledged what, when).

**Action catalogue → UI (SOURCE doc 04 §17.2).**

| Code | Title shown | Button | Human approval |
|---|---|---|---|
| `DO_NOT_REPEAT_CHECKOUT` | "Don't pay again" | none (it's a banner/state) | no |
| `AWAIT_PROVIDER_CONFIRMATION` | "Wait for PayPal's confirmation" | none; shows next check time | no |
| `REFRESH_ORIGINAL_ORDER` | "Re-check the original PayPal order" | **Re-check now** (R12) — read-only | no |
| `RECONCILE_MERCHANT_ORDER` | "Fix our own merchant record" | **Review and approve** (R31) → dialog | **yes** |
| `CHECK_OFFICIAL_ACCOUNT_ACTIVITY` | "Check your PayPal account activity" | external link + **I've checked** (R25) | no |
| `COMPLETE_OFFICIAL_VERIFICATION` | "Complete the verification PayPal asked for" | external link + **Done** (R25) | no |
| `CONTACT_PROVIDER_SUPPORT` | "Contact PayPal support with this evidence" | **Prepare evidence pack** (copyable summary; no send) | yes |
| `REQUEST_HUMAN_REVIEW` | "Ask a person to review" | **Request review** (R13/R25) | no |
| `NONE_REQUIRED` | "Nothing more to do" | none | no |

**Blocked reasons (words).** `CREATE_NEW_ORDER` — "A payment for this case is already in progress." (`ACTIVE_ATTEMPT_EXISTS`) · "We're still confirming the original payment." (`OUTCOME_UNKNOWN_RECONCILE_FIRST`) · "PayPal already captured this payment." (`PAYMENT_ALREADY_CAPTURED`) · `CAPTURE_AGAIN` — "A capture was already requested." · `REFUND/VOID/REVERSE` — "PayBridge doesn't move money." (`NOT_SUPPORTED_NO_MONEY_ACTIONS`) · `MARK_FAILED` — "Unknown is not the same as failed." (`UNKNOWN_IS_NOT_FAILED`) · `MARK_BANK_SETTLED` — "We can't observe this stage." (`NOT_OBSERVABLE`).

**Internal reconciliation dialog (R31).** Title "Fix our own merchant record". Body: "This updates your own merchant record. It does not contact PayPal or move money." Shows the evidence it relies on (seals), a required note field, and **Approve** / **Cancel**. After success the application lane ticks to `Records agree` and the Passport bumps a version.

**The "no checkout control" rule.** When `CREATE_NEW_ORDER` is blocked the UI renders **no** control that creates an order — not even disabled (§9). `START_NEW_ATTEMPT` appears only when the prior attempt is terminal by E1/E2 and requires an explicit confirmation with a fresh intent key; its copy: "Start a new payment for this invoice. This is a separate payment from the previous one."

**Composition.** Desktop: headline + allowed (7 cols) | not-available + log (5 cols). Tablet: stacked. Mobile: sticky bar = headline action; full panel in a sheet.

**States.** Loading skeleton; empty (no actions): "Nothing needs doing right now." only when engine returns `NONE_REQUIRED`; error: reads retry; unknown: headline `AWAIT_PROVIDER_CONFIRMATION` with live check line; restricted: a non-designated viewer sees actions as information with "{Owner} can do this."; simulated: ribbon; actions are inert demonstrations.

**Accessibility.** Blocked rows are `aria-disabled` *list items with text* (not disabled buttons); reasons read as text; dialogs focus-managed; acknowledgements announced via the live region.

**Acceptance.** UX-S10-1 blocked codes always show a reason in words · UX-S10-2 no element in the DOM can start a new order when `CREATE_NEW_ORDER` is blocked · UX-S10-3 `MARK_FAILED`/`MARK_BANK_SETTLED` never appear as allowed · UX-S10-4 internal-reconcile dialog states it does not contact PayPal · UX-S10-5 action text is template/server copy, never AI text · UX-S10-6 non-designated viewers see no action buttons.

---

### Screen 11 — Counterfactual Recovery Mirror (simulation lab)

**Route:** `/app/lab/mirror` · **Roles:** view: demo_operator, merchant_analyst, sender; **run controls only with `permissions.canRunScenario`** (demo_operator; D-15) · **Mode:** **SIMULATED** — real API, simulated data · **API:** R36 scenarios, R37 run, R38 reset; sim case reads R28/R11.

**Purpose & goal.** Let anyone *replay a failure and watch what a naive app would do versus what PayBridge does* — to build intuition and trust — while it is **impossible** to mistake it for a real PayPal event.

**Visual signature.** The *entire* lab sits inside a violet-hatched frame with a permanent top ribbon: **`SIMULATED — not a real PayPal event · LOCAL_SIMULATION`**. All ids are `SIM-…`. The environment badge reads `SIMULATED`. No green "verified" seals appear for simulated events: evidence is tier **E6** hatched (doc 04 I11), and verified-green is reserved for real E1/E2.

**Regions.**
1. **Scenario picker** (left): list from R36 (e.g., lost capture response, duplicate webhook, out-of-order webhook, API unavailable, unknown outcome, reconciliation) with a `ModeChip` (`SIMULATED` / `REAL SANDBOX + FAULT` where applicable) and one-line description. Selecting loads the stage.
2. **Stage** (centre): two parallel tracks labelled **"A naive app"** and **"PayBridge"** replaying the *same* timeline (client click → provider call → response lost → …). Each step is an event card with a seal (E6) and a one-line effect.
3. **Ledger counters** (below stage): for each track `orders created · captures created · customer told`; e.g. naive: `orders 2 · captures 2 · "failed, try again"`; PayBridge: `orders 1 · captures 1 · "unknown — don't pay again"`. Numbers come from the simulated run's ledger (R37 result), never typed in.
4. **Replay scrubber**: step back/forward, play/pause (no autoplay), speed 1×/2×; keyboard stepping.
5. **Verdict strip**: "In this simulation PayBridge re-read the original order and created nothing new." — plus "This is a simulation. It shows the design's behaviour, not a measured result."
6. **Run controls** (only if `canRunScenario`): **Run scenario**, **Reset lab** (R38), **Toggle fault proxy** (R39; only available in sandbox demo mode — hidden otherwise).

**Rules.** The naive track is a *teaching construct*, labelled "Illustrative naive behaviour". Real and simulated data never mix on this screen; if a real case id is passed, show the uniform not-found state. Simulated notifications are never sent to real addresses.

**Composition.** Desktop: picker 3 cols · stage 9 cols; tablet: picker becomes a select; mobile: tracks stack (A naive on top, PayBridge below) with a shared scrubber pinned at the bottom.

**States.** Loading: stage skeleton; empty (no runs): "Choose a scenario to replay." ; running: steps tick forward (reduced motion: step buttons only); error: reads retry/run errors shown with correlation id; unknown: scenario outcomes where the *simulated* provider is `UNKNOWN` use amber + open-ring + hatch; restricted: viewers without run permission see completed runs read-only with "Running scenarios needs a demo operator account."

**Accessibility.** The stage is two ordered lists synchronised to the scrubber with `aria-live="polite"` for the current step (throttled); the ribbon is announced once on entry; counters are real text.

**Copy.** Ribbon: "Simulation — not a real PayPal event." · Naive track caption: "A naive app that retries on timeout." · Footer: "Simulated results illustrate the design. They are not evidence about PayPal's production behaviour."

**Acceptance.** UX-S11-1 every surface carries the SIMULATED ribbon/hatch (AO-11) · UX-S11-2 no verified-green seal on any simulated event · UX-S11-3 counters derive from the run's ledger payload · UX-S11-4 no autoplay; stepping works by keyboard · UX-S11-5 run/reset/fault controls absent without `canRunScenario` · UX-S11-6 ids in the lab are `SIM-` namespaced and never link to real case routes.

---

### Screen 12 — Evidence explorer and provenance

**Route:** `/app/cases/:id/evidence` · **Roles:** merchant_analyst, merchant_admin (+ role-scoped subset for sender/receiver) · **Mode:** real · **API:** R29 list, R30 detail.

**Purpose & goal.** Let an analyst answer "*Where exactly did this fact come from, can I trust it, and what depends on it?*" — and export that reasoning into a support ticket.

**Regions.** (1) **Filter bar**: tier (E1–E6), kind, trust, environment, verified-only, text search over *redacted* fields; saved filter chips ("Provider only", "Simulated only"). (2) **EvidenceTable** (compact): `Seal · Kind · Source trust · Environment · Occurred (UTC) · Observed (UTC) · Verified · Used by n`. (3) **ProvenanceDrawer** (right/overlay) per §8.7. (4) **Edge strip** at the bottom of the drawer: `CORROBORATES / CONTRADICTS / PRECEDES / DERIVED_FROM / REQUIRES / RESOLVES / INDEPENDENT_OF` with the originating `rule_id` and `origin` (`DETERMINISTIC` vs `LLM_HYPOTHESIS`).

**Interactions.** Row `Enter` opens the drawer; `[`/`]` move between rows while the drawer is open; **Copy evidence reference** (id + hash) for tickets; "Used by" items link to the Passport version/claim; "Open in Truth Gap" focuses the node. Webhook evidence shows **Signature: Verified / Rejected** and the inbox duplicate count.

**Composition.** Desktop: table 8 cols + drawer 4 (pinned); tablet: drawer overlays; mobile: rows as cards (seal + kind + time) and drawer as full-screen sheet.

**States.** Loading skeleton; empty: "No evidence matches these filters."; error: retry reads; unknown: `UNKNOWN`-origin evidence is labelled by its kind (e.g., `APP_TRACE: timeout`), never "failure"; restricted: role-scoped subset with the note "Some evidence is only visible to the merchant team." (no count); simulated: hatched rows, filter chip "Simulated".

**Accessibility.** Table with sortable headers announcing sort; drawer is a dialog (modal on small screens, non-modal complementary region on desktop) with focus return; redacted markers have text "redacted: contains personal data".

**Acceptance.** UX-S12-1 payloads are never shown unredacted · UX-S12-2 every row shows tier + trust + environment · UX-S12-3 "Used by" lists claims/Passport versions · UX-S12-4 simulated rows are visually distinct and filterable · UX-S12-5 copying a reference never copies payload content · UX-S12-6 keyboard-only row→drawer→close round trip.

---

### Screen 13 — Notification center

**Route:** `/app/notifications` · **Roles:** any · **Mode:** real · **API:** R40 list, R41 mark read.

**Purpose & goal.** A calm inbox of *what changed that affects you*, with enough context to act — and no urgency theatre.

**Regions.** Tabs: **All · Needs your action · Updates**; list of notification rows: `[icon] Title · one-line summary · case reference · time · unread dot`. Selecting a row expands inline with the **server-generated** text (from validated Passport claims) and an **Open case** link. **Safety-critical** items (e.g., do-not-repeat, access withdrawn) are *pinned at the top with a persistent style*, not dismissible toasts (doc 05 §16.2).

**Rules.** Copy is server-generated; the client never composes status text. Notifications for revoked receivers are never shown. Delivered email/other-channel notices cannot be recalled (stated in Privacy settings).

**Composition.** Desktop: list 7 cols + preview 5; tablet: list with inline expansion; mobile: list only, row opens a sheet.

**States.** Loading skeleton; empty: "You're all caught up. We'll tell you when a payment changes."; error: retry reads; unknown: unknown-outcome notices use the amber unknown tone and never use "failed"; restricted: n/a (own inbox); simulated: simulated notices hatched with ribbon and clearly "Simulation only".

**Accessibility.** List is `role="list"`; unread state has text ("Unread"), not just a dot; **Mark all read** has a confirmation-free undo toast (reads only change state, not payment data).

**Acceptance.** UX-S13-1 unread conveyed by text + dot · UX-S13-2 safety-critical items pinned and non-dismissible while the condition holds · UX-S13-3 all notification text comes from the API · UX-S13-4 simulated notices are marked · UX-S13-5 empty state present · UX-S13-6 keyboard + SR pass.

---

### Screen 14 — Operations / sandbox diagnostics / webhook status

**Routes:** `/app/ops` (queue), `/app/ops/health` (diagnostics), `/app/ops/webhooks` (inbox) · **Roles:** merchant_analyst, merchant_admin (platform_operator JIT) · **Mode:** real (with simulated rows labelled) · **API:** R32 queue, R33 health, R34 inbox, R47 provider status, R12 reconcile, R31 internal reconciliation, R30 evidence.

**Purpose & goal.** Triage by *uncertainty*, investigate with *three columns*, and verify the *event pipeline is healthy* — separating **our own measurements** from **public provider incidents**.

**14A. Queue (`/app/ops`).** Ledger table sorted by `uncertainty` (unknown age) with columns: reference · tenant-bound owner · `case_status` · unknown age · verification age · webhook health (dot + text) · human-review owner. Clicking opens the **ThreeColumnInvestigator** (M15 §15.5):
* **Evidence** (left): compact evidence list with seals; click → drawer.
* **Causal timeline** (centre): the ledger tape (Screen 6) + the **Compare provider vs app** toggle (first-class).
* **Diagnosis & actions** (right): validated claims (confirmed / inference / unknown, each with seals and evidence chips), AI findings labelled **"AI hypothesis" + E6 seal** where applicable, the Safe Action set, and **Approve internal reconciliation** (R31) behind the confirm dialog.
Every alert shows evidence count, severity reasoning, and whether the signal is **measured / reported / simulated** (M15 §15.5).

**14B. Diagnostics (`/app/ops/health`).** `ServiceHealthStrip` with three *separately labelled* panels (M20 §19):
1. **PayBridge application health** (measured): API latency distribution, queue ages, error counts, agent runs — each with `n` and window.
2. **PayPal integration (our own traffic)** (measured): per-endpoint latency, 429/5xx *we* observed, token refresh count, webhook lag, reconcile success.
3. **Public provider incidents** (external, advisory): incident title, component, public timestamp, retrieval freshness; if unavailable → `STATUS_UNKNOWN` with last-success time. Never infers a PayPal outage from our timeouts; the fault-proxy demo only increases *our* error counts and the panel says so.
Charts follow §8.10. A **Verified vs synthetic events** counter distinguishes real vs simulated.

**14C. Webhooks (`/app/ops/webhooks`).** `WebhookInboxTable`: received time · event type · `processing_status` chip (`RECEIVED / REJECTED / CLAIMED / UNCORRELATED / PROCESSED / DUPLICATE / DEAD`) · signature verified (Verified / **Rejected** / n/a) · duplicate count · correlated case · lease/attempts. Filters: status, verified vs simulated. A **Dead-letter** subview lists `DEAD` rows with **Re-queue** (admin only, audited, confirm dialog). `REJECTED` rows show "Rejected: signature check failed" with a spike indicator (security signal) — rendered `attention`/`error` only for *confirmed* signature failures.

**Composition.** Desktop: queue 12 cols → investigator opens as a full-width three-column workspace (4/5/3); tablet: investigator columns become tabs; mobile: investigator is a stepper (Evidence → Timeline → Diagnosis).

**States.** Loading skeleton; empty queue: "No cases need attention."; error: retry reads; unknown: `STATUS_UNKNOWN`, backlog with age ("oldest 14 min"); restricted: non-ops roles get the uniform not-found; simulated: simulated rows hatched; simulated metrics never labelled live.

**Accessibility.** Tables with headers/sort; three-column investigator is `region`s with skip links between columns; chart table equivalents; status chips icon+text.

**Acceptance.** UX-S14-1 own vs public vs simulated metrics are visibly separate panels with sources · UX-S14-2 provider status absence shows `STATUS_UNKNOWN`, never "operational" · UX-S14-3 rejected-signature rows distinguishable from duplicates · UX-S14-4 re-queue requires confirm + is admin-only · UX-S14-5 internal-reconcile dialog copy present · UX-S14-6 every chart has n, window, source and a table view.

---

### Screen 15 — Settings / privacy / consent / integrations

**Routes:** `/app/settings/privacy`, `/app/settings/integration` (team settings are **LATER**) · **Roles:** privacy: any; integration/audit: merchant_admin · **Mode:** real · **API:** R42 privacy, R43 export, R44 deletion request, R23/R24 consents, R17/R18 shares, R45 provider config (booleans), R46 audit.

**15A. Privacy & consent.** Sections: **Consents** (list of `ConsentCard`s: case, who, fields shared, notifications on/off, granted date; **Withdraw**), **Share links** (`ShareLinkManager`: role, expiry, created by, **Revoke**; token shown only at creation), **Retention** (plain-language summary of default handling from R42; "subject to legal review" note), **Your data** — **Request export** (R43) and **Request deletion** (R44) with an explanation that audit records may be retained in restricted form after counsel's direction and that deletion is not instant for all logs (doc 04 §19.3). Revocation confirmation copy: "Access ends now. Messages already delivered can't be recalled."

**15B. Integration (admin).** `IntegrationStatus`: environment `REAL SANDBOX` (fixed; `PRODUCTION` row shown as "Disabled in this version"); booleans *client id · secret · webhook id* (present/missing with setup pointer to doc 08, **never values**); webhook endpoint URL (`PUBLIC_BASE_URL`-derived) with copy; "Last webhook received" time; LLM provider presence (boolean) and mode (`ON`/`DETERMINISTIC ONLY`). **Audit** tab (R46): list with **`chain_ok`** indicator and the sentence "Hash-chain detects tampering under stated assumptions. It is not legal certification."

**Composition.** Desktop: left subnav + content column; tablet/mobile: subnav becomes a top segmented control.

**States.** Loading skeleton; empty consents: "You haven't shared or been shown any cases."; error: retry reads (deletion/export show job id and status); unknown: provider-config unknown → "We couldn't check the setup" (amber); restricted: non-admin sees "Only administrators can see integration settings." (uniform); simulated: demo-mode banner "Demo environment — data may be reset."

**Accessibility.** Destructive actions use dialogs with explicit consequence text and focus return; token reveal has "Copied" live announcement; booleans have text.

**Acceptance.** UX-S15-1 no secret/token value rendered anywhere · UX-S15-2 `PRODUCTION` shown as disabled · UX-S15-3 revocation updates access within one request cycle (cross-checked with FE-06) · UX-S15-4 export/deletion copy states retention limits · UX-S15-5 `chain_ok` text never claims legal certification · UX-S15-6 share token shown once.

---

### Screen 16 — Demo walkthrough / scenario selector

**Route:** `/app/demo` · **Roles:** any in demo mode (run controls need `canRunScenario`, D-15) · **Mode:** real API driving `REAL SANDBOX`, `REAL SANDBOX + FAULT` or `SIMULATED` runs · **API:** R36 list, R37 run, R38 reset, R39 fault proxy.

**Purpose & goal.** Let a judge reproduce the golden demo in a few guided clicks and see *proof* — while always knowing what is real, sandbox, or simulated.

**Regions.**
1. **Walkthrough header** — "The 3-minute story": the **StepRail** with the 12 storyboard steps of doc 10 (case created → buyer approves → capture succeeds → confirmation dropped → unknown shown → read-only reconcile → verified evidence → Passport → Truth Gap → two lenses → safe action → Mirror). Each step shows its **mode chip** (`REAL SANDBOX`, `REAL SANDBOX + FAULT`, `SIMULATED`) and, once reached, **links to the exact screen** and evidence.
2. **Scenario selector** — cards: Demo A (normal sandbox checkout), **Demo B (hero: payment succeeded, our response disappeared)**, Demo C (real webhook then duplicate replay), Demo D (route readiness uncertainty), Demo E (AI unavailable; truth still works). Each card: one-line story, **ModeChip**, what it proves, and **Run** (if permitted) or "View last run".
3. **Run console** — current run's state, `case_id`, environment, and the **ledger counters** (`orders: 1 · captures: 1`) taken from the DB ledger export, plus links: *Open case · Open timeline · Open Truth Gap · Open lenses*.
4. **Controls** — **Reset demo** (R38; confirm dialog: "This resets demo data only. It never touches PayPal."), **Fault proxy** (R39; visible only in sandbox demo mode: "Drops *our* response after PayPal succeeds. PayPal does nothing wrong.").
5. **Fallback** — if sandbox credentials are missing, an amber banner: "Real sandbox isn't configured. Running simulation only." and every card switches to `SIMULATED` automatically (never silently).
6. **Reproducibility checklist** — read-only panel from the run: "orders created: 1 · captures created: 1 · AI money actions: 0 · mutator calls by reconciler: 0" (values from the run's report; shown as *measured in this run* with timestamp — never preset).

**Composition.** Desktop: StepRail left rail (3 cols) + selector/console (9 cols); tablet: rail collapses to a horizontal scroll strip; mobile: rail becomes a vertical accordion; run controls in a sticky bar.

**Interactions.** Running Demo B is a *single guided action* that advances the StepRail in real time via SSE; any step can be replayed; jumping to a screen preserves a "Back to walkthrough" chip.

**States.** Loading skeleton; empty (no runs): "Choose a scenario to begin."; error: run errors with correlation id; unknown: step 5 is rendered *as unknown* (amber) by design; restricted: no run permission → read-only history; simulated: ribbon on any simulated run; mixed (real + sim) runs show both marks per step.

**Accessibility.** StepRail is an ordered list with `aria-current`; mode chips icon+text; run progress announced via live region; confirm dialogs.

**Acceptance.** UX-S16-1 every step and card shows its mode chip · UX-S16-2 missing sandbox credentials force a visible simulation-only banner · UX-S16-3 counters/checklist values come from the run payload · UX-S16-4 Demo B can be completed with keyboard only · UX-S16-5 reset never calls PayPal · UX-S16-6 judge account can reach all 16 screens' relevant demos.

---

## 14. Visual QA checklist

Use for every screen before marking it done (doc 02 §9). Record results in `docs/evidence/` (**NOT RUN** until executed).

| # | Check | Pass criterion |
|---|---|---|
| V1 | Hierarchy | One dominant sentence, one dominant state, one dominant action visible above the fold at 1440/1024/390 |
| V2 | Unknown | Any unknown/timeout state is amber + open-ring + text; never red; no "failed" |
| V3 | Colour independence | Greyscale screenshot still distinguishes every status by icon/label/shape |
| V4 | Seals | Every confirmed claim has a seal with tier text; no percentages anywhere |
| V5 | Environment | `EnvironmentBadge` on every case surface; simulated hatched + ribboned + `SIM-` ids |
| V6 | No false progress | No filled bars/rings implying completion; unobserved stretches dashed |
| V7 | Tokens | No raw hex in components (`no-raw-color` lint); both themes render; forced-colors fallback checked |
| V8 | Typography | Tabular numerals for money/time; serif only for display; min 12 px |
| V9 | Motion | Reduced-motion removes non-essential animation; no flashing |
| V10 | Focus | 2 px ring visible on every control, both themes; not obscured by sticky UI |
| V11 | Targets | Interactive targets ≥ 24 px (≥ 44 px preferred on touch) |
| V12 | Reflow | 320 px / 400 % zoom without horizontal page scroll; graph falls back to stack |
| V13 | States | loading/empty/slow/error/unknown/restricted/simulated all designed and screenshotted |
| V14 | Copy | Passes forbidden-copy lint; financial verbs correct; restricted copy matches §12 |
| V15 | Restricted | Silhouettes contain no hidden content, no length hints |
| V16 | Honesty | No claim of bank settlement, production readiness, or certification in copy |

---

## 15. Acceptance and open items

### 15.1 Design acceptance (feeds doc 09/10)

| ID | Criterion | Verification |
|---|---|---|
| DS-01 | All 16 screens implemented to this spec at 1440/1024/390 in light and dark | Visual regression baselines (doc 09) |
| DS-02 | Every `UX-S*` acceptance item in §13 has a test (unit, component, e2e or manual) | Traceability matrix (doc 09) |
| DS-03 | Token file matches §3 names; no raw colours | Lint |
| DS-04 | Unknown never red / never "failed"; no confidence percentages | Component + e2e + copy lint |
| DS-05 | WCAG 2.2 AA target: axe clean + manual keyboard/SR pass recorded | `docs/evidence/a11y-report.md` |
| DS-06 | Simulated surfaces always marked | UI + API audit (AO-11) |
| DS-07 | Design validation with ≥ 5 task scenarios (M08 §8.6) run and *both passes and failures recorded* | Moderated notes (**NOT DONE**) |

### 15.2 Open items / VERIFY

* **VERIFY** font licences (Newsreader, Inter, IBM Plex Mono) and self-hosting terms; fallback stacks are acceptable if unavailable.
* **VERIFY** all contrast ratios in a real browser after implementation (values here are *computed*, not measured on rendered UI).
* **ASSUMPTION** reading-level target (≈ grade 8) — check with a readability tool; revise copy if needed.
* **ASSUMPTION** the demo password is supplied out-of-band (README/submission notes); not part of the UI.
* **DESIGN DECISION (reversible)** "Ask for help" (receiver) has no server route in the MVP inventory; shown as a copyable reference. A route would need a doc 04 amendment.
* **NOT DONE:** no usability test, no prototype, no screenshots exist. Do not claim user validation (doc 01 §22).

*End of document 06. Next: `07_AI_MULTI_AGENT_ARCHITECTURE_AND_PROMPT_SPEC.md`.*
