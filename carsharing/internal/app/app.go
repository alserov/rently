package app

import (
	"fmt"
	"github.com/alserov/rently/carsharing/internal/cache/redis"
	"github.com/alserov/rently/carsharing/internal/clients"
	"github.com/alserov/rently/carsharing/internal/config"
	"github.com/alserov/rently/carsharing/internal/db/postgres"
	"github.com/alserov/rently/carsharing/internal/log"
	"github.com/alserov/rently/carsharing/internal/metrics"
	"github.com/alserov/rently/carsharing/internal/notifications"
	"github.com/alserov/rently/carsharing/internal/payment"
	"github.com/alserov/rently/carsharing/internal/server"
	"github.com/alserov/rently/carsharing/internal/service"
	"github.com/alserov/rently/carsharing/internal/storage"
	"github.com/alserov/rently/carsharing/internal/utils/broker/rabbit"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func MustStart(cfg *config.Config) {
	l := log.MustSetup(cfg.Env)

	defer func() {
		err := recover()
		if err != nil {
			l.Error("panic recovery: ", err)
		}
	}()

	l.Info("starting app", slog.Int("port", cfg.Port))

	ch, conn := rabbit.NewProducer(cfg.Broker.Addr)
	defer conn.Close()
	defer ch.Close()

	cls := clients.DialServices(clients.Services{
		UserAddr: cfg.Services.User,
	})
	defer func() {
		if err := cls.CloseConns(); err != nil {
			l.Error("failed to close clients connection(s)", slog.String("error", err.Error()))
		}
	}()

	serv := service.NewService(service.Params{
		Repo:          postgres.NewRepo(postgres.MustConnect(cfg.DB.GetDsn())),
		Notifications: notifications.NewNotifier(),
		Payment:       payment.NewPayer("sk_test_51OU56CDOnc0MdcTNBwddO2cn8NrEebjfuAGjBjj9xSyKmiUO4ajJ1vZ0yBoOsAMq0HjHqCmis2niwoj2EZYCDLOA00lcCUlWxh"),
		ImageStorage:  storage.NewImageStorage(),
		UserClient:    clients.NewUserClient(cls.UserClient),
	})

	gRPCServer := grpc.NewServer()

	server.RegisterGRPCServer(gRPCServer, server.Params{
		Service: serv,
		Cache: redis.NewCache(redis.MustConnect(redis.Params{
			Addr:     cfg.Cache.Addr,
			Password: cfg.Cache.Password,
		})),
		Metrics: metrics.NewMetrics(ch, cfg.Broker.Topics.Metrics),
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
	s.GracefulStop()
}
