package mrparser

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrreq"
)

type (
	// UUID - парсер UUID.
	UUID struct {
		pathFunc mrserver.RequestParserParamFunc
	}
)

// NewUUID - создаёт объект UUID.
func NewUUID(pathFunc mrserver.RequestParserParamFunc) *UUID {
	return &UUID{
		pathFunc: pathFunc,
	}
}

// PathParamUUID - возвращает именованный UUID содержащийся в URL пути.
func (p *UUID) PathParamUUID(r *http.Request, name string) uuid.UUID {
	value, err := uuid.Parse(p.pathFunc(r, name))
	if err != nil {
		mrlog.Ctx(r.Context()).Warn().Err(err).Send()

		return uuid.Nil
	}

	return value
}

// FilterUUID - возвращает UUID поступивший из внешнего запроса.
// Если ключ key не найден или возникнет ошибка, то возвращается uuid.Nil.
func (p *UUID) FilterUUID(r *http.Request, key string) uuid.UUID {
	value, err := mrreq.ParseUUID(r.URL.Query(), key, false)
	if err != nil {
		mrlog.Ctx(r.Context()).Warn().Err(err).Send()

		return uuid.Nil
	}

	return value
}
