//nolint:dupl
package parser

import (
	"net/http"

	"github.com/mondegor/go-core/mrlog"
	"github.com/mondegor/go-core/mrtype"
	"github.com/mondegor/go-core/mrtype/parse"
)

type (
	// ListPager - парсер параметров для выборки части списка элементов.
	ListPager struct {
		logger             mrlog.Logger
		paramNamePageIndex string
		paramNamePageSize  string
		pageSizeMax        int
		pageSizeDefault    int
	}

	// ListPagerOptions - опции для создания ListPager.
	ListPagerOptions struct {
		ParamNamePageIndex string
		ParamNamePageSize  string
		PageSizeMax        int
		PageSizeDefault    int
	}
)

// NewListPager - создаёт объект ListPager.
func NewListPager(logger mrlog.Logger, opts ListPagerOptions) *ListPager {
	lp := ListPager{
		logger:             logger,
		paramNamePageIndex: "pageIndex",
		paramNamePageSize:  "pageSize",
		pageSizeMax:        1000,
		pageSizeDefault:    10,
	}

	if opts.ParamNamePageIndex != "" {
		lp.paramNamePageIndex = opts.ParamNamePageIndex
	}

	if opts.ParamNamePageSize != "" {
		lp.paramNamePageSize = opts.ParamNamePageSize
	}

	if opts.PageSizeMax > 0 {
		lp.pageSizeMax = opts.PageSizeMax
	}

	if opts.PageSizeDefault > 0 {
		lp.pageSizeDefault = opts.PageSizeDefault
	}

	if lp.pageSizeDefault > lp.pageSizeMax {
		if opts.PageSizeDefault > 0 {
			mrlog.Warn(
				logger, "ListPager: PageSizeDefault is greater than PageSizeMax, PageSizeMax is used",
				"page_size_default", opts.PageSizeDefault,
				"page_size_max", lp.pageSizeMax,
			)
		}

		lp.pageSizeDefault = lp.pageSizeMax
	}

	return &lp
}

// PageParams - возвращает распарсенные параметры выборки части списка элементов.
// Если размер страницы не указан или равен 0, то берётся размер по умолчанию, а если больше максимального - максимальный,
// индекс страницы при этом сохраняется; если параметры не удалось распарсить, то возвращается первая страница с размером по умолчанию.
func (p *ListPager) PageParams(r *http.Request) mrtype.PageParams {
	query := r.URL.Query()

	value, err := parse.PageParams(query.Get(p.paramNamePageIndex), query.Get(p.paramNamePageSize))
	if err != nil {
		p.logger.Warn(
			r.Context(), "PageParams",
			"index_key", p.paramNamePageIndex,
			"size_key", p.paramNamePageSize,
			"error", err,
		)

		return mrtype.PageParams{
			Size: p.pageSizeDefault,
		}
	}

	if value.Size < 1 {
		value.Size = p.pageSizeDefault
	} else if value.Size > p.pageSizeMax {
		value.Size = p.pageSizeMax
	}

	return value
}
