package app

import (
	"github.com/4udiwe/avito-pvz/pkg/validator"
	"github.com/labstack/echo/v4"
)

func (app *App) EchoHandler() *echo.Echo {
	if app.echoHandler != nil {
		return app.echoHandler
	}

	handler := echo.New()
	handler.Validator = validator.NewCustomValidator()

	app.configureRouter(handler)

	app.echoHandler = handler
	return app.echoHandler
}

func (app *App) configureRouter(handler *echo.Echo) {
	pvzGroup := handler.Group("pvz")
	{
		pvzGroup.POST("", app.PostPointHandler().Handle)
		pvzGroup.GET("", app.GetPointsHandler().Handle)
		pvzGroup.POST(":pvzId/close_last_reception", app.CloseReceptionHandler().Handle)
		pvzGroup.POST(":pvzId/delete_last_product", app.DeleteProductHandler().Handle)
	}

	receptionsGroup := handler.Group("receptions")
	{
		receptionsGroup.POST("", app.PostReceptionHandler().Handle)
	}

	productsGroup := handler.Group("products")
	{
		productsGroup.POST("", app.PostProductHandler().Handle)
	}
}
