package get_album

import (
	"context"

	"github.com/4udiwe/musicshop/internal/entity"
)

type AlbumsService interface {
	FindById(ctx context.Context, id int64) (entity.Album, error)
}