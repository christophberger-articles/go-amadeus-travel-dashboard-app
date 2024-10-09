package main

import (
	_ "embed"
	"net/http"
)

// HomeHandler renders the initial search page from home.html
func (a *app) CitySearchHandler(w http.ResponseWriter, r *http.Request) {

	//cityname := r.FormValue("city")

	// call the Amadeus City Search API with cityname as input.

	//result :=

	// Set the content type to HTML
	w.Header().Set("Content-Type", "text/html")
	defer r.Body.Close()

	_, err := w.Write([]byte(""))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
