package parser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrserver/request"
	"github.com/mondegor/go-webcore/mrserver/request/parser"
)

// Make sure the Validator conforms with the request.ParserValidate interface.
func TestValidatorImplementsRequestParserValidate(t *testing.T) {
	assert.Implements(t, (*request.ParserValidate)(nil), &parser.Validator{})
}
