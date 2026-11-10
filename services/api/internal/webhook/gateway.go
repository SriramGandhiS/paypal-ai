package webhook

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/paybridge/paybridge-ai/services/api/internal/httpx"
	"github.com/paybridge/paybridge-ai/services/api/internal/paypal"
	"github.com/paybridge/paybridge-ai/services/api/internal/store"
)

// Gateway handles raw PayPal webhook reception.
type Gateway struct {
	verifier  paypal.WebhookVerifier
	store     *store.Store
	webhookID string
	logger    *slog.Logger
}

// NewGateway constructs a new webhook Gateway.
func NewGateway(v paypal.WebhookVerifier, st *store.Store, webhookID string, logger *slog.Logger) *Gateway {
	if logger == nil {
		logger = slog.Default()
	}
	return &Gateway{
		verifier:  v,
		store:     st,
		webhookID: webhookID,
		logger:    logger,
	}
}

// Handler returns the HTTP handler for POST /api/v1/webhooks/paypal.
func (g *Gateway) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Enforce size limit and read raw bytes
		r.Body = http.MaxBytesReader(w, r.Body, 1024*1024) // 1MB limit
		rawBody, err := io.ReadAll(r.Body)
		if err != nil {
			httpx.WriteProblem(w, httpx.BadRequest("failed to read webhook body", "INVALID_BODY"))
			return
		}
		defer r.Body.Close()

		// 2. Extract PayPal signature headers
		headers := paypal.WebhookHeaders{
			TransmissionID:   r.Header.Get("PAYPAL-TRANSMISSION-ID"),
			TransmissionTime: r.Header.Get("PAYPAL-TRANSMISSION-TIME"),
			TransmissionSig:  r.Header.Get("PAYPAL-TRANSMISSION-SIG"),
			CertURL:          r.Header.Get("PAYPAL-CERT-URL"),
			AuthAlgo:         r.Header.Get("PAYPAL-AUTH-ALGO"),
			WebhookID:        g.webhookID,
		}

		if headers.TransmissionID == "" || headers.TransmissionSig == "" {
			httpx.WriteProblem(w, httpx.BadRequest("missing required PayPal signature headers", "MISSING_HEADERS"))
			return
		}

		// 3. Verify signature
		verified, err := g.verifier.VerifySignature(r.Context(), headers, rawBody)
		if err != nil {
			g.logger.Error("webhook signature verification service unavailable", "error", err)
			httpx.WriteProblem(w, httpx.ServiceUnavailable("signature verification unavailable; retry later"))
			return
		}

		if !verified {
			g.logger.Warn("webhook signature verification rejected", "transmission_id", headers.TransmissionID)
			httpx.WriteProblem(w, httpx.BadRequest("invalid webhook signature", "SIGNATURE_REJECTED"))
			return
		}

		// 4. Compute SHA-256 hash of raw payload
		h := sha256.Sum256(rawBody)
		payloadHash := hex.EncodeToString(h[:])

		// Parse basic event metadata
		var meta struct {
			ID        string `json:"id"`
			EventType string `json:"event_type"`
		}
		_ = json.Unmarshal(rawBody, &meta)

		// 5. Store in webhook_inbox if database connection is available
		if g.store != nil && g.store.DB() != nil {
			query := `
				INSERT INTO webhook_inbox (
					environment, provider, provider_event_id, transmission_id,
					event_type, signature_verified, payload_sha256, raw_payload, processing_status
				) VALUES (
					'REAL_SANDBOX', 'PAYPAL', $1, $2, $3, TRUE, $4, $5, 'RECEIVED'
				)
				ON CONFLICT (environment, provider, provider_event_id)
				DO UPDATE SET duplicate_count = webhook_inbox.duplicate_count + 1;
			`
			_, err = g.store.DB().ExecContext(r.Context(), query,
				meta.ID, headers.TransmissionID, meta.EventType, payloadHash, string(rawBody),
			)
			if err != nil {
				g.logger.Error("failed to record webhook in inbox", "error", err)
				httpx.WriteProblem(w, httpx.InternalError("failed to record webhook"))
				return
			}
		}

		// 6. Return 200 OK only after successful verification and storage
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status":          "ACCEPTED",
			"transmission_id": headers.TransmissionID,
		})
	}
}
