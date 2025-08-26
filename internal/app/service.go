package app

import (
	"github.com/4udiwe/avito-pvz/internal/service/point"
	"github.com/4udiwe/avito-pvz/internal/service/product"
	"github.com/4udiwe/avito-pvz/internal/service/reception"
	"github.com/4udiwe/avito-pvz/internal/service/user"
)

func (app *App) PointService() *point.Service {
	if app.pointService != nil {
		return app.pointService
	}
	app.pointService = point.New(app.PointRepo(), app.ReceptionRepo(), app.ProductRepo(), app.Postgres())
	return app.pointService
}

func (app *App) ProductService() *product.Service {
	if app.productService != nil {
		return app.productService
	}
	app.productService = product.New(app.ProductRepo(), app.ReceptionRepo(), app.Postgres())
	return app.productService
}

func (app *App) ReceptionService() *reception.Service {
	if app.receptionService != nil {
		return app.receptionService
	}
	app.receptionService = reception.New(app.ReceptionRepo(), app.Postgres())
	return app.receptionService
}

func (app *App) UserService() *user.Service {
	if app.userService != nil {
		return app.userService
	}
	app.userService = user.New(app.UserRepo(), app.Postgres(), app.Auth(), app.Hasher())
	return app.userService
}
