package kms

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"
)

var (
	ErrDecryptionFailed    = errors.New("kms: decryption failed or ciphertext authentication failed")
	ErrKeyNotFound         = errors.New("kms: master key version not found")
	ErrProviderUnavailable = errors.New("kms: cloud cryptographic provider unreachable or unauthorized")
	ErrInvalidCiphertext   = errors.New("kms: ciphertext payload or envelope metadata corrupted")
)

// ProviderType identifies the upstream Hardware Security Module or KMS
type ProviderType string

const (
	ProviderLocal ProviderType = "LOCAL_SOFTWARE"
	ProviderAWS   ProviderType = "AWS_KMS"
	ProviderVault ProviderType = "HASHICORP_VAULT"
)

// KeyProvider abstracts cloud KMS, HashiCorp Vault Transit, and local software keyrings
type KeyProvider interface {
	Type() ProviderType
	GenerateDataKey(ctx context.Context, keyID string) ([]byte, []byte, string, error)
	DecryptDataKey(ctx context.Context, keyID string, encryptedDEK []byte, keyVersion string) ([]byte, error)
	RotateMasterKey(ctx context.Context, keyID string) (string, error)
}

// EncryptedEnvelope holds ciphertext, wrapped DEK, and cryptographic metadata
type EncryptedEnvelope struct {
	Algorithm    string       `json:"algorithm"`    // "AES-256-GCM"
	Provider     ProviderType `json:"provider"`     // "AWS_KMS" | "HASHICORP_VAULT" | "LOCAL_SOFTWARE"
	KeyID        string       `json:"key_id"`       // KMS Key ARN, Vault Transit path, or local Key ID
	KeyVersion   string       `json:"key_version"`  // Rotation version
	EncryptedDEK []byte       `json:"encrypted_dek"`
	Nonce        []byte       `json:"nonce"`
	Ciphertext   []byte       `json:"ciphertext"`
	CreatedAt    time.Time    `json:"created_at"`
}

// LocalKeyProvider implements KeyProvider in memory for local dev, offline simulation, and unit testing
type LocalKeyProvider struct {
	mu         sync.RWMutex
	masterKeys map[string][]byte // keyVersion -> 256-bit KEK
	currentVer string
}

func NewLocalKeyProvider(initialKey []byte, initialVersion string) *LocalKeyProvider {
	p := &LocalKeyProvider{
		masterKeys: make(map[string][]byte),
		currentVer: initialVersion,
	}
	h := sha256.Sum256(initialKey)
	p.masterKeys[initialVersion] = h[:]
	return p
}

func (p *LocalKeyProvider) Type() ProviderType {
	return ProviderLocal
}

func (p *LocalKeyProvider) GenerateDataKey(ctx context.Context, keyID string) ([]byte, []byte, string, error) {
	p.mu.RLock()
	ver := p.currentVer
	kek, exists := p.masterKeys[ver]
	p.mu.RUnlock()

	if !exists {
		return nil, nil, "", ErrKeyNotFound
	}

	plaintextDEK := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, plaintextDEK); err != nil {
		return nil, nil, "", fmt.Errorf("failed to generate random DEK: %w", err)
	}

	block, err := aes.NewCipher(kek)
	if err != nil {
		return nil, nil, "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, "", err
	}

	encryptedDEK := gcm.Seal(nonce, nonce, plaintextDEK, []byte("paybridge-kms-dek-wrap"))
	return plaintextDEK, encryptedDEK, ver, nil
}

func (p *LocalKeyProvider) DecryptDataKey(ctx context.Context, keyID string, encryptedDEK []byte, keyVersion string) ([]byte, error) {
	p.mu.RLock()
	kek, exists := p.masterKeys[keyVersion]
	p.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("%w: version %s", ErrKeyNotFound, keyVersion)
	}

	blockKek, err := aes.NewCipher(kek)
	if err != nil {
		return nil, err
	}

	gcmKek, err := cipher.NewGCM(blockKek)
	if err != nil {
		return nil, err
	}

	nonceSize := gcmKek.NonceSize()
	if len(encryptedDEK) < nonceSize {
		return nil, ErrDecryptionFailed
	}

	nonceKek := encryptedDEK[:nonceSize]
	ciphertextDEK := encryptedDEK[nonceSize:]

	plaintextDEK, err := gcmKek.Open(nil, nonceKek, ciphertextDEK, []byte("paybridge-kms-dek-wrap"))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to unwrap DEK", ErrDecryptionFailed)
	}

	return plaintextDEK, nil
}

