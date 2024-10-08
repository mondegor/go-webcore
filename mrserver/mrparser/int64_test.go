package mrparser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrparser"
)

// Make sure the Int64 conforms with the mrserver.RequestParserInt64 interface.
func TestInt64ImplementsRequestParserInt64(t *testing.T) {
	assert.Implements(t, (*mrserver.RequestParserInt64)(nil), &mrparser.Int64{})
}
