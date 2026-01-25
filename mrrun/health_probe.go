package mrrun

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/mrlog"
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
		logger  mrlog.Logger
	}
)

// NewHealthProbe - создаёт объект AppRunner.
func NewHealthProbe(logger mrlog.Logger, caption string, check func(ctx context.Context) error, timeout time.Duration) *HealthProbe {
	if timeout == 0 {
		timeout = defaultTimeout
	}

	return &HealthProbe{
		caption: caption,
		check:   check,
		timeout: timeout,
		logger:  logger,
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
			p.logger.Error(
				ctx,
				"HealthProbe",
				"error",
				errors.ErrInternalCaughtPanic.New(
					"source", "probe: "+p.caption,
					"recover", rvr,
					"stack_trace", string(debug.Stack()),
				),
			)

			err = fmt.Errorf("probe has panic (caption='%s')", p.caption)
		}
	}()

	return p.check(ctx)
}
