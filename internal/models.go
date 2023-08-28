package internal

import (
	"database/sql"
)

// Struct that represents a Db
type Db struct {
	Location string
	Type     string
	Con      *sql.DB
	Err      error
	Rows     *sql.Rows
	Result   sql.Result
}

// Struct that represents a webpage
type Page struct {
	Title  string
	Active map[string]string
	Any    any //movies, series, animes, ...
	Users  []User
	Error  any
}

type Upload struct {
	Title       string
	Sinopse     string
	Genre       string
	Imdb_Rating string
	Launch_Date string
	View_Date   string
	Image       string
	Category    string
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
	MovieRating []MovieRating
}

// Struct that represents an user
type User struct {
	Id       int
	Username string
}

// Struct that represents Movie Rating
type MovieRating struct {
	UserName string
	Rating   string
	Comments string
}

// Struct that represents a Category
type Category struct {
	Id   int
	Name string
}

// Struct that represents Plan to Watch
type Ptw struct {
	Id       int
	Name     string
	Category Category
}
