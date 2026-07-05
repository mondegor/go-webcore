package parser

import (
	"net/http"

	"github.com/mondegor/go-core/mrlog"
	"github.com/mondegor/go-core/mrtype/parse"
)

type (
	// Bool - парсер bool значений.
	Bool struct {
		logger mrlog.Logger
	}
)

// NewBool - создаёт парсер логических значений.
func NewBool(
	logger mrlog.Logger,
) *Bool {
	return &Bool{
		logger: logger,
	}
}

// FilterNullableBool - извлекает nullable bool значение из query-параметра.
// Поддерживает значения: "true", "false", "1", "0", "yes", "no" (регистронезависимо).
func (p *Bool) FilterNullableBool(r *http.Request, key string) *bool {
	value, err := parse.NullableBool(r.URL.Query().Get(key))
	if err != nil {
		p.logger.Warn(r.Context(), "FilterNullableBool", "key", key, "error", err)

		return nil
	}

	return value
}
