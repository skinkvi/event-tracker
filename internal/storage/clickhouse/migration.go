package clickhouse

import (
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/skinkvi/event-tracker/internal/config"
)

func (s *Storage) Migrations(cfg config.Config) {
	m, err := migrate.New("file://migrations/clickhouse", cfg.ClickHouseDSN)
	if err != nil {
		s.log.Error("failed to migrate", slog.Any("error", err))
		return
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		s.log.Error("failed to up migration %w", slog.Any("error", err))
		return
	}
}
