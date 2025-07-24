package app

import (
	"github.com/4udiwe/musicshop/pkg/validator"
	"github.com/labstack/echo/v4"
)

func (app *App) EchoHandler() *echo.Echo {
	if app.echoHandler != nil {
		return app.echoHandler
	}

	handler := echo.New()
	handler.Validator = validator.NewCustomValidator()

	app.configureRouter(handler)

	app.echoHandler = handler
	return app.echoHandler
}

func (app *App) configureRouter(handler *echo.Echo) {
	albumsGroup := handler.Group("/albums")
	{
		albumsGroup.GET("", app.GetAlbumsHandler().Handle)
		albumsGroup.GET("/:id", app.GetAlbumHandler().Handle)
		albumsGroup.POST("", app.PostAlbumGenreHandler().Handle)
		albumsGroup.POST("/genres", app.PostAlbumGenreHandler().Handle)
		albumsGroup.DELETE("/:id", app.DeleteAlbumHandler().Handle)
	}

	genresGroup := handler.Group("/genres")
	{
		genresGroup.GET("", app.GetGenresHandler().Handle)
		genresGroup.POST("", app.PostGenreHandler().Handle)
		genresGroup.DELETE("", app.DeleteGenreHandler().Handle)
	}
}
