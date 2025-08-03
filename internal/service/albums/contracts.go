package albums

import (
	"context"

	"github.com/4udiwe/musicshop/internal/entity"
)

type AlbumRepository interface {
	Create(ctx context.Context, album entity.Album) (int64, error)
	FindAll(ctx context.Context) ([]entity.Album, error)
	FindById(ctx context.Context, id int64) (entity.Album, error)
	Delete(ctx context.Context, id int64) error
}

type GenreRepository interface {
	Create(ctx context.Context, genre entity.Genre) (int64, error)
	AddGenreToAlbum(ctx context.Context, albumID int64, genreID int64) error
	FindAll(ctx context.Context) ([]entity.Genre, error)
	Delete(ctx context.Context, id int64) error
}
