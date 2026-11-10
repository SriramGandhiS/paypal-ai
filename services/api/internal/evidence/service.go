package evidence

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/paybridge/paybridge-ai/services/api/internal/policy"
	"github.com/paybridge/paybridge-ai/services/api/internal/store"
	"github.com/paybridge/paybridge-ai/services/api/internal/tenancy"
)

// Service provides evidence store operations.
type Service struct {
	store *store.Store
}

// NewService constructs an evidence.Service instance.
func NewService(st *store.Store) *Service {
	return &Service{store: st}
}

// Ingest records an immutable evidence item under active tenant RLS.
func (s *Service) Ingest(ctx context.Context, in IngestInput) (*EvidenceItem, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	tenantID, err := tenancy.TenantFromContext(ctx)
	if err != nil {
		return nil, err
	}

	_, role, err := tenancy.UserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if err := policy.Authorize(role, policy.ActionIngestEvidence); err != nil {
		return nil, err
	}

	hash := ComputePayloadHash(in.RawPayload)
	var item EvidenceItem

	err = s.store.WithTx(ctx, func(tx *sql.Tx) error {
		if err := s.store.SetTenantContext(ctx, tx, tenantID); err != nil {
			return err
		}

		query := `
			INSERT INTO evidence_items (
				tenant_id, case_id, kind, tier, source_trust, environment,
				verified, untrusted_text, occurred_at, redacted_extract,
				payload_sha256, access_class
			) VALUES (
				$1, $2, $3, $4, $5, 'REAL_SANDBOX',
				TRUE, $6, $7, $8, $9, $10
			)
			RETURNING id, tenant_id, case_id, kind, tier, source_trust, environment,
			          verified, untrusted_text, occurred_at, observed_at, redacted_extract,
			          payload_sha256, access_class;
		`
		redactedBytes := in.RedactedJSON
		if len(redactedBytes) == 0 {
			redactedBytes = []byte("{}")
		}

		return tx.QueryRowContext(ctx, query,
			tenantID, in.CaseID, in.Kind, in.Tier, in.SourceTrust,
			in.UntrustedText, in.OccurredAt, redactedBytes, hash, in.AccessClass,
		).Scan(
			&item.ID, &item.TenantID, &item.CaseID, &item.Kind, &item.Tier,
			&item.SourceTrust, &item.Environment, &item.Verified, &item.UntrustedText,
			&item.OccurredAt, &item.ObservedAt, &item.RedactedJSON, &item.PayloadSHA256,
			&item.AccessClass,
		)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to ingest evidence: %w", err)
	}
	return &item, nil
}

// GetEvidence retrieves an individual evidence item by ID.
func (s *Service) GetEvidence(ctx context.Context, id string) (*EvidenceItem, error) {
	tenantID, err := tenancy.TenantFromContext(ctx)
	if err != nil {
		return nil, err
	}

	_, role, err := tenancy.UserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if err := policy.Authorize(role, policy.ActionReadEvidence); err != nil {
		return nil, err
	}

	query := `
		SELECT id, tenant_id, case_id, kind, tier, source_trust, environment,
		       verified, untrusted_text, occurred_at, observed_at, redacted_extract,
		       payload_sha256, access_class
		  FROM evidence_items
		 WHERE id = $1 AND tenant_id = $2;
	`
	var item EvidenceItem
	err = s.store.DB().QueryRowContext(ctx, query, id, tenantID).Scan(
		&item.ID, &item.TenantID, &item.CaseID, &item.Kind, &item.Tier,
		&item.SourceTrust, &item.Environment, &item.Verified, &item.UntrustedText,
		&item.OccurredAt, &item.ObservedAt, &item.RedactedJSON, &item.PayloadSHA256,
		&item.AccessClass,
	)
	if err == sql.ErrNoRows {
		return nil, ErrEvidenceNotFound
	} else if err != nil {
		return nil, fmt.Errorf("failed to query evidence: %w", err)
	}

	return &item, nil
}

// ListEvidenceForCase returns all evidence items associated with a case.
func (s *Service) ListEvidenceForCase(ctx context.Context, caseID string) ([]EvidenceItem, error) {
	tenantID, err := tenancy.TenantFromContext(ctx)
	if err != nil {
		return nil, err
	}

	_, role, err := tenancy.UserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if err := policy.Authorize(role, policy.ActionReadEvidence); err != nil {
		return nil, err
	}

	query := `
		SELECT id, tenant_id, case_id, kind, tier, source_trust, environment,
		       verified, untrusted_text, occurred_at, observed_at, redacted_extract,
		       payload_sha256, access_class
		  FROM evidence_items
		 WHERE case_id = $1 AND tenant_id = $2
		 ORDER BY observed_at ASC;
	`
	rows, err := s.store.DB().QueryContext(ctx, query, caseID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to list evidence: %w", err)
	}
	defer rows.Close()

	var items []EvidenceItem
	for rows.Next() {
		var item EvidenceItem
		if err := rows.Scan(
			&item.ID, &item.TenantID, &item.CaseID, &item.Kind, &item.Tier,
			&item.SourceTrust, &item.Environment, &item.Verified, &item.UntrustedText,
			&item.OccurredAt, &item.ObservedAt, &item.RedactedJSON, &item.PayloadSHA256,
			&item.AccessClass,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}
