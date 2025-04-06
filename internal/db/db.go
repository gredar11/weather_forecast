package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

func ConnectDb() (*pgxpool.Pool, error) {
	conConfig, err := pgxpool.ParseConfig("postgresql://postgres:12345@127.0.0.1:5432/weather")
	if err != nil {
		return nil, err
	}

	conConfig.MaxConns = 10
	conConfig.MaxConnIdleTime = 10 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), conConfig)
	if err != nil {
		log.Fatalln("Failed to connect to database:", err)
	}
	return pool, err
}
