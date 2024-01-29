package workers

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Worker interface {
	MustStart()
	Metrics() []prometheus.Collector
}
