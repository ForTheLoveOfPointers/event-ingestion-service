package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metric struct {
	EventRate   *prometheus.CounterVec
	ReqDuration *prometheus.HistogramVec
}

const (
	StatusSuccess string = "success"
	StatusFailure string = "failure"
)

func NewMetric() *Metric {
	m := &Metric{
		EventRate: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "event_ingestion",
				Name:      "event_rate",
				Help:      "Number of ingested events.",
			},
			[]string{"status"},
		),
		ReqDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "event_ingestion",
				Name:      "req_duration_milli",
				Help:      "Duration of ingestion requests in ms.",
			},
			[]string{"status"},
		),
	}

	prometheus.MustRegister(m.EventRate, m.ReqDuration)
	return m
}
