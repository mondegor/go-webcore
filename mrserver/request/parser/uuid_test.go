package parser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrserver/request"
	"github.com/mondegor/go-webcore/mrserver/request/parser"
)

// Make sure the UUID conforms with the request.ParserUUID interface.
func TestUUIDImplementsRequestParserUUID(t *testing.T) {
	assert.Implements(t, (*request.ParserUUID)(nil), &parser.UUID{})
}
