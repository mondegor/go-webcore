package parser_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/request"
	"github.com/mondegor/go-webcore/mrserver/request/parser"
)

// Make sure the User conforms with the request.ParserUser interface.
func TestUserImplementsRequestParserUser(t *testing.T) {
	t.Parallel()

	assert.Implements(t, (*request.ParserUser)(nil), &parser.User{})
}

func TestUser_SessionID(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name      string
		headerSet bool
		header    string
		want      string
	}

	tests := []testCase{
		{
			name:      "header is set",
			headerSet: true,
			header:    "8f14e45f",
			want:      "8f14e45f",
		},
		{
			name:      "header is empty",
			headerSet: true,
			header:    "",
			want:      "",
		},
		{
			name:      "header is absent",
			headerSet: false,
			want:      "",
		},
	}

	p := parser.NewUser(mrlog.NopLogger())

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			r := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
			if tc.headerSet {
				r.Header.Set(mrserver.HeaderKeySessionID, tc.header)
			}

			assert.Equal(t, tc.want, p.SessionID(r))
		})
	}
}
