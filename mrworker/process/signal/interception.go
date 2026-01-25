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
	signalBufferLen     = 10
	processCaption      = "SignalInterceptor"
	processReadyTimeout = 30 * time.Second
)

type (
	// Interceptor - сервис перехвата системных событий с целью
	// корректной (graceful) остановки всего приложения.
	Interceptor struct {
		logger mrlog.Logger
		wg     sync.WaitGroup
		done   chan struct{}
	}
)

// NewInterceptor - создаёт объект Interceptor и возвращает контекст,
// в котором установлена его отмена при перехвате системного события.
func NewInterceptor(logger mrlog.Logger) *Interceptor {
	return &Interceptor{
		logger: logger,
		wg:     sync.WaitGroup{},
		done:   make(chan struct{}),
	}
}

// Caption - возвращает название сервиса.
func (p *Interceptor) Caption() string {
	return processCaption
}

// ReadyTimeout - возвращает максимальное время, за которое должен быть запущен сервис.
func (p *Interceptor) ReadyTimeout() time.Duration {
	return processReadyTimeout
}

// Start - запуск сервиса перехвата системных событий.
// Отмена внешнего контекста приведёт к аварийному завершению процесса,
// для корректной остановки следует использовать Shutdown.
// Повторный запуск метода одно и того же объекта не предусмотрен, даже после вызова Shutdown.
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

// Shutdown - корректная остановка сервиса перехвата системных событий.
// При повторном вызове метода произойдёт panic.
func (p *Interceptor) Shutdown(ctx context.Context) error {
	p.logger.Debug(ctx, "Shutting down the signal interceptor...")
	close(p.done)

	p.wg.Wait()
	p.logger.Debug(ctx, "The signal interceptor has been shut down")

	return nil
}
