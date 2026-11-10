# PAYBRIDGE AI — FRONTEND BUTTON & INTERACTION ACTION MAP

**Document Version:** 2.0  
**Target Branch:** `feature/frontend-complete`  
**Standard:** **ZERO DEAD BUTTONS · ZERO FABRICATED SUCCESS**

---

## 1. COMPREHENSIVE INTERACTION REGISTRY

Every visible button, link, selector, and form element across the 16 canonical screens is mapped below with its exact action, state change, and disabled/failure behavior.

### Screen 1: Landing Showcase (Distill Comparison Hero)

| Element Label | Type | Action Triggered | Target / API Call | Failure / Disabled Behavior |
| :--- | :--- | :--- | :--- | :--- |
| **"Explore Distill"** | Pill Button | Transitions view to Screen 3 (Dashboard Matrix) | `onNavigate('dashboard')` | Always active; smooth scroll & tab switch |
| **"Start a free trial"** | Nav Pill | Navigates to Screen 4 (New Payment Checkout) | `onNavigate('new-payment')` | Always active |
| **"Book a demo"** | Nav Pill | Navigates to Screen 16 (Golden Path Stepper) | `onNavigate('demo')` | Always active |
| **"Sign in"** | Nav Link | Navigates to Screen 2 (RBAC Auth) | `onNavigate('auth')` | Always active |
| **"PayBridge Suite ↗"**| Nav Pill | Navigates to Screen 3 (Dashboard Matrix) | `onNavigate('dashboard')` | Always active |
| **Conveyor Belt** | Drag Area | Inertial scrub left/right; dragging right rebuilds cards | Local requestAnimationFrame loop | Inertia dampens smoothly; capped at period |
| **"Launch Golden Path Stepper"** | Hero Button | Navigates to Golden Path Stepper | `onNavigate('demo')` | Always active |
| **"View Cases Register"** | Hero Button | Navigates to Cases Dashboard | `onNavigate('cases')` | Always active |

### Screen 2: Login & RBAC Authentication

| Element Label | Type | Action Triggered | Target / API Call | Failure / Disabled Behavior |
| :--- | :--- | :--- | :--- | :--- |
| **"Sender / Business Payer"** | Identity Card | Authenticates as `usr-sender-991` (Acme Corp Payer) | Stores session token, sets role `sender`, navigates to dashboard | Emits feedback toast notification |
| **"Receiver / Contractor"** | Identity Card | Authenticates as `usr-receiver-102` (Priya Kasthuri) | Stores session token, sets role `receiver`, navigates to dashboard | Gated behind consent dialog |
| **"Merchant Operations Analyst"** | Identity Card | Authenticates as `usr-ops-881` (Triage Analyst) | Stores session token, sets role `merchant_analyst`, navigates to dashboard | Full administrative access |
| **"Judge / Demo Operator"** | Identity Card | Authenticates as `usr-judge-demo` with sandbox bypass | Sets role `demo_operator`, navigates to Golden Demo | Enables chaos proxy injection |

### Screen 3: Cases Dashboard Matrix

