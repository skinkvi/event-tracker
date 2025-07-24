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

func (s *Storage) CountEventByType(ctx context.Context) (map[string]int, error) {
	const fn = "storage.clickhouse.CountEventByType"

	query := `SELECT event_type, COUNT(*) as count FROM events WHERE timestamp => now() - INTERVAL 24 HOUR GROUP BY event_type`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	defer rows.Close()

	counts := make(map[string]int)
	for rows.Next() {
		var eventType string
		var count int
		err := rows.Scan(&eventType, &count)
		if err != nil {
			return nil, fmt.Errorf("%s, %w", fn, err)
		}

		counts[eventType] = count
	}

	return counts, nil
}
