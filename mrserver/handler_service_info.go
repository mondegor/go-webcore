package mrserver

import (
	"net/http"
	"os"
	"time"

	"github.com/mondegor/go-webcore/mrcore"
)

type (
	ConfigServiceInfo struct {
		Name      string
		Version   string
		StartedAt time.Time
	}

	serviceInfoResponse struct {
		Name      string `json:"name"`
		Version   string `json:"version"`
		HostName  string `json:"hostName"`
		IsDebug   bool   `json:"debug"`
		StartedAt string `json:"startedAt"`
	}
)

func HandlerGetServiceInfoAsJson(cfg ConfigServiceInfo) (func(w http.ResponseWriter, r *http.Request), error) {
	hostName, err := os.Hostname()

	if err != nil {
		return nil, err
	}

	return HandlerGetStructAsJson(
		serviceInfoResponse{
			Name:      cfg.Name,
			Version:   cfg.Version,
			HostName:  hostName,
			IsDebug:   mrcore.Debug(),
			StartedAt: cfg.StartedAt.Format(time.RFC3339Nano),
		},
		http.StatusOK,
	)
}
