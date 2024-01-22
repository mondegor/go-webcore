package mrparser

import (
	"net/http"
	"strconv"

	"github.com/mondegor/go-webcore/mrctx"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrreq"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	KeyInt32 struct {
		path mrserver.RequestParserPath
	}
)

// Make sure the KeyInt32 conforms with the mrserver.RequestParserKeyInt32 interface
var _ mrserver.RequestParserKeyInt32 = (*KeyInt32)(nil)

func NewKeyInt32(path mrserver.RequestParserPath) *KeyInt32 {
	return &KeyInt32{
		path: path,
	}
}

func (p *KeyInt32) PathKeyInt32(r *http.Request, name string) mrtype.KeyInt32 {
	value, err := strconv.ParseInt(p.path.PathParam(r, name), 10, 32)

	if err != nil {
		mrctx.Logger(r.Context()).Warn(err)
		return 0
	}

	return mrtype.KeyInt32(value)
}

func (p *KeyInt32) FilterKeyInt32(r *http.Request, key string) mrtype.KeyInt32 {
	value, err := mrreq.ParseInt64(r, key, false)

	if err != nil {
		mrctx.Logger(r.Context()).Warn(err)
		return 0
	}

	return mrtype.KeyInt32(value)
}

func (p *KeyInt32) FilterKeyInt32List(r *http.Request, key string) []mrtype.KeyInt32 {
	value, err := mrreq.ParseInt64List(r, key)

	if err != nil {
		mrctx.Logger(r.Context()).Warn(err)
		return []mrtype.KeyInt32{}
	}

	items := make([]mrtype.KeyInt32, len(value))

	for i := range value {
		items[i] = mrtype.KeyInt32(value[i])
	}

	return items
}
