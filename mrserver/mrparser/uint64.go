package mrparser

import (
	"net/http"
	"strconv"

	"github.com/mondegor/go-sysmess/mrlog"

	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrreq"
)

type (
	// Uint64 - парсер uint64 числа.
	Uint64 struct {
		pathFunc mrserver.RequestParserParamFunc
		logger   mrlog.Logger
	}
)

// NewUint64 - создаёт объект Uint64.
func NewUint64(pathFunc mrserver.RequestParserParamFunc, logger mrlog.Logger) *Uint64 {
	return &Uint64{
		pathFunc: pathFunc,
		logger:   logger,
	}
}

// PathParamUint64 - возвращает именованное uint64 число содержащееся в URL пути.
func (p *Uint64) PathParamUint64(r *http.Request, name string) uint64 {
	value, err := strconv.ParseUint(p.pathFunc(r, name), 10, 64)
	if err != nil {
		p.logger.Warn(r.Context(), "PathParamUint64", "error", err)

		return 0
	}

	return value
}

// FilterUint64 - возвращает uint64 число поступившее из внешнего запроса.
// Если ключ key не найден или возникнет ошибка, то возвращается нулевое значение.
func (p *Uint64) FilterUint64(r *http.Request, key string) uint64 {
	value, err := mrreq.ParseUint64(r.URL.Query(), key, false)
	if err != nil {
		p.logger.Warn(r.Context(), "FilterUint64", "error", err)

		return 0
	}

	return value
}

// FilterUint64List - возвращает массив uint64 чисел поступивший из внешнего запроса.
// Если ключ key не найден или возникнет ошибка, то возвращается пустой массив.
func (p *Uint64) FilterUint64List(r *http.Request, key string) []uint64 {
	items, err := mrreq.ParseUint64List(r.URL.Query(), key)
	if err != nil {
		p.logger.Warn(r.Context(), "FilterUint64List", "error", err)

		return nil
	}

	return items
}
