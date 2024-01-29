package main

import (
	"github.com/alserov/rently/metrics/internal/app"
	"github.com/alserov/rently/metrics/internal/config"
)

func main() {
	cfg := config.MustLoad()
	app.MustStart(cfg)
}
