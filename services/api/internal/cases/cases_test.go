package cases_test

import (
	"context"
	"testing"

	"github.com/paybridge/paybridge-ai/services/api/internal/cases"
	"github.com/paybridge/paybridge-ai/services/api/internal/store"
	"github.com/paybridge/paybridge-ai/services/api/internal/tenancy"
)

func TestCreateCaseInputValidation(t *testing.T) {
	tests := []struct {
		name    string
		input   cases.CreateCaseInput
		wantErr error
	}{
		{
			name: "valid input",
			input: cases.CreateCaseInput{
				AmountCents:     15000,
				Currency:        "USD",
				SenderCountry:   "US",
				ReceiverCountry: "IN",
				PaymentPurpose:  "SOFTWARE_SERVICES",
			},
			wantErr: nil,
		},
		{
			name: "zero amount error",
			input: cases.CreateCaseInput{
				AmountCents:     0,
				Currency:        "USD",
				SenderCountry:   "US",
				ReceiverCountry: "IN",
			},
			wantErr: cases.ErrInvalidAmount,
		},
		{
			name: "negative amount error",
			input: cases.CreateCaseInput{
				AmountCents:     -500,
				Currency:        "USD",
				SenderCountry:   "US",
				ReceiverCountry: "IN",
			},
			wantErr: cases.ErrInvalidAmount,
		},
		{
			name: "invalid currency length",
			input: cases.CreateCaseInput{
				AmountCents:     1000,
				Currency:        "US",
				SenderCountry:   "US",
				ReceiverCountry: "IN",
			},
			wantErr: cases.ErrInvalidCurrency,
		},
		{
			name: "invalid country code",
			input: cases.CreateCaseInput{
				AmountCents:     1000,
				Currency:        "USD",
				SenderCountry:   "USA",
				ReceiverCountry: "IN",
			},
			wantErr: cases.ErrInvalidCountry,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Validate()
			if err != tt.wantErr {
				t.Errorf("expected %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestCasesContextRequirement(t *testing.T) {
	st := store.NewFromDB(nil)
	svc := cases.NewService(st)

	t.Run("missing tenant context returns error", func(t *testing.T) {
		ctx := context.Background()
		_, err := svc.CreateCase(ctx, cases.CreateCaseInput{
			AmountCents:     5000,
			Currency:        "USD",
			SenderCountry:   "US",
			ReceiverCountry: "CA",
			PaymentPurpose:  "SERVICES",
		})
		if err != tenancy.ErrNoTenantInContext {
			t.Errorf("expected ErrNoTenantInContext, got %v", err)
		}
	})

	t.Run("unauthorized role returns policy error", func(t *testing.T) {
		ctx := tenancy.WithTenant(context.Background(), "018f3a21-9b34-7a12-89cd-456789000000")
		ctx = tenancy.WithUser(ctx, "018f3a21-9b34-7a12-89cd-456789111111", "receiver") // receiver cannot create case

		_, err := svc.CreateCase(ctx, cases.CreateCaseInput{
			AmountCents:     5000,
			Currency:        "USD",
			SenderCountry:   "US",
			ReceiverCountry: "CA",
			PaymentPurpose:  "SERVICES",
		})
		if err == nil {
			t.Fatal("expected error for unauthorized role, got nil")
		}
	})
}
