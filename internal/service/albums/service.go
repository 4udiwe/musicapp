package albums

import (
	"context"
	"errors"

	"github.com/4udiwe/musicshop/internal/entity"
	repo "github.com/4udiwe/musicshop/internal/repo"
	"github.com/4udiwe/musicshop/pkg/transactor"
	"github.com/sirupsen/logrus"
)

type Service struct {
	albumRepository AlbumRepository
	genreRepository GenreRepository
	txManager       transactor.Transactor
}

func New(
	a AlbumRepository,
	g GenreRepository,
	t transactor.Transactor,
) *Service {
	return &Service{
		albumRepository: a,
		genreRepository: g,
		txManager:       t,
	}
}

func (s *Service) Create(ctx context.Context, a entity.Album) (int64, error) {
	var id int64

	err := s.txManager.WithinTransaction(ctx, func(ctx context.Context) error {
		var err error
		id, err = s.albumRepository.Create(ctx, a)
		if err != nil {
			return err
		}
		for _, genre := range a.Genres {
			err = s.genreRepository.AddGenreToAlbum(ctx, id, genre.ID)
			if err != nil {
				return err
			}
		}
		return err

	})

	if err != nil {
		logrus.Infof("Result err = %v", err.Error())
		if errors.Is(err, repo.ErrAlbumAlreadyExists) {
			return 0, ErrAlbumAlreadyExists
		}
		if errors.Is(err, repo.ErrAddAlbumGenreConstraintFail) {
			return 0, ErrGenreNotExists
		}
		return 0, ErrCannotCreateAlbum
	}

	return id, nil
}

func (s *Service) FindAll(ctx context.Context) ([]entity.Album, error) {
	albums, err := s.albumRepository.FindAll(ctx)
	if err != nil {
		return nil, ErrCannotFetchAlbums
	}

	return albums, nil
}

func (s *Service) FindById(ctx context.Context, id int64) (entity.Album, error) {
	album, err := s.albumRepository.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, repo.ErrAlbumNotFound) {
			return entity.Album{}, ErrAlbumNotFound
		}
		return entity.Album{}, ErrFindingAlbum
	}
	return album, nil
}

func (s *Service) DeleteById(ctx context.Context, id int64) error {
	if err := s.albumRepository.Delete(ctx, id); err != nil {
		if errors.Is(err, repo.ErrAlbumNotFound) {
			return ErrAlbumNotFound
		}
		return ErrFindingAlbum
	}
	return nil
}
