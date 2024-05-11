package mrreq

import (
	"net"
	"net/http"
)

func ParseUserIP(r *http.Request) (net.IP, error) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return nil, FactoryErrHTTPRequestUserIP.New(r.RemoteAddr)
	}

	parsedIP := net.ParseIP(ip)

	if parsedIP == nil {
		return nil, FactoryErrHTTPRequestParseUserIP.New(ip)
	}

	return parsedIP, nil
}
