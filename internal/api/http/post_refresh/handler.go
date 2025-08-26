package post_refresh

import (
	"errors"
	"net/http"

	api "github.com/4udiwe/avito-pvz/internal/api/http"
	"github.com/4udiwe/avito-pvz/internal/api/http/decorator"
	"github.com/4udiwe/avito-pvz/internal/service/user"
	"github.com/labstack/echo/v4"
)

type handler struct {
	s UserService
}

func New(userService UserService) api.Handler {
	return decorator.NewBindAndValidateDerocator(&handler{s: userService})
}

type Request struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

func (h *handler) Handle(ctx echo.Context, in Request) error {
	tokens, err := h.s.RefreshTokens(ctx.Request().Context(), in.RefreshToken)

	if err != nil {
		if errors.Is(err, user.ErrInvalidRefreshToken) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusCreated, tokens)
}
