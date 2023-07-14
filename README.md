# What Am I Watching?

Go webserver to store movies/series/animes ratings and comments

## Using localy

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

## Using with gokrazy

- Its possible to use waiw in [gokrazy](https://gokrazy.org/) as a go appliance (thats how I use it)
- My current setup is
  - LXC running debian server (proxmox) with a postgresql service running on port 5432 (hosting the `waiw` database)
  - Raspberry pi 4 running gokrazy

![image](https://github.com/BrunoTeixeira1996/waiw/assets/12052283/9a0426af-6093-48de-a732-319d2c977fbb)

- To build and deploy on gokrazy do the following
  - Add the repo to gokrazy instance
    - `git clone git@github.com:BrunoTeixeira1996/waiw.git && gok -i <YOUR_GOKRAZY_INSTANCE_ NAME> add ./waiw`
  - Edit the `config.json`
``` json
"github.com/BrunoTeixeira1996/waiw": {
    "CommandLineFlags": [
        "-gokrazy",
		"-ip='IP_OF_THE_SERVER_THAT_HOSTS_POSTGRESQL_DB'",
		"-user='YOUR_DB_USERNAME'",
		"-password='YOUR_DB_PASSWORD'",
		"-dbname='YOUR_DB_NAME'"
    ],
	"ExtraFilePaths": {
       "/etc/waiw/assets": "/home/brun0/Desktop/personal/waiw/assets",
	   "/etc/waiw/templates": "/home/brun0/Desktop/personal/waiw/templates"
    }
```
  - Update the gokrazy instance
    - `cd ~/gokrazy/<YOUR_GOKRAZY_INSTANCE_NAME> && gok update`


## Uploading

- You can upload using the upload button or the youtube link
- Use what is best for you but I recommend the youtube link since all images stay the same size


## Note

This project was done because me and my gf wanted to watch poorly rated movies and laugh on how bad they were so I deciced to create this in order to save our ratings and our comments so we can laugh even more in the future.

There are open issues that I plan to work on but I don't have a deadline for this since the project does what I want.


![image](https://user-images.githubusercontent.com/12052283/230114199-517c115a-8804-4031-9b7e-dda5487b5535.png)
