package genres_test

import (
	"context"
	"errors"
	"testing"

	"github.com/4udiwe/musicshop/internal/entity"
	repo "github.com/4udiwe/musicshop/internal/repo"
	"github.com/4udiwe/musicshop/internal/repo/mock_genres"
	service "github.com/4udiwe/musicshop/internal/service/genres"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreate(t *testing.T) {
	var (
		arbitraryErr = errors.New("arbitrary error")
		ctx          = context.Background()
	)

	type MockBehavior func(r *mock_genres.MockGenreRepository)

	genre := entity.Genre{
		Name: "genre",
	}

	for _, tc := range []struct {
		name         string
		mockBehavior MockBehavior
		want         int64
		wantErr      error
	}{
		{
			name: "success",
			mockBehavior: func(r *mock_genres.MockGenreRepository) {
				r.EXPECT().Create(ctx, genre).Return(int64(1), nil)
			},
			want:    1,
			wantErr: nil,
		},
		{
			name: "genre already exists",
			mockBehavior: func(r *mock_genres.MockGenreRepository) {
				r.EXPECT().Create(ctx, genre).Return(int64(0), repo.ErrGenreAlreadyExists)
			},
			want:    0,
			wantErr: service.ErrGenreAlreadyExists,
		},
		{
			name: "cannot create genre",
			mockBehavior: func(r *mock_genres.MockGenreRepository) {
				r.EXPECT().Create(ctx, genre).Return(int64(0), arbitraryErr)
			},
			want:    0,
			wantErr: service.ErrCannotCreateGenre,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			mockGenreRepository := mock_genres.NewMockGenreRepository(ctrl)

			tc.mockBehavior(mockGenreRepository)

			s := service.New(mockGenreRepository)

			out, err := s.Create(ctx, genre)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.want, out)
		})
	}
}

func TestFindAll(t *testing.T) {
	var (
		arbitraryErr = errors.New("arbitrary error")
		ctx          = context.Background()
	)

	type MockBehavior func(r *mock_genres.MockGenreRepository)

	genres := []entity.Genre{
		{
			Name: "metal",
		},
		{
			Name: "rock",
		},
		{
			Name: "rap",
		},
	}

	for _, tc := range []struct {
		name         string
		mockBehavior MockBehavior
		want         []entity.Genre
		wantErr      error
	}{
		{
			name: "success",
			mockBehavior: func(r *mock_genres.MockGenreRepository) {
				r.EXPECT().FindAll(ctx).Return(genres, nil)
			},
			want:    genres,
			wantErr: nil,
		},
		{
			name: "cannot fetch genres",
			mockBehavior: func(r *mock_genres.MockGenreRepository) {
				r.EXPECT().FindAll(ctx).Return(nil, arbitraryErr)
			},
			want:    nil,
			wantErr: service.ErrCannotFetchGenres,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			mockGenreRepository := mock_genres.NewMockGenreRepository(ctrl)

			tc.mockBehavior(mockGenreRepository)

			s := service.New(mockGenreRepository)

			out, err := s.FindAll(ctx)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.want, out)
		})
	}
}

func TestDeleteById(t *testing.T) {
	var (
		arbitraryErr = errors.New("arbitrary error")
		ctx          = context.Background()
		id           = int64(1)
	)

	type MockBehavior func(r *mock_genres.MockGenreRepository)

	for _, tc := range []struct {
		name         string
		mockBehavior MockBehavior
		wantErr      error
	}{
		{
			name: "success",
			mockBehavior: func(r *mock_genres.MockGenreRepository) {
				r.EXPECT().Delete(ctx, id).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "album not found",
			mockBehavior: func(r *mock_genres.MockGenreRepository) {
				r.EXPECT().Delete(ctx, id).Return(repo.ErrGenreNotFound)
			},
			wantErr: service.ErrGenreNotFound,
		},
		{
			name: "cannot fetch album",
			mockBehavior: func(r *mock_genres.MockGenreRepository) {
				r.EXPECT().Delete(ctx, id).Return(arbitraryErr)
			},
			wantErr: service.ErrCannotDeleteGenre,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			mockGenreRepository := mock_genres.NewMockGenreRepository(ctrl)

			tc.mockBehavior(mockGenreRepository)

			s := service.New(mockGenreRepository)

			err := s.DeleteGenre(ctx, id)

			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestAddGenreToAlbum(t *testing.T) {
	var (
		arbitraryErr = errors.New("arbitrary error")
		ctx          = context.Background()
		albumID      = int64(111)
		genreID      = int64(222)
	)

	type MockBehavior func(r *mock_genres.MockGenreRepository)

	for _, tc := range []struct {
		name         string
		mockBehavior MockBehavior
		wantErr      error
	}{
		{
			name: "success",
			mockBehavior: func(r *mock_genres.MockGenreRepository) {
				r.EXPECT().AddGenreToAlbum(ctx, albumID, genreID).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "cannot add constraint between album and genre",
			mockBehavior: func(r *mock_genres.MockGenreRepository) {
				r.EXPECT().AddGenreToAlbum(ctx, albumID, genreID).Return(arbitraryErr)
			},
			wantErr: service.ErrCannotAddConstraintAlbumGenre,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			mockGenreRepository := mock_genres.NewMockGenreRepository(ctrl)

			tc.mockBehavior(mockGenreRepository)

			s := service.New(mockGenreRepository)

			err := s.AddGenreToAlbum(ctx, albumID, genreID)

			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}
