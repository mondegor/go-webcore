package mrparser

import (
	"net/http"

	"github.com/mondegor/go-sysmess/mrlog"

	"github.com/mondegor/go-webcore/mrserver/mrreq"
)

type (
	// Bool - парсер bool значений.
	Bool struct {
		logger mrlog.Logger
	}
)

// NewBool - создаёт объект Bool.
func NewBool(logger mrlog.Logger) *Bool {
	return &Bool{
		logger: logger,
	}
}

// FilterNullableBool - возвращает bool значение поступившее из внешнего запроса.
// Если ключ key не найден или возникнет ошибка, то возвращается nil значение.
func (p *Bool) FilterNullableBool(r *http.Request, key string) *bool {
	value, err := mrreq.ParseNullableBool(r.URL.Query(), key)
	if err != nil {
		p.logger.Warn(r.Context(), "FilterNullableBool", "error", err)

		return nil
	}

	return value
}
