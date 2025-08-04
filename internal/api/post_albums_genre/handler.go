package post_albums_genre

import (
	"errors"
	"net/http"

	"github.com/4udiwe/musicshop/internal/api"
	"github.com/4udiwe/musicshop/internal/api/decorator"
	service "github.com/4udiwe/musicshop/internal/service/genres"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

type handler struct {
	genreService GenreService
}

func New(gs GenreService) api.Handler {
	return decorator.NewBindAndValidateDerocator(&handler{
		genreService: gs,
	})
}

type Genre struct {
	ID int64 `json:"genre_id" validate:"required"`
}

type Request struct {
	AlbumID int64   `param:"id" validate:"required"`
	Genres  []Genre `json:"genres"`
}

func (h *handler) Handle(c echo.Context, in Request) error {
	genreIDs := lo.Map(in.Genres, func(g Genre, i int) int64 {
		return g.ID
	})
	err := h.genreService.AddGenresToAlbum(c.Request().Context(), in.AlbumID, genreIDs...)
	if err != nil {
		if errors.Is(err, service.ErrCannotAddConstraintAlbumGenre) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusCreated)
}
