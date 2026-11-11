package compliance_test

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/paybridge/paybridge-ai/services/api/internal/compliance"
)

func TestCompliance_RBIPurposeCode(t *testing.T) {
	eng := compliance.NewEngine()

	// Eligible FEMA service export code P0802
	err := eng.ValidateCorridorAndPurpose("US", "IN", "P0802")
	if err != nil {
		t.Fatalf("expected P0802 to be accepted, got error: %v", err)
	}

	// Ineligible / missing code for India
	err = eng.ValidateCorridorAndPurpose("US", "IN", "INVALID_CODE")
	if err == nil {
		t.Fatal("expected error on invalid RBI purpose code for India corridor")
	}

	// Non-India corridor should pass without RBI purpose requirement
	err = eng.ValidateCorridorAndPurpose("US", "GB", "")
	if err != nil {
		t.Errorf("unexpected error for US->GB corridor: %v", err)
	}
}

func TestCompliance_SanctionsHitBlocksImmediately(t *testing.T) {
	eng := compliance.NewEngine()

	// Clean entity passes
	err := eng.CheckSanctionsScreening("Acme Global Software", "Acme Global")
	if err != nil {
		t.Fatalf("unexpected sanctions hit on clean entity: %v", err)
	}

	// OFAC hit blocks immediately
	err = eng.CheckSanctionsScreening("John Doe", "BLOCKED ENTITY LTD")
	if err == nil {
		t.Fatal("CRITICAL INVARIANT VIOLATION: Sanctioned entity was not blocked")
	}
}

func TestCompliance_TravelRuleEnforcement(t *testing.T) {
	eng := compliance.NewEngine()

	// $2,500 (< $3,000 threshold)
	err := eng.EvaluateTravelRule(decimal.NewFromInt(2500), "USD", false, false)
	if err != nil {
		t.Fatalf("unexpected travel rule error below threshold: %v", err)
	}

	// $5,000 without verified identifiers
	err = eng.EvaluateTravelRule(decimal.NewFromInt(5000), "USD", false, true)
	if err == nil {
		t.Fatal("expected Travel Rule violation for transaction >= $3,000 without full ID")
	}

	// $5,000 with complete identifiers
	err = eng.EvaluateTravelRule(decimal.NewFromInt(5000), "USD", true, true)
	if err != nil {
		t.Fatalf("unexpected Travel Rule error with complete IDs: %v", err)
	}
}

func TestCompliance_EFIRCIngestion(t *testing.T) {
	eng := compliance.NewEngine()
	mockPDF := []byte("%PDF-1.4 Mock AD-1 Bank e-FIRC Inward Remittance")

	doc, err := eng.IngestEFIRC(mockPDF, "State Bank of India", "0021004", "SBI9912004", "P0802", "104250.00")
	if err != nil {
		t.Fatalf("e-FIRC ingestion failed: %v", err)
	}

	if doc.PDFSHA256 == "" {
		t.Error("expected non-empty SHA-256 hash for e-FIRC PDF")
	}

	// Re-ingestion of identical document must be rejected as duplicate
	_, err = eng.IngestEFIRC(mockPDF, "State Bank of India", "0021004", "SBI9912004", "P0802", "104250.00")
	if err == nil {
		t.Fatal("expected duplicate e-FIRC document rejection, got nil error")
	}
}

func TestCompliance_EFIRCFileValidation(t *testing.T) {
	eng := compliance.NewEngine()

	// Non-PDF data (corrupted / missing %PDF- magic bytes)
	nonPDF := []byte("INVALID_DATA_NOT_A_PDF_FILE")
	_, err := eng.IngestEFIRC(nonPDF, "HDFC", "001", "UTR001", "P0802", "5000")
	if err == nil {
		t.Fatal("expected invalid PDF file rejection, got nil error")
	}

	// Empty bytes
	_, err = eng.IngestEFIRC([]byte{}, "HDFC", "001", "UTR002", "P0802", "5000")
	if err == nil {
		t.Fatal("expected empty bytes rejection, got nil error")
	}
}

func TestCompliance_PersonaKYCIntegration(t *testing.T) {
	client := compliance.NewPersonaClient("persona_test_key_123")
	ctx := context.Background()

	inquiryID, err := client.CreateInquiry(ctx, "user-ref-99", "Dinny Paul Navis")
	if err != nil {
		t.Fatalf("failed to create inquiry: %v", err)
	}

	status, err := client.GetInquiryStatus(ctx, inquiryID)
	if err != nil {
		t.Fatalf("failed to get inquiry status: %v", err)
	}
	if status != "approved" {
		t.Errorf("expected approved status in sandbox, got %s", status)
	}
}

func TestCompliance_PreAuthorizationGate(t *testing.T) {
	eng := compliance.NewEngine()
	ctx := context.Background()

	// 1. Clean transaction passes pre-auth gate
	err := eng.PreAuthorizationGate(
		ctx,
		"US", "IN", "P0802",
		decimal.NewFromInt(1500), "USD",
		"Clean Merchant LLC", "Indian Dev Studio",
		true, true,
	)
	if err != nil {
		t.Fatalf("expected clean pre-auth gate to pass, got: %v", err)
	}

	// 2. Sanctioned actor hits compliance hold
	err = eng.PreAuthorizationGate(
		ctx,
		"US", "IN", "P0802",
		decimal.NewFromInt(1500), "USD",
		"BLOCKED ENTITY LTD", "Indian Dev Studio",
		true, true,
	)
	if err == nil {
		t.Fatal("CRITICAL INVARIANT VIOLATION: Sanctioned actor bypassed PreAuthorizationGate")
	}
}
