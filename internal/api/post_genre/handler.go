package post_genre

import (
	"errors"
	"net/http"

	"github.com/4udiwe/musicshop/internal/api"
	"github.com/4udiwe/musicshop/internal/api/decorator"
	"github.com/4udiwe/musicshop/internal/entity"
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
	Name string `json:"name" validate:"required,min=3"`
}

type Response struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (h *handler) Handle(c echo.Context, in Request) error {
	id, err := h.genreService.Create(c.Request().Context(), entity.Genre{Name: in.Name})
	if err != nil {
		if errors.Is(err, service.ErrGenreAlreadyExists) {
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, Response{ID: id, Name: in.Name})
}
