package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/paybridge/paybridge-ai/services/api/db/migrations"
	"github.com/paybridge/paybridge-ai/services/api/internal/store"
)

const defaultDevDatabaseURL = "postgres://paybridge_admin:paybridge_dev_secret@localhost:5432/paybridge_dev?sslmode=disable"

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	action := flag.String("action", "up", "Migration action: up, status")
	dbURLFlag := flag.String("db-url", "", "PostgreSQL database URL (defaults to DATABASE_URL env or local dev default)")
	flag.Parse()

	dbURL := *dbURLFlag
	if dbURL == "" {
		dbURL = os.Getenv("DATABASE_URL")
	}
	if dbURL == "" {
		dbURL = defaultDevDatabaseURL
	}

	st, err := store.NewStore(store.DefaultStoreConfig(dbURL))
	if err != nil {
		logger.Error("failed to initialize store configuration", "error", err)
		os.Exit(1)
	}
	defer st.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	if err := st.Ping(ctx); err != nil {
		logger.Error("database connection ping failed (PostgreSQL engine appears unavailable on port 5432)",
			"error", err,
			"target_url", st.String(),
		)
		os.Exit(1)
	}

	switch *action {
	case "up":
		if err := runMigrationsUp(ctx, st.DB(), logger); err != nil {
			logger.Error("migration up failed", "error", err)
			os.Exit(1)
		}
		logger.Info("all migrations applied successfully")
	case "status":
		if err := reportStatus(ctx, st.DB(), logger); err != nil {
			logger.Error("migration status check failed", "error", err)
			os.Exit(1)
		}
	default:
		logger.Error("unknown migration action", "action", *action)
		os.Exit(1)
	}
}

func runMigrationsUp(ctx context.Context, db *sql.DB, logger *slog.Logger) error {
	entries, err := fs.ReadDir(migrations.FS, ".")
	if err != nil {
		return fmt.Errorf("failed to read embedded migrations: %w", err)
	}

	var sqlFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			sqlFiles = append(sqlFiles, entry.Name())
		}
	}
	sort.Strings(sqlFiles)

	// Ensure schema_migrations table exists
	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version_id BIGINT PRIMARY KEY,
			is_applied BOOLEAN NOT NULL DEFAULT TRUE,
			tstamp TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp()
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to initialize schema_migrations: %w", err)
	}

	for i, filename := range sqlFiles {
		versionID := int64(i + 1)

		var applied bool
		err := db.QueryRowContext(ctx, "SELECT is_applied FROM schema_migrations WHERE version_id = $1", versionID).Scan(&applied)
		if err == nil && applied {
			logger.Info("skipping already applied migration", "file", filename, "version", versionID)
			continue
		}

		logger.Info("applying migration...", "file", filename, "version", versionID)
		content, err := fs.ReadFile(migrations.FS, filename)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", filename, err)
		}

		// Extract '-- +goose Up' segment
		sqlScript := extractUpSegment(string(content))

		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}

		if _, err := tx.ExecContext(ctx, sqlScript); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("migration %s failed: %w", filename, err)
		}

		if _, err := tx.ExecContext(ctx, "INSERT INTO schema_migrations (version_id, is_applied) VALUES ($1, TRUE)", versionID); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to record version %d: %w", versionID, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", filename, err)
		}

		logger.Info("applied migration successfully", "file", filename, "version", versionID)
	}

	return nil
}

func reportStatus(ctx context.Context, db *sql.DB, logger *slog.Logger) error {
	rows, err := db.QueryContext(ctx, "SELECT version_id, is_applied, tstamp FROM schema_migrations ORDER BY version_id ASC")
	if err != nil {
		return err
	}
	defer rows.Close()

	logger.Info("=== CURRENT MIGRATION STATUS ===")
	for rows.Next() {
		var ver int64
		var applied bool
		var tstamp time.Time
		if err := rows.Scan(&ver, &applied, &tstamp); err != nil {
			return err
		}
		logger.Info("version", "id", ver, "applied", applied, "timestamp", tstamp)
	}
	return rows.Err()
}

func extractUpSegment(content string) string {
	upMarker := "-- +goose Up"
	downMarker := "-- +goose Down"

	upIdx := strings.Index(content, upMarker)
	if upIdx == -1 {
		return content
	}

	remaining := content[upIdx+len(upMarker):]
	downIdx := strings.Index(remaining, downMarker)
	if downIdx == -1 {
		return remaining
	}

	return remaining[:downIdx]
}
