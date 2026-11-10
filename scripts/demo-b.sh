#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
API_DIR="$SCRIPT_DIR/../services/api"

echo "================================================================="
echo "PayBridge AI — Golden Demo B Automated Scenario Runner"
echo "Executing Go Runtime Lifecycle: cmd/demob/main.go"
echo "================================================================="

cd "$API_DIR"
go run cmd/demob/main.go
