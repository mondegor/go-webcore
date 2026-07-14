package parser

import (
	"net/http"
	"net/netip"
	"strings"

	"github.com/mondegor/go-core/mrlog"
	"github.com/mondegor/go-core/mrtype"
	"github.com/mondegor/go-core/mrtype/parse"
)

type (
	// ClientIP - определяет реальный и прокси IP-адрес клиента.
	ClientIP struct {
		headers []string
		logger  mrlog.Logger
	}
)

// NewClientIP - создаёт объект ClientIP.
func NewClientIP(
	logger mrlog.Logger,
	headers ...string,
) *ClientIP {
	if len(headers) == 0 {
		headers = append(headers, "X-Client-IP", "X-Cluster-Client-IP", "X-Forwarded-For")
	}

	return &ClientIP{
		headers: headers,
		logger:  logger,
	}
}

// RealIP - возвращает реальный IP адрес клиента из RemoteAddr.
func (p *ClientIP) RealIP(r *http.Request) netip.Addr {
	ip, err := parse.IP(r.RemoteAddr, true)
	if err != nil {
		p.logger.Warn(r.Context(), "remote address parse error", "addr", r.RemoteAddr, "error", err)

		return netip.Addr{}
	}

	return ip
}

// DetailedIP - возвращает детальную информацию об IP (реальный и прокси).
func (p *ClientIP) DetailedIP(r *http.Request) mrtype.DetailedIP {
	realIP := p.RealIP(r)

	for _, key := range p.headers {
		header := r.Header.Get(key)

		if header == "" {
			continue
		}

		for _, value := range strings.Split(header, ",") {
			ip, err := parse.IP(value, true)
			if err != nil || !p.isClientGlobalIP(ip) || ip == realIP {
				continue
			}

			return mrtype.DetailedIP{
				Real:  realIP,
				Proxy: ip,
			}
		}
	}

	return mrtype.DetailedIP{
		Real: realIP,
	}
}

func (p *ClientIP) isClientGlobalIP(ip netip.Addr) bool {
	return ip.IsGlobalUnicast() &&
		!ip.IsPrivate() &&
		!ip.IsInterfaceLocalMulticast() &&
		!ip.IsLinkLocalMulticast()
}
