package utils

import (
	"html/template"
	"net/http"
)

// Struct that represents a webpage
type Page struct {
	Title string
}

// Handles "/"
func IndexHandle(baseTemplate *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page := Page{
			Title: "Index",
		}
		baseTemplate.Execute(w, page)
	}
}

// Handles "/movies"
func MoviesHandle(baseTemplate *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page := Page{
			Title: "Movies",
		}
		baseTemplate.Execute(w, page)
	}
}
