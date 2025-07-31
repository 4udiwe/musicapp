package post_albums_genre

import (
	"errors"
	"net/http"

	"github.com/4udiwe/musicshop/internal/api"
	"github.com/4udiwe/musicshop/internal/api/decorator"
	service "github.com/4udiwe/musicshop/internal/service/genres"
	"github.com/labstack/echo/v4"
)

type handler struct {
	genreService GenreService
}

func New(gs GenreService) api.Handler {
	return decorator.NewBindAndValidateDerocator(&handler{
		genreService: gs,
	})
}

type Request struct {
	AlbumID int64 `param:"id" validate:"required"`
	GenreID int64 `json:"genre_id" validate:"required"`
}

func (h *handler) Handle(c echo.Context, in Request) error {
	err := h.genreService.AddGenreToAlbum(c.Request().Context(), in.AlbumID, in.GenreID)
	if err != nil {
		if errors.Is(err, service.ErrCannotAddConstraintAlbumGenre) {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusCreated)
}
