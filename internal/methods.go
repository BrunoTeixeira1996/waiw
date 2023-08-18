package internal

import (
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// Choose active endpoint to use selection in navbar
func (p *Page) LoadActiveEndpoint(endpoint string) error {
	switch endpoint {
	case "Home":
		p.Active = map[string]string{
			"Home": "active",
		}
	case "Movies":
		p.Active = map[string]string{
			"Movies": "active",
		}
	case "Series":
		p.Active = map[string]string{
			"Series": "active",
		}
	case "Upload":
		p.Active = map[string]string{
			"Upload": "active",
		}
	case "PlanToWatch":
		p.Active = map[string]string{
			"PlanToWatch": "active",
		}

	default:
		return fmt.Errorf("%s does not exist\n", endpoint)
	}

	return nil
}

// Connect to database
func (c *Db) Connect() error {
	c.Con, c.Err = sql.Open(c.Type, c.Location)
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
func (c *Db) QueryMovie(movieId string, title *string, movies *[]Movie, movieRating []MovieRating) error {
	if regexp.MustCompile(`\d`).MatchString(movieId) {
		if err := c.QueryAllFromMovies("select * from movies where id = $1", movies, movieId); err != nil {
			return fmt.Errorf("Error while QueryAllFromMovies for movie id=%s\n", movieId)
		}

		// Gathers comments and ratings about specific movie
		if err := c.QueryCommentsAndRatings("select users.username, ratings.value, movie_ratings.comments from ratings, movie_ratings, movies, users where ratings.id = movie_ratings.rating_id and movies.id = movie_ratings.movie_id and users.id = movie_ratings.user_id and movie_id = $1", &movieRating, movieId); err != nil {
			return fmt.Errorf("Error while QueryCommentsAndRatings for movie id=%s\n", movieId)
		}
	}

	// Adds movieRating to the rating of a certain movie
	(*movies)[0].MovieRating = movieRating
	// Adds title of the page according to the respective movie
	*title = (*movies)[0].Title

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

// Get available users
func (c *Db) GetAvailableUsers(users *[]User) error {
	if err := c.Connect(); err != nil {
		return err
	}
	defer c.Con.Close()

	if c.Rows, c.Err = c.Con.Query("select * from users"); c.Err != nil {
		return fmt.Errorf("Error while getting available users from database: %w", c.Err)
	}

	for c.Rows.Next() {
		var user User
		if c.Err = c.Rows.Scan(&user.Id, &user.Username); c.Err == sql.ErrNoRows {
			return fmt.Errorf("Error while scanning rows: %w", c.Err)
		}

		*users = append(*users, user)
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

// Gets CategoryName from Id
func (c *Db) GetCategoryName(categoryId int, category *Category) error {
	if err := c.Connect(); err != nil {
		return err
	}
	defer c.Con.Close()

	if c.Err = c.Con.QueryRow("select * from category where id = $1", categoryId).Scan(&category.Id, &category.Name); c.Err != nil {
		return fmt.Errorf("Error while querying the catory with id %d: %w", categoryId, c.Err)
	}

	return nil
}

// Gets CategoryId from name
func (c *Db) GetCategoryId(categoryName string, category *Category) error {
	if err := c.Connect(); err != nil {
		return err
	}
	defer c.Con.Close()

	if c.Err = c.Con.QueryRow("select * from category where name = $1", categoryName).Scan(&category.Id, &category.Name); c.Err != nil {
		return fmt.Errorf("Error while querying the catory with name %s: %w", categoryName, c.Err)
	}

	return nil
}

// Gets  plan to watch entries
func (c *Db) GetPlanToWatch(sptw *[]Ptw) error {
	if err := c.Connect(); err != nil {
		return err
	}
	defer c.Con.Close()

	if c.Rows, c.Err = c.Con.Query("select * from plan_to_watch order by category_id"); c.Err == sql.ErrNoRows {
		return fmt.Errorf("Error while getting plan to watch entries: %w", c.Err)
	}

	for c.Rows.Next() {
		var ptw Ptw
		var category Category

		if c.Err = c.Rows.Scan(&ptw.Id, &ptw.Name, &ptw.Category.Id); c.Err == sql.ErrNoRows {
			return fmt.Errorf("Error while scanning rows: %w", c.Err)
		}

		if err := c.GetCategoryName(ptw.Category.Id, &category); err != nil {
			return err
		}

		ptw.Category = category

		*sptw = append(*sptw, ptw)
	}

	return nil
}

func (c *Db) InsertPlanToWatch(q string, params ...any) error {
	if err := c.Connect(); err != nil {
		return err
	}
	defer c.Con.Close()

	if c.Result, c.Err = c.Con.Exec(q, params...); c.Err != nil {
		return fmt.Errorf("Error while inserting a new plan to watch: %w", c.Err)
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

// Validates fields of movie for the Upload Handle
func (m *Movie) ValidateFieldsInUpload() error {
	genreHasNumber := regexp.MustCompile(`\d`).MatchString(m.Genre)
	if genreHasNumber {
		return fmt.Errorf("Genre must be a string, not a number")
	}
	if _, err := strconv.Atoi(m.Imdb_Rating); err != nil {
		return fmt.Errorf("Imdb Rating must be a number, not a string")
	}

	currentYear := time.Now().Year()
	intLaunchDate, _ := strconv.Atoi(m.Launch_Date)
	intViewDate, _ := strconv.Atoi(m.View_Date)

	if intLaunchDate > currentYear {
		return fmt.Errorf("Not a valid launch date year")
	}
	if intViewDate > currentYear || intViewDate < intLaunchDate {
		return fmt.Errorf("Not a valid view date year")
	}
	return nil
}
