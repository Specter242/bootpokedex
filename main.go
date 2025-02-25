package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	commands := getCommands()

	// Use a single Scanner to process input
	scanner := bufio.NewScanner(os.Stdin)
	for {
		// Prompt the user
		fmt.Print("Pokedex > ")

		// Check if input is available or EOF is encountered
		if !scanner.Scan() {
			fmt.Println("No more input to scan. Exiting REPL...")
			break
		}

		// Process the input (command)
		text := scanner.Text()
		// fmt.Println("Input received:", text)

		// Clean and parse the command
		words := cleanInput(text)
		if len(words) == 0 {
			fmt.Println("No command entered, continuing...")
			continue
		}

		command := words[0]
		// fmt.Println("Parsed command:", command)

		// Execute the command if it exists
		if cmd, ok := commands[command]; ok {
			// fmt.Println("Command found, executing:", command)
			err := cmd.callback()
			if err != nil {
				fmt.Println("Error:", err)
			}

			// Explicitly exit the loop if the command was `exit`
			if command == "exit" {
				break
			}
		} else {
			// Handle unknown commands
			fmt.Println("Unknown command")
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
