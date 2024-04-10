package mrresp

import (
	"encoding/json"
	"net/http"
	"time"
)

func HandlerGetFatalErrorAsJson() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusTeapot)

		bytes, err := json.Marshal(
			ErrorDetailsResponse{
				Title:   "Internal server error",
				Details: "The server encountered an unexpected condition that prevented it from fulfilling the request",
				Request: r.URL.Path,
				Time:    time.Now().UTC().Format(time.RFC3339),
			},
		)

		if err == nil {
			w.Write(bytes)
		}
	}
}
