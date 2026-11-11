# Environment Variable Reference Guide

**Project:** PayBridge AI  
**Author / Module Owner:** Dinny Paul Navis (Enterprise Infrastructure, Banking & Compliance)  
**Recipient / Platform Lead:** Sriram  
**Reference Document:** `docs/implementation/dinny/ENVIRONMENT_VARIABLE_REFERENCE.md`  

---

## 1. Overview & Security Policy

All environment-specific configuration values in the PayBridge AI enterprise stack are externalized. 

### Security Rules:
1. **Never Commit Secrets:** Real API keys, private passwords, and provider access tokens must never be checked into version control.
2. **Mandatory Fail-Fast:** If a required variable is missing, the service must fail during initialization rather than falling back to insecure defaults.
3. **Loopback & Allowlist Boundaries:** Kafka Connect variables are strictly constrained by regex allowlisting (`CONNECT_CONFIG_PROVIDERS_ENV_PARAM_ALLOWLIST_PATTERN: "^(POSTGRES_USER|POSTGRES_PASSWORD|POSTGRES_DB)$"`).

---

## 2. Complete Environment Variable Inventory

### A. Database / PostgreSQL Persistence
| Variable | Description | Used By | Required? | Who Supplies | Default / Fallback |
|---|---|---|---|---|---|
| `POSTGRES_HOST` | Hostname or service name for PostgreSQL. | Backend, Debezium, Worker | Yes | Sriram / Infra | `localhost` |
| `POSTGRES_PORT` | Port for PostgreSQL. | Backend, Debezium, Worker | Yes | Sriram / Infra | `5432` |
| `POSTGRES_DB` | PostgreSQL database name. | All backend services | Yes | Sriram / Infra | `paybridge` |
| `POSTGRES_USER` | Admin/Replication user for PostgreSQL. | PostgreSQL, Debezium | Yes | Sriram / Infra | `paybridge_admin` |
| `POSTGRES_PASSWORD` | Strong password for PostgreSQL. | PostgreSQL, Debezium, API | **YES (SECRET)** | Sriram / Infra | *None (Fails if unset)* |
| `DATABASE_URL` | Standard PostgreSQL connection DSN. | `services/api` | **YES (SECRET)** | Sriram / Infra | *None (Fails if unset)* |

---

### B. Redis Semantic Cache & Rate Limiting
| Variable | Description | Used By | Required? | Who Supplies | Default / Fallback |
|---|---|---|---|---|---|
| `REDIS_HOST` | Hostname or service name for Redis. | `services/api`, SLM Cache | Yes | Sriram / Infra | `localhost` |
| `REDIS_PORT` | Port for Redis instance. | `services/api`, SLM Cache | Yes | Sriram / Infra | `6379` |
| `REDIS_PASSWORD` | Password for Redis `AUTH`. | Redis container, API | **YES (SECRET)** | Sriram / Infra | *None (Fails if unset)* |
| `REDIS_URL` | Redis connection URL with auth. | `services/api/internal/slm` | **YES (SECRET)** | Sriram / Infra | *None (Fails if unset)* |
| `REDIS_CACHE_ENABLED`| Flag to enable semantic embedding caching. | SLM Module | No | Sriram / Operator | `true` |

---

### C. D1: Distributed Event Streaming (Apache Kafka & Debezium)
| Variable | Description | Used By | Required? | Who Supplies | Default / Fallback |
|---|---|---|---|---|---|
| `KAFKA_BOOTSTRAP_SERVERS`| Comma-separated list of Kafka broker endpoints. | `services/api/internal/kafka` | Yes | Sriram / Infra | `localhost:9092` |
| `KAFKA_BROKERS` | Alias for bootstrap servers. | Kafka consumer/producer | Yes | Sriram / Infra | `localhost:9092` |
| `KAFKA_TOPIC` | Canonical payment events topic name. | Kafka bus & Debezium | Yes | Dinny (Code) | `paybridge.payment.events.v1` |
| `KAFKA_DLQ_TOPIC` | Dead-letter queue topic for poison messages. | Kafka consumer | Yes | Dinny (Code) | `paybridge.deadletter.v1` |
| `KAFKA_GROUP_ID` | Consumer group ID for payment processors. | Background workers | Yes | Sriram / Ops | `paybridge-core-worker` |
| `KAFKA_CLIENT_ID` | Identifier for producer telemetry. | Kafka bus | No | Sriram / Ops | `paybridge-core-worker` |
| `KAFKA_PARTITIONS` | Physical partition count for topic ordering. | `kafka-init` service | Yes | Dinny (Spec) | `32` |

