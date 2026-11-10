package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/paybridge/paybridge-ai/services/api/internal/api"
	"github.com/paybridge/paybridge-ai/services/api/internal/auth"
	"github.com/paybridge/paybridge-ai/services/api/internal/cases"
	"github.com/paybridge/paybridge-ai/services/api/internal/checkout"
	"github.com/paybridge/paybridge-ai/services/api/internal/evidence"
	"github.com/paybridge/paybridge-ai/services/api/internal/httpx"
	"github.com/paybridge/paybridge-ai/services/api/internal/paypal"
	"github.com/paybridge/paybridge-ai/services/api/internal/reconciliation"
	"github.com/paybridge/paybridge-ai/services/api/internal/webhook"
)

type mockSessionReader struct {
	session *auth.Session
}

func (m *mockSessionReader) LookupSession(ctx context.Context, tokenHash []byte) (*auth.Session, error) {
	if m.session != nil {
		return m.session, nil
	}
	return nil, auth.ErrInvalidSessionToken
}

type mockPayPalService struct {
	createOrderFunc func(ctx context.Context, ord paypal.OrderRequest) (*paypal.OrderResponse, error)
	captureFunc     func(ctx context.Context, orderID string) (*paypal.OrderResponse, error)
	getOrderFunc    func(ctx context.Context, orderID string) (*paypal.OrderResponse, error)
	getCaptureFunc  func(ctx context.Context, captureID string) (*paypal.CaptureGroundTruth, error)
}

func (m *mockPayPalService) GetOrder(ctx context.Context, orderID string) (*paypal.OrderResponse, error) {
	if m.getOrderFunc != nil {
		return m.getOrderFunc(ctx, orderID)
	}
	return &paypal.OrderResponse{
		ID:     orderID,
		Status: "APPROVED",
	}, nil
}

func (m *mockPayPalService) CreateOrder(ctx context.Context, ord paypal.OrderRequest) (*paypal.OrderResponse, error) {
	if m.createOrderFunc != nil {
		return m.createOrderFunc(ctx, ord)
	}
	return &paypal.OrderResponse{
		ID:          "ORDER-123",
		Status:      "CREATED",
		ApprovalURL: "https://sandbox.paypal.com/checkoutnow?token=ORDER-123",
	}, nil
}

func (m *mockPayPalService) CaptureOrder(ctx context.Context, orderID string) (*paypal.OrderResponse, error) {
	if m.captureFunc != nil {
		return m.captureFunc(ctx, orderID)
	}
	return &paypal.OrderResponse{
		ID:            orderID,
		Status:        "COMPLETED",
		CaptureID:     "CAP-456",
		CaptureStatus: "COMPLETED",
		GrossAmount:   "25.00",
		Currency:      "USD",
	}, nil
}

func (m *mockPayPalService) GetCapture(ctx context.Context, captureID string) (*paypal.CaptureGroundTruth, error) {
	if m.getCaptureFunc != nil {
		return m.getCaptureFunc(ctx, captureID)
	}
	return &paypal.CaptureGroundTruth{
		CaptureID:   captureID,
		Status:      "COMPLETED",
		AmountCents: 2500,
		Currency:    "USD",
	}, nil
}

func setupTestRouter(t *testing.T) (*httpx.Router, *mockSessionReader, *mockPayPalService) {
	cfg := httpx.DefaultServerConfig(":8080", nil)
	r := httpx.NewRouter(cfg)

	sessReader := &mockSessionReader{
		session: &auth.Session{
			UserID:    "00000000-0000-0000-0000-000000000001",
			TenantID:  "11111111-1111-1111-1111-111111111111",
			Role:      "sender",
			ExpiresAt: time.Now().Add(24 * time.Hour),
		},
	}
	authMw := auth.NewMiddleware(auth.MiddlewareConfig{
		Reader:         sessReader,
		AllowDevBypass: false,
	})

	// In the absence of a live DB, we pass nil store which triggers safe error/mock branches
	casesSvc := cases.NewService(nil)
	checkoutSvc := checkout.NewService(nil)
	evidenceSvc := evidence.NewService(nil)
	reconciler := reconciliation.NewReconciler(nil, nil)
	mockPP := &mockPayPalService{}
	mockVerifier := &paypal.LocalMockVerifier{ShouldPass: true}
	whGateway := webhook.NewGateway(mockVerifier, nil, "WH-123", nil)

	api.RegisterRoutes(r, api.RouterDependencies{
		Cases:          casesSvc,
		Checkout:       checkoutSvc,
		Evidence:       evidenceSvc,
		Reconciler:     reconciler,
		PayPal:         mockPP,
		WebhookGateway: whGateway,
		AuthMiddleware: authMw,
	})

	return r, sessReader, mockPP
}

