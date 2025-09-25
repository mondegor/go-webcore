package mrparser

import (
	"net/http"

	"github.com/mondegor/go-sysmess/mrlog"

	"github.com/mondegor/go-webcore/mrserver/mrreq"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	// ListPager - парсер параметров для выборки части списка элементов.
	ListPager struct {
		logger             mrlog.Logger
		paramNamePageIndex string
		paramNamePageSize  string
		pageSizeMax        uint64
		pageSizeDefault    uint64
	}

	// ListPagerOptions - опции для создания ListPager.
	ListPagerOptions struct {
		ParamNamePageIndex string
		ParamNamePageSize  string
		PageSizeMax        uint64
		PageSizeDefault    uint64
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

	return &lp
}

// PageParams - возвращает распарсенные параметры выборки части списка элементов.
func (p *ListPager) PageParams(r *http.Request) mrtype.PageParams {
	value, err := mrreq.ParsePageParams(
		r.URL.Query(),
		p.paramNamePageIndex,
		p.paramNamePageSize,
	)

	if err != nil || value.Size == 0 || value.Size > p.pageSizeMax {
		if err != nil {
			p.logger.Warn(r.Context(), "PageParams", "error", err)
		}

		return mrtype.PageParams{
			Size: p.pageSizeDefault,
		}
	}

	return value
}
