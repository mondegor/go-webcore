package mrresp

const (
	ErrorAttributeIDByDefault = "generalError" // ErrorAttributeIDByDefault - название пользовательской ошибки по умолчанию
)

type (
	// ErrorListResponse - application/json (400).
	ErrorListResponse []ErrorAttribute

	// ErrorAttribute - comment struct.
	ErrorAttribute struct {
		ID        string `json:"id"`
		Value     string `json:"value"`
		DebugInfo string `json:"debugInfo,omitempty"`
	}

	// ErrorDetailsResponse - application/problem+json (401, 403, 404, 418, 5XX).
	ErrorDetailsResponse struct {
		Title        string `json:"title"`
		Details      string `json:"details"`
		Request      string `json:"request"`
		Time         string `json:"time"`
		ErrorTraceID string `json:"errorTraceId,omitempty"`
	}
)
