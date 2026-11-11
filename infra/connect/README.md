# Hardened Kafka Connect Image & CVE Mitigation Guide

**Artifact:** `infra/connect/Dockerfile.patched`  
**SBOM:** `infra/connect/sbom.json`  
**Gateway Proxy:** `infra/connect/nginx-auth-proxy.conf`  
**Owner:** Dinny Paul Navis (Enterprise Infrastructure)  

---

## 1. Overview & Problem Statement

This directory provides the production-grade, reproducible custom build and administrative proxy configuration for Apache Kafka Connect within the PayBridge AI ecosystem.

### Vulnerabilities Addressed:
1. **CVE-2026-116593 (Kafka Connect REST Secret Leak)**:
   - Affects Kafka Connect versions prior to `4.2.2`, between `4.3.0` and `4.3.2`, and prior to `4.4.0`.
   - The `/connector-plugins/{connector-type}/config/validate` endpoint evaluates configuration providers and returns resolved plaintext secrets (such as `${env:POSTGRES_PASSWORD}`) in the JSON response body.
   - **Fix:** Built directly on **Apache Kafka Connect 4.4.0**, which sanitizes validation responses and prevents plaintext secret disclosure.
2. **CVE-2024-31141 (ConfigProvider Arbitrary Read)**:
   - Automatic `ConfigProvider` plugins allow unauthorized file or environment variable reads.
   - **Fix:** Handled in Kafka 4.4.0 alongside strict regex parameter allowlisting (`CONNECT_CONFIG_PROVIDERS_ENV_PARAM_ALLOWLIST_PATTERN: "^(POSTGRES_USER|POSTGRES_PASSWORD|POSTGRES_DB)$"`).

---

## 2. Image Build Instructions

To build the patched container image locally:

```bash
cd infra/connect
docker build -t paybridge/kafka-connect-hardened:4.4.0-debezium-2.5.4 -f Dockerfile.patched .
```

### Verification Commands:
Verify the Kafka runtime version inside the container:
```bash
docker run --rm paybridge/kafka-connect-hardened:4.4.0-debezium-2.5.4 kafka-run-class.sh org.apache.kafka.common.utils.AppInfoParser
```
*Expected Output:* `4.4.0`

Verify Debezium connector plugin discovery:
```bash
docker run --rm paybridge/kafka-connect-hardened:4.4.0-debezium-2.5.4 ls -l /opt/kafka/plugins/debezium-postgres
```
*Expected Output:* `debezium-connector-postgres-2.5.4.Final.jar`, `debezium-core-2.5.4.Final.jar`, and dependencies.

---

## 3. Maintenance & Dependency Management Policy

1. **No Manual In-Place Library Overwrites:** Never manually replace Kafka core JARs inside an existing upstream Debezium image. Rebuild from source or from the pinned `Dockerfile.patched`.
2. **Pinning & Upstream Certification:** When Debezium officially releases a `.Final` container image built on Kafka 4.x (e.g. `debezium/connect:3.x` with Kafka 4.x), compare SHA-256 digests and migrate smoothly.
3. **Network Boundary Enforcement:** Even with a patched runtime, always deploy behind an authenticated gateway (`nginx-auth-proxy.conf`) and restrict host port bindings to localhost loopback (`127.0.0.1:8083`).
