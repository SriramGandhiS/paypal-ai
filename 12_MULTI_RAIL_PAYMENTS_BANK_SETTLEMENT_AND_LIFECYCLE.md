# 12 — Multi-Rail Payments, Bank Settlement Observability, and Full Financial Lifecycle

| Field | Value |
|---|---|
| Document | `12_MULTI_RAIL_PAYMENTS_BANK_SETTLEMENT_AND_LIFECYCLE.md` |
| Version | 2.0.0-ENTERPRISE |
| Date | 10 November 2026 |
| Status | **SPECIFICATION — Enterprise Production Extension** |
| Audience | Payment Engineers, Core Banking Architects, Integration Engineers, Compliance Teams |
| Extends | `01_PRODUCT_BLUEPRINT_AND_MASTER_VISION.md`, `03_SYSTEM_DESIGN_AND_ARCHITECTURE.md`, `04_BACKEND_ENGINEERING_SPECIFICATION.md`, `08_PAYPAL_INTEGRATION_SECURITY_PRIVACY_AND_COMPLIANCE.md` |
| Purpose | Expands single-provider sandbox checkout into an enterprise multi-rail payment engine (PayPal Live, Stripe, Wise, Banking Rails), implements full dispute/refund lifecycle, and solves receiver bank-settlement observability. |

---

