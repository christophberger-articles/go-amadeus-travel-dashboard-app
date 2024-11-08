package main

import (
	"bytes"
	_ "embed"
	"net/http"
	"text/template"

	"github.com/christophberger-articles/go-amadeus-travel-dashboard-app/internal/amadeus"
)

//go:embed cityinfo.gotmpl
var tmpl string

// TravelInfoHandler renders the HTML snippet containing the city's travel information.
func (a *app) TravelInfoHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	query := r.URL.Query()
	cityname := query.Get("name")
	citycode := query.Get("iata")
	latitude := query.Get("lat")
	longitude := query.Get("lon")

	// Create a template
	t := template.Must(template.New("cityinfo").Parse(tmpl))

	// Prepare data for the template
	data := struct {
		CityName           string
		IATACode           string
		Latitude           string
		Longitude          string
		Airport            amadeus.Airport
		AirportError       error
		BusiestPeriod      string
		BusiestPeriodError error
		MostTraveled       string
		MostTraveledError  error
		Hotels             []string
		HotelsError        error
	}{}

	// Fetch data and populate the struct
	data.CityName = cityname
	data.IATACode = citycode
	data.Latitude = latitude
	data.Longitude = longitude
	data.Airport, data.AirportError = a.amadeusClient.Airports(latitude, longitude)
	if len(data.IATACode) == 0 {
		citycode = data.Airport.IATACode
	}
	data.BusiestPeriod, data.BusiestPeriodError = a.amadeusClient.BusiestPeriod(citycode)
	data.MostTraveled, data.MostTraveledError = a.amadeusClient.MostTraveledDestinations(citycode)
	data.Hotels, data.HotelsError = a.amadeusClient.Hotels(citycode)

	// Execute the template
	var buf bytes.Buffer
	err := t.Execute(&buf, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if buf.Len() == 0 {
		buf.WriteString("<div id=\"cityinfo\">No data available</div>")
	}

	w.Header().Set("Content-Type", "text/html")

	_, err = w.Write([]byte(buf.Bytes()))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
