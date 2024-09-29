package mrcore

import (
	"github.com/mondegor/go-sysmess/mrerr"
)

var (
	// ErrHttpResponseParseData - response data is not valid.
	ErrHttpResponseParseData = mrerr.NewProto(
		"errHttpResponseParseData", mrerr.ErrorKindInternal, "response data is not valid")

	// ErrHttpFileUpload - the file with the specified key was not uploaded.
	ErrHttpFileUpload = mrerr.NewProto(
		"errHttpFileUpload", mrerr.ErrorKindUser, "the file with the specified key '{{ .key }}' was not uploaded")

	// ErrHttpMultipartFormFile - the file with the specified key cannot be processed.
	ErrHttpMultipartFormFile = mrerr.NewProto(
		"errHttpMultipartFormFile", mrerr.ErrorKindSystem, "the file with the specified key '{{ .key }}' cannot be processed")

	// ErrHttpClientUnauthorized - 401. client is unauthorized.
	ErrHttpClientUnauthorized = mrerr.NewProto(
		"errHttpClientUnauthorized", mrerr.ErrorKindUser, "401. client is unauthorized")

	// ErrHttpAccessForbidden - 403. access forbidden.
	ErrHttpAccessForbidden = mrerr.NewProto(
		"errHttpAccessForbidden", mrerr.ErrorKindUser, "403. access forbidden")

	// ErrHttpResourceNotFound - 404. resource not found.
	ErrHttpResourceNotFound = mrerr.NewProto(
		"errHttpResourceNotFound", mrerr.ErrorKindUser, "404. resource not found")

	// ErrHttpRequestParseData - request body is not valid.
	ErrHttpRequestParseData = mrerr.NewProto(
		"errHttpRequestParseData", mrerr.ErrorKindInternal, "request body is not valid")

	// ErrHttpRequestParseParam - request param with key of type contains incorrect value.
	ErrHttpRequestParseParam = mrerr.NewProto(
		"errHttpRequestParseParam", mrerr.ErrorKindUser, "request param with key '{{ .key }}' of type '{{ .type }}' contains incorrect value '{{ .value }}'")

	// ErrHttpRequestParamEmpty - request param with key is empty.
	ErrHttpRequestParamEmpty = mrerr.NewProto(
		"errHttpRequestParamEmpty", mrerr.ErrorKindUser, "request param with key '{{ .key }}' is empty")

	// ErrHttpRequestParamMax - request param with key contains value greater than max.
	ErrHttpRequestParamMax = mrerr.NewProto(
		"errHttpRequestParamMax", mrerr.ErrorKindUser, "request param with key '{{ .key }}' contains value greater then max '{{ .max }}'")

	// ErrHttpRequestParamLenMax - request param with key has value length greater than max characters.
	ErrHttpRequestParamLenMax = mrerr.NewProto(
		"errHttpRequestParamLenMax", mrerr.ErrorKindUser, "request param with key '{{ .key }}' has value length greater then max '{{ .maxLength }}' characters")
)
