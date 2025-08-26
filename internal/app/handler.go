package app

import (
	api "github.com/4udiwe/avito-pvz/internal/api/http"
	"github.com/4udiwe/avito-pvz/internal/api/http/delete_product"
	"github.com/4udiwe/avito-pvz/internal/api/http/get_points"
	"github.com/4udiwe/avito-pvz/internal/api/http/patch_reception"
	"github.com/4udiwe/avito-pvz/internal/api/http/post_dummy_login"
	"github.com/4udiwe/avito-pvz/internal/api/http/post_login"
	"github.com/4udiwe/avito-pvz/internal/api/http/post_point"
	"github.com/4udiwe/avito-pvz/internal/api/http/post_product"
	"github.com/4udiwe/avito-pvz/internal/api/http/post_reception"
	"github.com/4udiwe/avito-pvz/internal/api/http/post_refresh"
	"github.com/4udiwe/avito-pvz/internal/api/http/post_register"
)

func (app *App) DeleteProductHandler() api.Handler {
	if app.deleteProductHandler != nil {
		return app.deleteProductHandler
	}
	app.deleteProductHandler = delete_product.New(app.ProductService())
	return app.deleteProductHandler
}

func (app *App) GetPointsHandler() api.Handler {
	if app.getPointsHandler != nil {
		return app.getPointsHandler
	}
	app.getPointsHandler = get_points.New(app.PointService())
	return app.getPointsHandler
}

func (app *App) CloseReceptionHandler() api.Handler {
	if app.closeReceptionHandler != nil {
		return app.closeReceptionHandler
	}
	app.closeReceptionHandler = patch_reception.New(app.ReceptionService())
	return app.closeReceptionHandler
}

func (app *App) PostPointHandler() api.Handler {
	if app.postPointHandler != nil {
		return app.postPointHandler
	}
	app.postPointHandler = post_point.New(app.PointService())
	return app.postPointHandler
}

func (app *App) PostProductHandler() api.Handler {
	if app.postProductHandler != nil {
		return app.postProductHandler
	}
	app.postProductHandler = post_product.New(app.ProductService())
	return app.postProductHandler
}

func (app *App) PostReceptionHandler() api.Handler {
	if app.postReceptionHandler != nil {
		return app.postReceptionHandler
	}
	app.postReceptionHandler = post_reception.New(app.ReceptionService())
	return app.postReceptionHandler
}

func (app *App) PostDummyLoginHandler() api.Handler {
	if app.postDummyLoginHandler != nil {
		return app.postDummyLoginHandler
	}
	app.postDummyLoginHandler = post_dummy_login.New(app.UserService())
	return app.postDummyLoginHandler
}

func (app *App) PostLoginHandler() api.Handler {
	if app.postLoginHandler != nil {
		return app.postLoginHandler
	}
	app.postLoginHandler = post_login.New(app.UserService())
	return app.postLoginHandler
}

func (app *App) PostRegisterHandler() api.Handler {
	if app.postRegisterHandler != nil {
		return app.postRegisterHandler
	}
	app.postRegisterHandler = post_register.New(app.UserService())
	return app.postRegisterHandler
}

func (app *App) PostRefreshHandler() api.Handler {
	if app.postRefreshHandler != nil {
		return app.postRefreshHandler
	}
	app.postRefreshHandler = post_refresh.New(app.UserService())
	return app.postRefreshHandler
}
