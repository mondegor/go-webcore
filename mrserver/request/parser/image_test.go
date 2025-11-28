package parser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrserver/request"
	"github.com/mondegor/go-webcore/mrserver/request/parser"
)

// Make sure the Image conforms with the request.ParserImage interface.
func TestImageImplementsRequestParserImage(t *testing.T) {
	assert.Implements(t, (*request.ParserImage)(nil), &parser.Image{})
}
