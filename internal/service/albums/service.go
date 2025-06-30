package albums

import (
	"context"
	"errors"

	"github.com/4udiwe/musicshop/internal/entity"
	repo "github.com/4udiwe/musicshop/internal/repo/albums"
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
		if errors.Is(err, repo.ErrAlbumAlreadyExists) {
			return 0, ErrAlbumAlreadyExists
		}
		return 0, ErrCannotCreateAlbum
	}
	return id, nil
}

func (s *Service) FindAll(ctx context.Context) ([]entity.Album, error) {
	albums, err := s.albumRepository.FindAll(ctx)
	if err != nil {
		if errors.Is(err, repo.ErrDatabase) {
			return nil, ErrCannotFetchAlbums
		}
		return nil, ErrCannotFetchAlbums
	}
	return albums, nil
}

func (s *Service) FindById(ctx context.Context, id int64) (entity.Album, error) {
	album, err := s.albumRepository.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, repo.ErrDatabase) {
			return entity.Album{}, ErrFindingAlbum
		}
		if errors.Is(err, repo.ErrAlbumNotFound) {
			return entity.Album{}, ErrAlbumNotFound
		}
	}
	return album, nil
}

func (s *Service) DeleteById(ctx context.Context, id int64) error {
	if err := s.albumRepository.Delete(ctx, id); err != nil {
		if errors.Is(err, repo.ErrDatabase) {
			return ErrFindingAlbum
		}
		if errors.Is(err, repo.ErrAlbumNotFound) {
			return ErrAlbumNotFound
		}
	}
	return nil
}
