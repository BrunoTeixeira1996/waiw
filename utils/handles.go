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
	var ratings []models.MovieRating
	var title string

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			movieId := r.URL.Query().Get("id")

			// List respective movie
			if movieId != "" {
				if regexp.MustCompile(`\d`).MatchString(movieId) {
					if err := db.QueryAllFromMovies("select * from movies where id = ?", &movies, movieId); err != nil {
						fmt.Printf("Error while QueryAllFromMovies for movie id=%s\n", movieId)
					}

					if err := db.QueryAllFromRatings("select users.username, ratings.value, movie_ratings.comments from ratings, movie_ratings, movies, users where ratings.id = movie_ratings.rating_id and movies.id = movie_ratings.movie_id and users.id = movie_ratings.user_id and movie_id = ?", &ratings, movieId); err != nil {
						fmt.Printf("Error while QueryAllFromRatings for movie id=%s\n", movieId)
					}
				}
				movies[0].MovieRating = ratings
				title = movies[0].Title

			} else {
				// List all movies
				if err := db.QueryAllFromMovies("select * from movies", &movies); err != nil {
					fmt.Println("Error while handling QueryAllFromMovies")
				}
				title = "Movies"
			}

			page := models.Page{
				Title: title,
				Any:   movies,
			}
			baseTemplate.Execute(w, page)

			// Since I am using a pointer, I need to clean this slice
			// or there will be dups
			movies = nil
			ratings = nil

		case "POST":
			// comment := r.FormValue("area_1")
			// author := r.FormValue("group_1")

			// TODO: insert in database

			movieId := r.URL.Query().Get("id")
			if regexp.MustCompile(`\d`).MatchString(movieId) {
				if err := db.QueryAllFromMovies("select * from movies where id = ?", &movies, movieId); err != nil {
					fmt.Printf("Error while QueryAllFromMovies for movie id=%s\n", movieId)
				}

				if err := db.QueryAllFromRatings("select users.username, ratings.value, movie_ratings.comments from ratings, movie_ratings, movies, users where ratings.id = movie_ratings.rating_id and movies.id = movie_ratings.movie_id and users.id = movie_ratings.user_id and movie_id = ?", &ratings, movieId); err != nil {
					fmt.Printf("Error while QueryAllFromRatings for movie id=%s\n", movieId)
				}

			}
			movies[0].MovieRating = ratings

			page := models.Page{
				Title: movies[0].Title,
				Any:   movies,
			}
			baseTemplate.Execute(w, page)
		}

		// Since I am using a pointer, I need to clean this slice
		// or there will be dups
		movies = nil
		ratings = nil
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
