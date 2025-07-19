package albums

import (
	"github.com/4udiwe/musicshop/internal/entity"
	"github.com/samber/lo"
)

type albumsGenreRow struct {
	ID         int64    `db:"id"`
	Title      string   `db:"title"`
	Artist     string   `db:"artist"`
	Price      float64  `db:"price"`
	GenreIDs   []int64  `db:"genre_ids"`
	GenreNames []string `db:"genre_names"`
}

func (r *albumsGenreRow) convertRowToAlbum() entity.Album {
	genres := lo.Map(r.GenreIDs, func(id int64, i int) entity.Genre {
		return entity.Genre{
			ID:   id,
			Name: r.GenreNames[i],
		}
	})

	return entity.Album{
		ID:     r.ID,
		Title:  r.Title,
		Artist: r.Artist,
		Price:  r.Price,
		Genres: genres,
	}
}

func convertRowsToAlbums(rows []albumsGenreRow) []entity.Album {
	return lo.Map(rows, func(row albumsGenreRow, _ int) entity.Album {
		return row.convertRowToAlbum()
	})
}
