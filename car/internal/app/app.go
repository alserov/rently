package app

import (
	"fmt"
	broker2 "github.com/alserov/rently/car/internal/utils/broker"

	"github.com/alserov/rently/car/internal/config"
	"github.com/alserov/rently/car/internal/db/postgres"
	"github.com/alserov/rently/car/internal/log"
	mtrcs "github.com/alserov/rently/car/internal/metrics"
	"github.com/alserov/rently/car/internal/server"
	"github.com/alserov/rently/car/internal/service"

	"google.golang.org/grpc"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	port int
	log  *slog.Logger

	dsn string

	broker broker2.Broker

	gRPCServer *grpc.Server
}

func NewApp(cfg *config.Config) *App {
	return &App{
		port: cfg.Port,

		broker: broker2.Broker{
			Addr: cfg.Broker.Addr,
			Metrics: mtrcs.MetricTopics{
				DecreaseActiveRentsAmount: cfg.Broker.Metrics.Topics.DecreaseActiveRentsAmount,
				IncreaseActiveRentsAmount: cfg.Broker.Metrics.Topics.IncreaseActiveRentsAmount,
			},
		},

		log:        log.MustSetup(cfg.Env),
		gRPCServer: grpc.NewServer(),
	}
}

func (a *App) MustStart() {
	defer func() {
		err := recover()
		if err != nil {
			a.log.Error("panic recovery: ", err)
		}
	}()

	a.log.Info("starting app", slog.Int("port", a.port))

	db := postgres.MustConnect(a.dsn)
	repo := postgres.NewRepo(db)

	producer := broker2.NewProducer(a.broker.Addr)
	metr := mtrcs.NewMetrics(producer, a.broker.Metrics, a.log)

	serv := service.NewService(repo, metr, a.log)

	server.RegisterGRPCServer(a.gRPCServer, serv)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		panic("failed to listen: " + err.Error())
	}

	a.run(l)
}

func (a *App) run(l net.Listener) {
	chStop := make(chan os.Signal)
	signal.Notify(chStop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := a.gRPCServer.Serve(l); err != nil {
			panic("failed to serve: " + err.Error())
		}
	}()

	sign := <-chStop
	a.gRPCServer.GracefulStop()
	a.log.Info("app was stopped", slog.String("signal", sign.String()))
}
