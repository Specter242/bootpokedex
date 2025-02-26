package pokeapi

import (
	"testing"
)

// MockClient implements APIClient interface for testing
type MockClient struct {
	GetLocationsFunc func(directionFWD bool) (*LocationResponse, error)
	ExploreFunc      func(locationName string) (*PokeList, error)
}

func (m *MockClient) GetLocations(directionFWD bool) (*LocationResponse, error) {
	return m.GetLocationsFunc(directionFWD)
}

func (m *MockClient) Explore(locationName string) (*PokeList, error) {
	return m.ExploreFunc(locationName)
}

func TestClientImplementsInterface(t *testing.T) {
	var _ APIClient = (*Client)(nil)
	var _ APIClient = (*MockClient)(nil)
}
