package albums

import (
	"context"
	"errors"
	"fmt"

	"github.com/4udiwe/musicshop/internal/entity"
	"github.com/4udiwe/musicshop/internal/repo"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool    *pgxpool.Pool
	builder squirrel.StatementBuilderType
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool:    pool,
		builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *Repository) Create(ctx context.Context, album entity.Album) (id int64, err error) {
	query, args, err := r.builder.
		Insert("albums").
		Columns("title", "artist", "price").
		Values(album.Title, album.Artist, album.Price).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("%w: failed to build query: %v", repo.ErrDatabase, err)
	}

	err = r.pool.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" { // unique_violation
				return 0, fmt.Errorf("%w: album with title '%s' and artist '%s' already exists",
					repo.ErrAlbumAlreadyExists, album.Title, album.Artist)
			}
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("%w: no returned id after insert", repo.ErrDatabase)
		}
		return 0, fmt.Errorf("%w: failed to execute query: %v", repo.ErrDatabase, err)
	}
	return id, nil
}

func (r *Repository) FindAll(ctx context.Context) (albums []entity.Album, err error) {
	query, args, err := r.builder.
		Select("id", "title", "artist", "price").
		From("albums").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("%w: failed to build query: %v", repo.ErrDatabase, err)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return nil, fmt.Errorf("%w: database error code %s: %v",
				repo.ErrDatabase, pgErr.Code, pgErr.Message)
		}
		return nil, fmt.Errorf("%w: failed to execute query: %v", repo.ErrDatabase, err)
	}
	defer rows.Close()

	albums = make([]entity.Album, 0)
	for rows.Next() {
		var album entity.Album
		if err := rows.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
			return nil, fmt.Errorf("%w: failed to scan row: %v", repo.ErrDatabase, err)
		}
		albums = append(albums, album)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: rows iteration error: %v", repo.ErrDatabase, err)
	}

	if len(albums) == 0 {
		return []entity.Album{}, nil // Явно возвращаем пустой слайс вместо nil
	}

	return albums, nil
}

func (r *Repository) FindById(ctx context.Context, id int64) (album entity.Album, err error) {
	query, args, err := r.builder.
		Select("id", "title", "artist", "price").
		From("albums").
		Where(squirrel.Eq{"id": id}).
		Limit(1).
		ToSql()

	if err != nil {
		return entity.Album{}, fmt.Errorf("%w: failed to build query: %v", repo.ErrDatabase, err)
	}
	err = r.pool.QueryRow(ctx, query, args...).Scan(&album.ID, &album.Title, &album.Artist, &album.Price)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Album{}, fmt.Errorf("%w: album with id '%d' not found", repo.ErrAlbumNotFound, id)
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return entity.Album{}, fmt.Errorf("%w: database error code %s: %v",
				repo.ErrDatabase, pgErr.Code, pgErr.Message)
		}

		return entity.Album{}, fmt.Errorf("%w: failed to execute query: %v", repo.ErrDatabase, err)
	}

	return album, nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	query, args, err := r.builder.
		Delete("albums").
		Where(squirrel.Eq{"id": id}).
		ToSql()

	if err != nil {
		return fmt.Errorf("%w: failed to build delete query: %v", repo.ErrDatabase, err)
	}

	// Используем Exec вместо Query, так как нам не нужны возвращаемые строки
	result, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return fmt.Errorf("%w: database error code %s: %v",
				repo.ErrDatabase, pgErr.Code, pgErr.Message)
		}
		return fmt.Errorf("%w: failed to execute delete query: %v", repo.ErrDatabase, err)
	}

	// Проверяем, была ли удалена хотя бы одна запись
	if result.RowsAffected() == 0 {
		return fmt.Errorf("%w: album with id %d not found", repo.ErrAlbumNotFound, id)
	}

	return nil
}
