# 13 — Regulatory Compliance, KYC/AML, and Enterprise Data Governance

| Field | Value |
|---|---|
| Document | `13_REGULATORY_COMPLIANCE_KYC_AML_AND_DATA_GOVERNANCE.md` |
| Version | 2.0.0-ENTERPRISE |
| Date | 10 November 2026 |
| Status | **SPECIFICATION — Enterprise Production Extension** |
| Audience | Chief Compliance Officers, General Counsel, Fintech Architects, Security Auditors |
| Extends | `01_PRODUCT_BLUEPRINT_AND_MASTER_VISION.md`, `04_BACKEND_ENGINEERING_SPECIFICATION.md`, `08_PAYPAL_INTEGRATION_SECURITY_PRIVACY_AND_COMPLIANCE.md` |
| Purpose | Formalizes statutory cross-border compliance (RBI, e-FIRC, OFAC, FinCEN), enterprise IAM (SAML/OIDC/SCIM), automated KYC/AML screening, and institutional data protection (SOC 2, PCI-DSS, GDPR, DPDP). |

---

## Table of Contents
1. [Cross-Border Statutory Framework (US ↔ India & Global Corridors)](#1-cross-border-statutory-framework-us--india--global-corridors)
2. [RBI Purpose Codes & Automated e-FIRC Ingestion](#2-rbi-purpose-codes--automated-e-firc-ingestion)
3. [Automated KYC/KYB & Sanctions/AML Screening](#3-automated-kyckyb--sanctionsaml-screening)
4. [Enterprise IAM: SAML 2.0, OIDC, WebAuthn & SCIM 2.0](#4-enterprise-iam-saml-20-oidc-webauthn--scim-20)
5. [Data Residency, Localization & Cryptographic Shredding](#5-data-residency-localization--cryptographic-shredding)
6. [Formal Compliance Controls: SOC 2 Type II & PCI-DSS](#6-formal-compliance-controls-soc-2-type-ii--pci-dss)

---

## 1. Cross-Border Statutory Framework (US ↔ India & Global Corridors)

Cross-border fund movements are strictly governed by financial authorities (Federal Reserve / FinCEN in the US, Reserve Bank of India in India, FATF globally). 

```
[US Sender: Business]
         │
         ▼
[FinCEN Travel Rule Check] ───(Over $3,000 USD)───▶ Collect Verified Sender/Beneficiary Info
         │
         ▼
[Cross-Border Payment Processor (PayPal/Stripe/Bank)]
         │
         ▼
[Inward Remittance to India (AD Category-I Bank)]
         │
         ├──▶ Mandatory RBI Purpose Code Assigned (e.g. P0802)
         ├──▶ Foreign Exchange Management Act (FEMA) Declaration
         └──▶ Electronic Foreign Inward Remittance Certificate (e-FIRC) Issued
                 │
                 ▼
     [PayBridge Compliance Engine]
     Validates & Binds e-FIRC to Payment Truth Passport
```

---

## 2. RBI Purpose Codes & Automated e-FIRC Ingestion

### 2.1 Standardized Purpose Code Registry
Inward payments to Indian service providers, freelancers, and businesses must be tagged with valid Reserve Bank of India purpose codes.

| Code | Export Service Category | Requirements |
|---|---|---|
| **P0802** | Software Consultancy & Implementation | Software export declaration (SOFTEX / GST Invoicing) |
| **P0803** | Data Processing & Hosted Services | Service contract + GST LUT number |
| **P1007** | Advertising, Trade Fairs, Market Research | Commercial invoice match |
| **P1009** | Legal, Accounting & Business Management | Professional service agreement |

The backend strictly enforces purpose code assignment during case creation:
```go
type CreateCaseInput struct {
    TenantID        uuid.UUID       `json:"tenant_id"`
    Amount          decimal.Decimal `json:"amount"`
    Currency        string          `json:"currency"`
    SenderCountry   string          `json:"sender_country"`
    ReceiverCountry string          `json:"receiver_country"`
    PurposeCode     string          `json:"purpose_code"` // Mandatory if ReceiverCountry == "IN"
}
```

### 2.2 Automated e-FIRC Ingestion Pipeline
To establish undisputed legal proof that foreign currency arrived via authorized banking channels:
1. PayBridge monitors AD Bank or Provider API (e.g., PayPal / Stripe India e-FIRC webhooks).
2. The signed Electronic FIRC (containing Unique Transaction Reference, Inward Remittance Value, and Nostro Clearing Account) is downloaded.
3. The cryptographic SHA-256 hash of the e-FIRC PDF is embedded directly into the **Payment Truth Passport** as Evidence Tier E1.

---

## 3. Automated KYC/KYB & Sanctions/AML Screening

No commercial organization can execute or investigate cross-border transactions without deterministic counterparty verification.

```
                  [NEW USER / MERCHANT ONBOARDING]
                                 │
                 ┌───────────────┴───────────────┐
                 ▼                               ▼
       [INDIVIDUAL (KYC)]               [BUSINESS (KYB)]
       Persona / Trulioo                Trulioo / Middesk
       - Government ID Verification     - Articles of Incorporation
       - Liveness Biometrics            - Ultimate Beneficial Owners (UBO >25%)
       - Address Proof Verification     - Corporate Tax ID (EIN / PAN / GSTIN)
                 │                               │
                 └───────────────┬───────────────┘
                                 ▼
               [REAL-TIME SANCTIONS & AML SCREENING]
                     ComplyAdvantage / Dow Jones
       - OFAC Specially Designated Nationals (SDN)
       - Politically Exposed Persons (PEP Tier 1-3)
       - EU / UN / UK Financial Sanctions Lists
       - Adverse Media Automated Analysis
                                 │
                     ┌───────────┴───────────┐
                     ▼                       ▼
              [CLEAR / PASS]          [MATCH DETECTED]
              Account Active          - Case Frozen Immediately
                                      - Zero Money Movement
                                      - Suspicious Activity Report (SAR) Enqueued
```

### 3.1 FinCEN Travel Rule Implementation
For any cross-border commercial case exceeding **$3,000 USD** (or local currency equivalent):
- PayBridge automatically records and verifies:
  1. Sender's full legal name, physical registered address, and bank routing code.
  2. Beneficiary's full legal name, account number, and recipient financial institution.
- This data is isolated in the encrypted tenant vault (`encrypted_payloads`) and never exposed in public UI links or LLM prompts.

---

## 4. Enterprise IAM: SAML 2.0, OIDC, WebAuthn & SCIM 2.0

Enterprise deployments require integration with institutional Identity Providers (IdP). The demo cookie-based login (`argon2id`) is replaced with an enterprise authentication architecture.

```
               [ENTERPRISE CUSTOMER IT / OKTA / AZURE AD]
                                    │
           ┌────────────────────────┴────────────────────────┐
           ▼                                                 ▼
[SAML 2.0 / OIDC AUTH FLOW]                         [SCIM 2.0 DIRECTORY SYNC]
1. User hits PayBridge enterprise login             1. Okta pushes user/group changes
2. Redirected to customer SSO IdP                   2. Automatic JIT role mapping:
3. Enforces WebAuthn / Hardware FIDO2                  - Okta `Fintech-Admins` -> `merchant_admin`
4. Returns cryptographically signed assertion          - Okta `Ops-Team` -> `merchant_analyst`
5. PayBridge issues secure ephemeral session        3. Instant revocation on employee departure
```

### 4.1 SCIM 2.0 Provisioning Endpoints
PayBridge implements RFC 7644 SCIM endpoints:
- `POST /scim/v2/Users`
- `PUT /scim/v2/Users/{id}`
- `DELETE /scim/v2/Users/{id}` (Immediately terminates active sessions, invalidates case access tokens, and rotates tenant API keys).

---

## 5. Data Residency, Localization & Cryptographic Shredding

### 5.1 Jurisdictional Data Isolation
1. **India (RBI Circular DPSS.CO.OD No. 2785/06.08.005/2017-18)**:
   - All end-to-end transaction data, logs, and customer credentials relating to Indian beneficiaries must be stored on domestic data systems within India.
   - Secondary region or replica in AWS `ap-south-1` (Mumbai) maintains localized copies.
2. **European Union (GDPR Chapter V)**:
   - Data relating to EU participants remains hosted within the EU (Frankfurt / Dublin) or moves under Standard Contractual Clauses (SCCs).

### 5.2 Cryptographic Shredding on Case Deletion
When a participant exercises the statutory Right to Erasure (GDPR Art. 17 / DPDP Sec. 12):
1. Financial audit records (transaction IDs, ledger entries) are preserved for statutory retention periods (7 years) with personal identifiers anonymized.
2. The specific AES-256 Data Encryption Key (DEK) used to encrypt the user's PII is permanently purged from AWS KMS / Vault.
3. This renders all historical backups, archives, and unredacted logs mathematically unrecoverable (**Cryptographic Erasure**).

---

## 6. Formal Compliance Controls: SOC 2 Type II & PCI-DSS

### 6.1 SOC 2 Type II Trust Services Criteria Mapping

| TSC ID | Control Activity | Verification Mechanism |
|---|---|---|
| **CC6.1** | Logical access strictly managed via IAM | Enforced SSO + SCIM; quarterly automated access reviews |
| **CC6.6** | Perimeter boundaries protected against unauthorized access | Cloudflare Enterprise WAF + automated vulnerability scans |
| **CC7.1** | Infrastructure configuration monitored for changes | Terraform / GitOps state monitoring; drift alerts |
| **CC7.2** | System anomalies detected and alerted | PagerDuty alerting on unexpected state transitions |
| **CC8.1** | Software development lifecycle changes tested | Mandatory CI branch protection, 2 human reviews, zero bypass |

### 6.2 PCI-DSS Compliance Scope (v4.0)
- **SAQ-A / SAQ-A-EP Boundary**: PayBridge AI **never touches, transmits, or stores Primary Account Numbers (PAN) or CVVs**.
- All checkout interfaces utilize Hosted Fields or Provider Native iFrames (PayPal Hosted Fields, Stripe Elements).
- Content Security Policy (CSP) headers strictly forbid unauthorized third-party script injection on checkout domains.
