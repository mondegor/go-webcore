package nopprovider

import (
	"context"
	"time"

	"github.com/mondegor/go-sysmess/mrtrace"

	"github.com/mondegor/go-webcore/mridempotency"
	"github.com/mondegor/go-webcore/mridempotency/nopresponser"
)

const (
	nopProviderName = "IdempotencyNopProvider"
	defaultExpiry   = time.Second
)

type (
	// Provider - заглушка реализующая интерфейс идемпотентности запроса.
	Provider struct {
		tracer mrtrace.Tracer
	}
)

// New - создаёт объект Provider.
func New(tracer mrtrace.Tracer) *Provider {
	return &Provider{
		tracer: tracer,
	}
}

// Validate - эмулирует успешную валидацию данных.
func (l *Provider) Validate(_ string) error {
	return nil
}

// Lock - эмулирует блокировку указанного ключа идемпотентности и возвращает функцию разблокировки.
func (l *Provider) Lock(ctx context.Context, key string) (unlock func(), err error) {
	l.traceCmd(ctx, "Lock:"+defaultExpiry.String(), key)

	return func() {
		l.traceCmd(ctx, "Unlock", key)
	}, nil
}

// Get - всегда возвращает пустой ответ.
func (l *Provider) Get(ctx context.Context, key string) (mridempotency.Responser, error) {
	l.traceCmd(ctx, "Get:"+key, key)

	return nopresponser.New(), nil
}

// Store - эмулирует сохранение результата по указанному ключу.
func (l *Provider) Store(ctx context.Context, key string, result mridempotency.Responser) error {
	l.traceCmd(ctx, "Store:"+key, key)

	l.tracer.Trace(
		ctx,
		"statusCode", result.StatusCode(),
		"body", result.Content(),
	)

	return nil
}

func (l *Provider) traceCmd(ctx context.Context, command, key string) {
	l.tracer.Trace(
		ctx,
		"source", nopProviderName,
		"cmd", command,
		"key", key,
	)
}
