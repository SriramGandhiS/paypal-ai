package api

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/paybridge/paybridge-ai/services/api/internal/auth"
	"github.com/paybridge/paybridge-ai/services/api/internal/cases"
	"github.com/paybridge/paybridge-ai/services/api/internal/checkout"
	"github.com/paybridge/paybridge-ai/services/api/internal/evidence"
	"github.com/paybridge/paybridge-ai/services/api/internal/httpx"
	"github.com/paybridge/paybridge-ai/services/api/internal/paypal"
	"github.com/paybridge/paybridge-ai/services/api/internal/policy"
	"github.com/paybridge/paybridge-ai/services/api/internal/reconciliation"
	"github.com/paybridge/paybridge-ai/services/api/internal/tenancy"
	"github.com/paybridge/paybridge-ai/services/api/internal/webhook"
)

// PayPalService abstracts PayPal orders and captures for payment flow execution.
type PayPalService interface {
	CreateOrder(ctx context.Context, ord paypal.OrderRequest) (*paypal.OrderResponse, error)
	CaptureOrder(ctx context.Context, orderID string) (*paypal.OrderResponse, error)
	GetOrder(ctx context.Context, orderID string) (*paypal.OrderResponse, error)
	GetCapture(ctx context.Context, captureID string) (*paypal.CaptureGroundTruth, error)
}

// RouterDependencies groups the required application services.
type RouterDependencies struct {
	Cases          *cases.Service
	Checkout       *checkout.Service
	Evidence       *evidence.Service
	Reconciler     *reconciliation.Reconciler
	PayPal         PayPalService
	WebhookGateway *webhook.Gateway
	AuthMiddleware *auth.Middleware
	Logger         *slog.Logger
}

// RegisterRoutes registers all PayBridge AI REST API endpoints onto the router.
func RegisterRoutes(router *httpx.Router, deps RouterDependencies) {
	if deps.Logger == nil {
		deps.Logger = slog.Default()
	}

	h := &handler{
		deps:   deps,
		logger: deps.Logger,
	}

	// Unauthenticated Webhook Gateway
	if deps.WebhookGateway != nil {
		router.HandleFunc("POST /api/v1/webhooks/paypal", deps.WebhookGateway.Handler())
	}

	// Authenticated Cases API
	authWrap := func(hf http.HandlerFunc) http.HandlerFunc {
		if deps.AuthMiddleware == nil {
			return hf
		}
		return deps.AuthMiddleware.Authenticate(hf).ServeHTTP
	}

	router.HandleFunc("POST /api/v1/cases", authWrap(h.handleCreateCase))
	router.HandleFunc("GET /api/v1/cases", authWrap(h.handleListCases))
	router.HandleFunc("GET /api/v1/cases/{id}", authWrap(h.handleGetCase))

	// Payment Attempts
	router.HandleFunc("POST /api/v1/cases/{id}/attempts", authWrap(h.handleCreateAttempt))
	router.HandleFunc("GET /api/v1/cases/{id}/attempts", authWrap(h.handleListAttempts))

	// Evidence Ingestion & Retrieval
	router.HandleFunc("POST /api/v1/cases/{id}/evidence", authWrap(h.handleIngestEvidence))
	router.HandleFunc("GET /api/v1/cases/{id}/evidence", authWrap(h.handleListEvidence))

	// Reconciliation Execution
	router.HandleFunc("POST /api/v1/cases/{id}/reconcile", authWrap(h.handleReconcileCase))

	// PayPal Golden Path (Create Order + Capture)
	router.HandleFunc("POST /api/v1/cases/{id}/checkout/paypal", authWrap(h.handlePayPalCheckout))
	router.HandleFunc("POST /api/v1/cases/{id}/checkout/paypal/capture", authWrap(h.handlePayPalCapture))
}

type handler struct {
	deps   RouterDependencies
	logger *slog.Logger
}

