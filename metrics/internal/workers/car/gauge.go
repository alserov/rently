package car

type GaugeTopics struct {
	ActiveCarRentsTopic string
}

func NewGaugeMetric(topics *GaugeTopics) *GaugeMetric {
	return &GaugeMetric{}
}

type GaugeMetric struct {
}

func (g *GaugeMetric) Run() {

}
