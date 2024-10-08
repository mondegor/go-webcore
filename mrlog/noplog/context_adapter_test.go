package noplog_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrlog/noplog"
)

// Make sure the contextAdapter conforms with the mrlog.LoggerContext interface.
func TestContextAdapterImplementsLoggerContext(t *testing.T) {
	assert.Implements(t, (*mrlog.LoggerContext)(nil), (&noplog.LoggerAdapter{}).With())
}
