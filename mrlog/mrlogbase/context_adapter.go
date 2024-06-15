package mrlogbase

import (
	"github.com/mondegor/go-webcore/mrlog"
)

type (
	contextAdapter struct {
		logger LoggerAdapter
	}
)

// Make sure the Image conforms with the mrlog.LoggerContext interface.
var _ mrlog.LoggerContext = (*contextAdapter)(nil)

// Logger - comment method.
func (c *contextAdapter) Logger() mrlog.Logger {
	return &c.logger
}

// Str - comment method.
func (c *contextAdapter) Str(key, value string) mrlog.LoggerContext {
	cp := *c
	cp.logger.context = appendValue(appendKey(cp.logger.context, key), value)

	return &cp
}

// Bytes - comment method.
func (c *contextAdapter) Bytes(key string, value []byte) mrlog.LoggerContext {
	cp := *c
	cp.logger.context = appendValue(appendKey(cp.logger.context, key), value)

	return &cp
}

// Int - comment method.
func (c *contextAdapter) Int(key string, value int) mrlog.LoggerContext {
	cp := *c
	cp.logger.context = appendValue(appendKey(cp.logger.context, key), value)

	return &cp
}

// Any - comment method.
func (c *contextAdapter) Any(key string, value any) mrlog.LoggerContext {
	cp := *c
	cp.logger.context = appendValue(appendKey(cp.logger.context, key), value)

	return &cp
}
