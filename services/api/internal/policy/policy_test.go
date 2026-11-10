package policy_test

import (
	"testing"

	"github.com/paybridge/paybridge-ai/services/api/internal/policy"
)

func TestPolicyAuthorization(t *testing.T) {
	tests := []struct {
		name        string
		role        string
		action      policy.Action
		expectAllow bool
	}{
		{"sender can create case", "sender", policy.ActionCreateCase, true},
		{"receiver cannot create case", "receiver", policy.ActionCreateCase, false},
		{"merchant_admin can create case", "merchant_admin", policy.ActionCreateCase, true},
		{"sender can capture payment", "sender", policy.ActionCapturePayment, true},
		{"analyst cannot capture payment", "merchant_analyst", policy.ActionCapturePayment, false},
		{"admin can execute purge", "merchant_admin", policy.ActionExecutePurge, true},
		{"sender cannot execute purge", "sender", policy.ActionExecutePurge, false},
		{"demo_operator can run demo scenario", "demo_operator", policy.ActionRunDemoScenario, true},
		{"unknown role denied", "untrusted_role", policy.ActionReadCase, false},
		{"ai_agent strictly denied from capturing payment", "ai_agent", policy.ActionCapturePayment, false},
		{"ai_agent strictly denied from creating orders", "ai_agent", policy.ActionCreateOrder, false},
		{"investigator_agent strictly denied from capturing payment", "investigator_agent", policy.ActionCapturePayment, false},
		{"ai_agent strictly denied from purge", "ai_agent", policy.ActionExecutePurge, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := policy.Authorize(tt.role, tt.action)
			if tt.expectAllow && err != nil {
				t.Errorf("expected allowed, got error: %v", err)
			}
			if !tt.expectAllow && err == nil {
				t.Errorf("expected denial, but was authorized")
			}
		})
	}
}
