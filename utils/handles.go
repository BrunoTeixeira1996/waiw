package utils

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"

	"github.com/BrunoTeixeira1996/waiw/models"
)

// Handles "/"
func IndexHandle(baseTemplate *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page := models.Page{
			Title: "Home",
		}
		baseTemplate.Execute(w, page)
	}
}

// Handles "/movies"
func MoviesHandle(baseTemplate *template.Template, db *models.Db) http.HandlerFunc {
	var movies []models.Movie

	return func(w http.ResponseWriter, r *http.Request) {
		movieId := r.URL.Query().Get("id")

		// List respective movie
		if movieId != "" {
			if regexp.MustCompile(`\d`).MatchString(movieId) {
				if err := db.QueryAllFromMovies("select * from movies where id = ?", &movies, movieId); err != nil {
					fmt.Printf("Error while QueryAllFromMovies for movie id=%s\n", movieId)
				}
			}

		} else {
			// List all movies
			if err := db.QueryAllFromMovies("select * from movies", &movies); err != nil {
				fmt.Println("Error while handling QueryAllFromMovies")
			}
		}

		page := models.Page{
			Title: "models.Movies",
			Any:   movies,
		}
		baseTemplate.Execute(w, page)

		// Since I am using a pointer, I need to clean this slice
		// or there will be dups
		movies = nil
	}
}

// Handles "/movies"
func SeriesHandle(baseTemplate *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page := models.Page{
			Title: "Series",
		}
		baseTemplate.Execute(w, page)
	}
}
