package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	// "github.com/pelletier/go-toml/query"
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

// func searchAlbumById(id string) (int, *Album, error) {
// 	for i, album := range albums {
// 		if id == album.ID {
// 			return i, &albums[i], nil
// 		}
// 	}
// 	return -1, nil, &SearchError{Message: "album not found"}
// }

func getAlbums(c *gin.Context) {
	rows, err := db.Query("SELECT id, title, artist, price FROM albums")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var albums []Album
	for rows.Next() {
		var a Album
		err := rows.Scan(&a.ID, &a.Title, &a.Artist, &a.Price)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		albums = append(albums, a)
	}

	c.JSON(http.StatusOK, albums)
}

func postAlbum(c *gin.Context) {
	var newAlbum Album

	// 1. Парсим JSON из тела запроса в структуру Album
	if err := c.BindJSON(&newAlbum); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid JSON data",
			"details": err.Error(),
		})
		return
	}

	// 2. Валидация данных
	if newAlbum.Title == "" || newAlbum.Artist == "" || newAlbum.Price == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "All fields (title, artist, price) are required",
		})
		return
	}

	// 3. SQL-запрос для вставки (возвращает ID новой записи)
	query := `
        INSERT INTO albums (title, artist, price)
        VALUES ($1, $2, $3)
        RETURNING id
    `

	var albumID int
	err := db.QueryRow(query, newAlbum.Title, newAlbum.Artist, newAlbum.Price).Scan(&albumID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create album",
			"details": err.Error(),
		})
		return
	}

	// 4. Возвращаем созданный альбом с ID
	newAlbum.ID = strconv.Itoa(albumID)
	c.JSON(http.StatusCreated, newAlbum)
}

func findAlbum(c *gin.Context) {

	id := c.Param("id")

	query := `
		SELECT * FROM albums 
		WHERE id = $1
	`
	row := db.QueryRow(query, id)
	if row == nil {
		c.JSON(http.StatusNotFound, id)
		return
	}

	var album Album
	err := row.Scan(&album.ID, &album.Title, &album.Artist, &album.Price)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Unable to parse album data")
		return
	}
	c.IndentedJSON(http.StatusFound, album)
}

func deleteAlbumById(c *gin.Context) {
	id := c.Param("id")

	query := `
		DELETE FROM albums
		WHERE id = $1
	`
	_, err := db.Query(query, id)
	if err != nil {
		c.JSON(http.StatusNotFound, id)
		return
	}
	c.IndentedJSON(http.StatusFound, "Deleted album by id: "+id)
}

func createTable() {
	query := `CREATE TABLE IF NOT EXISTS albums (
		id SERIAL PRIMARY KEY,
		title VARCHAR(100) NOT NULL,
		artist VARCHAR(100) NOT NULL,
		price REAL
	);`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}
}

var db *sql.DB

func main() {
	// Инициализация БД
	db = initDB()
	defer db.Close()

	// Создаём таблицу (если её нет)
	createTable()

	router := gin.Default()

	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", findAlbum)
	router.POST("/albums", postAlbum)
	router.DELETE("/albums/:id", deleteAlbumById)

	router.Run(":8080")
}
