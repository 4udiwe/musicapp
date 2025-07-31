package get_genres

import (
	"net/http"

	"github.com/4udiwe/musicshop/internal/api"
	"github.com/4udiwe/musicshop/internal/api/decorator"
	"github.com/4udiwe/musicshop/internal/entity"
	"github.com/labstack/echo/v4"
)

type handler struct {
	genreService GenreService
}

func New(gs GenreService) api.Handler {
	return decorator.NewBindAndValidateDerocator(&handler{
		genreService: gs,
	})
}

type Request struct{}

type Response struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// ToResponse преобразует entity в DTO
func ToResponse(g entity.Genre) Response {
	return Response{
		ID:   g.ID,
		Name: g.Name,
	}
}

// ToResponseList преобразует список entities в список DTO
func ToResponseList(genres []entity.Genre) []Response {
	result := make([]Response, len(genres))
	for i, g := range genres {
		result[i] = ToResponse(g)
	}
	return result
}

func (h *handler) Handle(c echo.Context, in Request) error {
	genres, err := h.genreService.FindAll(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, ToResponseList(genres))
}
