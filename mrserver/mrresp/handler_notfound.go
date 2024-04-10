package mrresp

import (
	"encoding/json"
	"net/http"
	"time"
)

func HandlerGetNotFoundAsJson() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		bytes, err := json.Marshal(
			ErrorDetailsResponse{
				Title:   "404 Not Found",
				Details: "The server cannot find the requested resource",
				Request: r.URL.Path,
				Time:    time.Now().UTC().Format(time.RFC3339),
			},
		)

		if err == nil {
			w.Write(bytes)
		}
	}
}
