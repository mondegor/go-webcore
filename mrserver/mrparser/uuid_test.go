package mrparser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrparser"
)

// Make sure the UUID conforms with the mrserver.RequestParserUUID interface.
func TestUUIDImplementsRequestParserUUID(t *testing.T) {
	assert.Implements(t, (*mrserver.RequestParserUUID)(nil), &mrparser.UUID{})
}
