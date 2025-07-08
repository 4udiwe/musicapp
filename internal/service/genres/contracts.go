package genres

import (
	"context"

	"github.com/4udiwe/musicshop/internal/entity"
)

type GenreRepository interface {
	Create(ctx context.Context, genre entity.Genre) (int64, error)
	AddGenreToAlbum(ctx context.Context, albumID int64, genreID int64) error
	FindAll(ctx context.Context) ([]entity.Genre, error)
	Delete(ctx context.Context, id int64) error
}
