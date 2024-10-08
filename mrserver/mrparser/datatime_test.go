package mrparser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrparser"
)

// Make sure the DateTime conforms with the mrserver.RequestParserDateTime interface.
func TestDateTimeImplementsRequestParserDateTime(t *testing.T) {
	assert.Implements(t, (*mrserver.RequestParserDateTime)(nil), &mrparser.DateTime{})
}
