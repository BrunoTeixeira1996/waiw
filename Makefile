dev_db:
	rm -f /dev/dev_database.db
	cat /dev/dev_db.sql | sqlite3 /dev/dev_database.db

rundev:
	go run "$(CURDIR)/cmd/waiw/waiw.go" -db "$(CURDIR)/dev/dev_database.db"

prod:
	go run "$(CURDIR)/cmd/waiw/waiw.go" -db "$(CURDIR)/prod/prod_database.db"

build:
	go build -o waiw main.go
