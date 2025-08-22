package app

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/4udiwe/avito-pvz/config"
	api "github.com/4udiwe/avito-pvz/internal/api/http"
	"github.com/4udiwe/avito-pvz/internal/database"
	repo_point "github.com/4udiwe/avito-pvz/internal/repository/point"
	repo_product "github.com/4udiwe/avito-pvz/internal/repository/product"
	repo_reception "github.com/4udiwe/avito-pvz/internal/repository/reception"
	repo_user "github.com/4udiwe/avito-pvz/internal/repository/user"
	"github.com/4udiwe/avito-pvz/internal/service/point"
	"github.com/4udiwe/avito-pvz/internal/service/product"
	"github.com/4udiwe/avito-pvz/internal/service/reception"
	"github.com/4udiwe/avito-pvz/pkg/httpserver"
	"github.com/4udiwe/avito-pvz/pkg/postgres"
	"github.com/labstack/echo/v4"
)

type App struct {
	cfg       *config.Config
	interrupt <-chan os.Signal

	// DB
	postgres *postgres.Postgres

	// Echo
	echoHandler *echo.Echo

	// Repositories
	userRepo      *repo_user.Repository
	pointRepo     *repo_point.Repository
	productRepo   *repo_product.Repository
	receptionRepo *repo_reception.Repository

	// Handlers
	deleteProductHandler  api.Handler
	getPointsHandler      api.Handler
	closeReceptionHandler api.Handler
	postPointHandler      api.Handler
	postProductHandler    api.Handler
	postReceptionHandler  api.Handler

	// Services
	pointService     *point.Service
	productService   *product.Service
	receptionService *reception.Service
}

func New(configPath string) *App {
	cfg, err := config.New(configPath)
	if err != nil {
		log.Fatalf("app - New - config.New: %v", err)
	}

	initLogger(cfg.Log.Level)

	return &App{
		cfg: cfg,
	}
}

func (app *App) Start() {
	// Postgres
	log.Info("Connecting to PostgreSQL...")

	postgres, err := postgres.New(app.cfg.Postgres.URL, postgres.ConnAttempts(5))

	if err != nil {
		log.Fatalf("app - Start - Postgres failed:%v", err)
	}
	app.postgres = postgres

	defer postgres.Close()

	// Migrations
	if err := database.RunMigrations(context.Background(), app.postgres.Pool); err != nil {
		log.Errorf("app - Start - Migrations failed: %v", err)
	}

	// Server
	log.Info("Start server...")
	httpServer := httpserver.New(app.EchoHandler(), httpserver.Port(app.cfg.HTTP.Port))
	httpServer.Start()

	defer func() {
		if err := httpServer.Shutdown(); err != nil {
			log.Errorf("HTTP server shutdown error: %v", err)
		}
	}()

	select {
	case s := <-app.interrupt:
		log.Infof("app - Start - signal: %v", s)
	case err := <-httpServer.Notify():
		log.Errorf("app - Start - server error: %v", err)
	}

	log.Info("Shutting down...")
}
