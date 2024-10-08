package mrparser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrparser"
)

// Make sure the Bool conforms with the mrserver.RequestParserBool interface.
func TestBoolImplementsRequestParserBool(t *testing.T) {
	assert.Implements(t, (*mrserver.RequestParserBool)(nil), &mrparser.Bool{})
}
