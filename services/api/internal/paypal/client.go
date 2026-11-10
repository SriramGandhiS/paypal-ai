package paypal

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	ErrPayPalAuthFailed     = errors.New("paypal oauth authentication failed")
	ErrPayPalOrderFailed    = errors.New("paypal order creation failed")
	ErrPayPalCaptureFailed  = errors.New("paypal capture failed")
	ErrPayPalResourceFailed = errors.New("paypal resource fetch failed")
)

// ClientConfig holds configuration for the PayPal Sandbox client.
type ClientConfig struct {
	APIBase      string
	ClientID     string
	ClientSecret string
	HTTPClient   *http.Client
}

// Client interacts with the PayPal Sandbox REST API.
type Client struct {
	cfg         ClientConfig
	tokenMutex  sync.RWMutex
	accessToken string
	tokenExpiry time.Time
}

// NewClient creates a new PayPal client.
func NewClient(cfg ClientConfig) *Client {
	if cfg.HTTPClient == nil {
		cfg.HTTPClient = &http.Client{
			Timeout: 15 * time.Second,
		}
	}
	cfg.APIBase = strings.TrimSuffix(cfg.APIBase, "/")
	return &Client{
		cfg: cfg,
	}
}

// OrderRequest specifies order creation parameters.
type OrderRequest struct {
	ReferenceID string `json:"reference_id"`
	AmountCents int64  `json:"amount_cents"`
	Currency    string `json:"currency"`
	ReturnURL   string `json:"return_url,omitempty"`
	CancelURL   string `json:"cancel_url,omitempty"`
}

// OrderResponse represents a created or fetched PayPal order.
type OrderResponse struct {
	ID            string `json:"id"`
	Status        string `json:"status"` // CREATED, SAVED, APPROVED, VOIDED, COMPLETED, PAYER_ACTION_REQUIRED
	ApprovalURL   string `json:"approval_url,omitempty"`
	CaptureID     string `json:"capture_id,omitempty"`
	GrossAmount   string `json:"gross_amount,omitempty"`
	Currency      string `json:"currency,omitempty"`
	CaptureStatus string `json:"capture_status,omitempty"`
}

// CaptureGroundTruth represents the verified provider state obtained directly via GET.
type CaptureGroundTruth struct {
	CaptureID   string
	Status      string // COMPLETED, PENDING, DENIED, REFUNDED
	AmountCents int64
	Currency    string
}

// GetAccessToken obtains or returns a valid cached OAuth 2.0 bearer token.
func (c *Client) GetAccessToken(ctx context.Context) (string, error) {
	c.tokenMutex.RLock()
	if c.accessToken != "" && time.Now().Add(60*time.Second).Before(c.tokenExpiry) {
		token := c.accessToken
		c.tokenMutex.RUnlock()
		return token, nil
	}
	c.tokenMutex.RUnlock()

	c.tokenMutex.Lock()
	defer c.tokenMutex.Unlock()

	// Double-check after acquiring write lock
	if c.accessToken != "" && time.Now().Add(60*time.Second).Before(c.tokenExpiry) {
		return c.accessToken, nil
	}

	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	tokenURL := fmt.Sprintf("%s/v1/oauth2/token", c.cfg.APIBase)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to build token request: %w", err)
	}

	req.SetBasicAuth(c.cfg.ClientID, c.cfg.ClientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := c.cfg.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("%w: network error: %v", ErrPayPalAuthFailed, err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%w: status %d: %s", ErrPayPalAuthFailed, resp.StatusCode, string(bodyBytes))
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}

	if err := json.Unmarshal(bodyBytes, &tokenResp); err != nil {
		return "", fmt.Errorf("failed to parse token response: %w", err)
	}

	c.accessToken = tokenResp.AccessToken
	c.tokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	return c.accessToken, nil
}

// CreateOrder initiates a PayPal checkout order with intent=CAPTURE.
func (c *Client) CreateOrder(ctx context.Context, ord OrderRequest) (*OrderResponse, error) {
	token, err := c.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	amountDecimal := fmt.Sprintf("%.2f", float64(ord.AmountCents)/100.0)

	payload := map[string]any{
		"intent": "CAPTURE",
		"purchase_units": []map[string]any{
			{
				"reference_id": ord.ReferenceID,
				"amount": map[string]any{
					"currency_code": ord.Currency,
					"value":         amountDecimal,
				},
			},
		},
	}

	if ord.ReturnURL != "" || ord.CancelURL != "" {
		appCtx := map[string]any{
			"user_action": "PAY_NOW",
		}
		if ord.ReturnURL != "" {
			appCtx["return_url"] = ord.ReturnURL
		}
		if ord.CancelURL != "" {
			appCtx["cancel_url"] = ord.CancelURL
		}
		payload["application_context"] = appCtx
	}

	reqBody, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	orderURL := fmt.Sprintf("%s/v2/checkout/orders", c.cfg.APIBase)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, orderURL, bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "return=representation")

	resp, err := c.cfg.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrPayPalOrderFailed, err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: status %d: %s", ErrPayPalOrderFailed, resp.StatusCode, string(bodyBytes))
	}

	var res struct {
		ID     string `json:"id"`
		Status string `json:"status"`
		Links  []struct {
			Href string `json:"href"`
			Rel  string `json:"rel"`
		} `json:"links"`
	}

	if err := json.Unmarshal(bodyBytes, &res); err != nil {
		return nil, fmt.Errorf("failed to parse order response: %w", err)
	}

	orderResp := &OrderResponse{
		ID:          res.ID,
		Status:      res.Status,
		Currency:    ord.Currency,
		GrossAmount: amountDecimal,
	}

	for _, link := range res.Links {
		if link.Rel == "approve" || link.Rel == "payer-action" {
			orderResp.ApprovalURL = link.Href
			break
		}
	}

	return orderResp, nil
}

