package tenancy

import (
	"context"
	"errors"
)

type contextKey string

const (
	tenantIDKey contextKey = "paybridge.tenant_id"
	userIDKey   contextKey = "paybridge.user_id"
	userRoleKey contextKey = "paybridge.user_role"
)

var (
	ErrNoTenantInContext = errors.New("no tenant found in request context")
	ErrNoUserInContext   = errors.New("no user found in request context")
)

// WithTenant attaches the active tenant ID to the context.
func WithTenant(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, tenantIDKey, tenantID)
}

// TenantFromContext retrieves the active tenant ID from the context.
func TenantFromContext(ctx context.Context) (string, error) {
	val := ctx.Value(tenantIDKey)
	if val == nil {
		return "", ErrNoTenantInContext
	}
	s, ok := val.(string)
	if !ok || s == "" {
		return "", ErrNoTenantInContext
	}
	return s, nil
}

// WithUser attaches the authenticated user ID and role to the context.
func WithUser(ctx context.Context, userID, role string) context.Context {
	ctx = context.WithValue(ctx, userIDKey, userID)
	return context.WithValue(ctx, userRoleKey, role)
}

// UserFromContext retrieves the authenticated user ID and role from the context.
func UserFromContext(ctx context.Context) (string, string, error) {
	uVal := ctx.Value(userIDKey)
	rVal := ctx.Value(userRoleKey)
	if uVal == nil || rVal == nil {
		return "", "", ErrNoUserInContext
	}
	u, uOk := uVal.(string)
	r, rOk := rVal.(string)
	if !uOk || !rOk || u == "" || r == "" {
		return "", "", ErrNoUserInContext
	}
	return u, r, nil
}
