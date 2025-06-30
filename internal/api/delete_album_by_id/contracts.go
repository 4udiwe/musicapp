package delete_album_by_id

import "context"

type AlbumsService interface {
	DeleteById(ctx context.Context, id int64) error
}
