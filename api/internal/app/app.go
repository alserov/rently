package app

import (
	"fmt"
	"github.com/alserov/rently/api/internal/clients"
	"github.com/alserov/rently/api/internal/config"
	"github.com/alserov/rently/api/internal/log"
	"github.com/alserov/rently/api/internal/routes"
	"github.com/alserov/rently/api/internal/server"
	"github.com/gofiber/fiber/v2"
	"log/slog"
	"os"
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

	l.Info("starting app", slog.Int("port", cfg.Port))

	cls := clients.DialServices(clients.Services{
		CarAddr:  cfg.Services.CarAddr,
		UserAddr: cfg.Services.UserAddr,
	})
	defer func() {
		if err := cls.CloseConns(); err != nil {
			l.Error("failed to close user connection(s)", slog.String("error", err.Error()))
		}
	}()

	serv := server.NewServer(server.Params{
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		Clients:      cls,
	})

	router := fiber.New()

	routes.Setup(router, serv)

	l.Info("app is running")
	run(router, cfg.Port)
	l.Info("app was stopped")
}

func run(s *fiber.App, port int) {
	chDone := make(chan os.Signal, 1)
	signal.Notify(chDone, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := s.Listen(fmt.Sprintf(":%d", port)); err != nil {
			panic("failed to listen: " + err.Error())
		}
	}()

	<-chDone
	if err := s.Shutdown(); err != nil {
		panic("failed to shutdown: " + err.Error())
	}
}
