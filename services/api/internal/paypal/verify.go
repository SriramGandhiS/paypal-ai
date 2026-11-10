package paypal

import (
	"context"
	"errors"
)

var (
	ErrSignatureInvalid    = errors.New("webhook signature verification failed")
	ErrVerifierUnavailable = errors.New("upstream PayPal signature verification service unavailable")
)

// WebhookHeaders contains the HTTP headers required for PayPal signature verification.
type WebhookHeaders struct {
	TransmissionID   string
	TransmissionTime string
	TransmissionSig  string
	CertURL          string
	AuthAlgo         string
	WebhookID        string
}

// WebhookVerifier validates incoming webhook signatures against PayPal's verification API.
type WebhookVerifier interface {
	VerifySignature(ctx context.Context, headers WebhookHeaders, rawBody []byte) (bool, error)
}

// LocalMockVerifier provides deterministic verification for unit and simulated tests.
type LocalMockVerifier struct {
	ShouldPass bool
	ShouldFail bool
}

// VerifySignature implements WebhookVerifier for testing.
func (m *LocalMockVerifier) VerifySignature(ctx context.Context, headers WebhookHeaders, rawBody []byte) (bool, error) {
	if m.ShouldFail {
		return false, ErrVerifierUnavailable
	}
	return m.ShouldPass, nil
}
