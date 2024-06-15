package noplocker

import (
	"context"
	"time"

	"github.com/mondegor/go-webcore/mrlock"
	"github.com/mondegor/go-webcore/mrlog"
)

const (
	nopLockerName = "NopLocker"
)

type (
	// Locker - comment struct.
	Locker struct{}
)

// Make sure the Image conforms with the mrlock.Locker interface.
var _ mrlock.Locker = (*Locker)(nil)

// New - создаёт объект Locker.
func New() *Locker {
	return &Locker{}
}

// Lock - comment method.
func (l *Locker) Lock(ctx context.Context, key string) (unlock func(), err error) {
	return l.LockWithExpiry(ctx, key, 0)
}

// LockWithExpiry - if expiry = 0 then set expiry by default.
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
	mrlog.Ctx(ctx).
		Trace().
		Str("source", nopLockerName).
		Str("cmd", command).
		Str("key", key).
		Send()
}
