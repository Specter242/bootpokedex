package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	commands := getCommands()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")

		if !scanner.Scan() {
			fmt.Println("No more input to scan. Exiting REPL...")
			break
		}

		words := cleanInput(scanner.Text())
		if len(words) == 0 {
			continue
		}

		commandName := words[0]
		cmd, ok := commands[commandName]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		var arg string
		if len(words) > 1 && cmd.requiresArg {
			arg = words[1]
		}

		if cmd.requiresArg && arg == "" {
			fmt.Printf("Command '%s' requires an argument\n", commandName)
			continue
		}

		if !cmd.requiresArg && len(words) > 1 {
			fmt.Printf("Command '%s' doesn't accept arguments\n", commandName)
			continue
		}

		err := cmd.callback(arg)
		if err != nil {
			fmt.Println("Error:", err)
		}

		if commandName == "exit" {
			break
		}
	}
}

func cleanInput(text string) []string {
	// Remove any leading or trailing whitespace
	text = strings.TrimSpace(text)

	// Split the input text into words
	words := strings.Fields(text)

	// Create output slice with correct capacity
	cleanOutput := make([]string, len(words))

	// Lowercase the words and remove any commas
	for i, word := range words {
		word = strings.ToLower(word)
		word = strings.Trim(word, ",")
		cleanOutput[i] = word
	}

	return cleanOutput
}
