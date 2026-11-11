#!/bin/bash
set -eo pipefail

echo "=== PayBridge Enterprise Database Initializer ==="

POSTGRES_USER="${POSTGRES_USER:-paybridge_admin}"
POSTGRES_DB="${POSTGRES_DB:-paybridge}"

# Check for official Goose CLI binary
if ! command -v goose &> /dev/null; then
    echo "Checking for goose binary in /usr/local/bin..."
    if command -v apk &> /dev/null; then
        apk add --no-cache curl ca-certificates bash || true
    fi
    if command -v curl &> /dev/null; then
        echo "Attempting to install Goose CLI..."
        curl -fsSL https://raw.githubusercontent.com/pressly/goose/master/install.sh | sh -s -- -b /usr/local/bin || true
    fi
fi

if command -v goose &> /dev/null; then
    echo "Running migrations with official Goose CLI..."
    goose -dir /migrations postgres "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@127.0.0.1:5432/${POSTGRES_DB}?sslmode=disable" up
else
    echo "Goose CLI not found or offline. Executing verified deterministic Up-migration runner..."
    
    # Initialize Goose tracking table to match Goose CLI semantics
    psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -v ON_ERROR_STOP=1 <<-EOSQL
        CREATE TABLE IF NOT EXISTS goose_db_version (
            id SERIAL PRIMARY KEY,
            version_id BIGINT NOT NULL,
            is_applied BOOLEAN NOT NULL,
            tstamp TIMESTAMP DEFAULT now()
        );
EOSQL

    for f in $(ls -1 /migrations/*.sql | sort -V); do
        filename=$(basename "$f")
        # Extract numerical version prefix (e.g., 0001 -> 1)
        version_id=$(echo "$filename" | grep -oE '^[0-9]+' | sed 's/^0*//')
        
        # Check if already applied
        applied=$(psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -t -A -c \
            "SELECT COUNT(*) FROM goose_db_version WHERE version_id = ${version_id} AND is_applied = true;")

        if [ "$applied" -eq "0" ]; then
            echo "Applying migration ${filename} (Up only)..."
            # Stream directly into psql to prevent bash variable expansion of $$ in PL/pgSQL
            {
                echo "BEGIN;"
                awk '/-- \+goose Up/{flag=1;next}/-- \+goose Down/{flag=0}flag' "$f"
                echo ""
                echo "INSERT INTO goose_db_version (version_id, is_applied) VALUES (${version_id}, true);"
                echo "COMMIT;"
            } | psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -v ON_ERROR_STOP=1
            echo "Successfully applied ${filename}"
        else
            echo "Skipping ${filename} (already applied)"
        fi
    done
fi

echo "=== PayBridge Migrations Successfully Initialized ==="
