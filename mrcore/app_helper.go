package mrcore

import (
    "context"
    "fmt"
    "io"
    "os"
    "os/signal"
    "syscall"
)

type (
	AppHelper struct {
        logger Logger
    }
)

func NewAppHelper(logger Logger) *AppHelper {
    return &AppHelper{
        logger: logger,
    }
}

func (h *AppHelper) ExitOnError(err error) {
    if err != nil {
        h.logger.Err(err)
        os.Exit(1)
    }
}

func (h *AppHelper) Close(c io.Closer) {
    err := c.Close()

    if err != nil {
        h.logger.Err(FactoryErrInternalFailedToClose.Caller(1).Wrap(err, fmt.Sprintf("%v", c)))
    }
}

func (h *AppHelper) GracefulShutdown(cancel context.CancelFunc) {
    signalAppChan := make(chan os.Signal, 1)

    signal.Notify(
        signalAppChan,
        syscall.SIGABRT,
        syscall.SIGQUIT,
        syscall.SIGHUP,
        os.Interrupt,
        syscall.SIGTERM,
    )

    signalApp := <-signalAppChan
    h.logger.Info("Application shutdown, signal: " + signalApp.String())
    cancel()
}
