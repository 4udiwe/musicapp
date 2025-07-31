package post_album

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
	Title  string  `json:"title" validate:"required,min=2"`
	Artist string  `json:"artist" validate:"required,min=2"`
	Price  float64 `json:"price" validate:"required"`
}

type Response struct {
	ID int64 `json:"id"`
}

func (h *handler) Handle(c echo.Context, in Request) error {
	id, err := h.albumsService.Create(c.Request().Context(), entity.Album{Title: in.Title, Artist: in.Artist, Price: in.Price})
	if err != nil {
		if errors.Is(err, service.ErrAlbumAlreadyExists) {
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, Response{ID: id})
}
