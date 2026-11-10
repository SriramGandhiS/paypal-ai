package store_test

import (
	"context"
	"testing"
	"time"

	"github.com/paybridge/paybridge-ai/services/api/internal/store"
)

func TestStoreConfigValidation(t *testing.T) {
	t.Run("empty database URL returns error", func(t *testing.T) {
		cfg := store.StoreConfig{DatabaseURL: ""}
		_, err := store.NewStore(cfg)
		if err == nil {
			t.Fatal("expected error on empty database URL, got nil")
		}
	})

	t.Run("default config sets sensible values", func(t *testing.T) {
		cfg := store.DefaultStoreConfig("postgres://user:pass@localhost:5432/db")
		if cfg.MaxOpenConns != 25 {
			t.Errorf("expected MaxOpenConns 25, got %d", cfg.MaxOpenConns)
		}
		if cfg.MaxIdleConns != 10 {
			t.Errorf("expected MaxIdleConns 10, got %d", cfg.MaxIdleConns)
		}
		if cfg.ConnMaxLifetime != 15*time.Minute {
			t.Errorf("expected ConnMaxLifetime 15m, got %v", cfg.ConnMaxLifetime)
		}
	})

	t.Run("nil db Ping returns ErrStoreClosed", func(t *testing.T) {
		s := store.NewFromDB(nil)
		err := s.Ping(context.Background())
		if err != store.ErrStoreClosed {
			t.Errorf("expected ErrStoreClosed, got %v", err)
		}
	})

	t.Run("nil db WithTx returns ErrStoreClosed", func(t *testing.T) {
		s := store.NewFromDB(nil)
		err := s.WithTx(context.Background(), nil)
		if err != store.ErrStoreClosed {
			t.Errorf("expected ErrStoreClosed, got %v", err)
		}
	})
}
