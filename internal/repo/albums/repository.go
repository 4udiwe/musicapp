package albums

import (
	"context"
	"database/sql"

	"github.com/4udiwe/musicshop/internal/entity"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, album entity.Album) (id int64, _ error) {
	query := `
		INSERT INTO albums (title, artist, price)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	if err := r.db.QueryRowContext(ctx, query, album.Title, album.Artist, album.Price).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}
