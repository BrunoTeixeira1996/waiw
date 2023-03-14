# What Am I Watching?

Go webserver to store movies/series ratings by me and my gf

# TODO

- Create tables for db

	- movies table
      | id | title | image | sinopse | launch_date | imdb_rating | rottentomatoes_rating | view_date |

	- rating table
      | id | rating |

	- user table
	  | id | username |

	- movie_rating table
	  | id | movie_id | user_id | rating_id | comments |

- Connect with go into db
- Insert dummy data into the db
- Visualize data in website

- Create `uploadmovie` endpoint to upload a movie
  - Create a form and upload to db
  - The image will be saved in a `images` folder with a random generated name and that name will be saved in the database


## Useful stuff

- Use `create_db.sql` to create a clean dummy database

```console
$ cat create_db.sql | sqlite3 database.db
```

