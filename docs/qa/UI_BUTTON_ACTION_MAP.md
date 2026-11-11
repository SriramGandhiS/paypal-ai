# PayBridge AI — UI Button & Interactive Action Map

**Component:** Golden Path Demo (`apps/web/index.html` & `apps/web/src/screens/GoldenPathDemoScreen.tsx`)  
**Specification Reference:** Section 15 — Master UI Button Action Map  
**Verification Date:** November 11, 2026  

---

## 1. Interactive Button Inventory & State Transitions

| UI Element / Button | Location | Trigger Event | Destination / Result | Verification Status |
| :--- | :--- | :--- | :--- | :--- |
| **`Pay with PayPal`** | Main Checkout Card | `onclick="gpdOpenSandboxModal()"` | Opens official PayPal Sandbox approval modal with funding options (`Visa •••• 4096`, `PayPal Balance`) | **VERIFIED (Playwright)** |
| **`Approve & Complete Payment`** | PayPal Sandbox Modal | `onclick="gpdApprovePayment()"` | Transitions to Processing animation (1.6s) -> Subtle Failure state (1.4s) -> Automatic AI Investigation | **VERIFIED (Playwright)** |
| **`Cancel and return to merchant`** | PayPal Sandbox Modal | `onclick="gpdCancelSandbox()"` | Dismisses modal without charging; returns cleanly to merchant checkout card | **VERIFIED (Playwright)** |
| **`View/Hide Cryptographic Evidence Vault`** | Diagnosed View | `onclick="gpdToggleEvidence()"` | Toggles expandable drawer displaying SHA-256 provenance hashes without cluttering conversational stream | **VERIFIED (Playwright)** |
| **`Copy Summary`** | Diagnosed View Action Bar | `onclick="gpdCopySummary()"` | Copies formatted investigation diagnosis to system clipboard; displays toast `✓ Investigation summary copied` | **VERIFIED (Playwright)** |
| **`Export Evidence JSON`** | Diagnosed View Action Bar | `onclick="gpdExportJSON()"` | Downloads `paybridge-investigation-<orderRef>.json` containing signed evidence & financial safety invariants | **VERIFIED (Playwright)** |
| **`Request Human Review`** | Diagnosed View Action Bar | `onclick="gpdShowToast(...)"` | Flags case for Tier-2 Human Financial Review; displays toast `✓ Case #<orderRef> flagged for Review` | **VERIFIED (Playwright)** |
| **`New Payment Demo`** | Assistant Header | `onclick="gpdStartNewDemo()"` | Randomly selects a new scenario fixture (avoiding immediate repetition); resets state machine to clean checkout | **VERIFIED (Playwright)** |
| **Quick Chip: `Was I charged twice?`** | Chat Input Bar | `onclick="gpdSendChatQuery(...)"` | Answers with authoritative $0.00 duplicate exposure guarantee & attempt lock confirmation | **VERIFIED (Playwright)** |
| **Quick Chip: `Should I try again?`** | Chat Input Bar | `onclick="gpdSendChatQuery(...)"` | Advises safe next steps based on specific decline type (avoids blind duplicate retries) | **VERIFIED (Playwright)** |
| **Quick Chip: `Where is my money?`** | Chat Input Bar | `onclick="gpdSendChatQuery(...)"` | Explains current funds location across buyer account, provider ledger, and merchant state | **VERIFIED (Playwright)** |
| **Quick Chip: `Why did my bank reject it?`**| Chat Input Bar | `onclick="gpdSendChatQuery(...)"` | Provides specific ISO 8583 response code and issuer rationale | **VERIFIED (Playwright)** |
| **Chat Send Button (`.gpd-send-btn`)** | Chat Input Bar | `onclick="gpdSendChatQuery()"` / `Enter` | Submits free-form user query; AI responds conversationally using active case evidence | **VERIFIED (Playwright)** |

---

## 2. Invariant Protections
- **Zero Thread-Blocking:** Native browser `alert()` is strictly prohibited; all confirmations use non-blocking toasts.
- **Strict Mutex Lock:** No button permits re-executing checkout while an active investigation is running.
- **Read-Only Copilot:** AI interface contains zero payment mutation or transfer buttons.
