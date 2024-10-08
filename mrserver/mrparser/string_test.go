package mrparser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrparser"
)

// Make sure the String conforms with the mrserver.RequestParserString interface.
func TestStringImplementsRequestParserString(t *testing.T) {
	assert.Implements(t, (*mrserver.RequestParserString)(nil), &mrparser.String{})
}
