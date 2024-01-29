package main

import (
	"github.com/alserov/rently/carsharing/internal/app"
	"github.com/alserov/rently/carsharing/internal/config"
)

func main() {
	cfg := config.MustLoad()
	app.MustStart(cfg)
}
