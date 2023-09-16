package mrcore

import (
    "context"
    "time"
)

type (
    Locker interface {
        Lock(ctx context.Context, key string) (UnlockFunc, error)
        LockWithExpiry(ctx context.Context, key string, expiry time.Duration) (UnlockFunc, error)
    }

    UnlockFunc func()
)
