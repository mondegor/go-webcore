package mrcore

import (
	"github.com/mondegor/go-sysmess/mrerr"
)

var (
	// ErrInternal - internal error,
	// обобщённая внутренняя ошибка системы которая может быть решена только силами разработки.
	// Для неё всегда должен формироваться стек вызовов и посылаться событие о её создании.
	ErrInternal = mrerr.NewProto(
		mrerr.ErrorCodeInternal, mrerr.ErrorKindInternal, "internal error")

	// ErrInternalWithDetails - internal error с дополнительными подробностями.
	// Для неё всегда должен формироваться стек вызовов и посылаться событие о её создании.
	ErrInternalWithDetails = mrerr.NewProto(
		mrerr.ErrorCodeInternal, mrerr.ErrorKindInternal, "internal error: {{ .details }}")

	// ErrUnexpectedInternal - unexpected internal error,
	// особая ошибка, в которую система заворачивает все ошибки отличные от типов AppError, ProtoAppError.
	// Для неё не имеет смысла формировать стек вызовов, но всегда должно посылаться событие о её создании.
	// При возникновении этой ошибки нужно найти место причины её возникновения и написать для него обработку с указанием конкретной ошибки.
	ErrUnexpectedInternal = mrerr.NewProto(
		"errUnexpectedInternal", mrerr.ErrorKindInternal, "unexpected internal error")

	// ErrInternalNilPointer - unexpected nil pointer.
	ErrInternalNilPointer = mrerr.NewProto(
		"errInternalNilPointer", mrerr.ErrorKindInternal, "unexpected nil pointer")

	// ErrInternalCaughtPanic - caught panic.
	ErrInternalCaughtPanic = mrerr.NewProto(
		"errInternalCaughtPanic", mrerr.ErrorKindInternal, "{{ .source }}; panic: {{ .recover }}; callstack: {{ .callstack }}")

	// ErrInternalTypeAssertion - invalid type assertion.
	ErrInternalTypeAssertion = mrerr.NewProto(
		"errInternalTypeAssertion", mrerr.ErrorKindInternal, "invalid type '{{ .type }}' assertion (value: {{ .value }})")

	// ErrInternalInvalidType - invalid type, expected.
	ErrInternalInvalidType = mrerr.NewProto(
		"errInternalInvalidType", mrerr.ErrorKindInternal, "invalid type '{{ .currentType }}', expected: '{{ .expectedType }}'")

	// ErrInternalUnhandledDefaultCase - unhandled default case.
	ErrInternalUnhandledDefaultCase = mrerr.NewProto(
		"errInternalUnhandledDefaultCase", mrerr.ErrorKindInternal, "unhandled default case")

	// ErrInternalKeyNotFoundInSource - key is not found in source.
	ErrInternalKeyNotFoundInSource = mrerr.NewProto(
		"errInternalKeyNotFoundInSource", mrerr.ErrorKindInternal, "key '{{ .key }}' is not found in source {{ .source }}")

	// ErrInternalTimeoutPeriodHasExpired - the timeout period has expired.
	ErrInternalTimeoutPeriodHasExpired = mrerr.NewProto(
		"errInternalTimeoutPeriodHasExpired", mrerr.ErrorKindSystem, "the timeout period has expired")

	// ErrInternalFailedToOpen - failed to open object.
	ErrInternalFailedToOpen = mrerr.NewProto(
		"errInternalFailedToOpen", mrerr.ErrorKindInternal, "failed to open object")

	// ErrInternalFailedToClose - failed to close object.
	ErrInternalFailedToClose = mrerr.NewProto(
		"errInternalFailedToClose", mrerr.ErrorKindInternal, "failed to close object")

	// ErrInternalUnexpectedEOF - unexpected EOF.
	ErrInternalUnexpectedEOF = mrerr.NewProto(
		"errInternalUnexpectedEOF", mrerr.ErrorKindInternal, "unexpected EOF")

	// ErrInternalValueLenMax - value has length greater than max characters.
	ErrInternalValueLenMax = mrerr.NewProto(
		"errInternalValueLenMax", mrerr.ErrorKindInternal, "value has length '{{ .curLength }}' greater then max '{{ .maxLength }}' characters")

	// ErrInternalValueNotMatchRegexpPattern - specified value does not match regexp pattern.
	ErrInternalValueNotMatchRegexpPattern = mrerr.NewProto(
		"errInternalValueNotMatchRegexpPattern", mrerr.ErrorKindInternal, "specified value '{{ .value }}' does not match regexp pattern '{{ .pattern }}'")
)
