package clickhouse

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/skinkvi/event-tracker/internal/config"
)

func (s *Storage) Migrations(cfg config.Config) {
	m, err := migrate.New("file://migrations/clickhouse", cfg.ClickHouseDSN)
	if err != nil {
		s.log.Error("failed to migrate", err)
		return
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		s.log.Error("failed to up migration %w", err.Error())
		return
	}
}
