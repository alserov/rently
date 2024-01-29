package main

import "github.com/alserov/rently/user/internal/app"

func main() {
	cfg := config.MustLoad()
	app.MustStart(cfg)
}
