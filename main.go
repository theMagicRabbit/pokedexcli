package main

import (
	"fmt"
	"strings"
	"bufio"
	"os"
)

var commands map[string]cliCommand

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	for k, v := range commands {
		fmt.Printf("%s: %s\n", k, v.description)
	}
	return nil
}

func cleanInput(text string) []string {
	words := strings.Fields(text)
	return words
}

func main() {
	prompt := "Pokedex > "
	scanner := bufio.NewScanner(os.Stdin)
	commands = map[string]cliCommand {
		"exit": {
			name:		"exit",
			description:	"Exit the program",
			callback:	commandExit,
		},
		"help": {
			name:		"help",
			description:	"Print help",
			callback:	commandHelp,
		},
	}
	for {
		fmt.Print(prompt)
		scanner.Scan()
		input := scanner.Text()
		cleanedInput := cleanInput(input)
		cmd, ok := commands[cleanedInput[0]]
		if ok {
			cmd.callback()
		} else {
			fmt.Println("Unknown command")
		}
	}

}
