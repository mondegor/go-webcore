package stat

import (
	"net/http"
	"time"

	"github.com/mondegor/go-webcore/mrsender"
)

type (
	// RequestTracer - comment struct.
	RequestTracer struct {
		tracer mrsender.Tracer
	}
)

// NewRequestTracer - создаёт объект RequestTracer.
func NewRequestTracer(tracer mrsender.Tracer) *RequestTracer {
	return &RequestTracer{
		tracer: tracer,
	}
}

// Enabled - comment method.
func (rs *RequestTracer) Enabled() bool {
	return rs.tracer.Enabled()
}

// Emit - функция трассировки http запроса.
func (rs *RequestTracer) Emit(r *http.Request, body []byte, size int, responseBody []byte, responseSize int, duration time.Duration, status int) {
	ctx := r.Context()

	rs.tracer.Trace(
		ctx,
		"source", "REQUEST",
		"method", r.Method,
		"uri", r.RequestURI,
		"remoteAddr", r.RemoteAddr,
		"userAgent", r.UserAgent(),
		"status", status,
		"requestSize", size,
		"size", responseSize,
		"elapsed_µs", duration.Microseconds(),
	)

	if len(body) > 0 {
		rs.tracer.Trace(ctx, "source", "REQUEST", "body", body)
	}

	rs.tracer.Trace(ctx, "source", "RESPONSE", "body", responseBody)
}
