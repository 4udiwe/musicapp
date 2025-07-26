package app

import (
	"github.com/4udiwe/musicshop/internal/repo/albums"
	"github.com/4udiwe/musicshop/internal/repo/genres"
	"github.com/4udiwe/musicshop/pkg/postgres"
)

func (app *App) Postgres() *postgres.Postgres {
	return app.postgres
}

func (app *App) AlbumsRepo() *albums.Repository {
	if app.albumsRepo != nil {
		return app.albumsRepo
	}
	app.albumsRepo = albums.New(app.Postgres())
	return app.albumsRepo
}

func (app *App) GenresRepo() *genres.Repository {
	if app.genresRepo != nil {
		return app.genresRepo
	}
	app.genresRepo = genres.New(app.Postgres())
	return app.genresRepo
}
