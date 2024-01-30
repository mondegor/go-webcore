package mrlock

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/mondegor/go-webcore/mrlog"
)

const (
	mutexLockerName = "MutexLocker"
)

type (
	mutexLocker struct {
		keysMu sync.Mutex
		keys   map[string]int64
	}
)

func NewMutexLocker() Locker {
	return &mutexLocker{
		keys: make(map[string]int64, 16),
	}
}

func (l *mutexLocker) Lock(ctx context.Context, key string) (UnlockFunc, error) {
	return l.LockWithExpiry(ctx, key, 0)
}

// LockWithExpiry - if expiry = 0 then set expiry by default
func (l *mutexLocker) LockWithExpiry(ctx context.Context, key string, expiry time.Duration) (UnlockFunc, error) {
	if expiry == 0 {
		expiry = DefaultExpiry
	}

	l.debugCmd(ctx, "Lock:"+expiry.String()+", keys-len="+strconv.Itoa(len(l.keys)), key)

	l.keysMu.Lock()
	defer l.keysMu.Unlock()

	if exp, ok := l.keys[key]; ok && exp > time.Now().UnixNano() {
		return nil, fmt.Errorf("%s: key %s is blocked", mutexLockerName, key)
	}

	l.keys[key] = time.Now().UnixNano() + expiry.Nanoseconds()

	return func() {
		l.debugCmd(ctx, "Unlock", key)

		l.keysMu.Lock()
		delete(l.keys, key)
		l.keysMu.Unlock()
	}, nil
}

func (l *mutexLocker) debugCmd(ctx context.Context, command, key string) {
	mrlog.Ctx(ctx).Debug().Msgf(
		"%s: cmd=%s, key=%s",
		mutexLockerName,
		command,
		key,
	)
}
