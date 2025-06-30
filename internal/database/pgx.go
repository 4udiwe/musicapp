package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/4udiwe/musicshop/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPgxPool(ctx context.Context, c *config.Config) (*pgxpool.Pool, error) {
	var pool *pgxpool.Pool
    var err error

	pgxConf, err := pgxpool.ParseConfig(c.PostgresURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse conn string: %w", err)
	}

	// Настройка retry
	retryAttempts := 5
	retryDelay := 5 * time.Second

	for i := 0; i < retryAttempts; i++ {
		pool, err = pgxpool.NewWithConfig(ctx, pgxConf)
		if err == nil {
			// Проверяем подключение
			if err = pool.Ping(ctx); err == nil {
				return pool, nil
			}
		}

		log.Printf("Attempt %d: failed to connect to database: %v", i+1, err)
		if i < retryAttempts-1 {
			log.Printf("Retrying in %v...", retryDelay)
			time.Sleep(retryDelay)
		}
	}

	return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", retryAttempts, err)
}
