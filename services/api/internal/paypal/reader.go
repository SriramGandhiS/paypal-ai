package paypal

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Reader is the strictly GET-only interface for PayPal queries.
// Compile-time invariant: No POST, PUT, PATCH, or DELETE methods exist on this interface.
type Reader interface {
	ShowOrder(ctx context.Context, env Environment, orderID string) (Order, RawEvidence, error)
	ShowCapture(ctx context.Context, env Environment, captureID string) (Capture, RawEvidence, error)
}

type readerClient struct {
	apiBase      string
	tokenManager *TokenManager
	httpClient   *http.Client
}

// Compile-time assertion of interface implementation
var _ Reader = (*readerClient)(nil)

func NewReader(apiBase string, tm *TokenManager, client *http.Client) Reader {
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	return &readerClient{
		apiBase:      strings.TrimRight(apiBase, "/"),
		tokenManager: tm,
		httpClient:   client,
	}
}

func (r *readerClient) ShowOrder(ctx context.Context, env Environment, orderID string) (Order, RawEvidence, error) {
	endpoint := fmt.Sprintf("%s/v2/checkout/orders/%s", r.apiBase, orderID)
	bodyBytes, statusCode, err := r.doGet(ctx, endpoint)
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

func (r *readerClient) ShowCapture(ctx context.Context, env Environment, captureID string) (Capture, RawEvidence, error) {
	endpoint := fmt.Sprintf("%s/v2/payments/captures/%s", r.apiBase, captureID)
	bodyBytes, statusCode, err := r.doGet(ctx, endpoint)
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

func (r *readerClient) doGet(ctx context.Context, endpoint string) ([]byte, int, error) {
	tok, err := r.tokenManager.GetToken(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("oauth token error: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Authorization", "Bearer "+tok)
	req.Header.Set("Accept", "application/json")

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	if resp.StatusCode == http.StatusUnauthorized {
		r.tokenManager.Invalidate()
		return nil, resp.StatusCode, fmt.Errorf("unauthorized response (401)")
	}

	return bodyBytes, resp.StatusCode, nil
}

func sanitizePayload(in map[string]any) map[string]any {
	out := make(map[string]any, len(in))
	for k, v := range in {
		lower := strings.ToLower(k)
		if strings.Contains(lower, "password") ||
			strings.Contains(lower, "secret") ||
			strings.Contains(lower, "token") ||
			strings.Contains(lower, "card_number") {
			out[k] = "[REDACTED]"
		} else {
			out[k] = v
		}
	}
	return out
}
