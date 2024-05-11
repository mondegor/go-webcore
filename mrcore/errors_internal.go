package mrcore

import . "github.com/mondegor/go-sysmess/mrerr"

var (
	FactoryErrInternal = NewFactory(
		ErrorCodeInternal, ErrorTypeInternal, "internal server error")

	FactoryErrInternalNotice = NewFactory(
		ErrorCodeInternal, ErrorTypeInternalNotice, "internal server error")

	FactoryErrInternalNilPointer = NewFactory(
		"errInternalNilPointer", ErrorTypeInternal, "nil pointer")

	FactoryErrInternalTypeAssertion = NewFactory(
		"errInternalTypeAssertion", ErrorTypeInternal, "invalid type '{{ .type }}' assertion (value: {{ .value }})")

	FactoryErrInternalInvalidType = NewFactory(
		"errInternalInvalidType", ErrorTypeInternal, "invalid type '{{ .currentType }}', expected: '{{ .expectedType }}'")

	FactoryErrInternalFailedToOpen = NewFactory(
		"errInternalFailedToOpen", ErrorTypeInternal, "failed to open object")

	FactoryErrInternalFailedToClose = NewFactory(
		"errInternalFailedToClose", ErrorTypeInternal, "failed to close object")

	FactoryErrWithData = NewFactory(
		"errWithData", ErrorTypeInternalNotice, "{{ .key }}={{ .data }}")
)
