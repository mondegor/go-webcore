package parser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrserver/request"
	"github.com/mondegor/go-webcore/mrserver/request/parser"
)

// Make sure the Bool conforms with the request.ParserBool interface.
func TestBoolImplementsRequestParserBool(t *testing.T) {
	assert.Implements(t, (*request.ParserBool)(nil), &parser.Bool{})
}
