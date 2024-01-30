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

	l.debugCmd(ctx, "Lock:"+expiry.String(), key)

	return func() {
		l.debugCmd(ctx, "Unlock", key)
	}, nil
}

func (l *nopLocker) debugCmd(ctx context.Context, command, key string) {
	mrlog.Ctx(ctx).Debug().Msgf(
		"%s: cmd=%s, key=%s",
		nopLockerName,
		command,
		key,
	)
}
