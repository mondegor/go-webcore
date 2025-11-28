package parser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrserver/request"
	"github.com/mondegor/go-webcore/mrserver/request/parser"
)

// Make sure the Float64 conforms with the request.ParserFloat64 interface.
func TestFloat64ImplementsRequestParserFloat64(t *testing.T) {
	assert.Implements(t, (*request.ParserFloat64)(nil), &parser.Float64{})
}
