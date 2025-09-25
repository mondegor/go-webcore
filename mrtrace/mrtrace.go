package mrtrace

import "context"

type (
	// ContextEmbedder - comment interface.
	ContextEmbedder interface {
		WithProcessIDContext(ctx context.Context) context.Context
		WithWorkerIDContext(ctx context.Context) context.Context
		WithTaskIDContext(ctx context.Context) context.Context
		WithCorrelationIDContext(ctx context.Context) context.Context
		WithRequestIDContext(ctx context.Context) context.Context

		NewContextWithIDs(originalCtx context.Context) context.Context
		ExtractCorrelationID(ctx context.Context) string
	}
)
