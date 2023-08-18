package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	cp "github.com/otiai10/copy"

	"github.com/BrunoTeixeira1996/waiw/internal"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// Handles the exit signal
func handleExit(exit chan bool) {
	ch := make(chan os.Signal, 5)
	signal.Notify(ch, os.Interrupt)
	<-ch
	log.Println("Closing web server")
	exit <- true
}

// Function that logs every request
func requestLogger(targetMux http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		targetMux.ServeHTTP(w, r)

		// log request by who(IP address)
		requesterIP := r.RemoteAddr

		log.Printf(
			"%s\t\t%s\t\t%s\t\t%v",
			r.Method,
			r.RequestURI,
			requesterIP,
			time.Since(start),
		)
	})
}

// Starts the web server
func startServer(currentPath string, databasePath string, debugFlag bool, dbType string) error {

	db := &internal.Db{}

	db.Location = databasePath
	db.Type = dbType

	// Handle exit
	exit := make(chan bool)
	go handleExit(exit)

	mux := http.NewServeMux()

	baseTemplate := template.Must(template.ParseFiles(currentPath+"/templates/base.html", currentPath+"/templates/index.html"))
	uploadTemplate := template.Must(template.ParseFiles(currentPath+"/templates/base.html", currentPath+"/templates/upload.html"))

	moviesTemplate := template.Must(template.ParseFiles(currentPath+"/templates/base.html", currentPath+"/templates/movies.html"))
	movieTemplate := template.Must(template.ParseFiles(currentPath+"/templates/base.html", currentPath+"/templates/movie.html"))
	seriesTemplate := template.Must(template.ParseFiles(currentPath+"/templates/base.html", currentPath+"/templates/series.html"))
	ptwTemplate := template.Must(template.ParseFiles(currentPath+"/templates/base.html", currentPath+"/templates/ptw.html"))

	fs := http.FileServer(http.Dir("assets"))

	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))
	mux.HandleFunc("/", internal.IndexHandle(baseTemplate))
	mux.HandleFunc("/upload", internal.UploadHandle(uploadTemplate, db))

	mux.HandleFunc("/movies", internal.MoviesHandle(moviesTemplate, db))
	mux.HandleFunc("/movie", internal.MoviesHandle(movieTemplate, db))
	mux.HandleFunc("/series", internal.SeriesHandle(seriesTemplate))
	mux.HandleFunc("/ptw", internal.PtwHandle(ptwTemplate))

	mux.HandleFunc("/api/ptw", internal.PtwApiHandle(db))

	// HTTP Server
	go func() {
		switch {
		// DEBUG Mode
		case debugFlag:
			err := http.ListenAndServe(":8080", requestLogger(mux))
			if err != nil && err != http.ErrServerClosed {
				panic("Error trying to start http server: " + err.Error())
			}

		case !debugFlag:
			err := http.ListenAndServe(":8080", mux)
			if err != nil && err != http.ErrServerClosed {
				panic("Error trying to start http server: " + err.Error())
			}
		}
	}()

	log.Println("Serving at 127.0.0.1:8080")
	<-exit

	return nil
}

// Function that handles the errors
func run() error {
	// host=192.168.30.171 port=5432 user=root password=toor dbname=waiw sslmode=disable
	var gokrazyFlag = flag.Bool("gokrazy", false, "use this if you are using gokrazy (note that this works with postgresql not sqlite)")
	var userFlag = flag.String("user", "", "-user='root'")
	var passwordFlag = flag.String("password", "", "-password='12345'")
	var ipFlag = flag.String("ip", "", "-ip='<ip where is postgresql database>'")
	var dbNameFlag = flag.String("dbname", "", "-dbname='<name of database to connect>'")

	var localFlag = flag.String("db", "", "-db='/prod/database.db' (note that this works with sqlite3 not postgresql)")

	var debugFlag = flag.String("debug", "", "-debug='/dev/database.db' to enter in debug mode with dev database (note that this works with sqlite3 not postgresql)")

	flag.Parse()

	var dbString string
	var dbType string

	isDebug := false

	currentPath, err := os.Getwd()
	if err != nil {
		return err
	}

	// if its gokrazy (postgresql)
	if *gokrazyFlag {
		// TODO: It would be better to use move instead of copy
		// copy required folders to /pem
		errAssets := cp.Copy("/etc/waiw/assets", "/perm/home/waiw/assets")
		errTemplates := cp.Copy("/etc/waiw/templates", "/perm/home/waiw/templates")
		if errAssets != nil || errTemplates != nil {
			return err
		}

		dbString = "host=" + *ipFlag + " port=5432 user=" + *userFlag + " password=" + *passwordFlag + " dbname=" + *dbNameFlag + " sslmode=disable"
		dbType = "postgres"

		// if its a local db file (sqlite)
	} else if len(*localFlag) > 0 {
		if _, err := os.Stat(*localFlag); err != nil {
			return fmt.Errorf("Database file does not exist\n")
		}

		dbType = "sqlite3"
		dbString = *localFlag

		// if its debug (sqlite)
	} else {
		if _, err := os.Stat(*debugFlag); err != nil {
			return fmt.Errorf("Database file does not exist\n")
		}

		dbType = "sqlite"
		dbString = *debugFlag
		isDebug = true
	}

	err = startServer(currentPath, dbString, isDebug, dbType)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
