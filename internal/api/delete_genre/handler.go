package delete_genre

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
	ID int64 `param:"id" validate:"required"`
}

func (h *handler) Handle(c echo.Context, in Request) error {
	err := h.genreService.DeleteGenre(c.Request().Context(), in.ID)
	if err != nil {
		if errors.Is(err, service.ErrGenreNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusOK)
}
