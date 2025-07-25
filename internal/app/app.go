package app

import (
	"context"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/4udiwe/musicshop/config"
	"github.com/4udiwe/musicshop/internal/api"
	"github.com/4udiwe/musicshop/internal/database"
	albums_repo "github.com/4udiwe/musicshop/internal/repo/albums"
	genres_repo "github.com/4udiwe/musicshop/internal/repo/genres"
	albums_service "github.com/4udiwe/musicshop/internal/service/albums"
	genres_service "github.com/4udiwe/musicshop/internal/service/genres"
	"github.com/4udiwe/musicshop/pkg/httpserver"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type App struct {
	cfg       *config.Config
	interrupt <-chan os.Signal

	// DB
	pgxPool *pgxpool.Pool

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

	pgxConf, err := pgxpool.ParseConfig(app.cfg.Postgres.URL)
	if err != nil {
		log.Fatalf("failed to parse conn string: %v", err)
	}

	retryAttempts := 5
	retryDelay := 5 * time.Second
	var pool *pgxpool.Pool

	timeout := app.cfg.Postgres.ConnectTimeout
	if timeout == 0 {
		timeout = 5 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for i := 0; i < retryAttempts; i++ {
		pool, err = pgxpool.NewWithConfig(ctx, pgxConf)
		if err == nil {
			if err = pool.Ping(ctx); err != nil {
				log.Printf("Connection established but ping failed: %v", err)
				continue
			}
			app.pgxPool = pool
			break
		}

		if i < retryAttempts-1 {
			log.Printf("Attempt %d failed: %v. Retrying in %v...", i+1, err, retryDelay)
			time.Sleep(retryDelay)
		}
	}

	if err != nil {
		log.Fatalf("failed to connect to PostgreSQL after %d attempts: %v", retryAttempts, err)
	}

	defer func() {
		if pool != nil {
			log.Info("Closing PostgreSQL connection pool...")
			pool.Close()
		}
	}()

	// Migrations

	if err := database.RunMigrations(context.Background(), app.pgxPool); err != nil {
		log.Fatalf("Migrations failed: %v", err)
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
