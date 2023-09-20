package mrcore

import . "github.com/mondegor/go-sysmess/mrerr"

var (
    FactoryErrHttpRequestParamEmpty = NewFactory(
        "errHttpRequestParamEmpty", ErrorKindUser, "request param with key '{{ .key }}' is empty'")

    FactoryErrHttpRequestParamLen = NewFactory(
        "errHttpRequestParamLen", ErrorKindUser, "request param with key '{{ .key }}' has value length greater then max '{{ .maxLength }}'")

    FactoryErrHttpRequestParseParam = NewFactory(
        "errHttpRequestParseParam", ErrorKindUser, "request param of type '{{ .type }}' with key '{{ .key }}' contains incorrect value '{{ .value }}'")

    FactoryErrHttpRequestParseData = NewFactory(
        "errHttpRequestParseData", ErrorKindUser, "request body is not valid")

    FactoryErrHttpResponseParseData = NewFactory(
        "errHttpResponseParseData", ErrorKindInternal, "response data is not valid")

    FactoryErrHttpResponseSendData = NewFactory(
        "errHttpResponseSendData", ErrorKindInternal, "response data is not send")

    FactoryErrHttpResponseSystemTemporarilyUnableToProcess = NewFactory(
       "errHttpResponseSystemTemporarilyUnableToProcess", ErrorKindUser, "the system is temporarily unable to process your request. please try again")

    FactoryErrHttpResourceNotFound = NewFactory(
        "errHttpResourceNotFound", ErrorKindUser, "resource not found")
)
