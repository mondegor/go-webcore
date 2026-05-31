package parser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrserver/request"
	"github.com/mondegor/go-webcore/mrserver/request/parser"
)

// Make sure the ListSorter conforms with the request.ParserListSorter interface.
func TestListSorterImplementsRequestParserListSorter(t *testing.T) {
	t.Parallel()

	assert.Implements(t, (*request.ParserListSorter)(nil), &parser.ListSorter{})
}
