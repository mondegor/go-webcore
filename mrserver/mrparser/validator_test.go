package mrparser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrparser"
)

// Make sure the Validator conforms with the mrserver.RequestParserValidate interface.
func TestValidatorImplementsRequestParserValidate(t *testing.T) {
	assert.Implements(t, (*mrserver.RequestParserValidate)(nil), &mrparser.Validator{})
}
