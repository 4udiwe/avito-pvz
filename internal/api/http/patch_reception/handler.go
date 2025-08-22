package patch_reception

import (
	"errors"
	"net/http"

	api "github.com/4udiwe/avito-pvz/internal/api/http"
	"github.com/4udiwe/avito-pvz/internal/api/http/decorator"
	service "github.com/4udiwe/avito-pvz/internal/service/reception"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type handler struct {
	s ReceptionService
}

func New(receptionService ReceptionService) api.Handler {
	return decorator.NewBindAndValidateDerocator(&handler{s: receptionService})
}

type Request struct {
	PointID uuid.UUID `param:"pvzId" validate:"required"`
}

func (h *handler) Handle(
	ctx echo.Context,
	in Request,
) error {
	err := h.s.CloseReception(ctx.Request().Context(), in.PointID)

	if err != nil {
		if errors.Is(err, service.ErrNoPointFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		if errors.Is(err, service.ErrNoReceptionFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		if errors.Is(err, service.ErrLastReceptionAlreadyClosed) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if errors.Is(err, service.ErrCannotCloseEmptyReception) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.NoContent(http.StatusAccepted)
}
