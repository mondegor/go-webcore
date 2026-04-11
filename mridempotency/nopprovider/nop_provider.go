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
	// Provider - заглушка (no-op) реализующая интерфейс Provider идемпотентности.
	// Все методы всегда возвращают успешный результат, не выполняя реальной работы.
	// Используется для тестирования, отладки или когда идемпотентность не требуется.
	Provider struct {
		tracer mrtrace.Tracer
	}
)

// New - создаёт no-op провайдер идемпотентности для тестирования и отладки.
func New(tracer mrtrace.Tracer) *Provider {
	return &Provider{
		tracer: tracer,
	}
}

// Validate - всегда возвращает nil, эмулируя успешную валидацию любого ключа.
func (l *Provider) Validate(_ string) error {
	return nil
}

// Lock - эмулирует блокировку ключа идемпотентности без реальной синхронизации.
func (l *Provider) Lock(ctx context.Context, key string) (unlock func(), err error) {
	l.traceCmd(ctx, "Lock:"+defaultExpiry.String(), key)

	return func() {
		l.traceCmd(ctx, "Unlock", key)
	}, nil
}

// Get - всегда возвращает пустой ответ (nopresponser.New).
func (l *Provider) Get(ctx context.Context, key string) (mridempotency.Responser, error) {
	l.traceCmd(ctx, "Get:"+key, key)

	return nopresponser.New(), nil
}

// Save - эмулирует сохранение ответа без реальной записи в хранилище.
func (l *Provider) Save(ctx context.Context, key string, result mridempotency.Responser) error {
	l.traceCmd(ctx, "Save:"+key, key)

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
