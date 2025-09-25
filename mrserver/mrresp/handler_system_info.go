package mrresp

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/mondegor/go-sysmess/mrerr/mr"
	"github.com/mondegor/go-sysmess/mrlog"

	"github.com/mondegor/go-webcore/mrlib"
)

type (
	// SystemInfoConfig - информация о запущенной системе.
	SystemInfoConfig struct {
		Name        string
		Version     string
		Environment string
		IsDebug     bool
		LogLevel    string
		StartedAt   time.Time
		Processes   func(ctx context.Context) map[string]string
	}

	systemInfoResponse struct {
		Name        string            `json:"name"`
		Version     string            `json:"version"`
		Environment string            `json:"environment"`
		HostName    string            `json:"hostName"`
		IsDebug     bool              `json:"isDebug"`
		LogLevel    string            `json:"logLevel"`
		StartedAt   string            `json:"startedAt"`
		Processes   map[string]string `json:"processes"`
	}
)

// HandlerGetSystemInfoAsJSON - возвращает обработчик для формирования информации о запущенной системе.
func HandlerGetSystemInfoAsJSON(logger mrlog.Logger, cfg SystemInfoConfig) (http.HandlerFunc, error) {
	hostName, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	staticResponse := systemInfoResponse{
		Name:        cfg.Name,
		Version:     cfg.Version,
		Environment: cfg.Environment,
		HostName:    hostName,
		IsDebug:     cfg.IsDebug,
		LogLevel:    cfg.LogLevel,
		StartedAt:   cfg.StartedAt.Format(time.RFC3339Nano),
	}

	return func(w http.ResponseWriter, r *http.Request) {
		status := http.StatusOK
		ctx := r.Context()

		response := staticResponse
		response.Processes = cfg.Processes(ctx)

		bytes, err := json.Marshal(response)
		if err != nil {
			status = http.StatusUnprocessableEntity
			bytes = nil

			logger.Error(r.Context(), "marshal failed", "error", mr.ErrHttpResponseParseData.Wrap(err))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		mrlib.Write(r.Context(), logger, w, bytes)
	}, nil
}
