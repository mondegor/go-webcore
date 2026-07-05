package parser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrserver/request"
	"github.com/mondegor/go-webcore/mrserver/request/parser"
)

// Make sure the String conforms with the request.ParserString interface.
func TestStringImplementsRequestParserString(t *testing.T) {
	t.Parallel()

	assert.Implements(t, (*request.ParserString)(nil), &parser.String{})
}
