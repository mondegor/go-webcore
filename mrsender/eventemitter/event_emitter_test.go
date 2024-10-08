package eventemitter_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrsender"
	"github.com/mondegor/go-webcore/mrsender/eventemitter"
)

// Make sure the Emitter conforms with the mrsender.EventEmitter interface.
func TestEmitterImplementsEventEmitter(t *testing.T) {
	assert.Implements(t, (*mrsender.EventEmitter)(nil), &eventemitter.Emitter{})
}
