package nopresponser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mridempotency"
	"github.com/mondegor/go-webcore/mridempotency/nopresponser"
)

// Make sure the Responser conforms with the mridempotency.Responser interface.
func TestResponserImplementsResponser(t *testing.T) {
	assert.Implements(t, (*mridempotency.Responser)(nil), &nopresponser.Responser{})
}
