package stat

import (
	"net/http"
	"time"

	"github.com/mondegor/go-sysmess/mrlog"
)

type (
	// RequestLogger - comment struct.
	RequestLogger struct {
		logger mrlog.Logger
	}
)

// NewRequestLogger - создаёт объект RequestLogger.
func NewRequestLogger(logger mrlog.Logger) *RequestLogger {
	return &RequestLogger{
		logger: logger,
	}
}

// Enabled - comment method.
func (rs *RequestLogger) Enabled() bool {
	return rs.logger.Enabled(mrlog.LevelInfo)
}

// Emit - comment method.
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
