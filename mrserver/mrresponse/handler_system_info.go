package mrresponse

import (
	"net/http"
	"os"
	"time"

	"github.com/mondegor/go-webcore/mrdebug"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver"
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

func HandlerGetSystemInfoAsJson(cfg SystemInfoConfig) (func(w http.ResponseWriter, r *http.Request), error) {
	hostName, err := os.Hostname()

	if err != nil {
		return nil, err
	}

	return mrserver.HandlerGetStructAsJson(
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
