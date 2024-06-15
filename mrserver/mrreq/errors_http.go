package mrreq

import "github.com/mondegor/go-sysmess/mrerr"

var (
	// ErrHttpRequestCorrelationID - comment error.
	ErrHttpRequestCorrelationID = mrerr.NewProto(
		"errHttpRequestCorrelationID", mrerr.ErrorKindInternal, "header 'X-Correlation-Id' contains incorrect value '{{ .value }}'")

	// ErrHttpRequestUserIP - comment error.
	ErrHttpRequestUserIP = mrerr.NewProto(
		"errHttpRequestUserIP", mrerr.ErrorKindInternal, "UserIP '{{ .value }}' is not IP:port")

	// ErrHttpRequestParseUserIP - comment error.
	ErrHttpRequestParseUserIP = mrerr.NewProto(
		"errHttpRequestParseUserIP", mrerr.ErrorKindInternal, "UserIP contains incorrect value '{{ .value }}'")
)
