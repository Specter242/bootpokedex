package pokeapi

import (
	"testing"
)

// MockClient implements APIClient interface for testing
type MockClient struct {
	GetLocationsFunc func(directionFWD bool) (*LocationResponse, error)
	ExploreFunc      func(locationName string) (*PokeList, error)
	CatchFunc        func(pokemonName string) (bool, error)
}

func (m *MockClient) GetLocations(directionFWD bool) (*LocationResponse, error) {
	return m.GetLocationsFunc(directionFWD)
}

func (m *MockClient) Explore(locationName string) (*PokeList, error) {
	return m.ExploreFunc(locationName)
}

func (m *MockClient) Catch(pokemonName string) (bool, error) {
	return m.CatchFunc(pokemonName)
}

func (m *MockClient) InspectPokemon(pokemonName string) (*Pokemon, error) {
	return nil, nil
}

func TestClientImplementsInterface(t *testing.T) {
	var _ APIClient = (*Client)(nil)
	var _ APIClient = (*MockClient)(nil)
}
