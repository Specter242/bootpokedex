package main

import (
	"fmt"
	"testing"

	"github.com/Specter242/bootpokedex/internal/pokeapi"
)

// MockClient implements the APIClient interface for testing
type MockClient struct {
	locations   *pokeapi.LocationResponse
	pokeList    *pokeapi.PokeList
	shouldError bool
	callHistory []bool // tracks forward/backward calls
}

func (m *MockClient) GetLocations(directionFWD bool) (*pokeapi.LocationResponse, error) {
	if m.shouldError {
		return nil, fmt.Errorf("mock error")
	}
	m.callHistory = append(m.callHistory, directionFWD)
	return m.locations, nil
}

func (m *MockClient) Explore(locationName string) (*pokeapi.PokeList, error) {
	if m.shouldError {
		return nil, fmt.Errorf("mock error")
	}
	if locationName == "" {
		return nil, fmt.Errorf("location name required")
	}
	return m.pokeList, nil
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
			input:    "Hello World  ",
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

func TestCommandMap(t *testing.T) {
	mockClient := &MockClient{
		locations: &pokeapi.LocationResponse{
			Count: 1,
			Results: []pokeapi.Location{
				{Name: "test-location", URL: "test-url"},
			},
		},
		pokeList: &pokeapi.PokeList{
			Name: "test-pokemon",
			URL:  "test-url",
		},
	}

	originalClient := pokeClient
	pokeClient = mockClient
	defer func() { pokeClient = originalClient }()

	err := commandMap("")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !mockClient.callHistory[0] {
		t.Error("Expected forward direction call")
	}
}

func TestCommandMapb(t *testing.T) {
	mockClient := &MockClient{
		locations: &pokeapi.LocationResponse{
			Count: 1,
			Results: []pokeapi.Location{
				{Name: "test-location", URL: "test-url"},
			},
		},
	}

	originalClient := pokeClient
	pokeClient = mockClient
	defer func() { pokeClient = originalClient }()

	err := commandMapb("")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if mockClient.callHistory[0] {
		t.Error("Expected backward direction call")
	}
}

func TestCommandHelp(t *testing.T) {
	err := commandHelp("")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
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

func TestCommandExplore(t *testing.T) {
	tests := []struct {
		name    string
		arg     string
		wantErr bool
	}{
		{
			name:    "With valid location",
			arg:     "test-location",
			wantErr: false,
		},
		{
			name:    "With empty location",
			arg:     "",
			wantErr: true,
		},
	}

	mockClient := &MockClient{
		pokeList: &pokeapi.PokeList{
			Name: "test-location",
			PokemonEncounters: []pokeapi.PokemonEncounter{
				{
					Pokemon: pokeapi.Pokemon{
						Name: "pikachu",
						URL:  "https://pokeapi.co/api/v2/pokemon/25/",
					},
				},
			},
		},
	}

	// Replace the global client with our mock
	pokeClient = mockClient

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := commandExplore(tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("commandExplore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCommandRequirements(t *testing.T) {
	commands := getCommands()

	expectedArgs := map[string]bool{
		"help":    false,
		"exit":    false,
		"map":     false,
		"mapb":    false,
		"explore": true,
	}

	for cmd, requiresArg := range expectedArgs {
		if command, exists := commands[cmd]; exists {
			if command.requiresArg != requiresArg {
				t.Errorf("Command %q: expected requiresArg=%v, got %v",
					cmd, requiresArg, command.requiresArg)
			}
		} else {
			t.Errorf("Expected command %q not found", cmd)
		}
	}
}
