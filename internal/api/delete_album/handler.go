package delete_album

import (
	"errors"
	"net/http"

	"github.com/4udiwe/musicshop/internal/api"
	"github.com/4udiwe/musicshop/internal/api/decorator"
	service "github.com/4udiwe/musicshop/internal/service/albums"
	"github.com/labstack/echo/v4"
)

type handler struct {
	albumsService AlbumsService
}

func New(albumsService AlbumsService) api.Handler {
	return decorator.NewBindAndValidateDerocator(&handler{albumsService: albumsService})
}

type Request struct {
	ID int64 `param:"id" validate:"required"`
}

func (h *handler) Handle(c echo.Context, in Request) error {
	err := h.albumsService.DeleteById(c.Request().Context(), in.ID)
	if err != nil {
		if errors.Is(err, service.ErrAlbumNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
