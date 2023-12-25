package metrics

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
)

type Metric interface {
	Get() prometheus.Collector
	Run(ctx context.Context)
}

func Register(reg prometheus.Registerer, m ...[]prometheus.Collector) {
	var c []prometheus.Collector
	for _, coll := range m {
		c = append(c, coll...)
	}

	reg.MustRegister(c...)
}
