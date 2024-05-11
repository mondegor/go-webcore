package mrresp

import (
	"encoding/json"
	"net/http"
)

func HandlerGetHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

func HandlerGetStatusOKAsJSON() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{\"status\": \"OK\"}"))
	}
}

func HandlerGetStructAsJSON(data any, status int) (http.HandlerFunc, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write(bytes)
	}, nil
}
