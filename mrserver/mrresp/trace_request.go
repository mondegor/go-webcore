package mrresp

import (
	"time"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver"
)

// TraceRequest - comment func.
func TraceRequest(l mrlog.Logger, start time.Time, sr *mrserver.StatRequest, sw *mrserver.StatResponseWriter) {
	r := sr.Request()

	l.Trace().
		Str("method", r.Method).
		Str("url", r.RequestURI).
		Str("remoteAddr", r.RemoteAddr).
		Str("userAgent", r.UserAgent()).
		Int("status", sw.StatusCode()).
		Int("requestSize", sr.Bytes()).
		Int("size", sw.Bytes()).
		Int("elapsed_μs", int(time.Since(start).Microseconds())).
		Msg("incoming request")
}
