package post_dummy_login

import (
	"net/http"

	api "github.com/4udiwe/avito-pvz/internal/api/http"
	"github.com/4udiwe/avito-pvz/internal/api/http/decorator"
	"github.com/4udiwe/avito-pvz/internal/dto"
	"github.com/4udiwe/avito-pvz/internal/entity"
	"github.com/labstack/echo/v4"
)

type handler struct {
	s UserService
}

func New(userService UserService) api.Handler {
	return decorator.NewBindAndValidateDerocator(&handler{s: userService})
}

type Request dto.PostDummyLoginJSONBody

func (h *handler) Handle(ctx echo.Context, in Request) error {
	token, err := h.s.DummyLogin(ctx.Request().Context(), entity.UserRole(in.Role))

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusCreated, token)
}
