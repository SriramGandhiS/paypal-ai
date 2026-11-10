package config

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"
)

// CanonicalPayPalSandboxHost is the only approved PayPal host for the MVP.
const CanonicalPayPalSandboxHost = "api-m.sandbox.paypal.com"

// ProhibitedLivePayPalHost is the live PayPal production host.
const ProhibitedLivePayPalHost = "api-m.paypal.com"

// MinSessionSecretLength defines the minimum character length required for session signing.
const MinSessionSecretLength = 32

// RedactedPlaceholder masks sensitive credentials in formatted representations.
const RedactedPlaceholder = "[REDACTED]"

// Config contains the validated application configuration.
// Sensitive fields must never be serialized or printed in logs/errors.
type Config struct {
	DatabaseURL        string
	SessionSecret      string
	EncryptionKey      string
	PayPalEnv          string
	PayPalAPIBase      string
	PayPalClientID     string
	PayPalClientSecret string
	AppEnv             string
	AllowDevAuth       bool
}

// String implements fmt.Stringer to prevent accidental credential leakage in logs and formatters.
func (c Config) String() string {
	safeDB := c.DatabaseURL
	if u, err := url.Parse(c.DatabaseURL); err == nil && u.User != nil {
		safeDB = u.Redacted()
	}
	return fmt.Sprintf("Config{DatabaseURL:%s, SessionSecret:%s, EncryptionKey:%s, PayPalEnv:%s, PayPalAPIBase:%s, PayPalClientID:%s, PayPalClientSecret:%s, AppEnv:%s, AllowDevAuth:%t}",
		safeDB, RedactedPlaceholder, RedactedPlaceholder, c.PayPalEnv, c.PayPalAPIBase, c.PayPalClientID, RedactedPlaceholder, c.AppEnv, c.AllowDevAuth)
}

// GoString implements fmt.GoStringer to prevent credential leakage via %#v format specifiers.
func (c Config) GoString() string {
	return c.String()
}

// MarshalJSON implements json.Marshaler to prevent accidental credential leakage during JSON encoding.
func (c Config) MarshalJSON() ([]byte, error) {
	safeDB := c.DatabaseURL
	if u, err := url.Parse(c.DatabaseURL); err == nil && u.User != nil {
		safeDB = u.Redacted()
	}
	type safeConfig struct {
		DatabaseURL        string `json:"database_url"`
		SessionSecret      string `json:"session_secret"`
		EncryptionKey      string `json:"encryption_key"`
		PayPalEnv          string `json:"paypal_env"`
		PayPalAPIBase      string `json:"paypal_api_base"`
		PayPalClientID     string `json:"paypal_client_id"`
		PayPalClientSecret string `json:"paypal_client_secret"`
		AppEnv             string `json:"app_env"`
		AllowDevAuth       string `json:"allow_dev_auth"`
	}
	return json.Marshal(safeConfig{
		DatabaseURL:        safeDB,
		SessionSecret:      RedactedPlaceholder,
		EncryptionKey:      RedactedPlaceholder,
		PayPalEnv:          c.PayPalEnv,
		PayPalAPIBase:      c.PayPalAPIBase,
		PayPalClientID:     c.PayPalClientID,
		PayPalClientSecret: RedactedPlaceholder,
		AppEnv:             c.AppEnv,
		AllowDevAuth:       fmt.Sprintf("%t", c.AllowDevAuth),
	})
}

// ValidationError represents an individual configuration validation failure.
// Invariant: Message and Field must never contain secret values.
type ValidationError struct {
	Field   string
	Message string
}

// ConfigError represents one or more configuration validation failures.
type ConfigError struct {
	Errors []ValidationError
}

func (e *ConfigError) Error() string {
	if len(e.Errors) == 0 {
		return "configuration validation failed"
	}
	var sb strings.Builder
	sb.WriteString("configuration validation failed: ")
	for i, err := range e.Errors {
		if i > 0 {
			sb.WriteString("; ")
		}
		sb.WriteString(err.Field)
		sb.WriteString(" ")
		sb.WriteString(err.Message)
	}
	return sb.String()
}

// Load loads configuration from environment variables and performs strict validation.
// It returns an error if any required setting is missing or invalid.
// Invariant: Secret values are never included in returned error messages.
func Load() (Config, error) {
	return LoadFromEnv(os.Getenv)
}

