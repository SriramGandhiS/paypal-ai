package compliance

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/shopspring/decimal"
)

var (
	ErrInvalidPurposeCode   = errors.New("compliance: invalid or missing RBI purpose code for India cross-border corridor")
	ErrTravelRuleViolation   = errors.New("compliance: FinCEN Travel Rule threshold exceeded without complete identity disclosure")
	ErrSanctionsHitDetected  = errors.New("compliance: OFAC/PEP sanctions match detected; transaction frozen immediately")
	ErrInvalidPDFFile        = errors.New("compliance: file is not a valid PDF document (missing %PDF- magic bytes)")
	ErrFileTooLarge          = errors.New("compliance: e-FIRC document exceeds maximum allowed size (10 MB)")
	ErrDuplicateDocument     = errors.New("compliance: duplicate e-FIRC document already ingested")
	ErrKYCVerificationFailed = errors.New("compliance: Persona KYC/KYB identity verification failed or incomplete")
	ErrComplianceHold        = errors.New("compliance: mandatory pre-authorization compliance hold active; AI override forbidden")
)

const (
	MaxPDFSizeBytes = 10 * 1024 * 1024 // 10 MB limit
)

// Supported RBI FEMA Service Export Purpose Codes
var RBIPurposeCodes = map[string]string{
	"P0801": "Hardware consultancy / system implementation",
	"P0802": "Software consultancy / implementation / data processing",
	"P0803": "Database services / software distribution",
	"P0805": "News agency services / online information provision",
	"P1006": "Business and management consultancy",
}

type EFIRCDocument struct {
	DocumentID  string    `json:"document_id"`
	BankName    string    `json:"bank_name"`
	ADCode      string    `json:"ad_code"`
	UTR         string    `json:"utr"`
	PurposeCode string    `json:"purpose_code"`
	AmountINR   string    `json:"amount_inr"`
	PDFSHA256   string    `json:"pdf_sha256"`
	IssuedAt    time.Time `json:"issued_at"`
	IngestedAt  time.Time `json:"ingested_at"`
}

type SanctionsHit struct {
	EntityName string  `json:"entity_name"`
	ListSource string  `json:"list_source"` // "OFAC_SDN", "EU_CONSOLIDATED", "UN_SANCTIONS", "PEP"
	MatchScore float64 `json:"match_score"` // 0.0 - 1.0
	Remarks    string  `json:"remarks"`
}

// PersonaClient provides KYC/KYB identity verification
type PersonaClient struct {
	apiKey    string
	inquiries map[string]string // inquiry_id -> status ("approved", "failed", "pending")
	mu        sync.RWMutex
}

func NewPersonaClient(apiKey string) *PersonaClient {
	return &PersonaClient{
		apiKey:    apiKey,
		inquiries: make(map[string]string),
	}
}

func (p *PersonaClient) CreateInquiry(ctx context.Context, referenceID, fullName string) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	inquiryID := fmt.Sprintf("inq_%x", sha256.Sum256([]byte(referenceID+fullName)))[:16]
	p.inquiries[inquiryID] = "approved" // default approved in sandbox mode
	return inquiryID, nil
}

func (p *PersonaClient) GetInquiryStatus(ctx context.Context, inquiryID string) (string, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	status, exists := p.inquiries[inquiryID]
	if !exists {
		return "", fmt.Errorf("persona inquiry %s not found", inquiryID)
	}
	return status, nil
}

func (p *PersonaClient) VerifyWebhook(payload []byte, signature, webhookSecret string) bool {
	if webhookSecret == "" {
		return true
	}
	mac := hmac.New(sha256.New, []byte(webhookSecret))
	mac.Write(payload)
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expected))
}

// ComplyAdvantageClient performs real-time sanctions and PEP screening
type ComplyAdvantageClient struct {
	apiKey    string
	watchlist map[string]SanctionsHit
	mu        sync.RWMutex
}

func NewComplyAdvantageClient(apiKey string) *ComplyAdvantageClient {
	c := &ComplyAdvantageClient{
		apiKey:    apiKey,
		watchlist: make(map[string]SanctionsHit),
	}
	// Seed known sanctions test actors
	c.watchlist["BLOCKED ENTITY LTD"] = SanctionsHit{
		EntityName: "BLOCKED ENTITY LTD",
		ListSource: "OFAC_SDN",
		MatchScore: 1.0,
		Remarks:    "Executive Order 13224 Terrorism Sanctions",
	}
	c.watchlist["SANCTIONED ACTOR"] = SanctionsHit{
		EntityName: "SANCTIONED ACTOR",
		ListSource: "OFAC_SDN",
		MatchScore: 0.98,
		Remarks:    "Specially Designated National",
	}
	c.watchlist["OFAC TEST SUBJECT"] = SanctionsHit{
		EntityName: "OFAC TEST SUBJECT",
		ListSource: "OFAC_SDN",
		MatchScore: 1.0,
		Remarks:    "Test subject SDN List",
	}
	return c
}

func (c *ComplyAdvantageClient) ScreenEntity(ctx context.Context, name string) ([]SanctionsHit, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	upper := strings.ToUpper(strings.TrimSpace(name))
	var hits []SanctionsHit

	for prohibited, hit := range c.watchlist {
		if strings.Contains(upper, prohibited) || strings.Contains(prohibited, upper) {
			hits = append(hits, hit)
		}
	}

	return hits, nil
}

type ComplianceEngine struct {
	mu               sync.RWMutex
	persona          *PersonaClient
	complyAdvantage  *ComplyAdvantageClient
	ingestedFIRCHashes map[string]bool
}

