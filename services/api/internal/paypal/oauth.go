package paypal

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type TokenManager struct {
	clientID     string
	clientSecret string
	apiBase      string
	httpClient   *http.Client

	mu        sync.RWMutex
	token     string
	expiresAt time.Time

	refreshMu sync.Mutex
}

func NewTokenManager(clientID, clientSecret, apiBase string, client *http.Client) *TokenManager {
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	return &TokenManager{
		clientID:     clientID,
		clientSecret: clientSecret,
		apiBase:      strings.TrimRight(apiBase, "/"),
		httpClient:   client,
	}
}

// GetToken returns a valid cached OAuth bearer token or fetches a fresh one with singleflight
func (tm *TokenManager) GetToken(ctx context.Context) (string, error) {
	tm.mu.RLock()
	// 60-second pre-expiration safety window
	if tm.token != "" && time.Now().Add(60*time.Second).Before(tm.expiresAt) {
		tok := tm.token
		tm.mu.RUnlock()
		return tok, nil
	}
	tm.mu.RUnlock()

	tm.refreshMu.Lock()
	defer tm.refreshMu.Unlock()

	// Double-check after acquiring write lock
	tm.mu.RLock()
	if tm.token != "" && time.Now().Add(60*time.Second).Before(tm.expiresAt) {
		tok := tm.token
		tm.mu.RUnlock()
		return tok, nil
	}
	tm.mu.RUnlock()

	tokenURL := fmt.Sprintf("%s/v1/oauth2/token", tm.apiBase)
	form := url.Values{}
	form.Set("grant_type", "client_credentials")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create oauth request: %w", err)
	}

	auth := base64.StdEncoding.EncodeToString([]byte(tm.clientID + ":" + tm.clientSecret))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := tm.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("oauth network failure: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("oauth token rejected with status %d", resp.StatusCode)
	}

	var res struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", fmt.Errorf("failed to decode oauth response: %w", err)
	}

	if res.AccessToken == "" {
		return "", errors.New("empty access token received from paypal")
	}

	tm.mu.Lock()
	tm.token = res.AccessToken
	tm.expiresAt = time.Now().Add(time.Duration(res.ExpiresIn) * time.Second)
	tm.mu.Unlock()

	return res.AccessToken, nil
}

// Invalidate invalidates cached token on 401
func (tm *TokenManager) Invalidate() {
	tm.mu.Lock()
	tm.token = ""
	tm.expiresAt = time.Time{}
	tm.mu.Unlock()
}
