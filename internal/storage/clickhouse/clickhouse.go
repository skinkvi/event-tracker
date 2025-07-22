package clickhouse

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/skinkvi/event-tracker/internal/config"
)

type Storage struct {
	db  *sql.DB
	log *slog.Logger
}

func New(ctx context.Context, cfg config.Config, log *slog.Logger) (*Storage, error) {
	const fn = "storage.clickhouse.New"

	db, err := sql.Open("clickhouse", cfg.ClickHouseDSN)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &Storage{
		db:  db,
		log: log,
	}, nil
}
