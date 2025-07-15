package main

import (
	"fmt"
	"strings"
	"bufio"
	"os"
)

func main() {
	prompt := "Pokedex > "
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(prompt)
		scanner.Scan()
		input := scanner.Text()
		cleanedInput := cleanInput(input)
		fmt.Printf("Your command was: %s\n", strings.ToLower(cleanedInput[0]))
	}

}

func cleanInput(text string) []string {
	words := strings.Fields(text)
	return words
}

