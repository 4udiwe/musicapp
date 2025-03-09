package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var albums = []album{
	{ID: "1", Title: "first", Artist: "first artist", Price: 100.50},
	{ID: "2", Title: "second", Artist: "second artist", Price: 200.50},
	{ID: "3", Title: "third", Artist: "third artist", Price: 300.50},
}

func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

func postAlbum(c *gin.Context) {
	var na album
	if err := c.BindJSON(&na); err != nil {
		c.IndentedJSON(http.StatusBadRequest, na)
		return
	}

	albums = append(albums, na)
	c.IndentedJSON(http.StatusCreated, na)
}

func findAlbum(c *gin.Context) {

	param := c.Param("id")

	for _, album := range albums {
		if album.ID == param {
			c.IndentedJSON(http.StatusFound, album)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, param)
}

func deleteAlbumById(c *gin.Context) {
	param := c.Param("id")

	for i, album := range albums {
		if album.ID == param {
			albums = append(albums[:i], albums[i+1:]...)
			c.IndentedJSON(http.StatusNoContent, album)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, param)
}

func main() {
	router := gin.Default()

	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", findAlbum)
	router.POST("/albums", postAlbum)
	router.DELETE("/albums/:id", deleteAlbumById)

	router.Run("localhost:8080")
}
