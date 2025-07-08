package genres

import "errors"

var (
	ErrCannotCreateGenre             = errors.New("cannot create genre")
	ErrGenreAlreadyExists            = errors.New("genre already exists")
	ErrCannotFetchGenres             = errors.New("cannot fetch genres")
	ErrCannotAddConstraintAlbumGenre = errors.New("cannot add constraint benween album and genre")
	ErrGenreNotFound                 = errors.New("genre not found")
	ErrCannotDeleteGenre             = errors.New("cannot delete genre")
)
