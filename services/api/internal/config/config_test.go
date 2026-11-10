package config_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/paybridge/paybridge-ai/services/api/internal/config"
)

// validEnv returns a map of environment variables containing valid development configuration.
func validEnv() map[string]string {
	return map[string]string{
		"DATABASE_URL":         "postgres://paybridge_admin:paybridge_dev_secret@localhost:5432/paybridge_dev?sslmode=disable",
		"SESSION_SECRET":       "dev_session_secret_change_in_production_min32chars",
		"ENCRYPTION_KEY":       "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
		"PAYPAL_ENV":           "sandbox",
		"PAYPAL_API_BASE":      "https://api-m.sandbox.paypal.com",
		"PAYPAL_CLIENT_ID":     "sandbox_client_id_test_12345",
		"PAYPAL_CLIENT_SECRET": "sandbox_client_secret_test_67890",
	}
}

// makeLookup creates a lookup function from a map.
func makeLookup(env map[string]string) func(string) string {
	return func(key string) string {
		return env[key]
	}
}

func TestValidConfiguration(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		modifyEnv     func(map[string]string)
		expectedKey   string
		expectedBase  string
		expectedDBUrl string
	}{
		{
			name: "Standard valid hex encryption key",
			modifyEnv: func(m map[string]string) {
				// defaults are valid hex
			},
			expectedKey:   "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
			expectedBase:  "https://api-m.sandbox.paypal.com",
			expectedDBUrl: "postgres://paybridge_admin:paybridge_dev_secret@localhost:5432/paybridge_dev?sslmode=disable",
		},
		{
			name: "Valid base64-encoded 32-byte encryption key (Doc 08 specification)",
			modifyEnv: func(m map[string]string) {
				// 32 bytes base64 encoded
				m["ENCRYPTION_KEY"] = "MDEyMzQ1Njc4OWFiY2RlZjAxMjM0NTY3ODlhYmNkZWY="
			},
			expectedKey:   "MDEyMzQ1Njc4OWFiY2RlZjAxMjM0NTY3ODlhYmNkZWY=",
			expectedBase:  "https://api-m.sandbox.paypal.com",
			expectedDBUrl: "postgres://paybridge_admin:paybridge_dev_secret@localhost:5432/paybridge_dev?sslmode=disable",
		},
		{
			name: "Valid postgresql:// scheme with trailing slash on API base",
			modifyEnv: func(m map[string]string) {
				m["DATABASE_URL"] = "postgresql://user:pass@db.example.com:5432/dbname"
				m["PAYPAL_API_BASE"] = "https://api-m.sandbox.paypal.com/"
			},
			expectedKey:   "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
			expectedBase:  "https://api-m.sandbox.paypal.com/",
			expectedDBUrl: "postgresql://user:pass@db.example.com:5432/dbname",
		},
		{
			name: "Uppercase canonical PayPal sandbox hostname normalized",
			modifyEnv: func(m map[string]string) {
				m["PAYPAL_API_BASE"] = "https://API-M.SANDBOX.PAYPAL.COM"
			},
			expectedKey:   "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
			expectedBase:  "https://API-M.SANDBOX.PAYPAL.COM",
			expectedDBUrl: "postgres://paybridge_admin:paybridge_dev_secret@localhost:5432/paybridge_dev?sslmode=disable",
		},
		{
			name: "Case-insensitive PAYPAL_ENV normalization",
			modifyEnv: func(m map[string]string) {
				m["PAYPAL_ENV"] = "  SANDBOX  "
			},
			expectedKey:   "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
			expectedBase:  "https://api-m.sandbox.paypal.com",
			expectedDBUrl: "postgres://paybridge_admin:paybridge_dev_secret@localhost:5432/paybridge_dev?sslmode=disable",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			env := validEnv()
			tc.modifyEnv(env)

			cfg, err := config.LoadFromEnv(makeLookup(env))
			if err != nil {
				t.Fatalf("expected valid config, got error: %v", err)
			}

			if cfg.PayPalEnv != "sandbox" {
				t.Errorf("expected PayPalEnv 'sandbox', got '%s'", cfg.PayPalEnv)
			}
			if cfg.PayPalAPIBase != tc.expectedBase {
				t.Errorf("expected PayPalAPIBase '%s', got '%s'", tc.expectedBase, cfg.PayPalAPIBase)
			}
			if cfg.EncryptionKey != tc.expectedKey {
				t.Errorf("expected EncryptionKey '%s', got '%s'", tc.expectedKey, cfg.EncryptionKey)
			}
			if cfg.DatabaseURL != tc.expectedDBUrl {
				t.Errorf("expected DatabaseURL '%s', got '%s'", tc.expectedDBUrl, cfg.DatabaseURL)
			}
		})
	}
}