// POST /api/v1/cases
func (h *handler) handleCreateCase(w http.ResponseWriter, r *http.Request) {
	var in cases.CreateCaseInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httpx.WriteProblem(w, httpx.BadRequest("malformed json request body", "INVALID_BODY"))
		return
	}

	// Check idempotency key header
	if in.IdempotencyKey == "" {
		in.IdempotencyKey = r.Header.Get("Idempotency-Key")
	}

	c, err := h.deps.Cases.CreateCase(r.Context(), in)
	if err != nil {
		if errors.Is(err, cases.ErrInvalidAmount) || errors.Is(err, cases.ErrInvalidCurrency) || errors.Is(err, cases.ErrInvalidCountry) {
			httpx.WriteProblem(w, httpx.BadRequest(err.Error(), "VALIDATION_FAILED"))
			return
		}
		if errors.Is(err, policy.ErrForbidden) {
			httpx.WriteProblem(w, httpx.Forbidden("insufficient permissions to create case"))
			return
		}
		h.logger.Error("failed to create case", "error", err)
		httpx.WriteProblem(w, httpx.InternalError("failed to create payment case"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(c)
}

// GET /api/v1/cases
func (h *handler) handleListCases(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	limit, _ := strconv.Atoi(query.Get("limit"))
	offset, _ := strconv.Atoi(query.Get("offset"))

	filter := cases.ListFilter{
		Status: query.Get("status"),
		Limit:  limit,
		Offset: offset,
	}

	items, total, err := h.deps.Cases.ListCases(r.Context(), filter)
	if err != nil {
		if errors.Is(err, policy.ErrForbidden) {
			httpx.WriteProblem(w, httpx.Forbidden("insufficient permissions to view cases"))
			return
		}
		h.logger.Error("failed to list cases", "error", err)
		httpx.WriteProblem(w, httpx.InternalError("failed to query payment cases"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"items": items,
		"total": total,
	})
}

// GET /api/v1/cases/{id}
func (h *handler) handleGetCase(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		httpx.WriteProblem(w, httpx.BadRequest("missing case ID in path", "MISSING_PATH_PARAM"))
		return
	}

	c, err := h.deps.Cases.GetCase(r.Context(), id)
	if err != nil {
		if errors.Is(err, cases.ErrCaseNotFound) {
			httpx.WriteProblem(w, httpx.NotFound("payment case"))
			return
		}
		if errors.Is(err, policy.ErrForbidden) {
			httpx.WriteProblem(w, httpx.Forbidden("insufficient permissions to view case"))
			return
		}
		h.logger.Error("failed to get case", "id", id, "error", err)
		httpx.WriteProblem(w, httpx.InternalError("failed to retrieve case"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(c)
}

// POST /api/v1/cases/{id}/attempts
func (h *handler) handleCreateAttempt(w http.ResponseWriter, r *http.Request) {
	caseID := r.PathValue("id")
	if caseID == "" {
		httpx.WriteProblem(w, httpx.BadRequest("missing case ID in path", "MISSING_PATH_PARAM"))
		return
	}

	var body struct {
		OrderOperationRequestID string `json:"order_operation_request_id"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)

	reqID := body.OrderOperationRequestID
	if reqID == "" {
		reqID = generateRandomHex(16)
	}

	attempt, err := h.deps.Checkout.CreateAttempt(r.Context(), caseID, reqID)
	if err != nil {
		if errors.Is(err, checkout.ErrActiveAttemptExists) {
			httpx.WriteProblem(w, httpx.Conflict("an active payment attempt already exists for this case", "ACTIVE_ATTEMPT_EXISTS"))
			return
		}
		if errors.Is(err, policy.ErrForbidden) {
			httpx.WriteProblem(w, httpx.Forbidden("insufficient permissions to initiate payment"))
			return
		}
		h.logger.Error("failed to create payment attempt", "case_id", caseID, "error", err)
		httpx.WriteProblem(w, httpx.InternalError("failed to create payment attempt"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(attempt)
}

// GET /api/v1/cases/{id}/attempts
func (h *handler) handleListAttempts(w http.ResponseWriter, r *http.Request) {
	caseID := r.PathValue("id")
	if caseID == "" {
		httpx.WriteProblem(w, httpx.BadRequest("missing case ID in path", "MISSING_PATH_PARAM"))
		return
	}

	attempts, err := h.deps.Checkout.ListAttemptsForCase(r.Context(), caseID)
	if err != nil {
		if errors.Is(err, policy.ErrForbidden) {
			httpx.WriteProblem(w, httpx.Forbidden("insufficient permissions to view attempts"))
			return
		}
		h.logger.Error("failed to list attempts", "case_id", caseID, "error", err)
		httpx.WriteProblem(w, httpx.InternalError("failed to query attempts"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"items": attempts,
		"total": len(attempts),
	})
}

// POST /api/v1/cases/{id}/evidence
func (h *handler) handleIngestEvidence(w http.ResponseWriter, r *http.Request) {
	caseID := r.PathValue("id")
	if caseID == "" {
		httpx.WriteProblem(w, httpx.BadRequest("missing case ID in path", "MISSING_PATH_PARAM"))
		return
	}

	// Enforce 1MB size limit
	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)
	rawBytes, err := io.ReadAll(r.Body)
	if err != nil {
		httpx.WriteProblem(w, httpx.BadRequest("failed to read evidence payload", "INVALID_PAYLOAD"))
		return
	}

	var req struct {
		Kind          string          `json:"kind"`
		Tier          int             `json:"tier"`
		SourceTrust   string          `json:"source_trust"`
		RedactedJSON  json.RawMessage `json:"redacted_extract"`
		UntrustedText bool            `json:"untrusted_text"`
		AccessClass   string          `json:"access_class"`
	}
	_ = json.Unmarshal(rawBytes, &req)

	now := time.Now().UTC()
	in := evidence.IngestInput{
		CaseID:        caseID,
		Kind:          req.Kind,
		Tier:          req.Tier,
		SourceTrust:   req.SourceTrust,
		OccurredAt:    &now,
		RawPayload:    rawBytes,
		RedactedJSON:  req.RedactedJSON,
		UntrustedText: req.UntrustedText,
		AccessClass:   req.AccessClass,
	}

	item, err := h.deps.Evidence.Ingest(r.Context(), in)
	if err != nil {
		if errors.Is(err, evidence.ErrInvalidTier) || errors.Is(err, evidence.ErrSimulatedTierRule) || errors.Is(err, evidence.ErrMissingPayload) {
			httpx.WriteProblem(w, httpx.BadRequest(err.Error(), "INVALID_EVIDENCE_SPEC"))
			return
		}
		if errors.Is(err, policy.ErrForbidden) {
			httpx.WriteProblem(w, httpx.Forbidden("insufficient permissions to ingest evidence"))
			return
		}
		h.logger.Error("failed to ingest evidence", "case_id", caseID, "error", err)
		httpx.WriteProblem(w, httpx.InternalError("failed to record evidence item"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(item)
}

// GET /api/v1/cases/{id}/evidence
func (h *handler) handleListEvidence(w http.ResponseWriter, r *http.Request) {
	caseID := r.PathValue("id")
	if caseID == "" {
		httpx.WriteProblem(w, httpx.BadRequest("missing case ID in path", "MISSING_PATH_PARAM"))
		return
	}

	items, err := h.deps.Evidence.ListEvidenceForCase(r.Context(), caseID)
	if err != nil {
		if errors.Is(err, policy.ErrForbidden) {
			httpx.WriteProblem(w, httpx.Forbidden("insufficient permissions to read evidence"))
			return
		}
		h.logger.Error("failed to list evidence", "case_id", caseID, "error", err)
		httpx.WriteProblem(w, httpx.InternalError("failed to query evidence items"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"items": items,
		"total": len(items),
	})
}

// POST /api/v1/cases/{id}/reconcile
func (h *handler) handleReconcileCase(w http.ResponseWriter, r *http.Request) {
	caseID := r.PathValue("id")
	if caseID == "" {
		httpx.WriteProblem(w, httpx.BadRequest("missing case ID in path", "MISSING_PATH_PARAM"))
		return
	}

	tenantID, err := tenancy.TenantFromContext(r.Context())
	if err != nil {
		httpx.WriteProblem(w, httpx.Unauthorized("tenant context required"))
		return
	}

	var truth reconciliation.CaptureGroundTruth
	if err := json.NewDecoder(r.Body).Decode(&truth); err != nil {
		httpx.WriteProblem(w, httpx.BadRequest("malformed ground truth body", "INVALID_TRUTH"))
		return
	}

	outcome, err := h.deps.Reconciler.ReconcileCase(r.Context(), tenantID, caseID, truth)
	if err != nil {
		if errors.Is(err, reconciliation.ErrCaseNotFound) {
			httpx.WriteProblem(w, httpx.NotFound("payment case"))
			return
		}
		h.logger.Error("reconciliation failed", "case_id", caseID, "error", err)
		httpx.WriteProblem(w, httpx.InternalError("failed to reconcile case"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(outcome)
}

// POST /api/v1/cases/{id}/checkout/paypal
func (h *handler) handlePayPalCheckout(w http.ResponseWriter, r *http.Request) {
	caseID := r.PathValue("id")
	if caseID == "" {
		httpx.WriteProblem(w, httpx.BadRequest("missing case ID in path", "MISSING_PATH_PARAM"))
		return
	}

	_, role, err := tenancy.UserFromContext(r.Context())
	if err != nil {
		httpx.WriteProblem(w, httpx.Unauthorized("user context required"))
		return
	}
	if err := policy.Authorize(role, policy.ActionCreateOrder); err != nil {
		httpx.WriteProblem(w, httpx.Forbidden("role not authorized to initiate checkout"))
		return
	}

	if h.deps.PayPal == nil {
		httpx.WriteProblem(w, httpx.ServiceUnavailable("paypal client is not configured"))
		return
	}

	c, err := h.deps.Cases.GetCase(r.Context(), caseID)
	if err != nil {
		if errors.Is(err, cases.ErrCaseNotFound) {
			httpx.WriteProblem(w, httpx.NotFound("payment case"))
			return
		}
		httpx.WriteProblem(w, httpx.InternalError("failed to retrieve case"))
		return
	}

	// Guard against re-checkout on terminal/resolved cases
	if c.CaseStatus == "RESOLVED" || c.ProviderState == "CAPTURE_COMPLETED" {
		httpx.WriteProblem(w, httpx.Conflict("payment case is already captured or resolved", "CASE_ALREADY_RESOLVED"))
		return
	}

	var reqBody struct {
		AmountCents int64  `json:"amount_cents"`
		ReturnURL   string `json:"return_url"`
		CancelURL   string `json:"cancel_url"`
	}
	_ = json.NewDecoder(r.Body).Decode(&reqBody)

	amtCents := reqBody.AmountCents
	if amtCents <= 0 {
		amtFloat, _ := strconv.ParseFloat(c.Amount, 64)
		amtCents = int64(amtFloat*100 + 0.5)
	}

	// 1. Create order on PayPal Sandbox
	ord, err := h.deps.PayPal.CreateOrder(r.Context(), paypal.OrderRequest{
		ReferenceID: c.ExternalReference,
		AmountCents: amtCents,
		Currency:    c.Currency,
		ReturnURL:   reqBody.ReturnURL,
		CancelURL:   reqBody.CancelURL,
	})
	if err != nil {
		h.logger.Error("paypal order creation failed", "case_id", caseID, "error", err)
		httpx.WriteProblem(w, httpx.ServiceUnavailable(fmt.Sprintf("paypal order creation failed: %v", err)))
		return
	}

	// 2. Record payment attempt (blind-retry guard)
	opReqID := generateRandomHex(16)
	attempt, _ := h.deps.Checkout.CreateAttempt(r.Context(), caseID, opReqID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"order":   ord,
		"attempt": attempt,
	})
}

// POST /api/v1/cases/{id}/checkout/paypal/capture
func (h *handler) handlePayPalCapture(w http.ResponseWriter, r *http.Request) {
	caseID := r.PathValue("id")
	if caseID == "" {
		httpx.WriteProblem(w, httpx.BadRequest("missing case ID in path", "MISSING_PATH_PARAM"))
		return
	}

	tenantID, err := tenancy.TenantFromContext(r.Context())
	if err != nil {
		httpx.WriteProblem(w, httpx.Unauthorized("tenant context required"))
		return
	}

	userID, role, err := tenancy.UserFromContext(r.Context())
	if err != nil {
		httpx.WriteProblem(w, httpx.Unauthorized("user context required"))
		return
	}

	if err := policy.Authorize(role, policy.ActionCapturePayment); err != nil {
		httpx.WriteProblem(w, httpx.Forbidden("role not authorized to capture payment"))
		return
	}

	if h.deps.PayPal == nil {
		httpx.WriteProblem(w, httpx.ServiceUnavailable("paypal client is not configured"))
		return
	}

	var reqBody struct {
		OrderID string `json:"order_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil || reqBody.OrderID == "" {
		httpx.WriteProblem(w, httpx.BadRequest("order_id is required", "MISSING_ORDER_ID"))
		return
	}

	// 0. Concurrency Guard: enforce distributed atomic in-flight capture lock per order ID
	if h.deps.Checkout != nil {
		release, acquired, lockErr := h.deps.Checkout.AcquireDistributedCaptureLock(r.Context(), tenantID, userID, reqBody.OrderID)
		if lockErr != nil {
			h.logger.Error("distributed capture lock error", "order_id", reqBody.OrderID, "error", lockErr)
		}
		if !acquired {
			h.logger.Warn("concurrent capture attempt rejected; lock held by another in-flight request",
				"order_id", reqBody.OrderID,
			)
			httpx.WriteProblem(w, httpx.Conflict("a capture operation is already in progress for this order", "CAPTURE_IN_FLIGHT"))
			return
		}
		defer release()
	}

	// 1. Idempotency Check: inspect existing provider order state before mutating
	existingOrder, getErr := h.deps.PayPal.GetOrder(r.Context(), reqBody.OrderID)
	var capResp *paypal.OrderResponse

	if getErr == nil && existingOrder != nil && existingOrder.Status == "COMPLETED" {
		// Order was ALREADY captured (e.g. earlier timeout or duplicate client request)
		h.logger.Info("idempotent capture: paypal order already completed", "order_id", reqBody.OrderID)
		capResp = existingOrder
	} else {
		// 2. Perform Capture on PayPal
		var capErr error
		capResp, capErr = h.deps.PayPal.CaptureOrder(r.Context(), reqBody.OrderID)
		if capErr != nil {
			// On capture error/timeout, recover by querying GET /orders/{id}
			h.logger.Warn("paypal capture returned error; inspecting provider state for recovery",
				"order_id", reqBody.OrderID, "error", capErr,
			)
			recoveredOrder, recErr := h.deps.PayPal.GetOrder(r.Context(), reqBody.OrderID)
			if recErr == nil && recoveredOrder != nil && recoveredOrder.Status == "COMPLETED" {
				h.logger.Info("recovered capture from provider GET inspection", "order_id", reqBody.OrderID)
				capResp = recoveredOrder
			} else {
				h.logger.Error("paypal capture order failed", "order_id", reqBody.OrderID, "error", capErr)
				httpx.WriteProblem(w, httpx.ServiceUnavailable(fmt.Sprintf("paypal capture failed: %v", capErr)))
				return
			}
		}
	}

	// 3. Fetch ground truth directly via GET /v2/payments/captures/{id} if capture ID available
	var truth reconciliation.CaptureGroundTruth
	if capResp.CaptureID != "" {
		pt, err := h.deps.PayPal.GetCapture(r.Context(), capResp.CaptureID)
		if err == nil && pt != nil {
			truth = reconciliation.CaptureGroundTruth{
				CaptureID:   pt.CaptureID,
				Status:      pt.Status,
				AmountCents: pt.AmountCents,
				Currency:    pt.Currency,
			}
		} else {
			truth = reconciliation.CaptureGroundTruth{
				CaptureID:   capResp.CaptureID,
				Status:      capResp.CaptureStatus,
				AmountCents: 0,
				Currency:    capResp.Currency,
			}
		}
	} else {
		truth = reconciliation.CaptureGroundTruth{
			CaptureID:   capResp.ID,
			Status:      capResp.Status,
			AmountCents: 0,
			Currency:    capResp.Currency,
		}
	}

	// 4. Ingest immutable capture evidence
	if h.deps.Evidence != nil {
		now := time.Now().UTC()
		rawPayload, _ := json.Marshal(capResp)
		_, _ = h.deps.Evidence.Ingest(r.Context(), evidence.IngestInput{
			CaseID:        caseID,
			Kind:          "PAYPAL_CAPTURE_CONFIRMATION",
			Tier:          1,
			SourceTrust:   "PROVIDER_VERIFIED",
			OccurredAt:    &now,
			RawPayload:    rawPayload,
			RedactedJSON:  rawPayload,
			UntrustedText: false,
			AccessClass:   "TENANT",
		})
	}

	// 5. Reconcile case preserving financial ground truth
	outcome, _ := h.deps.Reconciler.ReconcileCase(r.Context(), tenantID, caseID, truth)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"capture": capResp,
		"outcome": outcome,
	})
}

func generateRandomHex(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
