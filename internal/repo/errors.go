package repo

import "errors"

var (
	ErrDatabase = errors.New("database error")

	ErrAlbumAlreadyExists = errors.New("album already exists")
	ErrAlbumNotFound      = errors.New("album not found")

	ErrGenreNotFound               = errors.New("genre not found")
	ErrGenreAlreadyExists          = errors.New("genre already exists")
	ErrAddAlbumGenreConstraintFali = errors.New("cannot add constraint album_genre")
)
