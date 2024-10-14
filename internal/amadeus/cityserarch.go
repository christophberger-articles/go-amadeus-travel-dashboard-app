package amadeus

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// City represents a city with its geographical and identification information
type City struct {
	Name      string
	State     string
	IATACode  string
	Latitude  float64
	Longitude float64
}

// amadeusResponse represents the structure of the JSON response from the Amadeus API
// for city search queries.
type amadeusResponse struct {
	Data []struct {
		Name     string `json:"name"`
		Subtype  string `json:"subtype"`
		IATACode string `json:"iataCode"`
		Address  struct {
			StateCode string `json:"stateCode"`
		} `json:"address"`
		GeoCode struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"geoCode"`
	} `json:"data"`
}

// SearchCity queries the Amadeus API for cities matching the given cityName.
// It returns a slice of City structs containing information about the matched cities,
// including name, state, IATA code, latitude, and longitude.
// An error is returned if there's any issue with the API request or response processing.
func (c *Client) SearchCity(cityName string) ([]City, error) {
	// Construct the API URL
	apiURL := fmt.Sprintf("%s/reference-data/locations/cities?keyword=%s",
		c.baseURL, url.QueryEscape(cityName))

	// Create a new request
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Perform the request with authentication
	resp, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// Unmarshal the JSON response
	var amadeusResp amadeusResponse
	if err := json.Unmarshal(body, &amadeusResp); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	// Create and populate the locations slice
	locations := make([]City, len(amadeusResp.Data))
	for i, item := range amadeusResp.Data {
		// State codes are of the form CC-SC, where CC is the country code and SC is the state code. If a country has no state code, the state code form is CC-ZZZ. For nicer display, remove the -ZZZ suffix.
		state, _ := strings.CutSuffix(item.Address.StateCode, "-ZZZ")
		locations[i] = City{
			Name:      item.Name,
			State:     state,
			IATACode:  item.IATACode,
			Latitude:  item.GeoCode.Latitude,
			Longitude: item.GeoCode.Longitude,
		}
	}

	return locations, nil
}
