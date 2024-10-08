package mrlog_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrlog"
)

// Make sure the BaseLogger conforms with the Logger interface.
func TestBaseLoggerImplementsLogger(t *testing.T) {
	assert.Implements(t, (*mrlog.Logger)(nil), &mrlog.BaseLogger{})
}