---

### D. D2: KMS Envelope Encryption (AWS KMS / HashiCorp Vault)
| Variable | Description | Used By | Required? | Who Supplies | Default / Fallback |
|---|---|---|---|---|---|
| `ENCRYPTION_KEY` | 32-byte Base64-encoded seed key for local AES-256-GCM. | `services/api/internal/kms` | **YES (SECRET)** | Sriram / Security | Generates random in test |
| `ENCRYPTION_KEY_ID` | Active Key Encryption Key (KEK) version identifier. | KMS envelope cipher | Yes | Sriram / Security | `2026-Q1` |
| `AWS_REGION` | AWS region for KMS Key. | `services/api/internal/kms` | If using KMS | Sriram / AWS Admin | `us-east-1` |
| `AWS_ACCESS_KEY_ID` | AWS IAM credentials (when not using IAM roles). | AWS KMS SDK v2 | Conditional | Sriram / AWS Admin | IAM Role / Environment |
| `AWS_SECRET_ACCESS_KEY`| AWS IAM secret key. | AWS KMS SDK v2 | Conditional | Sriram / AWS Admin | IAM Role / Environment |
| `AWS_KMS_KEY_ID` | Customer Master Key (CMK) ARN. | AWS KMS SDK v2 | If using KMS | Sriram / AWS Admin | *None (Uses local fallback)* |
| `VAULT_ADDR` | HashiCorp Vault HTTP URL. | `services/api/internal/kms` | If using Vault | Sriram / Infra | `http://localhost:8200` |
| `VAULT_TOKEN` | HashiCorp Vault authentication token. | Vault Client | Conditional | Sriram / Security | *None (Uses local fallback)* |
| `VAULT_TRANSIT_KEY` | Key name inside Vault Transit secrets engine. | Vault Client | If using Vault | Sriram / Security | `paybridge-evidence` |

---

### E. D3: Multi-Rail Payment Gateways (Stripe & Wise)
| Variable | Description | Used By | Required? | Who Supplies | Default / Fallback |
|---|---|---|---|---|---|
| `PAYPAL_MODE` | Payment engine mode (`real` or `fake`). | PayPal client | Yes | Sriram / Core | `fake` |
| `PAYPAL_ENV` | Target PayPal environment (`sandbox` or `live`). | PayPal client | Yes | Sriram / Core | `sandbox` |
| `PAYPAL_CLIENT_ID` | PayPal REST API client ID. | PayPal adapter | Yes (Real mode) | Sriram | Placeholder |
| `PAYPAL_CLIENT_SECRET` | PayPal REST API secret. | PayPal adapter | **YES (SECRET)** | Sriram | Placeholder |
| `STRIPE_SECRET_KEY` | Stripe Sandbox secret key (`sk_test_...`). | `services/api/internal/multirail` | **YES (SECRET)** | Sriram / Stripe | Contract fallback in tests |
| `STRIPE_WEBHOOK_SECRET`| Stripe Webhook signing secret (`whsec_...`). | Webhook HMAC validator | **YES (SECRET)** | Sriram / Stripe | Strict validation enabled |
| `WISE_API_TOKEN` | Wise sandbox personal/business API token. | Wise adapter | **YES (SECRET)** | Sriram / Wise | Contract fallback in tests |
| `WISE_WEBHOOK_SECRET` | Wise webhook signature token. | Wise webhook listener | **YES (SECRET)** | Sriram / Wise | *None* |
| `WISE_PROFILE_ID` | Wise business profile identifier. | Wise wire transfers | Yes | Sriram / Wise | Placeholder |

---

