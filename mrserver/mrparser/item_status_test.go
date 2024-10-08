package mrparser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrparser"
)

// Make sure the ItemStatus conforms with the mrserver.RequestParserItemStatus interface.
func TestItemStatusImplementsRequestParserItemStatus(t *testing.T) {
	assert.Implements(t, (*mrserver.RequestParserItemStatus)(nil), &mrparser.ItemStatus{})
}
