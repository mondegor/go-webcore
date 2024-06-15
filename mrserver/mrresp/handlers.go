package mrresp

import (
	"encoding/json"
	"net/http"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlib"
)

// HandlerGetHealth  - comment func.
func HandlerGetHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

// HandlerGetStatusOkAsJSON  - comment func.
func HandlerGetStatusOkAsJSON() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		mrlib.Write(r.Context(), w, []byte("{\"status\": \"OK\"}"))
	}
}

// HandlerGetStructAsJSON  - comment func.
func HandlerGetStructAsJSON(data any, status int) (http.HandlerFunc, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, mrcore.ErrHttpResponseParseData.Wrap(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		mrlib.Write(r.Context(), w, bytes)
	}, nil
}
