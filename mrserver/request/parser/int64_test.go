package parser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrserver/request"
	"github.com/mondegor/go-webcore/mrserver/request/parser"
)

// Make sure the Int64 conforms with the request.ParserInt64 interface.
func TestInt64ImplementsRequestParserInt64(t *testing.T) {
	assert.Implements(t, (*request.ParserInt64)(nil), &parser.Int64{})
}
