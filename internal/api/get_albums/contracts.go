package get_albums

import (
	"context"

	"github.com/4udiwe/musicshop/internal/entity"
)

type AlbumsService interface {
	FindAll(ctx context.Context) ([]entity.Album, error)
}
