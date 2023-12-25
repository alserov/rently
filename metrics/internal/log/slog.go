package log

import (
	"log/slog"
	"os"
)

const (
	ENV_LOCAL = "local"
	ENV_PROD  = "prod"
)

func MustSetup(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case ENV_LOCAL:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case ENV_PROD:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	default:
		panic("unknown env")
	}

	return log
}
