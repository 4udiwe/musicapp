package get_genres

import (
	"context"

	"github.com/4udiwe/musicshop/internal/entity"
)

type GenreService interface {
	FindAll(ctx context.Context) (genres []entity.Genre, err error)
}
