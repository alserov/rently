package workers

import (
	"encoding/json"
	"github.com/alserov/rently/metrics/internal/config"
	"github.com/alserov/rently/metrics/internal/log"
	"github.com/alserov/rently/metrics/internal/utils/broker"
	"github.com/prometheus/client_golang/prometheus"
	"log/slog"
	"time"
)

func NewCarsharingWorker(p config.Broker) Worker {
	return &worker{
		topics:     p.Topics,
		brokerAddr: p.Addr,
		responseTime: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "carsharing",
			Name:      "response_time",
			Help:      "time it takes to send response",
		}, []string{"request_type", "operation"}),
	}
}

type worker struct {
	log log.Logger

	brokerAddr string
	topics     config.Topics

	responseTime *prometheus.GaugeVec
}

func (w *worker) Metrics() []prometheus.Collector {
	return []prometheus.Collector{w.responseTime}
}

func (w *worker) MustStart() {
	go w.responseTimeMetric()
	select {}
}

func (w *worker) responseTimeMetric() {
	c := broker.NewConsumer(w.brokerAddr, w.topics.CarsharingResponseTime)

	msgs, err := c.Subscribe(w.topics.CarsharingResponseTime)
	if err != nil {
		panic("failed to subscribe to queue: " + err.Error())
	}

	for m := range msgs {
		var data ResponseTimeData
		if err = json.Unmarshal(m.Body, &data); err != nil {
			w.log.Error("failed to unmarshal message", slog.String("error", err.Error()))
		}

		w.responseTime.With(prometheus.Labels{"request_type": data.RequestType, "operation": data.Operation}).Set(float64(data.Duration.Milliseconds()))
	}
}

type ResponseTimeData struct {
	RequestType string        `json:"requestType"`
	Operation   string        `json:"operation"`
	Duration    time.Duration `json:"duration"`
}
