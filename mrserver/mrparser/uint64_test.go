package mrparser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrparser"
)

// Make sure the Uint64 conforms with the mrserver.RequestParserUint64 interface.
func TestUint64ImplementsRequestParserUint64(t *testing.T) {
	assert.Implements(t, (*mrserver.RequestParserUint64)(nil), &mrparser.Uint64{})
}
