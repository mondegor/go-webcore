package mrresp

import (
	"net/http"
	"os"
	"time"

	"github.com/mondegor/go-webcore/mrdebug"
	"github.com/mondegor/go-webcore/mrlog"
)

type (
	SystemInfoConfig struct {
		Name      string
		Version   string
		StartedAt time.Time
	}

	systemInfoResponse struct {
		Name      string `json:"name"`
		Version   string `json:"version"`
		HostName  string `json:"hostName"`
		IsDebug   bool   `json:"isDebug"`
		LogLevel  string `json:"logLevel"`
		StartedAt string `json:"startedAt"`
	}
)

func HandlerGetSystemInfoAsJSON(cfg SystemInfoConfig) (http.HandlerFunc, error) {
	hostName, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	return HandlerGetStructAsJSON(
		systemInfoResponse{
			Name:      cfg.Name,
			Version:   cfg.Version,
			HostName:  hostName,
			IsDebug:   mrdebug.IsDebug(),
			LogLevel:  mrlog.Default().Level().String(),
			StartedAt: cfg.StartedAt.Format(time.RFC3339Nano),
		},
		http.StatusOK,
	)
}
