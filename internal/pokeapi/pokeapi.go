package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/Specter242/bootpokedex/internal/pokecache"
)

const cacheInterval = 30 * time.Second

// APIClient interface defines the methods that need to be implemented
type APIClient interface {
	GetLocations(directionFWD bool) (*LocationResponse, error)
	Explore(locationName string) (*PokeList, error)
	Catch(pokemonName string) (bool, error)
	InspectPokemon(pokemonName string) (*Pokemon, error)
	GetPokedex() (*Pokedex, error)
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
	Name           string     `json:"name"`
	URL            string     `json:"url"`
	BaseExperience int        `json:"base_experience"`
	Height         int        `json:"height"`
	Weight         int        `json:"weight"`
	Stats          []PokeStat `json:"stats"`
	Types          []PokeType `json:"types"`
}

type PokeStat struct {
	BaseStat int `json:"base_stat"`
	Stat     struct {
		Name string `json:"name"`
	} `json:"stat"`
}

type PokeType struct {
	Type struct {
		Name string `json:"name"`
	} `json:"type"`
}

type Pokedex struct {
	Pokemon []Pokemon `json:"pokemon"`
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

func (c *Client) Catch(pokemonName string) (bool, error) {
	url := c.BaseURL + "/pokemon/" + pokemonName

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return false, fmt.Errorf("error fetching pokemon: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("unexpected status %s: %s", resp.Status, body)
	}

	var pokemon Pokemon
	if err := json.NewDecoder(resp.Body).Decode(&pokemon); err != nil {
		return false, fmt.Errorf("error decoding pokemon data: %w", err)
	}

	// Calculate catch chance
	// Higher base experience means harder to catch
	catchRate := 0
	if pokemon.BaseExperience > 0 {
		catchRate = (100 - pokemon.BaseExperience/4)
		if catchRate < 10 {
			catchRate = 10 // Minimum 10% chance
		}
	} else {
		catchRate = 50 // Default 50% if no base experience
	}

	// Random roll (0-99)
	roll := rand.Intn(100)

	caught := roll < catchRate
	if caught {
		// Add to Pokedex
		jsonData, err := json.Marshal(&pokemon)
		if err != nil {
			return false, fmt.Errorf("error serializing pokemon data: %w", err)
		}
		c.cache.Add("pokedex/"+pokemonName, jsonData)
	}

	return caught, nil
}

func (c *Client) InspectPokemon(pokemonName string) (*Pokemon, error) {
	// Check if pokemon exists in pokedex (cache)
	if cachedData, exists := c.cache.Get("pokedex/" + pokemonName); exists {
		var pokemon Pokemon
		if err := json.Unmarshal(cachedData, &pokemon); err != nil {
			return nil, fmt.Errorf("error decoding cached pokemon data: %w", err)
		}
		return &pokemon, nil
	}

	return nil, fmt.Errorf("you haven't caught %s yet", pokemonName)
}

func (c *Client) GetPokedex() (*Pokedex, error) {
	var pokedex Pokedex
	const prefix = "pokedex/"

	// Get all pokemon keys from cache
	pokemonNames := c.cache.GetKeysWithPrefix(prefix)

	for _, name := range pokemonNames {
		if name == "" {
			continue
		}
		if pokemon, err := c.InspectPokemon(name); err == nil {
			pokedex.Pokemon = append(pokedex.Pokemon, *pokemon)
		}
	}

	return &pokedex, nil
}
