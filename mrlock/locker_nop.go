package mrlock

import (
	"context"
	"time"

	"github.com/mondegor/go-webcore/mrlog"
)

const (
	nopLockerName = "NopLocker"
)

type (
	nopLocker struct {
	}
)

func NewNopLocker() Locker {
	return &nopLocker{}
}

func (l *nopLocker) Lock(ctx context.Context, key string) (UnlockFunc, error) {
	return l.LockWithExpiry(ctx, key, 0)
}

// LockWithExpiry - if expiry = 0 then set expiry by default
func (l *nopLocker) LockWithExpiry(ctx context.Context, key string, expiry time.Duration) (UnlockFunc, error) {
	if expiry == 0 {
		expiry = DefaultExpiry
	}

	l.traceCmd(ctx, "Lock:"+expiry.String(), key)

	return func() {
		l.traceCmd(ctx, "Unlock", key)
	}, nil
}

func (l *nopLocker) traceCmd(ctx context.Context, command, key string) {
	mrlog.Ctx(ctx).
		Trace().
		Str("source", nopLockerName).
		Str("cmd", command).
		Str("key", key).
		Send()
}
