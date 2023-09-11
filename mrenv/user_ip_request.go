package mrenv

import (
    "net"
    "net/http"
)

func UserIpFromRequest(r *http.Request) (net.IP, error) {
    return parseIpString(r.RemoteAddr)
}

func parseIpString(value string) (net.IP, error) {
    ip, _, err := net.SplitHostPort(value)

    if err != nil {
        return nil, factoryErrHttpRequestUserIP.New(value)
    }

    parsedIp := net.ParseIP(ip)

    if parsedIp == nil {
        return nil, factoryErrHttpRequestParseUserIP.New(ip)
    }

    return parsedIp, nil
}
