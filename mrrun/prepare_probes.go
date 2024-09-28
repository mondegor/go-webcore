package mrrun

import (
	"context"
	"errors"
	"net/http"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlog"
)

type (
	// ProbeChecker - интерфейс проверки какой-либо пробы.
	ProbeChecker interface {
		Caption() string
		Check(ctx context.Context) error
	}

	// FinishedProbe - проба, которая была проведена и содержит статус выполнения.
	FinishedProbe struct {
		Caption string
		Status  int
	}
)

// PrepareProbesForCheck - возвращает функцию с заряженными пробами, для проверки работоспособности сервиса.
// Если хотя бы одна проба не завершится успешно, то возвращаемая функция вернёт false.
func PrepareProbesForCheck(probes ...ProbeChecker) func(ctx context.Context) bool {
	return func(ctx context.Context) bool {
		for _, probe := range probes {
			if err := probe.Check(ctx); err != nil {
				mrlog.Ctx(ctx).Error().Err(err).Send()

				return false
			}
		}

		return true
	}
}

// PrepareProbes - возвращает функцию с заряженными пробами, для проверки работоспособности сервиса.
// Сама возвращаемая функция возвращает список проведённых проб с их статусами выполнения.
func PrepareProbes(probes ...ProbeChecker) func(ctx context.Context) []FinishedProbe {
	return func(ctx context.Context) []FinishedProbe {
		info := make([]FinishedProbe, len(probes))

		for i, probe := range probes {
			status := http.StatusOK

			if err := probe.Check(ctx); err != nil {
				mrlog.Ctx(ctx).Error().Err(err).Send()

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

// WithAppReadyProbe - возвращает пробу готовности приложения к приёму запросов.
func WithAppReadyProbe(app *AppHealth) func(ctx context.Context) error {
	return func(_ context.Context) error {
		if app.IsReady() {
			return nil
		}

		return mrcore.ErrUseCaseTemporarilyUnavailable.Wrap(errors.New("app is not ready"))
	}
}
