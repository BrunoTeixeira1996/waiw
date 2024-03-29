package handles

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/BrunoTeixeira1996/waiw/internal/metandmod"
)

// Handles "/"
func IndexHandle(baseTemplate *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page := metandmod.Page{
			Title: "Home",
		}
		page.LoadActiveEndpoint("Home")

		baseTemplate.Execute(w, page)
	}
}

// Handles "/movies"
func MoviesHandle(baseTemplate *template.Template) http.HandlerFunc {
	var (
		movies       []metandmod.Movie
		movieRating  []metandmod.MovieRating
		users        []metandmod.User
		title        string
		alertDanger  string
		emptyInputs  bool
		hasCommented bool
	)

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
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
				if err := metandmod.QueryMovie(movieId, &title, &movies, movieRating); err != nil {
					log.Println("Error while querying a movie:", err)
					return
				}
				log.Println("Opened movie:", title)

				// Get users in database
				if err := metandmod.GetAvailableUsers(&users); err != nil {
					log.Println("Error while querying users:", err)
					return
				}

			} else {
				// List all movies
				if err := metandmod.QueryAllFromMovie("select * from movies", &movies); err != nil {
					log.Println("Error while handling QueryAllFromMovies:", err)
				}
				title = "Movies"
			}

			page := metandmod.Page{
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

		case http.MethodPost:
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
				log.Println("There are empty attributes when inserting a rting:", emptyAttr)
				return
			}

			var user metandmod.User

			if regexp.MustCompile(`\d`).MatchString(movieId) {
				if err := metandmod.SetUser("select * from users where username = $1", author, &user); err != nil {
					log.Println("Error while seting user:", err)
					return
				}

				// Verify if this user already commented
				userHasCommented := func() bool {
					yes, err := metandmod.UserAlreadyCommented("select movie_ratings.id from movie_ratings, movies, users where movie_ratings.movie_id = movies.id and movie_ratings.user_id = users.id and movies.id = $1 and users.id = $2", movieId, user.Id)
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
				if err := metandmod.GenericQuery("inserting movie comments", "insert into movie_ratings (movie_id, user_id, rating_id, comments) VALUES ($1,$2,$3,$4)", movieId, user.Id, choosenRating, comments); err != nil {
					log.Println("Error while inserting movie comment:", err)
					return
				}
			}

			log.Printf("User %s added comment to movieid %s\n", user.Username, movieId)

			// Redirects to GET
			http.Redirect(w, r, r.Header.Get("Referer"), 302)
		}
	}
}

// Handles "/upload"
func UploadHandle(baseTemplate *template.Template) http.HandlerFunc {
	var allowedImageTypes = map[string]int{
		"image/png":  1,
		"image/jpeg": 2,
		"image/jpg":  3,
	}

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			page := metandmod.Page{
				Title: "Upload",
			}

			page.LoadActiveEndpoint("Upload")

			baseTemplate.Execute(w, page)

		case http.MethodPost:
			upload := metandmod.Upload{
				Title:       r.FormValue("title"),
				Sinopse:     r.FormValue("area_1"),
				Genre:       r.FormValue("genre"),
				Imdb_Rating: r.FormValue("imdb"),
				Launch_Date: r.FormValue("ldate"),
				View_Date:   r.FormValue("vdate"),
				Category:    r.Form["categories"][0],
			}

			// validate all fields
			if err := upload.ValidateFieldsInUpload(upload.Category); err != nil {
				alertDanger := fmt.Sprintf("<p class='alert alert-danger'>  %s </p>", err)
				page := metandmod.Page{
					Title: "Upload",
					Error: template.HTML(alertDanger),
				}
				baseTemplate.Execute(w, page)
				log.Println("Error while validating all fields in the upload:", err)
				return
			}

			// if everything works fine, download the image otherwise we would have dead images even when we failed uploading a movie
			// Get image name and save in /assets/image/ folder
			// it returns the path of the image saved and the image link provided in the upload endpoint
			im := func() (string, string) {
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
						return "", ""
					}
					defer imageFile.Close()

					if _, ok := allowedImageTypes[handler.Header.Get("Content-Type")]; !ok {
						log.Println("Error, file type not allowed when uploading image")
						return "", ""
					}
				}

				path, err := os.Getwd()
				if err != nil {
					log.Println("Error while getting the current path:", err)
					return "", ""
				}

				newImage, err := os.CreateTemp(path+"/assets/images", "upload-*.png")
				if err != nil {
					log.Println("Error while creating the new image:", err)
					return "", ""
				}
				defer newImage.Close()

				// Image will be based on upload button or link input
				switch {

				case len(imageLink) > 0:
					res, err := http.Get(imageLink)
					if err != nil {
						log.Println("Error while querying the image:", imageLink)
						return "", ""
					}

					defer res.Body.Close()
					image = res.Body

				case len(handler.Filename) > 0:
					image = imageFile
				}

				imageBytes, err := io.ReadAll(image)
				if err != nil {
					log.Println("Error while reading the contents of the uploaded image:", err)
					return "", ""
				}
				if _, err := newImage.Write(imageBytes); err != nil {
					log.Println("Error while writting the new image:", err)
					return "", ""
				}
				im := strings.Split(newImage.Name(), "/")
				return im[len(im)-1], imageLink
			}

			// If its a movie we output the image path in order to get the movie comments
			// otherwise we just output the image link since we don't rage series nor animes
			if upload.Category == "Movie" {
				upload.Image, _ = im()
			} else {
				_, upload.Image = im()
			}

			if hasEmpty, emptyAttr := upload.HasEmptyAttr(upload.Category); hasEmpty {
				alertDanger := fmt.Sprintf("<p class='alert alert-danger'> Missing: %s </p>", emptyAttr)
				page := metandmod.Page{
					Title: "Upload",
					Error: template.HTML(alertDanger),
				}
				baseTemplate.Execute(w, page)
				log.Println("There are empty attributes when uploading:", emptyAttr)
				return
			}

			switch upload.Category {
			case "Movie":
				if err := metandmod.GenericQuery("inserting new movie", "insert into movies (title, image, sinopse, genre, imdb_rating, launch_date, view_date) VALUES ($1,$2,$3,$4,$5,$6,$7)", upload.Title, upload.Image, upload.Sinopse, upload.Genre, upload.Imdb_Rating, upload.Launch_Date, upload.View_Date); err != nil {
					log.Println("Error while inserting new movie:", err)
					return
				}
				log.Println("Added movie:", upload.Title)

			case "Serie":
				if err := metandmod.GenericQuery("inserting new serie", "insert into series (title, image, genre, imdb_rating, launch_date) VALUES ($1,$2,$3,$4,$5)", upload.Title, upload.Image, upload.Genre, upload.Imdb_Rating, upload.Launch_Date); err != nil {
					log.Println("Error while inserting new serie:", err)
					return
				}
				log.Println("Added serie:", upload.Title)

			case "Anime":
				log.Println("Not implemented yet")
			}

			page := metandmod.Page{
				Title: "Upload",
			}
			baseTemplate.Execute(w, page)
		}
	}
}

