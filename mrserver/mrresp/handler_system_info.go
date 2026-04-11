package mrresp

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/util/xio"
)

const (
	// processCacheTTL - время жизни кэша списка процессов.
	// Используется для снижения нагрузки при частых запросах информации о системе.
	processCacheTTL = 15 * time.Second
)

type (
	// SystemInfoConfig - конфигурация с информацией о запущенной системе.
	SystemInfoConfig struct {
		// Caption - название приложения в свободной форме.
		Caption string

		// Version - версия приложения.
		Version string

		// Environment - окружение (development, staging, production).
		Environment string

		// IsDebug - флаг режима отладки.
		IsDebug bool

		// LogLevel - текущий уровень логирования.
		LogLevel string

		// StartedAt - время запуска приложения.
		StartedAt time.Time

		// ProcessesFunc - функция получения списка текущих процессов.
		// Может возвращать nil, если информация о процессах не требуется.
		ProcessesFunc func(ctx context.Context) []SystemInfoProcess
	}

	// SystemInfoProcess - информация об отдельном процессе.
	// Используется в ответе эндпоинта /system/info.
	SystemInfoProcess struct {
		// Caption - название процесса в свободной форме.
		Caption string `json:"name"`

		// Status - состояние процесса (например: "running", "stopped").
		Status string `json:"status"`
	}
)

type (
	// systemInfoResponse - внутренний формат ответа с информацией о системе.
	systemInfoResponse struct {
		Caption     string              `json:"name"`
		Version     string              `json:"version"`
		Environment string              `json:"environment"`
		HostName    string              `json:"host_name"`
		IsDebug     bool                `json:"is_debug"`
		LogLevel    string              `json:"log_level"`
		StartedAt   string              `json:"started_at"`
		Processes   []SystemInfoProcess `json:"processes"`
	}
)

// HandlerGetSystemInfoAsJSON - создаёт обработчик для формирования информации о запущенной системе.
func HandlerGetSystemInfoAsJSON(logger mrlog.Logger, cfg SystemInfoConfig) (http.HandlerFunc, error) {
	hostName, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	staticResponse := systemInfoResponse{
		Caption:     cfg.Caption,
		Version:     cfg.Version,
		Environment: cfg.Environment,
		HostName:    hostName,
		IsDebug:     cfg.IsDebug,
		LogLevel:    cfg.LogLevel,
		StartedAt:   cfg.StartedAt.Format(time.RFC3339Nano),
	}

	processCache := newCachedProcesses(cfg.ProcessesFunc)

	return func(w http.ResponseWriter, r *http.Request) {
		status := http.StatusOK
		ctx := r.Context()

		response := staticResponse
		response.Processes = processCache.get(ctx)

		bytes, err := json.Marshal(response)
		if err != nil {
			status = http.StatusUnprocessableEntity
			bytes = nil

			logger.Error(r.Context(), "marshal failed", "error", errors.ErrInternalHttpResponseParseData.Wrap(err))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		xio.Write(r.Context(), logger, w, bytes)
	}, nil
}

type (
	// cachedProcesses - кэшированный список процессов с TTL.
	cachedProcesses struct {
		mu         sync.Mutex
		fn         func(ctx context.Context) []SystemInfoProcess
		cachedData []SystemInfoProcess
		expiresAt  time.Time
	}
)

func newCachedProcesses(fn func(ctx context.Context) []SystemInfoProcess) *cachedProcesses {
	return &cachedProcesses{
		fn: fn,
	}
}

func (c *cachedProcesses) get(ctx context.Context) []SystemInfoProcess {
	c.mu.Lock()
	defer c.mu.Unlock()

	if time.Now().After(c.expiresAt) || c.cachedData == nil {
		c.cachedData = c.fn(ctx)
		c.expiresAt = time.Now().Add(processCacheTTL)
	}

	return c.cachedData
}
