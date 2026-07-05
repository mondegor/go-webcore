package parser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrserver/request"
	"github.com/mondegor/go-webcore/mrserver/request/parser"
)

// Make sure the ListPager conforms with the request.ParserListPager interface.
func TestListPagerImplementsRequestParserListPager(t *testing.T) {
	t.Parallel()

	assert.Implements(t, (*request.ParserListPager)(nil), &parser.ListPager{})
}
