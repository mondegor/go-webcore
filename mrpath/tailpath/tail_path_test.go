package tailpath_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrpath"
	"github.com/mondegor/go-webcore/mrpath/tailpath"
)

// Make sure the Builder conforms with the mrpath.PathBuilder interface.
func TestBuilderImplementsPathBuilder(t *testing.T) {
	assert.Implements(t, (*mrpath.PathBuilder)(nil), &tailpath.Builder{})
}
