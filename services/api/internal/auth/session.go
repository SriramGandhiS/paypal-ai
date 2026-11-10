package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	SessionCookieName = "paybridge_session"
	SessionDuration   = 24 * time.Hour
	AssertionTTL      = 60 * time.Second
)

var (
	ErrInvalidSessionToken = errors.New("invalid or malformed session token")
	ErrSessionExpired      = errors.New("session token has expired")
	ErrSignatureMismatch   = errors.New("HMAC assertion signature mismatch")
)

// Session represents an authenticated user session in PayBridge AI.
type Session struct {
	IDHash    []byte
	UserID    string
	TenantID  string
	Role      string
	ExpiresAt time.Time
	RevokedAt *time.Time
}

// GenerateSessionToken creates a cryptographically secure random session token.
func GenerateSessionToken() (string, []byte, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", nil, fmt.Errorf("failed to generate random token: %w", err)
	}
	token := hex.EncodeToString(raw)
	h := sha256.Sum256([]byte(token))
	return token, h[:], nil
}

// HashToken computes the SHA-256 digest of a session token string.
func HashToken(token string) []byte {
	h := sha256.Sum256([]byte(token))
	return h[:]
}

// SetSessionCookie attaches the HttpOnly, Secure, SameSite session cookie to response.
func SetSessionCookie(w http.ResponseWriter, token string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    token,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   false, // Set to true in TLS/production mode
		SameSite: http.SameSiteLaxMode,
	})
}

// ClearSessionCookie removes the session cookie from the client.
func ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

// TokenFromRequest extracts the session token from the HTTP request cookie.
func TokenFromRequest(r *http.Request) (string, error) {
	c, err := r.Cookie(SessionCookieName)
	if err != nil {
		return "", err
	}
	if c.Value == "" {
		return "", ErrInvalidSessionToken
	}
	return c.Value, nil
}

// AssertionClaims represents the signed payload passed to PostgreSQL stored procedures.
type AssertionClaims struct {
	KeyID      string `json:"key_id"`
	TenantID   string `json:"tenant_id"`
	UserID     string `json:"user_id"`
	Role       string `json:"role"`
	Op         string `json:"op"`
	ResourceID string `json:"resource_id"`
	Nonce      string `json:"nonce"`
	Expiry     string `json:"expiry"`
}

// SignAssertion creates a canonical JSON token and hex-encoded HMAC-SHA256 signature.
func SignAssertion(keyID string, keyMaterial []byte, tenantID, userID, role, op, resourceID, argsHash string) (string, string, error) {
	nonceBytes := make([]byte, 16)
	if _, err := rand.Read(nonceBytes); err != nil {
		return "", "", fmt.Errorf("failed to generate nonce: %w", err)
	}
	nonce := hex.EncodeToString(nonceBytes)
	expiryTime := time.Now().UTC().Add(AssertionTTL)
	expiryStr := expiryTime.Format(time.RFC3339)

	claims := AssertionClaims{
		KeyID:      keyID,
		TenantID:   tenantID,
		UserID:     userID,
		Role:       role,
		Op:         op,
		ResourceID: resourceID,
		Nonce:      nonce,
		Expiry:     expiryStr,
	}

	tokenBytes, err := json.Marshal(claims)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal assertion claims: %w", err)
	}
	tokenStr := string(tokenBytes)

	// Canonical assertion string matching PostgreSQL _fn_verify_and_consume_assertion:
	// tenant_id:user_id:role:op:resource_id:nonce:expiry:args_hash
	canonical := fmt.Sprintf("%s:%s:%s:%s:%s:%s:%s:%s",
		tenantID, userID, role, op, resourceID, nonce, expiryStr, argsHash)

	mac := hmac.New(sha256.New, keyMaterial)
	mac.Write([]byte(canonical))
	sigHex := hex.EncodeToString(mac.Sum(nil))

	return tokenStr, sigHex, nil
}
