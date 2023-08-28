CREATE TABLE movies (
id INTEGER PRIMARY KEY AUTOINCREMENT,
title VARCHAR NOT NULL,
image VARCHAR NOT NULL,
sinopse VARCHAR NOT NULL,
genre VARCHAR NOT NULL,
imdb_rating VARCHAR NOT NULL,
launch_date INTEGER NOT NULL,
view_date INTEGER NOT NULL);

CREATE TABLE series (
id INTEGER PRIMARY KEY AUTOINCREMENT,
title VARCHAR NOT NULL,
image VARCHAR NOT NULL,
genre VARCHAR NOT NULL,
imdb_rating VARCHAR NOT NULL,
launch_date INTEGER NOT NULL);

CREATE TABLE animes (
id INTEGER PRIMARY KEY AUTOINCREMENT,
title VARCHAR NOT NULL,
image VARCHAR NOT NULL,
genre VARCHAR NOT NULL,
imdb_rating VARCHAR NOT NULL,
launch_date INTEGER NOT NULL);


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


CREATE TABLE category (
id INTEGER PRIMARY KEY AUTOINCREMENT,
name VARCHAR NOT NULL);


CREATE TABLE plan_to_watch (
id INTEGER PRIMARY KEY AUTOINCREMENT,
name VARCHAR,
category_id INTEGER,
FOREIGN KEY(category_id) REFERENCES category(id));


INSERT INTO category VALUES(1,'Movie');
INSERT INTO category VALUES(2,'Serie');
INSERT INTO category VALUES(3,'Anime');

INSERT INTO plan_to_watch VALUES(1,'Movie123',1);
INSERT INTO plan_to_watch VALUES(2,'Serie321',2);
INSERT INTO plan_to_watch VALUES(3,'Serie1234',2);
INSERT INTO plan_to_watch VALUES(4,'Anime123',3);
INSERT INTO plan_to_watch VALUES(5,'Movie123444',1);

INSERT INTO movies VALUES(1,'Movie1','image_name','sinopse1','romance','1', 1, 1);
INSERT INTO movies VALUES(2,'Movie2','image_name','sinopse2', 'horror','2', 2, 2);
INSERT INTO movies VALUES(3,'Movie3','image_name','sinopse3', 'romance','3',3, 3);
INSERT INTO movies VALUES(4,'Movie4','image_name','sinopse4', 'romance','4',4, 4);
INSERT INTO movies VALUES(5,'Movie5','image_name','sinopse5', 'romance','5',5, 5);
INSERT INTO movies VALUES(6,'Movie6','image_name','sinopse6', 'romance','6',6, 6);


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

INSERT INTO movie_ratings VALUES(1,1,1,2,'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.');
INSERT INTO movie_ratings VALUES(2,1,2,3,'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.');
INSERT INTO movie_ratings VALUES(3,2,1,1,'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.');
INSERT INTO movie_ratings VALUES(4,2,2,3,'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.');
INSERT INTO movie_ratings VALUES(5,3,1,2,'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.');
INSERT INTO movie_ratings VALUES(6,3,2,1,'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.');

