dev_db:
	rm -f dev_database.db
	cat dev_db.sql | sqlite3 dev_database.db

build:
	go build -o waiw main.go

run:
	go run main.go

all: hello build
