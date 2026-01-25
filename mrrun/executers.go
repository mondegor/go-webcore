package mrrun

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/mrtrace"
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
	ctx = r.traceManager.WithGeneratedProcessID(ctx, mrtrace.KeyProcessID)
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
			if err := process.Shutdown(ctx); err != nil {
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
					close(chCurrentReady)

					return errors.ErrSystemTimeoutPeriodHasExpired.WithDetails(
						"the waiting time for the previous process has expired",
						"previousProcess", prev.Caption,
						"process", process.Caption(),
					)
				case <-prev.ready:
				}
			}

			if err := ctx.Err(); err != nil {
				close(chCurrentReady)

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
			// ожидается готовность функции Execute
			select {
			case <-time.NewTimer(process.ReadyTimeout() + prev.readyTimeout).C:
				r.logger.Error(
					ctx,
					"the waiting time to interrupt the process has expired",
					"process", process.Caption(),
				)
			case <-chCurrentReady:
			}

			if !isStartCalled.Load() {
				return
			}

			if err := process.Shutdown(ctx); err != nil {
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
