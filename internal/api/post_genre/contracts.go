package post_genre

import (
	"context"

	"github.com/4udiwe/musicshop/internal/entity"
)

type GenreService interface {
	Create(ctx context.Context, genre entity.Genre) (int64, error)
}
