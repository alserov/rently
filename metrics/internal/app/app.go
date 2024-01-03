package app

import (
	"context"
	"fmt"
	"github.com/alserov/rently/metrics/internal/config"
	"github.com/alserov/rently/metrics/internal/log"
	"github.com/alserov/rently/metrics/internal/metrics"
	"github.com/alserov/rently/metrics/internal/metrics/car"
	"github.com/alserov/rently/metrics/internal/utils/broker"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
)

type App struct {
	port int

	log *slog.Logger

	broker broker.Broker
}

func NewApp(cfg *config.Config) *App {
	return &App{
		port: cfg.Port,
		log:  log.MustSetup(cfg.Env),
		broker: broker.Broker{
			Addr: cfg.Broker.Addr,
			Topics: broker.Topics{
				IncreaseActiveRentsAmount: cfg.Broker.Topics.IncreaseActiveRentsAmount,
				DecreaseActiveRentsAmount: cfg.Broker.Topics.DecreaseActiveRentsAmount,
			},
		},
	}
}

func (a *App) MustStart() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	a.log.Info("starting app")

	carWorker, carMetrics := car.NewWorker(a.broker.Addr, &car.GaugeTopics{
		IncreaseActiveRentsAmount: a.broker.Topics.IncreaseActiveRentsAmount,
		DecreaseActiveRentsAmount: a.broker.Topics.DecreaseActiveRentsAmount,
	})
	go func() {
		carWorker.MustStart(ctx)
	}()

	reg := prometheus.NewRegistry()
	metrics.Register(reg, carMetrics)

	mux := http.NewServeMux()
	prometheusHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	mux.Handle("/metrics", prometheusHandler)

	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", a.port), mux); err != nil {
			panic("failed to serve: " + err.Error())
		}
	}()

	a.log.Info("app is running", slog.Int("port", a.port))
	select {
	case <-ctx.Done():
		a.log.Info("app was stopped")
	}
}
