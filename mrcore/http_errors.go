package mrcore

import (
	"github.com/mondegor/go-sysmess/mrerr"
)

var (
	// ErrHttpResponseParseData - comment error.
	ErrHttpResponseParseData = mrerr.NewProto(
		"errHttpResponseParseData", mrerr.ErrorKindInternal, "response data is not valid")

	// ErrHttpFileUpload - comment error.
	ErrHttpFileUpload = mrerr.NewProto(
		"errHttpFileUpload", mrerr.ErrorKindUser, "the file with the specified key '{{ .key }}' was not uploaded")

	// ErrHttpMultipartFormFile - comment error.
	ErrHttpMultipartFormFile = mrerr.NewProto(
		"errHttpMultipartFormFile", mrerr.ErrorKindSystem, "the file with the specified key '{{ .key }}' cannot be processed")

	// ErrHttpClientUnauthorized - comment error.
	ErrHttpClientUnauthorized = mrerr.NewProto(
		"errHttpClientUnauthorized", mrerr.ErrorKindUser, "401. client is unauthorized")

	// ErrHttpAccessForbidden - comment error.
	ErrHttpAccessForbidden = mrerr.NewProto(
		"errHttpAccessForbidden", mrerr.ErrorKindUser, "403. access forbidden")

	// ErrHttpResourceNotFound - comment error.
	ErrHttpResourceNotFound = mrerr.NewProto(
		"errHttpResourceNotFound", mrerr.ErrorKindUser, "404. resource not found")

	// ErrHttpRequestParseData - comment error.
	ErrHttpRequestParseData = mrerr.NewProto(
		"errHttpRequestParseData", mrerr.ErrorKindInternal, "request body is not valid")

	// ErrHttpRequestParseParam - comment error.
	ErrHttpRequestParseParam = mrerr.NewProto(
		"errHttpRequestParseParam", mrerr.ErrorKindUser, "request param with key '{{ .key }}' of type '{{ .type }}' contains incorrect value '{{ .value }}'")

	// ErrHttpRequestParamEmpty - comment error.
	ErrHttpRequestParamEmpty = mrerr.NewProto(
		"errHttpRequestParamEmpty", mrerr.ErrorKindUser, "request param with key '{{ .key }}' is empty'")

	// ErrHttpRequestParamMax - comment error.
	ErrHttpRequestParamMax = mrerr.NewProto(
		"errHttpRequestParamMax", mrerr.ErrorKindUser, "request param with key '{{ .key }}' contains value greater then max '{{ .max }}'")

	// ErrHttpRequestParamLenMax - comment error.
	ErrHttpRequestParamLenMax = mrerr.NewProto(
		"errHttpRequestParamLenMax", mrerr.ErrorKindUser, "request param with key '{{ .key }}' has value length greater then max '{{ .maxLength }}' characters")
)
