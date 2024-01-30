package mrlock

import (
	"context"
	"time"
)

const (
	DefaultExpiry = time.Second
)

type (
	Locker interface {
		Lock(ctx context.Context, key string) (UnlockFunc, error)
		LockWithExpiry(ctx context.Context, key string, expiry time.Duration) (UnlockFunc, error)
	}

	UnlockFunc func()
)
