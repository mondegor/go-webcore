package mrcore

import (
	"github.com/mondegor/go-sysmess/mrerr"
)

var (
	// ErrInternal - comment error.
	ErrInternal = mrerr.NewProto(
		"errInternal", mrerr.ErrorKindInternal, "internal error")

	// ErrInternalNilPointer - comment error.
	ErrInternalNilPointer = mrerr.NewProto(
		"errInternalNilPointer", mrerr.ErrorKindInternal, "nil pointer")

	// ErrInternalTypeAssertion - comment error.
	ErrInternalTypeAssertion = mrerr.NewProto(
		"errInternalTypeAssertion", mrerr.ErrorKindInternal, "invalid type '{{ .type }}' assertion (value: {{ .value }})")

	// ErrInternalInvalidType - comment error.
	ErrInternalInvalidType = mrerr.NewProto(
		"errInternalInvalidType", mrerr.ErrorKindInternal, "invalid type '{{ .currentType }}', expected: '{{ .expectedType }}'")

	// ErrInternalKeyNotFoundInSource - comment error.
	ErrInternalKeyNotFoundInSource = mrerr.NewProto(
		"errInternalKeyNotFoundInSource", mrerr.ErrorKindInternal, "key '{{ .key }}' is not found in source {{ .source }}")

	// ErrInternalFailedToOpen - comment error.
	ErrInternalFailedToOpen = mrerr.NewProto(
		"errInternalFailedToOpen", mrerr.ErrorKindInternal, "failed to open object")

	// ErrInternalFailedToClose - comment error.
	ErrInternalFailedToClose = mrerr.NewProto(
		"errInternalFailedToClose", mrerr.ErrorKindInternal, "failed to close object")

	// ErrInternalValueLenMax - comment error.
	ErrInternalValueLenMax = mrerr.NewProto(
		"errInternalValueLenMax", mrerr.ErrorKindInternal, "value has length '{{ .curLength }}' greater then max '{{ .maxLength }}' characters")

	// ErrInternalValueNotMatchRegexpPattern - comment error.
	ErrInternalValueNotMatchRegexpPattern = mrerr.NewProto(
		"errInternalValueNotMatchRegexpPattern", mrerr.ErrorKindInternal, "specified value '{{ .value }}' does not match regexp pattern '{{ .pattern }}'")
)
