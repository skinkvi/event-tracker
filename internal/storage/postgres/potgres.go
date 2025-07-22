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

	pool, err := pgxpool.New(ctx, cfg.PostgresURL)
	if err != nil {
		return nil, fmt.Errorf("fn: %v. failed connect to postgres, err: %s", fn, err)
	}

	return &Storage{db: pool}, nil
}
