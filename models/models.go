package models

import (
	"database/sql"
)

// Struct that represents a Db
type Db struct {
	Con  *sql.DB
	Err  error
	Rows *sql.Rows
}

// Struct that represents a webpage
type Page struct {
	Title string
	Any   any
}

// Struct that represents a movie
type Movie struct {
	Id          int
	Title       string
	Image       string
	Sinopse     string
	Genre       string
	Imdb_Rating string
	Launch_Date string
	View_Date   string
}
