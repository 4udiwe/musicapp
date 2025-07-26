package get_albums

import (
	"net/http"

	"github.com/4udiwe/musicshop/internal/entity"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

type Handler struct {
	albumsService AlbumsService
}

func New(albumsService AlbumsService) *Handler {
	return &Handler{
		albumsService: albumsService,
	}
}

type Request struct {}

type Genre struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Album struct {
	ID     int64   `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
	Genres []Genre `json:"genres,omitempty"`
}

type Response struct {
	Albums []Album `json:"albums"`
}

func (h *Handler) Handle(c echo.Context) error {
	var in Request

	if err := c.Bind(&in); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(in); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	out, err := h.albumsService.FindAll(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	albums := lo.Map(out, func(album entity.Album, i int) Album {
		genres := lo.Map(album.Genres, func(genre entity.Genre, i int) Genre {
			return Genre{
				ID:   genre.ID,
				Name: genre.Name,
			}
		})

		return Album{
			ID:     album.ID,
			Title:  album.Title,
			Artist: album.Artist,
			Price:  album.Price,
			Genres: genres,
		}
	})

	return c.JSON(http.StatusOK, Response{Albums: albums})
}
