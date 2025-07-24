package redis

import (
	"context"
	"log/slog"

	"github.com/skinkvi/event-tracker/internal/config"
)

type Storage struct {
	// разобраться что тут должно быть, пологаю что клиент редиса, типо redis.Client, но я не уверен
}

func New(ctx context.Context, cfg config.Config, log *slog.Logger) (*Storage, error) {
	// нужно инициализировать клиент редиса с конфига, потом пингануть, и вернуть storage
}
