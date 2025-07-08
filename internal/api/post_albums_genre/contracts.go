package post_albums_genre

import "context"

type GenreService interface {
	AddGenreToAlbum(ctx context.Context, albumID int64, genreID int64) error
}
