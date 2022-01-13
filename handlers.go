package main

import (
	"fmt"
	"net/http"
)

func (app *Application) rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "read instructions on the terminal where u started server\n")
}

func (app *Application) startHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Add("Allow", http.MethodGet)
		http.Error(w, "hint: look response status for hint", http.StatusMethodNotAllowed)
		return
	}
	app.writeResponse(w, http.StatusAccepted, startSuccess)
}

func (app *Application) writeResponse(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	fmt.Fprintf(w, "%s", message)
}
