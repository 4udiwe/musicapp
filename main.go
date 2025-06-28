package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

type SearchError struct {
	Message string
}

func (e *SearchError) Error() string {
	return e.Message
}

var albums = []Album{
	{ID: "1", Title: "first", Artist: "first artist", Price: 100.50},
	{ID: "2", Title: "second", Artist: "second artist", Price: 200.50},
	{ID: "3", Title: "third", Artist: "third artist", Price: 300.50},
}

func searchAlbumById(id string) (int, *Album, error) {
	for i, album := range albums {
		if id == album.ID {
			return i, &albums[i], nil
		}
	}
	return -1, nil, &SearchError{Message: "album not found"}
}

func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

func postAlbum(c *gin.Context) {
	var na Album
	if err := c.BindJSON(&na); err != nil {
		c.IndentedJSON(http.StatusBadRequest, na)
		return
	}

	albums = append(albums, na)
	c.IndentedJSON(http.StatusCreated, na)
}

func findAlbum(c *gin.Context) {

	param := c.Param("id")

	_, album, err := searchAlbumById(param)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, param)
		return
	}
	c.IndentedJSON(http.StatusFound, album)
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

func updateAlbumById(c *gin.Context) {
	param := c.Param("id")

	index, album, err := searchAlbumById(param)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, param)
		return
	}

	if err := c.BindJSON(album); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	albums[index] = *album

	c.IndentedJSON(http.StatusAccepted, album)
}

func main() {
	router := gin.Default()

	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", findAlbum)
	router.POST("/albums", postAlbum)
	router.DELETE("/albums/:id", deleteAlbumById)
	router.PUT("/albums/:id", updateAlbumById)

	router.Run(":8080")
}
