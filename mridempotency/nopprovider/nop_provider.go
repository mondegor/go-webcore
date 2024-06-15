package nopprovider

import (
	"context"
	"time"

	"github.com/mondegor/go-webcore/mridempotency"
	"github.com/mondegor/go-webcore/mridempotency/nopresponser"
	"github.com/mondegor/go-webcore/mrlog"
)

const (
	nopProviderName = "IdempotencyNopProvider"
	defaultExpiry   = time.Second
)

type (
	// Provider - comment struct.
	Provider struct{}
)

// Make sure the Image conforms with the mridempotency.Provider interface.
var _ mridempotency.Provider = (*Provider)(nil)

// New - создаёт объект Provider.
func New() *Provider {
	return &Provider{}
}

// Validate - comment method.
func (l *Provider) Validate(_ string) error {
	return nil
}

// Lock - comment method.
func (l *Provider) Lock(ctx context.Context, key string) (unlock func(), err error) {
	l.traceCmd(ctx, "Lock:"+defaultExpiry.String(), key)

	return func() {
		l.traceCmd(ctx, "Unlock", key)
	}, nil
}

// Get - comment method.
func (l *Provider) Get(ctx context.Context, key string) (mridempotency.Responser, error) { //nolint:ireturn,nolintlint
	l.traceCmd(ctx, "Get:"+key, key)

	return nopresponser.New(), nil
}

// Store - comment method.
func (l *Provider) Store(ctx context.Context, key string, result mridempotency.Responser) error {
	l.traceCmd(ctx, "Store:"+key, key)

	mrlog.Ctx(ctx).
		Trace().
		Int("statusCode", result.StatusCode()).
		Bytes("body", result.Body()).
		Send()

	return nil
}

func (l *Provider) traceCmd(ctx context.Context, command, key string) {
	mrlog.Ctx(ctx).
		Trace().
		Str("source", nopProviderName).
		Str("cmd", command).
		Str("key", key).
		Send()
}
