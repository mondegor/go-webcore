package mrresp

import (
	"net/http"
	"os"
	"time"

	"github.com/mondegor/go-webcore/mrlog"
)

type (
	// SystemInfoConfig - comment struct.
	SystemInfoConfig struct {
		Name        string
		Version     string
		Environment string
		IsDebug     bool
		StartedAt   time.Time
	}

	systemInfoResponse struct {
		Name        string `json:"name"`
		Version     string `json:"version"`
		Environment string `json:"environment"`
		HostName    string `json:"hostName"`
		IsDebug     bool   `json:"isDebug"`
		LogLevel    string `json:"logLevel"`
		StartedAt   string `json:"startedAt"`
	}
)

// HandlerGetSystemInfoAsJSON - comment func.
func HandlerGetSystemInfoAsJSON(cfg SystemInfoConfig) (http.HandlerFunc, error) {
	hostName, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	return HandlerGetStructAsJSON(
		systemInfoResponse{
			Name:        cfg.Name,
			Version:     cfg.Version,
			Environment: cfg.Environment,
			HostName:    hostName,
			IsDebug:     cfg.IsDebug,
			LogLevel:    mrlog.Default().Level().String(),
			StartedAt:   cfg.StartedAt.Format(time.RFC3339Nano),
		},
		http.StatusOK,
	)
}
