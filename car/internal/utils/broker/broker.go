package broker

import "github.com/alserov/rently/car/internal/metrics"

type Broker struct {
	Addr    string
	Metrics metrics.MetricTopics
}