func TestMissingConfiguration(t *testing.T) {
	t.Parallel()

	requiredVars := []string{
		"DATABASE_URL",
		"SESSION_SECRET",
		"ENCRYPTION_KEY",
		"PAYPAL_ENV",
		"PAYPAL_API_BASE",
		"PAYPAL_CLIENT_ID",
		"PAYPAL_CLIENT_SECRET",
	}

	for _, missingKey := range requiredVars {
		missingKey := missingKey
		t.Run("Missing "+missingKey, func(t *testing.T) {
			t.Parallel()
			env := validEnv()
			delete(env, missingKey)

			_, err := config.LoadFromEnv(makeLookup(env))
			if err == nil {
				t.Fatalf("expected error when %s is missing, got nil", missingKey)
			}

			errMsg := err.Error()
			if !strings.Contains(errMsg, missingKey) {
				t.Errorf("expected error message to mention missing field '%s', got: %s", missingKey, errMsg)
			}
			if !strings.Contains(errMsg, "is required") {
				t.Errorf("expected error message to indicate 'is required', got: %s", errMsg)
			}
		})
	}
}

func TestInvalidConfiguration(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		modifyEnv     func(map[string]string)
		expectedField string
		expectedMsg   string
	}{
		// Whitespace-only values
		{
			name: "Whitespace-only DATABASE_URL",
			modifyEnv: func(m map[string]string) {
				m["DATABASE_URL"] = "   "
			},
			expectedField: "DATABASE_URL",
			expectedMsg:   "is required",
		},
		{
			name: "Whitespace-only SESSION_SECRET",
			modifyEnv: func(m map[string]string) {
				m["SESSION_SECRET"] = "\t  \n "
			},
			expectedField: "SESSION_SECRET",
			expectedMsg:   "is required",
		},
		{
			name: "Whitespace-only ENCRYPTION_KEY",
			modifyEnv: func(m map[string]string) {
				m["ENCRYPTION_KEY"] = " "
			},
			expectedField: "ENCRYPTION_KEY",
			expectedMsg:   "is required",
		},
		{
			name: "Whitespace-only PAYPAL_CLIENT_ID",
			modifyEnv: func(m map[string]string) {
				m["PAYPAL_CLIENT_ID"] = " "
			},
			expectedField: "PAYPAL_CLIENT_ID",
			expectedMsg:   "is required",
		},
		{
			name: "Whitespace-only PAYPAL_CLIENT_SECRET",
			modifyEnv: func(m map[string]string) {
				m["PAYPAL_CLIENT_SECRET"] = " "
			},
			expectedField: "PAYPAL_CLIENT_SECRET",
			expectedMsg:   "is required",
		},

		// DATABASE_URL validation
		{
			name: "DATABASE_URL with invalid scheme (http)",
			modifyEnv: func(m map[string]string) {
				m["DATABASE_URL"] = "http://localhost:5432/db"
			},
			expectedField: "DATABASE_URL",
			expectedMsg:   "scheme must be postgres or postgresql",
		},
		{
			name: "DATABASE_URL missing host",
			modifyEnv: func(m map[string]string) {
				m["DATABASE_URL"] = "postgres:///dbname"
			},
			expectedField: "DATABASE_URL",
			expectedMsg:   "must specify a host",
		},

		// SESSION_SECRET validation
		{
			name: "SESSION_SECRET too short (< 32 chars)",
			modifyEnv: func(m map[string]string) {
				m["SESSION_SECRET"] = "short_secret_under_32_chars!"
			},
			expectedField: "SESSION_SECRET",
			expectedMsg:   "must be at least 32 characters long",
		},

		// ENCRYPTION_KEY validation
		{
			name: "ENCRYPTION_KEY invalid hex length (not 32 bytes)",
			modifyEnv: func(m map[string]string) {
				m["ENCRYPTION_KEY"] = "0123456789abcdef" // only 8 bytes
			},
			expectedField: "ENCRYPTION_KEY",
			expectedMsg:   "must be a valid 256-bit (32-byte) key in hex or base64 format",
		},
		{
			name: "ENCRYPTION_KEY non-hex characters in 64-char string",
			modifyEnv: func(m map[string]string) {
				m["ENCRYPTION_KEY"] = strings.Repeat("z", 64)
			},
			expectedField: "ENCRYPTION_KEY",
			expectedMsg:   "must be a valid 256-bit (32-byte) key in hex or base64 format",
		},
		{
			name: "ENCRYPTION_KEY raw 32-character ASCII string (rejected as ambiguous)",
			modifyEnv: func(m map[string]string) {
				m["ENCRYPTION_KEY"] = "12345678901234567890123456789012"
			},
			expectedField: "ENCRYPTION_KEY",
			expectedMsg:   "must be a valid 256-bit (32-byte) key in hex or base64 format",
		},

		// PAYPAL_ENV validation
		{
			name: "PAYPAL_ENV is live/production (prohibited)",
			modifyEnv: func(m map[string]string) {
				m["PAYPAL_ENV"] = "production"
			},
			expectedField: "PAYPAL_ENV",
			expectedMsg:   "must be exactly 'sandbox'",
		},
		{
			name: "PAYPAL_ENV is live",
			modifyEnv: func(m map[string]string) {
				m["PAYPAL_ENV"] = "live"
			},
			expectedField: "PAYPAL_ENV",
			expectedMsg:   "must be exactly 'sandbox'",
		},
		{
			name: "PAYPAL_ENV is arbitrary test",
			modifyEnv: func(m map[string]string) {
				m["PAYPAL_ENV"] = "test"
			},
			expectedField: "PAYPAL_ENV",
			expectedMsg:   "must be exactly 'sandbox'",
		},

		// PAYPAL_API_BASE validation
		{
			name: "PAYPAL_API_BASE is live production endpoint (PROHIBITED)",
			modifyEnv: func(m map[string]string) {
				m["PAYPAL_API_BASE"] = "https://api-m.paypal.com"
			},
			expectedField: "PAYPAL_API_BASE",
			expectedMsg:   "live PayPal host is strictly prohibited in sandbox mode",
		},
		{
			name: "PAYPAL_API_BASE live production endpoint with trailing dot (PROHIBITED)",
			modifyEnv: func(m map[string]string) {
				m["PAYPAL_API_BASE"] = "https://api-m.paypal.com."
			},
			expectedField: "PAYPAL_API_BASE",
			expectedMsg:   "live PayPal host is strictly prohibited in sandbox mode",
		},
		{
			name: "PAYPAL_API_BASE sandbox endpoint with trailing dot (rejected non-canonical)",
			modifyEnv: func(m map[string]string) {
				m["PAYPAL_API_BASE"] = "https://api-m.sandbox.paypal.com."
			},
			expectedField: "PAYPAL_API_BASE",
			expectedMsg:   "host must be 'api-m.sandbox.paypal.com'",
		},
		{
			name: "PAYPAL_API_BASE is insecure HTTP",
			modifyEnv: func(m map[string]string) {
				m["PAYPAL_API_BASE"] = "http://api-m.sandbox.paypal.com"
			},
			expectedField: "PAYPAL_API_BASE",
			expectedMsg:   "must use HTTPS protocol",
		},
		{
			name: "PAYPAL_API_BASE is arbitrary unknown host",
			modifyEnv: func(m map[string]string) {
				m["PAYPAL_API_BASE"] = "https://api.sandbox.paypal.com"
			},
			expectedField: "PAYPAL_API_BASE",
			expectedMsg:   "host must be 'api-m.sandbox.paypal.com'",
		},
		{
			name: "PAYPAL_API_BASE has embedded credentials",
			modifyEnv: func(m map[string]string) {
				m["PAYPAL_API_BASE"] = "https://user:pass@api-m.sandbox.paypal.com"
			},
			expectedField: "PAYPAL_API_BASE",
			expectedMsg:   "must not contain embedded user credentials",
		},
		{
			name: "PAYPAL_API_BASE specifies custom port",
			modifyEnv: func(m map[string]string) {
				m["PAYPAL_API_BASE"] = "https://api-m.sandbox.paypal.com:8443"
			},
			expectedField: "PAYPAL_API_BASE",
			expectedMsg:   "must not specify a custom port",
		},
		{
			name: "PAYPAL_API_BASE specifies default port explicitly (rejected per canonical rule)",
			modifyEnv: func(m map[string]string) {
				m["PAYPAL_API_BASE"] = "https://api-m.sandbox.paypal.com:443"
			},
			expectedField: "PAYPAL_API_BASE",
			expectedMsg:   "must not specify a custom port",
		},
		{
			name: "PAYPAL_API_BASE specifies query parameters",
			modifyEnv: func(m map[string]string) {
				m["PAYPAL_API_BASE"] = "https://api-m.sandbox.paypal.com?env=live"
			},
			expectedField: "PAYPAL_API_BASE",
			expectedMsg:   "must not contain query parameters",
		},
		{
			name: "PAYPAL_API_BASE specifies URL fragment",
			modifyEnv: func(m map[string]string) {
				m["PAYPAL_API_BASE"] = "https://api-m.sandbox.paypal.com#frag"
			},
			expectedField: "PAYPAL_API_BASE",
			expectedMsg:   "must not contain a URL fragment",
		},
		{
			name: "PAYPAL_API_BASE specifies path",
			modifyEnv: func(m map[string]string) {
				m["PAYPAL_API_BASE"] = "https://api-m.sandbox.paypal.com/v2/checkout"
			},
			expectedField: "PAYPAL_API_BASE",
			expectedMsg:   "must not specify a URL path",
		},
		{
			name: "PAYPAL_API_BASE specifies double slash path",
			modifyEnv: func(m map[string]string) {
				m["PAYPAL_API_BASE"] = "https://api-m.sandbox.paypal.com//"
			},
			expectedField: "PAYPAL_API_BASE",
			expectedMsg:   "must not specify a URL path",
		},
		{
			name: "PAYPAL_API_BASE is malformed URL",
			modifyEnv: func(m map[string]string) {
				m["PAYPAL_API_BASE"] = "https://invalid url with spaces"
			},
			expectedField: "PAYPAL_API_BASE",
			expectedMsg:   "must be a valid URL",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			env := validEnv()
			tc.modifyEnv(env)

			_, err := config.LoadFromEnv(makeLookup(env))
			if err == nil {
				t.Fatalf("expected validation error, got nil")
			}

			errMsg := err.Error()
			if !strings.Contains(errMsg, tc.expectedField) {
				t.Errorf("expected error to cite field '%s', got: %s", tc.expectedField, errMsg)
			}
			if !strings.Contains(errMsg, tc.expectedMsg) {
				t.Errorf("expected error to contain '%s', got: %s", tc.expectedMsg, errMsg)
			}
		})
	}
}