// Handles "/series"
func SeriesHandle(baseTemplate *template.Template) http.HandlerFunc {
	var series []metandmod.Serie
	return func(w http.ResponseWriter, r *http.Request) {
		if err := metandmod.QueryAllFromSeries("select * from series", &series); err != nil {
			log.Println("Error while handling QueryAllFromSeries:", err)
		}

		page := metandmod.Page{
			Title: "Series",
			Any:   series,
		}
		series = nil
		page.LoadActiveEndpoint("Series")
		baseTemplate.Execute(w, page)
	}
}

// Handles "/ptw"
func PtwHandle(baseTemplate *template.Template) http.HandlerFunc {
	var (
		alertDanger string
		emptyInputs bool
	)

	return func(w http.ResponseWriter, r *http.Request) {
		var (
			valid bool
			err   error
		)

		switch r.Method {
		case http.MethodGet:
			// Get list of plan to watch from database
			var sptw []metandmod.Ptw

			// Checks if theres a cookie about an error so we can display that in the html
			c, _ := r.Cookie("error_cookie")

			if emptyInputs {
				alertDanger = fmt.Sprintf("<p class='alert alert-danger'> Missing: %s </p>", c.Value)
				cookie := http.Cookie{Name: "error_cookie", Value: "", Expires: time.Unix(0, 0), HttpOnly: true}
				http.SetCookie(w, &cookie)
				emptyInputs = false
			}

			if err := metandmod.GetPlanToWatch(&sptw); err != nil {
				log.Println("Error while querying ptw:", err)
				return
			}

			ptwTemp := struct {
				Movies []metandmod.Ptw
				Series []metandmod.Ptw
				Animes []metandmod.Ptw
			}{}

			for _, v := range sptw {
				switch v.Category.Name {
				case "Movie":
					ptwTemp.Movies = append(ptwTemp.Movies, v)
				case "Serie":
					ptwTemp.Series = append(ptwTemp.Series, v)
				case "Anime":
					ptwTemp.Animes = append(ptwTemp.Animes, v)
				}
			}

			page := metandmod.Page{
				Title: "Plan to Watch",
				Any:   ptwTemp,
				Error: template.HTML(alertDanger),
			}
			page.LoadActiveEndpoint("PlanToWatch")

			baseTemplate.Execute(w, page)

			// Cleans alert cookie
			alertDanger = ""

		case http.MethodPost:
			ptwname := r.FormValue("ptwname")
			categoryId := r.Form["categories"]

			hasEmptyAttrs := func() (bool, string) {
				if ptwname == "" {
					return true, "Name"
				}
				if len(categoryId) == 0 {
					return true, "Category"
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
				log.Println("There are empty attributes while inserting plan to watch:", emptyAttr)
				return
			}

			if err := metandmod.GenericQuery("inserting plan to watch", "insert into plan_to_watch (name,category_id) VALUES ($1,$2)", ptwname, categoryId[0]); err != nil {
				log.Println("Error while inserting new plan to watch:", err)
				return
			}

			log.Printf("Added new plan to watch %s from the ui\n", ptwname)
			http.Redirect(w, r, r.Header.Get("Referer"), 302)

		case http.MethodDelete:
			d := json.NewDecoder(r.Body)
			d.DisallowUnknownFields() // error if user sends extra data

			deletePtw := struct {
				Id     *string `json:"id"`
				Origin *string `json:"origin"`
			}{}

			if err := d.Decode(&deletePtw); err != nil {
				// bad JSON or unrecognized json field
				log.Println("Error while decoding plan to watch from DELETE:", err)
				return
			}

			// Check if theres more than what we want
			if d.More() {
				log.Println("Extraneous data after JSON object from DELETE in plan to watch")
				return
			}

			if valid, err = metandmod.DeletePlanToWatch(*deletePtw.Id, "ui"); err != nil {
				log.Println("Error while deleting plan to watch (ui):", err)
				return
			}

			if valid {
				log.Printf("Deleted %s from plan to watch (%s)\n", *deletePtw.Id, *deletePtw.Origin)
				http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
				return

			}

		}
	}
}
