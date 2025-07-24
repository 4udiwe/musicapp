package delete_album

import (
	"errors"
	"net/http"

	service "github.com/4udiwe/musicshop/internal/service/albums"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	albumsService AlbumsService
}

func New(albumsService AlbumsService) *Handler {
	return &Handler{
		albumsService: albumsService,
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

	err := h.albumsService.DeleteById(c.Request().Context(), in.ID)
	if err != nil {
		if errors.Is(err, service.ErrAlbumNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
