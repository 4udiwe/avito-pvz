package post_point

import (
	"errors"
	"net/http"

	api "github.com/4udiwe/avito-pvz/internal/api/http"
	"github.com/4udiwe/avito-pvz/internal/api/http/decorator"
	"github.com/4udiwe/avito-pvz/internal/dto"
	service "github.com/4udiwe/avito-pvz/internal/service/point"
	"github.com/labstack/echo/v4"
)

type handler struct {
	s PointService
}

func New(pointService PointService) api.Handler {
	return decorator.NewBindAndValidateDerocator(&handler{s: pointService})
}

type Request dto.PostPvzJSONRequestBody

func (h *handler) Handle(ctx echo.Context, in Request) error {
	point, err := h.s.CreatePoint(ctx.Request().Context(), string(in.City))

	if err != nil {
		if errors.Is(err, service.ErrNoCityFound) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(
		http.StatusCreated,
		dto.EntityPointToDTO(&point),
	)
}
