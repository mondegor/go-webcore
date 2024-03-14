package mrcore

import . "github.com/mondegor/go-sysmess/mrerr"

var (
	FactoryErrInternal = NewFactory(
		ErrorCodeInternal, ErrorKindInternal, "internal server error")

	FactoryErrInternalNotice = NewFactory(
		ErrorCodeInternal, ErrorKindInternalNotice, "internal server error")

	FactoryErrInternalNilPointer = NewFactory(
		"errInternalNilPointer", ErrorKindInternal, "nil pointer")

	FactoryErrInternalTypeAssertion = NewFactory(
		"errInternalTypeAssertion", ErrorKindInternal, "invalid type '{{ .type }}' assertion (value: {{ .value }})")

	FactoryErrInternalInvalidType = NewFactory(
		"errInternalInvalidType", ErrorKindInternal, "invalid type '{{ .currentType }}', expected: '{{ .expectedType }}'")

	FactoryErrInternalFailedToOpen = NewFactory(
		"errInternalFailedToOpen", ErrorKindInternal, "failed to open object")

	FactoryErrInternalFailedToClose = NewFactory(
		"errInternalFailedToClose", ErrorKindInternal, "failed to close object")

	FactoryErrWithData = NewFactory(
		"errWithData", ErrorKindInternalNotice, "{{ .key }}={{ .data }}")
)
