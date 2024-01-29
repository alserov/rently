package app

import (
	"fmt"
	"github.com/alserov/rently/carsharing/internal/cache/redis"
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

	serv := service.NewService(service.Params{
		Metrics:       metrics.NewMetrics(ch, cfg.Broker.Topics.Metrics, l),
		Repo:          postgres.NewRepo(postgres.MustConnect(cfg.DB.GetDsn())),
		Notifications: notifications.NewNotifier(),
		Payment:       payment.NewPayer("sk_test_51OU56CDOnc0MdcTNBwddO2cn8NrEebjfuAGjBjj9xSyKmiUO4ajJ1vZ0yBoOsAMq0HjHqCmis2niwoj2EZYCDLOA00lcCUlWxh"),
		ImageStorage:  storage.NewImageStorage(),
	})

	gRPCServer := grpc.NewServer()

	server.RegisterGRPCServer(gRPCServer, server.Server{
		Service: serv,
		Cache: redis.NewCache(redis.MustConnect(redis.Params{
			Addr:     cfg.Cache.Addr,
			Password: cfg.Cache.Password,
		})),
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

	chStop := make(chan os.Signal)
	signal.Notify(chStop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err = s.Serve(l); err != nil {
			panic("failed to serve: " + err.Error())
		}
	}()

	<-chStop
	s.GracefulStop()
}
