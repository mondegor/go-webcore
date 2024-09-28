package mrcore

import (
	"github.com/mondegor/go-sysmess/mrerr"
)

var (
	// ErrInternal - internal error.
	ErrInternal = mrerr.NewProto(
		"errInternal", mrerr.ErrorKindInternal, "internal error")

	// ErrInternalNilPointer - unexpected nil pointer.
	ErrInternalNilPointer = mrerr.NewProto(
		"errInternalNilPointer", mrerr.ErrorKindInternal, "unexpected nil pointer")

	// ErrInternalTypeAssertion - invalid type assertion.
	ErrInternalTypeAssertion = mrerr.NewProto(
		"errInternalTypeAssertion", mrerr.ErrorKindInternal, "invalid type '{{ .type }}' assertion (value: {{ .value }})")

	// ErrInternalInvalidType - invalid type, expected.
	ErrInternalInvalidType = mrerr.NewProto(
		"errInternalInvalidType", mrerr.ErrorKindInternal, "invalid type '{{ .currentType }}', expected: '{{ .expectedType }}'")

	// ErrInternalKeyNotFoundInSource - key is not found in source.
	ErrInternalKeyNotFoundInSource = mrerr.NewProto(
		"errInternalKeyNotFoundInSource", mrerr.ErrorKindInternal, "key '{{ .key }}' is not found in source {{ .source }}")

	// ErrInternalFailedToOpen - failed to open object.
	ErrInternalFailedToOpen = mrerr.NewProto(
		"errInternalFailedToOpen", mrerr.ErrorKindInternal, "failed to open object")

	// ErrInternalFailedToClose - failed to close object.
	ErrInternalFailedToClose = mrerr.NewProto(
		"errInternalFailedToClose", mrerr.ErrorKindInternal, "failed to close object")

	// ErrInternalValueLenMax - value has length greater than max characters.
	ErrInternalValueLenMax = mrerr.NewProto(
		"errInternalValueLenMax", mrerr.ErrorKindInternal, "value has length '{{ .curLength }}' greater then max '{{ .maxLength }}' characters")

	// ErrInternalValueNotMatchRegexpPattern - specified value does not match regexp pattern.
	ErrInternalValueNotMatchRegexpPattern = mrerr.NewProto(
		"errInternalValueNotMatchRegexpPattern", mrerr.ErrorKindInternal, "specified value '{{ .value }}' does not match regexp pattern '{{ .pattern }}'")
)
