package utils

import (
	"fmt"

	"github.com/BrunoTeixeira1996/waiw/models"
	_ "github.com/mattn/go-sqlite3"
)

func DbDEBUG() {
	var movies []models.Movie
	db := &models.Db{}

	if err := db.QueryAllFromMovies("select * from movies", &movies); err != nil {
		fmt.Println(err)
	}

	fmt.Println(movies)

}
