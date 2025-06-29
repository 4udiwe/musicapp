package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/4udiwe/musicshop/config"
	"github.com/4udiwe/musicshop/internal/api/post_albums"
	"github.com/4udiwe/musicshop/internal/repo/albums"
	albums_service "github.com/4udiwe/musicshop/internal/service/albums"
	"github.com/4udiwe/musicshop/pkg/validator"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

func InitDB(c *config.Config) *sql.DB {
	db, err := sql.Open("postgres", c.PostgresURL)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	fmt.Println("Successfully connected to PostgreSQL!")
	return db
}

func main() {

	c := config.LoadConfig()

	e := echo.New()
	e.Validator = validator.NewCustomValidator()
	e.POST("/albums", post_albums.New(albums_service.New(albums.New(InitDB(c)))).Handle)
	e.Logger.Fatal(e.Start(c.ServerPort))
}
