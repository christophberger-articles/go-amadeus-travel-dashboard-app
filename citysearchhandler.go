package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"net/http"
	"strings"
	"text/template"

	"github.com/christophberger-articles/go-amadeus-travel-dashboard-app/internal/amadeus"
)

// HomeHandler renders the initial search page from home.html
func (a *app) CitySearchHandler(w http.ResponseWriter, r *http.Request) {
	citylist := ""

	// call the Amadeus City Search API with cityname as input.
	cities, err := a.amadeusClient.SearchCity(r.URL.Query().Get("city"))
	if err != nil {
		citylist = fmt.Errorf("No cities found: %v", err).Error()
	} else {
		citylist, err = citiesToHTMLList(cities)
		if err != nil {
			citylist = fmt.Errorf("Error generating HTML list: %v", err).Error()
		}
	}

	// Set the content type to HTML
	w.Header().Set("Content-Type", "text/html")
	defer r.Body.Close()

	_, err = w.Write([]byte(citylist))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func citiesToHTMLList(locations []amadeus.City) (string, error) {
	// Define the HTML template
	const tmpl = `
<ul id="cities">
{{range .}}
    <li><a href="/travelinfo?name={{.Name | urlquery}}&iata={{.IATACode}}&lat={{.Latitude}}&lon={{.Longitude}}#cityinfo" target=htmz>{{.Name}}{{- if .State}}, {{.State}}{{end}}</a></li>
{{end}}
</ul>`

	// Parse the template
	t, err := template.New("cityList").Parse(tmpl)
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
