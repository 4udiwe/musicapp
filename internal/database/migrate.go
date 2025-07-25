package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

// RunMigrations применяет все pending миграции
func RunMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	// Конвертируем pgxpool.Pool в *sql.DB (goose требует database/sql)
	db, err := pgxPoolToStdlib(ctx, pool)
	if err != nil {
		return fmt.Errorf("failed to convert pool: %w", err)
	}

	// Устанавливаем диалект БД
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	// Применяем миграции
	if err := goose.Up(db, "database/migrations"); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	// Логируем текущую версию
	version, err := goose.GetDBVersion(db)
	if err != nil {
		return fmt.Errorf("failed to get DB version: %w", err)
	}
	log.Printf("Migrations applied. Current version: %d", version)

	return nil
}

// pgxPoolToStdlib конвертирует pgxpool.Pool в *sql.DB
func pgxPoolToStdlib(ctx context.Context, pool *pgxpool.Pool) (*sql.DB, error) {
	// Используем адаптер из pgx
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	// Конвертируем в stdlib (database/sql)
	db := stdlib.OpenDB(*conn.Conn().Config())
	return db, nil
}
