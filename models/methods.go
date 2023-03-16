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

// Query database
func (c *Db) QueryAllFromMovies(q string, movies *[]Movie) error {
	if err := c.Connect(); err != nil {
		return err
	}
	defer c.Con.Close()

	if c.Rows, c.Err = c.Con.Query(q); c.Err != nil {
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
