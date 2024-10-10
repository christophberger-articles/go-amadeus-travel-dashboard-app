package main

import (
	_ "embed"
	"net/http"
)

// HomeHandler renders the initial search page from home.html
func (a *app) TravelInfoHandler(w http.ResponseWriter, r *http.Request) {

	//city := r.FormValue("citycode")


	// Set the content type to HTML
	w.Header().Set("Content-Type", "text/html")
	defer r.Body.Close()

	_, err := w.Write([]byte(""))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
