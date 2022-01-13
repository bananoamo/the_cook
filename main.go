package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type Application struct {
}

func main() {
	app := &Application{}
	siteMux := http.NewServeMux()

	// routes
	siteMux.HandleFunc("/", app.rootHandler)
	siteMux.HandleFunc("/start", app.startHandler)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      siteMux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	fmt.Printf("Starting server at %+v\n\n", srv.Addr)
	fmt.Printf("%s", welcom)
	fmt.Printf("\nTo stop server and quit press CTRL+C")
	log.Fatal(srv.ListenAndServe())
}
