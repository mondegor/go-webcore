package parser_test

import (
	"net/http"
	"net/http/httptest"
	"net/netip"
	"testing"

	"github.com/mondegor/go-core/mrlog"
	"github.com/mondegor/go-core/mrtype"
	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrserver/request"
	"github.com/mondegor/go-webcore/mrserver/request/parser"
)

// Make sure the ClientIP conforms with the request.ParserClientIP interface.
func TestClientIPImplementsRequestParserClientIP(t *testing.T) {
	t.Parallel()

	assert.Implements(t, (*request.ParserClientIP)(nil), &parser.ClientIP{})
}

func TestClientIP_RealIP(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name       string
		remoteAddr string
		want       netip.Addr
	}

	tests := []testCase{
		{
			name:       "IPv4 with port",
			remoteAddr: "192.0.2.1:1234",
			want:       netip.MustParseAddr("192.0.2.1"),
		},
		{
			name:       "IPv6 with port",
			remoteAddr: "[2001:db8::1]:80",
			want:       netip.MustParseAddr("2001:db8::1"),
		},
		{
			name:       "IPv4-mapped IPv6 is unmapped to IPv4",
			remoteAddr: "[::ffff:192.0.2.1]:80",
			want:       netip.MustParseAddr("192.0.2.1"),
		},
		{
			name:       "empty remote address returns invalid addr",
			remoteAddr: "",
			want:       netip.Addr{},
		},
		{
			name:       "malformed remote address returns invalid addr",
			remoteAddr: "not-an-ip",
			want:       netip.Addr{},
		},
	}

	p := parser.NewClientIP(mrlog.NopLogger())

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			r := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
			r.RemoteAddr = tc.remoteAddr

			got := p.RealIP(r)

			assert.Equal(t, tc.want, got)
		})
	}
}

func TestClientIP_DetailedIP(t *testing.T) {
	t.Parallel()

	const remoteAddr = "192.0.2.1:1234"

	realIP := netip.MustParseAddr("192.0.2.1")
	proxyIP := netip.MustParseAddr("203.0.113.5")

	type testCase struct {
		name       string
		remoteAddr string
		headers    map[string]string
		want       mrtype.DetailedIP
	}

	tests := []testCase{
		{
			name:       "no headers - only real IP",
			remoteAddr: remoteAddr,
			want:       mrtype.DetailedIP{Real: realIP},
		},
		{
			name:       "global proxy in X-Forwarded-For",
			remoteAddr: remoteAddr,
			headers:    map[string]string{"X-Forwarded-For": proxyIP.String()},
			want:       mrtype.DetailedIP{Real: realIP, Proxy: proxyIP},
		},
		{
			name:       "IPv6 global proxy in X-Forwarded-For",
			remoteAddr: remoteAddr,
			headers:    map[string]string{"X-Forwarded-For": "2001:db8::1"},
			want:       mrtype.DetailedIP{Real: realIP, Proxy: netip.MustParseAddr("2001:db8::1")},
		},
		{
			name:       "comma-separated XFF takes first global proxy",
			remoteAddr: remoteAddr,
			headers:    map[string]string{"X-Forwarded-For": "10.0.0.1, 203.0.113.5, 198.51.100.7"},
			want:       mrtype.DetailedIP{Real: realIP, Proxy: proxyIP},
		},
		{
			name:       "proxy equal to real is ignored",
			remoteAddr: remoteAddr,
			headers:    map[string]string{"X-Forwarded-For": realIP.String()},
			want:       mrtype.DetailedIP{Real: realIP},
		},
		{
			name:       "private proxy is ignored",
			remoteAddr: remoteAddr,
			headers:    map[string]string{"X-Forwarded-For": "10.0.0.1"},
			want:       mrtype.DetailedIP{Real: realIP},
		},
		{
			name:       "loopback proxy is ignored",
			remoteAddr: remoteAddr,
			headers:    map[string]string{"X-Forwarded-For": "127.0.0.1"},
			want:       mrtype.DetailedIP{Real: realIP},
		},
		{
			name:       "first global proxy is taken across default headers",
			remoteAddr: remoteAddr,
			headers: map[string]string{
				"X-Client-IP":     "10.0.0.1",
				"X-Forwarded-For": proxyIP.String(),
			},
			want: mrtype.DetailedIP{Real: realIP, Proxy: proxyIP},
		},
	}

	p := parser.NewClientIP(mrlog.NopLogger())

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			r := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
			r.RemoteAddr = tc.remoteAddr

			for key, value := range tc.headers {
				r.Header.Set(key, value)
			}

			assert.Equal(t, tc.want, p.DetailedIP(r))
		})
	}
}

func TestClientIP_DetailedIP_CustomHeaders(t *testing.T) {
	t.Parallel()

	realIP := netip.MustParseAddr("192.0.2.1")
	proxyIP := netip.MustParseAddr("203.0.113.5")

	p := parser.NewClientIP(mrlog.NopLogger(), "X-My-IP")

	r := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	r.RemoteAddr = "192.0.2.1:1234"
	// Заголовок из списка по умолчанию игнорируется, т.к. задан кастомный список.
	r.Header.Set("X-Forwarded-For", "198.51.100.7")
	r.Header.Set("X-My-IP", proxyIP.String())

	assert.Equal(t, mrtype.DetailedIP{Real: realIP, Proxy: proxyIP}, p.DetailedIP(r))
}
