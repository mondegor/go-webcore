package mrrun

import (
	"context"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"

	"github.com/mondegor/go-webcore/mrlog"
)

const (
	signalChanLen = 10
)

// MakeSignalHandler - возвращает контекст, в котором установлена его отмена при перехвате системного события.
// Это необходимо для корректной (graceful) остановки приложения. Также возвращаются функция для отслеживания
// сигналов системы, а также функция корректного прекращения отслеживания сигналов системы.
func MakeSignalHandler(ctx context.Context) (ctxWithCancel context.Context, execute func() error, interrupt func(error)) {
	ctx, cancel := context.WithCancel(ctx)

	signalStop := make(chan os.Signal, signalChanLen)

	signal.Notify(
		signalStop,
		syscall.SIGABRT,
		syscall.SIGQUIT,
		syscall.SIGHUP,
		os.Interrupt,
		syscall.SIGTERM,
	)

	logger := mrlog.Ctx(ctx).With().Str("process", "SignalHandler").Logger()

	return ctx,
		func() error {
			logger.Info().Msg("Waiting for requests. To exit press CTRL+C")

			select {
			case signalApp := <-signalStop:
				logger.Info().Msgf("Shutting down the application by signal: " + signalApp.String())

				return nil
			case <-ctx.Done():
				logger.Info().Msgf("Shutting down the application by a neighboring process")

				return ctx.Err()
			}
		}, func(_ error) {
			logger.Info().Msg("Cancel the main context of the application")

			signal.Stop(signalStop)
			cancel()
			close(signalStop)
		}
}

// MakeExecuter - возвращает функции запуска и остановки процесса.
// Запуск этого процесса не зависит от других процессов.
func MakeExecuter(ctx context.Context, process Process) (execute func() error, interrupt func(error)) {
	return func() error {
			return process.Start(ctx, nil)
		}, func(_ error) {
			logger := mrlog.Ctx(ctx)

			shutdownCtx := logger.WithContext(context.Background())
			if err := process.Shutdown(shutdownCtx); err != nil {
				logger.Error().Err(err).Send()
			}
		}
}

// MakeNextExecuter - возвращает канал, по которому будет передано событие, что процесс запущен.
// Также возвращает функции запуска и остановки этого процесса.
// Запуск процесса будет осуществлён только при получении события по каналу chPrev (если канал указан).
func MakeNextExecuter(ctx context.Context, process Process, chPrev chan struct{}) (chNext chan struct{}, execute func() error, interrupt func(error)) {
	chNext = make(chan struct{})
	isStarted := atomic.Bool{}
	closer := func() {
		select {
		case <-chNext:
		default:
			close(chNext)
		}
	}

	return chNext,
		func() error {
			if chPrev != nil {
				<-chPrev
			}

			isStarted.Store(true)

			if err := ctx.Err(); err != nil {
				// ignore errors from other processes (context canceled)
				return nil //nolint:nilerr
			}

			return process.Start(ctx, closer)
		}, func(_ error) {
			closer()

			if !isStarted.Load() {
				return
			}

			shutdownCtx := mrlog.Ctx(ctx).WithContext(context.Background())
			if err := process.Shutdown(shutdownCtx); err != nil {
				mrlog.Ctx(ctx).Error().Err(err).Str("process", process.Caption()).Send()
			}
		}
}
