package mrresp

import (
	"net/http"
	"time"
)

const (
	// ErrorAttributeCodeByDefault - код ошибки по умолчанию для неклассифицированных ошибок.
	// Используется когда у ошибки нет собственного кода (Code).
	ErrorAttributeCodeByDefault = "FailedToProcessError"

	// timeLayoutRFC3339Milli - формат RFC3339 с точностью до миллисекунд; для UTC-времени
	// смещение выводится как суффикс Z (например 2020-01-01T12:00:00.000Z).
	timeLayoutRFC3339Milli = "2006-01-02T15:04:05.000Z07:00"
)

type (
	// ErrorDetailsResponse - модель пользовательской ошибки в формате RFC 9457 (Problem Details for HTTP APIs).
	// Используется для ответов с кодами: 401, 403, 404, 409, 422, 5xx.
	// Content-Type: application/problem+json.
	ErrorDetailsResponse struct {
		// Type - URL с описанием типа проблемы.
		Type string `json:"type,omitempty"`

		// Title - краткое описание проблемы.
		Title string `json:"title"`

		// Status - HTTP-код ответа.
		Status int `json:"status"`

		// Detail - подробное описание проблемы.
		Detail string `json:"detail,omitempty"`

		// Instance - идентификатор конкретного запроса (METHOD path).
		Instance string `json:"instance"`

		// Time - время возникновения ошибки в RFC3339 с точностью до миллисекунд (UTC).
		Time string `json:"time"`

		// ErrorTraceID - идентификатор трассировки ошибки для поиска в логах.
		ErrorTraceID string `json:"error_trace_id,omitempty"`

		// DebugInfo - отладочная информация (только в debug-режиме).
		DebugInfo string `json:"debug_info,omitempty"`
	}

	// Error400Response - модель пользовательской ошибки для кода 400 Bad Request.
	// Используется для ответов с валидацией полей.
	// Content-Type: application/json.
	Error400Response struct {
		// Status - HTTP-код ответа (всегда 400).
		Status int `json:"status"`

		// Instance - идентификатор запроса (METHOD path).
		Instance string `json:"instance"`

		// Time - время возникновения ошибки в RFC3339 с точностью до миллисекунд (UTC).
		Time string `json:"time"`

		// Errors - список ошибок валидации с кодами и описаниями.
		Errors []ErrorAttribute `json:"errors"`

		// DebugInfo - отладочная информация (только в debug-режиме).
		DebugInfo string `json:"debug_info,omitempty"`
	}

	// ErrorAttribute - атрибут отдельной пользовательской ошибки с кодом и описанием.
	// Используется в Error400Response для перечисления ошибок валидации.
	ErrorAttribute struct {
		// Code - уникальный код ошибки (например: имя поля или параметра).
		Code string `json:"code"`

		// Detail - описание ошибки.
		Detail string `json:"detail"`

		// DebugInfo - отладочная информация (только в debug-режиме).
		DebugInfo string `json:"debug_info,omitempty"`
	}
)

// NewError400Response - создаёт ответ с ошибкой валидации полей.
//
// Автоматически устанавливает статус 400, Instance и текущее время UTC.
func NewError400Response(r *http.Request, errorAttrs ...ErrorAttribute) Error400Response {
	return Error400Response{
		Status:   http.StatusBadRequest,
		Instance: r.Method + " " + r.URL.Path,
		Time:     time.Now().UTC().Format(timeLayoutRFC3339Milli),
		Errors:   errorAttrs,
	}
}
