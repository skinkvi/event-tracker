package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/skinkvi/event-tracker/internal/config"
)

type Storage struct {
	db *pgxpool.Pool
}

func New(ctx context.Context, cfg config.Config) (*Storage, error) {
	const fn = "storage.postgres.New"

	config, err := pgxpool.ParseConfig(cfg.PostgresURL)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", fn, err)
	}

	db, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", fn, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Close() {
	s.db.Close()
}

func (s *Storage) Get() *pgxpool.Pool {
	return s.db
}
