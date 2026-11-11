# PayBridge AI — Enterprise Release Gates Automation (PRG-01 to PRG-10)
# Owner: Dinny Paul Navis

$rootDir = Split-Path -Parent $PSScriptRoot
Set-Location "$rootDir\services\api"

Write-Host "==========================================================" -ForegroundColor Cyan
Write-Host "PayBridge AI — Production Release Gates (PRG-01 to PRG-10)" -ForegroundColor Cyan
Write-Host "Working Directory: $(Get-Location)" -ForegroundColor Cyan
Write-Host "==========================================================" -ForegroundColor Cyan

$passed = 0
$total = 10

function Run-Gate {
    param(
        [string]$GateName,
        [string]$Command
    )
    Write-Host "`n[Evaluating $GateName]..." -ForegroundColor Yellow
    $start = Get-Date
    Invoke-Expression $Command
    if ($LASTEXITCODE -eq 0) {
        $dur = ((Get-Date) - $start).TotalSeconds
        Write-Host ">>> $GateName : PASSED in $([math]::Round($dur, 2))s" -ForegroundColor Green
        return $true
    } else {
        Write-Host ">>> $GateName : FAILED (Exit Code: $LASTEXITCODE)" -ForegroundColor Red
        return $false
    }
}

# PRG-01: Static Analysis, Dependency and Secrets
$r1 = Run-Gate -GateName "PRG-01: Static Analysis and Code Safety" -Command "go vet ./..."
if ($r1) { $passed++ }

# PRG-02: Multi-Rail Payment Gateway Adapters
$r2 = Run-Gate -GateName "PRG-02: Multi-Rail Adapters and Signatures" -Command "go test -v ./internal/multirail"
if ($r2) { $passed++ }

# PRG-03: 10,000 TPS Kafka Benchmark & Partition Ordering
$r3 = Run-Gate -GateName "PRG-03: 10k TPS Kafka Benchmark" -Command "go test -v -run TestKafka_10000EventsThroughputLoad ./internal/kafka"
if ($r3) { $passed++ }

# PRG-04: Double-Entry SQL Ledger Zero-Sum Invariant
$r4 = Run-Gate -GateName "PRG-04: Zero-Sum Ledger Postings" -Command "go test -v ./internal/ledger"
if ($r4) { $passed++ }

# PRG-05: KMS Multi-Version Envelope Encryption & Tamper Detection
$r5 = Run-Gate -GateName "PRG-05: KMS Envelope Encryption" -Command "go test -v ./internal/kms"
if ($r5) { $passed++ }

# PRG-06: Authenticated Bank Evidence Settlement Matcher
$r6 = Run-Gate -GateName "PRG-06: Bank Settlement Verification" -Command "go test -v ./internal/settlement"
if ($r6) { $passed++ }

# PRG-07: Statutory RBI Purpose Codes and Sanctions Screening
$r7 = Run-Gate -GateName "PRG-07: Compliance and Sanctions Screening" -Command "go test -v ./internal/compliance"
if ($r7) { $passed++ }

# PRG-08: Sub-Second SLM P99 Latency & Semantic Cache
$r8 = Run-Gate -GateName "PRG-08: SLM Sub-Second P99 Benchmark" -Command "go test -v -run TestSLM_500ConcurrentTriageBenchmark ./internal/slm"
if ($r8) { $passed++ }

# PRG-09: Enterprise IAM and SCIM 2.0 Provisioning
$r9 = Run-Gate -GateName "PRG-09: Enterprise IAM and SCIM 2.0" -Command "go test -v ./internal/iam"
if ($r9) { $passed++ }

# PRG-10: Chaos Engineering Invariant Scenarios (CHAOS-01 to CHAOS-12)
$r10 = Run-Gate -GateName "PRG-10: Chaos Engineering Invariants" -Command "go test -v ./tests/chaos"
if ($r10) { $passed++ }

Write-Host "`n==========================================================" -ForegroundColor Cyan
Write-Host "RELEASE GATE RESULTS: $passed / $total PASSED" -ForegroundColor Cyan
Write-Host "==========================================================" -ForegroundColor Cyan

if ($passed -eq $total) {
    Write-Host "ALL 10 ENTERPRISE RELEASE GATES MET." -ForegroundColor Green
    exit 0
} else {
    Write-Host "RELEASE GATES FAILED. DEPLOYMENT BLOCKED." -ForegroundColor Red
    exit 1
}
