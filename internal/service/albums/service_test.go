package albums_test

import (
	"context"
	"errors"
	"testing"

	"github.com/4udiwe/musicshop/internal/entity"
	"github.com/4udiwe/musicshop/internal/mocks/mock_albums"
	"github.com/4udiwe/musicshop/internal/mocks/mock_transactor"
	repo "github.com/4udiwe/musicshop/internal/repo"
	service "github.com/4udiwe/musicshop/internal/service/albums"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreate(t *testing.T) {
	var (
		arbitraryErr = errors.New("arbitrary error")
		ctx          = context.Background()
	)

	type MockBehavior func(r *mock_albums.MockAlbumRepository, t *mock_transactor.MockTransactor)

	album := entity.Album{
		Title:  "title",
		Artist: "artist",
		Price:  100.0,
	}

	for _, tc := range []struct {
		name         string
		mockBehavior MockBehavior
		want         int64
		wantErr      error
	}{
		{
			name: "success",
			mockBehavior: func(r *mock_albums.MockAlbumRepository, t *mock_transactor.MockTransactor) {
				r.EXPECT().Create(ctx, album).Return(int64(1), nil)
				t.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(
						func(ctx context.Context, fn func(ctx context.Context) error) error {
							return fn(ctx)
						})

			},
			want:    1,
			wantErr: nil,
		},
		{
			name: "album already exists",
			mockBehavior: func(r *mock_albums.MockAlbumRepository, t *mock_transactor.MockTransactor) {
				r.EXPECT().Create(ctx, album).Return(int64(0), repo.ErrAlbumAlreadyExists)
				t.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(
						func(ctx context.Context, fn func(ctx context.Context) error) error {
							return fn(ctx)
						})
			},
			want:    0,
			wantErr: service.ErrAlbumAlreadyExists,
		},
		{
			name: "cannot create album",
			mockBehavior: func(r *mock_albums.MockAlbumRepository, t *mock_transactor.MockTransactor) {
				r.EXPECT().Create(ctx, album).Return(int64(0), arbitraryErr)
				t.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(
						func(ctx context.Context, fn func(ctx context.Context) error) error {
							return fn(ctx)
						})
			},
			want:    0,
			wantErr: service.ErrCannotCreateAlbum,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			mockAlbumRepository := mock_albums.NewMockAlbumRepository(ctrl)
			mockTransactor := mock_transactor.NewMockTransactor(ctrl)

			tc.mockBehavior(mockAlbumRepository, mockTransactor)

			s := service.New(mockAlbumRepository, mockTransactor)

			out, err := s.Create(ctx, album)

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

	type MockBehavior func(r *mock_albums.MockAlbumRepository)

	albums := []entity.Album{
		{
			Title:  "title 1",
			Artist: "artist 1",
			Price:  100.0,
		},
		{
			Title:  "title 2",
			Artist: "artist 2",
			Price:  200.0,
		},
		{
			Title:  "title 3",
			Artist: "artist 3",
			Price:  300.0,
		},
	}

	for _, tc := range []struct {
		name         string
		mockBehavior MockBehavior
		want         []entity.Album
		wantErr      error
	}{
		{
			name: "success",
			mockBehavior: func(r *mock_albums.MockAlbumRepository) {
				r.EXPECT().FindAll(ctx).Return(albums, nil)
			},
			want:    albums,
			wantErr: nil,
		},
		{
			name: "cannot fetch albums",
			mockBehavior: func(r *mock_albums.MockAlbumRepository) {
				r.EXPECT().FindAll(ctx).Return(nil, arbitraryErr)
			},
			want:    nil,
			wantErr: service.ErrCannotFetchAlbums,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			mockAlbumRepository := mock_albums.NewMockAlbumRepository(ctrl)
			mockTransactor := mock_transactor.NewMockTransactor(ctrl)

			tc.mockBehavior(mockAlbumRepository)

			s := service.New(mockAlbumRepository, mockTransactor)

			out, err := s.FindAll(ctx)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.want, out)
		})
	}
}

func TestFindById(t *testing.T) {
	var (
		arbitraryErr = errors.New("arbitrary error")
		ctx          = context.Background()
		id           = int64(1)
		emptyAlbum   = entity.Album{}
	)

	type MockBehavior func(r *mock_albums.MockAlbumRepository)

	album := entity.Album{
		Title:  "title",
		Artist: "artist",
		Price:  100.0,
	}

	for _, tc := range []struct {
		name         string
		mockBehavior MockBehavior
		want         entity.Album
		wantErr      error
	}{
		{
			name: "success",
			mockBehavior: func(r *mock_albums.MockAlbumRepository) {
				r.EXPECT().FindById(ctx, id).Return(album, nil)
			},
			want:    album,
			wantErr: nil,
		},
		{
			name: "album not found",
			mockBehavior: func(r *mock_albums.MockAlbumRepository) {
				r.EXPECT().FindById(ctx, id).Return(emptyAlbum, repo.ErrAlbumNotFound)
			},
			want:    emptyAlbum,
			wantErr: service.ErrAlbumNotFound,
		},
		{
			name: "cannot fetch album",
			mockBehavior: func(r *mock_albums.MockAlbumRepository) {
				r.EXPECT().FindById(ctx, id).Return(emptyAlbum, arbitraryErr)
			},
			want:    emptyAlbum,
			wantErr: service.ErrFindingAlbum,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			mockAlbumRepository := mock_albums.NewMockAlbumRepository(ctrl)
			mockTransactor := mock_transactor.NewMockTransactor(ctrl)

			tc.mockBehavior(mockAlbumRepository)

			s := service.New(mockAlbumRepository, mockTransactor)

			out, err := s.FindById(ctx, id)

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

	type MockBehavior func(r *mock_albums.MockAlbumRepository)

	for _, tc := range []struct {
		name         string
		mockBehavior MockBehavior
		wantErr      error
	}{
		{
			name: "success",
			mockBehavior: func(r *mock_albums.MockAlbumRepository) {
				r.EXPECT().Delete(ctx, id).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "album not found",
			mockBehavior: func(r *mock_albums.MockAlbumRepository) {
				r.EXPECT().Delete(ctx, id).Return(repo.ErrAlbumNotFound)
			},
			wantErr: service.ErrAlbumNotFound,
		},
		{
			name: "cannot fetch album",
			mockBehavior: func(r *mock_albums.MockAlbumRepository) {
				r.EXPECT().Delete(ctx, id).Return(arbitraryErr)
			},
			wantErr: service.ErrFindingAlbum,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			mockAlbumRepository := mock_albums.NewMockAlbumRepository(ctrl)
			mockTransactor := mock_transactor.NewMockTransactor(ctrl)

			tc.mockBehavior(mockAlbumRepository)

			s := service.New(mockAlbumRepository, mockTransactor)

			err := s.DeleteById(ctx, id)

			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}
