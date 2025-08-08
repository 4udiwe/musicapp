package genres

import (
	"context"
	"errors"
	"fmt"

	"github.com/4udiwe/musicshop/internal/entity"
	"github.com/4udiwe/musicshop/internal/repo"
	"github.com/4udiwe/musicshop/pkg/postgres"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgconn"
)

type Repository struct {
	pg *postgres.Postgres
}

func New(postgres *postgres.Postgres) *Repository {
	return &Repository{
		pg: postgres,
	}
}

func (r *Repository) Create(ctx context.Context, genre entity.Genre) (id int64, err error) {
	query, args, err := r.pg.Builder.
		Insert("genres").
		Columns("name").
		Values(genre.Name).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build query: %w", err)
	}

	err = r.pg.GetTxManager(ctx).QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" { // unique_violation
				return 0, fmt.Errorf("%w: genre '%s' already exists",
					repo.ErrGenreAlreadyExists, genre.Name)
			}
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("%w: no returned id after insert", err)
		}
		return 0, fmt.Errorf("failed to execute query: %w", err)
	}
	return id, nil
}

func (r *Repository) AddGenresToAlbum(ctx context.Context, albumID int64, genreIDs ...int64) error {
	if len(genreIDs) < 1 {
		return repo.ErrCannotAddEmptyGenres
	}

	builder := r.pg.Builder.
		Insert("album_genres").
		Columns("album_id", "genre_id")

	for _, genreID := range genreIDs {
		builder = builder.Values(albumID, genreID)
	}

	query, args, err := builder.ToSql()

	if err != nil {
		return fmt.Errorf("%w: failed to build query", err)
	}

	result, err := r.pg.GetTxManager(ctx).Exec(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return fmt.Errorf("%w: database error code %s: %v",
				repo.ErrAddAlbumGenreConstraintFail, pgErr.Code, pgErr.Message)
		}
		return fmt.Errorf("%w: failed to execute insert query", err)
	}

	if int(result.RowsAffected()) != len(genreIDs) {
		return fmt.Errorf(
			"%w: failed to add all constraints, album id: %d, expected %d rows affected, got %d",
			repo.ErrAddAlbumGenreConstraintFail, albumID, len(genreIDs), result.RowsAffected(),
		)
	}

	return nil
}

func (r *Repository) FindAll(ctx context.Context) (genres []entity.Genre, err error) {
	query, args, err := r.pg.Builder.
		Select("id", "name").
		From("genres").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.pg.GetTxManager(ctx).Query(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return nil, fmt.Errorf("%w: database error code %s: %v",
				err, pgErr.Code, pgErr.Message)
		}
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	genres = make([]entity.Genre, 0)
	for rows.Next() {
		var genre entity.Genre
		if err := rows.Scan(
			&genre.ID,
			&genre.Name,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		genres = append(genres, genre)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	if len(genres) == 0 {
		return []entity.Genre{}, nil
	}

	return genres, nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	query, args, err := r.pg.Builder.
		Delete("genres").
		Where(squirrel.Eq{"id": id}).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build delete query: %w", err)
	}

	result, err := r.pg.GetTxManager(ctx).Exec(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return fmt.Errorf("%w: database error code %s: %v",
				err, pgErr.Code, pgErr.Message)
		}
		return fmt.Errorf("failed to execute delete query: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("%w: genre with id %d not found", repo.ErrGenreNotFound, id)
	}

	return nil
}