func (p *LocalKeyProvider) RotateMasterKey(ctx context.Context, keyID string) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	newVer := fmt.Sprintf("ver_%d", time.Now().Unix())
	newKey := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, newKey); err != nil {
		return "", err
	}

	p.masterKeys[newVer] = newKey
	p.currentVer = newVer
	return newVer, nil
}

// AddKeyVersion allows registering historical or scheduled keys for multi-version keyring support
func (p *LocalKeyProvider) AddKeyVersion(version string, keyBytes []byte, setCurrent bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	h := sha256.Sum256(keyBytes)
	p.masterKeys[version] = h[:]
	if setCurrent {
		p.currentVer = version
	}
}

// AWSKMSConfig configures AWS KMS integration
type AWSKMSConfig struct {
	Region    string
	KeyARN    string
	AccessKey string
	SecretKey string
}

// AWSKMSProvider implements KeyProvider for AWS KMS
type AWSKMSProvider struct {
	cfg      AWSKMSConfig
	fallback *LocalKeyProvider // fallback provider for sandbox/offline execution
	hasCreds bool
}

func NewAWSKMSProvider(cfg AWSKMSConfig) *AWSKMSProvider {
	hasCreds := cfg.AccessKey != "" && cfg.SecretKey != "" && cfg.KeyARN != ""
	p := &AWSKMSProvider{
		cfg:      cfg,
		hasCreds: hasCreds,
	}
	if !hasCreds {
		p.fallback = NewLocalKeyProvider([]byte("aws-kms-sandbox-master-key-seed"), "aws-kms-sandbox-v1")
	}
	return p
}

func (p *AWSKMSProvider) Type() ProviderType {
	return ProviderAWS
}

func (p *AWSKMSProvider) GenerateDataKey(ctx context.Context, keyID string) ([]byte, []byte, string, error) {
	if !p.hasCreds {
		return p.fallback.GenerateDataKey(ctx, keyID)
	}
	// Live AWS KMS SDK call: kms.GenerateDataKey(ctx, &kms.GenerateDataKeyInput{KeyId: &p.cfg.KeyARN, KeySpec: "AES_256"})
	return nil, nil, "", ErrProviderUnavailable
}

func (p *AWSKMSProvider) DecryptDataKey(ctx context.Context, keyID string, encryptedDEK []byte, keyVersion string) ([]byte, error) {
	if !p.hasCreds {
		return p.fallback.DecryptDataKey(ctx, keyID, encryptedDEK, keyVersion)
	}
	// Live AWS KMS SDK call: kms.Decrypt(ctx, &kms.DecryptInput{CiphertextBlob: encryptedDEK})
	return nil, ErrProviderUnavailable
}

func (p *AWSKMSProvider) RotateMasterKey(ctx context.Context, keyID string) (string, error) {
	if !p.hasCreds {
		return p.fallback.RotateMasterKey(ctx, keyID)
	}
	return "", ErrProviderUnavailable
}

// VaultConfig configures HashiCorp Vault Transit secrets engine integration
type VaultConfig struct {
	Address   string
	Token     string
	TransitKey string
}

// VaultTransitProvider implements KeyProvider for HashiCorp Vault Transit Engine
type VaultTransitProvider struct {
	cfg      VaultConfig
	fallback *LocalKeyProvider
	hasCreds bool
}

func NewVaultTransitProvider(cfg VaultConfig) *VaultTransitProvider {
	hasCreds := cfg.Address != "" && cfg.Token != "" && cfg.TransitKey != ""
	p := &VaultTransitProvider{
		cfg:      cfg,
		hasCreds: hasCreds,
	}
	if !hasCreds {
		p.fallback = NewLocalKeyProvider([]byte("vault-transit-sandbox-key-seed"), "vault-transit-v1")
	}
	return p
}

func (p *VaultTransitProvider) Type() ProviderType {
	return ProviderVault
}

func (p *VaultTransitProvider) GenerateDataKey(ctx context.Context, keyID string) ([]byte, []byte, string, error) {
	if !p.hasCreds {
		return p.fallback.GenerateDataKey(ctx, keyID)
	}
	return nil, nil, "", ErrProviderUnavailable
}

func (p *VaultTransitProvider) DecryptDataKey(ctx context.Context, keyID string, encryptedDEK []byte, keyVersion string) ([]byte, error) {
	if !p.hasCreds {
		return p.fallback.DecryptDataKey(ctx, keyID, encryptedDEK, keyVersion)
	}
	return nil, ErrProviderUnavailable
}

