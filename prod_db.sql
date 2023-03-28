CREATE TABLE movies (
id INTEGER PRIMARY KEY AUTOINCREMENT,
title VARCHAR NOT NULL,
image VARCHAR,
sinopse VARCHAR,
genre VARCHAR,
imdb_rating VARCHAR,
launch_date INTEGER,
view_date INTEGER NOT NULL);


CREATE TABLE ratings (
id INTEGER PRIMARY KEY AUTOINCREMENT,
value VARCHAR NOT NULL);


CREATE TABLE users (
id INTEGER PRIMARY KEY AUTOINCREMENT,
username VARCHAR NOT NULL);


CREATE TABLE movie_ratings (
id INTEGER PRIMARY KEY AUTOINCREMENT,
movie_id INTEGER,
user_id INTEGER,
rating_id INTEGER,
comments VARCHAR,
FOREIGN KEY(movie_id) REFERENCES movies(id),
FOREIGN KEY(user_id) REFERENCES users(id),
FOREIGN KEY(rating_id) REFERENCES ratings(id));


INSERT INTO ratings VALUES(1,'WTF');
INSERT INTO ratings VALUES(2,'Realy Realy Bad');
INSERT INTO ratings VALUES(3,'Realy Bad');
INSERT INTO ratings VALUES(4,'Bad');
INSERT INTO ratings VALUES(5,'Meh');
INSERT INTO ratings VALUES(6,'Not that bad');
INSERT INTO ratings VALUES(7,'Good');
INSERT INTO ratings VALUES(8,'Very good');
INSERT INTO ratings VALUES(9,'Awesome');
INSERT INTO ratings VALUES(10,'Masterpiece');


INSERT INTO users VALUES(1,'Bruno');
INSERT INTO users VALUES(2,'Rafaela');
