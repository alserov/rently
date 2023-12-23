package main

import (
	"github.com/alserov/rently/car/internal/app"
	"github.com/alserov/rently/car/internal/config"
)

func main() {
	cfg := config.MustLoad()

	application := app.NewApp(cfg)
	application.MustStart()
}