// LoadFromEnv loads and validates configuration using a custom environment lookup function.
// This enables deterministic unit testing without mutating OS environment variables.
func LoadFromEnv(lookup func(string) string) (Config, error) {
	var valErrors []ValidationError

	addErr := func(field, msg string) {
		valErrors = append(valErrors, ValidationError{Field: field, Message: msg})
	}

	// 1. DATABASE_URL
	rawDB := strings.TrimSpace(lookup("DATABASE_URL"))
	if rawDB == "" {
		addErr("DATABASE_URL", "is required")
	} else {
		u, err := url.Parse(rawDB)
		if err != nil || u.Scheme == "" {
			addErr("DATABASE_URL", "must be a valid URL")
		} else if u.Scheme != "postgres" && u.Scheme != "postgresql" {
			addErr("DATABASE_URL", "scheme must be postgres or postgresql")
		} else if u.Host == "" {
			addErr("DATABASE_URL", "must specify a host")
		}
	}

	// 2. SESSION_SECRET
	rawSession := strings.TrimSpace(lookup("SESSION_SECRET"))
	if rawSession == "" {
		addErr("SESSION_SECRET", "is required")
	} else if len(rawSession) < MinSessionSecretLength {
		addErr("SESSION_SECRET", fmt.Sprintf("must be at least %d characters long", MinSessionSecretLength))
	}

	// 3. ENCRYPTION_KEY
	rawKey := strings.TrimSpace(lookup("ENCRYPTION_KEY"))
	if rawKey == "" {
		addErr("ENCRYPTION_KEY", "is required")
	} else if !isValidEncryptionKey(rawKey) {
		addErr("ENCRYPTION_KEY", "must be a valid 256-bit (32-byte) key in hex or base64 format")
	}

	// 4. PAYPAL_ENV (Mandatory Sandbox Guard)
	rawEnv := strings.TrimSpace(strings.ToLower(lookup("PAYPAL_ENV")))
	if rawEnv == "" {
		addErr("PAYPAL_ENV", "is required")
	} else if rawEnv != "sandbox" {
		addErr("PAYPAL_ENV", "must be exactly 'sandbox'")
	}

	// 5. PAYPAL_API_BASE (Strict Endpoint Allowlist)
	rawAPIBase := strings.TrimSpace(lookup("PAYPAL_API_BASE"))
	if rawAPIBase == "" {
		addErr("PAYPAL_API_BASE", "is required")
	} else {
		u, err := url.Parse(rawAPIBase)
		if err != nil || u.Scheme == "" {
			addErr("PAYPAL_API_BASE", "must be a valid URL")
		} else if u.Scheme != "https" {
			addErr("PAYPAL_API_BASE", "must use HTTPS protocol")
		} else if u.User != nil {
			addErr("PAYPAL_API_BASE", "must not contain embedded user credentials")
		} else if u.Port() != "" {
			addErr("PAYPAL_API_BASE", "must not specify a custom port")
		} else if u.RawQuery != "" {
			addErr("PAYPAL_API_BASE", "must not contain query parameters")
		} else if u.Fragment != "" {
			addErr("PAYPAL_API_BASE", "must not contain a URL fragment")
		} else if u.Path != "" && u.Path != "/" {
			addErr("PAYPAL_API_BASE", "must not specify a URL path")
		} else {
			rawHost := strings.ToLower(u.Hostname())
			normalizedHost := strings.TrimSuffix(rawHost, ".")

			if normalizedHost == ProhibitedLivePayPalHost {
				addErr("PAYPAL_API_BASE", "live PayPal host is strictly prohibited in sandbox mode")
			} else if rawHost != CanonicalPayPalSandboxHost {
				addErr("PAYPAL_API_BASE", fmt.Sprintf("host must be '%s'", CanonicalPayPalSandboxHost))
			}
		}
	}

	// 6. PAYPAL_CLIENT_ID
	rawClientID := strings.TrimSpace(lookup("PAYPAL_CLIENT_ID"))
	if rawClientID == "" {
		addErr("PAYPAL_CLIENT_ID", "is required")
	}

	// 7. PAYPAL_CLIENT_SECRET
	rawClientSecret := strings.TrimSpace(lookup("PAYPAL_CLIENT_SECRET"))
	if rawClientSecret == "" {
		addErr("PAYPAL_CLIENT_SECRET", "is required")
	}

	// 8. APP_ENV (Optional, defaults to development)
	appEnv := strings.TrimSpace(strings.ToLower(lookup("APP_ENV")))
	if appEnv == "" {
		appEnv = "development"
	}

	// 9. ALLOW_DEV_AUTH (Strictly prohibited in production)
	allowDevAuth := false
	rawDevAuth := strings.TrimSpace(strings.ToLower(lookup("ALLOW_DEV_AUTH")))
	if rawDevAuth == "true" || rawDevAuth == "1" {
		if appEnv == "production" {
			addErr("ALLOW_DEV_AUTH", "cannot be enabled in production environment")
		} else {
			allowDevAuth = true
		}
	}

	if len(valErrors) > 0 {
		return Config{}, &ConfigError{Errors: valErrors}
	}

	return Config{
		DatabaseURL:        rawDB,
		SessionSecret:      rawSession,
		EncryptionKey:      rawKey,
		PayPalEnv:          rawEnv,
		PayPalAPIBase:      rawAPIBase,
		PayPalClientID:     rawClientID,
		PayPalClientSecret: rawClientSecret,
		AppEnv:             appEnv,
		AllowDevAuth:       allowDevAuth,
	}, nil
}

// isValidEncryptionKey checks if a string decodes to an exact 256-bit (32-byte) key.
// Accepted formats: 64 hex characters or standard base64 encoding of 32 bytes.
func isValidEncryptionKey(key string) bool {
	// 64-char hex encoding (32 bytes)
	if len(key) == 64 {
		decoded, err := hex.DecodeString(key)
		if err == nil && len(decoded) == 32 {
			return true
		}
	}

	// Base64 encoding of 32 bytes (typically 44 chars with padding)
	decodedBase64, err := base64.StdEncoding.DecodeString(key)
	if err == nil && len(decodedBase64) == 32 {
		return true
	}

	return false
}
