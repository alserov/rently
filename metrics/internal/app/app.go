package app

import (
	"fmt"
	"github.com/alserov/rently/metrics/internal/config"
	"github.com/alserov/rently/metrics/internal/log"
	"github.com/alserov/rently/metrics/internal/workers"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func MustStart(cfg *config.Config) {
	l := log.MustSetup(cfg.Env)

	l.Info("starting app", slog.Int("port", cfg.Port))

	manager := workers.NewWorkerManager()
	manager.Add(workers.NewCarsharingWorker(cfg.Broker))

	reg := prometheus.NewRegistry()
	reg.MustRegister(manager.Metrics()...)

	mux := http.NewServeMux()
	prometheusHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	mux.Handle("/metrics", prometheusHandler)

	l.Info("app is running")
	run(mux, cfg.Port)
	l.Info("app was stopped")
}

func run(h http.Handler, port int) {
	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", port), h); err != nil {
			panic("failed to serve: " + err.Error())
		}
	}()

	chStop := make(chan os.Signal)
	signal.Notify(chStop, syscall.SIGINT, syscall.SIGTERM)
	<-chStop
}
