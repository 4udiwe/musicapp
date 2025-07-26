package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/gommon/log"
)

const (
	defaultConnTimeout  = time.Second
	defaultConnAttempts = 10
)

type Postgres struct {
	connTimeout  time.Duration
	connAttempts int

	Pool    *pgxpool.Pool
	Builder squirrel.StatementBuilderType
}

func New(url string, opts ...Option) (*Postgres, error) {
	pg := &Postgres{
		connAttempts: defaultConnAttempts,
		connTimeout:  defaultConnTimeout,
	}

	// Custom options
	for _, opt := range opts {
		opt(pg)
	}

	pg.Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - pgxpool.ParseConfig: %w", err)
	}

	for pg.connAttempts > 0 {
		pg.Pool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err != nil {
			return nil, fmt.Errorf("postgres - NewPostgres - pgxpool.NewWithConfig: %w", err)
		}

		if err = pg.Pool.Ping(context.Background()); err == nil {
			break
		}

		log.Infof("Postgres is trying to connect, attempts left: %d", pg.connAttempts)
		pg.connAttempts--
		time.Sleep(pg.connTimeout)
	}

	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - connAtempts == 0: %w", err)
	}

	return pg, nil
}

func (pg *Postgres) Close() {
	if pg.Pool != nil {
		pg.Pool.Close()
	}
}
