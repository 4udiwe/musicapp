package app

import (
	"github.com/4udiwe/musicshop/internal/repo/albums"
	"github.com/4udiwe/musicshop/internal/repo/genres"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (app *App) Pool() *pgxpool.Pool {
	return app.pgxPool
}

func (app *App) AlbumsRepo() *albums.Repository {
	if app.albumsRepo != nil {
		return app.albumsRepo
	}
	app.albumsRepo = albums.New(app.pgxPool)
	return app.albumsRepo
}

func (app *App) GenresRepo() *genres.Repository {
	if app.genresRepo != nil {
		return app.genresRepo
	}
	app.genresRepo = genres.New(app.pgxPool)
	return app.genresRepo
}
