package parser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrserver/request"
	"github.com/mondegor/go-webcore/mrserver/request/parser"
)

// Make sure the ListCursor conforms with the request.ParserListCursor interface.
func TestListCursorImplementsRequestParserListCursor(t *testing.T) {
	assert.Implements(t, (*request.ParserListCursor)(nil), &parser.ListCursor{})
}