### F. D4: Open Banking & Settlement (Plaid & Setu)
| Variable | Description | Used By | Required? | Who Supplies | Default / Fallback |
|---|---|---|---|---|---|
| `PLAID_CLIENT_ID` | Plaid developer account client ID. | `services/api/internal/settlement` | Yes | Sriram / Plaid | Placeholder |
| `PLAID_SECRET` | Plaid sandbox secret. | Plaid API client | **YES (SECRET)** | Sriram / Plaid | Contract fallback in tests |
| `PLAID_ENV` | Plaid API environment. | Plaid client | Yes | Sriram / Plaid | `sandbox` |
| `SETU_CLIENT_ID` | Sahamati / Setu Account Aggregator client ID. | Setu AA client | Yes | Sriram / Setu | Placeholder |
| `SETU_CLIENT_SECRET` | Sahamati / Setu Account Aggregator secret. | Setu AA client | **YES (SECRET)** | Sriram / Setu | Contract fallback in tests |
| `SETU_PRODUCT_INSTANCE_ID`| Setu Consent product instance UUID. | Setu AA client | Yes | Sriram / Setu | Placeholder |

---

### G. D6: Regulatory Compliance & KYC/Sanctions
| Variable | Description | Used By | Required? | Who Supplies | Default / Fallback |
|---|---|---|---|---|---|
| `PERSONA_API_KEY` | Persona KYC sandbox API key. | `services/api/internal/compliance` | **YES (SECRET)** | Sriram / Persona | Mock gate in unit tests |
| `PERSONA_WEBHOOK_SECRET`| Persona webhook HMAC secret. | Compliance Webhooks | **YES (SECRET)** | Sriram / Persona | Placeholder |
| `COMPLYADVANTAGE_API_KEY`| ComplyAdvantage sanctions & PEP API key. | Sanctions screening engine | **YES (SECRET)** | Sriram / ComplyAdv | Mock gate in unit tests |

---

### H. D7: Enterprise IAM & Private SLM
| Variable | Description | Used By | Required? | Who Supplies | Default / Fallback |
|---|---|---|---|---|---|
| `SCIM_BEARER_TOKEN` | Secret bearer token authorizing SCIM 2.0 requests. | `services/api/internal/iam` | **YES (SECRET)** | Sriram / Okta | Checked at `/scim/v2/Users` |
| `OKTA_DOMAIN` | Okta tenant domain URL. | SAML / OIDC provider | Yes | Sriram / Okta Admin| `https://paybridge.okta.com` |
| `OKTA_ISSUER` | OIDC token issuer URL. | JWT validation | Yes | Sriram / Okta Admin| `https://paybridge.okta.com/oauth2/default` |
| `OKTA_CLIENT_ID` | Okta OIDC application client ID. | OIDC Flow | Yes | Sriram / Okta Admin| Placeholder |
| `OKTA_CLIENT_SECRET` | Okta OIDC application client secret. | OIDC Flow | **YES (SECRET)** | Sriram / Okta Admin| Placeholder |
| `SAML_ENTITY_ID` | SP Entity ID for SAML assertion consumer. | SAML 2.0 module | Yes | Dinny (Code) | `https://paybridge.internal/saml/metadata` |
| `VLLM_BASE_URL` | HTTP endpoint for private vLLM SLM inference. | `services/api/internal/slm` | Yes | Sriram / GPU Infra | `http://localhost:8000/v1` |
| `SLM_MODEL_NAME` | Model identifier loaded in vLLM. | SLM client | Yes | Sriram / GPU Infra | `meta-llama/Meta-Llama-3-8B-Instruct` |

---

## 3. Operational Handoff Responsibility

- **Dinny Paul Navis (Code Author):** Implemented all consumer, producer, client, and cryptographic adapters to read from these variables with zero plaintext fallbacks. Provided mock/unit contract stubs so test suites pass without network dependencies.
- **Sriram (Platform Lead):** Responsible for provisioning cloud infrastructure, generating production/sandbox credentials, populating `.env`, and configuring production secrets in AWS Secrets Manager or HashiCorp Vault.

