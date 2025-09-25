package mrreq

import (
	"net"
	"strings"
)

// ParseClientIP - возвращает валидный IP адрес из указанной строки или ошибку, если парсинг не удался.
func ParseClientIP(value string) (ip net.IP, err error) {
	val := value

	if strings.Contains(value, ":") {
		val, _, err = net.SplitHostPort(value)
		if err != nil {
			return nil, ErrHttpRequestParseClientIP.Wrap(err, value)
		}
	}

	if ip := net.ParseIP(val); ip != nil {
		return ip, nil
	}

	return nil, ErrHttpRequestParseClientIP.New(value)
}
