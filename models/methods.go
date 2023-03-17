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

// Query all info from ratings
func (c *Db) QueryAllFromRatings(q string, ratings *[]MovieRating, params ...any) error {
	if err := c.Connect(); err != nil {
		return err
	}
	defer c.Con.Close()

	if c.Rows, c.Err = c.Con.Query(q, params...); c.Err != nil {
		return fmt.Errorf("Error while doing query: %w", c.Err)
	}

	for c.Rows.Next() {
		var rating MovieRating
		if c.Err = c.Rows.Scan(
			&rating.UserName,
			&rating.Rating,
			&rating.Comments,
		); c.Err == sql.ErrNoRows {
			return fmt.Errorf("Error while scanning rows: %w", c.Err)
		}

		*ratings = append(*ratings, rating)
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
