// Description: This file contains the command registers for the CLI.
// The getCommands function returns a map of command names to command structs.
// Each command struct contains the name of the command, a description, and a callback function.

package main

import (
	"fmt"
	"os"

	"github.com/Specter242/bootpokedex/internal/pokeapi"
)

// Make pokeClient a package variable that can be modified for testing
var pokeClient pokeapi.APIClient = pokeapi.NewClient("https://pokeapi.co/api/v2")

func commandExit(arg string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Stdout.Sync() // Flush stdout to ensure the message is displayed
	os.Exit(0)
	return nil
}

func commandHelp(arg string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	commands := getCommands()
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMap(arg string) error {
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

func commandMapb(arg string) error {
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

func commandExplore(arg string) error {
	if arg == "" {
		return fmt.Errorf("missing location name")
	}

	pokeList, err := pokeClient.Explore(arg)
	if err != nil {
		return err
	}

	fmt.Printf("Exploring %s...\n", pokeList.Name)
	fmt.Println("Found Pokemon:")
	for _, encounter := range pokeList.PokemonEncounters {
		fmt.Printf("- %s\n", encounter.Pokemon.Name)
	}

	return nil
}

func commandCatch(arg string) error {
	if arg == "" {
		return fmt.Errorf("missing Pokemon name")
	}

	caughtPokemon, err := pokeClient.Catch(arg)
	if err != nil {
		return err
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", arg)

	if caughtPokemon {
		fmt.Printf("%s was caught!\n", arg)
	} else {
		fmt.Printf("%s escaped!\n", arg)
	}
	return nil
}

func commandInspect(arg string) error {
	if arg == "" {
		return fmt.Errorf("missing Pokemon name")
	}

	pokemon, err := pokeClient.InspectPokemon(arg)
	if err != nil {
		return err
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Printf("Stats:\n")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  - %s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Printf("Types:\n")
	for _, t := range pokemon.Types {
		fmt.Printf("  - %s\n", t.Type.Name)
	}

	return nil
}

func commandPokedex(arg string) error {
	pokedex, err := pokeClient.GetPokedex()
	if err != nil {
		return err
	}

	fmt.Println("Pokedex:")
	for _, pokemon := range pokedex.Pokemon {
		fmt.Printf("- %s\n", pokemon.Name)
	}

	return nil
}

type cliCommand struct {
	name        string
	description string
	callback    func(arg string) error
	requiresArg bool
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
			requiresArg: false,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
			requiresArg: false,
		},
		"map": {
			name:        "map",
			description: "Display the next 20 locations",
			callback:    commandMap,
			requiresArg: false,
		},
		"mapb": {
			name:        "mapb",
			description: "Display the previous 20 locations",
			callback:    commandMapb,
			requiresArg: false,
		},
		"explore": {
			name:        "explore",
			description: "Explore a specific location. Usage: explore <location_name>",
			callback:    commandExplore,
			requiresArg: true,
		},
		"catch": {
			name:        "catch",
			description: "Catch a Pokemon. Usage: catch <pokemon_name>",
			callback:    commandCatch,
			requiresArg: true,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a Pokemon. Usage: inspect <pokemon_name>",
			callback:    commandInspect,
			requiresArg: true,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Display the Pokedex",
			callback:    commandPokedex,
			requiresArg: false,
		},
	}
}
