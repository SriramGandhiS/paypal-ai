# PayBridge AI — Payment Scenario Fixture Coverage Matrix

**Component:** Golden Path Scenario Architecture  
**Specification Reference:** Sections 3, 8, 10  
**Verification Date:** November 11, 2026  

---

## 1. Automated Scenario Fixture Pool

| Fixture ID | Category | Order Reference | Merchant & Product | What Happened | Root Cause / ISO Code | What to Do Next | What to Avoid |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| **`INTL_TRANSACTIONS_DISABLED`** | Cross-Border Compliance | `8K110492` | Berlin Quantum Labs (€740.00 EUR) | Bank blocked cross-border transaction | ISO 8583 Code 57 (Transaction Not Permitted) | Enable international usage in banking mobile app | Do not retry checkouts without enabling toggle |
| **`HTTP_504_GATEWAY_TIMEOUT`** | Edge Network Timeout | `9K201824` | Apex Robotics ($480.00 USD) | PayPal captured funds, edge TCP socket dropped | Envoy Ingress 504 Gateway Timeout | Do NOT pay again; merchant auto-reconciles | Do not click checkout again; mutex lock active |
| **`INSUFFICIENT_FUNDS`** | Card Issuer Decline | `5M910482` | CloudScale Infra ($1,250.00 USD) | Issuer declined authorization | ISO 8583 Code 51 (Insufficient Funds / Over Limit) | Check balance or select alternative linked card | Avoid rapid successive re-attempts |
| **`TRANSACTION_LIMIT_EXCEEDED`** | Card Spending Velocity | `2A491820` | Precision Lab ($4,800.00 USD) | Purchase exceeded daily transaction limit | ISO 8583 Code 61 (Exceeds Withdrawal Limit) | Call bank card service line to elevate daily limit | Do not re-attempt immediately without calling bank |
| **`SPLIT_BRAIN_SUCCESS_CRASH`** | Merchant Infrastructure | `7L449102` | Atlas Cloud ($350.00 USD) | PayPal captured, merchant DB crashed before commit | Postgres SQLState 53300 (Connection Pool Exhaustion) | PayBridge dispatches capture proof to merchant | Do not re-purchase item; settlement is guaranteed |
| **`CURRENCY_MISMATCH_FX`** | Foreign Exchange Corridor | `3N889210` | Tokyo Precision (¥88,000 JPY) | FX conversion quote expired mid-checkout | 10-minute spot guarantee window exceeded (>14m) | Request refreshed real-time FX quote from PayPal | Do not authorize stale rate terms |

---

## 2. Dynamic Randomization Rules
- **Non-Repeating Selection:** `gpdPickRandomFixture()` selects from the supported fixture pool, strictly preventing immediate repetition between consecutive runs.
- **Pre-Payment Concealment:** The active scenario is completely invisible to the user until PayPal approval and AI diagnosis have completed.
- **Authoritative Provenance:** Each fixture provides authentic SHA-256 evidence records tied to Tier 1 provider reads and Tier 2 webhooks.
