package mrparser

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/mondegor/go-webcore/mrctx"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrreq"
)

type (
	UUID struct {
		pathFunc mrserver.RequestParserParamFunc
	}
)

// Make sure the UUID conforms with the mrserver.RequestParserUUID interface
var _ mrserver.RequestParserUUID = (*UUID)(nil)

func NewUUID(pathFunc mrserver.RequestParserParamFunc) *UUID {
	return &UUID{
		pathFunc: pathFunc,
	}
}

func (p *UUID) PathParamUUID(r *http.Request, name string) uuid.UUID {
	value, err := uuid.Parse(p.pathFunc(r, name))

	if err != nil {
		mrctx.Logger(r.Context()).Warn(err)
		return uuid.Nil
	}

	return value
}

func (p *UUID) FilterUUID(r *http.Request, key string) uuid.UUID {
	value, err := mrreq.ParseUUID(r, key, true)

	if err != nil {
		mrctx.Logger(r.Context()).Warn(err)
		return uuid.Nil
	}

	return value
}