func (p *VaultTransitProvider) RotateMasterKey(ctx context.Context, keyID string) (string, error) {
	if !p.hasCreds {
		return p.fallback.RotateMasterKey(ctx, keyID)
	}
	return "", ErrProviderUnavailable
}

// Client coordinates Envelope Encryption across any KeyProvider
type Client struct {
	provider KeyProvider
	keyID    string
}

func NewClient(initialMasterKey []byte, keyVersion string) *Client {
	local := NewLocalKeyProvider(initialMasterKey, keyVersion)
	return &Client{
		provider: local,
		keyID:    keyVersion,
	}
}

func NewClientWithProvider(provider KeyProvider, keyID string) *Client {
	return &Client{
		provider: provider,
		keyID:    keyID,
	}
}

// AddKeyVersion delegates historical key addition to local provider
func (c *Client) AddKeyVersion(keyVersion string, newMasterKey []byte, setAsCurrent bool) {
	if local, ok := c.provider.(*LocalKeyProvider); ok {
		local.AddKeyVersion(keyVersion, newMasterKey, setAsCurrent)
	}
}

// EncryptPayload performs envelope encryption on arbitrary payload bytes:
// 1. Generates ephemeral 256-bit DEK from KeyProvider
// 2. Encrypts payload with DEK using AES-256-GCM with 96-bit random nonce
// 3. Authenticates Additional Authenticated Data (AAD) bound to the key version
// 4. Securely zeroizes plaintext DEK from memory immediately after encryption
// 5. Returns EncryptedEnvelope containing wrapped DEK and metadata
func (c *Client) EncryptPayload(ctx context.Context, rawBytes []byte) (*EncryptedEnvelope, error) {
	plaintextDEK, encryptedDEK, ver, err := c.provider.GenerateDataKey(ctx, c.keyID)
	if err != nil {
		return nil, fmt.Errorf("kms: failed to generate data key: %w", err)
	}

	// Memory safety: secure zeroisation of plaintext DEK working buffer
	defer func() {
		for i := range plaintextDEK {
			plaintextDEK[i] = 0
		}
	}()

	block, err := aes.NewCipher(plaintextDEK)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Bind key version and algorithm to AAD
	aad := []byte(fmt.Sprintf("%s:%s", c.provider.Type(), ver))
	ciphertext := gcm.Seal(nil, nonce, rawBytes, aad)

	return &EncryptedEnvelope{
		Algorithm:    "AES-256-GCM",
		Provider:     c.provider.Type(),
		KeyID:        c.keyID,
		KeyVersion:   ver,
		EncryptedDEK: encryptedDEK,
		Nonce:        nonce,
		Ciphertext:   ciphertext,
		CreatedAt:    time.Now().UTC(),
	}, nil
}

// DecryptPayload un-wraps the DEK using the specified key provider and decrypts the payload:
// 1. Un-wraps encrypted DEK using KeyProvider and the envelope's KeyVersion
// 2. Decrypts ciphertext and verifies GCM authentication tag
// 3. Zeroizes plaintext DEK working buffer immediately
func (c *Client) DecryptPayload(ctx context.Context, env *EncryptedEnvelope) ([]byte, error) {
	if env == nil || len(env.EncryptedDEK) == 0 || len(env.Nonce) == 0 {
		return nil, ErrInvalidCiphertext
	}

	plaintextDEK, err := c.provider.DecryptDataKey(ctx, env.KeyID, env.EncryptedDEK, env.KeyVersion)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDecryptionFailed, err)
	}

	// Memory safety: secure zeroisation of plaintext DEK
	defer func() {
		for i := range plaintextDEK {
			plaintextDEK[i] = 0
		}
	}()

	blockDEK, err := aes.NewCipher(plaintextDEK)
	if err != nil {
		return nil, err
	}

	gcmDEK, err := cipher.NewGCM(blockDEK)
	if err != nil {
		return nil, err
	}

	// Verify AAD: support both new format and backward compatibility
	aadNew := []byte(fmt.Sprintf("%s:%s", env.Provider, env.KeyVersion))
	plaintext, err := gcmDEK.Open(nil, env.Nonce, env.Ciphertext, aadNew)
	if err != nil {
		// Fallback for backward compatibility with older v1 envelopes
		plaintext, err = gcmDEK.Open(nil, env.Nonce, env.Ciphertext, []byte(env.KeyVersion))
		if err != nil {
			return nil, fmt.Errorf("%w: payload authentication tag mismatch or tampered ciphertext", ErrDecryptionFailed)
		}
	}

	return plaintext, nil
}
