package mrrun

import (
	"context"
	"net/http"

	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/mrlog"
)

type (
	// ProbeChecker - интерфейс проверки работоспособности компонента (пробы).
	ProbeChecker interface {
		Caption() string
		Check(ctx context.Context) error
	}

	// FinishedProbe - результат выполнения пробы с её статусом.
	// Содержит название пробы и HTTP-код результата.
	FinishedProbe struct {
		Caption string
		Status  int
	}
)

// PrepareProbesForCheck - создаёт функцию для проверки работоспособности всех проб.
//
// Возвращаемая функция выполняет последовательную проверку всех проб.
// Возвращает true, если ВСЕ пробы завершились успешно.
// Возвращает false, если хотя бы одна проба завершилась с ошибкой.
func PrepareProbesForCheck(logger mrlog.Logger, probes ...ProbeChecker) func(ctx context.Context) bool {
	return func(ctx context.Context) bool {
		for _, probe := range probes {
			if err := probe.Check(ctx); err != nil {
				logger.Error(ctx, "PrepareProbesForCheck", "error", err)

				return false
			}
		}

		return true
	}
}

// PrepareProbes - создаёт функцию для детальной проверки работоспособности всех проб.
// В отличие от PrepareProbesForCheck, возвращает результаты всех проб,
// а не просто общее состояние.
//
// Возвращаемая функция:
//  1. Выполняет последовательную проверку всех проб;
//  2. Для каждой пробы записывает статус: 200 (OK) или 422 (UnprocessableEntity);
//  3. Возвращает срез FinishedProbe с результатами всех проверок;
func PrepareProbes(logger mrlog.Logger, probes ...ProbeChecker) func(ctx context.Context) []FinishedProbe {
	return func(ctx context.Context) []FinishedProbe {
		info := make([]FinishedProbe, len(probes))

		for i, probe := range probes {
			status := http.StatusOK

			if err := probe.Check(ctx); err != nil {
				logger.Error(ctx, "PrepareProbes", "error", err)

				status = http.StatusUnprocessableEntity
			}

			info[i] = FinishedProbe{
				Caption: probe.Caption(),
				Status:  status,
			}
		}

		return info
	}
}

// WithAppReadyProbe - создаёт пробу готовности приложения к приёму запросов.
// Проверяет состояние AppHealth через метод IsReady().
//
// Возвращает nil, если приложение готово.
// Возвращает ошибку ErrSystemServiceTemporarilyUnavailable, если приложение не готово.
func WithAppReadyProbe(app *AppHealth) func(ctx context.Context) error {
	return func(_ context.Context) error {
		if app.IsReady() {
			return nil
		}

		return errors.ErrSystemServiceTemporarilyUnavailable.WithDetails("app is not ready")
	}
}
