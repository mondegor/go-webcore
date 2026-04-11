package signal

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/mondegor/go-sysmess/mrlog"
)

const (
	// signalBufferLen - размер буфера канала сигналов.
	signalBufferLen = 10

	// processCaption - название сервиса.
	processCaption = "SignalInterceptor"

	// processReadyTimeout - таймаут готовности сервиса по умолчанию.
	processReadyTimeout = 30 * time.Second
)

type (
	// Interceptor - сервис перехвата системных сигналов для graceful shutdown.
	//
	// Принцип работы:
	//  1. Подписывается на системные сигналы (SIGABRT, SIGQUIT, SIGHUP, SIGINT, SIGTERM);
	//  2. Блокируется до получения сигнала или отмены контекста;
	//  3. При получении сигнала завершает работу, позволяя приложению корректно остановиться;
	//
	// Используется как первый процесс в цепочке процессов для управления жизненным циклом сервиса.
	Interceptor struct {
		logger mrlog.Logger
		wg     sync.WaitGroup
		done   chan struct{}
	}
)

// NewInterceptor - создаёт сервис перехвата системных сигналов.
func NewInterceptor(logger mrlog.Logger) *Interceptor {
	return &Interceptor{
		logger: logger,
		wg:     sync.WaitGroup{},
		done:   make(chan struct{}),
	}
}

// Caption - возвращает название сервиса в свободной форме.
func (p *Interceptor) Caption() string {
	return processCaption
}

// ReadyTimeout - возвращает максимальное время запуска сервиса.
func (p *Interceptor) ReadyTimeout() time.Duration {
	return processReadyTimeout
}

// Start - запуск сервиса перехвата системных сигналов.
//
// Процесс работы:
//  1. Подписывается на сигналы: SIGABRT, SIGQUIT, SIGHUP, os.Interrupt, SIGTERM;
//  2. Блокируется до получения сигнала, отмены контекста или вызова Shutdown;
//  3. При получении сигнала логирует его и завершает работу;
//
// Важно:
//   - Отмена внешнего контекста приведёт к завершению процесса;
//   - Для корректной остановки используйте Shutdown;
//   - Повторный запуск того же объекта не поддерживается.
func (p *Interceptor) Start(ctx context.Context, ready func()) error {
	p.wg.Add(1)
	defer p.wg.Done()

	signalStop := make(chan os.Signal, signalBufferLen)

	signal.Notify(
		signalStop,
		syscall.SIGABRT,
		syscall.SIGQUIT,
		syscall.SIGHUP,
		os.Interrupt,
		syscall.SIGTERM,
	)

	p.logger.Debug(ctx, "Starting the signal interceptor...")

	defer func() {
		signal.Stop(signalStop)
		close(signalStop)
		p.logger.Debug(ctx, "The signal interceptor has been stopped")
	}()

	if ready != nil {
		ready()
	}

	select {
	case <-p.done:
	case <-ctx.Done():
		p.logger.Debug(ctx, "The signal interceptor detected context 'Done'", "error", ctx.Err())
	case signalApp := <-signalStop:
		p.logger.Debug(ctx, "The signal interceptor detected an interrupting signal", "value", signalApp.String())
	}

	return nil
}

// Shutdown - корректная остановка сервиса перехвата сигналов.
// Отменяет подписку на сигналы и завершает работу.
//
// Важно: при повторном вызове произойдёт panic (закрытие закрытого канала done).
func (p *Interceptor) Shutdown(ctx context.Context) error {
	p.logger.Debug(ctx, "Shutting down the signal interceptor...")
	close(p.done)

	p.wg.Wait()
	p.logger.Debug(ctx, "The signal interceptor has been shut down")

	return nil
}
