package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Specter242/bootpokedex/internal/pokecache"
)

const cacheInterval = 30 * time.Second

// APIClient interface defines the methods that need to be implemented
type APIClient interface {
	GetLocations(directionFWD bool) (*LocationResponse, error)
	Explore(locationName string) (*PokeList, error)
}

// Client is a PokeAPI client that handles API requests.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	cache      *pokecache.Cache
}

// Ensure Client implements APIClient
var _ APIClient = (*Client)(nil)

// NewClient creates a new instance of the PokeAPI client.
func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		cache: pokecache.NewCache(cacheInterval),
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

type PokeList struct {
	Name              string             `json:"name"`
	URL               string             `json:"url"`
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

type PokemonEncounter struct {
	Pokemon Pokemon `json:"pokemon"`
}

type Pokemon struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

var (
	CurrentLocationURL  string = "https://pokeapi.co/api/v2/location-area"
	PreviousLocationURL string = ""
	NextLocationURL     string = ""
)

// GetLocations fetches a list of locations.
func (c *Client) GetLocations(directionFWD bool) (*LocationResponse, error) {
	var url string
	if directionFWD {
		url = NextLocationURL
		if url == "" {
			url = c.BaseURL + "/location-area"
		}
	} else {
		url = PreviousLocationURL
		if url == "" {
			url = c.BaseURL + "/location-area"
		}
	}

	// Try to get from cache first
	if cachedData, exists := c.cache.Get(url); exists {
		// fmt.Printf("Cache hit for URL: %s\n", url)
		var locationResp LocationResponse
		if err := json.Unmarshal(cachedData, &locationResp); err != nil {
			return nil, fmt.Errorf("error decoding cached data: %w", err)
		}
		return &locationResp, nil
	}

	// fmt.Printf("Cache miss for URL: %s\n", url)
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
	CurrentLocationURL = url
	NextLocationURL = locationResp.Next
	PreviousLocationURL = locationResp.Previous

	// Serialize to JSON before storing in cache
	jsonData, err := json.Marshal(&locationResp)
	if err != nil {
		return nil, fmt.Errorf("error serializing location data for cache: %w", err)
	}
	c.cache.Add(url, jsonData)

	return &locationResp, nil
}

func (c *Client) Explore(locationName string) (*PokeList, error) {
	var url string
	if locationName == "" {
		url = CurrentLocationURL
	} else {
		url = c.BaseURL + "/location-area/" + locationName
	}

	// Try to get from cache first
	if cachedData, exists := c.cache.Get(url); exists {
		// fmt.Printf("Cache hit for URL: %s\n", url)
		var pokeList PokeList
		if err := json.Unmarshal(cachedData, &pokeList); err != nil {
			return nil, fmt.Errorf("error decoding cached data: %w", err)
		}
		return &pokeList, nil
	}

	// fmt.Printf("Cache miss for URL: %s\n", url)
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching locations: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %s: %s", resp.Status, body)
	}

	var pokeList PokeList
	if err := json.NewDecoder(resp.Body).Decode(&pokeList); err != nil {
		return nil, fmt.Errorf("error decoding location data: %w", err)
	}

	// Serialize to JSON before storing in cache
	jsonData, err := json.Marshal(&pokeList)
	if err != nil {
		return nil, fmt.Errorf("error serializing location data for cache: %w", err)
	}
	c.cache.Add(url, jsonData)

	return &pokeList, nil
}
