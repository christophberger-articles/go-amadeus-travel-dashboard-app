package amadeus

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) Hotels(iataCode string) ([]string, error) {
	// Construct the API URL
	apiURL := fmt.Sprintf("%s/reference-data/locations/hotels/by-city?cityCode=%s&radius=5&radiusUnit=KM&hotelSource=ALL&max=10",
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
	}

	// Unmarshal the JSON response
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	// Format the hotels
	var hotels []string
	for _, hotel := range response.Data {
		hotelInfo := fmt.Sprintf("%s, %.2f %s", hotel.Name, hotel.Distance.Value, hotel.Distance.Unit)
		hotels = append(hotels, hotelInfo)
		if len(hotels) == 10 {
			break
		}
	}

	return hotels, nil
}
