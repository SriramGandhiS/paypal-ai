package paypal

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
)

// ClassifyResponse maps HTTP status, body, and error into the canonical 5-outcome classification
func ClassifyResponse(statusCode int, body map[string]any, err error) Classification {
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) ||
			errors.Is(err, context.Canceled) ||
			os.IsTimeout(err) {
			return ClassUnknownMissingResponse
		}
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			return ClassUnknownMissingResponse
		}
		// Connection resets or drops after request sent
		return ClassUnknownMissingResponse
	}

	if statusCode >= 500 {
		return ClassUnknownMissingResponse
	}

	if statusCode >= 200 && statusCode < 300 {
		if status, ok := body["status"].(string); ok {
			switch status {
			case "COMPLETED":
				return ClassSuccess
			case "PENDING":
				return ClassPending
			case "APPROVED", "CREATED", "SAVED":
				return ClassSuccess
			}
		}
		return ClassSuccess
	}

	if statusCode == http.StatusBadRequest || statusCode == http.StatusUnprocessableEntity {
		if name, ok := body["name"].(string); ok {
			if name == "INSTRUMENT_DECLINED" || name == "TRANSACTION_REFUSED" || name == "PAYER_ACTION_REQUIRED" {
				return ClassExplicitDecline
			}
		}
		return ClassRequestError
	}

	if statusCode == http.StatusUnauthorized || statusCode == http.StatusForbidden {
		return ClassRequestError
	}

	return ClassUnknownMissingResponse
}
