package mrparser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrparser"
)

// Make sure the Float64 conforms with the mrserver.RequestParserFloat64 interface.
func TestFloat64ImplementsRequestParserFloat64(t *testing.T) {
	assert.Implements(t, (*mrserver.RequestParserFloat64)(nil), &mrparser.Float64{})
}
