package main

import (
	"github.com/alserov/rently/notifications/internal/app"
	"github.com/alserov/rently/notifications/internal/config"
)

func main() {
	cfg := config.MustLoad()
	app.MustStart(cfg)
}
