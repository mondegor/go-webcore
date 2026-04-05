package mridempotency

import (
	"context"
)

type (
	// Provider - обеспечивает поддержку идемпотентности запросов,
	// предоставляя методы для валидации, блокировки, сохранения и получения ответов.
	Provider interface {
		Validate(key string) error
		Lock(ctx context.Context, key string) (unlock func(), err error)
		Get(ctx context.Context, key string) (Responser, error)
		Store(ctx context.Context, key string, response Responser) error
	}
)
