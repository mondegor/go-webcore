package mrresp

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlib"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver"
)

type (
	statInfoResponse struct {
		RequestCount   uint64 `json:"requestCount"`
		RequestSize    uint64 `json:"requestSize"`
		ResponseSize   uint64 `json:"responseSize"`
		AccessLastTime string `json:"accessLastTime"`
	}
)

var (
	statMutex      = &sync.Mutex{}
	requestCount   uint64
	requestSize    uint64
	responseSize   uint64
	accessLastTime time.Time
)

// HandlerGetStatInfoAsJSON - comment func.
func HandlerGetStatInfoAsJSON(unexpectedStatus int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		statMutex.Lock()
		response := statInfoResponse{
			RequestCount:   requestCount,
			RequestSize:    requestSize,
			ResponseSize:   responseSize,
			AccessLastTime: accessLastTime.Format(time.RFC3339),
		}
		statMutex.Unlock()

		status := http.StatusOK

		bytes, err := json.Marshal(response)
		if err != nil {
			status = unexpectedStatus
			bytes = nil

			mrlog.Ctx(r.Context()).Error().Err(mrcore.ErrHttpResponseParseData.Wrap(err)).Msg("marshal failed")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		mrlib.Write(r.Context(), w, bytes)
	}
}

// ApplyStatRequest - comment func.
func ApplyStatRequest(l mrlog.Logger, start time.Time, sr *mrserver.StatRequest, sw *mrserver.StatResponseWriter) {
	r := sr.Request()

	l.Trace().
		Str("method", r.Method).
		Str("url", r.RequestURI).
		Str("remoteAddr", r.RemoteAddr).
		Str("userAgent", r.UserAgent()).
		Int("status", sw.StatusCode()).
		Int("requestSize", sr.Bytes()).
		Int("size", sw.Bytes()).
		Int("elapsed_μs", int(time.Since(start).Microseconds())).
		Msg("incoming request")

	statMutex.Lock()
	requestCount++
	requestSize += uint64(sr.Bytes())
	responseSize += uint64(sw.Bytes())
	accessLastTime = time.Now().UTC()
	statMutex.Unlock()
}
