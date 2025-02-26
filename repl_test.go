package main

import (
	"testing"

	"github.com/Specter242/bootpokedex/internal/pokeapi"
)

// MockClient implements the APIClient interface for testing
type MockClient struct {
	locations *pokeapi.LocationResponse
}

// Ensure MockClient implements APIClient
var _ pokeapi.APIClient = (*MockClient)(nil)

func (m *MockClient) GetLocations(directionFWD bool) (*pokeapi.LocationResponse, error) {
	return m.locations, nil
}

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "Hello, world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "MAP",
			expected: []string{"map"},
		},
		{
			input:    "",
			expected: []string{},
		},
		{
			input:    "  ",
			expected: []string{},
		},
		{
			input:    "hello,world",
			expected: []string{"hello", "world"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("Input %q: Expected %d words, but got %d", c.input, len(c.expected), len(actual))
			continue
		}

		for i := range actual {
			if actual[i] != c.expected[i] {
				t.Errorf("Input %q: Expected word %q at position %d, but got %q",
					c.input, c.expected[i], i, actual[i])
			}
		}
	}
}

func TestGetCommands(t *testing.T) {
	commands := getCommands()

	expectedCommands := []string{"help", "exit", "map", "mapb"}

	for _, cmd := range expectedCommands {
		if _, exists := commands[cmd]; !exists {
			t.Errorf("Expected command %q to exist", cmd)
		}
	}

	// Test command descriptions are not empty
	for name, cmd := range commands {
		if cmd.description == "" {
			t.Errorf("Command %q has empty description", name)
		}
		if cmd.callback == nil {
			t.Errorf("Command %q has nil callback", name)
		}
	}
}

func TestCommandCallbacks(t *testing.T) {
	// Setup mock client
	mockLocations := &pokeapi.LocationResponse{
		Count: 1,
		Results: []pokeapi.Location{
			{Name: "test-location", URL: "test-url"},
		},
	}

	mockClient := &MockClient{locations: mockLocations}

	// Store the original client and restore it after the test
	originalClient := pokeClient
	defer func() { pokeClient = originalClient }()

	// Use type assertion to ensure we're using the interface type
	pokeClient = pokeapi.APIClient(mockClient)

	commands := getCommands()

	// Test map command
	err := commands["map"].callback()
	if err != nil {
		t.Errorf("map command failed: %v", err)
	}

	// Test mapb command
	err = commands["mapb"].callback()
	if err != nil {
		t.Errorf("mapb command failed: %v", err)
	}

	// Test help command
	err = commands["help"].callback()
	if err != nil {
		t.Errorf("help command failed: %v", err)
	}

	// Note: we don't test exit command as it calls os.Exit
}

func TestCommandDescriptions(t *testing.T) {
	commands := getCommands()

	expectedDescriptions := map[string]string{
		"help": "Displays a help message",
		"exit": "Exit the Pokedex",
		"map":  "Display the next 20 locations",
		"mapb": "Display the previous 20 locations",
	}

	for cmd, expectedDesc := range expectedDescriptions {
		if command, exists := commands[cmd]; exists {
			if command.description != expectedDesc {
				t.Errorf("Command %q: expected description %q, got %q",
					cmd, expectedDesc, command.description)
			}
		} else {
			t.Errorf("Expected command %q not found", cmd)
		}
	}
}

func TestCommandExit(t *testing.T) {
	err := commandExit()
	if err != nil {
		t.Errorf("Expected nil error, but got %v", err)
	}
}

func TestCommandHelp(t *testing.T) {
	err := commandHelp()
	if err != nil {
		t.Errorf("Expected nil error, but got %v", err)
	}
}
