package noplog

import (
	"github.com/mondegor/go-webcore/mrlog"
)

type (
	contextAdapter struct {
		logger LoggerAdapter
	}
)

// Logger - comment method.
func (c *contextAdapter) Logger() mrlog.Logger {
	return &c.logger
}

// Str - comment method.
func (c *contextAdapter) Str(_, _ string) mrlog.LoggerContext {
	return c
}

// Bytes - comment method.
func (c *contextAdapter) Bytes(_ string, _ []byte) mrlog.LoggerContext {
	return c
}

// Int - comment method.
func (c *contextAdapter) Int(_ string, _ int) mrlog.LoggerContext {
	return c
}

// Any - comment method.
func (c *contextAdapter) Any(_ string, _ any) mrlog.LoggerContext {
	return c
}
