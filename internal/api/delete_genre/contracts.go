package delete_genre

import "context"

type GenreService interface {
	DeleteGenre(ctx context.Context, genreID int64) error
}
