package utils

import (
	"html/template"
	"net/http"
)

// Struct that represents a webpage
type Page struct {
	Title string
	Any   any
}

type Movie struct {
	Id                    int
	Title                 string
	Image                 string
	Sinopse               string
	Imdb_Rating           string
	RottenTomatoes_Rating string
	Launch_Date           string
	View_Date             string
}

// Handles "/"
func IndexHandle(baseTemplate *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page := Page{
			Title: "Home",
		}
		baseTemplate.Execute(w, page)
	}
}

// Handles "/movies"
func MoviesHandle(baseTemplate *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		movies := &[]Movie{
			Movie{
				Id:                    1,
				Title:                 "1",
				Image:                 "1",
				Sinopse:               "1",
				Imdb_Rating:           "1",
				RottenTomatoes_Rating: "1",
				Launch_Date:           "1",
				View_Date:             "1",
			},
			Movie{
				Id:                    2,
				Title:                 "2",
				Image:                 "2",
				Sinopse:               "2",
				Imdb_Rating:           "2",
				RottenTomatoes_Rating: "2",
				Launch_Date:           "2",
				View_Date:             "2",
			},
			Movie{
				Id:                    3,
				Title:                 "3",
				Image:                 "3",
				Sinopse:               "3",
				Imdb_Rating:           "3",
				RottenTomatoes_Rating: "3",
				Launch_Date:           "3",
				View_Date:             "3",
			},
			Movie{
				Id:                    3,
				Title:                 "3",
				Image:                 "3",
				Sinopse:               "3",
				Imdb_Rating:           "3",
				RottenTomatoes_Rating: "3",
				Launch_Date:           "3",
				View_Date:             "3",
			},
			Movie{
				Id:                    4,
				Title:                 "4",
				Image:                 "4",
				Sinopse:               "4",
				Imdb_Rating:           "4",
				RottenTomatoes_Rating: "4",
				Launch_Date:           "4",
				View_Date:             "4",
			},
			Movie{
				Id:                    5,
				Title:                 "5",
				Image:                 "5",
				Sinopse:               "5",
				Imdb_Rating:           "5",
				RottenTomatoes_Rating: "5",
				Launch_Date:           "5",
				View_Date:             "5",
			},
		}

		page := Page{
			Title: "Movies",
			Any:   movies,
		}
		baseTemplate.Execute(w, page)
	}
}

// Handles "/movies"
func SeriesHandle(baseTemplate *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page := Page{
			Title: "Series",
		}
		baseTemplate.Execute(w, page)
	}
}
