package amadeus

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"
)

// BusiestPeriod returns the busiest period(s) for the given IATA code.
// It fetches air traffic data for the previous year from the Amadeus API,
// analyzes the data, and returns a comma-separated string of the busiest
// period(s). The function returns an error if there are any issues with
// the API request or data processing.
func (c *Client) BusiestPeriod(iataCode string) (string, error) {
	// Calculate the previous year
	previousYear := time.Now().AddDate(-1, 0, 0).Format("2006")

	// Construct the API URL
	apiURL := fmt.Sprintf("%s/travel/analytics/air-traffic/busiest-period?cityCode=%s&period=%s",
		c.baseURL, iataCode, previousYear)

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
			Period string `json:"period"`
			Score  int    `json:"score"`
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

	// Sort periods by score in descending order
	sort.Slice(response.Data, func(i, j int) bool {
		return response.Data[i].Score > response.Data[j].Score
	})

	// Extract periods with the highest score
	var busiestPeriods []string
	highestScore := 0
	for _, period := range response.Data {
		if len(busiestPeriods) == 0 {
			highestScore = period.Score
			busiestPeriods = append(busiestPeriods, period.Period)
		} else if period.Score == highestScore {
			busiestPeriods = append(busiestPeriods, period.Period)
		} else {
			break
		}
	}

	// Join the busiest periods into a comma-separated string
	result := strings.Join(busiestPeriods, ", ")

	return result, nil
}
