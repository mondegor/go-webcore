//nolint:dupl
package parser

import (
	"net/http"

	"github.com/mondegor/go-core/mrlog"
	"github.com/mondegor/go-core/mrtype"
	"github.com/mondegor/go-core/mrtype/parse"
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

	if lc.limitDefault > lc.limitMax {
		if opts.LimitDefault > 0 {
			mrlog.Warn(
				logger, "ListCursor: LimitDefault is greater than LimitMax, LimitMax is used",
				"limit_default", opts.LimitDefault,
				"limit_max", lc.limitMax,
			)
		}

		lc.limitDefault = lc.limitMax
	}

	return &lc
}

// CursorParams - возвращает распарсенные параметры выборки части списка элементов.
// Если limit не указан или равен 0, то берётся limit по умолчанию, а если больше максимального - максимальный,
// курсор при этом сохраняется; если параметры не удалось распарсить, то возвращается первая страница с limit по умолчанию.
func (c *ListCursor) CursorParams(r *http.Request) mrtype.CursorParams {
	query := r.URL.Query()

	value, err := parse.CursorParams(query.Get(c.paramNameCursor), query.Get(c.paramNameLimit))
	if err != nil {
		c.logger.Warn(
			r.Context(), "CursorParams",
			"value_key", c.paramNameCursor,
			"limit_key", c.paramNameLimit,
			"error", err,
		)

		return mrtype.CursorParams{
			Limit: c.limitDefault,
		}
	}

	if value.Limit < 1 {
		value.Limit = c.limitDefault
	} else if value.Limit > c.limitMax {
		value.Limit = c.limitMax
	}

	return value
}