func TestCasesAPI_Unauthenticated(t *testing.T) {
	r, sessReader, _ := setupTestRouter(t)
	sessReader.session = nil // No session

	req := httptest.NewRequest(http.MethodPost, "/api/v1/cases", bytes.NewReader([]byte("{}")))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized, got %d: %s", w.Code, w.Body.String())
	}

	var prob httpx.ProblemDetails
	if err := json.Unmarshal(w.Body.Bytes(), &prob); err != nil {
		t.Fatalf("failed to decode RFC 7807 problem: %v", err)
	}
	if prob.Code != "UNAUTHORIZED" {
		t.Errorf("expected UNAUTHORIZED code, got %s", prob.Code)
	}
}

func TestCasesAPI_CreateCase_ValidationErrors(t *testing.T) {
	r, _, _ := setupTestRouter(t)

	testCases := []struct {
		name       string
		body       map[string]any
		expectCode int
	}{
		{
			name: "zero_amount",
			body: map[string]any{
				"amount_cents":     0,
				"currency":         "USD",
				"sender_country":   "US",
				"receiver_country": "DE",
			},
			expectCode: http.StatusBadRequest,
		},
		{
			name: "invalid_currency",
			body: map[string]any{
				"amount_cents":     1000,
				"currency":         "US", // not 3 chars
				"sender_country":   "US",
				"receiver_country": "DE",
			},
			expectCode: http.StatusBadRequest,
		},
		{
			name: "invalid_country",
			body: map[string]any{
				"amount_cents":     1000,
				"currency":         "USD",
				"sender_country":   "USA", // not 2 chars
				"receiver_country": "DE",
			},
			expectCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b, _ := json.Marshal(tc.body)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/cases", bytes.NewReader(b))
			req.Header.Set("Authorization", "Bearer valid-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tc.expectCode {
				t.Fatalf("expected status %d, got %d: %s", tc.expectCode, w.Code, w.Body.String())
			}
		})
	}
}

func TestEvidenceAPI_ValidationRules(t *testing.T) {
	r, _, _ := setupTestRouter(t)

	testCases := []struct {
		name       string
		body       map[string]any
		expectCode int
	}{
		{
			name: "invalid_tier_0",
			body: map[string]any{
				"kind":         "PAYMENT_CONFIRMATION",
				"tier":         0,
				"source_trust": "MERCHANT_REPORTED",
			},
			expectCode: http.StatusBadRequest,
		},
		{
			name: "simulated_must_be_tier_6_violation",
			body: map[string]any{
				"kind":         "PAYMENT_CONFIRMATION",
				"tier":         1, // Tier 1 with SIMULATED violates I11
				"source_trust": "SIMULATED",
			},
			expectCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b, _ := json.Marshal(tc.body)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/cases/case-123/evidence", bytes.NewReader(b))
			req.Header.Set("Authorization", "Bearer valid-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tc.expectCode {
				t.Fatalf("expected status %d, got %d: %s", tc.expectCode, w.Code, w.Body.String())
			}
		})
	}
}

func TestReconciliationEndpoint(t *testing.T) {
	r, _, _ := setupTestRouter(t)

	truth := reconciliation.CaptureGroundTruth{
		CaptureID:   "CAP-777",
		Status:      "COMPLETED",
		AmountCents: 1000,
		Currency:    "USD",
	}
	b, _ := json.Marshal(truth)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/cases/case-123/reconcile", bytes.NewReader(b))
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}

	var outcome reconciliation.ReconcileOutcome
	if err := json.Unmarshal(w.Body.Bytes(), &outcome); err != nil {
		t.Fatalf("failed to decode reconciliation outcome: %v", err)
	}

	if outcome.ProviderState != "CAPTURE_COMPLETED" {
		t.Errorf("expected ProviderState CAPTURE_COMPLETED, got %s", outcome.ProviderState)
	}
	if outcome.CaseStatus != "RESOLVED" {
		t.Errorf("expected CaseStatus RESOLVED, got %s", outcome.CaseStatus)
	}
}

