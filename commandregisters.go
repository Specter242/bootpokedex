// Description: This file contains the command registers for the CLI.
// The getCommands function returns a map of command names to command structs.
// Each command struct contains the name of the command, a description, and a callback function.

package main

import (
	"fmt"
	"os"

	"github.com/Specter242/bootpokedex/internal/pokeapi" // fixed casing
)

// Make pokeClient a package variable that can be modified for testing
var pokeClient pokeapi.APIClient = pokeapi.NewClient("https://pokeapi.co/api/v2")

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Stdout.Sync() // Flush stdout to ensure the message is displayed
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	commands := getCommands()
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMap() error {
	locations, err := pokeClient.GetLocations(true)
	if err != nil {
		return err
	}

	fmt.Println("Location areas:")
	for _, loc := range locations.Results {
		fmt.Printf("- %s\n", loc.Name)
	}

	return nil
}

func commandMapb() error {
	locations, err := pokeClient.GetLocations(false)
	if err != nil {
		return err
	}

	fmt.Println("Location areas:")
	for _, loc := range locations.Results {
		fmt.Printf("- %s\n", loc.Name)
	}

	return nil
}

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Display the next 20 locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Display the previous 20 locations",
			callback:    commandMapb,
		},
	}
}
