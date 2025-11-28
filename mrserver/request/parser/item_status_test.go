package parser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrserver/request"
	"github.com/mondegor/go-webcore/mrserver/request/parser"
)

// Make sure the ItemStatus conforms with the request.ParserItemStatus interface.
func TestItemStatusImplementsRequestParserItemStatus(t *testing.T) {
	assert.Implements(t, (*request.ParserItemStatus)(nil), &parser.ItemStatus{})
}
