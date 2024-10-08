package mrlog_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrlog"
)

// Make sure the eventAdapter conforms with the LoggerEvent interface.
func TestEventAdapterImplementsLoggerEvent(t *testing.T) {
	assert.Implements(t, (*mrlog.LoggerEvent)(nil), (&mrlog.BaseLogger{}).Error())
}
