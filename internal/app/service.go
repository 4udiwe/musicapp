package app

import (
	"github.com/4udiwe/musicshop/internal/service/albums"
	"github.com/4udiwe/musicshop/internal/service/genres"
)

func (app *App) AlbumsService() *albums.Service {
	if app.albumsService != nil {
		return app.albumsService
	}
	app.albumsService = albums.New(app.AlbumsRepo())
	return app.albumsService
}

func (app *App) GenresService() *genres.Service {
	if app.genresService != nil {
		return app.genresService
	}
	app.genresService = genres.New(app.GenresRepo())
	return app.genresService
}