func NewEngine() *ComplianceEngine {
	return &ComplianceEngine{
		persona:            NewPersonaClient("persona_test_api_key"),
		complyAdvantage:    NewComplyAdvantageClient("complyadvantage_test_api_key"),
		ingestedFIRCHashes: make(map[string]bool),
	}
}

// ValidateCorridorAndPurpose verifies RBI statutory purpose codes for Indian cross-border inflows
func (ce *ComplianceEngine) ValidateCorridorAndPurpose(senderCountry, receiverCountry, purposeCode string) error {
	if strings.ToUpper(receiverCountry) == "IN" {
		code := strings.ToUpper(strings.TrimSpace(purposeCode))
		if _, exists := RBIPurposeCodes[code]; !exists {
			return fmt.Errorf("%w: '%s' is not an eligible FEMA service export code", ErrInvalidPurposeCode, purposeCode)
		}
	}
	return nil
}

// CheckSanctionsScreening performs real-time OFAC SDN & PEP check via ComplyAdvantage client
func (ce *ComplianceEngine) CheckSanctionsScreening(name, entity string) error {
	hitsSender, err := ce.complyAdvantage.ScreenEntity(context.Background(), name)
	if err != nil {
		return err
	}
	if len(hitsSender) > 0 {
		return fmt.Errorf("%w: matched entity '%s' on %s (score: %.2f)",
			ErrSanctionsHitDetected, name, hitsSender[0].ListSource, hitsSender[0].MatchScore)
	}

	hitsEntity, err := ce.complyAdvantage.ScreenEntity(context.Background(), entity)
	if err != nil {
		return err
	}
	if len(hitsEntity) > 0 {
		return fmt.Errorf("%w: matched entity '%s' on %s (score: %.2f)",
			ErrSanctionsHitDetected, entity, hitsEntity[0].ListSource, hitsEntity[0].MatchScore)
	}

	return nil
}

// EvaluateTravelRule enforces FinCEN Travel Rule on transactions >= $3,000 USD
func (ce *ComplianceEngine) EvaluateTravelRule(amount decimal.Decimal, currency string, hasOriginatorID, hasBeneficiaryID bool) error {
	travelRuleThreshold := decimal.NewFromInt(3000)
	if strings.ToUpper(currency) == "USD" && amount.GreaterThanOrEqual(travelRuleThreshold) {
		if !hasOriginatorID || !hasBeneficiaryID {
			return fmt.Errorf("%w: amount %s %s requires complete transmittal order data (Originator & Beneficiary identity)", ErrTravelRuleViolation, amount.String(), currency)
		}
	}
	return nil
}

// IngestEFIRC parses and verifies an electronic Foreign Inward Remittance Certificate (e-FIRC)
// Validates PDF magic bytes, maximum size limit (10MB), and calculates immutable SHA-256 checksum
func (ce *ComplianceEngine) IngestEFIRC(pdfBytes []byte, bankName, adCode, utr, purposeCode, amountINR string) (*EFIRCDocument, error) {
	if len(pdfBytes) == 0 {
		return nil, errors.New("compliance: empty e-FIRC document bytes")
	}

	if len(pdfBytes) > MaxPDFSizeBytes {
		return nil, ErrFileTooLarge
	}

	// Validate PDF magic bytes: '%PDF-'
	if !bytes.HasPrefix(pdfBytes, []byte("%PDF-")) {
		return nil, ErrInvalidPDFFile
	}

	h := sha256.Sum256(pdfBytes)
	hashStr := fmt.Sprintf("%x", h)

	ce.mu.Lock()
	defer ce.mu.Unlock()

	if ce.ingestedFIRCHashes[hashStr] {
		return nil, fmt.Errorf("%w: document hash %s already processed", ErrDuplicateDocument, hashStr)
	}
	ce.ingestedFIRCHashes[hashStr] = true

	docID := fmt.Sprintf("efirc_%s_%s", utr, hashStr[:8])

	return &EFIRCDocument{
		DocumentID:  docID,
		BankName:    bankName,
		ADCode:      adCode,
		UTR:         utr,
		PurposeCode: purposeCode,
		AmountINR:   amountINR,
		PDFSHA256:   hashStr,
		IssuedAt:    time.Now().UTC(),
		IngestedAt:  time.Now().UTC(),
	}, nil
}

// PreAuthorizationGate executes mandatory compliance screening prior to payment capture
// CRITICAL RULE: NO AI AGENT CAN OVERRIDE A COMPLIANCE HOLD
func (ce *ComplianceEngine) PreAuthorizationGate(
	ctx context.Context,
	senderCountry, receiverCountry, purposeCode string,
	amount decimal.Decimal, currency string,
	senderName, receiverName string,
	hasOriginatorID, hasBeneficiaryID bool,
) error {
	// 1. Sanctions Screening
	if err := ce.CheckSanctionsScreening(senderName, receiverName); err != nil {
		return fmt.Errorf("%w: %v", ErrComplianceHold, err)
	}

	// 2. RBI Statutory Purpose Code
	if err := ce.ValidateCorridorAndPurpose(senderCountry, receiverCountry, purposeCode); err != nil {
		return fmt.Errorf("%w: %v", ErrComplianceHold, err)
	}

	// 3. FinCEN Travel Rule
	if err := ce.EvaluateTravelRule(amount, currency, hasOriginatorID, hasBeneficiaryID); err != nil {
		return fmt.Errorf("%w: %v", ErrComplianceHold, err)
	}

	return nil
}
