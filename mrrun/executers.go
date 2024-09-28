package mrrun

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"

	"github.com/mondegor/go-webcore/mrlog"
)

const (
	signalChanLen = 10
)

type (
	// Executer - структура выполнения процесса с возможностью его
	// остановки и отправки сообщения, о запуске этого процесса.
	Executer struct {
		Execute   func() error
		Interrupt func(error)
		StartedOk chan struct{} // OPTIONAL
	}
)

// MakeSignalHandler - возвращает контекст, в котором установлена его отмена при перехвате системного события.
// Это необходимо для корректной (graceful) остановки приложения. Также возвращаются функция для отслеживания
// сигналов системы, а также функция корректного прекращения отслеживания сигналов системы.
func MakeSignalHandler(ctx context.Context) (context.Context, Executer) {
	ctx, cancel := context.WithCancel(ctx)

	once := sync.Once{}
	chFirst := make(chan struct{})
	chFirstCloser := func() {
		once.Do(func() {
			close(chFirst)
		})
	}

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

	return ctx, Executer{
		Execute: func() error {
			chFirstCloser()

			select {
			case signalApp := <-signalStop:
				logger.Info().Msgf("Shutting down the application by signal: " + signalApp.String())

				return nil
			case <-ctx.Done():
				logger.Info().Msgf("Shutting down the application by a neighboring process")

				return ctx.Err()
			}
		},
		Interrupt: func(_ error) {
			logger.Info().Msg("Cancel the main context of the application")

			signal.Stop(signalStop)
			chFirstCloser()
			cancel()
			close(signalStop)
		},
		StartedOk: chFirst,
	}
}

// MakeExecuter - возвращает функции запуска и остановки процесса.
// Запуск этого процесса не зависит от других процессов.
func MakeExecuter(ctx context.Context, process Process) Executer {
	return Executer{
		Execute: func() error {
			return process.Start(ctx, nil)
		},
		Interrupt: func(_ error) {
			logger := mrlog.Ctx(ctx)

			shutdownCtx := logger.WithContext(context.Background())
			if err := process.Shutdown(shutdownCtx); err != nil {
				logger.Error().Err(err).Send()
			}
		},
	}
}

// MakeNextExecuter - возвращает канал, по которому будет передано событие, что процесс запущен.
// Также возвращает функции запуска и остановки этого процесса.
// Запуск процесса будет осуществлён только при получении события по каналу chPrev (если канал указан).
func MakeNextExecuter(ctx context.Context, process Process, chPrev chan struct{}) Executer {
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
			if chPrev != nil {
				<-chPrev
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

			shutdownCtx := mrlog.Ctx(ctx).WithContext(context.Background())
			if err := process.Shutdown(shutdownCtx); err != nil {
				mrlog.Ctx(ctx).Error().Err(err).Str("process", process.Caption()).Send()
			}
		},
		StartedOk: chNext,
	}
}
