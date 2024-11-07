package amadeus

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Hotels retrieves a list of hotels for a given city IATA code.
// It returns up to 10 hotel names with their distances from the city center.
// The function makes an API request to Amadeus' reference data endpoint.
// It handles JSON unmarshaling and formats the hotel information.
// If successful, it returns a slice of strings with hotel details.
// An error is returned if any step in the process fails.
func (c *Client) Hotels(iataCode string) ([]string, error) {
	// Construct the API URL
	apiURL := fmt.Sprintf("%s/reference-data/locations/hotels/by-city?cityCode=%s",
		c.baseURL, iataCode)

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

	// Define a struct to unmarshal the JSON response
	var response struct {
		Data []struct {
			Name     string `json:"name"`
			Distance struct {
				Value float64 `json:"value"`
				Unit  string  `json:"unit"`
			} `json:"distance"`
		} `json:"data"`
		Errors []struct {
			Status int    `json:"status"`
			Code   int    `json:"code"`
			Title  string `json:"title"`
			Detail string `json:"detail"`
		} `json:"errors"`
	}

	// Unmarshal the JSON response
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	// Format the hotels
	var hotels []string
	for _, hotel := range response.Data {
		if len(hotel.Distance.Unit) == 0 {
			hotel.Distance.Unit = "km"
		}
		hotelInfo := fmt.Sprintf("%s", cases.Title(language.English).String(hotel.Name))
		hotels = append(hotels, hotelInfo)
		if len(hotels) == 10 {
			break
		}
	}

	return hotels, nil
}
