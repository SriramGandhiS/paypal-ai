package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"time"
)

// Common database errors
var (
	ErrNilDatabaseURL = errors.New("database URL cannot be empty")
	ErrStoreClosed    = errors.New("store connection is closed")
)

// StoreConfig controls the database connection pool settings.
type StoreConfig struct {
	DatabaseURL     string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// DefaultStoreConfig provides institutional-grade connection pool defaults.
func DefaultStoreConfig(databaseURL string) StoreConfig {
	return StoreConfig{
		DatabaseURL:     databaseURL,
		MaxOpenConns:    25,
		MaxIdleConns:    10,
		ConnMaxLifetime: 15 * time.Minute,
		ConnMaxIdleTime: 5 * time.Minute,
	}
}

// Store manages the PostgreSQL database connection pool and execution context.
type Store struct {
	db      *sql.DB
	driver  string
	safeURL string
}

// String returns a redacted representation of the database URL.
func (s *Store) String() string {
	if s == nil {
		return "<nil>"
	}
	return s.safeURL
}

// NewStore initializes a new Store instance from configuration.
func NewStore(cfg StoreConfig) (*Store, error) {
	if cfg.DatabaseURL == "" {
		return nil, ErrNilDatabaseURL
	}

	// Validate URL parseability without leaking credentials
	u, err := url.Parse(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid database URL: %w", err)
	}

	driver := "postgres"
	if u.Scheme != "" && u.Scheme != "postgres" && u.Scheme != "postgresql" {
		driver = u.Scheme
	}

	db, err := sql.Open(driver, cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Configure pool parameters
	if cfg.MaxOpenConns > 0 {
		db.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		db.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	}
	if cfg.ConnMaxIdleTime > 0 {
		db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	}

	safeURL := cfg.DatabaseURL
	if u.User != nil {
		safeURL = u.Redacted()
	}

	return &Store{
		db:      db,
		driver:  driver,
		safeURL: safeURL,
	}, nil
}

// NewFromDB constructs a Store from an existing sql.DB handle (useful for testing and mocks).
func NewFromDB(db *sql.DB) *Store {
	return &Store{
		db:     db,
		driver: "mock",
	}
}

// DB returns the underlying *sql.DB connection pool.
func (s *Store) DB() *sql.DB {
	return s.db
}

// Ping verifies that a connection to the database is active.
func (s *Store) Ping(ctx context.Context) error {
	if s.db == nil {
		return ErrStoreClosed
	}
	ctxTimeout, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return s.db.PingContext(ctxTimeout)
}

// Close gracefully releases all connection pool resources.
func (s *Store) Close() error {
	if s.db == nil {
		return nil
	}
	return s.db.Close()
}

// WithTx executes a function within an isolated database transaction.
// If fn returns an error or panics, the transaction is automatically rolled back.
func (s *Store) WithTx(ctx context.Context, fn func(tx *sql.Tx) error) (err error) {
	if s.db == nil {
		return ErrStoreClosed
	}

	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()

	if err = fn(tx); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

// SetTenantContext sets the transaction-local app.tenant_id GUC for RLS defense-in-depth.
func (s *Store) SetTenantContext(ctx context.Context, tx *sql.Tx, tenantID string) error {
	if tx == nil {
		return errors.New("cannot set tenant context on nil transaction")
	}
	_, err := tx.ExecContext(ctx, "SELECT set_config('app.tenant_id', $1, true)", tenantID)
	if err != nil {
		return fmt.Errorf("failed to establish tenant context: %w", err)
	}
	return nil
}
