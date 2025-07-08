package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Create(ctx context.Context, pool *pgxpool.Pool) (err error) {
	err = CreateAlbumsTable(ctx, pool)
	if err != nil {
		return fmt.Errorf("failed to create albums table: %w", err)
	}
	err = CreateGenresTable(ctx, pool)
	if err != nil {
		return fmt.Errorf("failed to create genres table: %w", err)
	}
	err = CreateAlbumGenresTable(ctx, pool)
	if err != nil {
		return fmt.Errorf("failed to create albums_genres table: %w", err)
	}
	return nil
}

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

func CreateGenresTable(ctx context.Context, pool *pgxpool.Pool) error {
	// Проверяем, существует ли таблица
	var tableExists bool
	err := pool.QueryRow(ctx,
		`SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'genres'
		)`).Scan(&tableExists)

	if err != nil {
		return fmt.Errorf("failed to check table existence: %w", err)
	}

	if tableExists {
		return nil
	}

	// Создаем таблицу
	query := `
		CREATE TABLE genres (
			id BIGSERIAL PRIMARY KEY,
			name TEXT NOT NULL UNIQUE
		);

		CREATE INDEX idx_genres_name ON genres(name);
		
		COMMENT ON TABLE genres IS 'Содержит музыкальные жанры для классификации альбомов';
		COMMENT ON COLUMN genres.name IS 'Название жанра';
	`

	_, err = pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create genres table: %w", err)
	}

	return nil
}

func CreateAlbumGenresTable(ctx context.Context, pool *pgxpool.Pool) error {
	var tableExists bool
	err := pool.QueryRow(ctx,
		`SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'album_genres'
		)`).Scan(&tableExists)

	if err != nil {
		return fmt.Errorf("failed to check table existence: %w", err)
	}

	if tableExists {
		return nil
	}

	query := `
		CREATE TABLE album_genres (
			album_id BIGINT NOT NULL REFERENCES albums(id) ON DELETE CASCADE,
			genre_id BIGINT NOT NULL REFERENCES genres(id) ON DELETE CASCADE,
			PRIMARY KEY (album_id, genre_id)
		);
		
		CREATE INDEX idx_album_genres_album_id ON album_genres(album_id);
		CREATE INDEX idx_album_genres_genre_id ON album_genres(genre_id);
		
		COMMENT ON TABLE album_genres IS 'Связь между альбомами и жанрами (многие-ко-многим)';
	`

	_, err = pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create album_genres table: %w", err)
	}

	return nil
}
