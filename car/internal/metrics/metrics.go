package metrics

type Metrics interface {
	IncreaseActiveRentsAmount() error
	DecreaseActiveRentsAmount() error
}

func NewMetrics() Metrics {
	return &metrics{}
}

type metrics struct {
}

func (m *metrics) DecreaseActiveRentsAmount() error {
	//TODO implement me
	panic("implement me")
}

func (m *metrics) IncreaseActiveRentsAmount() error {
}
