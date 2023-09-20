package mrreq

import . "github.com/mondegor/go-sysmess/mrerr"

var (
    FactoryErrHttpRequestCorrelationID = NewFactory(
        "errHttpRequestCorrelationID", ErrorKindInternal, "header 'CorrelationID' contains incorrect value '{{ .value }}'")

    FactoryErrHttpRequestUserIP = NewFactory(
        "errHttpRequestUserIP", ErrorKindInternal, "UserIP '{{ .value }}' is not IP:port")

    FactoryErrHttpRequestParseUserIP = NewFactory(
        "errHttpRequestParseUserIP", ErrorKindInternal, "UserIP contains incorrect value '{{ .value }}'")

    FactoryErrHttpRequestPlatformValue = NewFactory(
        "errHttpRequestPlatformValue", ErrorKindInternal, "header 'Platform' contains incorrect value '{{ .value }}'")
)
