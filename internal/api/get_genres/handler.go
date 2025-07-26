package get_genres

import (
	"net/http"

	"github.com/4udiwe/musicshop/internal/entity"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	genreService GenreService
}

func New(gs GenreService) *Handler {
	return &Handler{
		genreService: gs,
	}
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

func (h *Handler) Handle(c echo.Context) error {
	var in Request

	if err := c.Bind(&in); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(in); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	genres, err := h.genreService.FindAll(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, ToResponseList(genres))
}
