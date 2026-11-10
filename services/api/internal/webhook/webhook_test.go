package webhook_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/paybridge/paybridge-ai/services/api/internal/paypal"
	"github.com/paybridge/paybridge-ai/services/api/internal/store"
	"github.com/paybridge/paybridge-ai/services/api/internal/webhook"
)

func TestWebhookGateway(t *testing.T) {
	st := store.NewFromDB(nil)

	t.Run("missing headers returns 400 Bad Request", func(t *testing.T) {
		gw := webhook.NewGateway(&paypal.LocalMockVerifier{ShouldPass: true}, st, "wh-123", nil)
		handler := gw.Handler()

		req := httptest.NewRequest("POST", "/api/v1/webhooks/paypal", bytes.NewBufferString(`{"id":"evt-1"}`))
		w := httptest.NewRecorder()

		handler(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400 for missing headers, got %d", w.Code)
		}
	})

	t.Run("invalid signature returns 400 Bad Request", func(t *testing.T) {
		gw := webhook.NewGateway(&paypal.LocalMockVerifier{ShouldPass: false}, st, "wh-123", nil)
		handler := gw.Handler()

		req := httptest.NewRequest("POST", "/api/v1/webhooks/paypal", bytes.NewBufferString(`{"id":"evt-1"}`))
		req.Header.Set("PAYPAL-TRANSMISSION-ID", "tx-999")
		req.Header.Set("PAYPAL-TRANSMISSION-SIG", "invalid-sig")

		w := httptest.NewRecorder()
		handler(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400 for invalid signature, got %d", w.Code)
		}
	})

	t.Run("verifier network error returns 503 so PayPal retries", func(t *testing.T) {
		gw := webhook.NewGateway(&paypal.LocalMockVerifier{ShouldFail: true}, st, "wh-123", nil)
		handler := gw.Handler()

		req := httptest.NewRequest("POST", "/api/v1/webhooks/paypal", bytes.NewBufferString(`{"id":"evt-1"}`))
		req.Header.Set("PAYPAL-TRANSMISSION-ID", "tx-999")
		req.Header.Set("PAYPAL-TRANSMISSION-SIG", "any-sig")

		w := httptest.NewRecorder()
		handler(w, req)

		if w.Code != http.StatusServiceUnavailable {
			t.Errorf("expected 503 when verifier fails, got %d", w.Code)
		}
	})

	t.Run("valid verified webhook returns 200 OK", func(t *testing.T) {
		gw := webhook.NewGateway(&paypal.LocalMockVerifier{ShouldPass: true}, st, "wh-123", nil)
		handler := gw.Handler()

		req := httptest.NewRequest("POST", "/api/v1/webhooks/paypal", bytes.NewBufferString(`{"id":"evt-100","event_type":"PAYMENT.CAPTURE.COMPLETED"}`))
		req.Header.Set("PAYPAL-TRANSMISSION-ID", "tx-valid-01")
		req.Header.Set("PAYPAL-TRANSMISSION-SIG", "valid-sig-hex")

		w := httptest.NewRecorder()
		handler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200 for valid webhook, got %d", w.Code)
		}
	})

	t.Run("duplicate webhook delivery is accepted idempotently", func(t *testing.T) {
		gw := webhook.NewGateway(&paypal.LocalMockVerifier{ShouldPass: true}, st, "wh-123", nil)
		handler := gw.Handler()

		payload := `{"id":"evt-duplicate-100","event_type":"PAYMENT.CAPTURE.COMPLETED"}`

		// First delivery
		req1 := httptest.NewRequest("POST", "/api/v1/webhooks/paypal", bytes.NewBufferString(payload))
		req1.Header.Set("PAYPAL-TRANSMISSION-ID", "tx-dup-01")
		req1.Header.Set("PAYPAL-TRANSMISSION-SIG", "valid-sig-hex")
		w1 := httptest.NewRecorder()
		handler(w1, req1)

		if w1.Code != http.StatusOK {
			t.Fatalf("expected 200 for initial webhook, got %d", w1.Code)
		}

		// Duplicate delivery with same event ID
		req2 := httptest.NewRequest("POST", "/api/v1/webhooks/paypal", bytes.NewBufferString(payload))
		req2.Header.Set("PAYPAL-TRANSMISSION-ID", "tx-dup-02") // Provider retransmissions may carry new transmission IDs
		req2.Header.Set("PAYPAL-TRANSMISSION-SIG", "valid-sig-hex")
		w2 := httptest.NewRecorder()
		handler(w2, req2)

		if w2.Code != http.StatusOK {
			t.Fatalf("expected 200 for duplicate webhook retransmission, got %d", w2.Code)
		}
	})
}
