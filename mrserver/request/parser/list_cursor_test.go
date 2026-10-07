package parser_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mondegor/go-core/mrlog"
	"github.com/mondegor/go-core/mrtype"
	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrserver/request"
	"github.com/mondegor/go-webcore/mrserver/request/parser"
)

type (
	// warnSpyLogger - логгер, фиксирующий вызовы уровня Warn.
	warnSpyLogger struct {
		mrlog.Logger
		warnCount int
	}
)

// Warn - фиксирует факт логирования предупреждения.
func (l *warnSpyLogger) Warn(_ context.Context, _ string, _ ...any) {
	l.warnCount++
}

// Make sure the ListCursor conforms with the request.ParserListCursor interface.
func TestListCursorImplementsRequestParserListCursor(t *testing.T) {
	t.Parallel()

	assert.Implements(t, (*request.ParserListCursor)(nil), &parser.ListCursor{})
}

func TestListCursor_CursorParams(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name         string
		opts         parser.ListCursorOptions
		query        string
		want         mrtype.CursorParams
		wantCtorWarn bool
		wantWarn     bool
	}

	opts := parser.ListCursorOptions{LimitMax: 100, LimitDefault: 10}

	tests := []testCase{
		{name: "absent", opts: opts, query: "", want: mrtype.CursorParams{Value: "", Limit: 10}},
		{name: "cursor without limit", opts: opts, query: "cursor=abc", want: mrtype.CursorParams{Value: "abc", Limit: 10}},
		{name: "empty limit", opts: opts, query: "cursor=abc&limit=", want: mrtype.CursorParams{Value: "abc", Limit: 10}},
		{name: "blank limit", opts: opts, query: "cursor=abc&limit=%20", want: mrtype.CursorParams{Value: "abc", Limit: 10}},
		{name: "zero limit", opts: opts, query: "cursor=abc&limit=0", want: mrtype.CursorParams{Value: "abc", Limit: 10}},
		{name: "cursor with limit", opts: opts, query: "cursor=abc&limit=20", want: mrtype.CursorParams{Value: "abc", Limit: 20}},
		{name: "limit equal max", opts: opts, query: "cursor=abc&limit=100", want: mrtype.CursorParams{Value: "abc", Limit: 100}},
		{name: "limit over max", opts: opts, query: "cursor=abc&limit=101", want: mrtype.CursorParams{Value: "abc", Limit: 100}},
		{name: "blank cursor", opts: opts, query: "cursor=%20&limit=20", want: mrtype.CursorParams{Value: "", Limit: 20}},
		{
			name:     "incorrect limit",
			opts:     opts,
			query:    "cursor=abc&limit=x",
			want:     mrtype.CursorParams{Value: "", Limit: 10},
			wantWarn: true,
		},
		{
			name:     "negative limit",
			opts:     opts,
			query:    "cursor=abc&limit=-1",
			want:     mrtype.CursorParams{Value: "", Limit: 10},
			wantWarn: true,
		},
		{
			name:     "limit overflow",
			opts:     opts,
			query:    "cursor=abc&limit=99999999999999999999",
			want:     mrtype.CursorParams{Value: "", Limit: 10},
			wantWarn: true,
		},
		{
			name:     "too long cursor",
			opts:     opts,
			query:    "cursor=" + strings.Repeat("a", 257) + "&limit=20",
			want:     mrtype.CursorParams{Value: "", Limit: 10},
			wantWarn: true,
		},
		{name: "default options without limit", query: "cursor=abc", want: mrtype.CursorParams{Value: "abc", Limit: 10}},
		{name: "default options limit equal max", query: "cursor=abc&limit=1000", want: mrtype.CursorParams{Value: "abc", Limit: 1000}},
		{name: "default options limit over max", query: "cursor=abc&limit=1001", want: mrtype.CursorParams{Value: "abc", Limit: 1000}},
		{
			name:  "custom param names",
			opts:  parser.ListCursorOptions{ParamNameCursor: "c", ParamNameLimit: "l"},
			query: "c=abc&l=20&cursor=def&limit=30",
			want:  mrtype.CursorParams{Value: "abc", Limit: 20},
		},
		{
			name:         "default limit over max",
			opts:         parser.ListCursorOptions{LimitMax: 5, LimitDefault: 10},
			query:        "cursor=abc",
			want:         mrtype.CursorParams{Value: "abc", Limit: 5},
			wantCtorWarn: true,
		},
		{
			name:  "max below builtin default",
			opts:  parser.ListCursorOptions{LimitMax: 5},
			query: "cursor=abc",
			want:  mrtype.CursorParams{Value: "abc", Limit: 5},
		},
		{
			name:         "default limit over max with incorrect limit",
			opts:         parser.ListCursorOptions{LimitMax: 5, LimitDefault: 10},
			query:        "cursor=abc&limit=x",
			want:         mrtype.CursorParams{Value: "", Limit: 5},
			wantCtorWarn: true,
			wantWarn:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			logger := &warnSpyLogger{Logger: mrlog.NopLogger()}
			r := httptest.NewRequest(http.MethodGet, "/?"+tt.query, http.NoBody)

			lc := parser.NewListCursor(logger, tt.opts)

			assert.Equal(t, tt.wantCtorWarn, logger.warnCount > 0)

			logger.warnCount = 0
			got := lc.CursorParams(r)

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantWarn, logger.warnCount > 0)
		})
	}
}
