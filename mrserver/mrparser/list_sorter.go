package mrparser

import (
	"fmt"
	"net/http"

	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/mrtype"

	"github.com/mondegor/go-webcore/mrserver/mrreq"
)

type (
	// ListSorter - парсер параметров для сортировки списка элементов.
	ListSorter struct {
		logger                 mrlog.Logger
		paramNameSortField     string
		paramNameSortDirection string
	}

	// ListSorterOptions - опции для создания ListSorter.
	ListSorterOptions struct {
		ParamNameSortField     string
		ParamNameSortDirection string
	}
)

// NewListSorter - создаёт объект ListSorter.
func NewListSorter(logger mrlog.Logger, opts ListSorterOptions) *ListSorter {
	ls := ListSorter{
		logger:                 logger,
		paramNameSortField:     "sortField",
		paramNameSortDirection: "sortDirection",
	}

	if opts.ParamNameSortField != "" {
		ls.paramNameSortField = opts.ParamNameSortField
	}

	if opts.ParamNameSortDirection != "" {
		ls.paramNameSortDirection = opts.ParamNameSortDirection
	}

	return &ls
}

// SortParams - возвращает распарсенные параметры сортировки списка элементов.
func (p *ListSorter) SortParams(r *http.Request, sorter mrtype.ListSorter) mrtype.SortParams {
	value, err := mrreq.ParseSortParams(
		r.URL.Query(),
		p.paramNameSortField,
		p.paramNameSortDirection,
	)
	if err != nil {
		p.logger.Warn(r.Context(), "SortParams", "error", err)
	}

	if value.FieldName == "" {
		return sorter.DefaultSort()
	}

	if !sorter.CheckField(value.FieldName) {
		p.logger.Warn(r.Context(), fmt.Sprintf("sort field '%s' is not registered", value.FieldName), "error", err)

		return sorter.DefaultSort()
	}

	return value
}
