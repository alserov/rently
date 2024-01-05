package main

import (
	"github.com/alserov/file-storage/internal/app"
	"github.com/alserov/file-storage/internal/config"
)

func main() {
	cfg := config.MustLoad()

	application := app.NewApp(cfg)
	application.MustStart()
}
