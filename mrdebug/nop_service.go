package mrdebug

import (
	"context"
	"fmt"
	"time"

	"github.com/mondegor/go-sysmess/mrlog"
)

type (
	// nopServiceError - внутренняя ошибка сервиса NopService.
	// Используется для сигнализации о завершении сервиса по таймауту или внешней отмене.
	// Содержит имя сервиса и причину завершения.
	nopServiceError struct {
		name string
		err  error
	}
)

// Error - возвращает человеко-читаемое описание ошибки NopService.
func (e *nopServiceError) Error() string {
	if e.err != nil {
		return "NopService has error: name=" + e.name + ", " + e.err.Error()
	}

	return "NopService has no error: name=" + e.name
}

// PrepareNopServiceWithTimeoutToStart - создаёт тестовый сервис с таймаутом для проверки
// механизмов запуска и остановки процессов.
//
// Сервис ждёт истечения таймаута (expiry) или отмены контекста, после чего завершается
// с ошибкой nopServiceError. Используется для тестирования infrastructure-кода,
// например appRunner, без выполнения реальной работы.
//
// Возвращает две функции:
//   - execute - запускает сервис (блокируется до истечения таймаута или отмены);
//   - interrupt - принудительно останавливает сервис с указанной ошибкой;
//
// Пример использования:
//   - appRunner.Add(mrdebug.PrepareNopServiceWithTimeoutToStart(ctx, logger, "s1", 5*time.Second)).
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
