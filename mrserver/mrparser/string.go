package mrparser

import (
	"net/http"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrreq"
)

type (
	// String - comment struct.
	String struct {
		pathFunc mrserver.RequestParserParamFunc
	}
)

// Make sure the String conforms with the mrserver.RequestParserString interface.
var _ mrserver.RequestParserString = (*String)(nil)

// NewString - создаёт объект String.
func NewString(pathFunc mrserver.RequestParserParamFunc) *String {
	return &String{
		pathFunc: pathFunc,
	}
}

// PathParamString - comment method.
func (p *String) PathParamString(r *http.Request, name string) string {
	return p.pathFunc(r, name)
}

// RawParamString - returns nil if the param not found.
func (p *String) RawParamString(r *http.Request, key string) *string {
	if !r.URL.Query().Has(key) {
		return nil
	}

	value := r.URL.Query().Get(key)

	return &value
}

// FilterString - comment method.
func (p *String) FilterString(r *http.Request, key string) string {
	value, err := mrreq.ParseStr(r, key, false)
	if err != nil {
		mrlog.Ctx(r.Context()).Warn().Err(err).Send()

		return ""
	}

	return value
}
