package mutexlocker

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlock"
	"github.com/mondegor/go-webcore/mrlog"
)

const (
	mutexLockerName = "MutexLocker"
)

type (
	// Locker - реализует интерфейс блокировщика указанного ключа основанный на mutex.
	Locker struct {
		mu   sync.Mutex
		keys map[string]int64
	}
)

// New - создаёт объект Locker.
func New(minBufferSize int) *Locker {
	return &Locker{
		keys: make(map[string]int64, minBufferSize),
	}
}

// Lock - блокирует указанный ключ в рамках приложения с временем блокировки по умолчанию
// и возвращает функцию разблокировки.
func (l *Locker) Lock(ctx context.Context, key string) (unlock func(), err error) {
	return l.LockWithExpiry(ctx, key, 0)
}

// LockWithExpiry - блокирует указанный ключ в рамках приложения с указанием
// времени блокировки и возвращает функцию разблокировки.
// Если указана expiry равная нулю, то будет установлено время блокировки по умолчанию.
func (l *Locker) LockWithExpiry(ctx context.Context, key string, expiry time.Duration) (unlock func(), err error) {
	if expiry == 0 {
		expiry = mrlock.DefaultExpiry
	}

	l.traceCmd(ctx, "Lock:"+expiry.String()+", keys-len="+strconv.Itoa(len(l.keys)), key)

	l.mu.Lock()
	defer l.mu.Unlock()

	if exp, ok := l.keys[key]; ok && exp > time.Now().UnixNano() {
		return nil, mrcore.ErrInternal.Wrap(fmt.Errorf("%s: key %s is blocked", mutexLockerName, key))
	}

	l.keys[key] = time.Now().UnixNano() + expiry.Nanoseconds()

	return func() {
		l.traceCmd(ctx, "Unlock", key)

		l.mu.Lock()
		delete(l.keys, key)
		l.mu.Unlock()
	}, nil
}

func (l *Locker) traceCmd(ctx context.Context, command, key string) {
	mrlog.Ctx(ctx).
		Trace().
		Str("source", mutexLockerName).
		Str("cmd", command).
		Str("key", key).
		Send()
}
