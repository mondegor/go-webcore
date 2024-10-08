package mrparser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrparser"
)

// Make sure the File conforms with the mrserver.RequestParserFile interface.
func TestFileImplementsRequestParserFile(t *testing.T) {
	assert.Implements(t, (*mrserver.RequestParserFile)(nil), &mrparser.File{})
}