// CaptureOrder captures funds for an approved PayPal order.
func (c *Client) CaptureOrder(ctx context.Context, orderID string) (*OrderResponse, error) {
	token, err := c.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	captureURL := fmt.Sprintf("%s/v2/checkout/orders/%s/capture", c.cfg.APIBase, url.PathEscape(orderID))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, captureURL, bytes.NewReader([]byte("{}")))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "return=representation")

	resp, err := c.cfg.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrPayPalCaptureFailed, err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: status %d: %s", ErrPayPalCaptureFailed, resp.StatusCode, string(bodyBytes))
	}

	var res struct {
		ID            string `json:"id"`
		Status        string `json:"status"`
		PurchaseUnits []struct {
			Payments struct {
				Captures []struct {
					ID     string `json:"id"`
					Status string `json:"status"`
					Amount struct {
						CurrencyCode string `json:"currency_code"`
						Value        string `json:"value"`
					} `json:"amount"`
				} `json:"captures"`
			} `json:"payments"`
		} `json:"purchase_units"`
	}

	if err := json.Unmarshal(bodyBytes, &res); err != nil {
		return nil, fmt.Errorf("failed to parse capture response: %w", err)
	}

	ordResp := &OrderResponse{
		ID:     res.ID,
		Status: res.Status,
	}

	if len(res.PurchaseUnits) > 0 && len(res.PurchaseUnits[0].Payments.Captures) > 0 {
		cap0 := res.PurchaseUnits[0].Payments.Captures[0]
		ordResp.CaptureID = cap0.ID
		ordResp.CaptureStatus = cap0.Status
		ordResp.GrossAmount = cap0.Amount.Value
		ordResp.Currency = cap0.Amount.CurrencyCode
	}

	return ordResp, nil
}

// GetOrder retrieves current order state from PayPal.
func (c *Client) GetOrder(ctx context.Context, orderID string) (*OrderResponse, error) {
	token, err := c.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	orderURL := fmt.Sprintf("%s/v2/checkout/orders/%s", c.cfg.APIBase, url.PathEscape(orderID))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, orderURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.cfg.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrPayPalResourceFailed, err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: status %d: %s", ErrPayPalResourceFailed, resp.StatusCode, string(bodyBytes))
	}

	var res struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}
	if err := json.Unmarshal(bodyBytes, &res); err != nil {
		return nil, err
	}

	return &OrderResponse{
		ID:     res.ID,
		Status: res.Status,
	}, nil
}

// GetCapture retrieves provider ground truth for a capture ID.
func (c *Client) GetCapture(ctx context.Context, captureID string) (*CaptureGroundTruth, error) {
	token, err := c.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	captureURL := fmt.Sprintf("%s/v2/payments/captures/%s", c.cfg.APIBase, url.PathEscape(captureID))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, captureURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.cfg.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrPayPalResourceFailed, err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: status %d: %s", ErrPayPalResourceFailed, resp.StatusCode, string(bodyBytes))
	}

	var res struct {
		ID     string `json:"id"`
		Status string `json:"status"`
		Amount struct {
			CurrencyCode string `json:"currency_code"`
			Value        string `json:"value"`
		} `json:"amount"`
	}
	if err := json.Unmarshal(bodyBytes, &res); err != nil {
		return nil, err
	}

	cents, _ := parseAmountToCents(res.Amount.Value)

	return &CaptureGroundTruth{
		CaptureID:   res.ID,
		Status:      res.Status,
		AmountCents: cents,
		Currency:    res.Amount.CurrencyCode,
	}, nil
}

// VerifySignature validates an incoming PayPal webhook notification via PayPal's API.
// Satisfies the WebhookVerifier interface.
func (c *Client) VerifySignature(ctx context.Context, headers WebhookHeaders, rawBody []byte) (bool, error) {
	token, err := c.GetAccessToken(ctx)
	if err != nil {
		return false, err
	}

	var webhookEvent any
	if err := json.Unmarshal(rawBody, &webhookEvent); err != nil {
		return false, fmt.Errorf("invalid json webhook body: %w", err)
	}

	verifyPayload := map[string]any{
		"transmission_id":   headers.TransmissionID,
		"transmission_time": headers.TransmissionTime,
		"transmission_sig":  headers.TransmissionSig,
		"cert_url":          headers.CertURL,
		"auth_algo":         headers.AuthAlgo,
		"webhook_id":        headers.WebhookID,
		"webhook_event":     webhookEvent,
	}

	payloadBytes, err := json.Marshal(verifyPayload)
	if err != nil {
		return false, err
	}

	verifyURL := fmt.Sprintf("%s/v1/notifications/verify-webhook-signature", c.cfg.APIBase)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, verifyURL, bytes.NewReader(payloadBytes))
	if err != nil {
		return false, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.cfg.HTTPClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("%w: %v", ErrVerifierUnavailable, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("%w: upstream status %d", ErrVerifierUnavailable, resp.StatusCode)
	}

	var res struct {
		VerificationStatus string `json:"verification_status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return false, fmt.Errorf("failed to decode verification response: %w", err)
	}

	return res.VerificationStatus == "SUCCESS", nil
}

func parseAmountToCents(val string) (int64, error) {
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0, err
	}
	return int64(f*100 + 0.5), nil
}
