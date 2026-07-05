package stat

import (
	"net/http"
	"time"

	"github.com/mondegor/go-core/mrtrace"
)

type (
	// RequestTracer - выполняет трассировку HTTP-запросов для отладки и мониторинга.
	//
	// Записывает подробную информацию о каждом запросе через tracer:
	//  - HTTP-метод, URI, удалённый адрес, User-Agent;
	//  - Статус-код, размеры запроса и ответа, время выполнения;
	//  - Тело запроса и ответа (если присутствуют);
	RequestTracer struct {
		tracer mrtrace.Tracer
	}
)

// NewRequestTracer - создаёт трассировщик HTTP-запросов.
func NewRequestTracer(tracer mrtrace.Tracer) *RequestTracer {
	return &RequestTracer{
		tracer: tracer,
	}
}

// Enabled - сообщает, включена ли трассировка запросов.
func (rs *RequestTracer) Enabled() bool {
	return rs.tracer != mrtrace.NopTracer()
}

// Emit - выполняет трассировку HTTP-запроса.
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
