package embedded

import (
	"context"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrtrace/distribute"
)

type (
	// ContextGenerator - comment struct.
	ContextGenerator struct {
		idGenerator mrcore.IdentifierGenerator
	}
)

// NewContextGenerator - создаёт объект ContextGenerator.
func NewContextGenerator(idGenerator mrcore.IdentifierGenerator) *ContextGenerator {
	return &ContextGenerator{
		idGenerator: idGenerator,
	}
}

// WithProcessIDContext - comment method.
func (e *ContextGenerator) WithProcessIDContext(ctx context.Context) context.Context {
	return distribute.WithProcessID(ctx, e.idGenerator.GenID())
}

// WithWorkerIDContext - comment method.
func (e *ContextGenerator) WithWorkerIDContext(ctx context.Context) context.Context {
	return distribute.WithWorkerID(ctx, e.idGenerator.GenID())
}

// WithTaskIDContext - comment method.
func (e *ContextGenerator) WithTaskIDContext(ctx context.Context) context.Context {
	return distribute.WithTaskID(ctx, e.idGenerator.GenID())
}

// WithCorrelationIDContext - comment method.
func (e *ContextGenerator) WithCorrelationIDContext(ctx context.Context) context.Context {
	return distribute.WithCorrelationID(ctx, e.idGenerator.GenID())
}

// WithRequestIDContext - comment method.
func (e *ContextGenerator) WithRequestIDContext(ctx context.Context) context.Context {
	return distribute.WithRequestID(ctx, e.idGenerator.GenID())
}

// NewContextWithIDs - comment method.
func (e *ContextGenerator) NewContextWithIDs(originalCtx context.Context) context.Context {
	return distribute.NewContextWithIDs(originalCtx)
}

// ExtractCorrelationID - comment method.
func (e *ContextGenerator) ExtractCorrelationID(ctx context.Context) string {
	return distribute.FindCorrelationID(ctx)
}
