package mridempotency

import (
	"context"

	"github.com/mondegor/go-webcore/mrlock"
)

type (
	Provider interface {
		Validate(key string) error
		Lock(ctx context.Context, key string) (mrlock.UnlockFunc, error)
		Get(ctx context.Context, key string) (Response, error)
		Store(ctx context.Context, key string, response Response) error
	}

	Response interface {
		StatusCode() int
		Body() []byte
	}
)
