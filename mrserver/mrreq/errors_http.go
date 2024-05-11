package mrreq

import e "github.com/mondegor/go-sysmess/mrerr"

var (
	FactoryErrHTTPRequestCorrelationID = e.NewFactory(
		"errHttpRequestCorrelationID", e.ErrorTypeInternalNotice, "header 'X-Correlation-Id' contains incorrect value '{{ .value }}'")

	FactoryErrHTTPRequestUserIP = e.NewFactory(
		"errHttpRequestUserIP", e.ErrorTypeInternalNotice, "UserIP '{{ .value }}' is not IP:port")

	FactoryErrHTTPRequestParseUserIP = e.NewFactory(
		"errHttpRequestParseUserIP", e.ErrorTypeInternalNotice, "UserIP contains incorrect value '{{ .value }}'")
)
