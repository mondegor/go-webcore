package mrreq

import . "github.com/mondegor/go-sysmess/mrerr"

var (
	FactoryErrHttpRequestCorrelationID = NewFactory(
		"errHttpRequestCorrelationID", ErrorKindInternalNotice, "header 'X-Correlation-Id' contains incorrect value '{{ .value }}'")

	FactoryErrHttpRequestUserIP = NewFactory(
		"errHttpRequestUserIP", ErrorKindInternalNotice, "UserIP '{{ .value }}' is not IP:port")

	FactoryErrHttpRequestParseUserIP = NewFactory(
		"errHttpRequestParseUserIP", ErrorKindInternalNotice, "UserIP contains incorrect value '{{ .value }}'")
)