func TestZeroSecretExposureInValidationErrors(t *testing.T) {
	t.Parallel()

	canarySessionSecret := "CANARY_SESSION_SECRET_VAL_999888777666"
	canaryEncryptionKey := "CANARY_ENCRYPTION_KEY_VAL_111222333444"
	canaryClientSecret := "CANARY_PAYPAL_CLIENT_SECRET_555444333"

	// Trigger validation failures while canary secret values are present in the environment
	testCases := []struct {
		name      string
		modifyEnv func(map[string]string)
	}{
		{
			name: "Missing DATABASE_URL with secrets present",
			modifyEnv: func(m map[string]string) {
				delete(m, "DATABASE_URL")
				m["SESSION_SECRET"] = canarySessionSecret
				m["ENCRYPTION_KEY"] = canaryEncryptionKey
				m["PAYPAL_CLIENT_SECRET"] = canaryClientSecret
			},
		},
		{
			name: "Invalid PAYPAL_API_BASE (live endpoint) with secrets present",
			modifyEnv: func(m map[string]string) {
				m["PAYPAL_API_BASE"] = "https://api-m.paypal.com"
				m["SESSION_SECRET"] = canarySessionSecret
				m["ENCRYPTION_KEY"] = canaryEncryptionKey
				m["PAYPAL_CLIENT_SECRET"] = canaryClientSecret
			},
		},
		{
			name: "Invalid PAYPAL_ENV (production) with secrets present",
			modifyEnv: func(m map[string]string) {
				m["PAYPAL_ENV"] = "production"
				m["SESSION_SECRET"] = canarySessionSecret
				m["ENCRYPTION_KEY"] = canaryEncryptionKey
				m["PAYPAL_CLIENT_SECRET"] = canaryClientSecret
			},
		},
		{
			name: "Invalid SESSION_SECRET length with secret present",
			modifyEnv: func(m map[string]string) {
				m["SESSION_SECRET"] = "CANARY_SHORT"
				m["ENCRYPTION_KEY"] = canaryEncryptionKey
				m["PAYPAL_CLIENT_SECRET"] = canaryClientSecret
			},
		},
		{
			name: "Invalid ENCRYPTION_KEY format with secret present",
			modifyEnv: func(m map[string]string) {
				m["SESSION_SECRET"] = canarySessionSecret
				m["ENCRYPTION_KEY"] = canaryEncryptionKey // not valid 32-byte hex/base64
				m["PAYPAL_CLIENT_SECRET"] = canaryClientSecret
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			env := validEnv()
			tc.modifyEnv(env)

			_, err := config.LoadFromEnv(makeLookup(env))
			if err == nil {
				t.Fatalf("expected error, got nil")
			}

			errStr := err.Error()

			if strings.Contains(errStr, canarySessionSecret) {
				t.Errorf("CRITICAL SECURITY DEFECT: SESSION_SECRET value leaked in error message: %s", errStr)
			}
			if strings.Contains(errStr, "CANARY_SHORT") {
				t.Errorf("CRITICAL SECURITY DEFECT: Short SESSION_SECRET value leaked in error message: %s", errStr)
			}
			if strings.Contains(errStr, canaryEncryptionKey) {
				t.Errorf("CRITICAL SECURITY DEFECT: ENCRYPTION_KEY value leaked in error message: %s", errStr)
			}
			if strings.Contains(errStr, canaryClientSecret) {
				t.Errorf("CRITICAL SECURITY DEFECT: PAYPAL_CLIENT_SECRET value leaked in error message: %s", errStr)
			}
		})
	}
}

