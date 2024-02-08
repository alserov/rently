package workers

import (
	"context"
	"fmt"
	"github.com/alserov/rently/user/internal/config"
	"github.com/alserov/rently/user/internal/log"
	"github.com/alserov/rently/user/internal/utils/broker"
	"runtime"
	"time"
)

type StatMetricParams struct {
	Producer broker.Producer
	Topics   config.Topics
}

func NewStatMetricWorker(p StatMetricParams) Actor {
	return &statMetricsWorker{
		log:      log.GetLogger(),
		producer: p.Producer,
		topics:   p.Topics,
	}
}

type statMetricsWorker struct {
	log      log.Logger
	producer broker.Producer
	topics   config.Topics
}

const (
	GOROUTINE_METRIC_ID = "20"
)

func (smw *statMetricsWorker) Action() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	if err := smw.producer.Produce(ctx, runtime.Stack(make([]byte, 0), true), GOROUTINE_METRIC_ID, smw.topics.Metrics.Goroutines); err != nil {
		return fmt.Errorf("failed to produce message: %w", err)
	}

	return nil
}
