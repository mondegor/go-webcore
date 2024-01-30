package mrdebug

import (
	"context"
	"fmt"
	"time"

	"github.com/mondegor/go-webcore/mrlog"
)

type (
	errorNopService struct {
		name string
		err  error
	}
)

func (e *errorNopService) Error() string {
	return fmt.Errorf("NopService '%s' has error: %w", e.name, e.err).Error()
}

// PrepareNopServiceWithTimeoutToStart - предназначен для тестирования запуска и остановки процессов
// пример: appRunner.Add(mrdebug.PrepareNopServiceWithTimeoutToStart(ctx, "s1", 5 * time.Second))
func PrepareNopServiceWithTimeoutToStart(
	ctx context.Context,
	name string,
	expiry time.Duration,
) (execute func() error, interrupt func(error)) {
	ctx, cancel := context.WithTimeout(ctx, expiry)
	logger := mrlog.Ctx(ctx)

	return func() error {
			logger.Info().Msgf("Running the NopService '%s' with timeout", name)
			<-ctx.Done()

			err := ctx.Err()

			// если такого типа процесс успел первым
			if _, ok := err.(*errorNopService); !ok {
				err = &errorNopService{name: name, err: err}
			}

			return err
		}, func(err error) {
			cancel()

			if errService, ok := err.(*errorNopService); ok {
				if errService.name == name {
					logger.Info().Msgf("Shutting down the NopService '%s' by timeout", name)
				} else {
					logger.Info().Msgf("Shutting down the NopService '%s' by process '%s'", name, errService.name)
				}
			} else {
				logger.Info().Msgf("Shutting down the NopService '%s' by another process", name)
			}
		}
}
