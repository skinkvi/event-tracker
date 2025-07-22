package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/tern/v2/migrate"
	"github.com/skinkvi/event-tracker/internal/config"
	"github.com/skinkvi/event-tracker/internal/logger"
	"github.com/skinkvi/event-tracker/internal/storage/clickhouse"
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

	postgresStorage, err := postgres.New(ctx, *cfg)
	if err != nil {
		return
	}

	conn, err := postgresStorage.Get().Acquire(ctx)
	if err != nil {
		log.Error("failed get conn postgres")
		return
	}
	defer conn.Release()

	migrator, err := migrate.NewMigrator(ctx, conn.Conn(), "schema_version")
	if err != nil {
		log.Error("failed create a migrator for postgres")
		return
	}

	migrationsDir := "./migrations"
	if err := migrator.LoadMigrations(os.DirFS(migrationsDir)); err != nil {
		log.Error("failed to load migration")
		return
	}

	if err := migrator.Migrate(ctx); err != nil {
		log.Error("migration failed")
		return
	}
	log.Info("migrations postgres applied successfully")
	log.Info("postgres storage init")

	clickhouseStorage, err := clickhouse.New(ctx, *cfg, log)
	if err != nil {
		return
	}

	clickhouseStorage.Migrations(*cfg)
	log.Info("migrations clickhouse applied successfully")
	log.Info("clickhouse storage init")
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
