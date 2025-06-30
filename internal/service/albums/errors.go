package albums

import "errors"

var (
	ErrAlbumAlreadyExists = errors.New("album already exists")
	ErrCannotCreateAlbum  = errors.New("cannot create album")
	ErrCannotFetchAlbums  = errors.New("cannot fetch albums")
	ErrFindingAlbum       = errors.New("error finding album")
	ErrAlbumNotFound      = errors.New("album not found")
)
