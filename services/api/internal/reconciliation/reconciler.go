package reconciliation

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"math"

	"github.com/paybridge/paybridge-ai/services/api/internal/store"
)

var (
	ErrCaseNotFound = errors.New("case not found for reconciliation")
)

// CaptureGroundTruth represents the verified provider state obtained via GET.
type CaptureGroundTruth struct {
	CaptureID   string
	Status      string // COMPLETED, PENDING, DENIED
	AmountCents int64
	Currency    string
}

// Reconciler executes Reader-only financial state reconciliation.
type Reconciler struct {
	store  *store.Store
	logger *slog.Logger
}

// NewReconciler constructs a new Reconciler.
func NewReconciler(st *store.Store, logger *slog.Logger) *Reconciler {
	if logger == nil {
		logger = slog.Default()
	}
	return &Reconciler{store: st, logger: logger}
}

// ReconcileOutcome holds the result of a reconciliation evaluation.
type ReconcileOutcome struct {
	CaseID           string
	ProviderState    string
	ApplicationState string
	CaseStatus       string
	AmountMismatch   bool
}

// EvaluateReconciliation compares merchant invoice expectations with ground truth.
func (r *Reconciler) EvaluateReconciliation(expectedAmountCents int64, truth CaptureGroundTruth) ReconcileOutcome {
	outcome := ReconcileOutcome{
		ProviderState: "UNKNOWN",
	}

	switch truth.Status {
	case "COMPLETED":
		outcome.ProviderState = "CAPTURE_COMPLETED"
	case "PENDING":
		outcome.ProviderState = "CAPTURE_PENDING"
	case "DENIED":
		outcome.ProviderState = "CAPTURE_DENIED"
	default:
		outcome.ProviderState = "UNKNOWN"
	}

	// Financial Truth Invariant:
	// If the provider captured funds, provider_state IS CAPTURE_COMPLETED.
	// Any amount mismatch MUST be surfaced in application_state as INCONSISTENT.
	if truth.Status == "COMPLETED" {
		diff := int64(math.Abs(float64(expectedAmountCents - truth.AmountCents)))
		if diff > 0 {
			outcome.AmountMismatch = true
			outcome.ApplicationState = "INCONSISTENT"
			outcome.CaseStatus = "NEEDS_HUMAN_REVIEW"
		} else {
			outcome.ApplicationState = "SYNCHRONISED"
			outcome.CaseStatus = "RESOLVED"
		}
	} else if truth.Status == "DENIED" {
		outcome.ApplicationState = "SYNCHRONISED"
		outcome.CaseStatus = "CLOSED_UNRESOLVED"
	} else {
		outcome.ApplicationState = "RECONCILING"
		outcome.CaseStatus = "OBSERVING"
	}

	return outcome
}

// ReconcileCase evaluates and applies reconciliation to a case in the database.
func (r *Reconciler) ReconcileCase(ctx context.Context, tenantID, caseID string, truth CaptureGroundTruth) (*ReconcileOutcome, error) {
	if r.store == nil || r.store.DB() == nil {
		outcome := r.EvaluateReconciliation(truth.AmountCents, truth)
		outcome.CaseID = caseID
		return &outcome, nil
	}

	var expectedAmountCents int64
	var currVersion int64

	err := r.store.WithTx(ctx, func(tx *sql.Tx) error {
		if err := r.store.SetTenantContext(ctx, tx, tenantID); err != nil {
			return err
		}

		err := tx.QueryRowContext(ctx, `
			SELECT (amount * 100)::bigint, state_version
			  FROM payment_cases
			 WHERE id = $1 AND tenant_id = $2
			   FOR UPDATE;
		`, caseID, tenantID).Scan(&expectedAmountCents, &currVersion)

		if err == sql.ErrNoRows {
			return ErrCaseNotFound
		} else if err != nil {
			return err
		}

		outcome := r.EvaluateReconciliation(expectedAmountCents, truth)
		outcome.CaseID = caseID

		// Update case state vector
		_, err = tx.ExecContext(ctx, `
			UPDATE payment_cases
			   SET provider_state = $1,
			       application_state = $2,
			       case_status = $3,
			       state_version = state_version + 1,
			       updated_at = clock_timestamp()
			 WHERE id = $4 AND tenant_id = $5;
		`, outcome.ProviderState, outcome.ApplicationState, outcome.CaseStatus, caseID, tenantID)

		return err
	})

	if err != nil {
		return nil, fmt.Errorf("reconciliation failed: %w", err)
	}

	outcome := r.EvaluateReconciliation(expectedAmountCents, truth)
	outcome.CaseID = caseID
	return &outcome, nil
}
