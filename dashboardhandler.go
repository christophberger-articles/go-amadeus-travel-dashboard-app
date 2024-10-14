package main

import (
	_ "embed"
	"net/http"
)

// embed the home page in the binary
//
//go:embed dashboard.html
var dashboardHTML []byte

// HomeHandler renders the initial search page from home.html
func (a *app) DashboardHandler(w http.ResponseWriter, r *http.Request) {
	// Set the content type to HTML
	w.Header().Set("Content-Type", "text/html")
	defer r.Body.Close()

	// Write the embedded HTML file content to the response writer
	_, err := w.Write(dashboardHTML)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
