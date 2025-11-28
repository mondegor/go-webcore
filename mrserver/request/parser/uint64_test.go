package parser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrserver/request"
	"github.com/mondegor/go-webcore/mrserver/request/parser"
)

// Make sure the Uint64 conforms with the request.ParserUint64 interface.
func TestUint64ImplementsRequestParserUint64(t *testing.T) {
	assert.Implements(t, (*request.ParserUint64)(nil), &parser.Uint64{})
}
