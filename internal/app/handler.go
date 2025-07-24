package app

import (
	"github.com/4udiwe/musicshop/internal/api"
	"github.com/4udiwe/musicshop/internal/api/delete_album"
	"github.com/4udiwe/musicshop/internal/api/delete_genre"
	"github.com/4udiwe/musicshop/internal/api/get_album"
	"github.com/4udiwe/musicshop/internal/api/get_albums"
	"github.com/4udiwe/musicshop/internal/api/get_genres"
	"github.com/4udiwe/musicshop/internal/api/post_album"
	"github.com/4udiwe/musicshop/internal/api/post_albums_genre"
	"github.com/4udiwe/musicshop/internal/api/post_genre"
)

func (app *App) DeleteAlbumHandler() api.Handler {
	if app.deleteAlbumHandler != nil {
		return app.deleteAlbumHandler
	}
	app.deleteAlbumHandler = delete_album.New(app.AlbumsService())
	return app.deleteAlbumHandler
}

func (app *App) DeleteGenreHandler() api.Handler {
	if app.deleteGenreHandler != nil {
		return app.deleteGenreHandler
	}
	app.deleteGenreHandler = delete_genre.New(app.GenresService())
	return app.deleteGenreHandler
}

func (app *App) GetAlbumHandler() api.Handler {
	if app.getAlbumHandler != nil {
		return app.getAlbumHandler
	}
	app.getAlbumHandler = get_album.New(app.AlbumsService())
	return app.getAlbumHandler
}

func (app *App) GetAlbumsHandler() api.Handler {
	if app.getAlbumsHandler != nil {
		return app.getAlbumsHandler
	}
	app.getAlbumsHandler = get_albums.New(app.AlbumsService())
	return app.getAlbumsHandler
}

func (app *App) GetGenresHandler() api.Handler {
	if app.getGenresHandler != nil {
		return app.getGenresHandler
	}
	app.getGenresHandler = get_genres.New(app.GenresService())
	return app.getGenresHandler
}

func (app *App) PostAlbumsHandler() api.Handler {
	if app.postAlbumHandler != nil {
		return app.postAlbumHandler
	}
	app.postAlbumHandler = post_album.New(app.AlbumsService())
	return app.postAlbumHandler
}

func (app *App) PostAlbumGenreHandler() api.Handler {
	if app.postAlbumGenreHandler != nil {
		return app.postAlbumGenreHandler
	}
	app.postAlbumGenreHandler = post_albums_genre.New(app.GenresService())
	return app.postAlbumGenreHandler
}

func (app *App) PostGenreHandler() api.Handler {
	if app.postGenreHandler != nil {
		return app.postGenreHandler
	}
	app.postGenreHandler = post_genre.New(app.GenresService())
	return app.postGenreHandler
}
