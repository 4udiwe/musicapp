package delete_genres_by_id

import "context"

type GenreService interface {
	DeleteGenre(ctx context.Context, genreID int64) error
}
