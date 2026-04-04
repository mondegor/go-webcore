package mrresp

import (
	"net/http"
	"time"
)

const (
	// ErrorAttributeCodeByDefault - название пользовательской ошибки по умолчанию.
	ErrorAttributeCodeByDefault = "FailedToProcessError"
)

type (
	// ErrorDetailsResponse - RFC 9457, application/problem+json (401, 403, 404, 409, 422, 5XX).
	ErrorDetailsResponse struct {
		Type         string `json:"type,omitempty"`
		Title        string `json:"title"`
		Status       int    `json:"status"`
		Detail       string `json:"detail,omitempty"`
		Instance     string `json:"instance"`
		Time         string `json:"time"`
		ErrorTraceID string `json:"error_trace_id,omitempty"`
		DebugInfo    string `json:"debug_info,omitempty"`
	}

	// Error400Response - application/json (400).
	Error400Response struct {
		Status    int              `json:"status"`
		Instance  string           `json:"instance"`
		Time      string           `json:"time"`
		Errors    []ErrorAttribute `json:"errors"`
		DebugInfo string           `json:"debug_info,omitempty"`
	}

	// ErrorAttribute - пользовательская ошибка с идентификатором и её значением.
	ErrorAttribute struct {
		Code      string `json:"code"`
		Detail    string `json:"detail"`
		DebugInfo string `json:"debug_info,omitempty"`
	}
)

// NewError400Response - создаёт объект Error400Response.
func NewError400Response(r *http.Request, errorAttrs ...ErrorAttribute) Error400Response {
	return Error400Response{
		Status:   http.StatusBadRequest,
		Instance: r.Method + " " + r.URL.Path,
		Time:     time.Now().UTC().Format(time.RFC3339),
		Errors:   errorAttrs,
	}
}
