package models

import (
	"database/sql"
	"fmt"
)

// Connect to database
func (c *Db) Connect() error {
	c.Con, c.Err = sql.Open("sqlite3", "/home/brun0/Desktop/personal/waiw/database.db")
	if c.Err != nil {
		return fmt.Errorf("Error while connecting to database: %w", c.Err)
	}

	return nil
}

// Query all info from movies
func (c *Db) QueryAllFromMovies(q string, movies *[]Movie, params ...any) error {
	if err := c.Connect(); err != nil {
		return err
	}
	defer c.Con.Close()

	if c.Rows, c.Err = c.Con.Query(q, params...); c.Err != nil {
		return fmt.Errorf("Error while doing query: %w", c.Err)
	}

	for c.Rows.Next() {
		var m Movie
		if c.Err = c.Rows.Scan(
			&m.Id,
			&m.Title,
			&m.Image,
			&m.Sinopse,
			&m.Genre,
			&m.Imdb_Rating,
			&m.Launch_Date,
			&m.View_Date,
		); c.Err == sql.ErrNoRows {
			return fmt.Errorf("Error while scanning rows: %w", c.Err)
		}

		*movies = append(*movies, m)
	}

	return nil
}

// Query all info about rating and comments
func (c *Db) QueryCommentsAndRatings(q string, movieRatings *[]MovieRating, params ...any) error {
	if err := c.Connect(); err != nil {
		return err
	}
	defer c.Con.Close()

	if c.Rows, c.Err = c.Con.Query(q, params...); c.Err != nil {
		return fmt.Errorf("Error while doing query: %w", c.Err)
	}

	for c.Rows.Next() {
		var r MovieRating
		if c.Err = c.Rows.Scan(
			&r.UserName,
			&r.Rating,
			&r.Comments,
		); c.Err == sql.ErrNoRows {
			return fmt.Errorf("Error while scanning rows: %w", c.Err)
		}

		*movieRatings = append(*movieRatings, r)
	}

	return nil
}

// Set user values
func (c *Db) SetUser(q string, username string, user *User) error {
	if err := c.Connect(); err != nil {
		return err
	}
	defer c.Con.Close()

	if c.Rows, c.Err = c.Con.Query(q, username); c.Err != nil {
		return fmt.Errorf("Error while querying for user: %w", c.Err)
	}

	for c.Rows.Next() {
		if c.Err = c.Rows.Scan(&user.Id, &user.Username); c.Err == sql.ErrNoRows {
			return fmt.Errorf("Error while scanning rows: %w", c.Err)
		}
	}

	return nil
}

// Insert comments into a movie
func (c *Db) InsertMovieComments(q string, params ...any) error {
	if err := c.Connect(); err != nil {
		return err
	}
	defer c.Con.Close()

	if c.Result, c.Err = c.Con.Exec(q, params...); c.Err != nil {
		return fmt.Errorf("Error while inserting comment in movie: %w", c.Err)
	}
	return nil
}

// Insert new movie
func (c *Db) InsertNewMovie(q string, params ...any) error {
	if err := c.Connect(); err != nil {
		return err
	}
	defer c.Con.Close()

	if c.Result, c.Err = c.Con.Exec(q, params...); c.Err != nil {
		return fmt.Errorf("Error while inserting a new movie: %w", c.Err)
	}
	return nil
}

// Check if any of movie field is empty
func (m *Movie) HasEmptyAttr() (bool, string) {
	if m.Title == "" {
		return true, "Title"
	}
	if m.Image == "" {
		return true, "Image"
	}
	if m.Sinopse == "" {
		return true, "Sinopse"
	}
	if m.Genre == "" {
		return true, "Genre"
	}
	if m.Imdb_Rating == "" {
		return true, "Imdb_Rating"
	}
	if m.Launch_Date == "" {
		return true, "Launch_Date"
	}
	if m.View_Date == "" {
		return true, "View_Date"
	}

	return false, ""
}
