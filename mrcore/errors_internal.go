package mrcore

import . "github.com/mondegor/go-sysmess/mrerr"

var (
	FactoryErrInternal = NewFactory(
		ErrorInternalID, ErrorKindInternal, "internal server error")

	FactoryErrInternalNotice = NewFactory(
		ErrorInternalID, ErrorKindInternalNotice, "internal server error")

	FactoryErrInternalNilPointer = NewFactory(
		"errInternalNilPointer", ErrorKindInternal, "nil pointer")

	FactoryErrInternalTypeAssertion = NewFactory(
		"errInternalTypeAssertion", ErrorKindInternal, "invalid type '{{ .type }}' assertion (value: {{ .value }})")

	FactoryErrInternalInvalidType = NewFactory(
		"errInternalInvalidType", ErrorKindInternal, "invalid type '{{ .currentType }}', expected: '{{ .expectedType }}'")

	FactoryErrInternalFailedToClose = NewFactory(
		"errInternalFailedToClose", ErrorKindInternal, "failed to close '{{ .name }}'")

	FactoryErrInternalWithData = NewFactory(
		"errInternalWithData", ErrorKindInternal, "{{ .key }}={{ .data }}")

	FactoryErrWithData = NewFactory(
		"errWithData", ErrorKindInternalNotice, "{{ .key }}={{ .data }}")
)
