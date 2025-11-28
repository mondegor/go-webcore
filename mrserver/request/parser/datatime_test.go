package parser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrserver/request"
	"github.com/mondegor/go-webcore/mrserver/request/parser"
)

// Make sure the DateTime conforms with the request.ParserDateTime interface.
func TestDateTimeImplementsRequestParserDateTime(t *testing.T) {
	assert.Implements(t, (*request.ParserDateTime)(nil), &parser.DateTime{})
}
