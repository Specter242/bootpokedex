package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Specter242/bootpokedex/internal/pokecache"
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
var previousLocationURL string = ""
var nextLocationURL string = ""
var cache = pokecache.NewCache(5 * time.Second)

// GetLocations fetches a list of locations.
func (c *Client) GetLocations(directionFWD bool) (*LocationResponse, error) {
	var url string
	if directionFWD {
		url = nextLocationURL
		if url == "" {
			url = c.BaseURL + "/location-area"
		}
	} else {
		url = previousLocationURL
		if url == "" {
			url = c.BaseURL + "/location-area"
		}
	}

	// Try to get from cache first
	if cachedData, exists := cache.Get(url); exists {
		fmt.Println("Cache hit")
		var locationResp LocationResponse
		if err := json.Unmarshal(cachedData, &locationResp); err != nil {
			return nil, fmt.Errorf("error decoding cached data: %w", err)
		}
		return &locationResp, nil
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

	// Update URLs
	currentLocationURL = url
	nextLocationURL = locationResp.Next
	previousLocationURL = locationResp.Previous

	// Serialize to JSON before storing in cache
	jsonData, err := json.Marshal(&locationResp)
	if err != nil {
		return nil, fmt.Errorf("error serializing location data for cache: %w", err)
	}
	cache.Add(url, jsonData)

	return &locationResp, nil
}
