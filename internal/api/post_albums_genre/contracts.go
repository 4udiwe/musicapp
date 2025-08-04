package post_albums_genre

import "context"

type GenreService interface {
	AddGenresToAlbum(ctx context.Context, albumID int64, genreIDs ...int64) error
}
