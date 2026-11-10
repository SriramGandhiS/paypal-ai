package httpx

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ProblemDetails represents an RFC 7807 standard error envelope.
type ProblemDetails struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
	Code     string `json:"code,omitempty"`
}

// Error implements the error interface.
func (p ProblemDetails) Error() string {
	return fmt.Sprintf("[%d] %s: %s", p.Status, p.Title, p.Detail)
}

// WriteProblem sends an RFC 7807 formatted JSON error response.
func WriteProblem(w http.ResponseWriter, p ProblemDetails) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(p.Status)
	_ = json.NewEncoder(w).Encode(p)
}

// BadRequest constructs a 400 Bad Request problem.
func BadRequest(detail, code string) ProblemDetails {
	return ProblemDetails{
		Type:   "https://paybridge.ai/errors/bad-request",
		Title:  "Bad Request",
		Status: http.StatusBadRequest,
		Detail: detail,
		Code:   code,
	}
}

// Unauthorized constructs a 401 Unauthorized problem.
func Unauthorized(detail string) ProblemDetails {
	return ProblemDetails{
		Type:   "https://paybridge.ai/errors/unauthorized",
		Title:  "Unauthorized",
		Status: http.StatusUnauthorized,
		Detail: detail,
		Code:   "UNAUTHORIZED",
	}
}

// Forbidden constructs a 403 Forbidden problem.
func Forbidden(detail string) ProblemDetails {
	return ProblemDetails{
		Type:   "https://paybridge.ai/errors/forbidden",
		Title:  "Forbidden",
		Status: http.StatusForbidden,
		Detail: detail,
		Code:   "FORBIDDEN",
	}
}

// NotFound constructs a 404 Not Found problem.
func NotFound(resource string) ProblemDetails {
	return ProblemDetails{
		Type:   "https://paybridge.ai/errors/not-found",
		Title:  "Resource Not Found",
		Status: http.StatusNotFound,
		Detail: fmt.Sprintf("The requested resource '%s' was not found", resource),
		Code:   "NOT_FOUND",
	}
}

// Conflict constructs a 409 Conflict problem.
func Conflict(detail, code string) ProblemDetails {
	return ProblemDetails{
		Type:   "https://paybridge.ai/errors/conflict",
		Title:  "Conflict",
		Status: http.StatusConflict,
		Detail: detail,
		Code:   code,
	}
}

// InternalError constructs a 500 Internal Server Error problem.
func InternalError(detail string) ProblemDetails {
	return ProblemDetails{
		Type:   "https://paybridge.ai/errors/internal-server-error",
		Title:  "Internal Server Error",
		Status: http.StatusInternalServerError,
		Detail: detail,
		Code:   "INTERNAL_ERROR",
	}
}

// ServiceUnavailable constructs a 503 Service Unavailable problem.
func ServiceUnavailable(detail string) ProblemDetails {
	return ProblemDetails{
		Type:   "https://paybridge.ai/errors/service-unavailable",
		Title:  "Service Unavailable",
		Status: http.StatusServiceUnavailable,
		Detail: detail,
		Code:   "SERVICE_UNAVAILABLE",
	}
}
