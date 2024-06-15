package mrreq

import (
	"net"
	"net/http"
)

// ParseUserIP  - comment func.
func ParseUserIP(r *http.Request) (net.IP, error) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return nil, ErrHttpRequestUserIP.Wrap(err, r.RemoteAddr)
	}

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return nil, ErrHttpRequestParseUserIP.New(ip)
	}

	return parsedIP, nil
}
