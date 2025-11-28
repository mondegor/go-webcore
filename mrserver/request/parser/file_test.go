package parser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrserver/request"
	"github.com/mondegor/go-webcore/mrserver/request/parser"
)

// Make sure the File conforms with the request.ParserFile interface.
func TestFileImplementsRequestParserFile(t *testing.T) {
	assert.Implements(t, (*request.ParserFile)(nil), &parser.File{})
}
