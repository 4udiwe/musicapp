package post_genre

import (
	"errors"
	"net/http"

	"github.com/4udiwe/musicshop/internal/entity"
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
	Name string `json:"name" validate:"required,min=3"`
}

type Response struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) Handle(c echo.Context) error {
	var in Request

	if err := c.Bind(&in); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(in); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	id, err := h.genreService.Create(c.Request().Context(), entity.Genre{Name: in.Name})
	if err != nil {
		if errors.Is(err, service.ErrGenreAlreadyExists) {
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}
		if errors.Is(err, service.ErrCannotCreateGenre) {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, Response{ID: id, Name: in.Name})

}
