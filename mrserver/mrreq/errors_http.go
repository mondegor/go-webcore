package mrreq

import "github.com/mondegor/go-sysmess/mrerr"

var (
	// ErrHttpRequestCorrelationID - header 'X-Correlation-Id' contains incorrect value.
	// Это вспомогательная ошибка, для неё необязательно формировать стек вызовов и отправлять событие о её создании.
	ErrHttpRequestCorrelationID = mrerr.NewProto(
		"errHttpRequestCorrelationID", mrerr.ErrorKindInternal, "header 'X-Correlation-Id' contains incorrect value '{{ .value }}'")

	// ErrHttpRequestUserIP - userIP is not IP:port.
	// Это вспомогательная ошибка, для неё необязательно формировать стек вызовов и отправлять событие о её создании.
	ErrHttpRequestUserIP = mrerr.NewProto(
		"errHttpRequestUserIP", mrerr.ErrorKindInternal, "userIP '{{ .value }}' is not IP:port")

	// ErrHttpRequestParseUserIP - userIP contains incorrect value.
	// Это вспомогательная ошибка, для неё необязательно формировать стек вызовов и отправлять событие о её создании.
	ErrHttpRequestParseUserIP = mrerr.NewProto(
		"errHttpRequestParseUserIP", mrerr.ErrorKindInternal, "userIP contains incorrect value '{{ .value }}'")
)
