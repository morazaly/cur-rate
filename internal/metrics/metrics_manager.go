package metrics

import (
	"context"
	"net"
	"net/http"

	"log/slog"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	log  *slog.Logger
	addr string
}

func NewMetrics(log *slog.Logger, addr string) *Metrics {
	return &Metrics{
		log:  log,
		addr: addr,
	}
}

func (m Metrics) Start(_ context.Context) error {
	listener, err := net.Listen("tcp", m.addr)
	if err != nil {
		return err
	}

	httpMux := http.NewServeMux()
	httpMux.Handle("/metrics", promhttp.Handler())

	m.log.Info("Metrics start on ", "MetricPort", m.addr)

	return http.Serve(listener, httpMux)
}
