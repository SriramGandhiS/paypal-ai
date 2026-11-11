package multirail

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/shopspring/decimal"
)

type StripeConfig struct {
	SecretKey      string
	WebhookSecret  string
	BaseURL        string
	ToleranceSec   int64
}

type StripeAdapter struct {
	mu                sync.RWMutex
	cfg               StripeConfig
	transactions      map[string]*RailTransaction
	disputes          map[string]*DisputeRecord
	processedEvents   map[string]time.Time // provider_event_id -> processed timestamp (replay defense)
	idempotencyKeys   map[string]*RailTransaction
}

var (
	_ RailReader      = (*StripeAdapter)(nil)
	_ RailMutator     = (*StripeAdapter)(nil)
	_ WebhookVerifier = (*StripeAdapter)(nil)
)

func NewStripeAdapter(cfg StripeConfig) *StripeAdapter {
	if cfg.ToleranceSec <= 0 {
		cfg.ToleranceSec = 300 // 5 minutes default tolerance
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.stripe.com/v1"
	}
	return &StripeAdapter{
		cfg:             cfg,
		transactions:    make(map[string]*RailTransaction),
		disputes:        make(map[string]*DisputeRecord),
		processedEvents: make(map[string]time.Time),
		idempotencyKeys: make(map[string]*RailTransaction),
	}
}

// VerifyWebhook validates Stripe-Signature header format: t=1614555555,v1=5257...
// Computes HMAC-SHA256(webhook_secret, timestamp + "." + rawBody) and checks replay and clock drift
func (s *StripeAdapter) VerifyWebhook(rawBody []byte, signatureHeader string) (string, error) {
	if s.cfg.WebhookSecret == "" {
		return "", errors.New("stripe: webhook secret not configured")
	}

	parts := strings.Split(signatureHeader, ",")
	var timestampStr string
	var signatures []string

	for _, part := range parts {
		kv := strings.SplitN(strings.TrimSpace(part), "=", 2)
		if len(kv) == 2 {
			switch kv[0] {
			case "t":
				timestampStr = kv[1]
			case "v1":
				signatures = append(signatures, kv[1])
			}
		}
	}

	if timestampStr == "" || len(signatures) == 0 {
		return "", fmt.Errorf("%w: missing timestamp or v1 signature in Stripe-Signature header", ErrInvalidSignature)
	}

	ts, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return "", fmt.Errorf("%w: invalid timestamp format", ErrInvalidSignature)
	}

	// Verify timestamp tolerance window (clock drift / replay prevention)
	now := time.Now().Unix()
	if now-ts > s.cfg.ToleranceSec || ts-now > s.cfg.ToleranceSec {
		return "", fmt.Errorf("%w: timestamp %d outside %ds tolerance window", ErrWebhookTimestampExpired, ts, s.cfg.ToleranceSec)
	}

	// Compute expected signature: HMAC-SHA256(secret, t.payload)
	mac := hmac.New(sha256.New, []byte(s.cfg.WebhookSecret))
	mac.Write([]byte(timestampStr))
	mac.Write([]byte("."))
	mac.Write(rawBody)
	expectedSig := hex.EncodeToString(mac.Sum(nil))

	// Constant-time signature comparison
	matched := false
	for _, sig := range signatures {
		if hmac.Equal([]byte(sig), []byte(expectedSig)) {
			matched = true
			break
		}
	}

	if !matched {
		return "", fmt.Errorf("%w: signature mismatch", ErrInvalidSignature)
	}

	// Calculate deterministic event ID from body SHA-256 for replay defense
	h := sha256.Sum256(rawBody)
	eventHash := fmt.Sprintf("evt_sha_%x", h[:16])

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.processedEvents[eventHash]; exists {
		return eventHash, fmt.Errorf("%w: event %s already processed", ErrWebhookReplay, eventHash)
	}
	s.processedEvents[eventHash] = time.Now().UTC()

	return eventHash, nil
}

func (s *StripeAdapter) InitiatePayment(ctx context.Context, auth HumanAuthorization, amount decimal.Decimal, currency, reference string) (*RailTransaction, error) {
	if auth.OperatorID == "" {
		return nil, fmt.Errorf("%w: payment initiation requires valid operator identity", ErrUnauthorizedMutation)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Idempotency check on reference
	if existing, found := s.idempotencyKeys[reference]; found {
		if !existing.Amount.Equal(amount) || existing.Currency != currency {
			return nil, ErrIdempotencyViolation
		}
		return existing, nil
	}

	txID := fmt.Sprintf("pi_%d", time.Now().UnixNano()%1000000)
	now := time.Now().UTC()
	h := sha256.Sum256([]byte(txID + amount.String() + currency))

	tx := &RailTransaction{
		Rail:           RailStripe,
		TransactionID:  txID,
		Amount:         amount,
		Currency:       currency,
		Status:         "requires_capture",
		OccurredAt:     now,
		RawPayloadHash: fmt.Sprintf("%x", h),
		IdempotencyKey: reference,
	}

	s.transactions[txID] = tx
	s.idempotencyKeys[reference] = tx
	return tx, nil
}

func (s *StripeAdapter) CapturePayment(ctx context.Context, auth HumanAuthorization, txID string, amount decimal.Decimal) (*RailTransaction, error) {
	if auth.OperatorID == "" {
		return nil, fmt.Errorf("%w: payment capture requires human operator confirmation", ErrUnauthorizedMutation)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	tx, exists := s.transactions[txID]
	if !exists {
		return nil, fmt.Errorf("stripe payment intent %s not found", txID)
	}

	if tx.Status == "succeeded" {
		return tx, nil // idempotent capture
	}

	tx.Status = "succeeded"
	return tx, nil
}

func (s *StripeAdapter) IssueRefund(ctx context.Context, auth HumanAuthorization, txID string, amount decimal.Decimal, reason string) (*RailTransaction, error) {
	if auth.OperatorID == "" {
		return nil, fmt.Errorf("%w: refund requires operator confirmation", ErrUnauthorizedMutation)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	tx, exists := s.transactions[txID]
	if !exists {
		return nil, fmt.Errorf("stripe payment intent %s not found", txID)
	}

	tx.Status = "refunded"
	return tx, nil
}

func (s *StripeAdapter) GetTransaction(ctx context.Context, txID string) (*RailTransaction, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tx, exists := s.transactions[txID]
	if !exists {
		return nil, fmt.Errorf("stripe payment intent %s not found", txID)
	}
	return tx, nil
}

func (s *StripeAdapter) GetDispute(ctx context.Context, disputeID string) (*DisputeRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	d, exists := s.disputes[disputeID]
	if !exists {
		return nil, fmt.Errorf("stripe dispute %s not found", disputeID)
	}
	return d, nil
}

func (s *StripeAdapter) IngestDisputeWebhook(disputeID, txID string, amount decimal.Decimal, currency, reason string) *DisputeRecord {
	s.mu.Lock()
	defer s.mu.Unlock()

	d := &DisputeRecord{
		DisputeID:     disputeID,
		TransactionID: txID,
		Rail:          RailStripe,
		Amount:        amount,
		Currency:      currency,
		Reason:        reason,
		Status:        DisputeStatusNeedsEvidence,
		EvidenceDueBy: time.Now().UTC().Add(7 * 24 * time.Hour),
	}
	s.disputes[disputeID] = d
	return d
}
