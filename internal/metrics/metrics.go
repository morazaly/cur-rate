package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	metricDurName = "cur_rate_postgres_duration_seconds"
	metricDurDesc = "Histogram for storing duration of postgres handling"
)

var (
	DurPGProcessed = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:        metricDurName,
		Help:        metricDurDesc,
		ConstLabels: map[string]string{"app_name": "cur_rate"},
		Buckets:     prometheus.DefBuckets,
	}, []string{
		"type",
	})
)
