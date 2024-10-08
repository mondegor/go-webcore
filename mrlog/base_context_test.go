package mrlog_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrlog"
)

// Make sure the contextAdapter conforms with the LoggerContext interface.
func TestContextAdapterImplementsLoggerContext(t *testing.T) {
	assert.Implements(t, (*mrlog.LoggerContext)(nil), (&mrlog.BaseLogger{}).With())
}
