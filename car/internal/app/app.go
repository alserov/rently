package app

import (
	"fmt"
	"github.com/alserov/rently/car/internal/utils/clients"
	"github.com/alserov/rently/car/internal/utils/files"

	"github.com/alserov/rently/car/internal/config"
	"github.com/alserov/rently/car/internal/db/postgres"
	"github.com/alserov/rently/car/internal/metrics"
	"github.com/alserov/rently/car/internal/server"
	"github.com/alserov/rently/car/internal/service"
	"github.com/alserov/rently/car/internal/utils/broker"
	"github.com/alserov/rently/car/internal/utils/log"

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

	broker broker.Broker

	services clients.Services

	gRPCServer *grpc.Server
}

func NewApp(cfg *config.Config) *App {
	return &App{
		port: cfg.Port,

		dsn: fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.DB.Name, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name),
		broker: broker.Broker{
			Addr: cfg.Broker.Addr,
			Topics: broker.Topics{
				Metrics: broker.MetricTopics{
					DecreaseActiveRentsAmount: cfg.Broker.Topics.Metrics.DecreaseActiveRentsAmount,
					IncreaseActiveRentsAmount: cfg.Broker.Topics.Metrics.IncreaseActiveRentsAmount,
					IncreaseRentsCancel:       cfg.Broker.Topics.Metrics.IncreaseRentsCancel,
					NotifyBrandDemand:         cfg.Broker.Topics.Metrics.NotifyBrandDemand,
				},
				Images: broker.ImageTopics{
					Save:   cfg.Broker.Topics.Files.SaveImages,
					Delete: cfg.Broker.Topics.Files.DeleteImages,
				},
			},
		},

		services: clients.Services{
			FileStorageAddr: cfg.Services.FileStorageAddr,
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

	asyncProducer := broker.NewAsyncProducer(a.broker.Addr)
	syncProducer := broker.NewSyncProducer(a.broker.Addr)

	// metrics
	metr := metrics.NewMetrics(syncProducer, a.broker.Topics.Metrics, a.log)

	// dials other services and returns their clients, connections and error
	cli, conn, closeConn := clients.SetupClients(clients.Services{
		FileStorageAddr: a.services.FileStorageAddr,
	})
	defer closeConn(conn) // closes all connections with other services

	// image functionality, save, delete, get
	imager := files.NewImager(asyncProducer, cli.FileStorage, a.broker.Topics.Images, a.log)

	// bll
	serv := service.NewService(repo, imager, a.log)

	server.RegisterGRPCServer(a.gRPCServer, server.Server{
		Service: serv,
		Metrics: metr,
		Clients: cli,
	}, a.log)

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

	a.log.Info("app is running")
	sign := <-chStop
	a.gRPCServer.GracefulStop()
	a.log.Info("app was stopped", slog.String("signal", sign.String()))
}
