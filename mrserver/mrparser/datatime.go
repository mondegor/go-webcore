package mrparser

import (
	"net/http"
	"time"

	"github.com/mondegor/go-sysmess/mrlog"

	"github.com/mondegor/go-webcore/mrserver/mrreq"
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

// FilterDateTime - возвращает дата и время поступившие из внешнего запроса.
// Если ключ key не найден или возникнет ошибка, то возвращается nil значение.
func (p *DateTime) FilterDateTime(r *http.Request, key string) time.Time {
	value, err := mrreq.ParseDateTime(r.URL.Query(), key, false)
	if err != nil {
		p.logger.Warn(r.Context(), "FilterDateTime", "error", err)

		return time.Time{}
	}

	return value
}
