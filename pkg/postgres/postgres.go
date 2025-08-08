package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
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

// Transaction management

type txKey struct{}

// injectTx добавляет транзакцию в контекст.
func injectTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// extractTx извлекает транзакцию из контекста, если она там присутствует.
func extractTx(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(pgx.Tx)
	return tx, ok
}

type TxManager interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func (pg *Postgres) GetTxManager(ctx context.Context) TxManager {
	if tx, ok := extractTx(ctx); ok {
		return tx
	}
	return pg.Pool
}

func (pg *Postgres) WithinTransaction(ctx context.Context, fn func(context.Context) error) error {
	tx, err := pg.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("postgres - Begin transaction: %w", err)
	}

	ctxTx := injectTx(ctx, tx)

	if err := fn(ctxTx); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}
