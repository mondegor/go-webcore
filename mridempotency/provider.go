package mridempotency

import (
	"context"
)

type (
	// Provider - провайдер предоставляющий методы, для поддержки идемпотентности запроса.
	Provider interface {
		Validate(key string) error
		Lock(ctx context.Context, key string) (unlock func(), err error)
		Get(ctx context.Context, key string) (Responser, error)
		Store(ctx context.Context, key string, response Responser) error
	}
)
