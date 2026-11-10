package migrations_test

import (
	"io/fs"
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

	if len(sqlFiles) != 10 {
		t.Fatalf("expected exactly 10 embedded migration files (0001-0010), found %d: %v", len(sqlFiles), sqlFiles)
	}

	expectedPrefixes := []string{
		"0001_", "0002_", "0003_", "0004_", "0005_",
		"0006_", "0007_", "0008_", "0009_", "0010_",
	}

	for i, prefix := range expectedPrefixes {
		if !strings.HasPrefix(sqlFiles[i], prefix) {
			t.Errorf("migration %d expected prefix %s, got: %s", i+1, prefix, sqlFiles[i])
		}
	}
}
