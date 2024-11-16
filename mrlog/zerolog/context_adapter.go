package zerolog

import (
	"github.com/rs/zerolog"

	"github.com/mondegor/go-webcore/mrlog"
)

type (
	// contextAdapter - comment struct.
	contextAdapter struct {
		zc   zerolog.Context
		opts loggerOptions
	}
)

// Logger - comment method.
func (c *contextAdapter) Logger() mrlog.Logger {
	return &LoggerAdapter{
		zl:   c.zc.Logger(),
		opts: c.opts,
	}
}

// Str - comment method.
func (c *contextAdapter) Str(key, value string) mrlog.LoggerContext {
	return &contextAdapter{
		zc:   c.zc.Str(key, value),
		opts: c.opts,
	}
}

// Bytes - comment method.
func (c *contextAdapter) Bytes(key string, value []byte) mrlog.LoggerContext {
	return &contextAdapter{
		zc:   c.zc.Bytes(key, value),
		opts: c.opts,
	}
}

// Int - comment method.
func (c *contextAdapter) Int(key string, value int) mrlog.LoggerContext {
	return &contextAdapter{
		zc:   c.zc.Int(key, value),
		opts: c.opts,
	}
}

// Int64 - comment method.
func (c *contextAdapter) Int64(key string, value int64) mrlog.LoggerContext {
	return &contextAdapter{
		zc:   c.zc.Int64(key, value),
		opts: c.opts,
	}
}

// Any - comment method.
func (c *contextAdapter) Any(key string, value any) mrlog.LoggerContext {
	return &contextAdapter{
		zc:   c.zc.Any(key, value),
		opts: c.opts,
	}
}
