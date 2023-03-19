# What Am I Watching?

Go webserver to store movies/series ratings and comments by me and my gf

# TODO

- [] Create `/upload` endpoint to upload a movie
  - Create a form and upload to db
  - The image will be saved in a `/assets/images/` folder with a random generated name and that name will be saved in the database
  - Create jquery to validate all inputs are filled
  - Check this https://tutorialedge.net/golang/go-file-upload-tutorial/
- [] Add correct images for correct movies
  - When inserting an image in db, use the name of an image
  - When showing that image, use `/assets/images/<correct_name>` because for now its hardcoded
- [] Change movie?id=1 layout
  - Comment once, if there's a comment for bruno and rafaela just don't show comment and rating html tags
  - Make editable comment and rating
  - Make deletable comment and rating
- [] Create Makefile
- [] Create tests

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
- [X] Connect with go into db
- [X] Visualize data from db in website
- [X] Create movie.html page (https://i.imgur.com/O6GyAM4.png)
  - [X] When press movie image in `movies` endpoint navigate to this page `movie?id=<movie_id>`
- [X] Insert comments and ratings in movies

## Useful stuff

- Use `create_db.sql` to create a clean dummy database

```console
$ cat create_db.sql | sqlite3 database.db
```

