package cases

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/paybridge/paybridge-ai/services/api/internal/policy"
	"github.com/paybridge/paybridge-ai/services/api/internal/store"
	"github.com/paybridge/paybridge-ai/services/api/internal/tenancy"
)

// Service provides payment case management logic.
type Service struct {
	store *store.Store
}

// NewService constructs a cases.Service instance.
func NewService(st *store.Store) *Service {
	return &Service{store: st}
}

// CreateCase initializes a new case under the active tenant context.
func (s *Service) CreateCase(ctx context.Context, in CreateCaseInput) (*Case, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	tenantID, err := tenancy.TenantFromContext(ctx)
	if err != nil {
		return nil, err
	}

	userID, role, err := tenancy.UserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if err := policy.Authorize(role, policy.ActionCreateCase); err != nil {
		return nil, err
	}

	var created Case
	err = s.store.WithTx(ctx, func(tx *sql.Tx) error {
		// Set transaction-local tenant context for RLS
		if err := s.store.SetTenantContext(ctx, tx, tenantID); err != nil {
			return err
		}

		query := `
			INSERT INTO payment_cases (
				tenant_id, external_reference, sender_user_id, sender_country,
				receiver_country, payment_purpose, amount, currency, environment
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, 'REAL_SANDBOX'
			)
			RETURNING id, tenant_id, external_reference, sender_user_id, sender_country,
			          receiver_country, payment_purpose, amount::text, currency, environment,
			          provider_state, last_confirmed_provider_state, application_state,
			          receiver_state, case_status, state_version, created_at, updated_at;
		`
		amountDecimal := float64(in.AmountCents) / 100.0

		return tx.QueryRowContext(ctx, query,
			tenantID, in.ExternalReference, userID, in.SenderCountry,
			in.ReceiverCountry, in.PaymentPurpose, amountDecimal, in.Currency,
		).Scan(
			&created.ID, &created.TenantID, &created.ExternalReference,
			&created.SenderUserID, &created.SenderCountry, &created.ReceiverCountry,
			&created.PaymentPurpose, &created.Amount, &created.Currency,
			&created.Environment, &created.ProviderState, &created.LastConfirmedProviderState,
			&created.ApplicationState, &created.ReceiverState, &created.CaseStatus,
			&created.StateVersion, &created.CreatedAt, &created.UpdatedAt,
		)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create payment case: %w", err)
	}
	return &created, nil
}

// GetCase retrieves a specific case within the active tenant.
func (s *Service) GetCase(ctx context.Context, id string) (*Case, error) {
	tenantID, err := tenancy.TenantFromContext(ctx)
	if err != nil {
		return nil, err
	}

	_, role, err := tenancy.UserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if err := policy.Authorize(role, policy.ActionReadCase); err != nil {
		return nil, err
	}

	var c Case
	query := `
		SELECT id, tenant_id, external_reference, sender_user_id, receiver_user_id,
		       sender_country, receiver_country, payment_purpose, amount::text, currency,
		       environment, provider_state, last_confirmed_provider_state,
		       application_state, receiver_state, case_status, state_version,
		       created_at, updated_at
		  FROM payment_cases
		 WHERE id = $1 AND tenant_id = $2;
	`
	err = s.store.DB().QueryRowContext(ctx, query, id, tenantID).Scan(
		&c.ID, &c.TenantID, &c.ExternalReference, &c.SenderUserID, &c.ReceiverUserID,
		&c.SenderCountry, &c.ReceiverCountry, &c.PaymentPurpose, &c.Amount, &c.Currency,
		&c.Environment, &c.ProviderState, &c.LastConfirmedProviderState,
		&c.ApplicationState, &c.ReceiverState, &c.CaseStatus, &c.StateVersion,
		&c.CreatedAt, &c.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrCaseNotFound
	} else if err != nil {
		return nil, fmt.Errorf("failed to query case: %w", err)
	}
	return &c, nil
}

// ListCases returns paginated summaries of payment cases for the active tenant.
func (s *Service) ListCases(ctx context.Context, filter ListFilter) ([]CaseSummary, int, error) {
	tenantID, err := tenancy.TenantFromContext(ctx)
	if err != nil {
		return nil, 0, err
	}

	_, role, err := tenancy.UserFromContext(ctx)
	if err != nil {
		return nil, 0, err
	}

	if err := policy.Authorize(role, policy.ActionReadCase); err != nil {
		return nil, 0, err
	}

	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	var countQuery string
	var listQuery string
	var args []any

	if filter.Status != "" {
		countQuery = "SELECT count(*) FROM payment_cases WHERE tenant_id = $1 AND case_status = $2"
		listQuery = `
			SELECT id, external_reference, amount::text, currency,
			       provider_state, application_state, case_status,
			       state_version, updated_at
			  FROM payment_cases
			 WHERE tenant_id = $1 AND case_status = $2
			 ORDER BY created_at DESC
			 LIMIT $3 OFFSET $4;
		`
		args = []any{tenantID, filter.Status, filter.Limit, filter.Offset}
	} else {
		countQuery = "SELECT count(*) FROM payment_cases WHERE tenant_id = $1"
		listQuery = `
			SELECT id, external_reference, amount::text, currency,
			       provider_state, application_state, case_status,
			       state_version, updated_at
			  FROM payment_cases
			 WHERE tenant_id = $1
			 ORDER BY created_at DESC
			 LIMIT $2 OFFSET $3;
		`
		args = []any{tenantID, filter.Limit, filter.Offset}
	}

	var total int
	if filter.Status != "" {
		err = s.store.DB().QueryRowContext(ctx, countQuery, tenantID, filter.Status).Scan(&total)
	} else {
		err = s.store.DB().QueryRowContext(ctx, countQuery, tenantID).Scan(&total)
	}
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count cases: %w", err)
	}

	rows, err := s.store.DB().QueryContext(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list cases: %w", err)
	}
	defer rows.Close()

	var items []CaseSummary
	for rows.Next() {
		var cs CaseSummary
		if err := rows.Scan(
			&cs.ID, &cs.ExternalReference, &cs.Amount, &cs.Currency,
			&cs.ProviderState, &cs.ApplicationState, &cs.CaseStatus,
			&cs.StateVersion, &cs.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		items = append(items, cs)
	}

	return items, total, rows.Err()
}
