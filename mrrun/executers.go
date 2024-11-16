package mrrun

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rs/xid"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrcore/mrapp"
	"github.com/mondegor/go-webcore/mrlog"
)

type (
	// Executer - структура выполнения процесса с возможностью его
	// остановки и отправки сообщения, о запуске этого процесса.
	Executer struct {
		Execute   func() error
		Interrupt func(error)
		Starting  StartingProcess // OPTIONAL
	}
)

// MakeExecuter - возвращает функции запуска и остановки процесса.
// Запуск этого процесса не зависит от других процессов.
func MakeExecuter(ctx context.Context, process Process) Executer {
	ctx = contextWithProcessID(ctx, process)

	return Executer{
		Execute: func() error {
			return process.Start(ctx, nil)
		},
		Interrupt: func(_ error) {
			// передаётся чистый контекст для исключения внешнего воздействия
			// при этом внутри Shutdown следует организовать персональный таймаут
			shutdownCtx := mrlog.Ctx(ctx).WithContext(context.Background())
			if err := process.Shutdown(shutdownCtx); err != nil {
				mrlog.Ctx(ctx).Error().Err(err).Send()
			}
		},
	}
}

// MakeNextExecuter - возвращает канал, по которому будет передано событие, что процесс запущен.
// Также возвращает функции запуска и остановки этого процесса.
// Запуск процесса будет осуществлён только при получении события по каналу chPrev (если канал указан).
func MakeNextExecuter(ctx context.Context, process Process, prev StartingProcess) Executer {
	ctx = contextWithProcessID(ctx, process)

	isStarted := atomic.Bool{}

	once := sync.Once{}
	chNext := make(chan struct{})
	chNextCloser := func() {
		once.Do(func() {
			close(chNext)
		})
	}

	return Executer{
		Execute: func() error {
			timeout := time.NewTimer(prev.ReadyTimeout)

			if prev.Ready != nil {
				select {
				case <-timeout.C:
					return mrcore.ErrInternalTimeoutPeriodHasExpired.Wrap(
						fmt.Errorf("the waiting time for the previous process '%s' has expired", prev.Caption),
					).WithAttr("process", process.Caption())
				case <-prev.Ready:
				}
			}

			isStarted.Store(true)

			if err := ctx.Err(); err != nil {
				// ignore errors from other processes (context canceled)
				return nil //nolint:nilerr
			}

			return process.Start(ctx, chNextCloser)
		},
		Interrupt: func(_ error) {
			chNextCloser()

			if !isStarted.Load() {
				return
			}

			// передаётся чистый контекст для исключения внешнего воздействия
			// при этом внутри Shutdown следует организовать персональный таймаут
			shutdownCtx := mrlog.Ctx(ctx).WithContext(context.Background())
			if err := process.Shutdown(shutdownCtx); err != nil {
				mrlog.Ctx(ctx).Error().Err(err).Send()
			}
		},
		Starting: StartingProcess{
			Caption:      process.Caption(),
			ReadyTimeout: process.ReadyTimeout(),
			Ready:        chNext,
		},
	}
}

func contextWithProcessID(ctx context.Context, process Process) context.Context {
	processID := xid.New().String()
	logger := mrlog.Ctx(ctx).With().Str("process", process.Caption()).Str(mrapp.KeyProcessID, processID).Logger()

	return mrlog.WithContext(mrapp.WithProcessContext(ctx, processID), logger)
}
