package mrparser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrparser"
)

// Make sure the ListSorter conforms with the mrserver.RequestParserListSorter interface.
func TestListSorterImplementsRequestParserListSorter(t *testing.T) {
	assert.Implements(t, (*mrserver.RequestParserListSorter)(nil), &mrparser.ListSorter{})
}
