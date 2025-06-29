package albums

import (
	"context"
	"errors"

	"github.com/4udiwe/musicshop/internal/entity"
	"github.com/4udiwe/musicshop/internal/repo/albums"
)

type Service struct {
	albumRepository AlbumRepository
}

func New(albumRepository AlbumRepository) *Service {
	return &Service{
		albumRepository: albumRepository,
	}
}

func (s *Service) Create(ctx context.Context, a entity.Album) (int64, error) {
	id, err := s.albumRepository.Create(ctx, a) 
	if err != nil {
		if errors.Is(err, albums.ErrAlbumAlreadyExists) {
			return 0, ErrAlbumAlreadyExists
		}
		return 0, ErrCannotCreateAlbum
	}
	return id, nil
}

