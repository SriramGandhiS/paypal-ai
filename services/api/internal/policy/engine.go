package policy

import (
	"errors"
	"strings"
)

var (
	ErrUnauthorizedRole = errors.New("user role is not authorized for this operation")
	ErrForbidden        = ErrUnauthorizedRole
)

type Action string

const (
	ActionCreateCase         Action = "CREATE_CASE"
	ActionReadCase           Action = "READ_CASE"
	ActionCreateOrder        Action = "CREATE_ORDER"
	ActionCapturePayment     Action = "CAPTURE_PAYMENT"
	ActionIngestEvidence     Action = "INGEST_EVIDENCE"
	ActionReadEvidence       Action = "READ_EVIDENCE"
	ActionTriggerReconcile   Action = "TRIGGER_RECONCILE"
	ActionTriggerInvestigate Action = "TRIGGER_INVESTIGATE"
	ActionManageConsent      Action = "MANAGE_CONSENT"
	ActionAcknowledgeAction  Action = "ACKNOWLEDGE_ACTION"
	ActionManageConnections  Action = "MANAGE_CONNECTIONS"
	ActionExecutePurge       Action = "EXECUTE_PURGE"
	ActionRunDemoScenario    Action = "RUN_DEMO_SCENARIO"
)

// AllowedRoles maps actions to the set of roles permitted to perform them.
var AllowedRoles = map[Action]map[string]bool{
	ActionCreateCase: {
		"sender":           true,
		"merchant_analyst": true,
		"merchant_admin":   true,
		"demo_operator":    true,
	},
	ActionReadCase: {
		"sender":            true,
		"receiver":          true,
		"merchant_analyst":  true,
		"merchant_admin":    true,
		"platform_operator": true,
		"demo_operator":     true,
	},
	ActionCreateOrder: {
		"sender":        true,
		"demo_operator": true,
	},
	ActionCapturePayment: {
		"sender":        true,
		"demo_operator": true,
	},
	ActionIngestEvidence: {
		"merchant_analyst":  true,
		"merchant_admin":    true,
		"platform_operator": true,
		"demo_operator":     true,
	},
	ActionReadEvidence: {
		"sender":            true,
		"receiver":          true,
		"merchant_analyst":  true,
		"merchant_admin":    true,
		"platform_operator": true,
		"demo_operator":     true,
	},
	ActionTriggerReconcile: {
		"sender":           true,
		"merchant_analyst": true,
		"merchant_admin":   true,
		"demo_operator":    true,
	},
	ActionTriggerInvestigate: {
		"sender":            true,
		"merchant_analyst":  true,
		"merchant_admin":    true,
		"platform_operator": true,
		"demo_operator":     true,
	},
	ActionManageConsent: {
		"sender":         true,
		"receiver":       true,
		"merchant_admin": true,
	},
	ActionAcknowledgeAction: {
		"sender":         true,
		"receiver":       true,
		"merchant_admin": true,
	},
	ActionManageConnections: {
		"merchant_admin":    true,
		"platform_operator": true,
	},
	ActionExecutePurge: {
		"merchant_admin": true,
	},
	ActionRunDemoScenario: {
		"demo_operator": true,
	},
}

// Authorize checks if a given user role is permitted to execute an action.
func Authorize(role string, action Action) error {
	normalizedRole := strings.ToLower(strings.TrimSpace(role))
	allowed, exists := AllowedRoles[action]
	if !exists || !allowed[normalizedRole] {
		return ErrUnauthorizedRole
	}
	return nil
}
