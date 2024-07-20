package mrprometheus

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver"
)

type (
	// ObserveRequest - comment struct.
	ObserveRequest struct {
		requestStatus *prometheus.HistogramVec
		requestSize   *prometheus.CounterVec
		responseSize  *prometheus.CounterVec
	}
)

// NewObserveRequest - создаёт объект StatRequest.
func NewObserveRequest(namespace string) *ObserveRequest {
	return &ObserveRequest{
		requestStatus: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: "request",
				Name:      "duration_seconds",
				Help:      "Request executed time",
				Buckets:   []float64{0.001, 0.005, 0.025, 0.2, 1, 4, 8},
			},
			[]string{"status"},
		),
		requestSize: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "request",
				Name:      "size_bytes_total",
				Help:      "Size in bytes of received information",
			},
			[]string{"method"},
		),
		responseSize: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "response",
				Name:      "size_bytes_total",
				Help:      "Size in bytes of sent information",
			},
			[]string{"method"},
		),
	}
}

// Collectors - comment method.
func (r *ObserveRequest) Collectors() []prometheus.Collector {
	return []prometheus.Collector{
		r.requestStatus,
		r.requestSize,
		r.responseSize,
	}
}

// SendMetrics - comment method.
func (r *ObserveRequest) SendMetrics(_ mrlog.Logger, start time.Time, sr *mrserver.StatRequest, sw *mrserver.StatResponseWriter) {
	r.requestStatus.WithLabelValues(strconv.Itoa(sw.StatusCode())).Observe(time.Since(start).Seconds())
	r.requestSize.WithLabelValues(sr.Request().Method).Add(float64(sr.Bytes()))
	r.responseSize.WithLabelValues(sr.Request().Method).Add(float64(sw.Bytes()))
}
