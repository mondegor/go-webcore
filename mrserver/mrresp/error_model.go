package mrresp

const (
	// ErrorAttributeIDByDefault - название пользовательской ошибки по умолчанию.
	ErrorAttributeIDByDefault = "FailedToProcessError"
)

type (
	// ErrorListResponse - используется для формирования ответа application/json (400).
	ErrorListResponse []ErrorAttribute

	// ErrorAttribute - пользовательская ошибка с идентификатором и её значением.
	ErrorAttribute struct {
		ID        string `json:"id"`
		Value     string `json:"value"`
		DebugInfo string `json:"debugInfo,omitempty"`
	}

	// ErrorDetailsResponse - application/problem+json (401, 403, 404, 409, 422, 5XX).
	ErrorDetailsResponse struct {
		Title        string `json:"title"`
		Details      string `json:"details"`
		Request      string `json:"request"`
		Time         string `json:"time"`
		ErrorTraceID string `json:"errorTraceId,omitempty"`
	}
)

// NewErrorAttribute - создаёт объект ErrorAttribute.
func NewErrorAttribute(code, value, debugInfo string) ErrorAttribute {
	if code == "" {
		code = ErrorAttributeIDByDefault
	}

	attr := ErrorAttribute{
		ID:        code,
		Value:     value,
		DebugInfo: debugInfo,
	}

	return attr
}
