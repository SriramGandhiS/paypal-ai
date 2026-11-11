package kms_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/paybridge/paybridge-ai/services/api/internal/kms"
)

func TestEnvelopeEncryption_RoundTrip(t *testing.T) {
	client := kms.NewClient([]byte("master-secret-key-32-bytes-long!"), "2026-Q1")
	ctx := context.Background()

	originalPayload := []byte(`{"order_id":"ORD-9912","buyer":"payer@example.com","secret_pan_token":"tok_8812"}`)

	env, err := client.EncryptPayload(ctx, originalPayload)
	if err != nil {
		t.Fatalf("encryption failed: %v", err)
	}

	if bytes.Equal(env.Ciphertext, originalPayload) {
		t.Fatal("ciphertext must not match plaintext payload")
	}

	decrypted, err := client.DecryptPayload(ctx, env)
	if err != nil {
		t.Fatalf("decryption failed: %v", err)
	}

	if !bytes.Equal(decrypted, originalPayload) {
		t.Fatalf("decrypted payload mismatch: got %s, want %s", string(decrypted), string(originalPayload))
	}
}

func TestEnvelopeEncryption_KeyRotation(t *testing.T) {
	client := kms.NewClient([]byte("master-secret-key-2026-Q1-seed!"), "2026-Q1")
	ctx := context.Background()

	payloadQ1 := []byte("Sensitive Q1 Evidence Record")

	// Encrypt under Q1 key
	envQ1, err := client.EncryptPayload(ctx, payloadQ1)
	if err != nil {
		t.Fatalf("Q1 encryption failed: %v", err)
	}

	// 90-day key rotation to Q2
	client.AddKeyVersion("2026-Q2", []byte("master-secret-key-2026-Q2-seed!"), true)

	payloadQ2 := []byte("Sensitive Q2 Evidence Record")
	envQ2, err := client.EncryptPayload(ctx, payloadQ2)
	if err != nil {
		t.Fatalf("Q2 encryption failed: %v", err)
	}

	if envQ2.KeyVersion != "2026-Q2" {
		t.Errorf("expected Q2 key version on new writes, got %s", envQ2.KeyVersion)
	}

	// Verify old Q1 payload can still be decrypted seamlessly via multi-version keyring
	decryptedQ1, err := client.DecryptPayload(ctx, envQ1)
	if err != nil {
		t.Fatalf("Q1 decryption after rotation failed: %v", err)
	}
	if !bytes.Equal(decryptedQ1, payloadQ1) {
		t.Errorf("decrypted Q1 mismatch: got %s, want %s", string(decryptedQ1), string(payloadQ1))
	}

	// Verify new Q2 payload decrypts
	decryptedQ2, err := client.DecryptPayload(ctx, envQ2)
	if err != nil {
		t.Fatalf("Q2 decryption failed: %v", err)
	}
	if !bytes.Equal(decryptedQ2, payloadQ2) {
		t.Errorf("decrypted Q2 mismatch: got %s, want %s", string(decryptedQ2), string(payloadQ2))
	}
}

func TestEnvelopeEncryption_TamperDetection(t *testing.T) {
	client := kms.NewClient([]byte("master-secret-key-32-bytes-long!"), "v1")
	ctx := context.Background()

	env, err := client.EncryptPayload(ctx, []byte("Immutable Financial Evidence"))
	if err != nil {
		t.Fatalf("encryption failed: %v", err)
	}

	// Tamper with one byte in the ciphertext
	env.Ciphertext[0] ^= 0xFF

	_, err = client.DecryptPayload(ctx, env)
	if err == nil {
		t.Fatal("expected decryption error on tampered ciphertext, got nil")
	}
}

func TestEnvelopeEncryption_AWSKMSProviderInterface(t *testing.T) {
	awsProvider := kms.NewAWSKMSProvider(kms.AWSKMSConfig{
		Region: "us-east-1",
		KeyARN: "arn:aws:kms:us-east-1:123456789012:key/12345678-1234-1234-1234-123456789012",
	})

	if awsProvider.Type() != kms.ProviderAWS {
		t.Fatalf("expected provider type %s, got %s", kms.ProviderAWS, awsProvider.Type())
	}

	client := kms.NewClientWithProvider(awsProvider, "arn:aws:kms:us-east-1:123456789012:key/default")
	ctx := context.Background()

	data := []byte(`{"account":"12345678","secret":"sandbox-test"}`)
	env, err := client.EncryptPayload(ctx, data)
	if err != nil {
		t.Fatalf("AWS KMS provider encryption failed: %v", err)
	}

	decrypted, err := client.DecryptPayload(ctx, env)
	if err != nil {
		t.Fatalf("AWS KMS provider decryption failed: %v", err)
	}

	if !bytes.Equal(decrypted, data) {
		t.Fatalf("AWS KMS round-trip mismatch: got %s, want %s", string(decrypted), string(data))
	}
}

func TestEnvelopeEncryption_VaultTransitProviderInterface(t *testing.T) {
	vaultProvider := kms.NewVaultTransitProvider(kms.VaultConfig{
		Address:    "http://vault.internal:8200",
		TransitKey: "paybridge-evidence",
	})

	if vaultProvider.Type() != kms.ProviderVault {
		t.Fatalf("expected provider type %s, got %s", kms.ProviderVault, vaultProvider.Type())
	}

	client := kms.NewClientWithProvider(vaultProvider, "transit/keys/paybridge-evidence")
	ctx := context.Background()

	data := []byte(`{"transit_token":"vault-evidence-secure"}`)
	env, err := client.EncryptPayload(ctx, data)
	if err != nil {
		t.Fatalf("Vault provider encryption failed: %v", err)
	}

	decrypted, err := client.DecryptPayload(ctx, env)
	if err != nil {
		t.Fatalf("Vault provider decryption failed: %v", err)
	}

	if !bytes.Equal(decrypted, data) {
		t.Fatalf("Vault round-trip mismatch: got %s, want %s", string(decrypted), string(data))
	}
}
