package delete_album

import "context"

type AlbumsService interface {
	DeleteById(ctx context.Context, id int64) error
}
