package noplocker_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrlock"
	"github.com/mondegor/go-webcore/mrlock/noplocker"
)

// Make sure the Locker conforms with the mrlock.Locker interface.
func TestLockerImplementsLocker(t *testing.T) {
	assert.Implements(t, (*mrlock.Locker)(nil), &noplocker.Locker{})
}
