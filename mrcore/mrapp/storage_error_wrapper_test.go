package mrapp_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrcore/mrapp"
)

// Make sure the mrapp.StorageErrorWrapper conforms with the mrcore.StorageErrorWrapper interface.
func TestStorageErrorWrapperImplementsStorageErrorWrapper(t *testing.T) {
	assert.Implements(t, (*mrcore.StorageErrorWrapper)(nil), &mrapp.StorageErrorWrapper{})
}
