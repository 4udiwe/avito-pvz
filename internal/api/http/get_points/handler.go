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
	Info []PointWithReceptions
}

type PointWithReceptions struct {
	Pvz        dto.PVZ                 `json:"pvz"`
	Receptions []ReceptionWithProducts `json:"receptions"`
}

type ReceptionWithProducts struct {
	Reception dto.Reception `json:"reception"`
	Products  []dto.Product `json:"products"`
}

func (h *handler) Handle(ctx echo.Context, in Request) error {
	pointsInfo, err := h.s.GetAllPointsFullInfo(ctx.Request().Context())

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusCreated,
		Response{
			Info: lo.Map(pointsInfo, func(item entity.PointFullInfo, _ int) PointWithReceptions {
				return PointWithReceptions{
					Pvz: *dto.EntityPointToDTO(&item.Point),
					Receptions: lo.Map(item.Receptions, func(e entity.ReceptionWithProducts, _ int) ReceptionWithProducts {
						return ReceptionWithProducts{
							Reception: *dto.EntityReceptionToDTO(&e.Reception),
							Products: lo.Map(e.Products, func(p entity.Product, _ int) dto.Product {
								return *dto.EntityProductToDTO(&p)
							}),
						}
					}),
				}
			}),
		})
}
