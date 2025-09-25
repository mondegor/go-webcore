package mrresp

import (
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrerr/mr"

	"github.com/mondegor/go-webcore/mrcore"
)

const (
	// ErrorAttributeIDByDefault - название пользовательской ошибки по умолчанию.
	ErrorAttributeIDByDefault = "GeneralError"
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

	// ErrorDetailsResponse - application/problem+json (401, 403, 404, 418, 422, 5XX).
	ErrorDetailsResponse struct {
		Title        string `json:"title"`
		Details      string `json:"details"`
		Request      string `json:"request"`
		Time         string `json:"time"`
		ErrorTraceID string `json:"errorTraceId,omitempty"`
	}
)

// NewErrorAttribute - создаёт объект ErrorAttribute.
func NewErrorAttribute(lz mrcore.Localizer, err error, withDebugInfo bool) ErrorAttribute {
	var (
		errCode    string
		customCode string
	)

	if e, ok := err.(*mrerr.CustomError); ok { //nolint:errorlint
		customCode = "/" + e.CustomCode()
		err = e.Err()
	}

	e := mr.CastOrWrapUnexpectedInternal(err)

	if e.Code() != "" {
		errCode = e.Code()
	} else {
		errCode = ErrorAttributeIDByDefault
	}

	attr := ErrorAttribute{
		ID:    errCode + customCode,
		Value: lz.TranslateError(e),
	}

	if withDebugInfo {
		attr.DebugInfo = e.Error()
	}

	return attr
}
