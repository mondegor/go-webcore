package zerolog_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrlog/zerolog"
)

// Make sure the eventAdapter conforms with the mrlog.LoggerEvent interface.
func TestEventAdapterImplementsLoggerEvent(t *testing.T) {
	assert.Implements(t, (*mrlog.LoggerEvent)(nil), (&zerolog.LoggerAdapter{}).Error())
}
