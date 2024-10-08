package zerolog_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrlog/zerolog"
)

// Make sure the LoggerAdapter conforms with the mrlog.Logger interface.
func TestLoggerAdapterImplementsLogger(t *testing.T) {
	assert.Implements(t, (*mrlog.Logger)(nil), &zerolog.LoggerAdapter{})
}
