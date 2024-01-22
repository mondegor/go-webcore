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
		path mrserver.RequestParserPath
	}
)

// Make sure the UUID conforms with the mrserver.RequestParserUUID interface
var _ mrserver.RequestParserUUID = (*UUID)(nil)

func NewUUID(path mrserver.RequestParserPath) *UUID {
	return &UUID{
		path: path,
	}
}

func (p *UUID) PathParamUUID(r *http.Request, name string) uuid.UUID {
	value, err := uuid.Parse(p.path.PathParam(r, name))

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
