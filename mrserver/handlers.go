package mrserver

import (
	"encoding/json"
	"net/http"
)

func HandlerGetHealth() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

func HandlerGetStatusOKAsJson() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{\"status\": \"OK\"}"))
	}
}

func HandlerGetStructAsJson(data any, status int) (func(w http.ResponseWriter, r *http.Request), error) {
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
