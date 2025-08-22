package delete_product

import (
	"errors"
	"net/http"

	api "github.com/4udiwe/avito-pvz/internal/api/http"
	"github.com/4udiwe/avito-pvz/internal/api/http/decorator"
	service "github.com/4udiwe/avito-pvz/internal/service/product"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type handler struct {
	s ProductService
}

func New(productService ProductService) api.Handler {
	return decorator.NewBindAndValidateDerocator(&handler{
		s: productService,
	})
}

type Request struct {
	PointID uuid.UUID `param:"pvzId" validate:"required"`
}

func (h *handler) Handle(ctx echo.Context, in Request) error {
	err := h.s.DeleteLastProductFromReception(ctx.Request().Context(), in.PointID)

	if err != nil {
		if errors.Is(err, service.ErrNoPointFound) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if errors.Is(err, service.ErrReceptionAlreadyClosed) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.NoContent(http.StatusOK)
}
