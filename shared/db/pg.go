package db

import (
	"context"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func MustOpen() *pgxpool.Pool {
	url := os.Getenv("POSTGRES_URL")
	if url == "" {
		panic("POSTGRES_URL missing")
	}

	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		panic(err)
	}
	cfg.MaxConns = 10
	cfg.MaxConnIdleTime = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		panic(err)
	}
	return pool
}
