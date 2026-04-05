package stat

import (
	"net/http"
	"time"

	"github.com/mondegor/go-sysmess/mrlog"
)

type (
	// RequestLogger - записывает информацию о HTTP-запросах в лог.
	RequestLogger struct {
		logger  mrlog.Logger
		enabled bool
	}
)

// NewRequestLogger - создаёт объект RequestLogger.
func NewRequestLogger(logger mrlog.Logger) *RequestLogger {
	return &RequestLogger{
		logger:  logger,
		enabled: mrlog.InfoEnabled(logger),
	}
}

// Enabled - сообщает, включено ли логирование запросов.
func (rs *RequestLogger) Enabled() bool {
	return rs.enabled
}

// Emit - записывает информацию о HTTP запросе в лог.
func (rs *RequestLogger) Emit(r *http.Request, _ []byte, size int, _ []byte, responseSize int, duration time.Duration, status int) {
	rs.logger.Info(
		r.Context(),
		"RESPONSE",
		"status", status,
		"requestSize", size,
		"size", responseSize,
		"elapsed_µs", duration.Microseconds(),
	)
}
