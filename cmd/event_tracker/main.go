package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/skinkvi/event-tracker/internal/config"
	"github.com/skinkvi/event-tracker/internal/logger"
	"github.com/skinkvi/event-tracker/internal/storage/postgres"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogget(cfg.Env)

	ctx := context.Background()

	storage, err := postgres.New(ctx, *cfg)
	if err != nil {
		return
	}

	log.Info("storage init")

	//TODO: something else
}

// мы выносим потому что логи должны быть разные в зависимости от окружения либо local, prod и т.д.
func setupLogget(env string) *slog.Logger {
	opts := &slog.HandlerOptions{Level: slog.LevelDebug}

	var log *slog.Logger

	switch env {
	case envLocal, envDev:
		log = slog.New(logger.NewCustomHandler(opts))

	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, opts))
	}

	return log
}
