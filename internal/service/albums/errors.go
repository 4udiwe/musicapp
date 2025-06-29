package albums

import "errors"

var (
	ErrAlbumAlreadyExists = errors.New("album already exists")
	ErrCannotCreateAlbum  = errors.New("cannot create album")
)
