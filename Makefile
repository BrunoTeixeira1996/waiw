db:
	rm database.db
	cat create_db.sql | sqlite3 database.db

build:
	go build -o waiw main.go

run:
	go run main.go

all: hello build
