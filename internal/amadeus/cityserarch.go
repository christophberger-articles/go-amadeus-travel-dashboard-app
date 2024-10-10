package amadeus

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Location struct {
	Name      string
	IATACode  string
	Latitude  float64
	Longitude float64
}

type amadeusResponse struct {
	Data []struct {
		Name     string `json:"name"`
		IATACode string `json:"iataCode"`
		GeoCode  struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"geoCode"`
	} `json:"data"`
}

func (c *Client) SearchCity(cityName string) ([]Location, error) {
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
	locations := make([]Location, len(amadeusResp.Data))
	for i, item := range amadeusResp.Data {
		locations[i] = Location{
			Name:      item.Name,
			IATACode:  item.IATACode,
			Latitude:  item.GeoCode.Latitude,
			Longitude: item.GeoCode.Longitude,
		}
	}

	return locations, nil
}
