package mrcore

import . "github.com/mondegor/go-sysmess/mrerr"

var (
	FactoryErrHttpResponseParseData = NewFactory(
		"errHttpResponseParseData", ErrorTypeInternal, "response data is not valid")

	FactoryErrHttpFileUpload = NewFactory(
		"errHttpFileUpload", ErrorTypeUser, "the file with the specified key '{{ .key }}' was not uploaded")

	FactoryErrHttpMultipartFormFile = NewFactory(
		"errHttpMultipartFormFile", ErrorTypeSystem, "the file with the specified key '{{ .key }}' cannot be processed")

	FactoryErrHttpClientUnauthorized = NewFactory(
		"errHttpClientUnauthorized", ErrorTypeUser, "401. client is unauthorized")

	FactoryErrHttpAccessForbidden = NewFactory(
		"errHttpAccessForbidden", ErrorTypeUser, "403. access forbidden")

	FactoryErrHttpResourceNotFound = NewFactory(
		"errHttpResourceNotFound", ErrorTypeUser, "404. resource not found")

	FactoryErrHttpRequestParseData = NewFactory(
		"errHttpRequestParseData", ErrorTypeUserWithCaller, "request body is not valid")

	FactoryErrHttpRequestParseParam = NewFactory(
		"errHttpRequestParseParam", ErrorTypeUser, "request param with key '{{ .key }}' of type '{{ .type }}' contains incorrect value '{{ .value }}'")

	FactoryErrHttpRequestParamEmpty = NewFactory(
		"errHttpRequestParamEmpty", ErrorTypeUser, "request param with key '{{ .key }}' is empty'")

	FactoryErrHttpRequestParamMax = NewFactory(
		"errHttpRequestParamMax", ErrorTypeUser, "request param with key '{{ .key }}' contains value greater then max '{{ .max }}'")

	FactoryErrHttpRequestParamLenMax = NewFactory(
		"errHttpRequestParamLenMax", ErrorTypeUser, "request param with key '{{ .key }}' has value length greater then max '{{ .maxLength }}' characters")
)
