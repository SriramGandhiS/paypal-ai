package paypal_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/paybridge/paybridge-ai/services/api/internal/paypal"
)

func TestPayPalClient_OAuthTokenCaching(t *testing.T) {
	callCount := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/oauth2/token" {
			callCount++
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"access_token": "mock-access-token-12345",
				"token_type":   "Bearer",
				"expires_in":   3600,
			})
			return
		}
		http.NotFound(w, r)
	}))
	defer ts.Close()

	client := paypal.NewClient(paypal.ClientConfig{
		APIBase:      ts.URL,
		ClientID:     "client-id",
		ClientSecret: "client-secret",
	})

	ctx := context.Background()

	// Call 1: Fetches from server
	tok1, err := client.GetAccessToken(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok1 != "mock-access-token-12345" {
		t.Fatalf("expected mock-access-token-12345, got %s", tok1)
	}

	// Call 2: Uses cached token
	tok2, err := client.GetAccessToken(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok2 != tok1 {
		t.Fatalf("expected cached token %s, got %s", tok1, tok2)
	}

	if callCount != 1 {
		t.Fatalf("expected 1 token call due to caching, got %d", callCount)
	}
}

func TestPayPalClient_CreateAndCaptureOrder(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/v1/oauth2/token":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"access_token": "token-xyz",
				"expires_in":   3600,
			})
		case "/v2/checkout/orders":
			if r.Method != http.MethodPost {
				http.Error(w, "bad method", http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id":     "ORDER-999",
				"status": "CREATED",
				"links": []map[string]any{
					{
						"href": "https://www.sandbox.paypal.com/checkoutnow?token=ORDER-999",
						"rel":  "approve",
					},
				},
			})
		case "/v2/checkout/orders/ORDER-999/capture":
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id":     "ORDER-999",
				"status": "COMPLETED",
				"purchase_units": []map[string]any{
					{
						"payments": map[string]any{
							"captures": []map[string]any{
								{
									"id":     "CAPTURE-888",
									"status": "COMPLETED",
									"amount": map[string]any{
										"currency_code": "USD",
										"value":         "50.00",
									},
								},
							},
						},
					},
				},
			})
		case "/v2/payments/captures/CAPTURE-888":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id":     "CAPTURE-888",
				"status": "COMPLETED",
				"amount": map[string]any{
					"currency_code": "USD",
					"value":         "50.00",
				},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	client := paypal.NewClient(paypal.ClientConfig{
		APIBase:      ts.URL,
		ClientID:     "cid",
		ClientSecret: "csec",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1. Create Order
	ord, err := client.CreateOrder(ctx, paypal.OrderRequest{
		ReferenceID: "ref-100",
		AmountCents: 5000,
		Currency:    "USD",
	})
	if err != nil {
		t.Fatalf("CreateOrder failed: %v", err)
	}
	if ord.ID != "ORDER-999" {
		t.Errorf("expected ID ORDER-999, got %s", ord.ID)
	}
	if ord.ApprovalURL != "https://www.sandbox.paypal.com/checkoutnow?token=ORDER-999" {
		t.Errorf("unexpected approval url: %s", ord.ApprovalURL)
	}

	// 2. Capture Order
	capResp, err := client.CaptureOrder(ctx, "ORDER-999")
	if err != nil {
		t.Fatalf("CaptureOrder failed: %v", err)
	}
	if capResp.CaptureID != "CAPTURE-888" {
		t.Errorf("expected capture ID CAPTURE-888, got %s", capResp.CaptureID)
	}
	if capResp.CaptureStatus != "COMPLETED" {
		t.Errorf("expected capture status COMPLETED, got %s", capResp.CaptureStatus)
	}

	// 3. Get Capture Ground Truth
	truth, err := client.GetCapture(ctx, "CAPTURE-888")
	if err != nil {
		t.Fatalf("GetCapture failed: %v", err)
	}
	if truth.AmountCents != 5000 {
		t.Errorf("expected 5000 cents, got %d", truth.AmountCents)
	}
	if truth.Status != "COMPLETED" {
		t.Errorf("expected status COMPLETED, got %s", truth.Status)
	}
}

func TestPayPalClient_WebhookSignatureVerification(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/v1/oauth2/token":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"access_token": "token-xyz",
				"expires_in":   3600,
			})
		case "/v1/notifications/verify-webhook-signature":
			var payload struct {
				TransmissionID string `json:"transmission_id"`
			}
			_ = json.NewDecoder(r.Body).Decode(&payload)
			if payload.TransmissionID == "valid-tx" {
				_ = json.NewEncoder(w).Encode(map[string]any{
					"verification_status": "SUCCESS",
				})
			} else {
				_ = json.NewEncoder(w).Encode(map[string]any{
					"verification_status": "FAILURE",
				})
			}
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	client := paypal.NewClient(paypal.ClientConfig{
		APIBase:      ts.URL,
		ClientID:     "cid",
		ClientSecret: "csec",
	})

	ctx := context.Background()
	headersValid := paypal.WebhookHeaders{
		TransmissionID:   "valid-tx",
		TransmissionTime: "2026-11-10T12:00:00Z",
		TransmissionSig:  "sig",
		CertURL:          "https://api.sandbox.paypal.com/cert",
		AuthAlgo:         "SHA256withRSA",
		WebhookID:        "WH-123",
	}

	pass, err := client.VerifySignature(ctx, headersValid, []byte(`{"id":"WH-EVT-1"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !pass {
		t.Errorf("expected signature verification to pass")
	}

	headersInvalid := headersValid
	headersInvalid.TransmissionID = "invalid-tx"

	fail, err := client.VerifySignature(ctx, headersInvalid, []byte(`{"id":"WH-EVT-2"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fail {
		t.Errorf("expected signature verification to fail")
	}
}
