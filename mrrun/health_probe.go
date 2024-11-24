package mrrun

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlog"
)

const (
	defaultTimeout = 5 * time.Second
)

type (
	// HealthProbe - обёртка для проверки работоспособности какого либо процесса.
	HealthProbe struct {
		caption string                          // название проверяемого процесса
		check   func(ctx context.Context) error // функция проверки работоспособности процесса
		timeout time.Duration                   // таймаут, после которого функция check должна остановить своё выполнение
	}
)

// NewHealthProbe - создаёт объект AppRunner.
func NewHealthProbe(caption string, check func(ctx context.Context) error, timeout time.Duration) *HealthProbe {
	if timeout == 0 {
		timeout = defaultTimeout
	}

	return &HealthProbe{
		caption: caption,
		check:   check,
		timeout: timeout,
	}
}

// Caption - возвращает название пробы.
func (p *HealthProbe) Caption() string {
	return p.caption
}

// Check - метод вызова проверки пробы.
func (p *HealthProbe) Check(ctx context.Context) (err error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer func() {
		cancel()

		if rvr := recover(); rvr != nil {
			mrlog.Ctx(ctx).
				Error().
				Err(
					mrcore.ErrInternalCaughtPanic.New(
						"probe: "+p.caption,
						rvr,
						string(debug.Stack()),
					),
				).Send()

			err = fmt.Errorf("probe '%s' has panic", p.caption)
		}
	}()

	return p.check(ctx)
}
