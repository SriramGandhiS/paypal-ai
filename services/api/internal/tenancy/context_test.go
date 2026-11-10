package tenancy_test

import (
	"context"
	"testing"

	"github.com/paybridge/paybridge-ai/services/api/internal/tenancy"
)

func TestTenancyContext(t *testing.T) {
	ctx := context.Background()

	// Missing tenant
	_, err := tenancy.TenantFromContext(ctx)
	if err != tenancy.ErrNoTenantInContext {
		t.Errorf("expected ErrNoTenantInContext, got: %v", err)
	}

	// Missing user
	_, _, err = tenancy.UserFromContext(ctx)
	if err != tenancy.ErrNoUserInContext {
		t.Errorf("expected ErrNoUserInContext, got: %v", err)
	}

	// Attached tenant
	ctx = tenancy.WithTenant(ctx, "tenant-123")
	tenantID, err := tenancy.TenantFromContext(ctx)
	if err != nil || tenantID != "tenant-123" {
		t.Errorf("expected tenant-123, got %s (err: %v)", tenantID, err)
	}

	// Attached user & role
	ctx = tenancy.WithUser(ctx, "user-456", "MERCHANT_ADMIN")
	userID, role, err := tenancy.UserFromContext(ctx)
	if err != nil || userID != "user-456" || role != "MERCHANT_ADMIN" {
		t.Errorf("expected user-456 / MERCHANT_ADMIN, got %s / %s (err: %v)", userID, role, err)
	}
}
