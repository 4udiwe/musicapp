package get_album_by_id

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

func (h *Handler) Handle(c echo.Context) error {
	var in Request

	if err := c.Bind(&in); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(in); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	album, err := h.albumsService.FindById(c.Request().Context(), in.ID)
	if err != nil {
		if errors.Is(err, service.ErrAlbumNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, ToResponse(album))
}
