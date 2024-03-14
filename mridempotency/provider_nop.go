package mridempotency

import (
	"context"
	"time"

	"github.com/mondegor/go-webcore/mrlock"
	"github.com/mondegor/go-webcore/mrlog"
)

const (
	nopProviderName = "IdempotencyNopProvider"
	defaultExpiry   = time.Second
)

type (
	nopProvider struct {
	}
)

func NewNopProvider() Provider {
	return &nopProvider{}
}

func (l *nopProvider) Validate(key string) error {
	return nil
}

func (l *nopProvider) Lock(ctx context.Context, key string) (mrlock.UnlockFunc, error) {
	l.traceCmd(ctx, "Lock:"+defaultExpiry.String(), key)

	return func() {
		l.traceCmd(ctx, "Unlock", key)
	}, nil
}

func (l *nopProvider) Get(ctx context.Context, key string) (Response, error) {
	l.traceCmd(ctx, "Get:"+key, key)

	return nil, nil
}

func (l *nopProvider) Store(ctx context.Context, key string, result Response) error {
	l.traceCmd(ctx, "Store:"+key, key)

	mrlog.Ctx(ctx).
		Trace().
		Int("statusCode", result.StatusCode()).
		Bytes("body", result.Body()).
		Send()

	return nil
}

func (l *nopProvider) traceCmd(ctx context.Context, command, key string) {
	mrlog.Ctx(ctx).
		Trace().
		Str("source", nopProviderName).
		Str("cmd", command).
		Str("key", key).
		Send()
}
