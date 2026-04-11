package stat

import (
	"net/http"
	"time"

	"github.com/mondegor/go-sysmess/mrlog"
)

type (
	// RequestLogger - логгер HTTP-запросов.
	// Записывает информацию о каждом запросе в лог на уровне INFO.
	// Логирует: статус-код, размер запроса и ответа, время выполнения.
	RequestLogger struct {
		// logger - логгер для записи информации о запросах.
		logger mrlog.Logger
		// enabled - флаг включения логирования (проверяется при создании).
		enabled bool
	}
)

// NewRequestLogger - создаёт логгер HTTP-запросов.
// Автоматически проверяет включён ли уровень INFO для оптимизации.
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

// Emit - записывает информацию о HTTP-запросе в лог на уровне INFO.
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
