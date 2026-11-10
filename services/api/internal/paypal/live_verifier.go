package paypal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// WebhookVerifier performs POST /v1/notifications/verify-webhook-signature.
// Deliberately NOT part of Reader so Reader remains proven GET-only (ADR-021).
type WebhookSignatureVerifier interface {
	VerifyWebhookSignature(ctx context.Context, in WebhookVerifyInput) (VerifyResult, error)
}

type webhookVerifierClient struct {
	apiBase      string
	tokenManager *TokenManager
	httpClient   *http.Client
}

var _ WebhookSignatureVerifier = (*webhookVerifierClient)(nil)

func NewWebhookVerifier(apiBase string, tm *TokenManager, client *http.Client) WebhookSignatureVerifier {
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	return &webhookVerifierClient{
		apiBase:      strings.TrimRight(apiBase, "/"),
		tokenManager: tm,
		httpClient:   client,
	}
}

func (v *webhookVerifierClient) VerifyWebhookSignature(ctx context.Context, in WebhookVerifyInput) (VerifyResult, error) {
	tok, err := v.tokenManager.GetToken(ctx)
	if err != nil {
		return VerifyResult{}, fmt.Errorf("oauth token error: %w", err)
	}

	endpoint := fmt.Sprintf("%s/v1/notifications/verify-webhook-signature", v.apiBase)

	payloadBytes, err := json.Marshal(in)
	if err != nil {
		return VerifyResult{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payloadBytes))
	if err != nil {
		return VerifyResult{}, err
	}
	req.Header.Set("Authorization", "Bearer "+tok)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := v.httpClient.Do(req)
	if err != nil {
		return VerifyResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return VerifyResult{VerificationStatus: "FAILURE"}, fmt.Errorf("verify webhook returned status %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return VerifyResult{}, err
	}

	var res VerifyResult
	if err := json.Unmarshal(bodyBytes, &res); err != nil {
		return VerifyResult{}, err
	}

	return res, nil
}
