package mrparser

import (
	"net/http"
	"time"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver/mrreq"
)

type (
	// DateTime - парсер даты и времени.
	DateTime struct{}
)

// NewDateTime - создаёт объект DateTime.
func NewDateTime() *DateTime {
	return &DateTime{}
}

// FilterDateTime - возвращает дата и время поступившие из внешнего запроса.
// Если ключ key не найден или возникнет ошибка, то возвращается nil значение.
func (p *DateTime) FilterDateTime(r *http.Request, key string) time.Time {
	value, err := mrreq.ParseDateTime(r.URL.Query(), key, false)
	if err != nil {
		mrlog.Ctx(r.Context()).Warn().Err(err).Send()

		return time.Time{}
	}

	return value
}
