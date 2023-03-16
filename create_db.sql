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
rating VARCHAR NOT NULL);


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

INSERT INTO movies VALUES(1,'Movie1','image_path','sinopse1','romance','1', 1, 1);
INSERT INTO movies VALUES(2,'Movie2','image_path','sinopse2', 'horror','2', 2, 2);
INSERT INTO movies VALUES(3,'Movie3','image_path','sinopse3', 'romance','3',3, 3);

INSERT INTO ratings VALUES(1,'Very very bad');
INSERT INTO ratings VALUES(2,'Very bad');
INSERT INTO ratings VALUES(3,'Bad');

INSERT INTO users VALUES(1,'Bruno');
INSERT INTO users VALUES(2,'Rafaela');

INSERT INTO movie_ratings VALUES(1,1,1,2,'so bad');
INSERT INTO movie_ratings VALUES(2,1,2,3,'bad');
INSERT INTO movie_ratings VALUES(3,2,1,1,'good');
INSERT INTO movie_ratings VALUES(4,2,2,3,'awesome');
INSERT INTO movie_ratings VALUES(5,3,1,2,'meh');
INSERT INTO movie_ratings VALUES(6,3,2,1,'not bad');

