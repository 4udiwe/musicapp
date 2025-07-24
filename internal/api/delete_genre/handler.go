package delete_genre

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
	ID int64 `param:"id" validate:"required"`
}

func (h *Handler) Handle(c echo.Context) error {
	var in Request

	if err := c.Bind(&in); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(in); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	err := h.genreService.DeleteGenre(c.Request().Context(), in.ID)
	if err != nil {
		if errors.Is(err, service.ErrGenreNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusOK)
}
