package paypal

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Mutator is the explicit checkout mutating interface.
// Architecture invariant: Injected ONLY into human-triggered checkout controller; never reachable from AI.
type Mutator interface {
	CreateOrder(ctx context.Context, in CreateOrderInput) (Order, RawEvidence, error)
	CaptureOrder(ctx context.Context, in CaptureInput) (Capture, RawEvidence, error)
}

type mutatorClient struct {
	apiBase      string
	tokenManager *TokenManager
	httpClient   *http.Client
}

var _ Mutator = (*mutatorClient)(nil)

func NewMutator(apiBase string, tm *TokenManager, client *http.Client) Mutator {
	if client == nil {
		client = &http.Client{Timeout: 15 * time.Second}
	}
	return &mutatorClient{
		apiBase:      strings.TrimRight(apiBase, "/"),
		tokenManager: tm,
		httpClient:   client,
	}
}

func (m *mutatorClient) CreateOrder(ctx context.Context, in CreateOrderInput) (Order, RawEvidence, error) {
	if in.RequestID == "" {
		return Order{}, RawEvidence{}, errors.New("PayPal-Request-Id is mandatory for order creation")
	}

	endpoint := fmt.Sprintf("%s/v2/checkout/orders", m.apiBase)

	bodyPayload := map[string]any{
		"intent": "CAPTURE",
		"purchase_units": []map[string]any{
			{
				"reference_id": in.InvoiceID,
				"amount": map[string]string{
					"currency_code": in.Currency,
					"value":         in.Amount,
				},
			},
		},
		"payment_source": map[string]any{
			"paypal": map[string]any{
				"experience_context": map[string]string{
					"return_url": in.ReturnURL,
					"cancel_url": in.CancelURL,
				},
			},
		},
	}

	rawReqBytes, err := json.Marshal(bodyPayload)
	if err != nil {
		return Order{}, RawEvidence{}, err
	}

	bodyBytes, statusCode, err := m.doPost(ctx, endpoint, in.RequestID, rawReqBytes)
	if err != nil {
		return Order{}, RawEvidence{StatusCode: statusCode, ObservedAt: time.Now()}, err
	}

	hash := sha256.Sum256(bodyBytes)
	evidence := RawEvidence{
		StatusCode:    statusCode,
		PayloadSHA256: hex.EncodeToString(hash[:]),
		ObservedAt:    time.Now(),
	}

	var rawMap map[string]any
	if err := json.Unmarshal(bodyBytes, &rawMap); err == nil {
		evidence.RedactedBody = sanitizePayload(rawMap)
	}

	var ord Order
	if err := json.Unmarshal(bodyBytes, &ord); err != nil {
		return Order{}, evidence, fmt.Errorf("failed to parse order response: %w", err)
	}

	return ord, evidence, nil
}

func (m *mutatorClient) CaptureOrder(ctx context.Context, in CaptureInput) (Capture, RawEvidence, error) {
	if in.RequestID == "" {
		return Capture{}, RawEvidence{}, errors.New("PayPal-Request-Id is mandatory for capture")
	}
	if in.OrderID == "" {
		return Capture{}, RawEvidence{}, errors.New("order_id is required for capture")
	}

	endpoint := fmt.Sprintf("%s/v2/checkout/orders/%s/capture", m.apiBase, in.OrderID)

	bodyBytes, statusCode, err := m.doPost(ctx, endpoint, in.RequestID, []byte("{}"))
	if err != nil {
		return Capture{}, RawEvidence{StatusCode: statusCode, ObservedAt: time.Now()}, err
	}

	hash := sha256.Sum256(bodyBytes)
	evidence := RawEvidence{
		StatusCode:    statusCode,
		PayloadSHA256: hex.EncodeToString(hash[:]),
		ObservedAt:    time.Now(),
	}

	var rawMap map[string]any
	if err := json.Unmarshal(bodyBytes, &rawMap); err == nil {
		evidence.RedactedBody = sanitizePayload(rawMap)
	}

	var cap Capture
	if err := json.Unmarshal(bodyBytes, &cap); err != nil {
		return Capture{}, evidence, fmt.Errorf("failed to parse capture response: %w", err)
	}

	return cap, evidence, nil
}

func (m *mutatorClient) doPost(ctx context.Context, endpoint, requestID string, body []byte) ([]byte, int, error) {
	tok, err := m.tokenManager.GetToken(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("oauth token error: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Authorization", "Bearer "+tok)
	req.Header.Set("PayPal-Request-Id", requestID)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	return respBytes, resp.StatusCode, nil
}
