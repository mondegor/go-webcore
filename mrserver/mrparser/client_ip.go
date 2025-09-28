package mrparser

import (
	"net"
	"net/http"
	"regexp"

	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/mrtype"

	"github.com/mondegor/go-webcore/mrserver/mrreq"
)

type (
	// ClientIP - comment struct.
	ClientIP struct {
		headers []string
		logger  mrlog.Logger
	}
)

var regexpClientIP = regexp.MustCompile(`\d+\.\d+\.\d+\.\d+`)

// NewClientIP - создаёт объект ClientIP.
func NewClientIP(logger mrlog.Logger, headers ...string) *ClientIP {
	if len(headers) == 0 {
		headers = append(headers, "X-Client-IP", "X-Cluster-Client-IP", "X-Forwarded-For")
	}

	return &ClientIP{
		headers: headers,
		logger:  logger,
	}
}

// RealIP - comment method.
func (p *ClientIP) RealIP(r *http.Request) net.IP {
	ip, err := mrreq.ParseClientIP(r.RemoteAddr)
	if err != nil {
		p.logger.Warn(r.Context(), "remote address parse error", "addr", r.RemoteAddr, "error", err)

		return net.IP{}
	}

	return ip
}

// DetailedIP - comment method.
func (p *ClientIP) DetailedIP(r *http.Request) mrtype.DetailedIP {
	realIP := p.RealIP(r)

	for _, key := range p.headers {
		header := r.Header.Get(key)

		if header == "" {
			continue
		}

		for _, value := range regexpClientIP.FindAllString(header, -1) {
			ip, err := mrreq.ParseClientIP(value)
			if err != nil || !p.isClientGlobalIP(ip) || ip.Equal(realIP) {
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

func (p *ClientIP) isClientGlobalIP(ip net.IP) bool {
	return ip.IsGlobalUnicast() &&
		!ip.IsPrivate() &&
		!ip.IsInterfaceLocalMulticast() &&
		!ip.IsLinkLocalMulticast()
}
