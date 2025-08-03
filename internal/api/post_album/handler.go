package post_album

import (
	"errors"
	"net/http"

	"github.com/4udiwe/musicshop/internal/api"
	"github.com/4udiwe/musicshop/internal/api/decorator"
	"github.com/4udiwe/musicshop/internal/entity"
	service "github.com/4udiwe/musicshop/internal/service/albums"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

type handler struct {
	albumsService AlbumsService
}

func New(albumsService AlbumsService) api.Handler {
	return decorator.NewBindAndValidateDerocator(&handler{
		albumsService: albumsService,
	})
}

type Genre struct {
	ID int64 `json:"id" validate:"required"`
}

type Request struct {
	Title  string  `json:"title" validate:"required,min=2"`
	Artist string  `json:"artist" validate:"required,min=2"`
	Price  float64 `json:"price" validate:"required"`
	Genres []Genre `json:"genres"`
}

type Response struct {
	ID int64 `json:"id"`
}

func (h *handler) Handle(c echo.Context, in Request) error {
	album := entity.Album{
		Title:  in.Title,
		Artist: in.Artist,
		Price:  in.Price,
		Genres: lo.Map(in.Genres, func(g Genre, i int) entity.Genre {
			return entity.Genre{
				ID: g.ID,
			}
		}),
	}
	id, err := h.albumsService.Create(c.Request().Context(), album)
	if err != nil {
		if errors.Is(err, service.ErrAlbumAlreadyExists) {
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}
		if errors.Is(err, service.ErrGenreNotExists) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, Response{ID: id})
}
