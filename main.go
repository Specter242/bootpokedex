package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	for {
		fmt.Print("Pokedex > ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		text := scanner.Text()
		words := cleanInput(text)
		firstWord := words[0]
		fmt.Printf("Your command was: %s", firstWord)
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
