package utils

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
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
		page.LoadActiveEndpoint("Home")

		baseTemplate.Execute(w, page)
	}
}

// Handles "/movies"
func MoviesHandle(baseTemplate *template.Template, db *models.Db) http.HandlerFunc {
	var (
		movies       []models.Movie
		movieRating  []models.MovieRating
		users        []models.User
		title        string
		alertDanger  string
		emptyInputs  bool
		hasCommented bool
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

			if hasCommented {
				alertDanger = fmt.Sprintf("<p class='alert alert-danger'> User already comment this </p>")
				cookie := http.Cookie{Name: "error_cookie", Value: "", Expires: time.Unix(0, 0), HttpOnly: true}
				http.SetCookie(w, &cookie)
				hasCommented = false
			}

			// Get movieId
			movieId := r.URL.Query().Get("id")

			// List respective movie based on movieId
			if movieId != "" {
				if err := db.QueryMovie(movieId, &title, &movies, movieRating); err != nil {
					log.Println("Error while querying a movie:", err)
					return
				}
				// Get users in database
				if err := db.GetAvailableUsers(&users); err != nil {
					log.Println("Error while querying users:", err)
					return
				}

			} else {
				// List all movies
				if err := db.QueryAllFromMovies("select * from movies", &movies); err != nil {
					log.Println("Error while handling QueryAllFromMovies:", err)
				}
				title = "Movies"
			}

			page := models.Page{
				Title: title,
				Any:   movies,
				Users: users,
				Error: template.HTML(alertDanger),
			}
			page.LoadActiveEndpoint("Movies")

			baseTemplate.Execute(w, page)

			// Cleaning slices since they are pointers, or they will get dup values as well as the alert
			func() {
				alertDanger = ""
				movies = nil
				movieRating = nil
				users = nil
			}()

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
					log.Println("Error while seting user:", err)
					return
				}

				// Verify if this user already commented
				userHasCommented := func() bool {
					yes, err := db.UserAlreadyCommented("select movie_ratings.id from movie_ratings, movies, users where movie_ratings.movie_id = movies.id and movie_ratings.user_id = users.id and movies.id = ? and users.id = ?", movieId, user.Id)
					if err != nil {
						log.Println("Error while checking if user already commented on movie:", err)
						return false
					}

					if yes {
						fmt.Println("User already commented")
						cookie := http.Cookie{Name: "error_cookie", Value: "User already commented this"}
						http.SetCookie(w, &cookie)
						return true
					}
					return false
				}

				if userHasCommented() {
					hasCommented = true
					http.Redirect(w, r, r.Header.Get("Referer"), 302)
					return
				}

				// Insert in database the comments and ratings
				if err := db.InsertMovieComments("insert into movie_ratings (movie_id, user_id, rating_id, comments) VALUES (?,?,?,?)", movieId, user.Id, choosenRating, comments); err != nil {
					log.Println("Error while inserting movie comment:", err)
					return
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

			page.LoadActiveEndpoint("Upload")

			baseTemplate.Execute(w, page)

		case "POST":
			// Get image name and save in /assets/image/ folder
			im := func() string {
				var (
					imageLink string
					imageFile multipart.File
					handler   *multipart.FileHeader
					image     io.Reader
					err       error
				)

				// Prefer youtube link
				if len(r.FormValue("limg")) > 0 {
					imageLink = r.FormValue("limg")

					// Otherwise, use the upload button
				} else {
					imageFile, handler, err = r.FormFile("myFile")
					if err != nil {
						log.Println("Error retrieved the image file:", err)
						return ""
					}
					defer imageFile.Close()

					if _, ok := allowedImageTypes[handler.Header.Get("Content-Type")]; !ok {
						log.Println("Error, file type not allowed")
						return ""
					}
				}

				path, err := os.Getwd()
				if err != nil {
					log.Println("Error while getting the current path:", err)
					return ""
				}

				newImage, err := os.CreateTemp(path+"/assets/images", "upload-*.png")
				if err != nil {
					log.Println("Error while creating the new image:", err)
					return ""
				}
				defer newImage.Close()

				// Image will be based on upload button or link input
				switch {

				case len(imageLink) > 0:
					res, err := http.Get(imageLink)
					if err != nil {
						log.Printf("Error while querying %s\n", imageLink)
					}

					defer res.Body.Close()
					image = res.Body

				case len(handler.Filename) > 0:
					image = imageFile
				}

				imageBytes, err := ioutil.ReadAll(image)
				if err != nil {
					log.Println("Error while reading the contents of the uploaded image:", err)
					return ""
				}
				if _, err := newImage.Write(imageBytes); err != nil {
					log.Println("Error while writting the new image:", err)
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
				log.Println("Error while inserting new movie:", err)
				return
			}

			page := models.Page{
				Title: "Upload",
			}
			baseTemplate.Execute(w, page)
		}
	}
}

// Handles "/series"
func SeriesHandle(baseTemplate *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page := models.Page{
			Title: "Series",
		}
		page.LoadActiveEndpoint("Series")
		baseTemplate.Execute(w, page)

	}
}
