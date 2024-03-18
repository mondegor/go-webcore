package mrreq

import . "github.com/mondegor/go-sysmess/mrerr"

var (
	FactoryErrHttpRequestCorrelationID = NewFactory(
		"errHttpRequestCorrelationID", ErrorKindInternal, "header 'X-Correlation-Id' contains incorrect value '{{ .value }}'")

	FactoryErrHttpRequestUserIP = NewFactory(
		"errHttpRequestUserIP", ErrorKindInternal, "UserIP '{{ .value }}' is not IP:port")

	FactoryErrHttpRequestParseUserIP = NewFactory(
		"errHttpRequestParseUserIP", ErrorKindInternal, "UserIP contains incorrect value '{{ .value }}'")
)
