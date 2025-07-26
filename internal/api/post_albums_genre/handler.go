package post_albums_genre

import (
	"errors"
	"net/http"

	service "github.com/4udiwe/musicshop/internal/service/genres"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	genreService GenreService
}

func New(gs GenreService) *Handler {
	return &Handler{
		genreService: gs,
	}
}

type Request struct {
	AlbumID int64 `param:"id" validate:"required"`
	GenreID int64 `json:"genre_id" validate:"required"`
}

func (h *Handler) Handle(c echo.Context) error {
	var in Request

	if err := c.Bind(&in); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(in); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	err := h.genreService.AddGenreToAlbum(c.Request().Context(), in.AlbumID, in.GenreID)
	if err != nil {
		if errors.Is(err, service.ErrCannotAddConstraintAlbumGenre) {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusCreated)
}
