package mrcore

import . "github.com/mondegor/go-sysmess/mrerr"

var (
	FactoryErrHttpResponseParseData = NewFactory(
		"errHttpResponseParseData", ErrorKindInternal, "response data is not valid")

	FactoryErrHttpResponseSendData = NewFactory(
		"errHttpResponseSendData", ErrorKindInternal, "response data is not send")

	FactoryErrHttpMultipartFormFile = NewFactory(
		"errHttpMultipartFormFile", ErrorKindSystem, "the file with the specified key '{{ .key }}' cannot be processed")

	FactoryErrHttpClientUnauthorized = NewFactory(
		"errHttpClientUnauthorized", ErrorKindUser, "401. client is unauthorized")

	FactoryErrHttpAccessForbidden = NewFactory(
		"errHttpAccessForbidden", ErrorKindUser, "403. access forbidden")

	FactoryErrHttpResourceNotFound = NewFactory(
		"errHttpResourceNotFound", ErrorKindUser, "404. resource not found")

	FactoryErrHttpRequestParseData = NewFactory(
		"errHttpRequestParseData", ErrorKindUser, "request body is not valid")

	FactoryErrHttpRequestParseParam = NewFactory(
		"errHttpRequestParseParam", ErrorKindUser, "request param with key '{{ .key }}' of type '{{ .type }}' contains incorrect value '{{ .value }}'")

	FactoryErrHttpRequestParamEmpty = NewFactory(
		"errHttpRequestParamEmpty", ErrorKindUser, "request param with key '{{ .key }}' is empty'")

	FactoryErrHttpRequestParamMax = NewFactory(
		"errHttpRequestParamMax", ErrorKindUser, "request param with key '{{ .key }}' contains value greater then max '{{ .max }}'")

	FactoryErrHttpRequestParamLenMax = NewFactory(
		"errHttpRequestParamLenMax", ErrorKindUser, "request param with key '{{ .key }}' has value length greater then max '{{ .maxLength }}' characters")
)
