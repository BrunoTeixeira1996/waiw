# What Am I Watching?

Go webserver to store movies/series ratings and comments by me and my gf

# TODO

- [] Fix selection in navbar, so when we are in /movies just select movies as blue and not home
- [X] Change movie?id=1 layout
  - [X] Comment once, if there's a comment for bruno and rafaela just don't show comment and rating html tags
  - [] Make editable comment and rating
  - [] Make deletable comment and rating
- [] Refactor code
  - [] change from fmt.Print to log
  - [] reduce code dupplication
- [] Create tests

# DONE

- [X] Create tables for db
- [X] Insert dummy data into the db
- [X] Build css and html base
- [X] Use dummy data to check if backend and frontend are communicating
- [X] Connect with go into db
- [X] Visualize data from db in website
- [X] Create movie.html page (https://i.imgur.com/O6GyAM4.png)
  - [X] When press movie image in `movies` endpoint navigate to this page `movie?id=<movie_id>`
- [X] Insert comments and ratings in movies
- [X] Create Makefile
- [X] Add correct images for correct movies
  - When inserting an image in db, use the name of an image
  - When showing that image, use `/assets/images/<correct_name>` because for now its hardcoded
- [X] Create `/upload` endpoint to upload a movie
  - [X] Create a form and upload to db
  - [X] The image will be saved in a `/assets/images/` folder with a random generated name and that name will be saved in the database
  - [X] Create jquery to validate all inputs are filled
- [X] Bug when inserting comments
  - If bruno inserts comments twice, rafaela can't insert more so I need to validate this
- [X] Create valid database and debug database
- [X] Make flag for debuging as well as a database for debugging
- [X] Upload button could be download image from (https://img.youtube.com/vi/9V0gyIZuUcg/hqdefault.jpg) since all movies are on youtube
- [X] Fix information in movies directory
