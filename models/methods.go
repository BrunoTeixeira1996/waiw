package models

import (
	"database/sql"
	"fmt"
	"regexp"
)

// Connect to database
func (c *Db) Connect() error {
	c.Con, c.Err = sql.Open("sqlite3", "/home/brun0/Desktop/personal/waiw/dev_database.db")
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

	defer c.Rows.Close()

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

// Query a movie
func (c *Db) QueryMovie(movieId string, title string, movies *[]Movie, movieRating []MovieRating) error {
	if regexp.MustCompile(`\d`).MatchString(movieId) {
		if err := c.QueryAllFromMovies("select * from movies where id = ?", movies, movieId); err != nil {
			return fmt.Errorf("Error while QueryAllFromMovies for movie id=%s\n", movieId)
		}

		// Gathers comments and ratings about specific movie
		if err := c.QueryCommentsAndRatings("select users.username, ratings.value, movie_ratings.comments from ratings, movie_ratings, movies, users where ratings.id = movie_ratings.rating_id and movies.id = movie_ratings.movie_id and users.id = movie_ratings.user_id and movie_id = ?", &movieRating, movieId); err != nil {
			return fmt.Errorf("Error while QueryCommentsAndRatings for movie id=%s\n", movieId)
		}
	}

	// Adds movieRating to the rating of a certain movie
	(*movies)[0].MovieRating = movieRating
	// Adds title of the page according to the respective movie
	title = (*movies)[0].Title

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

	defer c.Rows.Close()

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

	defer c.Rows.Close()
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

// Check if user already commented
func (c *Db) UserAlreadyCommented(q string, params ...any) (bool, error) {
	if err := c.Connect(); err != nil {
		return false, err
	}
	defer c.Con.Close()

	if c.Rows, c.Err = c.Con.Query(q, params...); c.Err != nil {
		return false, fmt.Errorf("Error while doing query: %w", c.Err)
	}

	defer c.Rows.Close()

	for c.Rows.Next() {
		return true, nil
	}

	return false, nil
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
