package app

import (
	"github.com/4udiwe/avito-pvz/internal/api/http/middleware"
	"github.com/4udiwe/avito-pvz/internal/auth"
	"github.com/4udiwe/avito-pvz/pkg/hasher"
)

func (app *App) Auth() *auth.Auth {
	if app.auth != nil {
		return app.auth
	}
	app.auth = auth.New()
	return app.auth
}

func (app *App) Hasher() *hasher.BcryptHasher {
	if app.hasher != nil {
		return app.hasher
	}
	app.hasher = hasher.New()
	return app.hasher
}

func (app *App) AuthMiddleware() *middleware.AuthMiddleware {
	if app.authMW != nil {
		return app.authMW
	}
	app.authMW = middleware.New(app.Auth())
	return app.authMW
}
