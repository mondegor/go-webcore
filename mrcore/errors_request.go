package mrcore

import . "github.com/mondegor/go-sysmess/mrerr"

var (
    FactoryErrHttpRequestParamLen = NewFactory(
        "errHttpRequestParamLen", ErrorKindUser, "request param with key '{{ .key }}' has value length greater then max '{{ .maxLength }}'")

    FactoryErrHttpRequestParseParam = NewFactory(
        "errHttpRequestParseParam", ErrorKindUser, "request param of type '{{ .type }}' with key '{{ .key }}' contains incorrect value '{{ .value }}'")
)