## Table of Contents
1. [Enterprise Multi-Rail Architecture](#1-enterprise-multi-rail-architecture)
2. [Unified Provider Adapter Specification (Go)](#2-unified-provider-adapter-specification-go)
3. [Live PayPal Production Integration](#3-live-paypal-production-integration)
4. [Secondary Payment Rails (Stripe, Wise, Direct Bank Rails)](#4-secondary-payment-rails-stripe-wise-direct-bank-rails)
5. [Receiver Bank Settlement Observability (Open Banking & Account Aggregators)](#5-receiver-bank-settlement-observability-open-banking--account-aggregators)
6. [Full Financial Lifecycle: Refunds, Disputes & Chargebacks](#6-full-financial-lifecycle-refunds-disputes--chargebacks)
7. [Immutable Double-Entry Financial Ledger Engine](#7-immutable-double-entry-financial-ledger-engine)
8. [Real-Time FX Quoting & Slippage Protection Engine](#8-real-time-fx-quoting--slippage-protection-engine)

---

## 1. Enterprise Multi-Rail Architecture

In production, no institutional platform relies solely on a single payment network or sandbox environment. PayBridge AI abstracts the payment network into a unified, high-reliability routing layer that preserves the platform's core truth: **no AI ever moves money, and reconciliation is strictly read-only.**

```
                                [PAYBRIDGE CORE]
                           Case & State Reducer Engine
                                        │
                         ┌──────────────┴──────────────┐
                         ▼                             ▼
                 [MUTATOR ENGINE]              [RECONCILER ENGINE]
              (Human Checkout Only)             (Read-Only GETs)
                         │                             │
    ┌────────────────────┼─────────────────────────────┼────────────────────┐
    ▼                    ▼                             ▼                    ▼
[PayPal Adapter]   [Stripe Adapter]              [Wise Adapter]      [Banking Rail Adapter]
- Orders v2        - PaymentIntents API          - Payouts API       - Open Banking / Plaid
- Payments v2      - Charge & Refund APIs        - Multi-Currency    - Account Aggregators (IN)
- Webhooks v1      - Webhook Event Gateway       - Wire Traces       - ISO 20022 Camt.053
```

---

## 2. Unified Provider Adapter Specification (Go)

The Go backend separates mutating operations from read-only reconciliation across all provider implementations.

```go
package provider

import (
    "context"
    "time"
    "github.com/shopspring/decimal"
)

// RailIdentifier specifies the payment rail
type RailIdentifier string

const (
    RailPayPal       RailIdentifier = "PAYPAL"
    RailStripe       RailIdentifier = "STRIPE"
    RailWise         RailIdentifier = "WISE"
    RailBankTransfer RailIdentifier = "BANK_WIRE"
)

// Reader is strictly GET-only across all provider adapters
type Reader interface {
    GetTransaction(ctx context.Context, id string) (TransactionDetails, error)
    GetSettlementStatus(ctx context.Context, id string) (SettlementStatus, error)
    VerifyWebhookSignature(ctx context.Context, payload []byte, headers map[string]string) (bool, error)
}

// Mutator is human-triggered ONLY and structurally forbidden from AI caller graphs
type Mutator interface {
    InitiateIntent(ctx context.Context, in CreateIntentInput) (IntentResult, error)
    ExecuteCapture(ctx context.Context, in CaptureInput) (CaptureResult, error)
    InitiateRefund(ctx context.Context, in RefundInput) (RefundResult, error)
}

// SettlementObserver queries banking rails and account aggregators for end-to-end receipt
type SettlementObserver interface {
    VerifyReceiverDeposit(ctx context.Context, in DepositVerificationQuery) (DepositEvidence, error)
}
```

---

## 3. Live PayPal Production Integration

Upgrading from PayPal Sandbox to Live Production requires implementing production security profiles and merchant credential governance:

### 3.1 Live Credentials & Configuration
- Live Host: `https://api-m.paypal.com` (strictly guarded by mutual environment flags).
- OAuth 2.0 with Hardware Token Storage: Client secrets are retrieved via KMS envelope encryption and cached in memory for 3,600 seconds.
- PayPal Partner Referrals: Support for multi-merchant marketplace payouts using PayPal Onboarding & Connected Merchant IDs.

### 3.2 PayPal Payouts v1 for Disbursals
When a merchant authorizes cross-border distribution, the Payouts API (`POST /v1/payments/payouts`) is utilized with guaranteed client request IDs (`PayPal-Request-Id`) to avoid double disbursements.

---

## 4. Secondary Payment Rails (Stripe, Wise, Direct Bank Rails)

### 4.1 Stripe Integration
- **Engine**: Stripe PaymentIntents API with automatic `idempotency_key` propagation.
- **Webhook Handlers**: Handling `payment_intent.succeeded`, `payment_intent.payment_failed`, `charge.dispute.created`, `charge.refunded`.
- **Reconciliation**: Strictly calls `GET /v1/payment_intents/:id`.

### 4.2 Wise Integration (For B2B Cross-Border Wire Transfers)
- **Engine**: Wise Payouts and Multi-Currency Accounts API.
- **Tracking**: Wise Transfer Status (`outgoing_payment_sent`, `bounced_back`, `funds_refunded`).
- **Proof of Payment**: Wise provides a cryptographically signed digital receipt containing banking clearing numbers (UETR / IMAD / Fedwire details), converted directly into Evidence Tier E1.

---

## 5. Receiver Bank Settlement Observability (Open Banking & Account Aggregators)

In the initial sandbox specification, the recipient bank settlement stage was defined as `receiver_state = NOT_OBSERVABLE`. In Tier-1 enterprise production, this gap is closed using **Open Banking and Regulated Account Aggregators**.

```
[Sender Pays $5,000] ──▶ [PayPal Captures Funds] ──▶ [Cross-Border Disbursal]
                                                               │
                                                               ▼
                                                  [Recipient Bank Account]
                                                               │
                                         ┌─────────────────────┴─────────────────────┐
                                         ▼                                           ▼
                               [US / EU: Plaid / Yapily]                   [India: Sahamati AA / Setu]
                               - Real-time Balance Checking                - Financial Information User (FIU)
                               - Plaid Transactions API                    - Consent-Artifact Cryptographic Flow
                               - Webhook: `TRANSACTIONS_SYNC`              - Bank Statement Signed JSON Fetch
                                         │                                           │
                                         └─────────────────────┬─────────────────────┘
                                                               ▼
                                               [PayBridge Settlement Observer]
                                               Matches: Date, Amount, Sender Ref
                                                               ▼
                                               State Updated to: `VERIFIED_BANK_CREDITED`
```

### 5.1 Account Aggregator Integration Protocol (India: Sahamati / RBI AA Framework)
1. **Consent Handle Request**: Merchant initiates a consent request via Setu or Sahamati AA specifying `TransactionHistory` for the target beneficiary account.
2. **User Notification**: Recipient approves consent on their mobile banking AA application.
3. **Signed FI Data Fetch**: PayBridge retrieves the encrypted Financial Information (FI) packet from the Financial Information Provider (FIP / Bank).
4. **Deterministic Match Verification**:
   - `deposit.amount == canonical_invoice.net_amount`
   - `deposit.occurred_at >= capture.occurred_at`
   - `deposit.narration CONTAINS payment_reference`
5. **State Transition**: Upon positive match, `receiver_state` transitions from `REPORTED_PENDING` to **`VERIFIED_BANK_CREDITED`** backed by Tier E1 cryptographic bank evidence.

### 5.2 Plaid Integration (US & Global Open Banking)
- Recipient connects account via Plaid Link (enforcing read-only `Transactions` scope).
- PayBridge listens for `SYNC_UPDATES_AVAILABLE` webhook, executes `/transactions/sync`, and validates transaction credit details against the open case record.

---

## 6. Full Financial Lifecycle: Refunds, Disputes & Chargebacks

A production system must handle reverse money movement without creating financial anomalies.

```
                                [PAYMENT STATE: CAPTURE_COMPLETED]
                                                │
                     ┌──────────────────────────┴──────────────────────────┐
                     ▼                                                     ▼
           [REFUND WORKFLOW]                                     [DISPUTE WORKFLOW]
    - Human Analyst Initiated                              - Provider Webhook Received
    - Input: Amount, Reason, Auth Token                    - State: `DISPUTE_OPENED`
                     │                                                     │
                     ▼                                                     ▼
    [Validation: Amount <= Net Balance]                   [Evidence Assembly Agent (A08)]
                     │                                    - Compiles Payment Truth Passport
                     ▼                                    - Compiles Bank Verification Proof
        [Provider Mutator Execute]                                         │
        (Idempotent Refund Request)                                        ▼
                     │                                    [Human Ops Review & Submission]
                     ▼                                                     │
      State: `PARTIALLY_REFUNDED`                                          ▼
             or `REFUNDED`                                 State: `DISPUTE_WON` or `LOST`
```

### 6.1 Dispute Evidence Generator (Agent A08 Extension)
When a chargeback is received (`webhook: charge.dispute.created` or `CUSTOMER_DISPUTE_CREATED`):
1. **Case Lock**: The payment case status moves to `NEEDS_HUMAN_REVIEW` with substate `DISPUTE_PENDING`.
2. **Automated Evidence Dossier**: PayBridge extracts the cryptographic audit chain:
   - Original checkout authorization payload.
   - Verified IP address, user agent, and timestamp.
   - Provider settlement receipt.
   - Receiver bank deposit confirmation (from Account Aggregator).
3. **Submission**: Generates a standardized PDF/JSON evidence packet formatted for PayPal Resolution Center or Stripe Dispute API.

---

## 7. Immutable Double-Entry Financial Ledger Engine

To meet financial audit requirements (GAAP / IFRS), PayBridge implements an append-only double-entry ledger. Every payment event must maintain a zero-sum balance:

```sql
CREATE TABLE ledger_accounts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL,
  account_type TEXT NOT NULL CHECK (account_type IN ('ASSET', 'LIABILITY', 'EQUITY', 'REVENUE', 'EXPENSE')),
  currency CHAR(3) NOT NULL,
  balance NUMERIC(20,6) NOT NULL DEFAULT 0.000000,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE ledger_entries (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL,
  case_id UUID REFERENCES payment_cases(id),
  attempt_id UUID REFERENCES payment_attempts(id),
  description TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE ledger_postings (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  entry_id UUID NOT NULL REFERENCES ledger_entries(id),
  account_id UUID NOT NULL REFERENCES ledger_accounts(id),
  amount NUMERIC(20,6) NOT NULL, -- Positive for Debit, Negative for Credit
  currency CHAR(3) NOT NULL,
  CHECK (amount != 0)
);

-- Constraint Invariant: Every entry MUST sum to exactly ZERO across its postings.
CREATE OR REPLACE FUNCTION verify_ledger_entry_balanced()
RETURNS TRIGGER AS $$
DECLARE
  v_sum NUMERIC(20,6);
BEGIN
  SELECT coalesce(SUM(amount), 0) INTO v_sum
  FROM ledger_postings
  WHERE entry_id = NEW.entry_id;
  
  IF v_sum != 0 THEN
    RAISE EXCEPTION 'Ledger entry is not balanced. Net sum must be 0, got %', v_sum;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;
```

---

## 8. Real-Time FX Quoting & Slippage Protection Engine

Cross-border payments between different currencies (e.g., USD -> INR or EUR -> USD) are subject to foreign exchange rate volatility:
1. **Guaranteed Rate Locks**: FX quotes are retrieved from authorized market makers (or provider APIs) with an explicit `quote_id`, `rate`, and a strict 15-minute expiration timestamp (`expires_at`).
2. **Slippage Threshold**: If the effective exchange rate deviates by more than **0.25%** at the moment of provider settlement, the checkout mutator halts and requests explicit customer re-approval.
3. **No Synthetic FX Rates**: If live liquidity rates are unavailable, the system renders `FX_RATE_UNAVAILABLE` rather than estimating or hallucinating values.
