package mrparser

import (
	"net/http"
	"strconv"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrreq"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	// KeyInt32 - comment struct.
	KeyInt32 struct {
		pathFunc mrserver.RequestParserParamFunc
	}
)

// Make sure the KeyInt32 conforms with the mrserver.RequestParserKeyInt32 interface.
var _ mrserver.RequestParserKeyInt32 = (*KeyInt32)(nil)

// NewKeyInt32 - создаёт объект KeyInt32.
func NewKeyInt32(pathFunc mrserver.RequestParserParamFunc) *KeyInt32 {
	return &KeyInt32{
		pathFunc: pathFunc,
	}
}

// PathKeyInt32 - comment method.
func (p *KeyInt32) PathKeyInt32(r *http.Request, name string) mrtype.KeyInt32 {
	value, err := strconv.ParseInt(p.pathFunc(r, name), 10, 32)
	if err != nil {
		mrlog.Ctx(r.Context()).Warn().Err(err).Send()

		return 0
	}

	return mrtype.KeyInt32(value)
}

// FilterKeyInt32 - comment method.
func (p *KeyInt32) FilterKeyInt32(r *http.Request, key string) mrtype.KeyInt32 {
	value, err := mrreq.ParseInt64(r, key, false)
	if err != nil {
		mrlog.Ctx(r.Context()).Warn().Err(err).Send()

		return 0
	}

	return mrtype.KeyInt32(value)
}

// FilterKeyInt32List - comment method.
func (p *KeyInt32) FilterKeyInt32List(r *http.Request, key string) []mrtype.KeyInt32 {
	value, err := mrreq.ParseInt64List(r, key)
	if err != nil {
		mrlog.Ctx(r.Context()).Warn().Err(err).Send()

		return nil
	}

	items := make([]mrtype.KeyInt32, len(value))

	for i := range value {
		items[i] = mrtype.KeyInt32(value[i])
	}

	return items
}
