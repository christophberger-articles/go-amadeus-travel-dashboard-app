package main

import (
	"log"
	"net/http"
)

func startServer(a *app) {
	mux := http.NewServeMux()

	// The dasboard page is on the base URL.
	// If the request includes a city, the dashboard will
	// display the travel info for that city.
	mux.HandleFunc("/{$}", a.DashboardHandler)
	mux.HandleFunc("/search", a.CitySearchHandler)
	mux.HandleFunc("/travelinfo/{latitude}/{longitude}", a.TravelInfoHandler)

	// Start the server
	go func() {
		log.Println("Listening on http://localhost:8020")
		log.Fatal(http.ListenAndServe("localhost:8020", mux))
	}()
}
