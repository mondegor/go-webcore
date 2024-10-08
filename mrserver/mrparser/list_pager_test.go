package mrparser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrparser"
)

// Make sure the ListPager conforms with the mrserver.RequestParserListPager interface.
func TestListPagerImplementsRequestParserListPager(t *testing.T) {
	assert.Implements(t, (*mrserver.RequestParserListPager)(nil), &mrparser.ListPager{})
}
