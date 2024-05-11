package mrreq

import . "github.com/mondegor/go-sysmess/mrerr"

var (
	FactoryErrHttpRequestCorrelationID = NewFactory(
		"errHttpRequestCorrelationID", ErrorTypeInternalNotice, "header 'X-Correlation-Id' contains incorrect value '{{ .value }}'")

	FactoryErrHttpRequestUserIP = NewFactory(
		"errHttpRequestUserIP", ErrorTypeInternalNotice, "UserIP '{{ .value }}' is not IP:port")

	FactoryErrHttpRequestParseUserIP = NewFactory(
		"errHttpRequestParseUserIP", ErrorTypeInternalNotice, "UserIP contains incorrect value '{{ .value }}'")
)
