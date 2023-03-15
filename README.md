# What Am I Watching?

Go webserver to store movies/series ratings and comments by me and my gf

# TODO

- [] Create movie/<movie_id> page ()
  - When press movie image in `movies` endpoint navigate to this page `movie/<movie_id>`
- [] Connect with go into db
- [] Visualize data from db in website
- [] Create `uploadmovie` endpoint to upload a movie
  - Create a form and upload to db
  - The image will be saved in a `images` folder with a random generated name and that name will be saved in the database


# DONE

- [X] Create tables for db
	- movies table
      | id | title | image | sinopse | genre | imdb_rating | launch_date | view_date |

	- rating table
      | id | rating |

	- user table
	  | id | username |

	- movie_rating table
	  | id | movie_id | user_id | rating_id | comments |

- [X] Insert dummy data into the db
- [X] Build css and html base
- [X] Use dummy data to check if backend and frontend are communicating

## Useful stuff

- Use `create_db.sql` to create a clean dummy database

```console
$ cat create_db.sql | sqlite3 database.db
```

