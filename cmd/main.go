package main

import (
	"context"
	"log"

	"github.com/4udiwe/musicshop/config"
	"github.com/4udiwe/musicshop/internal/api/post_albums"
	"github.com/4udiwe/musicshop/internal/database"
	"github.com/4udiwe/musicshop/internal/repo/albums"
	albums_service "github.com/4udiwe/musicshop/internal/service/albums"
	"github.com/4udiwe/musicshop/pkg/validator"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

func main() {
	ctx := context.Background()
	c := config.LoadConfig()

	pool, err := database.NewPgxPool(ctx, c)
	if err != nil {
		log.Fatalf("Failed to create pool: %v", err)
		return
	}
	defer pool.Close()

	if err := database.CreateAlbumsTable(ctx, pool); err != nil {
		log.Fatalf("Failed to create table: %v", err)
		return
	}

	e := echo.New()
	albumsRepository := albums.New(pool)
	albumsService := albums_service.New(albumsRepository)
	postHandler := post_albums.New(albumsService)

	e.Validator = validator.NewCustomValidator()
	e.POST("/albums", postHandler.Handle)
	e.Logger.Fatal(e.Start(c.ServerPort))
}
