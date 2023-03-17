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
	var movieRating []models.MovieRating
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

					if err := db.QueryCommentsAndRatings("select users.username, ratings.value, movie_ratings.comments from ratings, movie_ratings, movies, users where ratings.id = movie_ratings.rating_id and movies.id = movie_ratings.movie_id and users.id = movie_ratings.user_id and movie_id = ?", &movieRating, movieId); err != nil {
						fmt.Printf("Error while QueryCommentsAndRatings for movie id=%s\n", movieId)
					}
				}

				movies[0].MovieRating = movieRating
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
			movieRating = nil

		case "POST":
			// Insert in database
			comments := r.FormValue("area_1")
			author := r.FormValue("group_1")
			choosenRating := r.Form["ratings"][0]

			movieId := r.URL.Query().Get("id")
			var user models.User

			if regexp.MustCompile(`\d`).MatchString(movieId) {
				if err := db.SetUser("select * from users where username = ?", author, &user); err != nil {
					fmt.Println(err)
				}

				if err := db.InsertMovieComments("insert into movie_ratings (movie_id, user_id, rating_id, comments) VALUES (?,?,?,?)", movieId, user.Id, choosenRating, comments); err != nil {
					fmt.Println("Error while inserting movie comment %w", err)
				}
			}

			http.Redirect(w, r, r.Header.Get("Referer"), 302)

		}
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
