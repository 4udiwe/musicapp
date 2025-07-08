package genres

import (
	"context"
	"errors"

	"github.com/4udiwe/musicshop/internal/entity"
	repo "github.com/4udiwe/musicshop/internal/repo"
)

type Service struct {
	genreRepository GenreRepository
}

func New(r GenreRepository) *Service {
	return &Service{
		genreRepository: r,
	}
}

func (s *Service) Create(ctx context.Context, genre entity.Genre) (int64, error) {
	id, err := s.genreRepository.Create(ctx, genre)
	if err != nil {
		if errors.Is(err, repo.ErrGenreAlreadyExists) {
			return 0, ErrGenreAlreadyExists
		}
		return 0, ErrCannotCreateGenre
	}
	return id, nil
}

func (s *Service) FindAll(ctx context.Context) (genres []entity.Genre, err error) {
	genres, err = s.genreRepository.FindAll(ctx)
	if err != nil {
		return []entity.Genre{}, ErrCannotFetchGenres
	}
	return genres, nil
}

func (s *Service) DeleteGenre(ctx context.Context, genreID int64) error {
	err := s.genreRepository.Delete(ctx, genreID)
	if err != nil {
		if errors.Is(err, repo.ErrGenreNotFound) {
			return ErrGenreNotFound
		}
		return ErrCannotDeleteGenre
	}
	return nil
}

func (s *Service) AddGenreToAlbum(ctx context.Context, albumID int64, genreID int64) error {
	err := s.genreRepository.AddGenreToAlbum(ctx, albumID, genreID)
	if err != nil {
		return ErrCannotAddConstraintAlbumGenre
	}
	return nil
}
