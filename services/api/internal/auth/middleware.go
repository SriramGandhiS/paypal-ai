package auth

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/paybridge/paybridge-ai/services/api/internal/httpx"
	"github.com/paybridge/paybridge-ai/services/api/internal/store"
	"github.com/paybridge/paybridge-ai/services/api/internal/tenancy"
)

// SessionReader abstracts session lookup from database or test mock.
type SessionReader interface {
	LookupSession(ctx context.Context, tokenHash []byte) (*Session, error)
}

// DatabaseSessionStore implements SessionReader against PostgreSQL.
type DatabaseSessionStore struct {
	store *store.Store
}

// NewDatabaseSessionStore constructs a store-backed session reader.
func NewDatabaseSessionStore(st *store.Store) *DatabaseSessionStore {
	return &DatabaseSessionStore{store: st}
}

// LookupSession queries active sessions by sha256 id_hash.
func (d *DatabaseSessionStore) LookupSession(ctx context.Context, tokenHash []byte) (*Session, error) {
	if d.store == nil || d.store.DB() == nil {
		return nil, errors.New("database store unavailable")
	}

	query := `
		SELECT s.user_id, s.tenant_id, s.role, s.expires_at, s.revoked_at
		  FROM sessions s
		 WHERE s.id_hash = $1;
	`
	var s Session
	s.IDHash = tokenHash

	err := d.store.DB().QueryRowContext(ctx, query, tokenHash).Scan(
		&s.UserID, &s.TenantID, &s.Role, &s.ExpiresAt, &s.RevokedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrInvalidSessionToken
	} else if err != nil {
		return nil, err
	}

	if s.RevokedAt != nil {
		return nil, ErrInvalidSessionToken
	}
	if time.Now().UTC().After(s.ExpiresAt) {
		return nil, ErrSessionExpired
	}

	return &s, nil
}

// MiddlewareConfig configures authentication middleware behavior.
type MiddlewareConfig struct {
	Reader         SessionReader
	Logger         *slog.Logger
	AllowDevBypass bool
}

// Middleware provides tenant-aware authentication enforcement.
type Middleware struct {
	reader         SessionReader
	logger         *slog.Logger
	allowDevBypass bool
}

// NewMiddleware constructs an auth Middleware.
func NewMiddleware(cfg MiddlewareConfig) *Middleware {
	if cfg.Logger == nil {
		cfg.Logger = slog.Default()
	}
	return &Middleware{
		reader:         cfg.Reader,
		logger:         cfg.Logger,
		allowDevBypass: cfg.AllowDevBypass,
	}
}

// Authenticate wraps an http.Handler and extracts tenancy context.
func (m *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// 1. Check for Development / Testing bypass headers (X-Dev-Tenant-ID)
		devTenant := r.Header.Get("X-Dev-Tenant-ID")
		if devTenant != "" {
			if !m.allowDevBypass {
				m.logger.Warn("security violation: client attempted X-Dev-Tenant-ID bypass while disabled",
					"attempted_tenant", devTenant,
					"remote_addr", r.RemoteAddr,
				)
				httpx.WriteProblem(w, httpx.Forbidden("development authentication bypass is disabled"))
				return
			}

			m.logger.Warn("DEV AUTHENTICATION BYPASS ACTIVE - DO NOT USE IN PRODUCTION",
				"tenant_id", devTenant,
			)

			devUser := r.Header.Get("X-Dev-User-ID")
			if devUser == "" {
				devUser = "00000000-0000-0000-0000-000000000001"
			}
			devRole := r.Header.Get("X-Dev-Role")
			if devRole == "" {
				devRole = "MERCHANT_ADMIN"
			}

			ctx = tenancy.WithTenant(ctx, devTenant)
			ctx = tenancy.WithUser(ctx, devUser, devRole)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// 2. Extract token from cookie or Authorization Bearer header
		var rawToken string
		if cookie, err := r.Cookie(SessionCookieName); err == nil && cookie.Value != "" {
			rawToken = cookie.Value
		} else {
			authHeader := r.Header.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				rawToken = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if rawToken == "" {
			httpx.WriteProblem(w, httpx.Unauthorized("missing authentication session"))
			return
		}

		// 3. Hash and lookup session
		tokenHash := HashToken(rawToken)
		if m.reader == nil {
			httpx.WriteProblem(w, httpx.ServiceUnavailable("session verification unavailable"))
			return
		}

		sess, err := m.reader.LookupSession(ctx, tokenHash)
		if err != nil {
			if errors.Is(err, ErrSessionExpired) {
				httpx.WriteProblem(w, httpx.Unauthorized("session has expired"))
				return
			}
			httpx.WriteProblem(w, httpx.Unauthorized("invalid or revoked session"))
			return
		}

		// 4. Inject tenancy and user context
		ctx = tenancy.WithTenant(ctx, sess.TenantID)
		ctx = tenancy.WithUser(ctx, sess.UserID, sess.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
