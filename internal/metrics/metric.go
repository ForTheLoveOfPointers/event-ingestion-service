package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metric struct {
	EventRate   *prometheus.CounterVec
	ReqDuration *prometheus.HistogramVec
}

func NewMetric(reg prometheus.Registerer) *Metric {
	m := &Metric{
		EventRate: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "event_rate",
				Help: "Number of ingested events",
			},
			[]string{"events"},
		),
		ReqDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "ingest_request_duration_seconds",
				Help: "Duration of ingestion requests in seconds",
			},
			[]string{"seconds"},
		),
	}

	reg.MustRegister(m.EventRate, m.ReqDuration)
	return m
}
