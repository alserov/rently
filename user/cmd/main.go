package main

import (
	"github.com/alserov/rently/user/internal/app"
	"github.com/alserov/rently/user/internal/config"
)

func main() {
	cfg := config.MustLoad()
	app.MustStart(cfg)
}
