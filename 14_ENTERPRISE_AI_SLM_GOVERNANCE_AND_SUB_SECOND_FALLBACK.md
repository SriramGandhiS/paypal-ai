# 14 — Enterprise AI / SLM Governance, Sub-Second Fallback, and Model Risk Management

| Field | Value |
|---|---|
| Document | `14_ENTERPRISE_AI_SLM_GOVERNANCE_AND_SUB_SECOND_FALLBACK.md` |
| Version | 2.0.0-ENTERPRISE |
| Date | 10 November 2026 |
| Status | **SPECIFICATION — Enterprise Production Extension** |
| Audience | AI Engineers, LLM Ops Architects, Financial Model Risk Auditors, SREs |
| Extends | `01_PRODUCT_BLUEPRINT_AND_MASTER_VISION.md`, `03_SYSTEM_DESIGN_AND_ARCHITECTURE.md`, `07_AI_MULTI_AGENT_ARCHITECTURE_AND_PROMPT_SPEC.md` |
| Purpose | Establishes institutional AI SLAs (< 500 ms), hybrid SLM/LLM routing, semantic vector caching, advanced prompt-injection guardrails, and compliance with Federal Reserve SR 11-7 Model Risk Management guidelines. |

---

## Table of Contents
1. [Production AI Latency & SLA Guarantees](#1-production-ai-latency--sla-guarantees)
2. [Hybrid AI Routing Architecture (SLM + Hosted LLM)](#2-hybrid-ai-routing-architecture-slm--hosted-llm)
3. [Semantic Caching & Token Optimization Engine](#3-semantic-caching--token-optimization-engine)
4. [Sub-Second Circuit Breaker & Deterministic Fallback](#4-sub-second-circuit-breaker--deterministic-fallback)
5. [Enterprise Guardrails, Prompt-Injection & Canary Tokens](#5-enterprise-guardrails-prompt-injection--canary-tokens)
6. [Financial Model Risk Management (SR 11-7 Compliance)](#6-financial-model-risk-management-sr-11-7-compliance)
7. [Automated Shadow Evaluation & Drift Detection](#7-automated-shadow-evaluation--drift-detection)

---

## 1. Production AI Latency & SLA Guarantees

In an enterprise banking and payment ecosystem, waiting 15–30 seconds for a hosted multi-agent LLM chain is unacceptable for real-time customer support and automated checkout triage.

| Operation Tier | Target Latency | P99 Latency SLA | Model Tier | Fallback Trigger |
|---|---|---|---|---|
| **Tier 1: Instant Triage (A01, A05)** | < 150 ms | < 350 ms | Local SLM (Llama-3-8B / vLLM) | Timeout > 400 ms -> Go Deterministic |
| **Tier 2: Root Cause Hypothesis (A06)**| < 800 ms | < 1,500 ms | Quantized SLM + Semantic Cache | Cache Hit or Go Rule G01–G06 |
| **Tier 3: Dual-Side Synthesis (A08, A09)**| < 2,500 ms | < 4,500 ms | Frontier Hosted (Claude / Gemini) | Timeout > 5,000 ms -> Deterministic Templates |

---

## 2. Hybrid AI Routing Architecture (SLM + Hosted LLM)

Enterprise PayBridge implements a dual-tier execution engine:

```
                            [TRIGGER: CASE INVESTIGATION]
                                         │
                         [IMMUTABLE EVIDENCE SNAPSHOT]
                                         │
                 ┌───────────────────────┴───────────────────────┐
                 ▼                                               ▼
   [TIER 1 & 2: LOCAL SLM CLUSTER]                 [TIER 3: HOSTED FRONTIER LLM]
   (Private vLLM on NVIDIA L40S / A10G)           (Anthropic Claude / Google Gemini)
   - Zero Internet Egress for Raw Data             - Role-Projected, Anonymized Data Only
   - PII Scrubbed On-Premise                       - Dual-Side Narrative Generation
   - Runs: A01, A02, A05, A07                      - Complex Causal Counterfactuals
   - Latency: 80 - 180 ms                          - Runs: A06, A08, A09
                 │                                               │
                 └───────────────────────┬───────────────────────┘
                                         ▼
                             [SERVER-SIDE CLAIM VALIDATOR]
                      - Schema Enforcement (JSON-Schema)
                      - Cryptographic Citation Verification
                      - Cross-party DLP Privacy Gate
```

### 2.1 Private Small Language Model (SLM) Deployment
- **Model Baseline**: Quantized (AWQ/GPTQ) 8B parameter models (e.g., Llama-3-8B-Instruct or Mistral-7B-Fintech).
- **Inference Engine**: vLLM cluster running inside the private Kubernetes VPC with Triton Inference Server.
- **Benefits**: Fixed infrastructure cost, deterministic latency (< 150 ms), zero risk of third-party model vendor downtime, and zero cloud data leakage.

---

## 3. Semantic Caching & Token Optimization Engine

Repeated payment failure patterns (e.g., PayPal Sandbox capture timeouts, recurring gateway 504 errors) exhibit identical evidence topology.

```
[Evidence Snapshot] ──▶ [Embedding Generator (BGE-Small)] ──▶ [Redis Vector Store]
                                                                     │
                                         ┌───────────────────────────┴───────────────────────────┐
                                         ▼ (Cosine Similarity >= 0.96)                           ▼ (Cosine Similarity < 0.96)
                                   [CACHE HIT]                                             [CACHE MISS]
                       Extract Cached Explanation Template                       Forward to SLM / LLM Cluster
                       Hydrate with Fresh Case IDs & Timestamps                  Index New Result upon Validation
                       Latency: 12 ms | Cost: $0.00                              Latency: 1,800 ms
```

### 3.1 Invalidation & Freshness Invariant
The semantic cache key includes `(tenant_id, error_category, provider_state, application_state)`. Any change to the case's evidence tier or state vector instantly bypasses the cache.

---

## 4. Sub-Second Circuit Breaker & Deterministic Fallback

Under peak load, DDoS, or third-party LLM outages, PayBridge AI must continue operating seamlessly without degradation of financial truth:

```go
package llmgw

import (
    "context"
    "time"
    "github.com/sony/gobreaker"
)

type ProductionLLMGateway struct {
    slmClient       ModelClient
    hostedClient    ModelClient
    circuitBreaker  *gobreaker.CircuitBreaker
    deterministic   DeterministicEngine
}

func (gw *ProductionLLMGateway) GenerateSafeExplanation(ctx context.Context, snap Snapshot) (Explanation, error) {
    ctxTimeout, cancel := context.WithTimeout(ctx, 3500*time.Millisecond)
    defer cancel()

    // 1. Attempt Semantic Cache
    if hit, cached := gw.checkSemanticCache(ctx, snap); hit {
        return cached, nil
    }

    // 2. Execute via Circuit Breaker
    result, err := gw.circuitBreaker.Execute(func() (interface{}, error) {
        if snap.Complexity == Low {
            return gw.slmClient.Generate(ctxTimeout, snap)
        }
        return gw.hostedClient.Generate(ctxTimeout, snap)
    })

    if err != nil {
        // Instant Fallback to Deterministic Go Rule Engine (Demo E Pattern)
        // SLA: < 5 milliseconds
        return gw.deterministic.GenerateDeterministicPassport(snap), nil
    }

    return result.(Explanation), nil
}
```

---

## 5. Enterprise Guardrails, Prompt-Injection & Canary Tokens

Financial systems are prime targets for prompt-injection attacks aiming to induce false confirmation of funds or data exfiltration.

```
                           [UNTRUSTED USER / INVOICE INPUT]
                                         │
                                         ▼
                       [PASS 1: TOKEN PURIFICATION / NORMALIZATION]
                       - Unicode NFKC Normalization (Strips Homoglyphs)
                       - Zero-Width Character Stripping
                       - Markdown & Control Sequence Escaping
                                         │
                                         ▼
                       [PASS 2: NE-MO & LLAMA GUARD 3 FILTER]
                       - Jailbreak Pattern Recognition
                       - Financial Advice Manipulation Classifier
                       - Instruction Boundary Escape Detection
                                         │
                                         ▼
                       [PASS 3: SECRET CANARY TOKEN INJECTION]
                       System injects hidden random UUID into system block:
                       "<!-- canary_token: 7f8a3c9b-12d4-4a8e-b611-9e8a7c2b3d4e -->"
                                         │
                                         ▼
                                   [MODEL EXECUTION]
                                         │
                                         ▼
                       [PASS 4: OUTPUT DATA LOSS PREVENTION (DLP)]
                       - Assert Canary Token is ABSENT from Output (If present -> EXFILTRATION DETECTED)
                       - Regex scan for PAN (Credit Cards), SSN, Tax IDs, Private Auth Tokens
                       - Claim Citation Enforcement (Every claim must cite Evidence ID)
```

---

## 6. Financial Model Risk Management (SR 11-7 Compliance)

In the United States, the Federal Reserve and OCC mandate strict compliance under **Supervisory Guidance on Model Risk Management (SR 11-7 / OCC 2011-12)** for any automated model informing financial decisions or communications.

### 6.1 Compliance Requirements & Controls

| SR 11-7 Requirement | PayBridge Implementation Control |
|---|---|
| **Sound Model Development** | Clear separation between deterministic financial state (Reducer) and informational explanation (LLM). |
| **Model Inventory & Versioning** | Every model output stores immutable records: `model_provider`, `model_name`, `model_weights_hash`, `prompt_template_version`, and `temperature`. |
| **Rigorous Conceptual Validation** | Evaluation suite (T01–T10) run on every release; automated verification that AI outputs never initiate transactions. |
| **Ongoing Monitoring & Drift Detection**| Daily automated scoring of claim accuracy against provider reconciliation logs. |
| **Independent Model Review** | Quarterly third-party red-teaming and prompt-injection resistance auditing. |

---

## 7. Automated Shadow Evaluation & Drift Detection

In production, models can degrade over time due to upstream provider API changes, shifting prompt distributions, or model updates.

### 7.1 Shadow Pipeline
1. **Parallel Execution**: 5% of all live production investigations are asynchronously evaluated across an alternate secondary model (e.g., comparing Claude 3.5 Sonnet against Gemini 1.5 Pro).
2. **Discrepancy Scoring**:
   - Citation consistency delta.
   - Tone divergence (calm vs alarmist score).
   - Hallucination index (unsupported claims count).
3. **Automated Alerting**: If the discrepancy score between models exceeds **2.0%**, an automated P3 ticket is created in the engineering review queue with full anonymized evidence snapshots.
