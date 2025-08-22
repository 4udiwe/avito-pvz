package get_points

import (
	"net/http"

	api "github.com/4udiwe/avito-pvz/internal/api/http"
	"github.com/4udiwe/avito-pvz/internal/api/http/decorator"
	"github.com/4udiwe/avito-pvz/internal/dto"
	"github.com/4udiwe/avito-pvz/internal/entity"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

type handler struct {
	s PointService
}

func New(pointService PointService) api.Handler {
	return decorator.NewBindAndValidateDerocator(&handler{
		s: pointService,
	})
}

type Request struct{}

type Response struct {
	Points []dto.PVZ
}

func (h *handler) Handle(ctx echo.Context, in Request) error {
	points, err := h.s.GetAllPoints(ctx.Request().Context())

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusCreated,
		Response{
			Points: lo.Map(points, func(e entity.Point, i int) dto.PVZ {
				return *dto.EntityPointToDTO(&e)
			}),
		})
}
