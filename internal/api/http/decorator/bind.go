package decorator

import (
	"net/http"

	api "github.com/4udiwe/avito-pvz/internal/api/http"
	"github.com/labstack/echo/v4"
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
	var in T

	if err := c.Bind(&in); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(in); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return d.inner.Handle(c, in)
}
