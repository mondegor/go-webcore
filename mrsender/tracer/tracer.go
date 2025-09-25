package tracer

import (
	"context"

	"github.com/mondegor/go-sysmess/mrlog"
)

type (
	// Tracer - comment struct.
	Tracer struct {
		logger mrlog.Logger
	}
)

// New - создаёт объект Tracer.
func New(logger mrlog.Logger) *Tracer {
	return &Tracer{
		logger: logger,
	}
}

// Enabled - comment method.
func (e *Tracer) Enabled() bool {
	return e.logger != nil
}

// Trace - comment method.
func (e *Tracer) Trace(ctx context.Context, args ...any) {
	if e.logger != nil {
		e.logger.Log(ctx, mrlog.LevelTrace, "---", args...)
	}
}
