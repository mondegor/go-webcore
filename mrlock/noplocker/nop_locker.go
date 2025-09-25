package noplocker

import (
	"context"
	"time"

	"github.com/mondegor/go-webcore/mrlock"
	"github.com/mondegor/go-webcore/mrsender"
)

const (
	nopLockerName = "NopLocker"
)

type (
	// Locker - заглушка реализующая интерфейс блокировщика указанного ключа.
	Locker struct {
		tracer mrsender.Tracer
	}
)

// New - создаёт объект Locker.
func New(tracer mrsender.Tracer) *Locker {
	return &Locker{
		tracer: tracer,
	}
}

// Lock - эмулирует блокировку указанного ключа с временем блокировки по умолчанию
// и возвращает функцию разблокировки.
func (l *Locker) Lock(ctx context.Context, key string) (unlock func(), err error) {
	return l.LockWithExpiry(ctx, key, 0)
}

// LockWithExpiry - эмулирует блокировку указанного ключа
// с указанием времени блокировки и возвращает функцию разблокировки.
func (l *Locker) LockWithExpiry(ctx context.Context, key string, expiry time.Duration) (unlock func(), err error) {
	if expiry == 0 {
		expiry = mrlock.DefaultExpiry
	}

	l.traceCmd(ctx, "Lock:"+expiry.String(), key)

	return func() {
		l.traceCmd(ctx, "Unlock", key)
	}, nil
}

func (l *Locker) traceCmd(ctx context.Context, command, key string) {
	l.tracer.Trace(
		ctx,
		"source", nopLockerName,
		"cmd", command,
		"key", key,
	)
}
