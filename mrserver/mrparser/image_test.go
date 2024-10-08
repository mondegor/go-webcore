package mrparser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrparser"
)

// Make sure the Image conforms with the mrserver.RequestParserImage interface.
func TestImageImplementsRequestParserImage(t *testing.T) {
	assert.Implements(t, (*mrserver.RequestParserImage)(nil), &mrparser.Image{})
}
