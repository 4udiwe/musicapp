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
		return 0, fmt.Errorf("%w: failed to build query: %v", repo.ErrDatabase, err)
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
			return 0, fmt.Errorf("%w: no returned id after insert", repo.ErrDatabase)
		}
		return 0, fmt.Errorf("%w: failed to execute query: %v", repo.ErrDatabase, err)
	}
	return id, nil
}
func (r *Repository) AddGenreToAlbum(ctx context.Context, albumID int64, genreID int64) error {
	query, args, err := r.pg.Builder.
		Insert("album_genres").
		Columns("album_id", "genre_id").
		Values(albumID, genreID).
		ToSql()

	if err != nil {
		return fmt.Errorf("%w: failed to build query: %v", repo.ErrDatabase, err)
	}

	result, err := r.pg.GetTxManager(ctx).Exec(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return fmt.Errorf("%w: database error code %s: %v",
				repo.ErrDatabase, pgErr.Code, pgErr.Message)
		}
		return fmt.Errorf("%w: failed to execute insert query: %v", repo.ErrDatabase, err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf(
			"%w: failed to add constraint, album id: %d, genre id: %d",
			repo.ErrAddAlbumGenreConstraintFali, albumID, genreID,
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
		return nil, fmt.Errorf("%w: failed to build query: %v", repo.ErrDatabase, err)
	}

	rows, err := r.pg.GetTxManager(ctx).Query(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return nil, fmt.Errorf("%w: database error code %s: %v",
				repo.ErrDatabase, pgErr.Code, pgErr.Message)
		}
		return nil, fmt.Errorf("%w: failed to execute query: %v", repo.ErrDatabase, err)
	}
	defer rows.Close()

	genres = make([]entity.Genre, 0)
	for rows.Next() {
		var genre entity.Genre
		if err := rows.Scan(
			&genre.ID,
			&genre.Name,
		); err != nil {
			return nil, fmt.Errorf("%w: failed to scan row: %v", repo.ErrDatabase, err)
		}
		genres = append(genres, genre)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: rows iteration error: %v", repo.ErrDatabase, err)
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
		return fmt.Errorf("%w: failed to build delete query: %v", repo.ErrDatabase, err)
	}

	result, err := r.pg.GetTxManager(ctx).Exec(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return fmt.Errorf("%w: database error code %s: %v",
				repo.ErrDatabase, pgErr.Code, pgErr.Message)
		}
		return fmt.Errorf("%w: failed to execute delete query: %v", repo.ErrDatabase, err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("%w: genre with id %d not found", repo.ErrGenreNotFound, id)
	}

	return nil
}
