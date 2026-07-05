package parser

import (
	"net/http"
	"time"

	"github.com/mondegor/go-core/mrlog"
	"github.com/mondegor/go-core/mrtype/parse"
)

type (
	// DateTime - парсер даты и времени.
	DateTime struct {
		logger mrlog.Logger
	}
)

// NewDateTime - создаёт объект DateTime.
func NewDateTime(logger mrlog.Logger) *DateTime {
	return &DateTime{
		logger: logger,
	}
}

// FilterDateTime - возвращает дату и время поступившие из внешнего запроса.
// Если ключ key не найден или возникнет ошибка, то возвращается пустая структура.
func (p *DateTime) FilterDateTime(r *http.Request, key string) time.Time {
	value, err := parse.DateTime(r.URL.Query().Get(key), false)
	if err != nil {
		p.logger.Warn(r.Context(), "FilterDateTime", "key", key, "error", err)

		return time.Time{}
	}

	return value
}
