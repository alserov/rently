package app

import (
	"fmt"
	"github.com/alserov/rently/user/internal/config"
	"github.com/alserov/rently/user/internal/db/mysql"
	"github.com/alserov/rently/user/internal/log"
	"github.com/alserov/rently/user/internal/notifications"
	"github.com/alserov/rently/user/internal/server"
	"github.com/alserov/rently/user/internal/service"
	"github.com/alserov/rently/user/internal/workers"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func MustStart(cfg *config.Config) {
	l := log.MustSetup(cfg.Env)
	defer func() {
		err := recover()
		if err != nil {
			l.Error("panic recovery", slog.Any("error", err))
		}
	}()

	l.Info("starting app", slog.Int("port", cfg.Port))

	go workers.StartNotifier(time.NewTicker(time.Hour*24), workers.NewRentNotifier())

	srvc := service.NewService(service.Params{
		Repo:     mysql.NewRepository(mysql.MustConnect(cfg.DB.GetDSN())),
		Notifier: notifications.NewNotifier(),
	})

	gRPCServer := grpc.NewServer()

	server.RegisterGRPCServer(gRPCServer, server.Params{
		Service: srvc,
	})

	l.Info("app is running")
	run(gRPCServer, cfg.Port)
	l.Info("app was stopped")
}

func run(s *grpc.Server, port int) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic("failed to listen: " + err.Error())
	}

	chStop := make(chan os.Signal, 1)
	signal.Notify(chStop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err = s.Serve(l); err != nil {
			panic("failed to serve: " + err.Error())
		}
	}()

	<-chStop
}
