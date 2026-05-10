package db

import (
	"context"
	"fmt"

	"github.com/ESJ0/tienda-zapatos-backend/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(cfg *config.Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return pool, nil
}
