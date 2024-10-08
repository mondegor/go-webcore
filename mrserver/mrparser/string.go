package mrparser

import (
	"net/http"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrreq"
)

type (
	// String - парсер строки.
	String struct {
		pathFunc mrserver.RequestParserParamFunc
	}
)

// NewString - создаёт объект String.
func NewString(pathFunc mrserver.RequestParserParamFunc) *String {
	return &String{
		pathFunc: pathFunc,
	}
}

// PathParamString - возвращает именованная строка содержащаяся в URL пути.
// Если ключ name не найден, то возвращается пустое значение.
func (p *String) PathParamString(r *http.Request, name string) string {
	return p.pathFunc(r, name)
}

// RawParamString - возвращает именованная строка содержащаяся в URL пути.
// Если ключ name не найден, то возвращается nil значение.
func (p *String) RawParamString(r *http.Request, name string) *string {
	if !r.URL.Query().Has(name) {
		return nil
	}

	value := r.URL.Query().Get(name)

	return &value
}

// FilterString - возвращает строка поступившая из внешнего запроса.
// Если ключ key не найден или возникнет ошибка, то возвращается пустое значение.
func (p *String) FilterString(r *http.Request, key string) string {
	value, err := mrreq.ParseStr(r.URL.Query(), key, false)
	if err != nil {
		mrlog.Ctx(r.Context()).Warn().Err(err).Send()

		return ""
	}

	return value
}
