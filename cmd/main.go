package main

import (
	"context"
	"log"

	"github.com/4udiwe/musicshop/config"
	"github.com/4udiwe/musicshop/internal/api/delete_album_by_id"
	delete_genres_by_id "github.com/4udiwe/musicshop/internal/api/delete_genre_by_id"
	"github.com/4udiwe/musicshop/internal/api/get_album_by_id"
	"github.com/4udiwe/musicshop/internal/api/get_albums"
	"github.com/4udiwe/musicshop/internal/api/get_genres"
	"github.com/4udiwe/musicshop/internal/api/post_albums"
	"github.com/4udiwe/musicshop/internal/api/post_albums_genre"
	"github.com/4udiwe/musicshop/internal/api/post_genre"
	"github.com/4udiwe/musicshop/internal/database"
	"github.com/4udiwe/musicshop/internal/repo/albums"
	"github.com/4udiwe/musicshop/internal/repo/genres"
	albums_service "github.com/4udiwe/musicshop/internal/service/albums"
	genres_service "github.com/4udiwe/musicshop/internal/service/genres"
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

	if err := database.Create(ctx, pool); err != nil {
		log.Fatalf("Failed to create table: %v", err)
		return
	}

	e := echo.New()
	e.Validator = validator.NewCustomValidator()

	// albums
	albumsRepository := albums.New(pool)
	albumsService := albums_service.New(albumsRepository)

	postAlbumHandler := post_albums.New(albumsService)
	getAllAlbumsHandler := get_albums.New(albumsService)
	getAlbumByIdHandler := get_album_by_id.New(albumsService)
	deleteAlbumHandler := delete_album_by_id.New(albumsService)

	e.POST("/albums", postAlbumHandler.Handle)
	e.GET("/albums", getAllAlbumsHandler.Handle)
	e.GET("/albums/:id", getAlbumByIdHandler.Handle)
	e.DELETE("/albums/:id", deleteAlbumHandler.Handle)

	// genres
	genresRepository := genres.New(pool)
	genresSrivce := genres_service.New(genresRepository)

	postGenreHandler := post_genre.New(genresSrivce)
	deleteGenreHandler := delete_genres_by_id.New(genresSrivce)
	getAllGenresHandler := get_genres.New(genresSrivce)
	postAlbumsGenresHandler := post_albums_genre.New(genresSrivce)

	e.POST("/genres", postGenreHandler.Handle)
	e.DELETE("/genres", deleteGenreHandler.Handle)
	e.GET("/genres", getAllGenresHandler.Handle)
	e.POST("/albums/genres", postAlbumsGenresHandler.Handle)

	e.Logger.Fatal(e.Start(c.ServerPort))
}
