package post_album

import (
	"errors"
	"net/http"

	"github.com/4udiwe/musicshop/internal/entity"
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
	Title  string  `json:"title" validate:"required,min=2"`
	Artist string  `json:"artist" validate:"required"`
	Price  float64 `json:"price" validate:"required"`
}

type Response struct {
	ID int64 `json:"id"`
}

func (h *Handler) Handle(c echo.Context) error {
	var in Request

	if err := c.Bind(&in); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(in); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	id, err := h.albumsService.Create(c.Request().Context(), entity.Album{Title: in.Title, Artist: in.Artist, Price: in.Price})
	if err != nil {
		if errors.Is(err, service.ErrAlbumAlreadyExists) {
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, Response{ID: id})
}
