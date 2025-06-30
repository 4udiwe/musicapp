package albums

import "errors"

var (
	ErrAlbumAlreadyExists = errors.New("album already exists")
	ErrDatabase           = errors.New("database error")
	ErrAlbumNotFound      = errors.New("album not found")
)
