package mrserver

import (
	"time"

	"github.com/mondegor/go-webcore/mrlog"
)

// TraceRequest - функция трассировки http запроса.
func TraceRequest(l mrlog.Logger, start time.Time, sr *StatRequestReader, sw *StatResponseWriter) {
	r := sr.Request()

	if l.Level() == mrlog.TraceLevel {
		l.Trace().
			Str("method", r.Method).
			Str("uri", r.RequestURI).
			Str("remoteAddr", r.RemoteAddr).
			Str("userAgent", r.UserAgent()).
			Int("status", sw.StatusCode()).
			Int("requestSize", sr.Size()).
			Int("size", sw.Size()).
			Int64("elapsed_µs", time.Since(start).Microseconds()).
			Msg("response")

		if sr.HasContent() {
			l.Trace().Bytes("body", sr.Content()).Msg("request")
		}

		l.Trace().Bytes("body", sw.Content()).Msg("response")

		return
	}

	l.Info().
		Int("status", sw.StatusCode()).
		Int("requestSize", sr.Size()).
		Int("size", sw.Size()).
		Int64("elapsed_µs", time.Since(start).Microseconds()).
		Msg("response")
}
