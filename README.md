# What Am I Watching?

Go webserver to store movies/series/animes ratings and comments

# Using

- Clone the repo
- Create the database based on `/dev/dev_db.sql` (just remove the INSERTs)
- Run

``` console
$ git clone git@github.com:BrunoTeixeira1996/waiw.git
$ cd waiw
$ cat /dev/dev_db.sql | sqlite3 /<your database path>/database.db
$ go build
$ ./waiw -db "<database path>"
```

## Uploading

- You can upload using the upload button or the youtube link
- Use what is best for you but I recommend the youtube link since all images stay the same size


## Note

This project was done because me and my gf wanted to watch poorly rated movies and laugh on how bad they were so I deciced to create this in order to save our ratings and our comments so we can laugh even more in the future.

There are open issues that I plan to work on but I don't have a deadline for this since the project does what I want.


![image](https://user-images.githubusercontent.com/12052283/230114199-517c115a-8804-4031-9b7e-dda5487b5535.png)
