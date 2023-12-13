package mrlib

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrctx"
)

const (
	lockerMutexName = "LockerMutex"
)

type (
	lockerMutex struct {
		keysMu sync.Mutex
		keys   map[string]int64
	}
)

func NewLockerMutex() *lockerMutex {
	return &lockerMutex{
		keys: make(map[string]int64, 16),
	}
}

func (l *lockerMutex) Lock(ctx context.Context, key string) (mrcore.UnlockFunc, error) {
	return l.LockWithExpiry(ctx, key, 0)
}

// LockWithExpiry - if expiry = 0 then set expiry by default
func (l *lockerMutex) LockWithExpiry(ctx context.Context, key string, expiry time.Duration) (mrcore.UnlockFunc, error) {
	if expiry == 0 {
		expiry = mrcore.LockerDefaultExpiry
	}

	l.debugCmd(ctx, "Lock:"+expiry.String()+", keys-len="+strconv.Itoa(len(l.keys)), key)

	l.keysMu.Lock()
	defer l.keysMu.Unlock()

	if exp, ok := l.keys[key]; ok && exp > time.Now().UnixNano() {
		return nil, fmt.Errorf("%s: key %s is blocked", lockerMutexName, key)
	}

	l.keys[key] = time.Now().UnixNano() + expiry.Nanoseconds()

	return func() {
		l.debugCmd(ctx, "Unlock", key)

		l.keysMu.Lock()
		delete(l.keys, key)
		l.keysMu.Unlock()
	}, nil
}

func (l *lockerMutex) debugCmd(ctx context.Context, command, key string) {
	mrctx.Logger(ctx).Debug(
		"%s: cmd=%s, key=%s",
		lockerMutexName,
		command,
		key,
	)
}
