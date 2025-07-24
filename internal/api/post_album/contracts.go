package post_album

import (
	"context"

	"github.com/4udiwe/musicshop/internal/entity"
)

type AlbumsService interface {
	Create(ctx context.Context, album entity.Album) (int64, error)
}
