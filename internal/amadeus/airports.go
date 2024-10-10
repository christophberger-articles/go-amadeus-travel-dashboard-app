package amadeus

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"appliedgo.net/what"
)

type Airport struct {
	Name        string
	IATACode    string
	Performance string
	Probability string
}

func (c *Client) Airports(latitude, longitude string) (*Airport, error) {
	apiURL := fmt.Sprintf("%s/reference-data/locations/airports?latitude=%s&longitude=%s&sort=distance&page[limit]=1",
		c.baseURL, latitude, longitude)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var response struct {
		Data []struct {
			Name     string `json:"name"`
			IATACode string `json:"iataCode"`
		} `json:"data"`
		Errors []struct {
			Status int    `json:"status"`
			Code   int    `json:"code"`
			Title  string `json:"title"`
			Detail string `json:"detail"`
		} `json:"errors"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error unmarshaling airport response: %v", err)
	}

	what.Happens("response: %v", response)

	if len(response.Data) == 0 {
		return nil, fmt.Errorf("no airports found")
	}

	airport := response.Data[0]

	// 2. Call On-Time Performance API
	today := time.Now().Format("2006-01-02")
	performanceURL := fmt.Sprintf("%s/airport/predictions/on-time?airportCode=%s&date=%s",
		c.baseURL, airport.IATACode, today)

	performanceReq, err := http.NewRequest("GET", performanceURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating performance request: %v", err)
	}

	performanceResp, err := c.doRequest(performanceReq)
	if err != nil {
		return nil, fmt.Errorf("error making performance request: %v", err)
	}
	defer performanceResp.Body.Close()

	performanceBody, err := io.ReadAll(performanceResp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading performance response body: %v", err)
	}

	var performanceResponse struct {
		Data struct {
			Result      string `json:"result"`
			Probability string `json:"probability"`
		} `json:"data"`
	}

	if err := json.Unmarshal(performanceBody, &performanceResponse); err != nil {
		return nil, fmt.Errorf("error unmarshaling performance response: %v", err)
	}

	return &Airport{
		Name:        airport.Name,
		IATACode:    airport.IATACode,
		Performance: performanceResponse.Data.Result,
		Probability: performanceResponse.Data.Probability,
	}, nil
}
