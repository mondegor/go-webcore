package mrplayvalidator_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrview"
	"github.com/mondegor/go-webcore/mrview/mrplayvalidator"
)

// Make sure the ValidatorAdapter conforms with the mrview.Validator interface.
func TestValidatorAdapterImplementsValidator(t *testing.T) {
	assert.Implements(t, (*mrview.Validator)(nil), &mrplayvalidator.ValidatorAdapter{})
}
