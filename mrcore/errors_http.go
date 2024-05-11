package mrcore

import e "github.com/mondegor/go-sysmess/mrerr"

var (
	FactoryErrHTTPResponseParseData = e.NewFactory(
		"errHttpResponseParseData", e.ErrorTypeInternal, "response data is not valid")

	FactoryErrHTTPFileUpload = e.NewFactory(
		"errHttpFileUpload", e.ErrorTypeUser, "the file with the specified key '{{ .key }}' was not uploaded")

	FactoryErrHTTPMultipartFormFile = e.NewFactory(
		"errHttpMultipartFormFile", e.ErrorTypeSystem, "the file with the specified key '{{ .key }}' cannot be processed")

	FactoryErrHTTPClientUnauthorized = e.NewFactory(
		"errHttpClientUnauthorized", e.ErrorTypeUser, "401. client is unauthorized")

	FactoryErrHTTPAccessForbidden = e.NewFactory(
		"errHttpAccessForbidden", e.ErrorTypeUser, "403. access forbidden")

	FactoryErrHTTPResourceNotFound = e.NewFactory(
		"errHttpResourceNotFound", e.ErrorTypeUser, "404. resource not found")

	FactoryErrHTTPRequestParseData = e.NewFactory(
		"errHttpRequestParseData", e.ErrorTypeUserWithCaller, "request body is not valid")

	FactoryErrHTTPRequestParseParam = e.NewFactory(
		"errHttpRequestParseParam", e.ErrorTypeUser, "request param with key '{{ .key }}' of type '{{ .type }}' contains incorrect value '{{ .value }}'")

	FactoryErrHTTPRequestParamEmpty = e.NewFactory(
		"errHttpRequestParamEmpty", e.ErrorTypeUser, "request param with key '{{ .key }}' is empty'")

	FactoryErrHTTPRequestParamMax = e.NewFactory(
		"errHttpRequestParamMax", e.ErrorTypeUser, "request param with key '{{ .key }}' contains value greater then max '{{ .max }}'")

	FactoryErrHTTPRequestParamLenMax = e.NewFactory(
		"errHttpRequestParamLenMax", e.ErrorTypeUser, "request param with key '{{ .key }}' has value length greater then max '{{ .maxLength }}' characters")
)
