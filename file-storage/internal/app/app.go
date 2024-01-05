package app

import (
	"context"
	"fmt"
	"github.com/alserov/file-storage/internal/config"
	"github.com/alserov/file-storage/internal/server"
	"github.com/alserov/file-storage/internal/service"
	"github.com/alserov/file-storage/internal/utils/broker"
	"github.com/alserov/file-storage/internal/utils/log"
	"github.com/alserov/file-storage/internal/worker"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	"os/signal"
	"syscall"
)

type App struct {
	port int

	log *slog.Logger

	broker broker.Broker

	gRPCServer *grpc.Server
}

func NewApp(cfg *config.Config) *App {
	return &App{
		log: log.MustSetup(cfg.Env),
		broker: broker.Broker{
			Topics: broker.Topics{
				SaveImages:   cfg.Broker.Topics.SaveImages,
				DeleteImages: cfg.Broker.Topics.DeleteImages,
			},
		},
		gRPCServer: grpc.NewServer(),
	}
}

func (a *App) MustStart() {
	a.log.Info("starting app")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	imageWorker := worker.NewImageWorker(a.broker.Addr, a.broker.Topics, a.log)
	imageWorker.MustStart(ctx)

	server.RegisterGRPCServer(a.gRPCServer, service.NewService())

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		panic("failed to listen: " + err.Error())
	}

	go func() {
		err = a.gRPCServer.Serve(l)
		if err != nil {
			panic("failed to serve: " + err.Error())
		}
	}()

	a.log.Info("app is running")
	<-ctx.Done()
	a.gRPCServer.GracefulStop()
}
