package mrrun

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/mondegor/go-sysmess/mrerr/mr"
)

type (
	// Executer - структура выполнения процесса с возможностью его
	// остановки и отправки сообщения, о запуске этого процесса.
	Executer struct {
		Execute      func() error
		Interrupt    func(error)
		Synchronizer ProcessSync // OPTIONAL
	}
)

func (r *AppRunner) contextWithProcessID(ctx context.Context, process Process) context.Context {
	ctx = r.traceManager.WithGeneratedProcessID(ctx)
	r.logger.Info(ctx, "Start new process", "process_name", process.Caption())

	return ctx
}

// makeExecuter - возвращает функции запуска и остановки процесса.
// Запуск этого процесса не зависит от других процессов.
func (r *AppRunner) makeExecuter(ctx context.Context, process Process) Executer {
	ctx = r.contextWithProcessID(ctx, process)

	return Executer{
		Execute: func() error {
			return process.Start(ctx, func() {})
		},
		Interrupt: func(_ error) {
			// WARNING: создаётся новый контекст без возможности внешней отмены
			// для того чтобы метод Shutdown гарантированно отработал
			// при этом внутри Shutdown следует организовать персональный таймаут
			if err := process.Shutdown(r.traceManager.NewContextWithIDs(ctx)); err != nil {
				r.logger.Error(ctx, "AppRunner.makeExecuter", "error", err)
			}
		},
	}
}

// makeNextExecuter - возвращает канал, по которому будет передано событие, что процесс запущен.
// Также возвращает функции запуска и остановки этого процесса.
// Запуск процесса будет осуществлён только при получении события по каналу chPrev (если канал указан).
func (r *AppRunner) makeNextExecuter(ctx context.Context, process Process, prev ProcessSync) Executer {
	ctx = r.contextWithProcessID(ctx, process)

	isStartCalled := atomic.Bool{}
	chCurrentReady := make(chan struct{})

	return Executer{
		Execute: func() error {
			if prev.ready != nil {
				select {
				case <-time.NewTimer(prev.readyTimeout).C:
					return mr.ErrInternalTimeoutPeriodHasExpired.Wrap(
						fmt.Errorf("the waiting time for the previous process '%s' has expired", prev.Caption),
						"process_name", process.Caption(),
					)
				case <-prev.ready:
				}
			}

			if err := ctx.Err(); err != nil {
				// ignore errors from other processes (context canceled)
				return nil //nolint:nilerr
			}

			isStartCalled.Store(true)

			return process.Start(
				ctx,
				func() {
					close(chCurrentReady)
				},
			)
		},
		Interrupt: func(_ error) {
			if !isStartCalled.Load() {
				return
			}

			// WARNING: создаётся новый контекст без возможности внешней отмены
			// для того чтобы метод Shutdown гарантированно отработал
			// при этом внутри Shutdown следует организовать персональный таймаут
			if err := process.Shutdown(r.traceManager.NewContextWithIDs(ctx)); err != nil {
				r.logger.Error(ctx, "AppRunner.makeNextExecuter", "error", err)
			}
		},
		Synchronizer: ProcessSync{
			Caption:      process.Caption(),
			readyTimeout: process.ReadyTimeout(),
			ready:        chCurrentReady,
		},
	}
}
