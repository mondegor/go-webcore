package mrtool

import (
    "context"
    "fmt"
    "io"
    "os"
    "os/signal"
    "syscall"

    "github.com/mondegor/go-webcore/mrcore"
)

type (
    AppHelper struct {
        Logger mrcore.Logger
    }
)

func NewAppHelper(logger mrcore.Logger) *AppHelper {
    return &AppHelper{
        Logger: logger,
    }
}

func (h *AppHelper) Close(c io.Closer) {
    if err := c.Close(); err != nil {
        h.Logger.Err(mrcore.FactoryErrInternalFailedToClose.Caller(1).Wrap(err, fmt.Sprintf("%#v", c)))
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
    h.Logger.Info("Application shutdown, signal: " + signalApp.String())
    cancel()
}

func (h *AppHelper) ExitOnError(err error) {
    if err != nil {
        h.Logger.Err(mrcore.FactoryErrInternal.Caller(1).Wrap(err))
        os.Exit(1)
    }
}
