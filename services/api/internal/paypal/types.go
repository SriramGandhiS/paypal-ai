package paypal

import "time"

// Environment represents PayPal target environment
type Environment string

const (
	EnvSandbox    Environment = "sandbox"
	EnvProduction Environment = "production"
)

// Classification represents standard payment outcome taxonomy
type Classification string

const (
	ClassSuccess                Classification = "success"
	ClassExplicitDecline        Classification = "explicit_decline"
	ClassPending                Classification = "pending"
	ClassUnknownMissingResponse Classification = "unknown_missing_response"
	ClassRequestError           Classification = "request_error"
)

// RawEvidence encapsulates raw HTTP response metadata and SHA-256 for audit
type RawEvidence struct {
	StatusCode    int            `json:"status_code"`
	PayloadSHA256 string         `json:"payload_sha256"`
	RedactedBody  map[string]any `json:"redacted_body"`
	ObservedAt    time.Time      `json:"observed_at"`
}

// Order represents PayPal Order v2 schema
type Order struct {
	ID            string `json:"id"`
	Status        string `json:"status"` // CREATED, SAVED, APPROVED, VOIDED, COMPLETED, PAYER_ACTION_REQUIRED
	Intent        string `json:"intent"`
	CreateTime    string `json:"create_time"`
	UpdateTime    string `json:"update_time"`
	Links         []Link `json:"links,omitempty"`
	PurchaseUnits []struct {
		ReferenceID string `json:"reference_id"`
		Amount      struct {
			CurrencyCode string `json:"currency_code"`
			Value        string `json:"value"`
		} `json:"amount"`
		Payments *struct {
			Captures []Capture `json:"captures,omitempty"`
		} `json:"payments,omitempty"`
	} `json:"purchase_units,omitempty"`
}

// Link represents HATEOAS link
type Link struct {
	Href   string `json:"href"`
	Rel    string `json:"rel"`
	Method string `json:"method"`
}

// CaptureRef returns capture ID if present on order
func (o *Order) CaptureRef() string {
	if len(o.PurchaseUnits) > 0 && o.PurchaseUnits[0].Payments != nil && len(o.PurchaseUnits[0].Payments.Captures) > 0 {
		return o.PurchaseUnits[0].Payments.Captures[0].ID
	}
	return ""
}

// Capture represents PayPal Capture v2 schema
type Capture struct {
	ID            string `json:"id"`
	Status        string `json:"status"` // COMPLETED, DECLINED, PARTIALLY_REFUNDED, PENDING, REFUNDED, FAILED
	StatusDetails *struct {
		Reason string `json:"reason,omitempty"`
	} `json:"status_details,omitempty"`
	Amount struct {
		CurrencyCode string `json:"currency_code"`
		Value        string `json:"value"`
	} `json:"amount"`
	FinalCapture bool   `json:"final_capture"`
	CreateTime   string `json:"create_time"`
	UpdateTime   string `json:"update_time"`
}

// CreateOrderInput encapsulates order creation parameters
type CreateOrderInput struct {
	RequestID   string `json:"request_id"` // Stable PayPal-Request-Id
	Amount      string `json:"amount"`
	Currency    string `json:"currency"`
	InvoiceID   string `json:"invoice_id"`
	ReturnURL   string `json:"return_url"`
	CancelURL   string `json:"cancel_url"`
}

// CaptureInput encapsulates capture execution parameters
type CaptureInput struct {
	RequestID string `json:"request_id"` // Stable PayPal-Request-Id
	OrderID   string `json:"order_id"`
}

// WebhookVerifyInput encapsulates PayPal webhook signature verification parameters
type WebhookVerifyInput struct {
	AuthAlgo         string         `json:"auth_algo"`
	CertURL          string         `json:"cert_url"`
	TransmissionID   string         `json:"transmission_id"`
	TransmissionSig  string         `json:"transmission_sig"`
	TransmissionTime string         `json:"transmission_time"`
	WebhookID        string         `json:"webhook_id"`
	WebhookEvent     map[string]any `json:"webhook_event"`
}

// VerifyResult represents signature validation status
type VerifyResult struct {
	VerificationStatus string `json:"verification_status"` // SUCCESS or FAILURE
}