func TestConfigFormattingAndSerializationSecretsRedacted(t *testing.T) {
	t.Parallel()

	canaryPass := "super_secret_db_password_12345"
	canarySession := "session_secret_exact_32_chars_now!!"
	canaryKey := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	canaryClientSecret := "paypal_client_secret_xyz987"

	cfg := config.Config{
		DatabaseURL:        fmt.Sprintf("postgres://paybridge_admin:%s@localhost:5432/paybridge_dev?sslmode=disable", canaryPass),
		SessionSecret:      canarySession,
		EncryptionKey:      canaryKey,
		PayPalEnv:          "sandbox",
		PayPalAPIBase:      "https://api-m.sandbox.paypal.com",
		PayPalClientID:     "client_id_pub",
		PayPalClientSecret: canaryClientSecret,
	}

	// 1. fmt.Sprintf("%v", cfg)
	sDefault := fmt.Sprintf("%v", cfg)
	assertNoSecrets(t, "fmt.Sprintf(%v)", sDefault, canaryPass, canarySession, canaryKey, canaryClientSecret)
	if !strings.Contains(sDefault, config.RedactedPlaceholder) {
		t.Errorf("expected redacted placeholder in %%v output, got: %s", sDefault)
	}

	// 2. fmt.Sprintf("%+v", cfg)
	sDetailed := fmt.Sprintf("%+v", cfg)
	assertNoSecrets(t, "fmt.Sprintf(%+v)", sDetailed, canaryPass, canarySession, canaryKey, canaryClientSecret)

	// 3. fmt.Sprintf("%#v", cfg) (GoStringer)
	sGoDetailed := fmt.Sprintf("%#v", cfg)
	assertNoSecrets(t, "fmt.Sprintf(%#v)", sGoDetailed, canaryPass, canarySession, canaryKey, canaryClientSecret)

	// 4. json.Marshal(cfg)
	jsonBytes, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}
	sJSON := string(jsonBytes)
	assertNoSecrets(t, "json.Marshal", sJSON, canaryPass, canarySession, canaryKey, canaryClientSecret)

	var unmarshaled map[string]string
	if err := json.Unmarshal(jsonBytes, &unmarshaled); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}
	if unmarshaled["session_secret"] != config.RedactedPlaceholder {
		t.Errorf("expected session_secret to be redacted in JSON, got: %s", unmarshaled["session_secret"])
	}
	if unmarshaled["encryption_key"] != config.RedactedPlaceholder {
		t.Errorf("expected encryption_key to be redacted in JSON, got: %s", unmarshaled["encryption_key"])
	}
	if unmarshaled["paypal_client_secret"] != config.RedactedPlaceholder {
		t.Errorf("expected paypal_client_secret to be redacted in JSON, got: %s", unmarshaled["paypal_client_secret"])
	}
	if strings.Contains(unmarshaled["database_url"], canaryPass) {
		t.Errorf("CRITICAL: Database password found in JSON: %s", unmarshaled["database_url"])
	}
}

