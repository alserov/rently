package main

import (
	"github.com/alserov/rently/api/internal/app"
	"github.com/alserov/rently/api/internal/config"
)

func main() {
	cfg := config.MustLoad()
	app.MustStart(cfg)
}
