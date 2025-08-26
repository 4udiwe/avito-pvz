package post_login

import (
	"errors"
	"net/http"

	api "github.com/4udiwe/avito-pvz/internal/api/http"
	"github.com/4udiwe/avito-pvz/internal/api/http/decorator"
	"github.com/4udiwe/avito-pvz/internal/dto"
	"github.com/4udiwe/avito-pvz/internal/service/user"
	"github.com/labstack/echo/v4"
)

type handler struct {
	s UserService
}

func New(userService UserService) api.Handler {
	return decorator.NewBindAndValidateDerocator(&handler{s: userService})
}

type Request dto.PostLoginJSONBody

func (h *handler) Handle(ctx echo.Context, in Request) error {
	tokens, err := h.s.Authenticate(ctx.Request().Context(), string(in.Email), in.Password)

	if err != nil {
		if errors.Is(err, user.ErrNoUserFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		if errors.Is(err, user.ErrInvalidCredentials) {
			return echo.NewHTTPError(http.StatusForbidden)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusCreated, tokens)
}
