package parser

import (
	"net/http"

	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/mrtype/parse"
)

type (
	// Bool - парсер bool значений.
	Bool struct {
		logger mrlog.Logger
	}
)

// NewBool - создаёт объект Bool.
func NewBool(
	logger mrlog.Logger,
) *Bool {
	return &Bool{
		logger: logger,
	}
}

// FilterNullableBool - возвращает bool значение поступившее из внешнего запроса.
// Если ключ key не найден или возникнет ошибка, то возвращается nil значение.
func (p *Bool) FilterNullableBool(r *http.Request, key string) *bool {
	value, err := parse.NullableBool(r.URL.Query().Get(key))
	if err != nil {
		p.logger.Warn(r.Context(), "FilterNullableBool", "key", key, "error", err)

		return nil
	}

	return value
}
