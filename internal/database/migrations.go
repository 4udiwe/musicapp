package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateAlbumsTable(ctx context.Context, pool *pgxpool.Pool) error {
	// Проверяем, существует ли таблица
	var tableExists bool
	err := pool.QueryRow(ctx,
		`SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'albums'
		)`).Scan(&tableExists)

	if err != nil {
		return fmt.Errorf("failed to check table existence: %w", err)
	}

	if tableExists {
		return nil
	}

	// Создаем таблицу
	query := `
		CREATE TABLE albums (
			id BIGSERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			artist TEXT NOT NULL,
			price DECIMAL(10,2) NOT NULL CHECK (price >= 0),
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);

		CREATE INDEX idx_albums_artist ON albums(artist);
		CREATE INDEX idx_albums_price ON albums(price);
		
		COMMENT ON TABLE albums IS 'Содержит информацию о музыкальных альбомах';
		COMMENT ON COLUMN albums.title IS 'Название альбома';
		COMMENT ON COLUMN albums.artist IS 'Исполнитель или группа';
	`

	_, err = pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create albums table: %w", err)
	}

	return nil
}
