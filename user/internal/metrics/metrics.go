package metrics

type Metrics interface {
}

func NewMetrics() Metrics {
	return &metrics{}
}

type metrics struct {
}
