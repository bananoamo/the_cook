package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type Application struct {
	key         string
	guessTries  int
	guessNumber int
	Ingredients *Ingredients
}

func main() {
	app := &Application{
		guessTries: 1,
	}
	siteMux := http.NewServeMux()

	// routes
	siteMux.HandleFunc("/", app.rootHandler)
	siteMux.HandleFunc("/start", app.startHandler)
	siteMux.HandleFunc("/cook_island", app.cookIslandHandler)
	siteMux.HandleFunc("/become_dedicated_chef", app.becomeChef)
	siteMux.HandleFunc("/requirements", app.requirements)
	siteMux.HandleFunc("/make_potion", app.makePotion)

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
