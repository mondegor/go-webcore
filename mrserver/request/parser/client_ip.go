package parser

import (
	"net/http"
	"net/netip"
	"strings"

	"github.com/mondegor/go-core/mrlog"
	"github.com/mondegor/go-core/mrtype"
	"github.com/mondegor/go-core/mrtype/parse"

	"github.com/mondegor/go-webcore/mrserver"
)

type (
	// ClientIP - определяет реальный и прокси IP-адрес клиента.
	ClientIP struct {
		proxyHeaders []string
		logger       mrlog.Logger
	}
)

// NewClientIP - создаёт объект ClientIP.
// Заголовки просматриваются в порядке их перечисления, поэтому первым указывается
// наиболее достоверный источник. Если proxyHeaders не заданы, используется список
// по умолчанию: X-Real-Ip, X-Forwarded-For. Свой список полностью заменяет список
// по умолчанию, а не дополняет его.
func NewClientIP(
	logger mrlog.Logger,
	proxyHeaders ...string,
) *ClientIP {
	if len(proxyHeaders) == 0 {
		proxyHeaders = []string{mrserver.HeaderKeyRealIP, mrserver.HeaderKeyForwardedFor}
	}

	return &ClientIP{
		proxyHeaders: proxyHeaders,
		logger:       logger,
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

	for _, key := range p.proxyHeaders {
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
