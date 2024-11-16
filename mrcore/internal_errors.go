package mrcore

import (
	"github.com/mondegor/go-sysmess/mrerr"
)

var (
	// ErrUnexpectedInternal - unexpected internal error,
	// особая ошибка, в которую система заворачивает все ошибки отличные от типов AppError, ProtoAppError.
	ErrUnexpectedInternal = mrerr.NewProto(
		"errUnexpectedInternal", mrerr.ErrorKindInternal, "unexpected internal error")

	// ErrInternal - internal error,
	// обобщённая внутренняя ошибка системы которая формирует стек вызовов и посылает событие о своём создании.
	ErrInternal = mrerr.NewProto(
		mrerr.ErrorCodeInternal, mrerr.ErrorKindInternal, "internal error")

	// ErrSystem - system error,
	// обобщённая системная ошибка, которая сообщает о сетевых проблемах,
	// о работоспособности внешних ресурсов (БД, API, FileSystem).
	ErrSystem = mrerr.NewProto(
		mrerr.ErrorCodeSystem, mrerr.ErrorKindSystem, "system error")

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
