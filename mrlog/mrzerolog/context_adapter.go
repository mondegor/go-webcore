package mrzerolog

import (
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/rs/zerolog"
)

// go get -u github.com/rs/zerolog

type (
	ContextAdapter struct {
		zc zerolog.Context
	}
)

func (c *ContextAdapter) Logger() mrlog.Logger {
	l := &LoggerAdapter{
		zl: c.zc.Logger(),
	}

	l.level = mrlog.Level(l.zl.GetLevel())

	return l
}

func (c *ContextAdapter) CallerWithSkipFrame(count int) mrlog.LoggerContext {
	return &ContextAdapter{zc: c.zc.CallerWithSkipFrameCount(count)}
}

func (c *ContextAdapter) Str(key, value string) mrlog.LoggerContext {
	return &ContextAdapter{zc: c.zc.Str(key, value)}
}

func (c *ContextAdapter) Bytes(key string, value []byte) mrlog.LoggerContext {
	return &ContextAdapter{zc: c.zc.Bytes(key, value)}
}

func (c *ContextAdapter) Int(key string, value int) mrlog.LoggerContext {
	return &ContextAdapter{zc: c.zc.Int(key, value)}
}

func (c *ContextAdapter) Any(key string, value any) mrlog.LoggerContext {
	return &ContextAdapter{zc: c.zc.Any(key, value)}
}
