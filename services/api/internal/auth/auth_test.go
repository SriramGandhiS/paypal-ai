package auth_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/paybridge/paybridge-ai/services/api/internal/auth"
)

func TestSessionTokenGeneration(t *testing.T) {
	t.Run("generate and hash session token", func(t *testing.T) {
		token, hashBytes, err := auth.GenerateSessionToken()
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if len(token) != 64 {
			t.Errorf("expected 64 hex chars, got %d", len(token))
		}
		if len(hashBytes) != 32 {
			t.Errorf("expected 32 byte hash, got %d", len(hashBytes))
		}

		manualHash := sha256.Sum256([]byte(token))
		if hex.EncodeToString(hashBytes) != hex.EncodeToString(manualHash[:]) {
			t.Errorf("hash mismatch between generated and manual sha256")
		}
	})

	t.Run("cookie roundtrip", func(t *testing.T) {
		w := httptest.NewRecorder()
		expires := time.Now().Add(auth.SessionDuration)
		auth.SetSessionCookie(w, "test-token-12345", expires)

		resp := w.Result()
		cookies := resp.Cookies()
		if len(cookies) == 0 {
			t.Fatal("expected session cookie in response")
		}
		if cookies[0].Value != "test-token-12345" {
			t.Errorf("expected token value 'test-token-12345', got %s", cookies[0].Value)
		}

		req := httptest.NewRequest("GET", "/", nil)
		req.AddCookie(cookies[0])

		extracted, err := auth.TokenFromRequest(req)
		if err != nil {
			t.Fatalf("unexpected error extracting token: %v", err)
		}
		if extracted != "test-token-12345" {
			t.Errorf("expected extracted token 'test-token-12345', got %s", extracted)
		}
	})
}

func TestSignAssertion(t *testing.T) {
	t.Run("generates valid HMAC assertion", func(t *testing.T) {
		key := []byte("secret-hmac-key-0000000000000000")
		tokenStr, sigHex, err := auth.SignAssertion(
			"key-01", key,
			"018f3a21-9b34-7a12-89cd-456789000000",
			"018f3a21-9b34-7a12-89cd-456789111111",
			"merchant_admin",
			"CREATE_CASE",
			"00000000-0000-0000-0000-000000000000",
			"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		)
		if err != nil {
			t.Fatalf("unexpected error signing assertion: %v", err)
		}

		var claims auth.AssertionClaims
		if err := json.Unmarshal([]byte(tokenStr), &claims); err != nil {
			t.Fatalf("failed to parse token JSON: %v", err)
		}

		if claims.KeyID != "key-01" {
			t.Errorf("expected key_id 'key-01', got %s", claims.KeyID)
		}

		canonical := claims.TenantID + ":" + claims.UserID + ":" + claims.Role + ":" +
			claims.Op + ":" + claims.ResourceID + ":" + claims.Nonce + ":" +
			claims.Expiry + ":" + "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

		mac := hmac.New(sha256.New, key)
		mac.Write([]byte(canonical))
		expectedSig := hex.EncodeToString(mac.Sum(nil))

		if sigHex != expectedSig {
			t.Errorf("signature verification failed: expected %s, got %s", expectedSig, sigHex)
		}
	})
}
