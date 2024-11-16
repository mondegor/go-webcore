package signal

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mondegor/go-webcore/mrlog"
)

const (
	signalChanLen       = 10
	processCaption      = "SignalInterception"
	processReadyTimeout = 60 * time.Second
)

// Interception - сервис перехвата системных событий с целью
// корректной (graceful) остановки всего приложения.
type Interception struct {
	cancel     context.CancelFunc
	signalStop chan os.Signal
}

// NewInterception - создаёт объект Interception и возвращает контекст,
// в котором установлена его отмена при перехвате системного события.
func NewInterception(ctx context.Context) (context.Context, *Interception) {
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

	return ctx, &Interception{
		cancel:     cancel,
		signalStop: signalStop,
	}
}

// Caption - возвращает название сервиса.
func (p *Interception) Caption() string {
	return processCaption
}

// ReadyTimeout - возвращает максимальное время, за которое должен быть запущен сервис.
func (p *Interception) ReadyTimeout() time.Duration {
	return processReadyTimeout
}

// Start - запуск сервиса перехвата системных событий.
func (p *Interception) Start(ctx context.Context, ready func()) error {
	mrlog.Ctx(ctx).Info().Msg("Starting the application...")

	if ready != nil {
		ready()
	}

	select {
	case signalApp := <-p.signalStop:
		mrlog.Ctx(ctx).Info().Msgf("Shutting down the application by signal: " + signalApp.String())
	case <-ctx.Done():
		mrlog.Ctx(ctx).Info().Msg("Shutting down the application by a neighboring process")

		return ctx.Err()
	}

	return nil
}

// Shutdown - корректная остановка сервиса перехвата системных событий.
func (p *Interception) Shutdown(ctx context.Context) error {
	mrlog.Ctx(ctx).Debug().Msg("Cancel the main context of the application")

	signal.Stop(p.signalStop)
	p.cancel()
	close(p.signalStop)

	return nil
}
