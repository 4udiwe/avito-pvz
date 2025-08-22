package post_reception

import (
	"errors"
	"net/http"

	api "github.com/4udiwe/avito-pvz/internal/api/http"
	"github.com/4udiwe/avito-pvz/internal/api/http/decorator"
	"github.com/4udiwe/avito-pvz/internal/dto"
	service "github.com/4udiwe/avito-pvz/internal/service/reception"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type handler struct {
	s ReceptionService
}

func New(receptionService ReceptionService) api.Handler {
	return decorator.NewBindAndValidateDerocator(&handler{s: receptionService})
}

type Request dto.PostReceptionsJSONBody

func (h *handler) Handle(ctx echo.Context, in Request) error {
	reception, err := h.s.OpenReception(ctx.Request().Context(), in.PvzId)

	if err != nil {
		logrus.Errorf("POST reception - %v", err)
		if errors.Is(err, service.ErrNoPointFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		if errors.Is(err, service.ErrLastReceptionNotClosed) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(
		http.StatusCreated,
		dto.EntityReceptionToDTO(&reception),
	)
}
