package get_album

import (
	"errors"
	"net/http"

	"github.com/4udiwe/musicshop/internal/api"
	"github.com/4udiwe/musicshop/internal/api/decorator"
	"github.com/4udiwe/musicshop/internal/entity"
	service "github.com/4udiwe/musicshop/internal/service/albums"
	"github.com/labstack/echo/v4"
)

type handler struct {
	albumsService AlbumsService
}

func New(albumsService AlbumsService) api.Handler {
	return decorator.NewBindAndValidateDerocator(&handler{
		albumsService: albumsService,
	})
}

type Request struct {
	ID int64 `param:"id" validate:"required"`
}

type Response struct {
	ID     int64   `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// ToResponse преобразует entity в DTO
func ToResponse(a entity.Album) Response {
	return Response{
		ID:     a.ID,
		Title:  a.Title,
		Artist: a.Artist,
		Price:  a.Price,
	}
}

func (h *handler) Handle(c echo.Context, in Request) error {
	album, err := h.albumsService.FindById(c.Request().Context(), in.ID)
	if err != nil {
		if errors.Is(err, service.ErrAlbumNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, ToResponse(album))
}
