package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/BrunoTeixeira1996/waiw/utils"
)

// Handles the exit signal
func handleExit(exit chan bool) {
	ch := make(chan os.Signal, 5)
	signal.Notify(ch, os.Interrupt)
	<-ch
	log.Println("Closing web server")
	exit <- true
}

// Starts the web server
func startServer(currentPath string) error {

	// Handle exit
	exit := make(chan bool)
	go handleExit(exit)

	mux := http.NewServeMux()

	baseTemplate := template.Must(template.ParseFiles(currentPath+"/templates/base.html", currentPath+"/templates/index.html"))
	moviesTemplate := template.Must(template.ParseFiles(currentPath+"/templates/base.html", currentPath+"/templates/movies.html"))
	fs := http.FileServer(http.Dir("assets"))

	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))
	mux.HandleFunc("/", utils.IndexHandle(baseTemplate))
	mux.HandleFunc("/movies", utils.IndexHandle(moviesTemplate))

	httpServer := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: mux,
	}

	// HTTP Server
	go func() {
		err := httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			// If the server cannot start, just panic
			panic("Error trying to start http server: " + err.Error())
		}
	}()

	log.Println("Serving at 127.0.0.1:8080")
	<-exit

	return nil
}

// Function that handles the errors
func run() error {
	// TODO:
	// - Create tables for db
	// - Create struct to handle sqlite
	// - Get the sqlite db path
	// - Insert dummy data into the db
	// - Visualize data in website

	currentPath, err := os.Getwd()
	if err != nil {
		return err
	}

	err = startServer(currentPath)
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
