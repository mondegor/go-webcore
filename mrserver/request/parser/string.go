package parser

import (
	"net/http"

	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/mrtype/parse"
)

type (
	// String - парсер строки.
	String struct {
		pathFunc func(r *http.Request, key string) string
		logger   mrlog.Logger
	}
)

// NewString - создаёт объект String.
func NewString(
	pathFunc func(r *http.Request, key string) string,
	logger mrlog.Logger,
) *String {
	return &String{
		pathFunc: pathFunc,
		logger:   logger,
	}
}

// PathParamString - возвращает именованную строку, содержащуюся в URL пути.
// Если ключ name не найден, то возвращается пустое значение.
func (p *String) PathParamString(r *http.Request, name string) string {
	return p.pathFunc(r, name)
}

// RawParamString - возвращает именованную строку содержащуюся в URL пути.
// Если ключ name не найден, то возвращается nil значение.
func (p *String) RawParamString(r *http.Request, name string) *string {
	if !r.URL.Query().Has(name) {
		return nil
	}

	value := r.URL.Query().Get(name)

	return &value
}

// FilterString - возвращает строку поступившая из внешнего запроса.
// Если ключ key не найден или возникнет ошибка, то возвращается пустое значение.
func (p *String) FilterString(r *http.Request, key string) string {
	value, err := parse.String(r.URL.Query().Get(key), false)
	if err != nil {
		p.logger.Warn(r.Context(), "FilterString", "key", key, "error", err)

		return ""
	}

	return value
}

// FilterStringList - возвращает строку поступившую из внешнего запроса.
// Если ключ key не найден или возникнет ошибка, то возвращается nil значение.
func (p *String) FilterStringList(r *http.Request, key string) []string {
	items, err := parse.StringList(r.URL.Query().Get(key))
	if err != nil {
		p.logger.Warn(r.Context(), "FilterStringList", "key", key, "error", err)

		return nil
	}

	return items
}
