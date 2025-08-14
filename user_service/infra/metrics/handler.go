package metrics

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	server *http.Server
}

var (
	RequestCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "service_requests_total",
		Help: "Total number of requests to exchange service",
	})
	ResponseTimeSeconds = promauto.NewSummary(prometheus.SummaryOpts{
		Name: "service_response_time_seconds",
		Help: "Response time of exchange service",
		Objectives: map[float64]float64{
			0.1:  0.1,
			0.25: 0.1,
			0.5:  0.05,
			0.75: 0.01,
			0.9:  0.01,
			0.99: 0.001,
		},
	})
	ErrorCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "service_errors_total",
		Help: "Total number of errors in exchange service",
	})
)

func NewMetricsApp() *Metrics {
	return &Metrics{
		server: &http.Server{
			Addr: ":9100",
		},
	}
}

func (m *Metrics) Run() {
	http.Handle("/metrics", promhttp.Handler())
	if err := m.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

func (m *Metrics) Stop() {
	if err := m.server.Shutdown(context.Background()); err != nil {
		slog.Info("Http server for metrics is closed")
	}
}
