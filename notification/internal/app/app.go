package app

import (
	"context"
	"github.com/alserov/rently/notifications/internal/config"
	"github.com/alserov/rently/notifications/internal/log"
	"github.com/alserov/rently/notifications/internal/utils/broker/rabbit"
	"github.com/alserov/rently/notifications/internal/workers"
	"log/slog"
	"os/signal"
	"syscall"
)

func MustStart(cfg *config.Config) {
	l := log.MustSetup(cfg.Env)

	defer func() {
		err := recover()
		if err != nil {
			l.Error("panic recovery", slog.Any("error", err))
		}
	}()

	l.Info("starting app")
	l.Debug("config", slog.Any("config", *cfg))

	wm := workers.NewManager()

	rbtConn, rbtCh, err := rabbit.Dial(cfg.Broker.Addr)
	defer func() {
		if err = rbtConn.Close(); err != nil {
			l.Error("failed to close rabbit connection", slog.String("error", err.Error()))
		}
		if err = rbtCh.Close(); err != nil {
			l.Error("failed to close rabbit channel", slog.String("error", err.Error()))
		}
	}()
	if err != nil {
		panic(err)
	}

	cons := rabbit.NewConsumer(rbtCh)
	wm.Add(workers.NewEmailWorker(587, "smtp.gmail.com", cons.Consume(cfg.Broker.Topics.Email)))

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	l.Info("app is running")
	<-ctx.Done()
}
