package mrreq

import (
	"net"
	"net/http"
)

func ParseUserIp(r *http.Request) (net.IP, error) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)

	if err != nil {
		return nil, FactoryErrHttpRequestUserIP.New(r.RemoteAddr)
	}

	parsedIp := net.ParseIP(ip)

	if parsedIp == nil {
		return nil, FactoryErrHttpRequestParseUserIP.New(ip)
	}

	return parsedIp, nil
}
