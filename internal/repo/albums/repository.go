package albums

import (
	"context"

	"github.com/4udiwe/musicshop/internal/entity"
	"github.com/Masterminds/squirrel"
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

func (r *Repository) Create(ctx context.Context, album entity.Album) (id int64, _ error) {
	query, args, err := r.builder.
		Insert("albums").
		Columns("title", "artist", "price").
		Values(album.Title, album.Artist, album.Price).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return 0, err
	}

	if err := r.pool.QueryRow(ctx, query, args...).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}
