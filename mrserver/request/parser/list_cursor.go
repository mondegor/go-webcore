//nolint:dupl
package parser

import (
	"net/http"

	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/mrtype"
	"github.com/mondegor/go-sysmess/mrtype/parse"
)

type (
	// ListCursor - парсер параметров для выборки части списка элементов.
	ListCursor struct {
		logger          mrlog.Logger
		paramNameCursor string
		paramNameLimit  string
		limitMax        int
		limitDefault    int
	}

	// ListCursorOptions - опции для создания ListCursor.
	ListCursorOptions struct {
		ParamNameCursor string
		ParamNameLimit  string
		LimitMax        int
		LimitDefault    int
	}
)

// NewListCursor - создаёт объект ListCursor.
func NewListCursor(logger mrlog.Logger, opts ListCursorOptions) *ListCursor {
	lc := ListCursor{
		logger:          logger,
		paramNameCursor: "cursor",
		paramNameLimit:  "limit",
		limitMax:        1000,
		limitDefault:    10,
	}

	if opts.ParamNameCursor != "" {
		lc.paramNameCursor = opts.ParamNameCursor
	}

	if opts.ParamNameLimit != "" {
		lc.paramNameLimit = opts.ParamNameLimit
	}

	if opts.LimitMax > 0 {
		lc.limitMax = opts.LimitMax
	}

	if opts.LimitDefault > 0 {
		lc.limitDefault = opts.LimitDefault
	}

	return &lc
}

// CursorParams - возвращает распарсенные параметры выборки части списка элементов.
func (c *ListCursor) CursorParams(r *http.Request) mrtype.CursorParams {
	value, err := parse.CursorParams(
		r.URL.Query().Get(c.paramNameCursor),
		r.URL.Query().Get(c.paramNameLimit),
	)

	if err != nil || value.Limit == 0 || value.Limit > c.limitMax {
		if err != nil {
			c.logger.Warn(
				r.Context(), "CursorParams",
				"value_key", c.paramNameCursor,
				"limit_key", c.paramNameLimit,
				"error", err,
			)
		}

		return mrtype.CursorParams{
			Limit: c.limitDefault,
		}
	}

	return value
}
