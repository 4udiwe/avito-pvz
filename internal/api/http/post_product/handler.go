package post_product

import (
	"errors"
	"net/http"

	api "github.com/4udiwe/avito-pvz/internal/api/http"
	"github.com/4udiwe/avito-pvz/internal/api/http/decorator"
	"github.com/4udiwe/avito-pvz/internal/dto"
	"github.com/4udiwe/avito-pvz/internal/entity"
	service "github.com/4udiwe/avito-pvz/internal/service/product"
	"github.com/labstack/echo/v4"
)

type handler struct {
	s ProductService
}

func New(productService ProductService) api.Handler {
	return decorator.NewBindAndValidateDerocator(&handler{s: productService})
}

type Request dto.PostProductsJSONBody

func (h *handler) Handle(ctx echo.Context, in Request) error {
	product, err := h.s.AddProduct(ctx.Request().Context(), in.PvzId, entity.ProductType(in.Type))

	if err != nil {
		if errors.Is(err, service.ErrNoPointFound) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if errors.Is(err, service.ErrNoReceptionFound) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if errors.Is(err, service.ErrReceptionAlreadyClosed) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(
		http.StatusCreated,
		dto.EntityProductToDTO(&product),
	)
}
