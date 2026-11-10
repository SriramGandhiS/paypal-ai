package httpx

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

// ContextKey is a custom type for context values to prevent key collisions.
type ContextKey string

const (
	RequestIDKey ContextKey = "request_id"
)

// Pinger interface defines readiness check capabilities.
type Pinger interface {
	Ping(ctx context.Context) error
}

// ServerConfig configures the HTTP server instance.
type ServerConfig struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	Logger       *slog.Logger
	Pinger       Pinger
}

// DefaultServerConfig provides production-hardened timeout defaults.
func DefaultServerConfig(addr string, pinger Pinger) ServerConfig {
	return ServerConfig{
		Addr:         addr,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
		Logger:       slog.Default(),
		Pinger:       pinger,
	}
}

// Router wraps the HTTP handler and routes.
type Router struct {
	mux    *http.ServeMux
	cfg    ServerConfig
	logger *slog.Logger
}

// NewRouter constructs a configured HTTP router with health checks and middleware.
func NewRouter(cfg ServerConfig) *Router {
	if cfg.Logger == nil {
		cfg.Logger = slog.Default()
	}
	r := &Router{
		mux:    http.NewServeMux(),
		cfg:    cfg,
		logger: cfg.Logger,
	}
	r.registerRoutes()
	return r
}

func (r *Router) registerRoutes() {
	// 1. Healthz (Liveness check)
	r.mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status":    "UP",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"service":   "paybridge-api",
		})
	})

	// 2. Readyz (Readiness check - checks database connectivity)
	r.mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.cfg.Pinger != nil {
			ctx, cancel := context.WithTimeout(req.Context(), 2*time.Second)
			defer cancel()
			if err := r.cfg.Pinger.Ping(ctx); err != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				_ = json.NewEncoder(w).Encode(map[string]any{
					"status":  "DOWN",
					"reason":  "database unavailable",
					"service": "paybridge-api",
				})
				return
			}
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status":    "READY",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"service":   "paybridge-api",
		})
	})
}

// HandleFunc registers a new route on the router.
func (r *Router) HandleFunc(pattern string, handler http.HandlerFunc) {
	r.mux.HandleFunc(pattern, handler)
}

// ServeHTTP implements http.Handler with panic recovery and request correlation.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Correlation / Request ID extraction or generation
	reqID := req.Header.Get("X-Request-ID")
	if reqID == "" {
		b := make([]byte, 16)
		_, _ = rand.Read(b)
		reqID = hex.EncodeToString(b)
	}

	w.Header().Set("X-Request-ID", reqID)
	w.Header().Set("X-Correlation-ID", reqID)

	ctx := context.WithValue(req.Context(), RequestIDKey, reqID)
	reqWithContext := req.WithContext(ctx)

	// Panic Recovery Middleware
	defer func() {
		if rec := recover(); rec != nil {
			r.logger.Error("panic recovered in HTTP handler",
				"error", rec,
				"request_id", reqID,
				"path", req.URL.Path,
			)
			WriteProblem(w, InternalError("An unexpected server error occurred"))
		}
	}()

	r.mux.ServeHTTP(w, reqWithContext)
}
