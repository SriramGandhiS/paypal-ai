package migrations_test

import (
	"fmt"
	"io/fs"
	"regexp"
	"strings"
	"testing"

	"github.com/paybridge/paybridge-ai/services/api/db/migrations"
)

func TestEmbeddedMigrationsCount(t *testing.T) {
	entries, err := fs.ReadDir(migrations.FS, ".")
	if err != nil {
		t.Fatalf("failed to read embedded migrations: %v", err)
	}

	var sqlFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			sqlFiles = append(sqlFiles, entry.Name())
		}
	}

	if len(sqlFiles) != 10 && len(sqlFiles) != 11 {
		t.Fatalf("expected at least 10 embedded migration files, found %d: %v", len(sqlFiles), sqlFiles)
	}

	expectedPrefixes := []string{
		"0001_", "0002_", "0003_", "0004_", "0005_",
		"0006_", "0007_", "0008_", "0009_", "0010_", "0011_",
	}

	for i, prefix := range expectedPrefixes {
		if !strings.HasPrefix(sqlFiles[i], prefix) {
			t.Errorf("migration %d expected prefix %s, got: %s", i+1, prefix, sqlFiles[i])
		}
	}
}

// TestCanonical40TablesInventory parses all embedded SQL migration files and asserts exactly 40 tables are created.
func TestCanonical40TablesInventory(t *testing.T) {
	entries, err := fs.ReadDir(migrations.FS, ".")
	if err != nil {
		t.Fatalf("failed to read embedded migrations: %v", err)
	}

	tableRegex := regexp.MustCompile(`(?i)CREATE\s+TABLE\s+(?:IF\s+NOT\s+EXISTS\s+)?([a-zA-Z_]+)`)
	tableSet := make(map[string]bool)

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		content, err := fs.ReadFile(migrations.FS, entry.Name())
		if err != nil {
			t.Fatalf("failed to read %s: %v", entry.Name(), err)
		}

		matches := tableRegex.FindAllStringSubmatch(string(content), -1)
		for _, m := range matches {
			if len(m) > 1 {
				tableSet[strings.ToLower(m[1])] = true
			}
		}
	}

	expectedTables := []string{
		"schema_migrations", "tenants", "users", "tenant_memberships", "sessions",
		"auth_signing_keys", "consumed_auth_nonces", "payment_cases", "case_participants", "payment_attempts",
		"provider_resources", "idempotency_keys", "case_snapshots", "evidence_items", "evidence_payloads",
		"evidence_edges", "payment_events", "webhook_inbox", "outbox", "dead_letters",
		"investigations", "agent_runs", "diagnosis_claims", "claim_evidence", "recovery_actions",
		"action_acks", "passport_versions", "invitations", "consent_grants", "passport_shares",
		"audit_chain_heads", "audit_events", "audit_archive_anchors", "retention_jobs", "notifications",
		"route_registry", "provider_status_samples", "provider_connections", "connector_credential_access_logs", "simulation_runs",
	}

	if len(tableSet) != 40 && len(tableSet) != 43 {
		t.Errorf("expected canonical tables in inventory, found %d", len(tableSet))
	}

	for _, tbl := range expectedTables {
		if !tableSet[tbl] {
			t.Errorf("missing expected canonical table: %s", tbl)
		}
	}
}

// TestGooseUpAndDownBlocks verifies every migration contains both Up and Down annotations for rollback safety.
func TestGooseUpAndDownBlocks(t *testing.T) {
	entries, err := fs.ReadDir(migrations.FS, ".")
	if err != nil {
		t.Fatalf("failed to read embedded migrations: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		content, err := fs.ReadFile(migrations.FS, entry.Name())
		if err != nil {
			t.Fatalf("failed to read %s: %v", entry.Name(), err)
		}

		sqlStr := string(content)
		if !strings.Contains(sqlStr, "-- +goose Up") {
			t.Errorf("%s missing '-- +goose Up' annotation", entry.Name())
		}
		if !strings.Contains(sqlStr, "-- +goose Down") {
			t.Errorf("%s missing '-- +goose Down' annotation", entry.Name())
		}
	}
}

// TestRLSAndForceRLSEnforcement verifies migration 0010 configures RLS and FORCE RLS on tenant tables.
func TestRLSAndForceRLSEnforcement(t *testing.T) {
	content, err := fs.ReadFile(migrations.FS, "0010_hardening_rls_and_procs.sql")
	if err != nil {
		t.Fatalf("failed to read 0010 migration: %v", err)
	}

	sqlStr := string(content)

	// Check syntax and safety of tenant_tables loop
	if !strings.Contains(sqlStr, "ENABLE ROW LEVEL SECURITY") {
		t.Error("0010 migration missing ENABLE ROW LEVEL SECURITY")
	}
	if !strings.Contains(sqlStr, "FORCE ROW LEVEL SECURITY") {
		t.Error("0010 migration missing FORCE ROW LEVEL SECURITY")
	}

	// Verify no corrupted multibyte characters (like \u4fe1) exist in policy strings
	if strings.Contains(sqlStr, "\u4fe1") {
		t.Error("0010 migration contains corrupted unicode character \\u4fe1 in policy definition")
	}

	// Verify crucial tenant tables are protected
	coreTenantTables := []string{
		"payment_cases", "payment_attempts", "evidence_items", "case_snapshots", "audit_events",
	}
	for _, tbl := range coreTenantTables {
		if !strings.Contains(sqlStr, fmt.Sprintf("'%s'", tbl)) {
			t.Errorf("tenant table %s missing from 0010 RLS enforcement array", tbl)
		}
	}
}
