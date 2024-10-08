package mrreq

import (
	"net"
)

// ParseUserIP - возвращает валидный IP адрес из указанной строки или ошибку, если парсинг не удался.
func ParseUserIP(remoteAddr string) (net.IP, error) {
	ip, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return nil, ErrHttpRequestUserIP.Wrap(err, remoteAddr)
	}

	if parsedIP := net.ParseIP(ip); parsedIP != nil {
		return parsedIP, nil
	}

	return nil, ErrHttpRequestParseUserIP.New(ip)
}
