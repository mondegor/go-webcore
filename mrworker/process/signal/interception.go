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
		cancel     context.CancelFunc
		signalStop chan os.Signal
		logger     mrlog.Logger
		wgMain     sync.WaitGroup
	}
)

// NewInterceptor - создаёт объект Interceptor и возвращает контекст,
// в котором установлена его отмена при перехвате системного события.
func NewInterceptor(ctx context.Context, logger mrlog.Logger) (context.Context, *Interceptor) {
	ctx, cancel := context.WithCancel(ctx)
	signalStop := make(chan os.Signal, signalBufferLen)

	signal.Notify(
		signalStop,
		syscall.SIGABRT,
		syscall.SIGQUIT,
		syscall.SIGHUP,
		os.Interrupt,
		syscall.SIGTERM,
	)

	return ctx, &Interceptor{
		cancel:     cancel,
		signalStop: signalStop,
		logger:     logger,
		wgMain:     sync.WaitGroup{},
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
// Повторный запуск метода одно и того же объекта не предусмотрен, даже после вызова Shutdown.
func (p *Interceptor) Start(ctx context.Context, ready func()) error {
	p.wgMain.Add(1)
	defer p.wgMain.Done()

	p.logger.Debug(ctx, "Starting the signal interceptor...")
	defer p.logger.Debug(ctx, "The signal interceptor has been stopped")

	if ready != nil {
		ready()
	}

	select {
	case signalApp := <-p.signalStop:
		p.logger.Info(ctx, "Interceptor detected an interrupting signal", "value", signalApp.String())
	case <-ctx.Done():
		p.logger.Info(ctx, "Interceptor detected context signal 'cancel'")

		return ctx.Err()
	}

	return nil
}

// Shutdown - корректная остановка сервиса перехвата системных событий.
func (p *Interceptor) Shutdown(ctx context.Context) error {
	p.logger.Debug(ctx, "Shutting down the signal interceptor...")
	signal.Stop(p.signalStop)
	close(p.signalStop)

	p.wgMain.Wait()
	p.cancel()
	p.logger.Debug(ctx, "The signal interceptor has been shut down")

	return nil
}
