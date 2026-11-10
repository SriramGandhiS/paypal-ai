package checkout

import (
	"errors"
	"time"
)

var (
	ErrActiveAttemptExists = errors.New("an active, non-terminal payment attempt already exists for this case")
	ErrAttemptNotFound     = errors.New("payment attempt not found")
	ErrInvalidStatusChange = errors.New("invalid attempt status transition")
	ErrCaptureInFlight     = errors.New("capture operation already in progress for this order")
)

// Attempt represents a provider-level transaction attempt.
type Attempt struct {
	ID                        string    `json:"id"`
	TenantID                  string    `json:"tenant_id"`
	CaseID                    string    `json:"case_id"`
	AttemptNo                 int       `json:"attempt_no"`
	Provider                  string    `json:"provider"`
	Environment               string    `json:"environment"`
	OrderOperationRequestID   string    `json:"order_operation_request_id"`
	CaptureOperationRequestID *string   `json:"capture_operation_request_id,omitempty"`
	Status                    string    `json:"status"`
	Terminal                  bool      `json:"terminal"`
	CreatedAt                 time.Time `json:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at"`
}

// IsTerminal checks if the attempt status has reached a terminal state.
func IsTerminal(status string) bool {
	switch status {
	case "CANCELLED", "CAPTURE_DENIED", "APPLICATION_RECONCILED":
		return true
	default:
		return false
	}
}
