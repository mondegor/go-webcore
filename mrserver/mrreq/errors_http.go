package mrreq

import "github.com/mondegor/go-sysmess/mrerr"

var (
	// ErrHttpRequestCorrelationID - header 'X-Correlation-Id' contains incorrect value.
	ErrHttpRequestCorrelationID = mrerr.NewProto(
		"errHttpRequestCorrelationID", mrerr.ErrorKindInternal, "header 'X-Correlation-Id' contains incorrect value '{{ .value }}'")

	// ErrHttpRequestUserIP - userIP is not IP:port.
	ErrHttpRequestUserIP = mrerr.NewProto(
		"errHttpRequestUserIP", mrerr.ErrorKindInternal, "userIP '{{ .value }}' is not IP:port")

	// ErrHttpRequestParseUserIP - userIP contains incorrect value.
	ErrHttpRequestParseUserIP = mrerr.NewProto(
		"errHttpRequestParseUserIP", mrerr.ErrorKindInternal, "userIP contains incorrect value '{{ .value }}'")
)
