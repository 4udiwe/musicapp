package entity

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float64
	Genres []Genre
}
