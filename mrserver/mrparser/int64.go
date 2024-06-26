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
	// Int64 - comment struct.
	Int64 struct {
		pathFunc mrserver.RequestParserParamFunc
	}
)

// Make sure the Int64 conforms with the mrserver.RequestParserInt64 interface.
var _ mrserver.RequestParserInt64 = (*Int64)(nil)

// NewInt64 - создаёт объект Int64.
func NewInt64(pathFunc mrserver.RequestParserParamFunc) *Int64 {
	return &Int64{
		pathFunc: pathFunc,
	}
}

// PathParamInt64 - comment method.
func (p *Int64) PathParamInt64(r *http.Request, name string) int64 {
	value, err := strconv.ParseInt(p.pathFunc(r, name), 10, 64)
	if err != nil {
		mrlog.Ctx(r.Context()).Warn().Err(err).Send()

		return 0
	}

	return value
}

// FilterInt64 - comment method.
func (p *Int64) FilterInt64(r *http.Request, key string) int64 {
	value, err := mrreq.ParseInt64(r, key, false)
	if err != nil {
		mrlog.Ctx(r.Context()).Warn().Err(err).Send()

		return 0
	}

	return value
}

// FilterRangeInt64 - comment method.
func (p *Int64) FilterRangeInt64(r *http.Request, key string) mrtype.RangeInt64 {
	value, err := mrreq.ParseRangeInt64(r, key)
	if err != nil {
		mrlog.Ctx(r.Context()).Warn().Err(err).Send()

		return mrtype.RangeInt64{}
	}

	return value
}

// FilterInt64List - comment method.
func (p *Int64) FilterInt64List(r *http.Request, key string) []int64 {
	items, err := mrreq.ParseInt64List(r, key)
	if err != nil {
		mrlog.Ctx(r.Context()).Warn().Err(err).Send()

		return nil
	}

	return items
}
