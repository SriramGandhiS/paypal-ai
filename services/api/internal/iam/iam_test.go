package iam_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/paybridge/paybridge-ai/services/api/internal/iam"
)

func TestSCIM_AuthenticationAndUserLifecycle(t *testing.T) {
	service := iam.NewService("scim-test-secret-token")
	ctx := context.Background()

	// 1. Auth check
	err := service.AuthenticateSCIM("Bearer scim-test-secret-token")
	if err != nil {
		t.Fatalf("expected valid bearer token to pass: %v", err)
	}

	err = service.AuthenticateSCIM("Bearer invalid-token")
	if err == nil {
		t.Fatal("expected unauthorized on invalid bearer token")
	}

	// 2. Create user via SCIM POST /scim/v2/Users
	user, err := service.CreateUser(ctx, "tenant-alpha", "jdoe@example.com", "John Doe", iam.RoleMerchantAdmin, "jdoe@example.com")
	if err != nil {
		t.Fatalf("failed to create SCIM user: %v", err)
	}
	if !user.Active {
		t.Errorf("expected newly created user to be active")
	}

	// 3. Retrieve user
	fetched, err := service.GetUser(ctx, user.ID)
	if err != nil {
		t.Fatalf("failed to retrieve user: %v", err)
	}
	if fetched.DisplayName != "John Doe" {
		t.Errorf("expected DisplayName John Doe, got %s", fetched.DisplayName)
	}

	// 4. Deactivate user (employee offboarding)
	err = service.DeactivateUser(ctx, user.ID)
	if err != nil {
		t.Fatalf("failed to deactivate user: %v", err)
	}

	updated, _ := service.GetUser(ctx, user.ID)
	if updated.Active {
		t.Fatal("expected user to be deactivated")
	}
}

func TestSAML_AssertionValidation(t *testing.T) {
	service := iam.NewService("token")

	validSAML := `<saml2:Assertion xmlns:saml2="urn:oasis:names:tc:SAML:2.0:assertion"><saml2:AudienceRestriction><saml2:Audience>https://paybridge.internal/saml/sso</saml2:Audience></saml2:AudienceRestriction></saml2:Assertion>`

	session, err := service.ValidateSAMLAssertion(validSAML, "https://paybridge.internal/saml/sso")
	if err != nil {
		t.Fatalf("expected valid SAML assertion to pass: %v", err)
	}
	if session.Role != iam.RoleMerchantAdmin {
		t.Errorf("expected role merchant_admin, got %s", session.Role)
	}

	// Invalid / forged SAML
	_, err = service.ValidateSAMLAssertion("<corrupted></corrupted>", "https://paybridge.internal/saml/sso")
	if err == nil {
		t.Fatal("expected invalid SAML to be rejected, got nil error")
	}
}

func TestWebAuthn_FIDO2ChallengeFlow(t *testing.T) {
	service := iam.NewService("token")
	userID := "usr_fido2_test"

	challenge, err := service.GenerateWebAuthnChallenge(userID)
	if err != nil {
		t.Fatalf("failed to generate challenge: %v", err)
	}

	clientDataJSON := fmt.Sprintf(`{"type":"webauthn.get","challenge":"%s","origin":"https://paybridge.internal"}`, challenge)

	// Valid challenge response passes
	err = service.VerifyWebAuthnResponse(userID, challenge, clientDataJSON)
	if err != nil {
		t.Fatalf("expected valid WebAuthn challenge verification to pass: %v", err)
	}

	// Replay of consumed challenge is rejected
	err = service.VerifyWebAuthnResponse(userID, challenge, clientDataJSON)
	if err == nil {
		t.Fatal("expected replay of consumed WebAuthn challenge to fail")
	}
}

func TestIAM_CrossTenantIDORDefense(t *testing.T) {
	session := &iam.SessionContext{
		UserID:      "usr_alice",
		TenantID:    "tenant_acme",
		Role:        iam.RoleMerchantAdmin,
		MFAVerified: true,
	}

	// Same tenant allowed
	err := iam.AuthorizeTenantAction(session, "tenant_acme", iam.RoleMerchantAdmin)
	if err != nil {
		t.Fatalf("expected authorized access for same tenant: %v", err)
	}

	// Cross-tenant access blocked (IDOR protection)
	err = iam.AuthorizeTenantAction(session, "tenant_globex", iam.RoleMerchantAdmin)
	if err == nil {
		t.Fatal("CRITICAL INVARIANT VIOLATION: Cross-tenant access allowed (IDOR breach)")
	}
}

