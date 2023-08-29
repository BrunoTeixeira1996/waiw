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

// Connect to database once
func InitDb(dbType string, location string) error {
	d, err := sql.Open(dbType, location)
	if err != nil {
		return fmt.Errorf("Error while connecting to database: %w", err)
	}

	dbCon = d

	return nil
}

// Query all info from movie
func QueryAllFromMovie(q string, movies *[]Movie, params ...any) error {
	var (
		rows *sql.Rows
		err  error
	)

	if rows, err = dbCon.Query(q, params...); err != nil {
		return fmt.Errorf("Error while doing query: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var m Movie
		if err = rows.Scan(
			&m.Id,
			&m.Title,
			&m.Image,
			&m.Sinopse,
			&m.Genre,
			&m.Imdb_Rating,
			&m.Launch_Date,
			&m.View_Date,
		); err == sql.ErrNoRows {
			return fmt.Errorf("Error while scanning rows: %w", err)
		}

		*movies = append(*movies, m)
	}

	return nil
}

// Query a movie
func QueryMovie(movieId string, title *string, movies *[]Movie, movieRating []MovieRating) error {
	if regexp.MustCompile(`\d`).MatchString(movieId) {
		if err := QueryAllFromMovie("select * from movies where id = $1", movies, movieId); err != nil {
			return fmt.Errorf("Error while QueryAllFromMovies for movie id=%s\n", movieId)
		}

		// Gathers comments and ratings about specific movie
		if err := QueryCommentsAndRatings("select users.username, ratings.value, movie_ratings.comments from ratings, movie_ratings, movies, users where ratings.id = movie_ratings.rating_id and movies.id = movie_ratings.movie_id and users.id = movie_ratings.user_id and movie_id = $1", &movieRating, movieId); err != nil {
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
func QueryCommentsAndRatings(q string, movieRatings *[]MovieRating, params ...any) error {
	var (
		rows *sql.Rows
		err  error
	)

	if rows, err = dbCon.Query(q, params...); err != nil {
		return fmt.Errorf("Error while doing query: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var r MovieRating
		if err = rows.Scan(
			&r.UserName,
			&r.Rating,
			&r.Comments,
		); err == sql.ErrNoRows {
			return fmt.Errorf("Error while scanning rows: %w", err)
		}

		*movieRatings = append(*movieRatings, r)
	}

	return nil
}

// Get available users
func GetAvailableUsers(users *[]User) error {
	var (
		rows *sql.Rows
		err  error
	)

	if rows, err = dbCon.Query("select * from users"); err != nil {
		return fmt.Errorf("Error while getting available users from database: %w", err)
	}

	for rows.Next() {
		var user User
		if err = rows.Scan(&user.Id, &user.Username); err == sql.ErrNoRows {
			return fmt.Errorf("Error while scanning rows: %w", err)
		}

		*users = append(*users, user)
	}

	return nil
}

// Set user values
func SetUser(q string, username string, user *User) error {
	var (
		rows *sql.Rows
		err  error
	)

	if rows, err = dbCon.Query(q, username); err != nil {
		return fmt.Errorf("Error while querying for user: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&user.Id, &user.Username); err == sql.ErrNoRows {
			return fmt.Errorf("Error while scanning rows: %w", err)
		}
	}

	return nil
}

// Insert comments into a movie
func InsertMovieComments(q string, params ...any) error {
	if _, err := dbCon.Exec(q, params...); err != nil {
		return fmt.Errorf("Error while inserting comment in movie: %w", err)
	}
	return nil
}

// Insert new entry (movie,serie or anime)
func InsertNewEntry(q string, params ...any) error {
	if _, err := dbCon.Exec(q, params...); err != nil {
		return fmt.Errorf("Error while inserting a new movie: %w", err)
	}
	return nil
}

// Check if user already commented
func UserAlreadyCommented(q string, params ...any) (bool, error) {
	var (
		rows *sql.Rows
		err  error
	)

	if rows, err = dbCon.Query(q, params...); err != nil {
		return false, fmt.Errorf("Error while doing query: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		return true, nil
	}

	return false, nil
}

// Gets CategoryName from Id
func GetCategoryName(categoryId int, category *Category) error {
	if err := dbCon.QueryRow("select * from category where id = $1", categoryId).Scan(&category.Id, &category.Name); err != nil {
		return fmt.Errorf("Error while querying the catory with id %d: %w", categoryId, err)
	}

	return nil
}

// Gets CategoryId from name
func GetCategoryId(categoryName string, category *Category) error {
	if err := dbCon.QueryRow("select * from category where name = $1", categoryName).Scan(&category.Id, &category.Name); err != nil {
		return fmt.Errorf("Error while querying the catory with name %s: %w", categoryName, err)
	}

	return nil
}

func GetPlanToWatch(sptw *[]Ptw) error {
	var (
		rows *sql.Rows
		err  error
	)

	if rows, err = dbCon.Query("select * from plan_to_watch order by category_id"); err == sql.ErrNoRows {
		return fmt.Errorf("Error while getting plan to watch entries: %w", err)
	}

	for rows.Next() {
		var ptw Ptw
		var category Category

		if err = rows.Scan(&ptw.Id, &ptw.Name, &ptw.Category.Id); err == sql.ErrNoRows {
			return fmt.Errorf("Error while scanning rows: %w", err)
		}

		if err := GetCategoryName(ptw.Category.Id, &category); err != nil {
			return err
		}

		ptw.Category = category

		*sptw = append(*sptw, ptw)
	}

	return nil
}

func InsertPlanToWatch(q string, params ...any) error {
	if _, err := dbCon.Exec(q, params...); err != nil {
		return fmt.Errorf("Error while inserting a new plan to watch: %w", err)
	}
	return nil
}

func DeletePlanToWatch(name string, origin string) (bool, error) {
	var recordId sql.NullInt64
	switch origin {
	case "ui":
		if err := dbCon.QueryRow("delete from plan_to_watch where id = $1 returning id", name).Scan(&recordId); err != nil {
			return false, err
		}
	case "api":
		if err := dbCon.QueryRow("delete from plan_to_watch where name = $1 returning id", name).Scan(&recordId); err != nil {
			return false, err
		}
	}

	return recordId.Valid, nil
}

// Query all info from series
func QueryAllFromSeries(q string, series *[]Serie, params ...any) error {
	var (
		rows *sql.Rows
		err  error
	)

	if rows, err = dbCon.Query(q, params...); err != nil {
		return fmt.Errorf("Error while doing query: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var s Serie
		if err = rows.Scan(
			&s.Id,
			&s.Title,
			&s.Image,
			&s.Genre,
			&s.Imdb_Rating,
			&s.Launch_Date,
		); err == sql.ErrNoRows {
			return fmt.Errorf("Error while scanning rows: %w", err)
		}

		*series = append(*series, s)
	}

	return nil
}

// Check if any of upload field is empty
func (u *Upload) HasEmptyAttr(category string) (bool, string) {
	if u.Category == "" {
		return true, "Category"
	}

	switch category {
	case "Movie":
		if u.Title == "" {
			return true, "Title"
		}
		if u.Image == "" {
			return true, "Image"
		}
		if u.Sinopse == "" {
			return true, "Sinopse"
		}
		if u.Genre == "" {
			return true, "Genre"
		}
		if u.Imdb_Rating == "" {
			return true, "Imdb_Rating"
		}
		if u.Launch_Date == "" {
			return true, "Launch_Date"
		}
		if u.View_Date == "" {
			return true, "View_Date"
		}

	default:
		if u.Title == "" {
			return true, "Title"
		}
		if u.Image == "" {
			return true, "Image"
		}
		if u.Genre == "" {
			return true, "Genre"
		}
		if u.Imdb_Rating == "" {
			return true, "Imdb_Rating"
		}
		if u.Launch_Date == "" {
			return true, "Launch_Date"
		}
	}
	return false, ""
}

// Validates fields for the Upload Handle
func (u *Upload) ValidateFieldsInUpload(category string) error {
	genreHasNumber := regexp.MustCompile(`\d`).MatchString(u.Genre)
	if genreHasNumber {
		return fmt.Errorf("Genre must be a string, not a number")
	}
	if _, err := strconv.Atoi(u.Imdb_Rating); err != nil {
		return fmt.Errorf("Imdb Rating must be a number, not a string")
	}

	currentYear := time.Now().Year()
	intLaunchDate, _ := strconv.Atoi(u.Launch_Date)
	if intLaunchDate > currentYear {
		return fmt.Errorf("Not a valid launch date year")
	}

	if category == "Movie" {
		intViewDate, _ := strconv.Atoi(u.View_Date)
		if intViewDate > currentYear || intViewDate < intLaunchDate {
			return fmt.Errorf("Not a valid view date year")
		}
	}
	return nil
}
