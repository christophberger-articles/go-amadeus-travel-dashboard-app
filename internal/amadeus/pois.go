package amadeus

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"appliedgo.net/what"
)

// POI represents a Point of Interest
type POI struct {
	Name     string
	Category string
}

// Pois retrieves Points of Interest near a specified location. It takes
// latitude and longitude as string parameters representing the center
// point. The function makes an API request to fetch POIs within a 5km
// radius of the given coordinates. It returns a slice of POI structs
// containing the name and category of each point of interest, or an
// error if the request fails or the response cannot be processed.
func (c *Client) Pois(latitude, longitude string) ([]POI, error) {
	// Construct the API URL
	apiURL := fmt.Sprintf("%s/reference-data/locations/pois?latitude=%s&longitude=%s&radius=5",
		c.baseURL, latitude, longitude)

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
			Category struct {
				Name string `json:"name"`
			} `json:"category"`
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

	what.Happens("response: %v", response)

	// Create and populate the POIs slice
	pois := make([]POI, len(response.Data))
	for i, item := range response.Data {
		pois[i] = POI{
			Name:     item.Name,
			Category: item.Category.Name,
		}
	}

	return pois, nil
}