func TestPayPalGoldenPath(t *testing.T) {
	r, _, pp := setupTestRouter(t)

	// Mock capture call
	pp.captureFunc = func(ctx context.Context, orderID string) (*paypal.OrderResponse, error) {
		return &paypal.OrderResponse{
			ID:            orderID,
			Status:        "COMPLETED",
			CaptureID:     "CAP-REAL-123",
			CaptureStatus: "COMPLETED",
			GrossAmount:   "35.00",
			Currency:      "USD",
		}, nil
	}
	pp.getCaptureFunc = func(ctx context.Context, captureID string) (*paypal.CaptureGroundTruth, error) {
		return &paypal.CaptureGroundTruth{
			CaptureID:   captureID,
			Status:      "COMPLETED",
			AmountCents: 3500,
			Currency:    "USD",
		}, nil
	}

	// Capture step
	captureReqBody := map[string]any{
		"order_id": "ORDER-555",
	}
	b, _ := json.Marshal(captureReqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/cases/case-999/checkout/paypal/capture", bytes.NewReader(b))
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Capture *paypal.OrderResponse            `json:"capture"`
		Outcome *reconciliation.ReconcileOutcome `json:"outcome"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode capture response: %v", err)
	}

	if resp.Capture.CaptureID != "CAP-REAL-123" {
		t.Errorf("expected CaptureID CAP-REAL-123, got %s", resp.Capture.CaptureID)
	}
	if resp.Outcome.ProviderState != "CAPTURE_COMPLETED" {
		t.Errorf("expected ProviderState CAPTURE_COMPLETED, got %s", resp.Outcome.ProviderState)
	}
}

func TestDevAuthBypassDisabled_RejectsXDevHeaders(t *testing.T) {
	r, _, _ := setupTestRouter(t)

	// Attempt bypass with header when AllowDevBypass == false
	req := httptest.NewRequest(http.MethodGet, "/api/v1/cases/case-123", nil)
	req.Header.Set("X-Dev-Tenant-ID", "attacker-tenant-id")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden when X-Dev-Tenant-ID is used without authorization, got %d: %s", w.Code, w.Body.String())
	}
}

func TestNoAIInitiationMutation_Forbidden(t *testing.T) {
	r, sessReader, _ := setupTestRouter(t)

	// Simulate AI role attempting checkout
	sessReader.session.Role = "ai_agent"

	reqBody := map[string]any{"amount_cents": 5000}
	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/cases/case-123/checkout/paypal", bytes.NewReader(b))
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden for AI role attempting checkout, got %d: %s", w.Code, w.Body.String())
	}

	// Simulate AI role attempting capture
	capBody := map[string]any{"order_id": "ORDER-123"}
	cb, _ := json.Marshal(capBody)
	capReq := httptest.NewRequest(http.MethodPost, "/api/v1/cases/case-123/checkout/paypal/capture", bytes.NewReader(cb))
	capReq.Header.Set("Authorization", "Bearer valid-token")
	capW := httptest.NewRecorder()
	r.ServeHTTP(capW, capReq)

	if capW.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden for AI role attempting capture, got %d: %s", capW.Code, capW.Body.String())
	}
}

func TestPayPalCapture_IdempotentDuplicateCapture(t *testing.T) {
	r, _, pp := setupTestRouter(t)

	captureCalled := 0
	pp.getOrderFunc = func(ctx context.Context, orderID string) (*paypal.OrderResponse, error) {
		// Provider already reports COMPLETED
		return &paypal.OrderResponse{
			ID:            orderID,
			Status:        "COMPLETED",
			CaptureID:     "CAP-IDEMPOTENT-888",
			CaptureStatus: "COMPLETED",
			GrossAmount:   "50.00",
			Currency:      "USD",
		}, nil
	}
	pp.captureFunc = func(ctx context.Context, orderID string) (*paypal.OrderResponse, error) {
		captureCalled++
		return nil, errors.New("should not be called when order already completed")
	}

	capBody := map[string]any{"order_id": "ORDER-ALREADY-DONE"}
	b, _ := json.Marshal(capBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/cases/case-123/checkout/paypal/capture", bytes.NewReader(b))
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on idempotent retry, got %d: %s", w.Code, w.Body.String())
	}
	if captureCalled != 0 {
		t.Errorf("expected CaptureOrder not to be called, called %d times", captureCalled)
	}
}

func TestPayPalCapture_TimeoutRecoveryViaGetOrder(t *testing.T) {
	r, _, pp := setupTestRouter(t)

	// First GetOrder returns APPROVED
	firstGet := true
	pp.getOrderFunc = func(ctx context.Context, orderID string) (*paypal.OrderResponse, error) {
		if firstGet {
			firstGet = false
			return &paypal.OrderResponse{ID: orderID, Status: "APPROVED"}, nil
		}
		// Second GetOrder during timeout recovery returns COMPLETED!
		return &paypal.OrderResponse{
			ID:            orderID,
			Status:        "COMPLETED",
			CaptureID:     "CAP-RECOVERED-999",
			CaptureStatus: "COMPLETED",
			GrossAmount:   "75.00",
			Currency:      "USD",
		}, nil
	}

	// CaptureOrder returns timeout/network error
	pp.captureFunc = func(ctx context.Context, orderID string) (*paypal.OrderResponse, error) {
		return nil, errors.New("upstream gateway timeout (504)")
	}

	capBody := map[string]any{"order_id": "ORDER-TIMEOUT-RECOVER"}
	b, _ := json.Marshal(capBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/cases/case-123/checkout/paypal/capture", bytes.NewReader(b))
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK via timeout recovery, got %d: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Capture *paypal.OrderResponse            `json:"capture"`
		Outcome *reconciliation.ReconcileOutcome `json:"outcome"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Capture.CaptureID != "CAP-RECOVERED-999" {
		t.Errorf("expected recovered capture ID CAP-RECOVERED-999, got %s", resp.Capture.CaptureID)
	}
}

func TestPayPalCapture_ConcurrentRequestsRaceCondition(t *testing.T) {
	r, _, pp := setupTestRouter(t)

	captureCallCount := 0
	pp.getOrderFunc = func(ctx context.Context, orderID string) (*paypal.OrderResponse, error) {
		return &paypal.OrderResponse{ID: orderID, Status: "APPROVED"}, nil
	}
	pp.captureFunc = func(ctx context.Context, orderID string) (*paypal.OrderResponse, error) {
		captureCallCount++
		// Simulate network latency to PayPal
		time.Sleep(30 * time.Millisecond)
		return &paypal.OrderResponse{
			ID:            orderID,
			Status:        "COMPLETED",
			CaptureID:     "CAP-RACE-SUCCESS",
			CaptureStatus: "COMPLETED",
			GrossAmount:   "50.00",
			Currency:      "USD",
		}, nil
	}

	capBody := map[string]any{"order_id": "ORDER-CONCURRENT-RACE"}
	b, _ := json.Marshal(capBody)

	startGate := make(chan struct{})
	statusCodeChan := make(chan int, 2)

	// Launch two concurrent capture attempts simultaneously
	for i := 0; i < 2; i++ {
		go func() {
			<-startGate
			req := httptest.NewRequest(http.MethodPost, "/api/v1/cases/case-123/checkout/paypal/capture", bytes.NewReader(b))
			req.Header.Set("Authorization", "Bearer valid-token")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			statusCodeChan <- w.Code
		}()
	}

	// Release both requests at the exact same instant
	close(startGate)

	code1 := <-statusCodeChan
	code2 := <-statusCodeChan

	// Exactly one request must succeed (200 OK) and the concurrent duplicate must receive 409 Conflict
	has200 := (code1 == http.StatusOK || code2 == http.StatusOK)
	has409 := (code1 == http.StatusConflict || code2 == http.StatusConflict)

	if !has200 || !has409 {
		t.Fatalf("expected one 200 OK and one 409 Conflict, got codes: %d and %d", code1, code2)
	}

	if captureCallCount != 1 {
		t.Errorf("expected PayPal CaptureOrder to be executed exactly 1 time, got %d", captureCallCount)
	}
}
