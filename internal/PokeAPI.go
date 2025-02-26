package pokeapi

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

// location represents a subset of data returned for a location.
type location struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	// Add additional fields as needed.
}

// GetLocations fetches a list of locations.
func (c *Client) GetLocations() ([]location, error) {
	endpoint := fmt.Sprintf("%s/location-area", c.BaseURL)
	resp, err := c.HTTPClient.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("error fetching locations: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %s: %s", resp.Status, body)
	}

	var locations []location
	if err := json.NewDecoder(resp.Body).Decode(&locations); err != nil {
		return nil, fmt.Errorf("error decoding location data: %w", err)
	}

	return locations, nil
}

// Ability represents a subset of data returned for an ability.
type Ability struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	// Add additional fields as needed.
}

// GetAbility fetches an ability's information by name.
func (c *Client) GetAbility(name string) (*Ability, error) {
	endpoint := fmt.Sprintf("%s/ability/%s", c.BaseURL, name)
	resp, err := c.HTTPClient.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("error fetching ability: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %s: %s", resp.Status, body)
	}

	var ability Ability
	if err := json.NewDecoder(resp.Body).Decode(&ability); err != nil {
		return nil, fmt.Errorf("error decoding ability data: %w", err)
	}

	return &ability, nil
}