| Element Label | Type | Action Triggered | Target / API Call | Failure / Disabled Behavior |
| :--- | :--- | :--- | :--- | :--- |
| **Case Row (e.g. #PB-8821)** | Clickable Row | Selects active case and navigates to Screen 5 | `setSelectedCaseId(id)`, `onNavigate('cases')` | Accessible via Enter key on focus |
| **"New Case" / "Create Case"**| Modal Trigger | Opens new case submission dialog | Local React modal state | Escape or backdrop click dismisses |
| **Status Filter Dropdown** | Select Input | Filters cases table by state (`OBSERVING`, `NEEDS_REVIEW`)| Updates filtered table state | If 0 matches, renders empty state panel |
| **Search Input** | Text Input | Filters cases by external reference or invoice ID | Instant substring filter on table rows | Debounced; clear button restores view |
| **"Refresh" Button** | Action Button | Re-queries `/api/v1/cases` | `api.listCases()` | Shows loading spinner; falls back to seed |

### Screen 4: New Payment & Checkout Timeout Hero

| Element Label | Type | Action Triggered | Target / API Call | Failure / Disabled Behavior |
| :--- | :--- | :--- | :--- | :--- |
| **"Approve & Capture via PayPal"** | Primary Button | Submits single-attempt order capture | `POST /v2/checkout/orders` with idempotency key | Disabled while in-flight or after timeout |
| **Simulated 504 Timeout** | Fault Trigger | Emulates upstream network drop | Emits state transition `API_TIMEOUT`, `UNKNOWN` | Locks button into "Retry Blocked" |
| **"DO NOT REPEAT CHECKOUT"** | Banner Action | Displays safety rationale; navigates to investigation | `onNavigate('cases')` | High-contrast amber warning |

### Screen 5: Payment Detail & Truth Vector

| Element Label | Type | Action Triggered | Target / API Call | Failure / Disabled Behavior |
| :--- | :--- | :--- | :--- | :--- |
| **"Back to Cases"** | Nav Link | Returns to Case Dashboard | `setSelectedCaseId(null)`, `onNavigate('dashboard')` | Preserves table filters |
| **"Verify Truth Passport"**| Action Button | Navigates to Signature Truth Passport | `onNavigate('passport')` | Always active |
| **"Query Original Order (GET)"** | Safe Action | Dispatches read-only PayPal capture poll | `api.getCase(caseId)` | Read-only; zero mutation |
| **"Retry Checkout Mutation"** | Blocked Action | Strictly disabled by safety engine | N/A (Blocked) | Shows cursor-not-allowed + tooltip rationale |

### Screen 6: Lifecycle Timeline

| Element Label | Type | Action Triggered | Target / API Call | Failure / Disabled Behavior |
| :--- | :--- | :--- | :--- | :--- |
| **Event Step Row** | Accordion Item | Expands cryptographic event details and tier | Local expand state | Renders verified SHA-256 hash stamp |
| **"Filter Tiers"** | Tab Bar | Filters timeline to Tier 1 (Direct), Tier 2 (Webhook)| Filter state update | Renders only matching events |

### Screen 7: Truth Gap Graph

| Element Label | Type | Action Triggered | Target / API Call | Failure / Disabled Behavior |
| :--- | :--- | :--- | :--- | :--- |
| **Graph Nodes (P, A, R)**| SVG Node | Highlights node and renders explanatory evidence extract | Local selected node state | Outlines active node; updates side card |
| **"Reset Topology"** | Control Button | Re-centers and zooms multi-dimensional graph | Viewport reset | Smooth transition |

### Screen 8: Sender Two-Lens View

| Element Label | Type | Action Triggered | Target / API Call | Failure / Disabled Behavior |
| :--- | :--- | :--- | :--- | :--- |
| **"Copy Case Receipt"** | Action Button | Copies verified case reference to clipboard | `navigator.clipboard.writeText()` | Emits confirmation toast |
| **"Contact Merchant Support"**| Action Button | Opens calm support inquiry draft with case token | Mailto / Modal | Pre-fills invoice and transaction ID |

### Screen 9: Receiver Two-Lens View

| Element Label | Type | Action Triggered | Target / API Call | Failure / Disabled Behavior |
| :--- | :--- | :--- | :--- | :--- |
| **"Grant Consent & View Case"**| Consent Button | Grants dual-lens access under statutory privacy laws | Updates `consentGranted=true` | Persists in session state |
| **"Revoke Consent"** | Danger Action | Revokes consent and conceals counterparty banking view | Updates `consentGranted=false` | Immediately hides financial projection |

### Screen 10: Next Safe Action Engine

| Element Label | Type | Action Triggered | Target / API Call | Failure / Disabled Behavior |
| :--- | :--- | :--- | :--- | :--- |
| **"Execute Safe Re-read"** | Primary Action | Triggers read-only provider status poll | `GET /v2/checkout/orders/:id` | Prohibits money movement |
| **"Retry Payment"** | Blocked Action | Prohibited by state lattice | Disabled | Permanently disabled while outcome unknown |

### Screen 11: Payment Truth Passport

| Element Label | Type | Action Triggered | Target / API Call | Failure / Disabled Behavior |
| :--- | :--- | :--- | :--- | :--- |
| **"Download Passport (JSON)"**| Action Button | Exports sealed cryptographic certificate | Blob download `passport-:id.json` | Generates verified SHA-256 digest |
| **"Copy SHA-256 Seal"** | Icon Button | Copies certificate hash to clipboard | `navigator.clipboard.writeText()` | Shows "Copied!" badge |

### Screen 12: Evidence Explorer

| Element Label | Type | Action Triggered | Target / API Call | Failure / Disabled Behavior |
| :--- | :--- | :--- | :--- | :--- |
| **Evidence Card Row** | Clickable Card | Expands raw JSON payload extract and signatures | Toggle expanded item state | Displays redacted payload |
| **"Verify Integrity"** | Action Button | Recalculates payload SHA-256 against sealed hash | `crypto.subtle.digest('SHA-256')` | Displays green "VERIFIED MATCH" chip |

### Screen 13: Notification Center

| Element Label | Type | Action Triggered | Target / API Call | Failure / Disabled Behavior |
| :--- | :--- | :--- | :--- | :--- |
| **"Mark All Read"** | Action Button | Clears unread alert indicators | Local state | Updates badge counter to 0 |
| **Notification Row** | Alert Row | Navigates directly to cited case | `onNavigate('cases')` | Highlights relevant case |

### Screen 14: Operations & Webhook Ingestion

| Element Label | Type | Action Triggered | Target / API Call | Failure / Disabled Behavior |
| :--- | :--- | :--- | :--- | :--- |
| **"Replay Webhook Event"** | Action Button | Simulates duplicate webhook delivery to test idempotency | `POST /api/v1/webhooks/paypal` | Returns HTTP 200 without duplicate state |
| **"Clear Ingest Log"** | Action Button | Clears live webhook stream monitor | Resets event stream state | Emits confirmation prompt |

### Screen 15: Privacy & Consent Settings

| Element Label | Type | Action Triggered | Target / API Call | Failure / Disabled Behavior |
| :--- | :--- | :--- | :--- | :--- |
| **DLP Bank Account Masking**| Toggle Switch | Toggles automatic account number masking | Local state | Enforced in all Sender views |
| **Token Spend Limit Cap** | Slider Input | Sets daily LLM token cost ceiling ($5.00–$50.00) | Local state | Gated at $50.00 |

### Screen 16: Golden Path Stepper

| Element Label | Type | Action Triggered | Target / API Call | Failure / Disabled Behavior |
| :--- | :--- | :--- | :--- | :--- |
| **"Next Step" / "Advance"**| Primary Button | Advances to next scenario stage (1 through 7) | Updates stepper stage state | Disabled on Stage 7 (Complete) |
| **"Reset Simulation"** | Action Button | Resets stepper back to Stage 1 | Resets stage state to 1 | Restores initial order fixture |
| **Stage Pills (1–7)** | Stepper Tab | Direct navigation to completed stages | Stage selection | Future stages locked until reached |
