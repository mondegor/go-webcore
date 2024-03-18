package mrcore

import . "github.com/mondegor/go-sysmess/mrerr"

var (
	FactoryErrInternal = NewFactoryWithCaller(
		ErrorCodeInternal, ErrorKindInternal, "internal server error")

	FactoryErrInternalNotice = NewFactory(
		ErrorCodeInternal, ErrorKindInternal, "internal server error")

	FactoryErrInternalNilPointer = NewFactoryWithCaller(
		"errInternalNilPointer", ErrorKindInternal, "nil pointer")

	FactoryErrInternalTypeAssertion = NewFactoryWithCaller(
		"errInternalTypeAssertion", ErrorKindInternal, "invalid type '{{ .type }}' assertion (value: {{ .value }})")

	FactoryErrInternalInvalidType = NewFactoryWithCaller(
		"errInternalInvalidType", ErrorKindInternal, "invalid type '{{ .currentType }}', expected: '{{ .expectedType }}'")

	FactoryErrInternalFailedToOpen = NewFactoryWithCaller(
		"errInternalFailedToOpen", ErrorKindInternal, "failed to open object")

	FactoryErrInternalFailedToClose = NewFactoryWithCaller(
		"errInternalFailedToClose", ErrorKindInternal, "failed to close object")

	FactoryErrWithData = NewFactory(
		"errWithData", ErrorKindInternal, "{{ .key }}={{ .data }}")
)
