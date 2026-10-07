package parser_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mondegor/go-core/mrlog"
	"github.com/mondegor/go-core/mrtype"
	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrserver/request"
	"github.com/mondegor/go-webcore/mrserver/request/parser"
)

// Make sure the ListPager conforms with the request.ParserListPager interface.
func TestListPagerImplementsRequestParserListPager(t *testing.T) {
	t.Parallel()

	assert.Implements(t, (*request.ParserListPager)(nil), &parser.ListPager{})
}

func TestListPager_PageParams(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name         string
		opts         parser.ListPagerOptions
		query        string
		want         mrtype.PageParams
		wantCtorWarn bool
		wantWarn     bool
	}

	opts := parser.ListPagerOptions{PageSizeMax: 100, PageSizeDefault: 10}

	tests := []testCase{
		{name: "absent", opts: opts, query: "", want: mrtype.PageParams{Index: 0, Size: 10}},
		{name: "index without size", opts: opts, query: "pageIndex=3", want: mrtype.PageParams{Index: 3, Size: 10}},
		{name: "empty size", opts: opts, query: "pageIndex=3&pageSize=", want: mrtype.PageParams{Index: 3, Size: 10}},
		{name: "blank size", opts: opts, query: "pageIndex=3&pageSize=%20", want: mrtype.PageParams{Index: 3, Size: 10}},
		{name: "zero size", opts: opts, query: "pageIndex=3&pageSize=0", want: mrtype.PageParams{Index: 3, Size: 10}},
		{name: "index with size", opts: opts, query: "pageIndex=3&pageSize=20", want: mrtype.PageParams{Index: 3, Size: 20}},
		{name: "size equal max", opts: opts, query: "pageIndex=3&pageSize=100", want: mrtype.PageParams{Index: 3, Size: 100}},
		{name: "size over max", opts: opts, query: "pageIndex=3&pageSize=101", want: mrtype.PageParams{Index: 3, Size: 100}},
		{name: "size without index", opts: opts, query: "pageSize=20", want: mrtype.PageParams{Index: 0, Size: 20}},
		{
			name:     "incorrect size",
			opts:     opts,
			query:    "pageIndex=3&pageSize=x",
			want:     mrtype.PageParams{Index: 0, Size: 10},
			wantWarn: true,
		},
		{
			name:     "negative size",
			opts:     opts,
			query:    "pageIndex=3&pageSize=-1",
			want:     mrtype.PageParams{Index: 0, Size: 10},
			wantWarn: true,
		},
		{
			name:     "size overflow",
			opts:     opts,
			query:    "pageIndex=3&pageSize=99999999999999999999",
			want:     mrtype.PageParams{Index: 0, Size: 10},
			wantWarn: true,
		},
		{
			name:     "incorrect index",
			opts:     opts,
			query:    "pageIndex=x&pageSize=20",
			want:     mrtype.PageParams{Index: 0, Size: 10},
			wantWarn: true,
		},
		{
			name:     "index overflow",
			opts:     opts,
			query:    "pageIndex=99999999999999999999&pageSize=20",
			want:     mrtype.PageParams{Index: 0, Size: 10},
			wantWarn: true,
		},
		{name: "default options without size", query: "pageIndex=3", want: mrtype.PageParams{Index: 3, Size: 10}},
		{name: "default options size equal max", query: "pageIndex=3&pageSize=1000", want: mrtype.PageParams{Index: 3, Size: 1000}},
		{name: "default options size over max", query: "pageIndex=3&pageSize=1001", want: mrtype.PageParams{Index: 3, Size: 1000}},
		{
			name:  "custom param names",
			opts:  parser.ListPagerOptions{ParamNamePageIndex: "i", ParamNamePageSize: "s"},
			query: "i=3&s=20&pageIndex=4&pageSize=30",
			want:  mrtype.PageParams{Index: 3, Size: 20},
		},
		{
			name:         "default size over max",
			opts:         parser.ListPagerOptions{PageSizeMax: 5, PageSizeDefault: 10},
			query:        "pageIndex=3",
			want:         mrtype.PageParams{Index: 3, Size: 5},
			wantCtorWarn: true,
		},
		{
			name:  "max below builtin default",
			opts:  parser.ListPagerOptions{PageSizeMax: 5},
			query: "pageIndex=3",
			want:  mrtype.PageParams{Index: 3, Size: 5},
		},
		{
			name:         "default size over max with incorrect size",
			opts:         parser.ListPagerOptions{PageSizeMax: 5, PageSizeDefault: 10},
			query:        "pageIndex=3&pageSize=x",
			want:         mrtype.PageParams{Index: 0, Size: 5},
			wantCtorWarn: true,
			wantWarn:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			logger := &warnSpyLogger{Logger: mrlog.NopLogger()}
			r := httptest.NewRequest(http.MethodGet, "/?"+tt.query, http.NoBody)

			lp := parser.NewListPager(logger, tt.opts)

			assert.Equal(t, tt.wantCtorWarn, logger.warnCount > 0)

			logger.warnCount = 0
			got := lp.PageParams(r)

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantWarn, logger.warnCount > 0)
		})
	}
}
