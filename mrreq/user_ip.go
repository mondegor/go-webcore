package mrreq

import (
    "net"
    "net/http"
)

func UserIp(r *http.Request) (net.IP, error) {
    ip, _, err := net.SplitHostPort(r.RemoteAddr)

    if err != nil {
        return nil, factoryErrHttpRequestUserIP.New(r.RemoteAddr)
    }

    parsedIp := net.ParseIP(ip)

    if parsedIp == nil {
        return nil, factoryErrHttpRequestParseUserIP.New(ip)
    }

    return parsedIp, nil
}
