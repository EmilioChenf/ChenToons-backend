package database

import (
	"context"
	"time"

	"chentoons-backend/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ConectarPostgresChen(cfg config.Config) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL())
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}
