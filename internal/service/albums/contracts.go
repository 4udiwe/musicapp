package albums

import (
	"context"

	"github.com/4udiwe/musicshop/internal/entity"
)

type AlbumRepository interface {
	Create(ctx context.Context, album entity.Album) (int64, error)
	// FindAll(ctx context.Context) ([]entity.Album, error)
	// FindById(ctx context.Context, id string) (entity.Album, error)
	// Delete(ctx context.Context, id string) error
}
 