func assertNoSecrets(t *testing.T, context, output, canaryPass, canarySession, canaryKey, canaryClientSecret string) {
	t.Helper()
	if strings.Contains(output, canaryPass) {
		t.Errorf("CRITICAL LEAK in %s: Database password leaked: %s", context, output)
	}
	if strings.Contains(output, canarySession) {
		t.Errorf("CRITICAL LEAK in %s: Session secret leaked: %s", context, output)
	}
	if strings.Contains(output, canaryKey) {
		t.Errorf("CRITICAL LEAK in %s: Encryption key leaked: %s", context, output)
	}
	if strings.Contains(output, canaryClientSecret) {
		t.Errorf("CRITICAL LEAK in %s: PayPal client secret leaked: %s", context, output)
	}
}

func TestConfigErrorEmpty(t *testing.T) {
	t.Parallel()
	emptyErr := &config.ConfigError{}
	if emptyErr.Error() != "configuration validation failed" {
		t.Errorf("expected default message for empty ConfigError, got: %s", emptyErr.Error())
	}
}

func TestLoadWithOSEnv(t *testing.T) {
	env := validEnv()
	for k, v := range env {
		t.Setenv(k, v)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("expected Load() to succeed with valid OS environment, got error: %v", err)
	}

	if cfg.PayPalEnv != "sandbox" {
		t.Errorf("expected PayPalEnv 'sandbox', got '%s'", cfg.PayPalEnv)
	}
	if cfg.PayPalAPIBase != "https://api-m.sandbox.paypal.com" {
		t.Errorf("expected canonical sandbox base, got '%s'", cfg.PayPalAPIBase)
	}
}
