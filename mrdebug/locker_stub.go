package mrdebug

import (
	"context"
	"time"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrctx"
)

const (
	lockerStubName = "LockerStub"
)

type (
	lockerStub struct {
	}
)

func NewLockerStub() *lockerStub {
	return &lockerStub{}
}

func (l *lockerStub) Lock(ctx context.Context, key string) (mrcore.UnlockFunc, error) {
	return l.LockWithExpiry(ctx, key, 0)
}

// LockWithExpiry - if expiry = 0 then set expiry by default
func (l *lockerStub) LockWithExpiry(ctx context.Context, key string, expiry time.Duration) (mrcore.UnlockFunc, error) {
	if expiry == 0 {
		expiry = mrcore.LockerDefaultExpiry
	}

	l.debugCmd(ctx, "Lock:"+expiry.String(), key)

	return func() {
		l.debugCmd(ctx, "Unlock", key)
	}, nil
}

func (l *lockerStub) debugCmd(ctx context.Context, command, key string) {
	mrctx.Logger(ctx).Debug(
		"%s: cmd=%s, key=%s",
		lockerStubName,
		command,
		key,
	)
}
