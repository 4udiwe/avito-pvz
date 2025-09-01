package decorator

import (
	"net/http"

	api "github.com/4udiwe/avito-pvz/internal/api/http"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type handler[T any] interface {
	Handle(c echo.Context, in T) error
}

type bindAndValidateDecorator[T any] struct {
	inner handler[T]
}

func NewBindAndValidateDerocator[T any](inner handler[T]) api.Handler {
	return &bindAndValidateDecorator[T]{inner: inner}
}

func (d *bindAndValidateDecorator[T]) Handle(c echo.Context) error {
	logrus.Infof("HTTP %s %s from %s", c.Request().Method, c.Path(), c.Request().RemoteAddr)

	var in T

	if err := c.Bind(&in); err != nil {
		// каст ошибки производится для правильного тестирования и отображения только message в handler, без кода
		logrus.Errorf("Failed to bind request: %v", err.(*echo.HTTPError))
		return echo.NewHTTPError(http.StatusBadRequest, err.(*echo.HTTPError).Message)
	}

	if err := c.Validate(in); err != nil {
		logrus.Errorf("Failed to validate request: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.(*echo.HTTPError).Message)
	}

	return d.inner.Handle(c, in)
}
