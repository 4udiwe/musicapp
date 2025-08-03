package albums

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/4udiwe/musicshop/internal/entity"
	"github.com/4udiwe/musicshop/internal/repo"
	"github.com/4udiwe/musicshop/pkg/postgres"
)

type Repository struct {
	pg *postgres.Postgres
}

func New(postgres *postgres.Postgres) *Repository {
	return &Repository{
		pg: postgres,
	}
}

func (r *Repository) Create(ctx context.Context, album entity.Album) (id int64, err error) {
	query, args, err := r.pg.Builder.
		Insert("albums").
		Columns("title", "artist", "price").
		Values(album.Title, album.Artist, album.Price).
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

func (r *Repository) FindAll(ctx context.Context) ([]entity.Album, error) {
	query, args, _ := r.pg.Builder.
		Select(`
            a.id, 
            a.title, 
            a.artist, 
            a.price,
            ARRAY_AGG(g.id ORDER BY g.id) FILTER (WHERE g.id IS NOT NULL) as genre_ids,
            ARRAY_AGG(g.name ORDER BY g.id) FILTER (WHERE g.name IS NOT NULL) as genre_names
        `).
		From("albums a").
		LeftJoin("album_genres ag ON a.id = ag.album_id").
		LeftJoin("genres g ON ag.genre_id = g.id").
		GroupBy(`
			a.id, 
            a.title, 
            a.artist, 
            a.price`).
		ToSql()

	rows, err := r.pg.GetTxManager(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	rawAlbums, err := pgx.CollectRows(rows, pgx.RowToStructByName[albumsGenreRow])
	if err != nil {
		return nil, fmt.Errorf("failed to parse albums data: %w", err)
	}

	return convertRowsToAlbums(rawAlbums), nil
}

func (r *Repository) FindById(ctx context.Context, id int64) (album entity.Album, err error) {
	query, args, err := r.pg.Builder.
		Select("id", "title", "artist", "price").
		From("albums").
		Where(squirrel.Eq{"id": id}).
		Limit(1).
		ToSql()

	if err != nil {
		return entity.Album{}, fmt.Errorf("%w: failed to build query: %v", repo.ErrDatabase, err)
	}
	err = r.pg.GetTxManager(ctx).QueryRow(ctx, query, args...).Scan(&album.ID, &album.Title, &album.Artist, &album.Price)
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
	query, args, err := r.pg.Builder.
		Delete("albums").
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

	// Проверяем, была ли удалена хотя бы одна запись
	if result.RowsAffected() == 0 {
		return fmt.Errorf("%w: album with id %d not found", repo.ErrAlbumNotFound, id)
	}

	return nil
}
