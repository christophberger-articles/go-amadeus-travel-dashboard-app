package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"net/http"
	"strings"
	"text/template"

	"appliedgo.net/what"
	"github.com/christophberger-articles/go-amadeus-travel-dashboard-app/internal/amadeus"
)

// HomeHandler renders the initial search page from home.html
func (a *app) CitySearchHandler(w http.ResponseWriter, r *http.Request) {
	locationlist := ""

	// call the Amadeus City Search API with cityname as input.
	locations, err := a.amadeusClient.SearchCity(r.URL.Query().Get("city"))
	what.Is(locations)
	if err != nil {
		locationlist = fmt.Errorf("No cities found: %v", err).Error()
	} else {
		locationlist, err = locationsToHTMLList(locations)
		if err != nil {
			locationlist = fmt.Errorf("Error generating HTML list: %v", err).Error()
		}
	}

	// Set the content type to HTML
	w.Header().Set("Content-Type", "text/html")
	defer r.Body.Close()

	_, err = w.Write([]byte(locationlist))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func locationsToHTMLList(locations []amadeus.Location) (string, error) {
	// Define the HTML template
	const tmpl = `
<ul id="cities">
{{range .}}
    <li><a href="/travelinfo/{{.Latitude}}/{{.Longitude}}#cityinfo">{{.Name}}</a></li>
{{end}}
</ul>`

	// Parse the template
	t, err := template.New("locationList").Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("error parsing template: %v", err)
	}

	// Execute the template with the locations data
	var buf bytes.Buffer
	if err := t.Execute(&buf, locations); err != nil {
		return "", fmt.Errorf("error executing template: %v", err)
	}

	// Trim any leading or trailing whitespace
	return strings.TrimSpace(buf.String()), nil
}
