package mrcore

import . "github.com/mondegor/go-sysmess/mrerr"

var (
    FactoryErrInternal = NewFactory(
        ErrorInternalID, ErrorKindInternal, "internal server error")

    FactoryErrInternalNilPointer = NewFactory(
        "errInternalNilPointer", ErrorKindInternal, "nil pointer")

    FactoryErrInternalTypeAssertion = NewFactory(
        "errInternalTypeAssertion", ErrorKindInternal, "invalid type '{{ .type }}' assertion (value: {{ .value }})")

    FactoryErrInternalInvalidType = NewFactory(
        "errInternalInvalidType", ErrorKindInternal, "invalid type '{{ .type1 }}', expected: '{{ .type2 }}'")

    FactoryErrInternalInvalidData = NewFactory(
        "errInternalInvalidData", ErrorKindInternal, "invalid data '{{ .value }}'")

    FactoryErrInternalParseData = NewFactory(
        "errInternalParseData", ErrorKindInternal, "data '{{ .name1 }}' parsed to {{ .name2 }} with error")

    FactoryErrInternalFailedToClose = NewFactory(
        "errInternalFailedToClose", ErrorKindInternal, "failed to close '{{ .name }}'")

    FactoryErrInternalMapValueNotFound = NewFactory(
        "errInternalMapValueNotFound", ErrorKindInternal, "'{{ .value }}' is not found in map {{ .name }}")

    FactoryErrInternalNoticeDataContainer = NewFactory(
        "errInternalNoticeDataContainer", ErrorKindInternalNotice, "data: '{{ .value }}'")
)
