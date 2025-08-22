package app

import (
	repo_point "github.com/4udiwe/avito-pvz/internal/repository/point"
	repo_product "github.com/4udiwe/avito-pvz/internal/repository/product"
	repo_reception "github.com/4udiwe/avito-pvz/internal/repository/reception"
	repo_user "github.com/4udiwe/avito-pvz/internal/repository/user"
	"github.com/4udiwe/avito-pvz/pkg/postgres"
)

func (app *App) Postgres() *postgres.Postgres {
	return app.postgres
}

func (app *App) UserRepo() *repo_user.Repository {
	if app.userRepo != nil {
		return app.userRepo
	}
	app.userRepo = repo_user.New(app.Postgres())
	return app.userRepo
}

func (app *App) PointRepo() *repo_point.Repository {
	if app.pointRepo != nil {
		return app.pointRepo
	}
	app.pointRepo = repo_point.New(app.Postgres())
	return app.pointRepo
}

func (app *App) ProductRepo() *repo_product.Repository {
	if app.productRepo != nil {
		return app.productRepo
	}
	app.productRepo = repo_product.New(app.Postgres())
	return app.productRepo
}

func (app *App) ReceptionRepo() *repo_reception.Repository {
	if app.receptionRepo != nil {
		return app.receptionRepo
	}
	app.receptionRepo = repo_reception.New(app.Postgres())
	return app.receptionRepo
}
