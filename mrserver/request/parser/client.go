package parser

import (
	"net/http"
	"net/netip"
	"strings"

	"github.com/mondegor/go-core/mrlog"
	"github.com/mondegor/go-core/mrtype"
	"github.com/mondegor/go-core/mrtype/parse"
	"github.com/mondegor/go-core/util/xstrings"

	"github.com/mondegor/go-webcore/mrserver"
)

const (
	// defaultUserAgentMaxLength - предельная длина User-Agent по умолчанию (в символах).
	defaultUserAgentMaxLength = 512
)

type (
	// Client - определяет сведения о клиенте из запроса: реальный и прокси IP-адрес,
	// User-Agent. Все они - недоверенный ввод, контролируемый клиентом.
	Client struct {
		proxyHeaders       []string
		userAgentMaxLength int
		logger             mrlog.Logger
	}

	// ClientOptions - опции для создания Client.
	ClientOptions struct {
		// ProxyHeaders - заголовки с прокси IP-адресом клиента (пусто - список по умолчанию).
		ProxyHeaders []string

		// UserAgentMaxLength - предельная длина User-Agent в символах (0 - длина по умолчанию).
		UserAgentMaxLength int
	}
)

// NewClient - создаёт объект Client.
// Заголовки ProxyHeaders просматриваются в порядке их перечисления, поэтому первым указывается
// наиболее достоверный источник. Если они не заданы, используется список
// по умолчанию: X-Real-Ip, X-Forwarded-For. Свой список полностью заменяет список
// по умолчанию, а не дополняет его.
func NewClient(logger mrlog.Logger, opts ClientOptions) *Client {
	if len(opts.ProxyHeaders) == 0 {
		opts.ProxyHeaders = []string{mrserver.HeaderKeyRealIP, mrserver.HeaderKeyForwardedFor}
	}

	if opts.UserAgentMaxLength < 1 {
		opts.UserAgentMaxLength = defaultUserAgentMaxLength
	}

	return &Client{
		proxyHeaders:       opts.ProxyHeaders,
		userAgentMaxLength: opts.UserAgentMaxLength,
		logger:             logger,
	}
}

// RealIP - возвращает реальный IP адрес клиента из RemoteAddr.
func (p *Client) RealIP(r *http.Request) netip.Addr {
	ip, err := parse.IP(r.RemoteAddr, true)
	if err != nil {
		p.logger.Warn(r.Context(), "remote address parse error", "addr", r.RemoteAddr, "error", err)

		return netip.Addr{}
	}

	return ip
}

// DetailedIP - возвращает детальную информацию об IP (реальный и прокси).
func (p *Client) DetailedIP(r *http.Request) mrtype.DetailedIP {
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

func (p *Client) isClientGlobalIP(ip netip.Addr) bool {
	return ip.IsGlobalUnicast() &&
		!ip.IsPrivate() &&
		!ip.IsInterfaceLocalMulticast() &&
		!ip.IsLinkLocalMulticast()
}

// UserAgent - возвращает User-Agent клиента, приведённый к безопасному для хранения и вывода виду.
func (p *Client) UserAgent(r *http.Request) string {
	return xstrings.SanitizePrintable(r.UserAgent(), p.userAgentMaxLength)
}
