package utils

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

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
	var (
		movies      []models.Movie
		movieRating []models.MovieRating
		title       string
		alertDanger string
		emptyInputs bool
	)

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			// Checks if theres a cookie about an error so we can display that in the html
			c, _ := r.Cookie("error_cookie")

			if emptyInputs {
				alertDanger = fmt.Sprintf("<p class='alert alert-danger'> Missing: %s </p>", c.Value)
				cookie := http.Cookie{Name: "error_cookie", Value: "", Expires: time.Unix(0, 0), HttpOnly: true}
				http.SetCookie(w, &cookie)
				emptyInputs = false
			}

			// Get movieId
			movieId := r.URL.Query().Get("id")

			// List respective movie based on movieId
			if movieId != "" {
				if err := db.QueryMovie(movieId, title, &movies, movieRating); err != nil {
					fmt.Println("Error while querying a movie")
					return
				}

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
				Error: template.HTML(alertDanger),
			}

			baseTemplate.Execute(w, page)

			// Cleaning slices since they are pointers, or they will get dup values
			movies = nil
			movieRating = nil

		case "POST":
			// Gather user inputs
			comments := r.FormValue("area_1")
			author := r.FormValue("group_1")
			choosenRating := r.Form["ratings"][0]
			movieId := r.URL.Query().Get("id")
			hasEmptyAttrs := func() (bool, string) {
				if comments == "" {
					return true, "Comments"
				}
				if author == "" {
					return true, "Author"
				}
				if choosenRating == "" {
					return true, "Rating"
				}
				if movieId == "" {
					return true, "Movie ID"
				}

				return false, ""
			}

			// Check if all inputs are filled
			if hasEmpty, emptyAttr := hasEmptyAttrs(); hasEmpty {
				// Set cookie so GET knows there's an error
				emptyInputs = true
				cookie := http.Cookie{Name: "error_cookie", Value: emptyAttr}
				http.SetCookie(w, &cookie)
				http.Redirect(w, r, r.Header.Get("Referer"), 302)
				return
			}

			var user models.User

			if regexp.MustCompile(`\d`).MatchString(movieId) {
				if err := db.SetUser("select * from users where username = ?", author, &user); err != nil {
					fmt.Println(err)
				}
				// Insert in database the comments and ratings
				if err := db.InsertMovieComments("insert into movie_ratings (movie_id, user_id, rating_id, comments) VALUES (?,?,?,?)", movieId, user.Id, choosenRating, comments); err != nil {
					fmt.Println("Error while inserting movie comment %w", err)
				}
			}

			// Redirects to GET
			http.Redirect(w, r, r.Header.Get("Referer"), 302)
		}
	}
}

// Handles "/upload"
func UploadHandle(baseTemplate *template.Template, db *models.Db) http.HandlerFunc {
	var allowedImageTypes = map[string]int{
		"image/png":  1,
		"image/jpeg": 2,
		"image/jpg":  3,
	}

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			page := models.Page{
				Title: "Upload",
			}
			baseTemplate.Execute(w, page)

		case "POST":
			// Get image name and save in /assets/image/ folder
			im := func() string {
				image, handler, err := r.FormFile("myFile")
				if err != nil {
					fmt.Println("Error retrieved the image file ", err)
					return ""
				}
				defer image.Close()
				if _, ok := allowedImageTypes[handler.Header.Get("Content-Type")]; !ok {
					fmt.Println("Error, file type not allowed")
					return ""
				}

				path, err := os.Getwd()
				if err != nil {
					fmt.Println("Error while getting the current path ", err)
					return ""
				}

				newImage, err := os.CreateTemp(path+"/assets/images", "upload-*.png")
				if err != nil {
					fmt.Println("Error while creating the new image ", err)
					return ""
				}
				defer newImage.Close()

				imageBytes, err := ioutil.ReadAll(image)
				if err != nil {
					fmt.Println("Error while reading the contents of the uploaded image ", err)
					return ""
				}
				if _, err := newImage.Write(imageBytes); err != nil {
					fmt.Println("Error while writting the new image ", err)
					return ""
				}
				im := strings.Split(newImage.Name(), "/")
				return im[len(im)-1]
			}

			movie := models.Movie{
				Title:       r.FormValue("title"),
				Image:       im(),
				Sinopse:     r.FormValue("area_1"),
				Genre:       r.FormValue("genre"),
				Imdb_Rating: r.FormValue("imdb"),
				Launch_Date: r.FormValue("ldate"),
				View_Date:   r.FormValue("vdate"),
			}

			if hasEmpty, emptyAttr := movie.HasEmptyAttr(); hasEmpty {
				alertDanger := fmt.Sprintf("<p class='alert alert-danger'> Missing: %s </p>", emptyAttr)
				page := models.Page{
					Title: "Upload",
					Error: template.HTML(alertDanger),
				}
				baseTemplate.Execute(w, page)
				return
			}

			if err := db.InsertMovieComments("insert into movies (title, image, sinopse, genre, imdb_rating, launch_date, view_date) VALUES (?,?,?,?,?,?,?)", movie.Title, movie.Image, movie.Sinopse, movie.Genre, movie.Imdb_Rating, movie.Launch_Date, movie.View_Date); err != nil {
				fmt.Println("Error while inserting new movie %w", err)
				return
			}

			page := models.Page{
				Title: "Upload",
			}
			baseTemplate.Execute(w, page)
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
