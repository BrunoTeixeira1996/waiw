dev_db:
	rm -f dev_database.db
	cat dev_db.sql | sqlite3 dev_database.db

dev:
	go run . -debug

prod:
	go run .

build:
	go build -o waiw main.go
