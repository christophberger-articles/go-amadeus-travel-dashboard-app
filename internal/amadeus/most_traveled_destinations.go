package amadeus

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"appliedgo.net/what"
)

// MostTraveledDestinations retrieves the top 5 most traveled destinations
// from the given IATA code for the current month. It returns a string
// containing the destinations and their traveler scores, or an error if
// the request fails. The returned string format is "DEST1: SCORE1,
// DEST2: SCORE2, ...". 
func (c *Client) MostTraveledDestinations(iataCode string) (string, error) {
	// Get current month in YYYY-MM format
	currentMonth := time.Now().Format("2006-01")

	// Construct the API URL
	apiURL := fmt.Sprintf("%s/travel/analytics/air-traffic/traveled?originCityCode=%s&period=%s&max=5",
		c.baseURL, iataCode, currentMonth)

	// Create a new request
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	// Perform the request with authentication
	resp, err := c.doRequest(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	// Define a struct to unmarshal the JSON response
	var response struct {
		Data []struct {
			Destination string `json:"destination"`
			Analytics   struct {
				Travelers struct {
					Score int `json:"score"`
				} `json:"travelers"`
			} `json:"analytics"`
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
		return "", fmt.Errorf("error unmarshaling response: %v", err)
	}

	what.Happens("response: %v", response)

	// Format the destinations
	var destinations []string
	for _, item := range response.Data {
		destinations = append(destinations, fmt.Sprintf("%s: %d", item.Destination, item.Analytics.Travelers.Score))
	}

	// Join the destinations into a single string
	result := strings.Join(destinations, ", ")

	return result, nil
}
