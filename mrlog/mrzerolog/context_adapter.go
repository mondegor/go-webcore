package mrzerolog

import (
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/rs/zerolog"
)

// go get -u github.com/rs/zerolog

type (
	ContextAdapter struct {
		zc   zerolog.Context
		opts loggerOptions
	}
)

func (c *ContextAdapter) Logger() mrlog.Logger {
	return &LoggerAdapter{
		zl:   c.zc.Logger(),
		opts: c.opts,
	}
}

func (c *ContextAdapter) CallerWithSkipFrame(count int) mrlog.LoggerContext {
	return &ContextAdapter{
		zc: c.zc.CallerWithSkipFrameCount(count + 3), // +2 zerolog infra
		opts: loggerOptions{
			level: c.opts.level,
			// set isAutoCallerOnFunc = nil
		},
	}
}

func (c *ContextAdapter) Str(key, value string) mrlog.LoggerContext {
	return &ContextAdapter{
		zc:   c.zc.Str(key, value),
		opts: c.opts,
	}
}

func (c *ContextAdapter) Bytes(key string, value []byte) mrlog.LoggerContext {
	return &ContextAdapter{
		zc:   c.zc.Bytes(key, value),
		opts: c.opts,
	}
}

func (c *ContextAdapter) Int(key string, value int) mrlog.LoggerContext {
	return &ContextAdapter{
		zc:   c.zc.Int(key, value),
		opts: c.opts,
	}
}

func (c *ContextAdapter) Any(key string, value any) mrlog.LoggerContext {
	return &ContextAdapter{
		zc:   c.zc.Any(key, value),
		opts: c.opts,
	}
}
