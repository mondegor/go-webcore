package mrcore

import e "github.com/mondegor/go-sysmess/mrerr"

var (
	FactoryErrInternal = e.NewFactory(
		e.ErrorCodeInternal, e.ErrorTypeInternal, "internal server error")

	FactoryErrInternalNotice = e.NewFactory(
		e.ErrorCodeInternal, e.ErrorTypeInternalNotice, "internal server error")

	FactoryErrInternalNilPointer = e.NewFactory(
		"errInternalNilPointer", e.ErrorTypeInternal, "nil pointer")

	FactoryErrInternalTypeAssertion = e.NewFactory(
		"errInternalTypeAssertion", e.ErrorTypeInternal, "invalid type '{{ .type }}' assertion (value: {{ .value }})")

	FactoryErrInternalInvalidType = e.NewFactory(
		"errInternalInvalidType", e.ErrorTypeInternal, "invalid type '{{ .currentType }}', expected: '{{ .expectedType }}'")

	FactoryErrInternalFailedToOpen = e.NewFactory(
		"errInternalFailedToOpen", e.ErrorTypeInternal, "failed to open object")

	FactoryErrInternalFailedToClose = e.NewFactory(
		"errInternalFailedToClose", e.ErrorTypeInternal, "failed to close object")

	FactoryErrWithData = e.NewFactory(
		"errWithData", e.ErrorTypeInternalNotice, "{{ .key }}={{ .data }}")
)
