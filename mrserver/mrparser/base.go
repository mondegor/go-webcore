package mrparser

import (
	"net/http"
	"time"

	"github.com/mondegor/go-webcore/mrctx"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrreq"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	Base struct {
		path mrserver.RequestParserPath
	}
)

// Make sure the Base conforms with the mrserver.RequestParser interface
var _ mrserver.RequestParser = (*Base)(nil)

func NewBase(path mrserver.RequestParserPath) *Base {
	return &Base{
		path: path,
	}
}

func (p *Base) PathParamString(r *http.Request, name string) string {
	return p.path.PathParam(r, name)
}

func (p *Base) PathParamInt64(r *http.Request, name string) int64 {
	value, err := mrreq.ParseInt64(r, p.path.PathParam(r, name), false)

	if err != nil {
		mrctx.Logger(r.Context()).Warn(err)
		return 0
	}

	return value
}

// RawQueryParam - returns nil if the param not found
func (p *Base) RawQueryParam(r *http.Request, key string) *string {
	if !r.URL.Query().Has(key) {
		return nil
	}

	value := r.URL.Query().Get(key)

	return &value
}

func (p *Base) FilterString(r *http.Request, key string) string {
	value, err := mrreq.ParseStr(r, key, false)

	if err != nil {
		mrctx.Logger(r.Context()).Warn(err)
		return ""
	}

	return value
}

func (p *Base) FilterNullableBool(r *http.Request, key string) *bool {
	value, err := mrreq.ParseNullableBool(r, key)

	if err != nil {
		mrctx.Logger(r.Context()).Warn(err)
		return nil
	}

	return value
}

func (p *Base) FilterInt64(r *http.Request, key string) int64 {
	value, err := mrreq.ParseInt64(r, key, false)

	if err != nil {
		mrctx.Logger(r.Context()).Warn(err)
		return 0
	}

	return value
}

func (p *Base) FilterRangeInt64(r *http.Request, key string) mrtype.RangeInt64 {
	value, err := mrreq.ParseRangeInt64(r, key)

	if err != nil {
		mrctx.Logger(r.Context()).Warn(err)
		return mrtype.RangeInt64{}
	}

	return value
}

func (p *Base) FilterInt64List(r *http.Request, key string) []int64 {
	items, err := mrreq.ParseInt64List(r, key)

	if err != nil {
		mrctx.Logger(r.Context()).Warn(err)
		return []int64{}
	}

	return items
}

func (p *Base) FilterDateTime(r *http.Request, key string) time.Time {
	value, err := mrreq.ParseDateTime(r, key, false)

	if err != nil {
		mrctx.Logger(r.Context()).Warn(err)
		return time.Time{}
	}

	return value
}
