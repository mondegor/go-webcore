package mrenv

import . "github.com/mondegor/go-sysmess/mrerr"

var (
    factoryErrHttpRequestPlatformValue = NewFactory(
        "errHttpRequestPlatformValue", ErrorKindInternal, "header 'Platform' contains incorrect value '{{ .value }}'")

    factoryErrHttpRequestCorrelationID = NewFactory(
        "errHttpRequestCorrelationID", ErrorKindInternalNotice, "header 'CorrelationID' contains incorrect value '{{ .value }}'")

    factoryErrHttpRequestUserIP = NewFactory(
        "errHttpRequestUserIP", ErrorKindInternal, "UserIP '{{ .value }}' is not IP:port")

    factoryErrHttpRequestParseUserIP = NewFactory(
        "errHttpRequestParseUserIP", ErrorKindInternal, "UserIP contains incorrect value '{{ .value }}'")
)
