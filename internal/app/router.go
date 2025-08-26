package app

import (
	"fmt"

	"github.com/4udiwe/avito-pvz/internal/api/http/middleware"
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

	for _, r := range handler.Routes() {
		fmt.Printf("%s %s\n", r.Method, r.Path)
	}

	app.echoHandler = handler
	return app.echoHandler
}

func (app *App) configureRouter(handler *echo.Echo) {
	handler.POST("dummyLogin", app.PostDummyLoginHandler().Handle)
	handler.POST("register", app.PostRegisterHandler().Handle)
	handler.POST("login", app.PostLoginHandler().Handle)
	handler.POST("refresh", app.PostRefreshHandler().Handle)

	receptionsGroup := handler.Group("receptions", app.AuthMiddleware().Middleware)
	{
		receptionsGroup.POST("", app.PostReceptionHandler().Handle, middleware.EmployeeOnly)
	}

	productsGroup := handler.Group("products", app.AuthMiddleware().Middleware)
	{
		productsGroup.POST("", app.PostProductHandler().Handle, middleware.EmployeeOnly)
	}

	pvzGroup := handler.Group("pvz", app.AuthMiddleware().Middleware)
	{
		pvzGroup.POST("/:pvzId/close_last_reception", app.CloseReceptionHandler().Handle, middleware.EmployeeOnly)
		pvzGroup.POST("/:pvzId/delete_last_product", app.DeleteProductHandler().Handle, middleware.EmployeeOnly)
		pvzGroup.POST("", app.PostPointHandler().Handle, middleware.ModderatorOnly)
		pvzGroup.GET("", app.GetPointsHandler().Handle, middleware.EmployeeAndModerator)
	}
}
