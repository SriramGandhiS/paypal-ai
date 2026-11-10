package httpx_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/paybridge/paybridge-ai/services/api/internal/httpx"
)

type mockPinger struct {
	err error
}

func (m *mockPinger) Ping(ctx context.Context) error {
	return m.err
}

func TestHealthEndpoints(t *testing.T) {
	t.Run("healthz returns 200 OK", func(t *testing.T) {
		router := httpx.NewRouter(httpx.DefaultServerConfig(":8080", nil))
		req := httptest.NewRequest("GET", "/healthz", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
		if w.Header().Get("X-Request-ID") == "" {
			t.Error("expected non-empty X-Request-ID header")
		}

		var body map[string]any
		if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if body["status"] != "UP" {
			t.Errorf("expected status UP, got %v", body["status"])
		}
	})

	t.Run("readyz returns 200 when database is healthy", func(t *testing.T) {
		router := httpx.NewRouter(httpx.DefaultServerConfig(":8080", &mockPinger{err: nil}))
		req := httptest.NewRequest("GET", "/readyz", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
	})

	t.Run("readyz returns 503 when database is down", func(t *testing.T) {
		router := httpx.NewRouter(httpx.DefaultServerConfig(":8080", &mockPinger{err: errors.New("db connection refused")}))
		req := httptest.NewRequest("GET", "/readyz", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusServiceUnavailable {
			t.Errorf("expected 503, got %d", w.Code)
		}
	})
}

func TestProblemDetails(t *testing.T) {
	t.Run("WriteProblem formats RFC 7807 JSON correctly", func(t *testing.T) {
		w := httptest.NewRecorder()
		prob := httpx.BadRequest("Amount must be greater than zero", "INVALID_AMOUNT")

		httpx.WriteProblem(w, prob)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
		if w.Header().Get("Content-Type") != "application/problem+json" {
			t.Errorf("expected application/problem+json, got %s", w.Header().Get("Content-Type"))
		}

		var decoded httpx.ProblemDetails
		if err := json.NewDecoder(w.Body).Decode(&decoded); err != nil {
			t.Fatalf("failed to decode problem details: %v", err)
		}
		if decoded.Code != "INVALID_AMOUNT" {
			t.Errorf("expected code INVALID_AMOUNT, got %s", decoded.Code)
		}
	})

	t.Run("panic recovery returns 500 ProblemDetails", func(t *testing.T) {
		router := httpx.NewRouter(httpx.DefaultServerConfig(":8080", nil))
		router.HandleFunc("GET /panic", func(w http.ResponseWriter, r *http.Request) {
			panic("unexpected memory fault")
		})

		req := httptest.NewRequest("GET", "/panic", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", w.Code)
		}
		if w.Header().Get("Content-Type") != "application/problem+json" {
			t.Errorf("expected application/problem+json, got %s", w.Header().Get("Content-Type"))
		}
	})
}
