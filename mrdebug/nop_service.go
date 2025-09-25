package mrdebug

import (
	"context"
	"fmt"
	"time"

	"github.com/mondegor/go-sysmess/mrlog"
)

type (
	nopServiceError struct {
		name string
		err  error
	}
)

// Error - comment method.
func (e *nopServiceError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("NopService '%s' has error: %s", e.name, e.err.Error())
	}

	return fmt.Sprintf("NopService '%s' has no error", e.name)
}

// PrepareNopServiceWithTimeoutToStart - предназначен для тестирования запуска и остановки процессов
// Пример: appRunner.Add(mrdebug.PrepareNopServiceWithTimeoutToStart(ctx, "s1", 5 * time.Second)).
func PrepareNopServiceWithTimeoutToStart(
	ctx context.Context,
	logger mrlog.Logger,
	name string,
	expiry time.Duration,
) (execute func() error, interrupt func(error)) {
	ctx, cancel := context.WithTimeout(ctx, expiry)

	return func() error {
			logger.Info(ctx, fmt.Sprintf("Running the NopService '%s' with timeout", name))
			<-ctx.Done()

			err := ctx.Err()

			// если такого типа процесс успел первым
			if _, ok := err.(*nopServiceError); !ok { //nolint:errorlint
				err = &nopServiceError{
					name: name,
					err:  err,
				}
			}

			return err
		}, func(err error) {
			cancel()

			if errService, ok := err.(*nopServiceError); ok { //nolint:errorlint
				if errService.name == name {
					logger.Info(ctx, fmt.Sprintf("Shutting down the NopService '%s' by timeout", name))
				} else {
					logger.Info(ctx, fmt.Sprintf("Shutting down the NopService '%s' by process '%s'", name, errService.name))
				}
			} else {
				logger.Info(ctx, fmt.Sprintf("Shutting down the NopService '%s' by another process", name))
			}
		}
}
