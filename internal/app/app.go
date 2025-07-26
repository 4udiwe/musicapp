package app

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/4udiwe/musicshop/config"
	"github.com/4udiwe/musicshop/internal/api"
	"github.com/4udiwe/musicshop/internal/database"
	albums_repo "github.com/4udiwe/musicshop/internal/repo/albums"
	genres_repo "github.com/4udiwe/musicshop/internal/repo/genres"
	albums_service "github.com/4udiwe/musicshop/internal/service/albums"
	genres_service "github.com/4udiwe/musicshop/internal/service/genres"
	"github.com/4udiwe/musicshop/pkg/httpserver"
	"github.com/4udiwe/musicshop/pkg/postgres"
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
	albumsRepo *albums_repo.Repository
	genresRepo *genres_repo.Repository

	// Handlers
	deleteAlbumHandler api.Handler
	deleteGenreHandler api.Handler

	getAlbumHandler  api.Handler
	getAlbumsHandler api.Handler
	getGenresHandler api.Handler

	postAlbumHandler      api.Handler
	postAlbumGenreHandler api.Handler
	postGenreHandler      api.Handler

	// Services
	albumsService *albums_service.Service
	genresService *genres_service.Service
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
		log.Fatalf("app - Start - Migrations failed: %v", err)
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
