package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/BrunoTeixeira1996/waiw/models"
	"github.com/BrunoTeixeira1996/waiw/utils"
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
func startServer(currentPath string, databasePath string, debugFlag bool) error {

	db := &models.Db{}

	db.Location = databasePath

	// Check if its in debug mode
	// FIXME: need to create somethign so we know it's in debug mode
	// if debugFlag {
	//
	// } else {
	//
	// }

	// Handle exit
	exit := make(chan bool)
	go handleExit(exit)

	mux := http.NewServeMux()

	baseTemplate := template.Must(template.ParseFiles(currentPath+"/templates/base.html", currentPath+"/templates/index.html"))
	uploadTemplate := template.Must(template.ParseFiles(currentPath+"/templates/base.html", currentPath+"/templates/upload.html"))

	moviesTemplate := template.Must(template.ParseFiles(currentPath+"/templates/base.html", currentPath+"/templates/movies.html"))
	movieTemplate := template.Must(template.ParseFiles(currentPath+"/templates/base.html", currentPath+"/templates/movie.html"))
	seriesTemplate := template.Must(template.ParseFiles(currentPath+"/templates/base.html", currentPath+"/templates/series.html"))

	fs := http.FileServer(http.Dir("assets"))

	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))
	mux.HandleFunc("/", utils.IndexHandle(baseTemplate))
	mux.HandleFunc("/upload", utils.UploadHandle(uploadTemplate, db))

	mux.HandleFunc("/movies", utils.MoviesHandle(moviesTemplate, db))
	mux.HandleFunc("/movie", utils.MoviesHandle(movieTemplate, db))
	mux.HandleFunc("/series", utils.SeriesHandle(seriesTemplate))

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
	debugFlag := false

	checkArgs := func() error {
		if len(os.Args) < 3 {
			return fmt.Errorf("Wrong nº of args, use ./waiw -db '<path>'\n")
		}
		if os.Args[1] != "-db" {
			return fmt.Errorf("Please provide the database full path using -db '<path>'\n")
		}

		if _, err := os.Stat(os.Args[2]); err != nil {
			return fmt.Errorf("Database file does not exist\n")
		}

		if len(os.Args) > 4 {
			return fmt.Errorf("Wrong nº of args, use ./waiw -db '<path>'\n")
		}

		if len(os.Args) == 4 && os.Args[3] == "-debug" {
			log.Println("DEBUG MODE")
			debugFlag = true
		}

		return nil
	}

	if err := checkArgs(); err != nil {
		return err
	}

	currentPath, err := os.Getwd()
	if err != nil {
		return err
	}

	err = startServer(currentPath, os.Args[2], debugFlag)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}
