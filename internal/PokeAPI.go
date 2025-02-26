package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is a PokeAPI client that handles API requests.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a new instance of the PokeAPI client.
func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type LocationResponse struct {
	Count    int        `json:"count"`
	Next     string     `json:"next"`
	Previous string     `json:"previous"`
	Results  []Location `json:"results"`
}

type Location struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

var currentLocationURL string = "https://pokeapi.co/api/v2/location-area"

// GetLocations fetches a list of locations.
func (c *Client) GetLocations() (*LocationResponse, error) {
	url := currentLocationURL
	if url == "" {
		url = c.BaseURL + "/location-area"
	}

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching locations: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %s: %s", resp.Status, body)
	}

	var locationResp LocationResponse
	if err := json.NewDecoder(resp.Body).Decode(&locationResp); err != nil {
		return nil, fmt.Errorf("error decoding location data: %w", err)
	}

	currentLocationURL = locationResp.Next
	return &locationResp, nil
}
