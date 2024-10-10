package mrapp_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrcore/mrapp"
)

// Make sure the mrapp.UseCaseErrorWrapper conforms with the mrcore.UseCaseErrorWrapper interface.
func TestUseCaseErrorWrapperImplementsUseCaseErrorWrapper(t *testing.T) {
	assert.Implements(t, (*mrcore.UseCaseErrorWrapper)(nil), &mrapp.UseCaseErrorWrapper{})
}
