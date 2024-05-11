package mrresp

import (
	"encoding/json"
	"net/http"
	"time"
)

func HandlerGetMethodNotAllowedAsJSON() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)

		bytes, err := json.Marshal(
			ErrorDetailsResponse{
				Title:   "405 Method Not Allowed",
				Details: "The server knows the request method, but the target resource doesn't support this method",
				Request: r.URL.Path,
				Time:    time.Now().UTC().Format(time.RFC3339),
			},
		)

		if err == nil {
			w.Write(bytes)
		}
	}
}
