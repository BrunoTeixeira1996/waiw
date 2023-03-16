package utils

import (
	"html/template"
	"log"
	"net/http"

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
	return func(w http.ResponseWriter, r *http.Request) {
		var movies []models.Movie

		if err := db.QueryAllFromMovies("select * from movies", &movies); err != nil {
			log.Fatal("Error while handling QueryAllFromMovies")
		}

		page := models.Page{
			Title: "models.Movies",
			Any:   movies,
		}
		baseTemplate.Execute(w, page)
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
